package users

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/pquerna/otp/totp"
	"github.com/segmentio/ksuid"
	passwordvalidator "github.com/wagslane/go-password-validator"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// UserIDURIParamKey is used to refer to user IDs in router params.
	UserIDURIParamKey = "userID"

	totpIssuer             = "Prixfixe"
	base64ImagePrefix      = "data:image/jpeg;base64,"
	minimumPasswordEntropy = 60
	totpSecretSize         = 64
	passwordResetTokenSize = 32
)

// validateCredentialChangeRequest takes a user's credentials and determines
// if they match what is on record.
func (s *service) validateCredentialChangeRequest(ctx context.Context, userID, password, totpToken string) (user *types.User, httpStatus int) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.UserIDKey, userID)

	// fetch user data.
	user, err := s.userDataManager.GetUser(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, http.StatusNotFound
	} else if err != nil {
		logger.Error(err, "error encountered fetching user")
		return nil, http.StatusInternalServerError
	}

	// validate login.
	valid, validationErr := s.authenticator.ValidateLogin(ctx, user.HashedPassword, password, user.TwoFactorSecret, totpToken)
	if validationErr != nil {
		logger.WithValue("validation_error", validationErr).Debug("error validating credentials")
		return nil, http.StatusBadRequest
	} else if !valid {
		logger.WithValue("valid", valid).Error(err, "invalid credentials")
		return nil, http.StatusUnauthorized
	}

	return user, http.StatusOK
}

// UsernameSearchHandler is a handler for responding to username queries.
func (s *service) UsernameSearchHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	query := req.URL.Query().Get(types.SearchQueryKey)
	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// fetch user data.
	users, err := s.userDataManager.SearchForUsersByUsername(ctx, query)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for users")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode response.
	s.encoderDecoder.RespondWithData(ctx, res, users)
}

// ListHandler is a handler for responding with a list of users.
func (s *service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine desired filter.
	qf := types.ExtractQueryFilterFromRequest(req)

	// fetch user data.
	users, err := s.userDataManager.GetUsers(ctx, qf)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching users")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode response.
	s.encoderDecoder.RespondWithData(ctx, res, users)
}

