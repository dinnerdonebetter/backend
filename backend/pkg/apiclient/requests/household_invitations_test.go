package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetHouseholdInvitationRequest(T *testing.T) {
	T.Parallel()

	expectedPathFormat := "/api/v1/households/%s/invitations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		spec := newRequestSpec(false, http.MethodGet, "", expectedPathFormat, exampleHouseholdInvitation.DestinationHousehold.ID, exampleHouseholdInvitation.ID)

		actual, err := helper.builder.BuildGetHouseholdInvitationRequest(helper.ctx, exampleHouseholdInvitation.DestinationHousehold.ID, exampleHouseholdInvitation.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		actual, err := helper.builder.BuildGetHouseholdInvitationRequest(helper.ctx, "", exampleHouseholdInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid invitation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		actual, err := helper.builder.BuildGetHouseholdInvitationRequest(helper.ctx, exampleHouseholdInvitation.DestinationHousehold.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		actual, err := helper.builder.BuildGetHouseholdInvitationRequest(helper.ctx, exampleHouseholdInvitation.DestinationHousehold.ID, exampleHouseholdInvitation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetPendingHouseholdInvitationsFromUserRequest(T *testing.T) {
	T.Parallel()

	expectedPathFormat := "/api/v1/household_invitations/sent"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		filter := types.DefaultQueryFilter()

		spec := newRequestSpec(false, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetPendingHouseholdInvitationsFromUserRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		filter := types.DefaultQueryFilter()

		actual, err := helper.builder.BuildGetPendingHouseholdInvitationsFromUserRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetPendingHouseholdInvitationsForUserRequest(T *testing.T) {
	T.Parallel()

	expectedPathFormat := "/api/v1/household_invitations/received"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		filter := types.DefaultQueryFilter()

		spec := newRequestSpec(false, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetPendingHouseholdInvitationsForUserRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		filter := types.DefaultQueryFilter()

		actual, err := helper.builder.BuildGetPendingHouseholdInvitationsForUserRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildAcceptHouseholdInvitationRequest(T *testing.T) {
	T.Parallel()

	expectedPathFormat := "/api/v1/household_invitations/%s/accept"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleHouseholdInvitation.ID)

		actual, err := helper.builder.BuildAcceptHouseholdInvitationRequest(helper.ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, t.Name())
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		actual, err := helper.builder.BuildAcceptHouseholdInvitationRequest(helper.ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, t.Name())
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCancelHouseholdInvitationRequest(T *testing.T) {
	T.Parallel()

	expectedPathFormat := "/api/v1/household_invitations/%s/cancel"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleHouseholdInvitation.ID)

		actual, err := helper.builder.BuildCancelHouseholdInvitationRequest(helper.ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, t.Name())
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		actual, err := helper.builder.BuildCancelHouseholdInvitationRequest(helper.ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, t.Name())
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildRejectHouseholdInvitationRequest(T *testing.T) {
	T.Parallel()

	expectedPathFormat := "/api/v1/household_invitations/%s/reject"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleHouseholdInvitation.ID)

		actual, err := helper.builder.BuildRejectHouseholdInvitationRequest(helper.ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, t.Name())
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		actual, err := helper.builder.BuildRejectHouseholdInvitationRequest(helper.ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, t.Name())
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
