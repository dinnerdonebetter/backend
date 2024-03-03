package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildRunFinalizeMealPlansWorkerRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/workers/finalize_meal_plans"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleInput := fakes.BuildFakeFinalizeMealPlansRequest()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildRunFinalizeMealPlansWorkerRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildRunFinalizeMealPlansWorkerRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleInput := fakes.BuildFakeFinalizeMealPlansRequest()

		actual, err := helper.builder.BuildRunFinalizeMealPlansWorkerRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildRunMealPlanGroceryListInitializationWorkerRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/workers/meal_plan_grocery_list_init"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildRunMealPlanGroceryListInitializationWorkerRequest(helper.ctx)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		actual, err := helper.builder.BuildRunMealPlanGroceryListInitializationWorkerRequest(helper.ctx)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildRunMealPlanTaskCreationWorkerRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/workers/meal_plan_tasks"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildRunMealPlanTaskCreationWorkerRequest(helper.ctx)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		actual, err := helper.builder.BuildRunMealPlanTaskCreationWorkerRequest(helper.ctx)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