// CreateHandler is our user creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// in the event that we don't want new users to be able to sign up (a config setting) just decline the request from the get-go
	if !s.authSettings.EnableUserSignup {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "user creation is disabled", http.StatusForbidden)
		return
	}

	// decode the request.
	registrationInput := new(types.UserRegistrationInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, registrationInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	registrationInput.Username = strings.TrimSpace(registrationInput.Username)
	tracing.AttachUsernameToSpan(span, registrationInput.Username)
	registrationInput.EmailAddress = strings.TrimSpace(strings.ToLower(registrationInput.EmailAddress))
	tracing.AttachEmailAddressToSpan(span, registrationInput.EmailAddress)
	registrationInput.Password = strings.TrimSpace(registrationInput.Password)

	logger = logger.WithValues(map[string]interface{}{
		keys.UsernameKey:                 registrationInput.Username,
		keys.UserEmailAddressKey:         registrationInput.EmailAddress,
		keys.HouseholdInvitationIDKey:    registrationInput.InvitationID,
		keys.HouseholdInvitationTokenKey: registrationInput.InvitationToken,
	})

	if err := registrationInput.ValidateWithContext(ctx, s.authSettings.MinimumUsernameLength, s.authSettings.MinimumPasswordLength); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	// ensure the password is not garbage-tier
	if err := passwordvalidator.Validate(registrationInput.Password, minimumPasswordEntropy); err != nil {
		logger.WithValue("password_validation_error", err).Debug("weak password provided to user creation route")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "password too weak", http.StatusBadRequest)
		return
	}

	var invitation *types.HouseholdInvitation
	if registrationInput.InvitationID != "" && registrationInput.InvitationToken != "" {
		i, err := s.householdInvitationDataManager.GetHouseholdInvitationByTokenAndID(ctx, registrationInput.InvitationToken, registrationInput.InvitationID)
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
			return
		} else if err != nil {
			observability.AcknowledgeError(err, logger, span, "retrieving invitation")
			s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
			return
		}

		invitation = i
		logger = logger.WithValue(keys.HouseholdInvitationIDKey, invitation.ID)
		logger.Debug("retrieved household invitation")
	}

	logger.Debug("completed invitation check")

	// hash the password
	hp, err := s.authenticator.HashPassword(ctx, registrationInput.Password)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// generate a two factor secret.
	tfs, err := s.secretGenerator.GenerateBase32EncodedString(ctx, totpSecretSize)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "internal error", http.StatusInternalServerError)
		return
	}

	input := &types.UserDatabaseCreationInput{
		ID:              ksuid.New().String(),
		Username:        registrationInput.Username,
		EmailAddress:    registrationInput.EmailAddress,
		HashedPassword:  hp,
		TwoFactorSecret: tfs,
		InvitationToken: registrationInput.InvitationToken,
		BirthDay:        registrationInput.BirthDay,
		BirthMonth:      registrationInput.BirthMonth,
	}

	if invitation != nil {
		logger.Debug("supplementing user creation input with invitation data")
		input.DestinationHouseholdID = invitation.DestinationHousehold.ID
		input.InvitationToken = invitation.Token
	}

	// create the user.
	user, userCreationErr := s.userDataManager.CreateUser(ctx, input)
	if userCreationErr != nil {
		if errors.Is(userCreationErr, database.ErrUserAlreadyExists) {
			observability.AcknowledgeError(err, logger, span, "creating user")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "username taken", http.StatusBadRequest)
			return
		}

		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	logger.Debug("user created")

	// notify the relevant parties.
	tracing.AttachUserIDToSpan(span, user.ID)
	s.userCounter.Increment(ctx)

	// UserCreationResponse is a struct we can use to notify the user of their two factor secret, but ideally just this once and then never again.
	ucr := &types.UserCreationResponse{
		CreatedUserID:   user.ID,
		Username:        user.Username,
		EmailAddress:    user.EmailAddress,
		CreatedAt:       user.CreatedAt,
		TwoFactorSecret: user.TwoFactorSecret,
		BirthDay:        user.BirthDay,
		BirthMonth:      user.BirthMonth,
		TwoFactorQRCode: s.buildQRCode(ctx, user.Username, user.TwoFactorSecret),
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.UserDataType,
			EventType:            types.UserSignedUpCustomerEventType,
			AttributableToUserID: user.ID,
		}

		if publishErr := s.dataChangesPublisher.Publish(ctx, dcm); publishErr != nil {
			observability.AcknowledgeError(publishErr, logger, span, "publishing data change message")
		}
	}

	// encode and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, ucr, http.StatusCreated)
}

// buildQRCode builds a QR code for a given username and secret.
func (s *service) buildQRCode(ctx context.Context, username, twoFactorSecret string) string {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.UsernameKey, username)

	// "otpauth://totp/{{ .Issuer }}:{{ .EnsureUsername }}?secret={{ .Secret }}&issuer={{ .Issuer }}",
	otpString := fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		totpIssuer,
		username,
		twoFactorSecret,
		totpIssuer,
	)

	// encode two factor secret as authenticator-friendly QR code
	qrCode, err := qr.Encode(otpString, qr.L, qr.Auto)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "encoding OTP string")
		return ""
	}

	// scale the QR code so that it's not a PNG for ants.
	qrCode, err = barcode.Scale(qrCode, 256, 256)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "scaling QR code")
		return ""
	}

	// encode the QR code to PNG.
	var b bytes.Buffer
	if err = png.Encode(&b, qrCode); err != nil {
		observability.AcknowledgeError(err, logger, span, "encoding QR code to PNG")
		return ""
	}

	// base64 encode the image for easy HTML use.
	return fmt.Sprintf("%s%s", base64ImagePrefix, base64.StdEncoding.EncodeToString(b.Bytes()))
}

