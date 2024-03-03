package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetValidIngredientGroupRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_groups/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidIngredientGroup.ID)

		actual, err := helper.builder.BuildGetValidIngredientGroupRequest(helper.ctx, exampleValidIngredientGroup.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidIngredientGroupRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		actual, err := helper.builder.BuildGetValidIngredientGroupRequest(helper.ctx, exampleValidIngredientGroup.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidIngredientGroupsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_groups"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetValidIngredientGroupsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidIngredientGroupsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildSearchValidIngredientGroupsRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_ingredient_groups/search"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"
		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)

		actual, err := helper.builder.BuildSearchValidIngredientGroupsRequest(helper.ctx, exampleQuery, *limit)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		actual, err := helper.builder.BuildSearchValidIngredientGroupsRequest(helper.ctx, exampleQuery, *limit)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidIngredientGroupRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_ingredient_groups"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidIngredientGroupRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidIngredientGroupRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidIngredientGroupRequest(helper.ctx, &types.ValidIngredientGroupCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()

		actual, err := helper.builder.BuildCreateValidIngredientGroupRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidIngredientGroupRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_groups/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidIngredientGroup.ID)

		actual, err := helper.builder.BuildUpdateValidIngredientGroupRequest(helper.ctx, exampleValidIngredientGroup)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidIngredientGroupRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		actual, err := helper.builder.BuildUpdateValidIngredientGroupRequest(helper.ctx, exampleValidIngredientGroup)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidIngredientGroupRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_groups/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidIngredientGroup.ID)

		actual, err := helper.builder.BuildArchiveValidIngredientGroupRequest(helper.ctx, exampleValidIngredientGroup.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidIngredientGroupRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		actual, err := helper.builder.BuildArchiveValidIngredientGroupRequest(helper.ctx, exampleValidIngredientGroup.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
