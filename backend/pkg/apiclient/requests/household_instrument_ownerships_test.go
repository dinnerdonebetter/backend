package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetHouseholdInstrumentOwnershipRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleHouseholdInstrumentOwnership.ID)

		actual, err := helper.builder.BuildGetHouseholdInstrumentOwnershipRequest(helper.ctx, exampleHouseholdInstrumentOwnership.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetHouseholdInstrumentOwnershipRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()

		actual, err := helper.builder.BuildGetHouseholdInstrumentOwnershipRequest(helper.ctx, exampleHouseholdInstrumentOwnership.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetHouseholdInstrumentOwnershipsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/instruments"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetHouseholdInstrumentOwnershipsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetHouseholdInstrumentOwnershipsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateHouseholdInstrumentOwnershipRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/households/instruments"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeHouseholdInstrumentOwnershipCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateHouseholdInstrumentOwnershipRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateHouseholdInstrumentOwnershipRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateHouseholdInstrumentOwnershipRequest(helper.ctx, &types.HouseholdInstrumentOwnershipCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeHouseholdInstrumentOwnershipCreationRequestInput()

		actual, err := helper.builder.BuildCreateHouseholdInstrumentOwnershipRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateHouseholdInstrumentOwnershipRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleHouseholdInstrumentOwnership.ID)

		actual, err := helper.builder.BuildUpdateHouseholdInstrumentOwnershipRequest(helper.ctx, exampleHouseholdInstrumentOwnership)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateHouseholdInstrumentOwnershipRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()

		actual, err := helper.builder.BuildUpdateHouseholdInstrumentOwnershipRequest(helper.ctx, exampleHouseholdInstrumentOwnership)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveHouseholdInstrumentOwnershipRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleHouseholdInstrumentOwnership.ID)

		actual, err := helper.builder.BuildArchiveHouseholdInstrumentOwnershipRequest(helper.ctx, exampleHouseholdInstrumentOwnership.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveHouseholdInstrumentOwnershipRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()

		actual, err := helper.builder.BuildArchiveHouseholdInstrumentOwnershipRequest(helper.ctx, exampleHouseholdInstrumentOwnership.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