// SelfHandler returns information about the user making the request.
func (s *service) SelfHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// figure out who this is all for.
	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachRequestingUserIDToSpan(span, requester)

	// fetch user data.
	user, err := s.userDataManager.GetUser(ctx, requester)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("no such user")
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, user)
}

// PermissionsHandler returns information about the user making the request.
func (s *service) PermissionsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, sessionCtxDataErr := s.sessionContextDataFetcher(req)
	if sessionCtxDataErr != nil {
		observability.AcknowledgeError(sessionCtxDataErr, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// decode the request.
	permissionsInput := new(types.UserPermissionsRequestInput)
	if decodeErr := s.encoderDecoder.DecodeRequest(ctx, req, permissionsInput); decodeErr != nil {
		observability.AcknowledgeError(decodeErr, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	body := &types.UserPermissionsResponse{
		Permissions: make(map[string]bool),
	}

	for _, perm := range permissionsInput.Permissions {
		p := authorization.Permission(perm)
		hasHouseholdPerm := sessionCtxData.HouseholdPermissions[sessionCtxData.ActiveHouseholdID].HasPermission(p)
		hasServicePerm := sessionCtxData.Requester.ServicePermissions.HasPermission(p)
		body.Permissions[perm] = hasHouseholdPerm || hasServicePerm
	}

	// encode response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, body)
}

// ReadHandler is our read route.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// figure out who this is all for.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	// fetch user data.
	x, err := s.userDataManager.GetUser(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("no such user")
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user from database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// TOTPSecretVerificationHandler accepts a TOTP token as input and returns 200 if the TOTP token
// is validated by the user's TOTP secret.
func (s *service) TOTPSecretVerificationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// decode the request.
	input := new(types.TOTPSecretVerificationInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	logger = logger.WithValue(keys.UserIDKey, input.UserID)

	user, fetchUserErr := s.userDataManager.GetUserWithUnverifiedTwoFactorSecret(ctx, input.UserID)
	if fetchUserErr != nil {
		observability.AcknowledgeError(fetchUserErr, logger, span, "fetching user to verify two factor secret")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	tracing.AttachUserIDToSpan(span, user.ID)
	tracing.AttachUsernameToSpan(span, user.Username)

	if user.TwoFactorSecretVerifiedAt != nil {
		// I suppose if this happens too many times, we might want to keep track of that
		logger.Debug("two factor secret already verified")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "TOTP secret already verified", http.StatusAlreadyReported)
		return
	}

	totpValid := totp.Validate(input.TOTPToken, user.TwoFactorSecret)
	if !totpValid {
		s.encoderDecoder.EncodeInvalidInputResponse(ctx, res)
		return
	}

	if updateUserErr := s.userDataManager.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); updateUserErr != nil {
		observability.AcknowledgeError(updateUserErr, logger, span, "verifying user two factor secret")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.UserDataType,
			EventType:            types.TwoFactorSecretVerifiedCustomerEventType,
			AttributableToUserID: user.ID,
		}

		if publishErr := s.dataChangesPublisher.Publish(ctx, dcm); publishErr != nil {
			observability.AcknowledgeError(publishErr, logger, span, "publishing data change message")
		}
	}

	res.WriteHeader(http.StatusAccepted)
}

// NewTOTPSecretHandler fetches a user, and issues them a new TOTP secret, after validating
// that information received from TOTPSecretRefreshInputContextMiddleware is valid.
func (s *service) NewTOTPSecretHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// decode the request.
	input := new(types.TOTPSecretRefreshInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// fetch user
	user, err := s.userDataManager.GetUser(ctx, sessionCtxData.Requester.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("no such user")
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user from database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if user.TwoFactorSecretVerifiedAt != nil {
		// validate login.
		valid, validationErr := s.authenticator.ValidateLogin(ctx, user.HashedPassword, input.CurrentPassword, user.TwoFactorSecret, input.TOTPToken)
		if validationErr != nil {
			observability.AcknowledgeError(validationErr, logger, span, "validating credentials")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
			return
		} else if !valid {
			observability.AcknowledgeError(validationErr, logger, span, "invalid credentials")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
			return
		}
	}

	// document who this is for.
	tracing.AttachRequestingUserIDToSpan(span, sessionCtxData.Requester.UserID)
	tracing.AttachUsernameToSpan(span, user.Username)
	logger = logger.WithValue(keys.UserIDKey, user.ID)

	// set the two factor secret.
	tfs, err := s.secretGenerator.GenerateBase32EncodedString(ctx, totpSecretSize)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating 2FA secret")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	user.TwoFactorSecret = tfs
	user.TwoFactorSecretVerifiedAt = nil

	// update the user in the database.
	if err = s.userDataManager.UpdateUser(ctx, user); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating 2FA secret")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// let the requester know we're all good.
	result := &types.TOTPSecretRefreshResponse{
		TwoFactorSecret: user.TwoFactorSecret,
		TwoFactorQRCode: s.buildQRCode(ctx, user.Username, user.TwoFactorSecret),
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, result, http.StatusAccepted)
}

// UpdatePasswordHandler updates a user's password.
func (s *service) UpdatePasswordHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// decode the request.
	input := new(types.PasswordUpdateInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx, s.authSettings.MinimumPasswordLength); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	// determine relevant user ID.
	tracing.AttachRequestingUserIDToSpan(span, sessionCtxData.Requester.UserID)
	logger = sessionCtxData.AttachToLogger(logger)

	// make sure everything's on the up-and-up
	user, httpStatus := s.validateCredentialChangeRequest(
		ctx,
		sessionCtxData.Requester.UserID,
		input.CurrentPassword,
		input.TOTPToken,
	)

	// if the above function returns something other than 200, it means some error occurred.
	if httpStatus != http.StatusOK {
		res.WriteHeader(httpStatus)
		return
	}

	tracing.AttachUsernameToSpan(span, user.Username)

	// ensure the password isn't garbage-tier
	if err = passwordvalidator.Validate(input.NewPassword, minimumPasswordEntropy); err != nil {
		logger.WithValue("password_validation_error", err).Debug("invalid password provided")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "new password is too weak!", http.StatusBadRequest)
		return
	}

	// hash the new password.
	newPasswordHash, err := s.authenticator.HashPassword(ctx, input.NewPassword)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "hashing password")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the user.
	if err = s.userDataManager.UpdateUserPassword(ctx, user.ID, newPasswordHash); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// we're all good, log the user out
	http.SetCookie(res, &http.Cookie{MaxAge: -1})
}

