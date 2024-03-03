package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetValidPreparationInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidPreparationInstrument.ID)

		actual, err := helper.builder.BuildGetValidPreparationInstrumentRequest(helper.ctx, exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidPreparationInstrumentRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		actual, err := helper.builder.BuildGetValidPreparationInstrumentRequest(helper.ctx, exampleValidPreparationInstrument.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidPreparationInstrumentsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_instruments"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetValidPreparationInstrumentsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationInstrumentsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidPreparationInstrumentsForPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_instruments/by_preparation/%s"

	examplePreparation := fakes.BuildFakeValidPreparation()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, examplePreparation.ID)

		actual, err := helper.builder.BuildGetValidPreparationInstrumentsForPreparationRequest(helper.ctx, examplePreparation.ID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationInstrumentsForPreparationRequest(helper.ctx, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationInstrumentsForPreparationRequest(helper.ctx, examplePreparation.ID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidPreparationInstrumentsForInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_instruments/by_instrument/%s"

	examplePreparation := fakes.BuildFakeValidPreparation()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, examplePreparation.ID)

		actual, err := helper.builder.BuildGetValidPreparationInstrumentsForInstrumentRequest(helper.ctx, examplePreparation.ID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationInstrumentsForInstrumentRequest(helper.ctx, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationInstrumentsForInstrumentRequest(helper.ctx, examplePreparation.ID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidPreparationInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_preparation_instruments"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidPreparationInstrumentRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidPreparationInstrumentRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidPreparationInstrumentRequest(helper.ctx, &types.ValidPreparationInstrumentCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()

		actual, err := helper.builder.BuildCreateValidPreparationInstrumentRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidPreparationInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidPreparationInstrument.ID)

		actual, err := helper.builder.BuildUpdateValidPreparationInstrumentRequest(helper.ctx, exampleValidPreparationInstrument)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidPreparationInstrumentRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		actual, err := helper.builder.BuildUpdateValidPreparationInstrumentRequest(helper.ctx, exampleValidPreparationInstrument)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidPreparationInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidPreparationInstrument.ID)

		actual, err := helper.builder.BuildArchiveValidPreparationInstrumentRequest(helper.ctx, exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidPreparationInstrumentRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		actual, err := helper.builder.BuildArchiveValidPreparationInstrumentRequest(helper.ctx, exampleValidPreparationInstrument.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
