package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetValidMeasurementUnitConversionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidMeasurementUnitConversion.ID)

		actual, err := helper.builder.BuildGetValidMeasurementUnitConversionRequest(helper.ctx, exampleValidMeasurementUnitConversion.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidMeasurementUnitConversionRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		actual, err := helper.builder.BuildGetValidMeasurementUnitConversionRequest(helper.ctx, exampleValidMeasurementUnitConversion.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidMeasurementUnitConversionsFromUnitRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/from_unit/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidMeasurementUnit.ID)

		actual, err := helper.builder.BuildGetValidMeasurementUnitConversionsFromUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidMeasurementUnitConversionsFromUnitRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		actual, err := helper.builder.BuildGetValidMeasurementUnitConversionsFromUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidMeasurementUnitConversionsToUnitRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/to_unit/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidMeasurementUnit.ID)

		actual, err := helper.builder.BuildGetValidMeasurementUnitConversionsToUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidMeasurementUnitConversionsToUnitRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		actual, err := helper.builder.BuildGetValidMeasurementUnitConversionsToUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidMeasurementUnitConversionRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_measurement_conversions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidMeasurementUnitConversionRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidMeasurementUnitConversionRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidMeasurementUnitConversionRequest(helper.ctx, &types.ValidMeasurementUnitConversionCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()

		actual, err := helper.builder.BuildCreateValidMeasurementUnitConversionRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidMeasurementUnitConversionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidMeasurementUnitConversion.ID)

		actual, err := helper.builder.BuildUpdateValidMeasurementUnitConversionRequest(helper.ctx, exampleValidMeasurementUnitConversion)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidMeasurementUnitConversionRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		actual, err := helper.builder.BuildUpdateValidMeasurementUnitConversionRequest(helper.ctx, exampleValidMeasurementUnitConversion)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidMeasurementUnitConversionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidMeasurementUnitConversion.ID)

		actual, err := helper.builder.BuildArchiveValidMeasurementUnitConversionRequest(helper.ctx, exampleValidMeasurementUnitConversion.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidMeasurementUnitConversionRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		actual, err := helper.builder.BuildArchiveValidMeasurementUnitConversionRequest(helper.ctx, exampleValidMeasurementUnitConversion.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
