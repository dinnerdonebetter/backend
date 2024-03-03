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

func TestRecipeRatings(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeRatingsTestSuite))
}

type recipeRatingsBaseSuite struct {
	suite.Suite
	ctx                             context.Context
	exampleMeal                     *types.Meal
	exampleRecipeRating             *types.RecipeRating
	exampleRecipeRatingResponse     *types.APIResponse[*types.RecipeRating]
	exampleRecipeRatingListResponse *types.APIResponse[[]*types.RecipeRating]
	exampleRecipeRatingList         []*types.RecipeRating
}

var _ suite.SetupTestSuite = (*recipeRatingsBaseSuite)(nil)

func (s *recipeRatingsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMeal = fakes.BuildFakeMeal()
	s.exampleRecipeRating = fakes.BuildFakeRecipeRating()
	s.exampleRecipeRatingResponse = &types.APIResponse[*types.RecipeRating]{
		Data: s.exampleRecipeRating,
	}
	exampleList := fakes.BuildFakeRecipeRatingList()
	s.exampleRecipeRatingList = exampleList.Data
	s.exampleRecipeRatingListResponse = &types.APIResponse[[]*types.RecipeRating]{
		Data:       s.exampleRecipeRatingList,
		Pagination: &exampleList.Pagination,
	}
}

type recipeRatingsTestSuite struct {
	suite.Suite
	recipeRatingsBaseSuite
}

func (s *recipeRatingsTestSuite) TestClient_GetRecipeRating() {
	const expectedPathFormat = "/api/v1/recipes/%s/ratings/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMeal.ID, s.exampleRecipeRating.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeRatingResponse)
		actual, err := c.GetRecipeRating(s.ctx, s.exampleMeal.ID, s.exampleRecipeRating.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeRating, actual)
	})

	s.Run("with invalid valid instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeRating(s.ctx, s.exampleMeal.ID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeRating(s.ctx, s.exampleMeal.ID, s.exampleRecipeRating.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMeal.ID, s.exampleRecipeRating.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeRating(s.ctx, s.exampleMeal.ID, s.exampleRecipeRating.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeRatingsTestSuite) TestClient_GetRecipeRatings() {
	const expectedPath = "/api/v1/recipes/%s/ratings"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleMeal.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeRatingListResponse)
		actual, err := c.GetRecipeRatings(s.ctx, s.exampleMeal.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeRatingList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeRatings(s.ctx, s.exampleMeal.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleMeal.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeRatings(s.ctx, s.exampleMeal.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeRatingsTestSuite) TestClient_CreateRecipeRating() {
	const expectedPath = "/api/v1/recipes/%s/ratings"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeRatingCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleMeal.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeRatingResponse)

		actual, err := c.CreateRecipeRating(s.ctx, s.exampleMeal.ID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeRating, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeRating(s.ctx, s.exampleMeal.ID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeRatingCreationRequestInput{}

		actual, err := c.CreateRecipeRating(s.ctx, s.exampleMeal.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeRatingToRecipeRatingCreationRequestInput(s.exampleRecipeRating)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeRating(s.ctx, s.exampleMeal.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeRatingToRecipeRatingCreationRequestInput(s.exampleRecipeRating)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeRating(s.ctx, s.exampleMeal.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeRatingsTestSuite) TestClient_UpdateRecipeRating() {
	const expectedPathFormat = "/api/v1/recipes/%s/ratings/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeRating.RecipeID, s.exampleRecipeRating.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeRatingResponse)

		err := c.UpdateRecipeRating(s.ctx, s.exampleRecipeRating)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeRating(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeRating(s.ctx, s.exampleRecipeRating)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeRating(s.ctx, s.exampleRecipeRating)
		assert.Error(t, err)
	})
}

func (s *recipeRatingsTestSuite) TestClient_ArchiveRecipeRating() {
	const expectedPathFormat = "/api/v1/recipes/%s/ratings/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMeal.ID, s.exampleRecipeRating.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeRatingResponse)

		err := c.ArchiveRecipeRating(s.ctx, s.exampleMeal.ID, s.exampleRecipeRating.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeRating(s.ctx, s.exampleMeal.ID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeRating(s.ctx, s.exampleMeal.ID, s.exampleRecipeRating.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeRating(s.ctx, s.exampleMeal.ID, s.exampleRecipeRating.ID)
		assert.Error(t, err)
	})
}
