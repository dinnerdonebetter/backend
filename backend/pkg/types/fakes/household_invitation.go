package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeHouseholdInvitation builds a faked HouseholdInvitation.
func BuildFakeHouseholdInvitation() *types.HouseholdInvitation {
	return &types.HouseholdInvitation{
		FromUser:             *BuildFakeUser(),
		ToEmail:              fake.Email(),
		ToName:               buildUniqueString(),
		ToUser:               func(s string) *string { return &s }(buildUniqueString()),
		Note:                 buildUniqueString(),
		StatusNote:           buildUniqueString(),
		Token:                fake.UUID(),
		DestinationHousehold: *BuildFakeHousehold(),
		ID:                   BuildFakeID(),
		ExpiresAt:            BuildFakeTime(),
		Status:               string(types.PendingHouseholdInvitationStatus),
		CreatedAt:            BuildFakeTime(),
	}
}

// BuildFakeHouseholdInvitationList builds a faked HouseholdInvitationList.
func BuildFakeHouseholdInvitationList() *types.QueryFilteredResult[types.HouseholdInvitation] {
	var examples []*types.HouseholdInvitation
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeHouseholdInvitation())
	}

	return &types.QueryFilteredResult[types.HouseholdInvitation]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeHouseholdInvitationCreationRequestInput builds a faked HouseholdInvitationCreationRequestInput from a webhook.
func BuildFakeHouseholdInvitationCreationRequestInput() *types.HouseholdInvitationCreationRequestInput {
	invitation := BuildFakeHouseholdInvitation()
	return converters.ConvertHouseholdInvitationToHouseholdInvitationCreationInput(invitation)
}
