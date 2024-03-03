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

func TestHouseholds(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(householdsTestSuite))
}

type householdsTestSuite struct {
	suite.Suite
	ctx                          context.Context
	exampleUser                  *types.User
	exampleHousehold             *types.Household
	exampleHouseholdResponse     *types.APIResponse[*types.Household]
	exampleHouseholdListResponse *types.APIResponse[[]*types.Household]
	exampleHouseholdList         []*types.Household
}

var _ suite.SetupTestSuite = (*householdsTestSuite)(nil)

func (s *householdsTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleUser = fakes.BuildFakeUser()
	s.exampleHousehold = fakes.BuildFakeHousehold()
	s.exampleHousehold.WebhookEncryptionKey = ""
	s.exampleHousehold.BelongsToUser = s.exampleUser.ID
	exampleHouseholdList := fakes.BuildFakeHouseholdList()
	for i := range exampleHouseholdList.Data {
		exampleHouseholdList.Data[i].WebhookEncryptionKey = ""
	}

	s.exampleHouseholdList = exampleHouseholdList.Data
	s.exampleHouseholdListResponse = &types.APIResponse[[]*types.Household]{
		Data:       exampleHouseholdList.Data,
		Pagination: &exampleHouseholdList.Pagination,
	}
	s.exampleHouseholdResponse = &types.APIResponse[*types.Household]{
		Data: s.exampleHousehold,
	}
}

