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

func TestValidMeasurementUnits(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validMeasurementUnitsTestSuite))
}

type validMeasurementUnitsBaseSuite struct {
	suite.Suite
	ctx                                     context.Context
	exampleValidMeasurementUnit             *types.ValidMeasurementUnit
	exampleValidMeasurementUnitResponse     *types.APIResponse[*types.ValidMeasurementUnit]
	exampleValidMeasurementUnitListResponse *types.APIResponse[[]*types.ValidMeasurementUnit]
	exampleValidMeasurementUnitList         []*types.ValidMeasurementUnit
}

var _ suite.SetupTestSuite = (*validMeasurementUnitsBaseSuite)(nil)

func (s *validMeasurementUnitsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
	s.exampleValidMeasurementUnitResponse = &types.APIResponse[*types.ValidMeasurementUnit]{
		Data: s.exampleValidMeasurementUnit,
	}

	exampleList := fakes.BuildFakeValidMeasurementUnitList()
	s.exampleValidMeasurementUnitList = exampleList.Data
	s.exampleValidMeasurementUnitListResponse = &types.APIResponse[[]*types.ValidMeasurementUnit]{
		Data:       s.exampleValidMeasurementUnitList,
		Pagination: &exampleList.Pagination,
	}
}

type validMeasurementUnitsTestSuite struct {
	suite.Suite
	validMeasurementUnitsBaseSuite
}

func (s *validMeasurementUnitsTestSuite) TestClient_GetValidMeasurementUnit() {
	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitResponse)
		actual, err := c.GetValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidMeasurementUnit, actual)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidMeasurementUnit(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_GetValidMeasurementUnits() {
	const expectedPath = "/api/v1/valid_measurement_units"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitListResponse)
		actual, err := c.GetValidMeasurementUnits(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidMeasurementUnitList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementUnits(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementUnits(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_SearchValidMeasurementUnits() {
	const expectedPath = "/api/v1/valid_measurement_units/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitListResponse)
		actual, err := c.SearchValidMeasurementUnits(s.ctx, exampleQuery, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidMeasurementUnitList, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidMeasurementUnits(s.ctx, "", 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidMeasurementUnits(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidMeasurementUnits(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_SearchValidMeasurementUnitsByIngredientID() {
	const expectedPath = "/api/v1/valid_measurement_units/by_ingredient/%s"

	s.Run("standard", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidIngredientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitListResponse)

		actual, err := c.SearchValidMeasurementUnitsByIngredientID(s.ctx, exampleValidIngredientID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidMeasurementUnitList, actual.Data)
	})

	s.Run("with empty query", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()
		c, _ := buildSimpleTestClient(t)

		actual, err := c.SearchValidMeasurementUnitsByIngredientID(s.ctx, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientID := fakes.BuildFakeID()
		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidMeasurementUnitsByIngredientID(s.ctx, exampleValidIngredientID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidIngredientID)
		c := buildTestClientWithInvalidResponse(t, spec)

		actual, err := c.SearchValidMeasurementUnitsByIngredientID(s.ctx, exampleValidIngredientID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_CreateValidMeasurementUnit() {
	const expectedPath = "/api/v1/valid_measurement_units"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitResponse)

		actual, err := c.CreateValidMeasurementUnit(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidMeasurementUnit, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidMeasurementUnit(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidMeasurementUnitCreationRequestInput{}

		actual, err := c.CreateValidMeasurementUnit(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(s.exampleValidMeasurementUnit)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidMeasurementUnit(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(s.exampleValidMeasurementUnit)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidMeasurementUnit(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_UpdateValidMeasurementUnit() {
	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitResponse)

		err := c.UpdateValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidMeasurementUnit(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_ArchiveValidMeasurementUnit() {
	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitResponse)

		err := c.ArchiveValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidMeasurementUnit(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)
		assert.Error(t, err)
	})
}
