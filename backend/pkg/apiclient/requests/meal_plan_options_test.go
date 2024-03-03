package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetMealPlanOptionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)

		actual, err := helper.builder.BuildGetMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildGetMealPlanOptionRequest(helper.ctx, "", exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildGetMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetMealPlanOptionsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, exampleMealPlanID, exampleMealPlanEventID)

		actual, err := helper.builder.BuildGetMealPlanOptionsRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()
		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionsRequest(helper.ctx, "", exampleMealPlanEventID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionsRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateMealPlanOptionRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/meal_plans/%s/events/%s/options"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleMealPlanID, exampleMealPlanEventID)

		actual, err := helper.builder.BuildCreateMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateMealPlanOptionRequest(helper.ctx, exampleMealPlanEventID, exampleMealPlanEventID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateMealPlanOptionRequest(helper.ctx, exampleMealPlanEventID, exampleMealPlanEventID, &types.MealPlanOptionCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealPlanOptionRequest(helper.ctx, exampleMealPlanEventID, exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateMealPlanOptionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanOption.BelongsToMealPlanEvent, exampleMealPlanOption.ID)

		actual, err := helper.builder.BuildUpdateMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOption)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildUpdateMealPlanOptionRequest(helper.ctx, exampleMealPlanEventID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildUpdateMealPlanOptionRequest(helper.ctx, exampleMealPlanEventID, exampleMealPlanOption)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveMealPlanOptionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)

		actual, err := helper.builder.BuildArchiveMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildArchiveMealPlanOptionRequest(helper.ctx, "", exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		actual, err := helper.builder.BuildArchiveMealPlanOptionRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
