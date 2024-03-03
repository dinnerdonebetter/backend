package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, "", exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, "", exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetMealPlanOptionVotesRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, "", exampleMealPlanEventID, exampleMealPlanOptionID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/meal_plans/%s/events/%s/vote"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleMealPlanID, exampleMealPlanEventID)

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, "", exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, &types.MealPlanOptionVoteCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionVote.BelongsToMealPlanOption, exampleMealPlanOptionVote.ID)

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionVote)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, "", "", exampleMealPlanOptionVote)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionVote)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, "", exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, "", exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
