package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetUserIngredientPreferencesRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/user_ingredient_preferences"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetUserIngredientPreferencesRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetUserIngredientPreferencesRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateUserIngredientPreferenceRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/user_ingredient_preferences"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeUserIngredientPreferenceCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateUserIngredientPreferenceRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateUserIngredientPreferenceRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateUserIngredientPreferenceRequest(helper.ctx, &types.UserIngredientPreferenceCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeUserIngredientPreferenceCreationRequestInput()

		actual, err := helper.builder.BuildCreateUserIngredientPreferenceRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateUserIngredientPreferenceRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/user_ingredient_preferences/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleUserIngredientPreference.ID)

		actual, err := helper.builder.BuildUpdateUserIngredientPreferenceRequest(helper.ctx, exampleUserIngredientPreference)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateUserIngredientPreferenceRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		actual, err := helper.builder.BuildUpdateUserIngredientPreferenceRequest(helper.ctx, exampleUserIngredientPreference)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveUserIngredientPreferenceRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/user_ingredient_preferences/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleUserIngredientPreference.ID)

		actual, err := helper.builder.BuildArchiveUserIngredientPreferenceRequest(helper.ctx, exampleUserIngredientPreference.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveUserIngredientPreferenceRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()

		actual, err := helper.builder.BuildArchiveUserIngredientPreferenceRequest(helper.ctx, exampleUserIngredientPreference.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
