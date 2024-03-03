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

func TestRecipeStepProducts(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepProductsTestSuite))
}

type recipeStepProductsBaseSuite struct {
	suite.Suite
	ctx                                  context.Context
	exampleRecipeStepProduct             *types.RecipeStepProduct
	exampleRecipeStepProductResponse     *types.APIResponse[*types.RecipeStepProduct]
	exampleRecipeStepProductListResponse *types.APIResponse[[]*types.RecipeStepProduct]
	exampleRecipeID                      string
	exampleRecipeStepID                  string
	exampleRecipeStepProductList         []*types.RecipeStepProduct
}

var _ suite.SetupTestSuite = (*recipeStepProductsBaseSuite)(nil)

func (s *recipeStepProductsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipeStepID = fakes.BuildFakeID()
	s.exampleRecipeStepProduct = fakes.BuildFakeRecipeStepProduct()
	s.exampleRecipeStepProduct.BelongsToRecipeStep = s.exampleRecipeStepID
	s.exampleRecipeStepProductResponse = &types.APIResponse[*types.RecipeStepProduct]{
		Data: s.exampleRecipeStepProduct,
	}
	exampleList := fakes.BuildFakeRecipeStepProductList()
	s.exampleRecipeStepProductList = exampleList.Data
	s.exampleRecipeStepProductListResponse = &types.APIResponse[[]*types.RecipeStepProduct]{
		Data:       s.exampleRecipeStepProductList,
		Pagination: &exampleList.Pagination,
	}
}

type recipeStepProductsTestSuite struct {
	suite.Suite
	recipeStepProductsBaseSuite
}

func (s *recipeStepProductsTestSuite) TestClient_GetRecipeStepProduct() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/products/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepProductResponse)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepProduct, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepProduct.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step product ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_GetRecipeStepProducts() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/products"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepProductListResponse)
		actual, err := c.GetRecipeStepProducts(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepProductList, actual.Data)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProducts(s.ctx, "", s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProducts(s.ctx, s.exampleRecipeID, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepProducts(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepProducts(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_CreateRecipeStepProduct() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/products"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepProductResponse)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleRecipeStepProduct, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationRequestInput()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, "", s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepProductCreationRequestInput{}

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(s.exampleRecipeStepProduct)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(s.exampleRecipeStepProduct)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_UpdateRecipeStepProduct() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/products/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepProductResponse)

		err := c.UpdateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepProduct)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepProduct(s.ctx, "", s.exampleRecipeStepProduct)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepProduct(s.ctx, s.exampleRecipeID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepProduct)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepProduct)
		assert.Error(t, err)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_ArchiveRecipeStepProduct() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/products/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepProductResponse)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step product ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})
}
