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

func TestHouseholdInstrumentOwnerships(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(householdInstrumentOwnershipsTestSuite))
}

type householdInstrumentOwnershipsBaseSuite struct {
	suite.Suite
	ctx                                             context.Context
	exampleHouseholdInstrumentOwnership             *types.HouseholdInstrumentOwnership
	exampleHouseholdInstrumentOwnershipResponse     *types.APIResponse[*types.HouseholdInstrumentOwnership]
	exampleHouseholdInstrumentOwnershipListResponse *types.APIResponse[[]*types.HouseholdInstrumentOwnership]
	exampleHouseholdInstrumentOwnershipList         []*types.HouseholdInstrumentOwnership
}

var _ suite.SetupTestSuite = (*householdInstrumentOwnershipsBaseSuite)(nil)

func (s *householdInstrumentOwnershipsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleHouseholdInstrumentOwnership = fakes.BuildFakeHouseholdInstrumentOwnership()
	s.exampleHouseholdInstrumentOwnershipResponse = &types.APIResponse[*types.HouseholdInstrumentOwnership]{
		Data: s.exampleHouseholdInstrumentOwnership,
	}
	exampleHouseholdInstrumentOwnershipList := fakes.BuildFakeHouseholdInstrumentOwnershipList()
	s.exampleHouseholdInstrumentOwnershipList = exampleHouseholdInstrumentOwnershipList.Data
	s.exampleHouseholdInstrumentOwnershipListResponse = &types.APIResponse[[]*types.HouseholdInstrumentOwnership]{
		Data:       s.exampleHouseholdInstrumentOwnershipList,
		Pagination: &exampleHouseholdInstrumentOwnershipList.Pagination,
	}
}

type householdInstrumentOwnershipsTestSuite struct {
	suite.Suite
	householdInstrumentOwnershipsBaseSuite
}

func (s *householdInstrumentOwnershipsTestSuite) TestClient_GetHouseholdInstrumentOwnership() {
	const expectedPathFormat = "/api/v1/households/instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleHouseholdInstrumentOwnership.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInstrumentOwnershipResponse)
		actual, err := c.GetHouseholdInstrumentOwnership(s.ctx, s.exampleHouseholdInstrumentOwnership.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHouseholdInstrumentOwnership, actual)
	})

	s.Run("with invalid valid instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetHouseholdInstrumentOwnership(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHouseholdInstrumentOwnership(s.ctx, s.exampleHouseholdInstrumentOwnership.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleHouseholdInstrumentOwnership.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHouseholdInstrumentOwnership(s.ctx, s.exampleHouseholdInstrumentOwnership.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdInstrumentOwnershipsTestSuite) TestClient_GetHouseholdInstrumentOwnerships() {
	const expectedPath = "/api/v1/households/instruments"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInstrumentOwnershipListResponse)
		actual, err := c.GetHouseholdInstrumentOwnerships(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHouseholdInstrumentOwnershipList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHouseholdInstrumentOwnerships(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHouseholdInstrumentOwnerships(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdInstrumentOwnershipsTestSuite) TestClient_CreateHouseholdInstrumentOwnership() {
	const expectedPath = "/api/v1/households/instruments"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeHouseholdInstrumentOwnershipCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInstrumentOwnershipResponse)

		actual, err := c.CreateHouseholdInstrumentOwnership(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHouseholdInstrumentOwnership, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateHouseholdInstrumentOwnership(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.HouseholdInstrumentOwnershipCreationRequestInput{}

		actual, err := c.CreateHouseholdInstrumentOwnership(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipCreationRequestInput(s.exampleHouseholdInstrumentOwnership)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateHouseholdInstrumentOwnership(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipCreationRequestInput(s.exampleHouseholdInstrumentOwnership)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateHouseholdInstrumentOwnership(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdInstrumentOwnershipsTestSuite) TestClient_UpdateHouseholdInstrumentOwnership() {
	const expectedPathFormat = "/api/v1/households/instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleHouseholdInstrumentOwnership.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInstrumentOwnershipResponse)

		err := c.UpdateHouseholdInstrumentOwnership(s.ctx, s.exampleHouseholdInstrumentOwnership)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateHouseholdInstrumentOwnership(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateHouseholdInstrumentOwnership(s.ctx, s.exampleHouseholdInstrumentOwnership)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateHouseholdInstrumentOwnership(s.ctx, s.exampleHouseholdInstrumentOwnership)
		assert.Error(t, err)
	})
}

func (s *householdInstrumentOwnershipsTestSuite) TestClient_ArchiveHouseholdInstrumentOwnership() {
	const expectedPathFormat = "/api/v1/households/instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleHouseholdInstrumentOwnership.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInstrumentOwnershipResponse)

		err := c.ArchiveHouseholdInstrumentOwnership(s.ctx, s.exampleHouseholdInstrumentOwnership.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveHouseholdInstrumentOwnership(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveHouseholdInstrumentOwnership(s.ctx, s.exampleHouseholdInstrumentOwnership.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveHouseholdInstrumentOwnership(s.ctx, s.exampleHouseholdInstrumentOwnership.ID)
		assert.Error(t, err)
	})
}
