package mealplangrocerylistitems

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	// MealPlanGroceryListItemIDURIParamKey is a standard string that we'll use to refer to meal plan grocery list item IDs with.
	MealPlanGroceryListItemIDURIParamKey = "mealPlanGroceryListItemID"
)

// CreateHandler is our meal plan grocery list item creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanGroceryListItemCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	input := converters.ConvertMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItemDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, input.ID)

	logger = logger.WithValue("input", input)

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	input.BelongsToMealPlan = mealPlanID

	mealPlanGroceryListItem, err := s.mealPlanGroceryListItemDataManager.CreateMealPlanGroceryListItem(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating meal plan")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:               types.MealPlanGroceryListItemCreatedCustomerEventType,
		MealPlanID:              mealPlanID,
		MealPlanGroceryListItem: mealPlanGroceryListItem,
		HouseholdID:             sessionCtxData.ActiveHouseholdID,
		UserID:                  sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, mealPlanGroceryListItem, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a meal plan grocery list item.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan grocery list item ID.
	mealPlanGroceryListItemID := s.mealPlanGroceryListItemIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanGroceryListItemID)
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	// fetch meal plan grocery list item from database.
	x, err := s.mealPlanGroceryListItemDataManager.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan grocery list item")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ListByMealPlanHandler is our list route.
func (s *service) ListByMealPlanHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req)
	logger = filter.AttachToLogger(logger)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	mealPlanGroceryListItems, err := s.mealPlanGroceryListItemDataManager.GetMealPlanGroceryListItemsForMealPlan(ctx, mealPlanID)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		mealPlanGroceryListItems = []*types.MealPlanGroceryListItem{}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan grocery list items for meal plan")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, mealPlanGroceryListItems)
}

// UpdateHandler returns a handler that updates a meal plan grocery list item.
func (s *service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, sessionCtxFetchErr := s.sessionContextDataFetcher(req)
	if sessionCtxFetchErr != nil {
		observability.AcknowledgeError(sessionCtxFetchErr, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan grocery list item ID.
	mealPlanGroceryListItemID := s.mealPlanGroceryListItemIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanGroceryListItemID)
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanGroceryListItemUpdateRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err := providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	mealPlanGroceryListItem, fetchMealPlanGroceryListItemErr := s.mealPlanGroceryListItemDataManager.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)
	if fetchMealPlanGroceryListItemErr != nil {
		if errors.Is(fetchMealPlanGroceryListItemErr, sql.ErrNoRows) {
			errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
			return
		}
		observability.AcknowledgeError(fetchMealPlanGroceryListItemErr, logger, span, "checking meal plan grocery list item existence")

		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	mealPlanGroceryListItem.Update(providedInput)

	if err := s.mealPlanGroceryListItemDataManager.UpdateMealPlanGroceryListItem(ctx, mealPlanGroceryListItem); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan grocery list item")

		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:                 types.MealPlanGroceryListItemUpdatedCustomerEventType,
		MealPlanGroceryListItem:   mealPlanGroceryListItem,
		MealPlanGroceryListItemID: mealPlanGroceryListItemID,
		HouseholdID:               sessionCtxData.ActiveHouseholdID,
		UserID:                    sessionCtxData.Requester.UserID,
	}

	if err := s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	s.encoderDecoder.RespondWithData(ctx, res, mealPlanGroceryListItem)
}

// ArchiveHandler returns a GET handler that returns a meal plan grocery list item.
func (s *service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan grocery list item ID.
	mealPlanGroceryListItemID := s.mealPlanGroceryListItemIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanGroceryListItemID)
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	// check that meal plan grocery list item exists in database.
	x, err := s.mealPlanGroceryListItemDataManager.MealPlanGroceryListItemExists(ctx, mealPlanID, mealPlanGroceryListItemID)
	if errors.Is(err, sql.ErrNoRows) || !x {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "checking meal plan grocery list item existence")

		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// fetch meal plan grocery list item from database.
	err = s.mealPlanGroceryListItemDataManager.ArchiveMealPlanGroceryListItem(ctx, mealPlanGroceryListItemID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan grocery list item")

		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:                 types.MealPlanGroceryListItemArchivedCustomerEventType,
		MealPlanGroceryListItemID: mealPlanGroceryListItemID,
		HouseholdID:               sessionCtxData.ActiveHouseholdID,
		UserID:                    sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	responseValue := &types.APIResponse[*types.Webhook]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}
