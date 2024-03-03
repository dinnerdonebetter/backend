package converters

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertHouseholdInvitationCreationInputToHouseholdInvitationDatabaseCreationInput creates a HouseholdInvitationDatabaseCreationInput from a HouseholdInvitationCreationRequestInput.
func ConvertHouseholdInvitationCreationInputToHouseholdInvitationDatabaseCreationInput(input *types.HouseholdInvitationCreationRequestInput) *types.HouseholdInvitationDatabaseCreationInput {
	x := &types.HouseholdInvitationDatabaseCreationInput{
		ToEmail: input.ToEmail,
		ToName:  input.ToName,
	}

	if input.ExpiresAt != nil {
		x.ExpiresAt = *input.ExpiresAt
	}

	return x
}

// ConvertHouseholdInvitationToHouseholdInvitationCreationInput builds a faked HouseholdInvitationCreationRequestInput.
func ConvertHouseholdInvitationToHouseholdInvitationCreationInput(householdInvitation *types.HouseholdInvitation) *types.HouseholdInvitationCreationRequestInput {
	return &types.HouseholdInvitationCreationRequestInput{
		Note:      householdInvitation.Note,
		ToName:    householdInvitation.ToName,
		ToEmail:   householdInvitation.ToEmail,
		ExpiresAt: &householdInvitation.ExpiresAt,
	}
}

// ConvertHouseholdInvitationToHouseholdInvitationUpdateInput builds a faked HouseholdInvitationUpdateRequestInput.
func ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(householdInvitation *types.HouseholdInvitation) *types.HouseholdInvitationUpdateRequestInput {
	return &types.HouseholdInvitationUpdateRequestInput{
		Token: householdInvitation.Token,
		Note:  householdInvitation.Note,
	}
}

// ConvertHouseholdInvitationToHouseholdInvitationDatabaseCreationInput builds a faked HouseholdInvitationCreationRequestInput.
func ConvertHouseholdInvitationToHouseholdInvitationDatabaseCreationInput(householdInvitation *types.HouseholdInvitation) *types.HouseholdInvitationDatabaseCreationInput {
	return &types.HouseholdInvitationDatabaseCreationInput{
		ID:                     householdInvitation.ID,
		FromUser:               householdInvitation.FromUser.ID,
		ToUser:                 householdInvitation.ToUser,
		ToName:                 householdInvitation.ToName,
		Note:                   householdInvitation.Note,
		ToEmail:                householdInvitation.ToEmail,
		Token:                  householdInvitation.Token,
		DestinationHouseholdID: householdInvitation.DestinationHousehold.ID,
		ExpiresAt:              householdInvitation.ExpiresAt,
	}
}
