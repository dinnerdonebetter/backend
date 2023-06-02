package validpreparationinstruments

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidPreparationInstrumentsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManager.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidPreparationInstrumentDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidPreparationInstrument, nil)
		helper.service.validPreparationInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("without input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidPreparationInstrumentCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManager.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidPreparationInstrumentDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManager.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidPreparationInstrumentDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidPreparationInstrument, nil)
		helper.service.validPreparationInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidPreparationInstrumentsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidPreparationInstrument{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such valid preparation instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})
}

func TestValidPreparationInstrumentsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationInstrumentList, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.ValidPreparationInstrument]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidPreparationInstrument])(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.ValidPreparationInstrument]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error retrieving valid preparation instruments from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidPreparationInstrument])(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})
}

func TestValidPreparationInstrumentsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)

		dbManager.ValidPreparationInstrumentDataManager.On(
			"UpdateValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument,
		).Return(nil)
		helper.service.validPreparationInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidPreparationInstrumentUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("without input attached to context", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with no such valid preparation instrument", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error retrieving valid preparation instrument from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)

		dbManager.ValidPreparationInstrumentDataManager.On(
			"UpdateValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument,
		).Return(errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)

		dbManager.ValidPreparationInstrumentDataManager.On(
			"UpdateValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument,
		).Return(nil)
		helper.service.validPreparationInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidPreparationInstrumentsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(true, nil)

		dbManager.ValidPreparationInstrumentDataManager.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(nil)
		helper.service.validPreparationInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such valid preparation instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(false, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(false, errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(true, nil)

		dbManager.ValidPreparationInstrumentDataManager.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = dbManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(true, nil)

		dbManager.ValidPreparationInstrumentDataManager.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(nil)
		helper.service.validPreparationInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidPreparationInstrumentsService_SearchByPreparationHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleResponse := fakes.BuildFakeValidPreparationInstrumentList()

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrumentsForPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			testutils.QueryFilterMatcher,
		).Return(exampleResponse, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeResponseWithStatus",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			exampleResponse,
			http.StatusOK,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.SearchByPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = func(request *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.SearchByPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error fetching data from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrumentsForPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			testutils.QueryFilterMatcher,
		).Return((*types.QueryFilteredResult[types.ValidPreparationInstrument])(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.SearchByPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})
}

func TestValidPreparationInstrumentsService_SearchByInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleResponse := fakes.BuildFakeValidPreparationInstrumentList()

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrumentsForInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
			testutils.QueryFilterMatcher,
		).Return(exampleResponse, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeResponseWithStatus",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			exampleResponse,
			http.StatusOK,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.SearchByInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = func(request *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.SearchByInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error fetching data from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrumentsForInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
			testutils.QueryFilterMatcher,
		).Return((*types.QueryFilteredResult[types.ValidPreparationInstrument])(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.SearchByInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})
}
