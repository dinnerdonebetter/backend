package validpreparationinstruments

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// ValidPreparationInstrumentIDURIParamKey is a standard string that we'll use to refer to valid preparation instrument IDs with.
	ValidPreparationInstrumentIDURIParamKey = "validPreparationInstrumentID"
	// ValidPreparationIDURIParamKey is a standard string that we'll use to refer to valid preparation IDs with.
	ValidPreparationIDURIParamKey = "validPreparationID"
	// ValidInstrumentIDURIParamKey is a standard string that we'll use to refer to valid preparation IDs with.
	ValidInstrumentIDURIParamKey = "validInstrumentID"
)

// CreateHandler is our valid preparation instrument creation route.
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
	providedInput := new(types.ValidPreparationInstrumentCreationRequestInput)
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

	input := types.ValidPreparationInstrumentDatabaseCreationInputFromValidPreparationInstrumentCreationInput(providedInput)
	input.ID = ksuid.New().String()

	tracing.AttachValidPreparationInstrumentIDToSpan(span, input.ID)

	validPreparationInstrument, err := s.validPreparationInstrumentDataManager.CreateValidPreparationInstrument(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid preparation instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                   types.ValidPreparationInstrumentDataType,
			EventType:                  types.ValidPreparationInstrumentCreatedCustomerEventType,
			ValidPreparationInstrument: validPreparationInstrument,
			AttributableToUserID:       sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
		}
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validPreparationInstrument, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid preparation instrument.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine valid preparation instrument ID.
	validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	// fetch valid preparation instrument from database.
	x, err := s.validPreparationInstrumentDataManager.GetValidPreparationInstrument(ctx, validPreparationInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ListHandler is our list route.
func (s *service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	filter := types.ExtractQueryFilter(req)
	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, string(filter.SortBy))

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	validPreparationInstruments, err := s.validPreparationInstrumentDataManager.GetValidPreparationInstruments(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validPreparationInstruments = &types.ValidPreparationInstrumentList{ValidPreparationInstruments: []*types.ValidPreparationInstrument{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation instruments")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validPreparationInstruments)
}

// UpdateHandler returns a handler that updates a valid preparation instrument.
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

	// check for parsed input attached to session context data.
	input := new(types.ValidPreparationInstrumentUpdateRequestInput)
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

	// determine valid preparation instrument ID.
	validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	// fetch valid preparation instrument from database.
	validPreparationInstrument, err := s.validPreparationInstrumentDataManager.GetValidPreparationInstrument(ctx, validPreparationInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation instrument for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the valid preparation instrument.
	validPreparationInstrument.Update(input)

	if err = s.validPreparationInstrumentDataManager.UpdateValidPreparationInstrument(ctx, validPreparationInstrument); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid preparation instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                   types.ValidPreparationInstrumentDataType,
			EventType:                  types.ValidPreparationInstrumentUpdatedCustomerEventType,
			ValidPreparationInstrument: validPreparationInstrument,
			AttributableToUserID:       sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validPreparationInstrument)
}

// ArchiveHandler returns a handler that archives a valid preparation instrument.
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

	// determine valid preparation instrument ID.
	validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	exists, existenceCheckErr := s.validPreparationInstrumentDataManager.ValidPreparationInstrumentExists(ctx, validPreparationInstrumentID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking valid preparation instrument existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.validPreparationInstrumentDataManager.ArchiveValidPreparationInstrument(ctx, validPreparationInstrumentID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid preparation instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.ValidPreparationInstrumentDataType,
			EventType:            types.ValidPreparationInstrumentArchivedCustomerEventType,
			AttributableToUserID: sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// SearchByPreparationHandler is our valid preparation instrument search route for preparations.
func (s *service) SearchByPreparationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilter(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))

	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, string(filter.SortBy))

	validPreparationID := s.validPreparationIDFetcher(req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	validPreparationInstruments, err := s.validPreparationInstrumentDataManager.GetValidPreparationInstrumentsForPreparation(ctx, validPreparationID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid preparation instruments")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validPreparationInstruments, http.StatusOK)
}

// SearchByInstrumentHandler is our valid preparation instrument search route for instruments.
func (s *service) SearchByInstrumentHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilter(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))

	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, string(filter.SortBy))

	validInstrumentID := s.validInstrumentIDFetcher(req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	validPreparationInstruments, err := s.validPreparationInstrumentDataManager.GetValidPreparationInstrumentsForInstrument(ctx, validInstrumentID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid preparation instruments")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validPreparationInstruments, http.StatusOK)
}