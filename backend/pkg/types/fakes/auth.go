package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeSessionContextData builds a faked HouseholdUserMembershipCreationRequestInput.
func BuildFakeSessionContextData() *types.SessionContextData {
	return &types.SessionContextData{
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{},
		Requester: types.RequesterInfo{
			ServicePermissions:       nil,
			AccountStatus:            string(types.GoodStandingUserAccountStatus),
			AccountStatusExplanation: "fake",
			UserID:                   BuildFakeID(),
			EmailAddress:             fake.Email(),
			Username:                 buildUniqueString(),
		},
		ActiveHouseholdID: BuildFakeID(),
	}
}

// BuildFakeHouseholdUserMembershipCreationRequestInput builds a faked HouseholdUserMembershipCreationRequestInput.
func BuildFakeHouseholdUserMembershipCreationRequestInput() *types.HouseholdUserMembershipCreationRequestInput {
	return &types.HouseholdUserMembershipCreationRequestInput{
		Reason: fake.Sentence(10),
		UserID: BuildFakeID(),
	}
}

// BuildFakeHouseholdUserMembershipDatabaseCreationInput builds a faked HouseholdUserMembershipCreationRequestInput.
func BuildFakeHouseholdUserMembershipDatabaseCreationInput() *types.HouseholdUserMembershipDatabaseCreationInput {
	input := BuildFakeHouseholdUserMembershipCreationRequestInput()

	return converters.ConvertHouseholdUserMembershipCreationRequestInputToHouseholdUserMembershipDatabaseCreationInput(input)
}

// BuildFakeUserPermissionModificationInput builds a faked ModifyUserPermissionsInput.
func BuildFakeUserPermissionModificationInput() *types.ModifyUserPermissionsInput {
	return &types.ModifyUserPermissionsInput{
		Reason:  fake.Sentence(10),
		NewRole: authorization.HouseholdMemberRole.String(),
	}
}

// BuildFakeTransferHouseholdOwnershipInput builds a faked HouseholdOwnershipTransferInput.
func BuildFakeTransferHouseholdOwnershipInput() *types.HouseholdOwnershipTransferInput {
	return &types.HouseholdOwnershipTransferInput{
		Reason:       fake.Sentence(10),
		CurrentOwner: fake.UUID(),
		NewOwner:     fake.UUID(),
	}
}

// BuildFakeChangeActiveHouseholdInput builds a faked ChangeActiveHouseholdInput.
func BuildFakeChangeActiveHouseholdInput() *types.ChangeActiveHouseholdInput {
	return &types.ChangeActiveHouseholdInput{
		HouseholdID: fake.UUID(),
	}
}
