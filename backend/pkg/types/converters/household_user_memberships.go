package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertHouseholdUserMembershipToHouseholdUserMembershipDatabaseCreationInput builds a faked HouseholdUserMembershipCreationRequestInput.
func ConvertHouseholdUserMembershipToHouseholdUserMembershipDatabaseCreationInput(household *types.HouseholdUserMembership) *types.HouseholdUserMembershipDatabaseCreationInput {
	return &types.HouseholdUserMembershipDatabaseCreationInput{
		ID:            identifiers.New(),
		Reason:        "",
		UserID:        household.BelongsToUser,
		HouseholdID:   household.ID,
		HouseholdRole: household.HouseholdRole,
	}
}
