package mealplantasks

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMealPlanTasksService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := &mocktypes.MealPlanTaskDataManagerMock{}
		mealPlanTaskDataManager.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return(helper.exampleMealPlanTask, nil)
		helper.service.mealPlanTaskDataManager = mealPlanTaskDataManager

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanTask)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal plan in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := &mocktypes.MealPlanTaskDataManagerMock{}
		mealPlanTaskDataManager.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return((*types.MealPlanTask)(nil), sql.ErrNoRows)
		helper.service.mealPlanTaskDataManager = mealPlanTaskDataManager

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := &mocktypes.MealPlanTaskDataManagerMock{}
		mealPlanTaskDataManager.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return((*types.MealPlanTask)(nil), errors.New("blah"))
		helper.service.mealPlanTaskDataManager = mealPlanTaskDataManager

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})
}

func TestMealPlanTasksService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealPlanTaskList := fakes.BuildFakeMealPlanTaskList().Data

		mealPlanTaskDataManager := &mocktypes.MealPlanTaskDataManagerMock{}
		mealPlanTaskDataManager.On(
			"GetMealPlanTasksForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return(exampleMealPlanTaskList, nil)
		helper.service.mealPlanTaskDataManager = mealPlanTaskDataManager

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleMealPlanTaskList)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := &mocktypes.MealPlanTaskDataManagerMock{}
		mealPlanTaskDataManager.On(
			"GetMealPlanTasksForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.MealPlanTask(nil), sql.ErrNoRows)
		helper.service.mealPlanTaskDataManager = mealPlanTaskDataManager

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})

	T.Run("with error retrieving meal plans from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := &mocktypes.MealPlanTaskDataManagerMock{}
		mealPlanTaskDataManager.On(
			"GetMealPlanTasksForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.MealPlanTask(nil), errors.New("blah"))
		helper.service.mealPlanTaskDataManager = mealPlanTaskDataManager

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})
}

func TestMealPlanTasksService_StatusChangeHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleStatusChangeInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()
		exampleStatusChangeInput.ID = helper.exampleMealPlanTask.ID
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleStatusChangeInput)

		expectedPrepStep := helper.exampleMealPlanTask
		expectedPrepStep.Status = *exampleStatusChangeInput.Status
		expectedPrepStep.StatusExplanation = exampleStatusChangeInput.StatusExplanation

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return(expectedPrepStep, nil)

		dbManager.MealPlanTaskDataManagerMock.On(
			"ChangeMealPlanTaskStatus",
			testutils.ContextMatcher,
			exampleStatusChangeInput,
		).Return(nil)
		helper.service.mealPlanTaskDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.StatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanTask)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.StatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleCreationInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()
		exampleCreationInput.ID = helper.exampleMealPlanTask.ID
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return(helper.exampleMealPlanTask, nil)

		dbManager.MealPlanTaskDataManagerMock.On(
			"ChangeMealPlanTaskStatus",
			testutils.ContextMatcher,
			exampleCreationInput,
		).Return(errors.New("blah"))
		helper.service.mealPlanTaskDataManager = dbManager

		helper.service.StatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleCreationInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()
		exampleCreationInput.ID = helper.exampleMealPlanTask.ID
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		expectedPrepStep := helper.exampleMealPlanTask
		expectedPrepStep.Status = *exampleCreationInput.Status
		expectedPrepStep.StatusExplanation = exampleCreationInput.StatusExplanation

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return(expectedPrepStep, nil)

		dbManager.MealPlanTaskDataManagerMock.On(
			"ChangeMealPlanTaskStatus",
			testutils.ContextMatcher,
			exampleCreationInput,
		).Return(nil)
		helper.service.mealPlanTaskDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.StatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanTask)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
