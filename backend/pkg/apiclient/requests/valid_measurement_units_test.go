package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetValidMeasurementUnitRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidMeasurementUnit.ID)

		actual, err := helper.builder.BuildGetValidMeasurementUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid measurement ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidMeasurementUnitRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		actual, err := helper.builder.BuildGetValidMeasurementUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidMeasurementUnitsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_units"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetValidMeasurementUnitsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidMeasurementUnitsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildSearchValidMeasurementUnitsRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_measurement_units/search"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"
		spec := newRequestSpec(true, http.MethodGet, "limit=50&q=whatever", expectedPath)

		actual, err := helper.builder.BuildSearchValidMeasurementUnitsRequest(helper.ctx, exampleQuery, *limit)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		actual, err := helper.builder.BuildSearchValidMeasurementUnitsRequest(helper.ctx, exampleQuery, *limit)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildSearchValidMeasurementUnitsByIngredientIDRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_measurement_units/by_ingredient/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientID := fakes.BuildFakeID()
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidIngredientID)

		actual, err := helper.builder.BuildSearchValidMeasurementUnitsByIngredientIDRequest(helper.ctx, exampleValidIngredientID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildSearchValidMeasurementUnitsByIngredientIDRequest(helper.ctx, exampleValidIngredientID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidMeasurementUnitRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_measurement_units"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidMeasurementUnitRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidMeasurementUnitRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidMeasurementUnitRequest(helper.ctx, &types.ValidMeasurementUnitCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()

		actual, err := helper.builder.BuildCreateValidMeasurementUnitRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidMeasurementUnitRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidMeasurementUnit.ID)

		actual, err := helper.builder.BuildUpdateValidMeasurementUnitRequest(helper.ctx, exampleValidMeasurementUnit)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidMeasurementUnitRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		actual, err := helper.builder.BuildUpdateValidMeasurementUnitRequest(helper.ctx, exampleValidMeasurementUnit)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidMeasurementUnitRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidMeasurementUnit.ID)

		actual, err := helper.builder.BuildArchiveValidMeasurementUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid measurement ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidMeasurementUnitRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		actual, err := helper.builder.BuildArchiveValidMeasurementUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
