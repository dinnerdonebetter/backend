package authentication

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/o1egl/paseto"

	"github.com/prixfixeco/api_server/internal/authentication"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	// ErrUserNotFound indicates a user was not located.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserBanned indicates a user is banned from using the service.
	ErrUserBanned = errors.New("user is banned")
	// ErrInvalidCredentials indicates a user provided invalid credentials.
	ErrInvalidCredentials = errors.New("invalid credentials")

	customCookieDomainHeader = "X-PRIXFIXE-COOKIE-DOMAIN"

	allowedCookiesHat    sync.Mutex
	allowedCookieDomains = map[string]uint{
		".prixfixe.local": 0,
		".prixfixe.dev":   1,
		".prixfixe.app":   2,
	}
)

// determineCookieDomain determines which domain to assign a cookie.
func (s *service) determineCookieDomain(ctx context.Context, req *http.Request) string {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	requestedCookieDomain := s.config.Cookies.Domain
	if headerCookieDomain := req.Header.Get(customCookieDomainHeader); headerCookieDomain != "" {
		allowedCookiesHat.Lock()
		// if the requested domain is present in the map, and it has a lower score than the current domain, then
		if currentScore, ok1 := allowedCookieDomains[requestedCookieDomain]; ok1 {
			if newScore, ok2 := allowedCookieDomains[headerCookieDomain]; ok2 {
				if currentScore > newScore {
					requestedCookieDomain = headerCookieDomain
				}
			}
		}
		allowedCookiesHat.Unlock()
	}

	return requestedCookieDomain
}

func (s *service) authenticateUserAndBuildCookie(ctx context.Context, loginData *types.UserLoginInput, requestedCookieDomain string) (*types.User, *http.Cookie, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.UsernameKey, loginData.Username)

	user, err := s.userDataManager.GetUserByUsername(ctx, loginData.Username)
	if err != nil || user == nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, observability.PrepareError(err, logger, span, "fetching user")
	}

	logger = logger.WithValue(keys.UserIDKey, user.ID)
	tracing.AttachUserToSpan(span, user)

	if user.IsBanned() {
		return user, nil, ErrUserBanned
	}

	loginValid, err := s.validateLogin(ctx, user, loginData)
	logger.WithValue("login_valid", loginValid)

	if err != nil {
		if errors.Is(err, authentication.ErrInvalidTOTPToken) {
			return user, nil, ErrInvalidCredentials
		} else if errors.Is(err, authentication.ErrPasswordDoesNotMatch) {
			return user, nil, ErrInvalidCredentials
		}

		logger.Error(err, "error encountered validating login")

		return user, nil, observability.PrepareError(err, logger, span, "validating login")
	} else if !loginValid {
		logger.Debug("login was invalid")
		return user, nil, ErrInvalidCredentials
	}

	defaultHouseholdID, err := s.householdMembershipManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
	if err != nil {
		return user, nil, observability.PrepareError(err, logger, span, "fetching user memberships")
	}

	cookie, err := s.issueSessionManagedCookie(ctx, defaultHouseholdID, user.ID, requestedCookieDomain)
	if err != nil {
		return user, nil, observability.PrepareError(err, logger, span, "issuing cookie")
	}

	return user, cookie, nil
}

// BeginSessionHandler is our login route.
func (s *service) BeginSessionHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	loginData := new(types.UserLoginInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, loginData); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := loginData.ValidateWithContext(ctx, s.config.MinimumUsernameLength, s.config.MinimumPasswordLength); err != nil {
		observability.AcknowledgeError(err, logger, span, "validating input")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	logger = logger.WithValue(keys.UsernameKey, loginData.Username)

	requestedCookieDomain := s.determineCookieDomain(ctx, req)
	if requestedCookieDomain != "" {
		logger = logger.WithValue("cookie_domain", requestedCookieDomain)
		logger.Debug("setting alternative cookie domain")
	}

	user, cookie, err := s.authenticateUserAndBuildCookie(ctx, loginData, requestedCookieDomain)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		case errors.Is(err, ErrUserBanned):
			s.encoderDecoder.EncodeErrorResponse(ctx, res, user.ReputationExplanation, http.StatusForbidden)
		case errors.Is(err, ErrInvalidCredentials):
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "login was invalid", http.StatusUnauthorized)
		default:
			observability.AcknowledgeError(err, logger, span, "issuing cookie")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
		}
		return
	}

	if err = s.customerDataCollector.EventOccurred(ctx, "logged_in", user.ID, map[string]interface{}{}); err != nil {
		logger.Error(err, "notifying customer data platform of login")
	}

	http.SetCookie(res, cookie)

	statusResponse := &types.UserStatusResponse{
		UserIsAuthenticated:       true,
		UserReputation:            user.ServiceHouseholdStatus,
		UserReputationExplanation: user.ReputationExplanation,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, statusResponse, http.StatusAccepted)
	logger.Debug("user logged in")
}

