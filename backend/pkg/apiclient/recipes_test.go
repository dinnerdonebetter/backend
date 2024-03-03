package apiclient

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestRecipes(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipesTestSuite))
}

type recipesBaseSuite struct {
	suite.Suite
	ctx                       context.Context
	exampleRecipe             *types.Recipe
	exampleRecipeResponse     *types.APIResponse[*types.Recipe]
	exampleRecipeListResponse *types.APIResponse[[]*types.Recipe]
	exampleRecipeList         []*types.Recipe
}

var _ suite.SetupTestSuite = (*recipesBaseSuite)(nil)

func (s *recipesBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipe = fakes.BuildFakeRecipe()
	s.exampleRecipeResponse = &types.APIResponse[*types.Recipe]{
		Data: s.exampleRecipe,
	}

	exampleList := fakes.BuildFakeRecipeList()
	s.exampleRecipeList = exampleList.Data
	s.exampleRecipeListResponse = &types.APIResponse[[]*types.Recipe]{
		Data:       s.exampleRecipeList,
		Pagination: &exampleList.Pagination,
	}
}

type recipesTestSuite struct {
	suite.Suite
	recipesBaseSuite
}

func (s *recipesTestSuite) TestClient_GetRecipe() {
	const expectedPathFormat = "/api/v1/recipes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipe.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeResponse)
		actual, err := c.GetRecipe(s.ctx, s.exampleRecipe.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipe, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipe(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipe(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipe.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipe(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_GetRecipes() {
	const expectedPath = "/api/v1/recipes"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeListResponse)
		actual, err := c.GetRecipes(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipes(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipes(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_SearchForRecipes() {
	const expectedPath = "/api/v1/recipes/search"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&q=example&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeListResponse)
		actual, err := c.SearchForRecipes(s.ctx, "example", filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchForRecipes(s.ctx, "example", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&q=example&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchForRecipes(s.ctx, "example", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_CreateRecipe() {
	const expectedPath = "/api/v1/recipes"

	s.Run("standard", func() {
		t := s.T()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleInput := converters.ConvertRecipeToRecipeCreationRequestInput(exampleRecipe)
		examplePrepTask := fakes.BuildFakeRecipePrepTask()
		examplePrepTaskInput := converters.ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(exampleRecipe, examplePrepTask)
		examplePrepTaskInput.TaskSteps = []*types.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{
				BelongsToRecipeStepIndex: exampleInput.Steps[0].Index,
				SatisfiesRecipeStep:      false,
			},
		}
		exampleInput.PrepTasks = []*types.RecipePrepTaskWithinRecipeCreationRequestInput{
			examplePrepTaskInput,
		}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeResponse)

		actual, err := c.CreateRecipe(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipe, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipe(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeCreationRequestInput{}

		actual, err := c.CreateRecipe(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeToRecipeCreationRequestInput(s.exampleRecipe)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipe(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeToRecipeCreationRequestInput(s.exampleRecipe)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipe(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_UpdateRecipe() {
	const expectedPathFormat = "/api/v1/recipes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipe.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeResponse)

		err := c.UpdateRecipe(s.ctx, s.exampleRecipe)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipe(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipe(s.ctx, s.exampleRecipe)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipe(s.ctx, s.exampleRecipe)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_ArchiveRecipe() {
	const expectedPathFormat = "/api/v1/recipes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipe.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeResponse)

		err := c.ArchiveRecipe(s.ctx, s.exampleRecipe.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipe(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipe(s.ctx, s.exampleRecipe.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipe(s.ctx, s.exampleRecipe.ID)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_GetRecipeDAG() {
	const expectedPathFormat = "/api/v1/recipes/%s/dag"

	s.Run("standard", func() {
		t := s.T()

		exampleImage, imageBytes := buildArbitraryImagePNGBytes(15)
		require.NotNil(t, exampleImage)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipe.ID)
		c, _ := buildTestClientWithBytesResponse(t, spec, imageBytes)
		actual, err := c.GetRecipeDAG(s.ctx, s.exampleRecipe.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleImage, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeDAG(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeDAG(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipe.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeDAG(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_CloneRecipe() {
	const expectedPathFormat = "/api/v1/recipes/%s/clone"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodPost, "", expectedPathFormat, s.exampleRecipe.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeResponse)
		actual, err := c.CloneRecipe(s.ctx, s.exampleRecipe.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipe, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.CloneRecipe(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CloneRecipe(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodPost, "", expectedPathFormat, s.exampleRecipe.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CloneRecipe(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

// buildArbitraryImage builds an image with a bunch of colors in it.
func buildArbitraryImage(widthAndHeight int) image.Image {
	img := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: widthAndHeight, Y: widthAndHeight}})

	// Set color for each pixel.
	for x := 0; x < widthAndHeight; x++ {
		for y := 0; y < widthAndHeight; y++ {
			img.Set(x, y, color.RGBA{R: uint8(x % math.MaxUint8), G: uint8(y % math.MaxUint8), B: uint8(x + y%math.MaxUint8), A: math.MaxUint8})
		}
	}

	return img
}

// buildArbitraryImagePNGBytes builds an image with a bunch of colors in it.
func buildArbitraryImagePNGBytes(widthAndHeight int) (img image.Image, imgBytes []byte) {
	var b bytes.Buffer

	img = buildArbitraryImage(widthAndHeight)
	if err := png.Encode(&b, img); err != nil {
		panic(err)
	}

	return img, b.Bytes()
}

func (s *recipesTestSuite) TestClient_GetMealPlanTasksForRecipe() {
	const expectedPathFormat = "/api/v1/recipes/%s/prep_steps"

	s.Run("standard", func() {
		t := s.T()

		examplePrepSteps := &types.APIResponse[[]*types.MealPlanTaskDatabaseCreationInput]{
			Data: fakes.BuildFakeMealPlanTaskDatabaseCreationInputs(),
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipe.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, examplePrepSteps)
		actual, err := c.GetMealPlanTasksForRecipe(s.ctx, s.exampleRecipe.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, examplePrepSteps.Data, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanTasksForRecipe(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanTasksForRecipe(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipe.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanTasksForRecipe(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
