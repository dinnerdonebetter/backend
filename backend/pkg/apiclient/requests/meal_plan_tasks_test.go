package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetMealPlanTaskRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanTask := fakes.BuildFakeMealPlanTask()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanTask.ID)

		actual, err := helper.builder.BuildGetMealPlanTaskRequest(helper.ctx, exampleMealPlanID, exampleMealPlanTask.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan task ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetMealPlanTaskRequest(helper.ctx, exampleMealPlanID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanTask := fakes.BuildFakeMealPlanTask()

		actual, err := helper.builder.BuildGetMealPlanTaskRequest(helper.ctx, exampleMealPlanID, exampleMealPlanTask.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateMealPlanTaskRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanTaskCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, exampleMealPlanID)

		actual, err := helper.builder.BuildCreateMealPlanTaskRequest(helper.ctx, exampleMealPlanID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan task ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateMealPlanTaskRequest(helper.ctx, exampleMealPlanID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanTaskCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealPlanTaskRequest(helper.ctx, exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetMealPlanTasksRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		exampleMealPlanID := fakes.BuildFakeID()
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, exampleMealPlanID)

		actual, err := helper.builder.BuildGetMealPlanTasksRequest(helper.ctx, exampleMealPlanID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)
		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetMealPlanTasksRequest(helper.ctx, exampleMealPlanID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildChangeMealPlanTaskStatusRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()
		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, exampleMealPlanID, exampleInput.ID)

		actual, err := helper.builder.BuildChangeMealPlanTaskStatusRequest(helper.ctx, exampleMealPlanID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan task ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildChangeMealPlanTaskStatusRequest(helper.ctx, exampleMealPlanID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()

		actual, err := helper.builder.BuildChangeMealPlanTaskStatusRequest(helper.ctx, exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
