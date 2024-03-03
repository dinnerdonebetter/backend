package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetRecipePrepTaskRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleRecipeID, exampleRecipePrepTask.ID)

		actual, err := helper.builder.BuildGetRecipePrepTaskRequest(helper.ctx, exampleRecipeID, exampleRecipePrepTask.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		actual, err := helper.builder.BuildGetRecipePrepTaskRequest(helper.ctx, "", exampleRecipePrepTask.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetRecipePrepTaskRequest(helper.ctx, exampleRecipeID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		actual, err := helper.builder.BuildGetRecipePrepTaskRequest(helper.ctx, exampleRecipeID, exampleRecipePrepTask.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetRecipePrepTasksRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, exampleRecipeID)

		actual, err := helper.builder.BuildGetRecipePrepTasksRequest(helper.ctx, exampleRecipeID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipePrepTasksRequest(helper.ctx, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetRecipePrepTasksRequest(helper.ctx, exampleRecipeID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateRecipePrepTaskRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/recipes/%s/prep_tasks"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeRecipePrepTaskCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleInput.BelongsToRecipe)

		actual, err := helper.builder.BuildCreateRecipePrepTaskRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateRecipePrepTaskRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateRecipePrepTaskRequest(helper.ctx, &types.RecipePrepTaskCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeRecipePrepTaskCreationRequestInput()

		actual, err := helper.builder.BuildCreateRecipePrepTaskRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateRecipePrepTaskRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleRecipePrepTask.BelongsToRecipe, exampleRecipePrepTask.ID)

		actual, err := helper.builder.BuildUpdateRecipePrepTaskRequest(helper.ctx, exampleRecipePrepTask)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateRecipePrepTaskRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		actual, err := helper.builder.BuildUpdateRecipePrepTaskRequest(helper.ctx, exampleRecipePrepTask)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveRecipePrepTaskRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleRecipeID, exampleRecipePrepTask.ID)

		actual, err := helper.builder.BuildArchiveRecipePrepTaskRequest(helper.ctx, exampleRecipeID, exampleRecipePrepTask.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		actual, err := helper.builder.BuildArchiveRecipePrepTaskRequest(helper.ctx, "", exampleRecipePrepTask.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleRecipeID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveRecipePrepTaskRequest(helper.ctx, exampleRecipeID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		actual, err := helper.builder.BuildArchiveRecipePrepTaskRequest(helper.ctx, exampleRecipeID, exampleRecipePrepTask.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
