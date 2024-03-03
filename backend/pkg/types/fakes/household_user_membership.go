package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/pkg/types"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeHouseholdUserMembership builds a faked HouseholdUserMembership.
func BuildFakeHouseholdUserMembership() *types.HouseholdUserMembership {
	return &types.HouseholdUserMembership{
		ID:                 BuildFakeID(),
		BelongsToUser:      BuildFakeID(),
		BelongsToHousehold: fake.UUID(),
		HouseholdRole:      authorization.HouseholdMemberRole.String(),
		CreatedAt:          BuildFakeTime(),
		ArchivedAt:         nil,
	}
}

// BuildFakeHouseholdUserMembershipWithUser builds a faked HouseholdUserMembershipWithUser.
func BuildFakeHouseholdUserMembershipWithUser() *types.HouseholdUserMembershipWithUser {
	u := BuildFakeUser()
	u.TwoFactorSecret = ""

	return &types.HouseholdUserMembershipWithUser{
		ID:                 BuildFakeID(),
		BelongsToUser:      u,
		BelongsToHousehold: fake.UUID(),
		HouseholdRole:      authorization.HouseholdMemberRole.String(),
		CreatedAt:          BuildFakeTime(),
		ArchivedAt:         nil,
	}
}
