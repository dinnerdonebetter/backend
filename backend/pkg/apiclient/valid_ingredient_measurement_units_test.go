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

func TestValidIngredientMeasurementUnits(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validIngredientMeasurementUnitsTestSuite))
}

type validIngredientMeasurementUnitsBaseSuite struct {
	suite.Suite
	ctx                                               context.Context
	exampleValidIngredientMeasurementUnit             *types.ValidIngredientMeasurementUnit
	exampleValidIngredientMeasurementUnitResponse     *types.APIResponse[*types.ValidIngredientMeasurementUnit]
	exampleValidIngredientMeasurementUnitListResponse *types.APIResponse[[]*types.ValidIngredientMeasurementUnit]
	exampleValidIngredientMeasurementUnitList         []*types.ValidIngredientMeasurementUnit
}

var _ suite.SetupTestSuite = (*validIngredientMeasurementUnitsBaseSuite)(nil)

func (s *validIngredientMeasurementUnitsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidIngredientMeasurementUnit = fakes.BuildFakeValidIngredientMeasurementUnit()
	s.exampleValidIngredientMeasurementUnitResponse = &types.APIResponse[*types.ValidIngredientMeasurementUnit]{
		Data: s.exampleValidIngredientMeasurementUnit,
	}

	exampleList := fakes.BuildFakeValidIngredientMeasurementUnitList()
	s.exampleValidIngredientMeasurementUnitList = exampleList.Data
	s.exampleValidIngredientMeasurementUnitListResponse = &types.APIResponse[[]*types.ValidIngredientMeasurementUnit]{
		Data:       s.exampleValidIngredientMeasurementUnitList,
		Pagination: &exampleList.Pagination,
	}
}

type validIngredientMeasurementUnitsTestSuite struct {
	suite.Suite
	validIngredientMeasurementUnitsBaseSuite
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_GetValidIngredientMeasurementUnit() {
	const expectedPathFormat = "/api/v1/valid_ingredient_measurement_units/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientMeasurementUnitResponse)
		actual, err := c.GetValidIngredientMeasurementUnit(s.ctx, s.exampleValidIngredientMeasurementUnit.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientMeasurementUnit, actual)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientMeasurementUnit(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientMeasurementUnit(s.ctx, s.exampleValidIngredientMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientMeasurementUnit.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientMeasurementUnit(s.ctx, s.exampleValidIngredientMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_GetValidIngredientMeasurementUnits() {
	const expectedPath = "/api/v1/valid_ingredient_measurement_units"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientMeasurementUnitListResponse)

		filter := (*types.QueryFilter)(nil)
		actual, err := c.GetValidIngredientMeasurementUnits(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientMeasurementUnitList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientMeasurementUnits(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientMeasurementUnits(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_GetValidIngredientMeasurementUnitsForIngredient() {
	const expectedPath = "/api/v1/valid_ingredient_measurement_units/by_ingredient/%s"

	exampleValidIngredient := fakes.BuildFakeValidIngredient()

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientMeasurementUnitListResponse)

		filter := (*types.QueryFilter)(nil)
		actual, err := c.GetValidIngredientMeasurementUnitsForIngredient(s.ctx, exampleValidIngredient.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientMeasurementUnitList, actual.Data)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientMeasurementUnitsForIngredient(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientMeasurementUnitsForIngredient(s.ctx, exampleValidIngredient.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidIngredient.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientMeasurementUnitsForIngredient(s.ctx, exampleValidIngredient.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_GetValidIngredientMeasurementUnitsForMeasurementUnit() {
	const expectedPath = "/api/v1/valid_ingredient_measurement_units/by_measurement_unit/%s"

	exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientMeasurementUnitListResponse)

		filter := (*types.QueryFilter)(nil)
		actual, err := c.GetValidIngredientMeasurementUnitsForMeasurementUnit(s.ctx, exampleValidMeasurementUnit.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientMeasurementUnitList, actual.Data)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientMeasurementUnitsForMeasurementUnit(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientMeasurementUnitsForMeasurementUnit(s.ctx, exampleValidMeasurementUnit.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidMeasurementUnit.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientMeasurementUnitsForMeasurementUnit(s.ctx, exampleValidMeasurementUnit.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_CreateValidIngredientMeasurementUnit() {
	const expectedPath = "/api/v1/valid_ingredient_measurement_units"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientMeasurementUnitCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientMeasurementUnitResponse)

		actual, err := c.CreateValidIngredientMeasurementUnit(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientMeasurementUnit, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidIngredientMeasurementUnit(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidIngredientMeasurementUnitCreationRequestInput{}

		actual, err := c.CreateValidIngredientMeasurementUnit(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(s.exampleValidIngredientMeasurementUnit)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidIngredientMeasurementUnit(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(s.exampleValidIngredientMeasurementUnit)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidIngredientMeasurementUnit(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_UpdateValidIngredientMeasurementUnit() {
	const expectedPathFormat = "/api/v1/valid_ingredient_measurement_units/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidIngredientMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientMeasurementUnitResponse)

		err := c.UpdateValidIngredientMeasurementUnit(s.ctx, s.exampleValidIngredientMeasurementUnit)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidIngredientMeasurementUnit(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidIngredientMeasurementUnit(s.ctx, s.exampleValidIngredientMeasurementUnit)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidIngredientMeasurementUnit(s.ctx, s.exampleValidIngredientMeasurementUnit)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_ArchiveValidIngredientMeasurementUnit() {
	const expectedPathFormat = "/api/v1/valid_ingredient_measurement_units/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidIngredientMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientMeasurementUnitResponse)

		err := c.ArchiveValidIngredientMeasurementUnit(s.ctx, s.exampleValidIngredientMeasurementUnit.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidIngredientMeasurementUnit(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidIngredientMeasurementUnit(s.ctx, s.exampleValidIngredientMeasurementUnit.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidIngredientMeasurementUnit(s.ctx, s.exampleValidIngredientMeasurementUnit.ID)
		assert.Error(t, err)
	})
}
