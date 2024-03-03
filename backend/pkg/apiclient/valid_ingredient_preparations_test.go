package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestValidIngredientPreparations(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validIngredientPreparationsTestSuite))
}

type validIngredientPreparationsBaseSuite struct {
	suite.Suite
	ctx                                           context.Context
	exampleValidIngredientPreparation             *types.ValidIngredientPreparation
	exampleValidIngredientPreparationResponse     *types.APIResponse[*types.ValidIngredientPreparation]
	exampleValidIngredientPreparationListResponse *types.APIResponse[[]*types.ValidIngredientPreparation]
	exampleValidIngredientPreparationList         []*types.ValidIngredientPreparation
}

var _ suite.SetupTestSuite = (*validIngredientPreparationsBaseSuite)(nil)

func (s *validIngredientPreparationsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidIngredientPreparation = fakes.BuildFakeValidIngredientPreparation()
	s.exampleValidIngredientPreparationResponse = &types.APIResponse[*types.ValidIngredientPreparation]{
		Data: s.exampleValidIngredientPreparation,
	}

	exampleList := fakes.BuildFakeValidIngredientPreparationList()
	s.exampleValidIngredientPreparationList = exampleList.Data
	s.exampleValidIngredientPreparationListResponse = &types.APIResponse[[]*types.ValidIngredientPreparation]{
		Data:       s.exampleValidIngredientPreparationList,
		Pagination: &exampleList.Pagination,
	}
}

type validIngredientPreparationsTestSuite struct {
	suite.Suite
	validIngredientPreparationsBaseSuite
}

func (s *validIngredientPreparationsTestSuite) TestClient_GetValidIngredientPreparation() {
	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparationResponse)
		actual, err := c.GetValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientPreparation, actual)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientPreparation(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientPreparation.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_GetValidIngredientPreparations() {
	const expectedPath = "/api/v1/valid_ingredient_preparations"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparationListResponse)
		actual, err := c.GetValidIngredientPreparations(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientPreparationList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparations(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparations(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_GetValidIngredientPreparationsForIngredient() {
	const expectedPath = "/api/v1/valid_ingredient_preparations/by_ingredient/%s"

	exampleValidIngredient := fakes.BuildFakeValidIngredient()

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparationListResponse)
		actual, err := c.GetValidIngredientPreparationsForIngredient(s.ctx, exampleValidIngredient.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientPreparationList, actual.Data)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientPreparationsForIngredient(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparationsForIngredient(s.ctx, exampleValidIngredient.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidIngredient.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparationsForIngredient(s.ctx, exampleValidIngredient.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_GetValidIngredientPreparationsForPreparation() {
	const expectedPath = "/api/v1/valid_ingredient_preparations/by_preparation/%s"

	exampleValidPreparation := fakes.BuildFakeValidPreparation()

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparationListResponse)
		actual, err := c.GetValidIngredientPreparationsForPreparation(s.ctx, exampleValidPreparation.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientPreparationList, actual.Data)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientPreparationsForPreparation(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparationsForPreparation(s.ctx, exampleValidPreparation.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidPreparation.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparationsForPreparation(s.ctx, exampleValidPreparation.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_GetValidIngredientPreparationsForPreparationAndIngredientName() {
	const expectedPath = "/api/v1/valid_ingredients/by_preparation/%s"

	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	exampleQuery := "blah"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", exampleQuery), expectedPath, exampleValidPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparationListResponse)
		actual, err := c.GetValidIngredientPreparationsForPreparationAndIngredientName(s.ctx, exampleValidPreparation.ID, exampleQuery, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientPreparationList, actual.Data)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientPreparationsForPreparationAndIngredientName(s.ctx, "", exampleQuery, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparationsForPreparationAndIngredientName(s.ctx, exampleValidPreparation.ID, exampleQuery, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", exampleQuery), expectedPath, exampleValidPreparation.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparationsForPreparationAndIngredientName(s.ctx, exampleValidPreparation.ID, exampleQuery, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_CreateValidIngredientPreparation() {
	const expectedPath = "/api/v1/valid_ingredient_preparations"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparationResponse)

		actual, err := c.CreateValidIngredientPreparation(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientPreparation, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidIngredientPreparation(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidIngredientPreparationCreationRequestInput{}

		actual, err := c.CreateValidIngredientPreparation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(s.exampleValidIngredientPreparation)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidIngredientPreparation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(s.exampleValidIngredientPreparation)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidIngredientPreparation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_UpdateValidIngredientPreparation() {
	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidIngredientPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparationResponse)

		err := c.UpdateValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidIngredientPreparation(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_ArchiveValidIngredientPreparation() {
	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidIngredientPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparationResponse)

		err := c.ArchiveValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidIngredientPreparation(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)
		assert.Error(t, err)
	})
}