func stringPointer(storageProviderPath string) *string {
	return &storageProviderPath
}

// AvatarUploadHandler updates a user's avatar.
func (s *service) AvatarUploadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	logger.Debug("session context data data extracted")

	user, err := s.userDataManager.GetUser(ctx, sessionCtxData.Requester.UserID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching associated user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	logger = logger.WithValue(keys.UserIDKey, user.ID)
	logger.Debug("retrieved user from database")

	img, err := s.imageUploadProcessor.Process(ctx, req, "avatar")
	if err != nil || img == nil {
		observability.AcknowledgeError(err, logger, span, "processing provided avatar upload file")
		s.encoderDecoder.EncodeInvalidInputResponse(ctx, res)
		return
	}

	internalPath := fmt.Sprintf("avatars/%d_%s", time.Now().Unix(), img.Filename)
	logger = logger.WithValue("file_size", len(img.Data)).WithValue("internal_path", internalPath)

	if err = s.uploadManager.SaveFile(ctx, internalPath, img.Data); err != nil {
		observability.AcknowledgeError(err, logger, span, "saving provided avatar")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	user.AvatarSrc = stringPointer(internalPath)

	if err = s.userDataManager.UpdateUser(ctx, user); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user info")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}
}

// ArchiveHandler is a handler for archiving a user.
func (s *service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// figure out who this is for.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	// do the deed.
	err := s.userDataManager.ArchiveUser(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// inform the relatives.
	s.userCounter.Decrement(ctx)

	// we're all good.
	res.WriteHeader(http.StatusNoContent)
}

// RequestUsernameReminderHandler checks for a user with a given email address and notifies them via email if there is a username associated with it.
func (s *service) RequestUsernameReminderHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	input := new(types.UsernameReminderRequestInput)
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

	u, err := s.userDataManager.GetUserByEmail(ctx, input.EmailAddress)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusAccepted)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	msg, emailGenerationErr := email.BuildUsernameReminderEmail(u.EmailAddress, u.Username)
	if emailGenerationErr != nil {
		observability.AcknowledgeError(emailGenerationErr, logger, span, "building username reminder email")
	}

	if emailSendErr := s.emailer.SendEmail(ctx, msg); emailSendErr != nil {
		observability.AcknowledgeError(emailSendErr, logger, span, "sending username reminder email")
	}

	res.WriteHeader(http.StatusAccepted)
}

