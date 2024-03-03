package validmeasurementunitconversions

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
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

func TestValidMeasurementUnitConversionsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"CreateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversionDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidMeasurementUnitConversion, nil)
		helper.service.validMeasurementUnitConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnitConversion)
		assert.NoError(t, actual.Error.AsError())

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
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidMeasurementUnitConversionCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"CreateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversionDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidMeasurementUnitConversion)(nil), errors.New("blah"))
		helper.service.validMeasurementUnitConversionDataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"CreateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversionDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidMeasurementUnitConversion, nil)
		helper.service.validMeasurementUnitConversionDataManager = dbManager

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

func TestValidMeasurementUnitConversionsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(helper.exampleValidMeasurementUnitConversion, nil)
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnitConversion)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), sql.ErrNoRows)
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), errors.New("blah"))
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})
}

func TestValidMeasurementUnitConversionsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(helper.exampleValidMeasurementUnitConversion, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"UpdateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversion) bool { return true }),
		).Return(nil)
		helper.service.validMeasurementUnitConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnitConversion)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidMeasurementUnitConversionUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
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
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), sql.ErrNoRows)
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error retrieving valid measurement conversion from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), errors.New("blah"))
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with problem writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(helper.exampleValidMeasurementUnitConversion, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"UpdateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversion) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.validMeasurementUnitConversionDataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(helper.exampleValidMeasurementUnitConversion, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"UpdateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversion) bool { return true }),
		).Return(nil)
		helper.service.validMeasurementUnitConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnitConversion)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidMeasurementUnitConversionsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ValidMeasurementUnitConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ArchiveValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(nil)
		helper.service.validMeasurementUnitConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"ValidMeasurementUnitConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(false, nil)
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"ValidMeasurementUnitConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(false, errors.New("blah"))
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ValidMeasurementUnitConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ArchiveValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(errors.New("blah"))
		helper.service.validMeasurementUnitConversionDataManager = dbManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ValidMeasurementUnitConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ArchiveValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(nil)
		helper.service.validMeasurementUnitConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidMeasurementUnitConversionsService_FromMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		expected := []*types.ValidMeasurementUnitConversion{helper.exampleValidMeasurementUnitConversion}

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversionsFromUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(expected, nil)
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.FromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, expected)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.FromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversionsFromUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), sql.ErrNoRows)
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.FromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversionsFromUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), errors.New("blah"))
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.FromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})
}

func TestValidMeasurementUnitConversionsService_ToMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		expected := []*types.ValidMeasurementUnitConversion{helper.exampleValidMeasurementUnitConversion}

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversionsToUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(expected, nil)
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.ToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, expected)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversionsToUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), sql.ErrNoRows)
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.ToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := &mocktypes.ValidMeasurementUnitConversionDataManagerMock{}
		validMeasurementUnitConversionDataManager.On(
			"GetValidMeasurementUnitConversionsToUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), errors.New("blah"))
		helper.service.validMeasurementUnitConversionDataManager = validMeasurementUnitConversionDataManager

		helper.service.ToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})
}