// ChangeActiveHouseholdHandler is our login route.
func (s *service) ChangeActiveHouseholdHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	input := new(types.ChangeActiveHouseholdInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	householdID := input.HouseholdID
	logger = logger.WithValue("new_session_household_id", householdID)

	requesterID := sessionCtxData.Requester.UserID
	logger = logger.WithValue("user_id", requesterID)

	authorizedForHousehold, err := s.householdMembershipManager.UserIsMemberOfHousehold(ctx, requesterID, householdID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "checking permissions")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
		return
	}

	if !authorizedForHousehold {
		logger.Debug("invalid household ID requested for activation")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	requestedCookieDomain := s.determineCookieDomain(ctx, req)
	if requestedCookieDomain != "" {
		logger = logger.WithValue("cookie_domain", requestedCookieDomain)
		logger.Debug("setting alternative cookie domain")
	}

	cookie, err := s.issueSessionManagedCookie(ctx, householdID, requesterID, requestedCookieDomain)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "issuing cookie")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
		return
	}

	if err = s.customerDataCollector.EventOccurred(ctx, "changed_active_household", requesterID, map[string]interface{}{
		"old_household_id":        sessionCtxData.ActiveHouseholdID,
		keys.ActiveHouseholdIDKey: householdID,
	}); err != nil {
		logger.Error(err, "notifying customer data platform of login")
	}

	logger.Info("successfully changed active session household")
	http.SetCookie(res, cookie)

	res.WriteHeader(http.StatusAccepted)
}