// CreatePasswordResetTokenHandler rotates the cookie building secret with a new random secret.
func (s *service) CreatePasswordResetTokenHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	input := new(types.PasswordResetTokenCreationRequestInput)
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

	token, err := s.secretGenerator.GenerateBase32EncodedString(ctx, passwordResetTokenSize)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating secret")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := s.userDataManager.GetUserByEmail(ctx, input.EmailAddress)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusAccepted)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	dbInput := &types.PasswordResetTokenDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Token:         token,
		BelongsToUser: u.ID,
		ExpiresAt:     time.Now().Add(30 * time.Minute),
	}

	t, err := s.passwordResetTokenDataManager.CreatePasswordResetToken(ctx, dbInput)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating password reset token")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	msg, emailGenerationErr := email.BuildGeneratedPasswordResetTokenEmail(u.EmailAddress, t)
	if emailGenerationErr != nil {
		observability.AcknowledgeError(emailGenerationErr, logger, span, "building password reset token email")
	}

	if emailSendErr := s.emailer.SendEmail(ctx, msg); emailSendErr != nil {
		observability.AcknowledgeError(emailSendErr, logger, span, "sending password reset token email")
	}

	res.WriteHeader(http.StatusAccepted)
}

// PasswordResetTokenRedemptionHandler rotates the cookie building secret with a new random secret.
func (s *service) PasswordResetTokenRedemptionHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	input := new(types.PasswordResetTokenRedemptionRequestInput)
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

	t, err := s.passwordResetTokenDataManager.GetPasswordResetTokenByToken(ctx, input.Token)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching password reset token")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := s.userDataManager.GetUser(ctx, t.BelongsToUser)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	// ensure the password isn't garbage-tier
	if err = passwordvalidator.Validate(input.NewPassword, minimumPasswordEntropy); err != nil {
		logger.WithValue("password_validation_error", err).Debug("invalid password provided")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "new password is too weak!", http.StatusBadRequest)
		return
	}

	// hash the new password.
	newPasswordHash, err := s.authenticator.HashPassword(ctx, input.NewPassword)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "hashing password")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the user.
	if err = s.userDataManager.UpdateUserPassword(ctx, u.ID, newPasswordHash); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if redemptionErr := s.passwordResetTokenDataManager.RedeemPasswordResetToken(ctx, t.ID); redemptionErr != nil {
		observability.AcknowledgeError(err, logger, span, "redeeming password reset token")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	msg, emailGenerationErr := email.BuildPasswordResetTokenRedeemedEmail(u.EmailAddress)
	if emailGenerationErr != nil {
		observability.AcknowledgeError(emailGenerationErr, logger, span, "building password reset token redemption email")
	}

	if emailSendErr := s.emailer.SendEmail(ctx, msg); emailSendErr != nil {
		observability.AcknowledgeError(emailSendErr, logger, span, "sending password reset token redemption email")
	}

	res.WriteHeader(http.StatusAccepted)
}
