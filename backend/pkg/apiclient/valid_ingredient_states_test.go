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

func TestValidIngredientStates(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validIngredientStatesTestSuite))
}

type validIngredientStatesBaseSuite struct {
	suite.Suite
	ctx                                     context.Context
	exampleValidIngredientState             *types.ValidIngredientState
	exampleValidIngredientStateResponse     *types.APIResponse[*types.ValidIngredientState]
	exampleValidIngredientStateListResponse *types.APIResponse[[]*types.ValidIngredientState]
	exampleValidIngredientStateList         []*types.ValidIngredientState
}

var _ suite.SetupTestSuite = (*validIngredientStatesBaseSuite)(nil)

func (s *validIngredientStatesBaseSuite) SetupTest() {
	s.ctx = context.Background()
	exampleList := fakes.BuildFakeValidIngredientStateList()
	s.exampleValidIngredientState = fakes.BuildFakeValidIngredientState()
	s.exampleValidIngredientStateList = exampleList.Data
	s.exampleValidIngredientStateResponse = &types.APIResponse[*types.ValidIngredientState]{
		Data: s.exampleValidIngredientState,
	}
	s.exampleValidIngredientStateListResponse = &types.APIResponse[[]*types.ValidIngredientState]{
		Data:       s.exampleValidIngredientStateList,
		Pagination: &exampleList.Pagination,
	}
}

type validIngredientStatesTestSuite struct {
	suite.Suite
	validIngredientStatesBaseSuite
}

func (s *validIngredientStatesTestSuite) TestClient_GetValidIngredientState() {
	const expectedPathFormat = "/api/v1/valid_ingredient_states/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientState.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientStateResponse)
		actual, err := c.GetValidIngredientState(s.ctx, s.exampleValidIngredientState.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientState, actual)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientState(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientState(s.ctx, s.exampleValidIngredientState.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientState.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientState(s.ctx, s.exampleValidIngredientState.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientStatesTestSuite) TestClient_GetValidIngredientStates() {
	const expectedPath = "/api/v1/valid_ingredient_states"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientStateListResponse)
		actual, err := c.GetValidIngredientStates(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientStateList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientStates(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientStates(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientStatesTestSuite) TestClient_SearchValidIngredientStates() {
	const expectedPath = "/api/v1/valid_ingredient_states/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientStateListResponse)
		actual, err := c.SearchValidIngredientStates(s.ctx, exampleQuery, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientStateList, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidIngredientStates(s.ctx, "", 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidIngredientStates(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidIngredientStates(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientStatesTestSuite) TestClient_CreateValidIngredientState() {
	const expectedPath = "/api/v1/valid_ingredient_states"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientStateCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientStateResponse)

		actual, err := c.CreateValidIngredientState(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidIngredientState, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidIngredientState(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidIngredientStateCreationRequestInput{}

		actual, err := c.CreateValidIngredientState(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientStateToValidIngredientStateCreationRequestInput(s.exampleValidIngredientState)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidIngredientState(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientStateToValidIngredientStateCreationRequestInput(s.exampleValidIngredientState)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidIngredientState(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientStatesTestSuite) TestClient_UpdateValidIngredientState() {
	const expectedPathFormat = "/api/v1/valid_ingredient_states/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidIngredientState.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientStateResponse)

		err := c.UpdateValidIngredientState(s.ctx, s.exampleValidIngredientState)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidIngredientState(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidIngredientState(s.ctx, s.exampleValidIngredientState)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidIngredientState(s.ctx, s.exampleValidIngredientState)
		assert.Error(t, err)
	})
}

func (s *validIngredientStatesTestSuite) TestClient_ArchiveValidIngredientState() {
	const expectedPathFormat = "/api/v1/valid_ingredient_states/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidIngredientState.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientStateResponse)

		err := c.ArchiveValidIngredientState(s.ctx, s.exampleValidIngredientState.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidIngredientState(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidIngredientState(s.ctx, s.exampleValidIngredientState.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidIngredientState(s.ctx, s.exampleValidIngredientState.ID)
		assert.Error(t, err)
	})
}