// EndSessionHandler is our logout route.
func (s *service) EndSessionHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	ctx, loadErr := s.sessionManager.Load(ctx, "")
	if loadErr != nil {
		// this can literally never happen in this version of scs, because the token is empty
		observability.AcknowledgeError(err, logger, span, "loading token")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if destroyErr := s.sessionManager.Destroy(ctx); destroyErr != nil {
		observability.AcknowledgeError(err, logger, span, "destroying session")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	requestedCookieDomain := s.determineCookieDomain(ctx, req)
	if requestedCookieDomain != "" {
		logger = logger.WithValue("cookie_domain", requestedCookieDomain)
		logger.Debug("setting alternative cookie domain")
	}

	newCookie, cookieBuildingErr := s.buildCookie(requestedCookieDomain, "deleted", time.Time{})
	if cookieBuildingErr != nil || newCookie == nil {
		observability.AcknowledgeError(err, logger, span, "building cookie")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	newCookie.MaxAge = -1
	http.SetCookie(res, newCookie)

	if err = s.customerDataCollector.EventOccurred(ctx, "logged_out", sessionCtxData.Requester.UserID, map[string]interface{}{}); err != nil {
		logger.Error(err, "notifying customer data platform of login")
	}

	logger.Debug("user logged out")

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// StatusHandler returns the user info for the user making the request.
func (s *service) StatusHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	var statusResponse *types.UserStatusResponse

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	statusResponse = &types.UserStatusResponse{
		ActiveHousehold:           sessionCtxData.ActiveHouseholdID,
		UserReputation:            sessionCtxData.Requester.Reputation,
		UserReputationExplanation: sessionCtxData.Requester.ReputationExplanation,
		UserIsAuthenticated:       true,
	}

	s.encoderDecoder.RespondWithData(ctx, res, statusResponse)
}

const (
	pasetoRequestTimeThreshold = 2 * time.Minute
)

// PASETOHandler returns the user info for the user making the request.
func (s *service) PASETOHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	input := new(types.PASETOCreationInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	requestedHousehold := input.HouseholdID
	logger = logger.WithValue(keys.APIClientClientIDKey, input.ClientID)

	if requestedHousehold != "" {
		logger = logger.WithValue("requested_household", requestedHousehold)
	}

	reqTime := time.Unix(0, input.RequestTime)
	if time.Until(reqTime) > pasetoRequestTimeThreshold || time.Since(reqTime) > pasetoRequestTimeThreshold {
		logger.WithValue("provided_request_time", reqTime.String()).Debug("PASETO request denied because its time is out of threshold")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	sum, err := base64.RawURLEncoding.DecodeString(req.Header.Get(signatureHeaderKey))
	if err != nil || len(sum) == 0 {
		logger.WithValue("sum_length", len(sum)).Error(err, "invalid signature")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	client, clientRetrievalErr := s.apiClientManager.GetAPIClientByClientID(ctx, input.ClientID)
	if clientRetrievalErr != nil {
		observability.AcknowledgeError(err, logger, span, "fetching API client")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	mac := hmac.New(sha256.New, client.ClientSecret)
	if _, macWriteErr := mac.Write(s.encoderDecoder.MustEncodeJSON(ctx, input)); macWriteErr != nil {
		// sha256.digest.Write does not ever return an error, so this branch will remain "uncovered" :(
		observability.AcknowledgeError(err, logger, span, "writing HMAC message for comparison")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	if !hmac.Equal(sum, mac.Sum(nil)) {
		logger.Info("invalid credentials passed to PASETO creation route")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	user, err := s.userDataManager.GetUser(ctx, client.BelongsToUser)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving user")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	logger = logger.WithValue(keys.UserIDKey, user.ID)

	sessionCtxData, err := s.householdMembershipManager.BuildSessionContextDataForUser(ctx, user.ID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving perms for API client")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	var requestedHouseholdID string

	if requestedHousehold != "" {
		if _, isMember := sessionCtxData.HouseholdPermissions[requestedHousehold]; !isMember {
			logger.Debug("invalid household ID requested for token")
			s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
			return
		}

		logger.WithValue("requested_household", requestedHousehold).Debug("setting token household ID to requested household")
		requestedHouseholdID = requestedHousehold
		sessionCtxData.ActiveHouseholdID = requestedHousehold
	} else {
		requestedHouseholdID = sessionCtxData.ActiveHouseholdID
	}

	logger = logger.WithValue(keys.HouseholdIDKey, requestedHouseholdID)

	// Encrypt data
	tokenRes, err := s.buildPASETOResponse(ctx, sessionCtxData, client)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "encrypting PASETO")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	logger.Info("PASETO issued")

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, tokenRes, http.StatusAccepted)
}

func (s *service) buildPASETOToken(ctx context.Context, sessionCtxData *types.SessionContextData, client *types.APIClient) paseto.JSONToken {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	now := time.Now().UTC()
	lifetime := time.Duration(math.Min(float64(maxPASETOLifetime), float64(s.config.PASETO.Lifetime)))
	expiry := now.Add(lifetime)

	jsonToken := paseto.JSONToken{
		Audience:   client.BelongsToUser,
		Subject:    client.BelongsToUser,
		Jti:        uuid.NewString(),
		Issuer:     s.config.PASETO.Issuer,
		IssuedAt:   now,
		NotBefore:  now,
		Expiration: expiry,
	}

	jsonToken.Set(pasetoDataKey, base64.RawURLEncoding.EncodeToString(sessionCtxData.ToBytes()))

	return jsonToken
}

func (s *service) buildPASETOResponse(ctx context.Context, sessionCtxData *types.SessionContextData, client *types.APIClient) (*types.PASETOResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	jsonToken := s.buildPASETOToken(ctx, sessionCtxData, client)

	// Encrypt data
	token, err := paseto.NewV2().Encrypt(s.config.PASETO.LocalModeKey, jsonToken, "")
	if err != nil {
		return nil, observability.PrepareError(err, s.logger, span, "encrypting PASETO")
	}

	tokenRes := &types.PASETOResponse{
		Token:     token,
		ExpiresAt: jsonToken.Expiration.String(),
	}

	return tokenRes, nil
}

// CycleCookieSecretHandler rotates the cookie building secret with a new random secret.
func (s *service) CycleCookieSecretHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Info("cycling cookie secret!")

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	if !sessionCtxData.Requester.ServicePermissions.CanCycleCookieSecrets() {
		logger.Debug("invalid permissions")
		s.encoderDecoder.EncodeInvalidPermissionsResponse(ctx, res)
		return
	}

	s.cookieManager = securecookie.New(
		securecookie.GenerateRandomKey(cookieSecretSize),
		[]byte(s.config.Cookies.BlockKey),
	)

	res.WriteHeader(http.StatusAccepted)
}
