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

func TestValidMeasurementUnitConversions(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validMeasurementUnitConversionsTestSuite))
}

type validMeasurementUnitConversionsBaseSuite struct {
	suite.Suite
	ctx                                               context.Context
	exampleValidMeasurementUnit                       *types.ValidMeasurementUnit
	exampleValidMeasurementUnitConversion             *types.ValidMeasurementUnitConversion
	exampleValidMeasurementUnitConversionResponse     *types.APIResponse[*types.ValidMeasurementUnitConversion]
	exampleValidMeasurementUnitConversionListResponse *types.APIResponse[[]*types.ValidMeasurementUnitConversion]
	exampleValidMeasurementUnitConversionList         []*types.ValidMeasurementUnitConversion
}

var _ suite.SetupTestSuite = (*validMeasurementUnitConversionsBaseSuite)(nil)

func (s *validMeasurementUnitConversionsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
	s.exampleValidMeasurementUnitConversion = fakes.BuildFakeValidMeasurementUnitConversion()
	s.exampleValidMeasurementUnitConversionList = []*types.ValidMeasurementUnitConversion{
		s.exampleValidMeasurementUnitConversion,
	}

	s.exampleValidMeasurementUnitConversionResponse = &types.APIResponse[*types.ValidMeasurementUnitConversion]{
		Data: s.exampleValidMeasurementUnitConversion,
	}

	s.exampleValidMeasurementUnitConversionListResponse = &types.APIResponse[[]*types.ValidMeasurementUnitConversion]{
		Data: s.exampleValidMeasurementUnitConversionList,
	}
}

type validMeasurementUnitConversionsTestSuite struct {
	suite.Suite
	validMeasurementUnitConversionsBaseSuite
}

func (s *validMeasurementUnitConversionsTestSuite) TestClient_GetValidMeasurementUnitConversion() {
	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnitConversion.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitConversionResponse)
		actual, err := c.GetValidMeasurementUnitConversion(s.ctx, s.exampleValidMeasurementUnitConversion.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidMeasurementUnitConversion, actual)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidMeasurementUnitConversion(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementUnitConversion(s.ctx, s.exampleValidMeasurementUnitConversion.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnitConversion.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementUnitConversion(s.ctx, s.exampleValidMeasurementUnitConversion.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitConversionsTestSuite) TestClient_GetValidMeasurementUnitConversionsFromUnit() {
	const expectedPathFormat = "/api/v1/valid_measurement_conversions/from_unit/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitConversionListResponse)
		actual, err := c.GetValidMeasurementUnitConversionsFromUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidMeasurementUnitConversionList, actual)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidMeasurementUnitConversionsFromUnit(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementUnitConversionsFromUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementUnitConversionsFromUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitConversionsTestSuite) TestClient_CreateValidMeasurementUnitConversion() {
	const expectedPath = "/api/v1/valid_measurement_conversions"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitConversionResponse)

		actual, err := c.CreateValidMeasurementUnitConversion(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidMeasurementUnitConversion, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidMeasurementUnitConversion(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidMeasurementUnitConversionCreationRequestInput{}

		actual, err := c.CreateValidMeasurementUnitConversion(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionCreationRequestInput(s.exampleValidMeasurementUnitConversion)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidMeasurementUnitConversion(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionCreationRequestInput(s.exampleValidMeasurementUnitConversion)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidMeasurementUnitConversion(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitConversionsTestSuite) TestClient_UpdateValidMeasurementUnitConversion() {
	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidMeasurementUnitConversion.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitConversionResponse)

		err := c.UpdateValidMeasurementUnitConversion(s.ctx, s.exampleValidMeasurementUnitConversion)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidMeasurementUnitConversion(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidMeasurementUnitConversion(s.ctx, s.exampleValidMeasurementUnitConversion)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidMeasurementUnitConversion(s.ctx, s.exampleValidMeasurementUnitConversion)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitConversionsTestSuite) TestClient_ArchiveValidMeasurementUnitConversion() {
	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidMeasurementUnitConversion.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnitConversionResponse)

		err := c.ArchiveValidMeasurementUnitConversion(s.ctx, s.exampleValidMeasurementUnitConversion.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidMeasurementUnitConversion(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidMeasurementUnitConversion(s.ctx, s.exampleValidMeasurementUnitConversion.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidMeasurementUnitConversion(s.ctx, s.exampleValidMeasurementUnitConversion.ID)
		assert.Error(t, err)
	})
}
