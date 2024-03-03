package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetValidIngredientStateIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidIngredientStateIngredient.ID)

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientRequest(helper.ctx, exampleValidIngredientStateIngredient.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientRequest(helper.ctx, exampleValidIngredientStateIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidIngredientStateIngredientsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidIngredientStateIngredientsForIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/by_ingredient/%s"

	exampleIngredient := fakes.BuildFakeValidIngredient()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, exampleIngredient.ID)

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientsForIngredientRequest(helper.ctx, exampleIngredient.ID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientsForIngredientRequest(helper.ctx, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientsForIngredientRequest(helper.ctx, exampleIngredient.ID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidIngredientStateIngredientsForPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/%s"

	examplePreparation := fakes.BuildFakeValidPreparation()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, examplePreparation.ID)

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientsForPreparationRequest(helper.ctx, examplePreparation.ID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientsForPreparationRequest(helper.ctx, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidIngredientStateIngredientsForPreparationRequest(helper.ctx, examplePreparation.ID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidIngredientStateIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_ingredient_state_ingredients"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidIngredientStateIngredientCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidIngredientStateIngredientRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidIngredientStateIngredientRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidIngredientStateIngredientRequest(helper.ctx, &types.ValidIngredientStateIngredientCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidIngredientStateIngredientCreationRequestInput()

		actual, err := helper.builder.BuildCreateValidIngredientStateIngredientRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidIngredientStateIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidIngredientStateIngredient.ID)

		actual, err := helper.builder.BuildUpdateValidIngredientStateIngredientRequest(helper.ctx, exampleValidIngredientStateIngredient)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidIngredientStateIngredientRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		actual, err := helper.builder.BuildUpdateValidIngredientStateIngredientRequest(helper.ctx, exampleValidIngredientStateIngredient)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidIngredientStateIngredientRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidIngredientStateIngredient.ID)

		actual, err := helper.builder.BuildArchiveValidIngredientStateIngredientRequest(helper.ctx, exampleValidIngredientStateIngredient.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidIngredientStateIngredientRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		actual, err := helper.builder.BuildArchiveValidIngredientStateIngredientRequest(helper.ctx, exampleValidIngredientStateIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
