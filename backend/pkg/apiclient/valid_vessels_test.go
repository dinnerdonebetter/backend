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

func TestValidVessels(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validVesselsTestSuite))
}

type validVesselsBaseSuite struct {
	suite.Suite

	ctx                        context.Context
	exampleValidVessel         *types.ValidVessel
	exampleValidVesselResponse *types.APIResponse[*types.ValidVessel]
}

var _ suite.SetupTestSuite = (*validVesselsBaseSuite)(nil)

func (s *validVesselsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidVessel = fakes.BuildFakeValidVessel()
	s.exampleValidVesselResponse = &types.APIResponse[*types.ValidVessel]{
		Data: s.exampleValidVessel,
	}
}

type validVesselsTestSuite struct {
	suite.Suite
	validVesselsBaseSuite
}

func (s *validVesselsTestSuite) TestClient_GetValidVessel() {
	const expectedPathFormat = "/api/v1/valid_vessels/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidVessel.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidVesselResponse)
		actual, err := c.GetValidVessel(s.ctx, s.exampleValidVessel.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidVessel, actual)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidVessel(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidVessel(s.ctx, s.exampleValidVessel.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidVessel.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidVessel(s.ctx, s.exampleValidVessel.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validVesselsTestSuite) TestClient_GetRandomValidVessel() {
	const expectedPath = "/api/v1/valid_vessels/random"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidVesselResponse)
		actual, err := c.GetRandomValidVessel(s.ctx)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidVessel, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRandomValidVessel(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRandomValidVessel(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validVesselsTestSuite) TestClient_GetValidVessels() {
	const expectedPath = "/api/v1/valid_vessels"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidVesselList := fakes.BuildFakeValidVesselList()
		exampleValidVesselListAPIResponse := &types.APIResponse[[]*types.ValidVessel]{
			Data:       exampleValidVesselList.Data,
			Pagination: &exampleValidVesselList.Pagination,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidVesselListAPIResponse)
		actual, err := c.GetValidVessels(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidVesselList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidVessels(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidVessels(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validVesselsTestSuite) TestClient_SearchValidVessels() {
	const expectedPath = "/api/v1/valid_vessels/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		exampleValidVesselList := fakes.BuildFakeValidVesselList()
		exampleValidVesselListAPIResponse := &types.APIResponse[[]*types.ValidVessel]{
			Data: exampleValidVesselList.Data,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidVesselListAPIResponse)
		actual, err := c.SearchValidVessels(s.ctx, exampleQuery, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidVesselList.Data, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidVessels(s.ctx, "", 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidVessels(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidVessels(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validVesselsTestSuite) TestClient_CreateValidVessel() {
	const expectedPath = "/api/v1/valid_vessels"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidVesselCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidVesselResponse)

		actual, err := c.CreateValidVessel(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidVessel, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidVessel(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidVesselCreationRequestInput{}

		actual, err := c.CreateValidVessel(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidVesselToValidVesselCreationRequestInput(s.exampleValidVessel)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidVessel(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidVesselToValidVesselCreationRequestInput(s.exampleValidVessel)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidVessel(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validVesselsTestSuite) TestClient_UpdateValidVessel() {
	const expectedPathFormat = "/api/v1/valid_vessels/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidVessel.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidVesselResponse)

		err := c.UpdateValidVessel(s.ctx, s.exampleValidVessel)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidVessel(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidVessel(s.ctx, s.exampleValidVessel)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidVessel(s.ctx, s.exampleValidVessel)
		assert.Error(t, err)
	})
}

func (s *validVesselsTestSuite) TestClient_ArchiveValidVessel() {
	const expectedPathFormat = "/api/v1/valid_vessels/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidVessel.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidVesselResponse)

		err := c.ArchiveValidVessel(s.ctx, s.exampleValidVessel.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidVessel(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidVessel(s.ctx, s.exampleValidVessel.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidVessel(s.ctx, s.exampleValidVessel.ID)
		assert.Error(t, err)
	})
}
