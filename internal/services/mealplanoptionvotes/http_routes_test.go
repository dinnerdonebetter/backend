package mealplanoptionvotes

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/encoding/mock"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMealPlanOptionVotesService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanEventDataManager.On(
			"MealPlanEventIsEligibleForVoting",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanOptionVotesDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleMealPlanOptionVotes, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)

		dbManager.MealPlanDataManager.On(
			"AttemptToFinalizeMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)
		helper.service.dataManager = dbManager

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

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("without being eligible for voting", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanEventDataManager.On(
			"MealPlanEventIsEligibleForVoting",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(false, nil)
		helper.service.dataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error checking voting eligibility", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanEventDataManager.On(
			"MealPlanEventIsEligibleForVoting",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(false, errors.New("blah"))
		helper.service.dataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error writing create to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanEventDataManager.On(
			"MealPlanEventIsEligibleForVoting",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanOptionVotesDatabaseCreationInput) bool { return true }),
		).Return([]*types.MealPlanOptionVote(nil), errors.New("blah"))
		helper.service.dataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing first event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanEventDataManager.On(
			"MealPlanEventIsEligibleForVoting",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanOptionVotesDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleMealPlanOptionVotes, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)

		dbManager.MealPlanDataManager.On(
			"AttemptToFinalizeMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)
		helper.service.dataManager = dbManager

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

	T.Run("with error finalizing meal plan option", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanEventDataManager.On(
			"MealPlanEventIsEligibleForVoting",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanOptionVotesDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleMealPlanOptionVotes, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleHousehold.ID,
		).Return(false, errors.New("blah"))
		helper.service.dataManager = dbManager

		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing second event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanEventDataManager.On(
			"MealPlanEventIsEligibleForVoting",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanOptionVotesDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleMealPlanOptionVotes, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))

		dbManager.MealPlanDataManager.On(
			"AttemptToFinalizeMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)
		helper.service.dataManager = dbManager

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

	T.Run("with error attempting to finalize complete meal plan", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanEventDataManager.On(
			"MealPlanEventIsEligibleForVoting",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanOptionVotesDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleMealPlanOptionVotes, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		dbManager.MealPlanDataManager.On(
			"AttemptToFinalizeMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleHousehold.ID,
		).Return(false, errors.New("blah"))
		helper.service.dataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing final event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanEventDataManager.On(
			"MealPlanEventIsEligibleForVoting",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanOptionVotesDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleMealPlanOptionVotes, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)

		dbManager.MealPlanDataManager.On(
			"AttemptToFinalizeMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)
		helper.service.dataManager = dbManager

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

func TestMealPlanOptionVotesService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(helper.exampleMealPlanOptionVote, nil)
		helper.service.dataManager = dbManager

		encoderDecoder := mockencoding.mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.MealPlanOptionVote{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, encoderDecoder)
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

	T.Run("with no such meal plan option vote in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return((*types.MealPlanOptionVote)(nil), sql.ErrNoRows)
		helper.service.dataManager = dbManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return((*types.MealPlanOptionVote)(nil), errors.New("blah"))
		helper.service.dataManager = dbManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, encoderDecoder)
	})
}

func TestMealPlanOptionVotesService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealPlanOptionVoteList := fakes.BuildFakeMealPlanOptionVoteList()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVotes",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleMealPlanOptionVoteList, nil)
		helper.service.dataManager = dbManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.MealPlanOptionVote]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, encoderDecoder)
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

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVotes",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.MealPlanOptionVote])(nil), sql.ErrNoRows)
		helper.service.dataManager = dbManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.MealPlanOptionVote]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, encoderDecoder)
	})

	T.Run("with error retrieving meal plan option votes from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVotes",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.MealPlanOptionVote])(nil), errors.New("blah"))
		helper.service.dataManager = dbManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, encoderDecoder)
	})
}

func TestMealPlanOptionVotesService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(helper.exampleMealPlanOptionVote, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"UpdateMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlanOptionVote,
		).Return(nil)
		helper.service.dataManager = dbManager

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

		exampleCreationInput := &types.MealPlanOptionVoteUpdateRequestInput{}
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

	T.Run("with no such meal plan option vote", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return((*types.MealPlanOptionVote)(nil), sql.ErrNoRows)
		helper.service.dataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(helper.exampleMealPlanOptionVote, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"UpdateMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlanOptionVote,
		).Return(errors.New("blah"))
		helper.service.dataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error retrieving meal plan option vote from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return((*types.MealPlanOptionVote)(nil), errors.New("blah"))
		helper.service.dataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(helper.exampleMealPlanOptionVote, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"UpdateMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlanOptionVote,
		).Return(nil)
		helper.service.dataManager = dbManager

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

func TestMealPlanOptionVotesService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(nil)
		helper.service.dataManager = dbManager

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

	T.Run("with no such meal plan option vote in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(false, nil)
		helper.service.dataManager = dbManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, encoderDecoder)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(false, errors.New("blah"))
		helper.service.dataManager = dbManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(errors.New("blah"))
		helper.service.dataManager = dbManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(nil)
		helper.service.dataManager = dbManager

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
