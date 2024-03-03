package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetRecipeRatingRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/ratings/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealID := fakes.BuildFakeID()
		exampleRecipeRating := fakes.BuildFakeRecipeRating()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealID, exampleRecipeRating.ID)

		actual, err := helper.builder.BuildGetRecipeRatingRequest(helper.ctx, exampleMealID, exampleRecipeRating.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleMealID := fakes.BuildFakeID()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetRecipeRatingRequest(helper.ctx, exampleMealID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealID := fakes.BuildFakeID()
		exampleRecipeRating := fakes.BuildFakeRecipeRating()

		actual, err := helper.builder.BuildGetRecipeRatingRequest(helper.ctx, exampleMealID, exampleRecipeRating.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetRecipeRatingsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/ratings"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealID := fakes.BuildFakeID()
		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, exampleMealID)

		actual, err := helper.builder.BuildGetRecipeRatingsRequest(helper.ctx, exampleMealID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealID := fakes.BuildFakeID()
		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipeRatingsRequest(helper.ctx, exampleMealID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateRecipeRatingRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/recipes/%s/ratings"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeRecipeRatingCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleMealID)

		actual, err := helper.builder.BuildCreateRecipeRatingRequest(helper.ctx, exampleMealID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		exampleMealID := fakes.BuildFakeID()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateRecipeRatingRequest(helper.ctx, exampleMealID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleMealID := fakes.BuildFakeID()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateRecipeRatingRequest(helper.ctx, exampleMealID, &types.RecipeRatingCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeRecipeRatingCreationRequestInput()

		actual, err := helper.builder.BuildCreateRecipeRatingRequest(helper.ctx, exampleMealID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateRecipeRatingRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/ratings/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeRating := fakes.BuildFakeRecipeRating()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleRecipeRating.RecipeID, exampleRecipeRating.ID)

		actual, err := helper.builder.BuildUpdateRecipeRatingRequest(helper.ctx, exampleRecipeRating)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateRecipeRatingRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeRating := fakes.BuildFakeRecipeRating()

		actual, err := helper.builder.BuildUpdateRecipeRatingRequest(helper.ctx, exampleRecipeRating)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveRecipeRatingRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/ratings/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealID := fakes.BuildFakeID()
		exampleRecipeRating := fakes.BuildFakeRecipeRating()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleMealID, exampleRecipeRating.ID)

		actual, err := helper.builder.BuildArchiveRecipeRatingRequest(helper.ctx, exampleMealID, exampleRecipeRating.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleMealID := fakes.BuildFakeID()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveRecipeRatingRequest(helper.ctx, exampleMealID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealID := fakes.BuildFakeID()
		exampleRecipeRating := fakes.BuildFakeRecipeRating()

		actual, err := helper.builder.BuildArchiveRecipeRatingRequest(helper.ctx, exampleMealID, exampleRecipeRating.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
