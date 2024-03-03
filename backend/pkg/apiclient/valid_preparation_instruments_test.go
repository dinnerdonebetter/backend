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

func TestValidPreparationInstruments(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validPreparationInstrumentsTestSuite))
}

type validPreparationInstrumentsBaseSuite struct {
	suite.Suite
	ctx                                           context.Context
	exampleValidPreparationInstrument             *types.ValidPreparationInstrument
	exampleValidPreparationInstrumentResponse     *types.APIResponse[*types.ValidPreparationInstrument]
	exampleValidPreparationInstrumentListResponse *types.APIResponse[[]*types.ValidPreparationInstrument]
	exampleValidPreparationInstrumentList         []*types.ValidPreparationInstrument
}

var _ suite.SetupTestSuite = (*validPreparationInstrumentsBaseSuite)(nil)

func (s *validPreparationInstrumentsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidPreparationInstrument = fakes.BuildFakeValidPreparationInstrument()
	s.exampleValidPreparationInstrumentResponse = &types.APIResponse[*types.ValidPreparationInstrument]{
		Data: s.exampleValidPreparationInstrument,
	}

	exampleList := fakes.BuildFakeValidPreparationInstrumentList()

	s.exampleValidPreparationInstrumentList = exampleList.Data
	s.exampleValidPreparationInstrumentListResponse = &types.APIResponse[[]*types.ValidPreparationInstrument]{
		Data:       s.exampleValidPreparationInstrumentList,
		Pagination: &exampleList.Pagination,
	}
}

type validPreparationInstrumentsTestSuite struct {
	suite.Suite
	validPreparationInstrumentsBaseSuite
}

func (s *validPreparationInstrumentsTestSuite) TestClient_GetValidPreparationInstrument() {
	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidPreparationInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrumentResponse)
		actual, err := c.GetValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidPreparationInstrument, actual)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidPreparationInstrument(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidPreparationInstrument.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_GetValidPreparationInstruments() {
	const expectedPath = "/api/v1/valid_preparation_instruments"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrumentListResponse)
		actual, err := c.GetValidPreparationInstruments(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidPreparationInstrumentList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationInstruments(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationInstruments(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_GetValidPreparationInstrumentsForPreparation() {
	const expectedPath = "/api/v1/valid_preparation_instruments/by_preparation/%s"

	exampleValidPreparation := fakes.BuildFakeValidPreparation()

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrumentListResponse)
		actual, err := c.GetValidPreparationInstrumentsForPreparation(s.ctx, exampleValidPreparation.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidPreparationInstrumentList, actual.Data)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidPreparationInstrumentsForPreparation(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationInstrumentsForPreparation(s.ctx, exampleValidPreparation.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidPreparation.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationInstrumentsForPreparation(s.ctx, exampleValidPreparation.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_GetValidPreparationInstrumentsForInstrument() {
	const expectedPath = "/api/v1/valid_preparation_instruments/by_instrument/%s"

	exampleValidInstrument := fakes.BuildFakeValidInstrument()

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrumentListResponse)
		actual, err := c.GetValidPreparationInstrumentsForInstrument(s.ctx, exampleValidInstrument.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidPreparationInstrumentList, actual.Data)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidPreparationInstrumentsForInstrument(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationInstrumentsForInstrument(s.ctx, exampleValidInstrument.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidInstrument.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationInstrumentsForInstrument(s.ctx, exampleValidInstrument.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_CreateValidPreparationInstrument() {
	const expectedPath = "/api/v1/valid_preparation_instruments"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrumentResponse)

		actual, err := c.CreateValidPreparationInstrument(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidPreparationInstrument, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidPreparationInstrument(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidPreparationInstrumentCreationRequestInput{}

		actual, err := c.CreateValidPreparationInstrument(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(s.exampleValidPreparationInstrument)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidPreparationInstrument(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(s.exampleValidPreparationInstrument)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidPreparationInstrument(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_UpdateValidPreparationInstrument() {
	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidPreparationInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrumentResponse)

		err := c.UpdateValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidPreparationInstrument(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_ArchiveValidPreparationInstrument() {
	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidPreparationInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrumentResponse)

		err := c.ArchiveValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidPreparationInstrument(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
	})
}