func (s *householdsTestSuite) TestClient_SwitchActiveHousehold() {
	const expectedPath = "/api/v1/users/household/select"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdResponse)
		c.authMethod = cookieAuthMethod

		assert.NoError(t, c.SwitchActiveHousehold(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		c.authMethod = cookieAuthMethod

		assert.Error(t, c.SwitchActiveHousehold(s.ctx, ""))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		c.authMethod = cookieAuthMethod

		assert.Error(t, c.SwitchActiveHousehold(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		c.authMethod = cookieAuthMethod

		assert.Error(t, c.SwitchActiveHousehold(s.ctx, s.exampleHousehold.ID))
	})
}

func (s *householdsTestSuite) TestClient_GetCurrentHousehold() {
	const expectedPathFormat = "/api/v1/households/current"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdResponse)

		actual, err := c.GetCurrentHousehold(s.ctx)

		for i := range actual.Members {
			actual.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt = s.exampleHousehold.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt
		}

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHousehold, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetCurrentHousehold(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat)

		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetCurrentHousehold(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_GetHousehold() {
	const expectedPathFormat = "/api/v1/households/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdResponse)

		actual, err := c.GetHousehold(s.ctx, s.exampleHousehold.ID)

		for i := range actual.Members {
			actual.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt = s.exampleHousehold.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt
		}

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHousehold, actual)
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetHousehold(s.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHousehold(s.ctx, s.exampleHousehold.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleHousehold.ID)

		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHousehold(s.ctx, s.exampleHousehold.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_GetHouseholds() {
	const expectedPath = "/api/v1/households"

	spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdListResponse)
		actual, err := c.GetHouseholds(s.ctx, filter)

		for i, household := range actual.Data {
			for j := range household.Members {
				actual.Data[i].Members[j].BelongsToUser.TwoFactorSecretVerifiedAt = s.exampleHouseholdList[i].Members[j].BelongsToUser.TwoFactorSecretVerifiedAt
			}
		}

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHouseholdList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHouseholds(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHouseholds(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_CreateHousehold() {
	const expectedPath = "/api/v1/households"

	spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""
		exampleInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(s.exampleHousehold)

		c := buildTestClientWithRequestBodyValidation(t, spec, exampleInput, exampleInput, s.exampleHouseholdResponse)
		actual, err := c.CreateHousehold(s.ctx, exampleInput)

		for i := range actual.Members {
			actual.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt = s.exampleHousehold.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt
		}

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHousehold, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateHousehold(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		exampleInput := &types.HouseholdCreationRequestInput{}
		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateHousehold(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(s.exampleHousehold)
		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateHousehold(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""
		exampleInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(s.exampleHousehold)

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateHousehold(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_UpdateHousehold() {
	const expectedPathFormat = "/api/v1/households/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdResponse)

		assert.NoError(t, c.UpdateHousehold(s.ctx, s.exampleHousehold), "no error should be returned")
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.UpdateHousehold(s.ctx, nil), "error should be returned")
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateHousehold(s.ctx, s.exampleHousehold)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.UpdateHousehold(s.ctx, s.exampleHousehold), "error should be returned")
	})
}

func (s *householdsTestSuite) TestClient_ArchiveHousehold() {
	const expectedPathFormat = "/api/v1/households/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdResponse)

		assert.NoError(t, c.ArchiveHousehold(s.ctx, s.exampleHousehold.ID), "no error should be returned")
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.ArchiveHousehold(s.ctx, ""), "no error should be returned")
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		assert.Error(t, c.ArchiveHousehold(s.ctx, s.exampleHousehold.ID), "error should be returned")
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.ArchiveHousehold(s.ctx, s.exampleHousehold.ID), "no error should be returned")
	})
}

func (s *householdsTestSuite) TestClient_InviteUserToHousehold() {
	const expectedPathFormat = "/api/v1/households/%s/invite"

	s.Run("standard", func() {
		t := s.T()

		invitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdID := fakes.BuildFakeID()
		invitation.FromUser.TwoFactorSecret = ""
		invitation.DestinationHousehold.WebhookEncryptionKey = ""
		invitationResponse := &types.APIResponse[*types.HouseholdInvitation]{
			Data: invitation,
		}

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationCreationInput(invitation)
		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, exampleHouseholdID)
		c, _ := buildTestClientWithJSONResponse(t, spec, invitationResponse)

		householdInvite, err := c.InviteUserToHousehold(s.ctx, exampleHouseholdID, exampleInput)
		assert.Equal(t, invitation, householdInvite)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		exampleHouseholdID := fakes.BuildFakeID()

		householdInvite, err := c.InviteUserToHousehold(s.ctx, exampleHouseholdID, nil)
		assert.Nil(t, householdInvite)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		exampleHouseholdID := fakes.BuildFakeID()

		householdInvite, err := c.InviteUserToHousehold(s.ctx, exampleHouseholdID, &types.HouseholdInvitationCreationRequestInput{})
		assert.Nil(t, householdInvite)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		exampleHouseholdID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()

		householdInvite, err := c.InviteUserToHousehold(s.ctx, exampleHouseholdID, exampleInput)
		assert.Nil(t, householdInvite)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		exampleHouseholdID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()

		householdInvite, err := c.InviteUserToHousehold(s.ctx, exampleHouseholdID, exampleInput)
		assert.Nil(t, householdInvite)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_MarkAsDefault() {
	const expectedPathFormat = "/api/v1/households/%s/default"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodPost, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdResponse)

		assert.NoError(t, c.MarkAsDefault(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.MarkAsDefault(s.ctx, ""))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.MarkAsDefault(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.MarkAsDefault(s.ctx, s.exampleHousehold.ID))
	})
}

func (s *householdsTestSuite) TestClient_RemoveUserFromHousehold() {
	const expectedPathFormat = "/api/v1/households/%s/members/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleHousehold.ID, s.exampleUser.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdResponse)

		assert.NoError(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, "", s.exampleUser.ID))
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, ""))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID))
	})
}

func (s *householdsTestSuite) TestClient_ModifyMemberPermissions() {
	const expectedPathFormat = "/api/v1/households/%s/members/%s/permissions"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, s.exampleHousehold.ID, s.exampleUser.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdResponse)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.NoError(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, exampleInput))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, "", s.exampleUser.ID, exampleInput))
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, "", exampleInput))
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, nil))
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ModifyUserPermissionsInput{}

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, exampleInput))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, exampleInput))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, exampleInput))
	})
}

func (s *householdsTestSuite) TestClient_TransferHouseholdOwnership() {
	const expectedPathFormat = "/api/v1/households/%s/transfer"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdResponse)
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		assert.NoError(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, exampleInput))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, "", exampleInput))
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, nil))
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.HouseholdOwnershipTransferInput{}

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, exampleInput))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, exampleInput))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, exampleInput))
	})
}
