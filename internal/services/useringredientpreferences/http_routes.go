package useringredientpreferences

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	// UserIngredientPreferenceIDURIParamKey is a standard string that we'll use to refer to user ingredient preference IDs with.
	UserIngredientPreferenceIDURIParamKey = "userIngredientPreferenceID"
)

// CreateHandler is our user ingredient preference creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	providedInput := new(types.UserIngredientPreferenceCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	input := converters.ConvertUserIngredientPreferenceCreationRequestInputToUserIngredientPreferenceDatabaseCreationInput(providedInput)
	input.BelongsToUser = sessionCtxData.Requester.UserID

	userIngredientPreference, err := s.userIngredientPreferenceDataManager.CreateUserIngredientPreference(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating user ingredient preferences")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:                 types.UserIngredientPreferenceCreatedCustomerEventType,
		UserIngredientPreferences: userIngredientPreference,
		UserID:                    sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, userIngredientPreference, http.StatusCreated)
}

// ListHandler is our list route.
func (s *service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy)

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	requester := sessionCtxData.Requester.UserID

	userIngredientPreferences, err := s.userIngredientPreferenceDataManager.GetUserIngredientPreferences(ctx, requester, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		userIngredientPreferences = &types.QueryFilteredResult[types.UserIngredientPreference]{Data: []*types.UserIngredientPreference{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving user ingredient preferences")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, userIngredientPreferences)
}

// UpdateHandler returns a handler that updates a user ingredient preference.
func (s *service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	requester := sessionCtxData.Requester.UserID

	// check for parsed input attached to session context data.
	input := new(types.UserIngredientPreferenceUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		logger.Error(err, "error encountered decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.Error(err, "provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	// determine user ingredient preference ID.
	userIngredientPreferenceID := s.userIngredientPreferenceIDFetcher(req)
	tracing.AttachUserIngredientPreferenceIDToSpan(span, userIngredientPreferenceID)
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)

	// fetch user ingredient preference from database.
	userIngredientPreference, err := s.userIngredientPreferenceDataManager.GetUserIngredientPreference(ctx, userIngredientPreferenceID, requester)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving user ingredient preference for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the user ingredient preference.
	userIngredientPreference.Update(input)

	if err = s.userIngredientPreferenceDataManager.UpdateUserIngredientPreference(ctx, userIngredientPreference); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user ingredient preferences")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:                 types.UserIngredientPreferenceUpdatedCustomerEventType,
		UserIngredientPreferences: []*types.UserIngredientPreference{userIngredientPreference},
		UserID:                    sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, userIngredientPreference)
}

// ArchiveHandler returns a handler that archives a user ingredient preference.
func (s *service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	requester := sessionCtxData.Requester.UserID

	// determine user ingredient preference ID.
	userIngredientPreferenceID := s.userIngredientPreferenceIDFetcher(req)
	tracing.AttachUserIngredientPreferenceIDToSpan(span, userIngredientPreferenceID)
	logger = logger.WithValue(keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)

	exists, existenceCheckErr := s.userIngredientPreferenceDataManager.UserIngredientPreferenceExists(ctx, userIngredientPreferenceID, requester)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking user ingredient preference existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.userIngredientPreferenceDataManager.ArchiveUserIngredientPreference(ctx, userIngredientPreferenceID, requester); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving user ingredient preferences")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.UserIngredientPreferenceArchivedCustomerEventType,
		UserID:    sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
