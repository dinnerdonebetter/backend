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

func TestValidInstruments(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validInstrumentsTestSuite))
}

type validInstrumentsBaseSuite struct {
	suite.Suite
	ctx                                context.Context
	exampleValidInstrument             *types.ValidInstrument
	exampleValidInstrumentListResponse types.APIResponse[[]*types.ValidInstrument]
	exampleValidInstrumentResponse     types.APIResponse[*types.ValidInstrument]
	exampleValidInstrumentList         []*types.ValidInstrument
}

var _ suite.SetupTestSuite = (*validInstrumentsBaseSuite)(nil)

func (s *validInstrumentsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidInstrument = fakes.BuildFakeValidInstrument()
	s.exampleValidInstrumentResponse = types.APIResponse[*types.ValidInstrument]{
		Data: s.exampleValidInstrument,
	}

	exampleList := fakes.BuildFakeValidInstrumentList()
	s.exampleValidInstrumentList = exampleList.Data
	s.exampleValidInstrumentListResponse = types.APIResponse[[]*types.ValidInstrument]{
		Data:       s.exampleValidInstrumentList,
		Pagination: &exampleList.Pagination,
	}
}

type validInstrumentsTestSuite struct {
	suite.Suite
	validInstrumentsBaseSuite
}

func (s *validInstrumentsTestSuite) TestClient_GetValidInstrument() {
	const expectedPathFormat = "/api/v1/valid_instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidInstrumentResponse)
		actual, err := c.GetValidInstrument(s.ctx, s.exampleValidInstrument.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidInstrument, actual)
	})

	s.Run("with invalid valid instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidInstrument(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidInstrument(s.ctx, s.exampleValidInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidInstrument.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidInstrument(s.ctx, s.exampleValidInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validInstrumentsTestSuite) TestClient_GetRandomValidInstrument() {
	const expectedPath = "/api/v1/valid_instruments/random"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidInstrumentResponse)
		actual, err := c.GetRandomValidInstrument(s.ctx)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidInstrument, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRandomValidInstrument(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRandomValidInstrument(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validInstrumentsTestSuite) TestClient_GetValidInstruments() {
	const expectedPath = "/api/v1/valid_instruments"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidInstrumentListResponse)
		actual, err := c.GetValidInstruments(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidInstrumentList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidInstruments(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidInstruments(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validInstrumentsTestSuite) TestClient_SearchValidInstruments() {
	const expectedPath = "/api/v1/valid_instruments/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidInstrumentListResponse)
		actual, err := c.SearchValidInstruments(s.ctx, exampleQuery, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidInstrumentList, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidInstruments(s.ctx, "", 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidInstruments(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidInstruments(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validInstrumentsTestSuite) TestClient_CreateValidInstrument() {
	const expectedPath = "/api/v1/valid_instruments"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidInstrumentCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidInstrumentResponse)

		actual, err := c.CreateValidInstrument(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidInstrument, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidInstrument(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidInstrumentCreationRequestInput{}

		actual, err := c.CreateValidInstrument(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(s.exampleValidInstrument)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidInstrument(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(s.exampleValidInstrument)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidInstrument(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validInstrumentsTestSuite) TestClient_UpdateValidInstrument() {
	const expectedPathFormat = "/api/v1/valid_instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidInstrumentResponse)

		err := c.UpdateValidInstrument(s.ctx, s.exampleValidInstrument)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidInstrument(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidInstrument(s.ctx, s.exampleValidInstrument)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidInstrument(s.ctx, s.exampleValidInstrument)
		assert.Error(t, err)
	})
}

func (s *validInstrumentsTestSuite) TestClient_ArchiveValidInstrument() {
	const expectedPathFormat = "/api/v1/valid_instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidInstrumentResponse)

		err := c.ArchiveValidInstrument(s.ctx, s.exampleValidInstrument.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidInstrument(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidInstrument(s.ctx, s.exampleValidInstrument.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidInstrument(s.ctx, s.exampleValidInstrument.ID)
		assert.Error(t, err)
	})
}
