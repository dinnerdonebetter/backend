package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetValidPreparationVesselRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidPreparationVessel.ID)

		actual, err := helper.builder.BuildGetValidPreparationVesselRequest(helper.ctx, exampleValidPreparationVessel.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidPreparationVesselRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		actual, err := helper.builder.BuildGetValidPreparationVesselRequest(helper.ctx, exampleValidPreparationVessel.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidPreparationVesselsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetValidPreparationVesselsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationVesselsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidPreparationVesselsForPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels/by_preparation/%s"

	examplePreparation := fakes.BuildFakeValidPreparation()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, examplePreparation.ID)

		actual, err := helper.builder.BuildGetValidPreparationVesselsForPreparationRequest(helper.ctx, examplePreparation.ID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationVesselsForPreparationRequest(helper.ctx, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationVesselsForPreparationRequest(helper.ctx, examplePreparation.ID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidPreparationVesselsForVesselRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels/by_vessel/%s"

	examplePreparation := fakes.BuildFakeValidPreparation()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, examplePreparation.ID)

		actual, err := helper.builder.BuildGetValidPreparationVesselsForVesselRequest(helper.ctx, examplePreparation.ID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationVesselsForVesselRequest(helper.ctx, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationVesselsForVesselRequest(helper.ctx, examplePreparation.ID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidPreparationVesselRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_preparation_vessels"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidPreparationVesselCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidPreparationVesselRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidPreparationVesselRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidPreparationVesselRequest(helper.ctx, &types.ValidPreparationVesselCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidPreparationVesselCreationRequestInput()

		actual, err := helper.builder.BuildCreateValidPreparationVesselRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidPreparationVesselRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidPreparationVessel.ID)

		actual, err := helper.builder.BuildUpdateValidPreparationVesselRequest(helper.ctx, exampleValidPreparationVessel)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidPreparationVesselRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		actual, err := helper.builder.BuildUpdateValidPreparationVesselRequest(helper.ctx, exampleValidPreparationVessel)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidPreparationVesselRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_vessels/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidPreparationVessel.ID)

		actual, err := helper.builder.BuildArchiveValidPreparationVesselRequest(helper.ctx, exampleValidPreparationVessel.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidPreparationVesselRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		actual, err := helper.builder.BuildArchiveValidPreparationVesselRequest(helper.ctx, exampleValidPreparationVessel.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
