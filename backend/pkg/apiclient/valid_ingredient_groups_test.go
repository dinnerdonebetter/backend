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

func TestValidIngredientGroups(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validIngredientGroupsTestSuite))
}

type validIngredientGroupsBaseSuite struct {
	suite.Suite
	ctx                                     context.Context
	exampleValidIngredientGroup             *types.ValidIngredientGroup
	exampleValidIngredientGroupResponse     *types.APIResponse[*types.ValidIngredientGroup]
	exampleValidIngredientGroupListResponse *types.APIResponse[[]*types.ValidIngredientGroup]
	exampleValidIngredientGroupList         []*types.ValidIngredientGroup
}

var _ suite.SetupTestSuite = (*validIngredientGroupsBaseSuite)(nil)

func (s *validIngredientGroupsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidIngredientGroup = fakes.BuildFakeValidIngredientGroup()
	s.exampleValidIngredientGroupResponse = &types.APIResponse[*types.ValidIngredientGroup]{
		Data: s.exampleValidIngredientGroup,
	}

	exampleList := fakes.BuildFakeValidIngredientGroupList()
	s.exampleValidIngredientGroupList = exampleList.Data
	s.exampleValidIngredientGroupListResponse = &types.APIResponse[[]*types.ValidIngredientGroup]{
		Data:       s.exampleValidIngredientGroupList,
		Pagination: &exampleList.Pagination,
	}
}

type validIngredientGroupsTestSuite struct {
	suite.Suite
	validIngredientGroupsBaseSuite
}

func (s *validIngredientGroupsTestSuite) TestClient_GetValidIngredientGroup() {
	const expectedPathFormat = "/api/v1/valid_ingredient_groups/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientGroup.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientGroupResponse)
		actual, err := c.GetValidIngredientGroup(s.ctx, s.exampleValidIngredientGroup.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientGroup, actual)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientGroup(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientGroup(s.ctx, s.exampleValidIngredientGroup.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientGroup.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientGroup(s.ctx, s.exampleValidIngredientGroup.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientGroupsTestSuite) TestClient_GetValidIngredientGroups() {
	const expectedPath = "/api/v1/valid_ingredient_groups"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientGroupListResponse)

		filter := (*types.QueryFilter)(nil)
		actual, err := c.GetValidIngredientGroups(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientGroupList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientGroups(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientGroups(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientGroupsTestSuite) TestClient_SearchValidIngredientGroups() {
	const expectedPath = "/api/v1/valid_ingredient_groups/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientGroupListResponse)
		actual, err := c.SearchValidIngredientGroups(s.ctx, exampleQuery, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientGroupList, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidIngredientGroups(s.ctx, "", 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidIngredientGroups(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidIngredientGroups(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientGroupsTestSuite) TestClient_CreateValidIngredientGroup() {
	const expectedPath = "/api/v1/valid_ingredient_groups"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientGroupResponse)

		actual, err := c.CreateValidIngredientGroup(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidIngredientGroup, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidIngredientGroup(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidIngredientGroupCreationRequestInput{}

		actual, err := c.CreateValidIngredientGroup(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(s.exampleValidIngredientGroup)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidIngredientGroup(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(s.exampleValidIngredientGroup)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidIngredientGroup(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientGroupsTestSuite) TestClient_UpdateValidIngredientGroup() {
	const expectedPathFormat = "/api/v1/valid_ingredient_groups/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidIngredientGroup.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientGroupResponse)

		err := c.UpdateValidIngredientGroup(s.ctx, s.exampleValidIngredientGroup)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidIngredientGroup(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidIngredientGroup(s.ctx, s.exampleValidIngredientGroup)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidIngredientGroup(s.ctx, s.exampleValidIngredientGroup)
		assert.Error(t, err)
	})
}

func (s *validIngredientGroupsTestSuite) TestClient_ArchiveValidIngredientGroup() {
	const expectedPathFormat = "/api/v1/valid_ingredient_groups/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidIngredientGroup.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientGroupResponse)

		err := c.ArchiveValidIngredientGroup(s.ctx, s.exampleValidIngredientGroup.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidIngredientGroup(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidIngredientGroup(s.ctx, s.exampleValidIngredientGroup.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidIngredientGroup(s.ctx, s.exampleValidIngredientGroup.ID)
		assert.Error(t, err)
	})
}
