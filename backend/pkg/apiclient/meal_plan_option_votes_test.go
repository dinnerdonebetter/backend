package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestMealPlanOptionVotes(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlanOptionVotesTestSuite))
}

type mealPlanOptionVotesBaseSuite struct {
	suite.Suite
	ctx                                   context.Context
	exampleMealPlanOptionVote             *types.MealPlanOptionVote
	exampleMealPlanOptionVoteResponse     *types.APIResponse[*types.MealPlanOptionVote]
	exampleMealPlanOptionVoteListResponse *types.APIResponse[[]*types.MealPlanOptionVote]
	exampleMealPlanID                     string
	exampleMealPlanEventID                string
	exampleMealPlanOptionID               string
	exampleMealPlanOptionVoteList         []*types.MealPlanOptionVote
}

var _ suite.SetupTestSuite = (*mealPlanOptionVotesBaseSuite)(nil)

func (s *mealPlanOptionVotesBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlanID = fakes.BuildFakeID()
	s.exampleMealPlanEventID = fakes.BuildFakeID()
	s.exampleMealPlanOptionID = fakes.BuildFakeID()
	s.exampleMealPlanOptionVote = fakes.BuildFakeMealPlanOptionVote()
	s.exampleMealPlanOptionVote.BelongsToMealPlanOption = s.exampleMealPlanOptionID
	exampleMealPlanOptionVoteList := fakes.BuildFakeMealPlanOptionVoteList()
	s.exampleMealPlanOptionVoteList = exampleMealPlanOptionVoteList.Data
	s.exampleMealPlanOptionVoteResponse = &types.APIResponse[*types.MealPlanOptionVote]{
		Data: s.exampleMealPlanOptionVote,
	}
	s.exampleMealPlanOptionVoteListResponse = &types.APIResponse[[]*types.MealPlanOptionVote]{
		Data:       s.exampleMealPlanOptionVoteList,
		Pagination: &exampleMealPlanOptionVoteList.Pagination,
	}
}

type mealPlanOptionVotesTestSuite struct {
	suite.Suite
	mealPlanOptionVotesBaseSuite
}

func (s *mealPlanOptionVotesTestSuite) TestClient_GetMealPlanOptionVote() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionVoteResponse)
		actual, err := c.GetMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanOptionVote, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVote(s.ctx, "", s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, "", s.exampleMealPlanOptionVote.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option vote ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionVotesTestSuite) TestClient_GetMealPlanOptionVotes() {
	const expectedPath = "/api/v1/meal_plans/%s/events/%s/options/%s/votes"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionVoteListResponse)
		actual, err := c.GetMealPlanOptionVotes(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanOptionVoteList, actual.Data)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVotes(s.ctx, "", s.exampleMealPlanEventID, s.exampleMealPlanOptionID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVotes(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOptionVotes(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOptionVotes(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionVotesTestSuite) TestClient_CreateMealPlanOptionVote() {
	const expectedPath = "/api/v1/meal_plans/%s/events/%s/vote"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleMealPlanID, s.exampleMealPlanEventID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionVoteListResponse)

		actual, err := c.CreateMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanOptionVoteList, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlanOptionVote(s.ctx, "", s.exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealPlanOptionVoteCreationRequestInput{}

		actual, err := c.CreateMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanOptionVoteToMealPlanOptionVoteCreationRequestInput(s.exampleMealPlanOptionVote)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanOptionVoteToMealPlanOptionVoteCreationRequestInput(s.exampleMealPlanOptionVote)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionVotesTestSuite) TestClient_UpdateMealPlanOptionVote() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionVoteResponse)

		err := c.UpdateMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionVote)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealPlanOptionVote(s.ctx, "", s.exampleMealPlanEventID, s.exampleMealPlanOptionVote)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionVote)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionVote)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionVotesTestSuite) TestClient_ArchiveMealPlanOptionVote() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionVoteResponse)

		err := c.ArchiveMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanOptionVote(s.ctx, "", s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, "", s.exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option vote ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMealPlanOptionVote(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOptionID, s.exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
	})
}
