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

func TestRecipeSteps(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepsTestSuite))
}

type recipeStepsBaseSuite struct {
	suite.Suite
	ctx                           context.Context
	exampleRecipeStep             *types.RecipeStep
	exampleRecipeStepResponse     *types.APIResponse[*types.RecipeStep]
	exampleRecipeStepListResponse *types.APIResponse[[]*types.RecipeStep]
	exampleRecipeID               string
	exampleRecipeStepList         []*types.RecipeStep
}

var _ suite.SetupTestSuite = (*recipeStepsBaseSuite)(nil)

func (s *recipeStepsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipeStep = fakes.BuildFakeRecipeStep()
	s.exampleRecipeStep.BelongsToRecipe = s.exampleRecipeID
	s.exampleRecipeStepResponse = &types.APIResponse[*types.RecipeStep]{
		Data: s.exampleRecipeStep,
	}

	exampleList := fakes.BuildFakeRecipeStepList()
	s.exampleRecipeStepList = exampleList.Data
	s.exampleRecipeStepListResponse = &types.APIResponse[[]*types.RecipeStep]{
		Data:       s.exampleRecipeStepList,
		Pagination: &exampleList.Pagination,
	}
}

type recipeStepsTestSuite struct {
	suite.Suite
	recipeStepsBaseSuite
}

func (s *recipeStepsTestSuite) TestClient_GetRecipeStep() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStep.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepResponse)
		actual, err := c.GetRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStep, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStep(s.ctx, "", s.exampleRecipeStep.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStep(s.ctx, s.exampleRecipeID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStep.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepsTestSuite) TestClient_GetRecipeSteps() {
	const expectedPath = "/api/v1/recipes/%s/steps"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleRecipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepListResponse)
		actual, err := c.GetRecipeSteps(s.ctx, s.exampleRecipeID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepList, actual.Data)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeSteps(s.ctx, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeSteps(s.ctx, s.exampleRecipeID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleRecipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeSteps(s.ctx, s.exampleRecipeID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepsTestSuite) TestClient_CreateRecipeStep() {
	const expectedPath = "/api/v1/recipes/%s/steps"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepResponse)

		actual, err := c.CreateRecipeStep(s.ctx, s.exampleRecipeID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStep, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStep(s.ctx, s.exampleRecipeID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepCreationRequestInput{}

		actual, err := c.CreateRecipeStep(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepToRecipeStepCreationRequestInput(s.exampleRecipeStep)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStep(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepToRecipeStepCreationRequestInput(s.exampleRecipeStep)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStep(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepsTestSuite) TestClient_UpdateRecipeStep() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStep.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepResponse)

		err := c.UpdateRecipeStep(s.ctx, s.exampleRecipeStep)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStep(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeStep(s.ctx, s.exampleRecipeStep)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeStep(s.ctx, s.exampleRecipeStep)
		assert.Error(t, err)
	})
}

func (s *recipeStepsTestSuite) TestClient_ArchiveRecipeStep() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStep.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepResponse)

		err := c.ArchiveRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStep(s.ctx, "", s.exampleRecipeStep.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStep(s.ctx, s.exampleRecipeID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)
		assert.Error(t, err)
	})
}
