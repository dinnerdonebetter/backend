package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertHouseholdCreationInputToHouseholdDatabaseCreationInput creates a HouseholdDatabaseCreationInput from a HouseholdCreationRequestInput.
func ConvertHouseholdCreationInputToHouseholdDatabaseCreationInput(input *types.HouseholdCreationRequestInput) *types.HouseholdDatabaseCreationInput {
	x := &types.HouseholdDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         input.Name,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		City:         input.City,
		State:        input.State,
		ZipCode:      input.ZipCode,
		Country:      input.Country,
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		ContactPhone: input.ContactPhone,
	}

	return x
}

// ConvertHouseholdToHouseholdUpdateRequestInput creates a HouseholdUpdateRequestInput from a Household.
func ConvertHouseholdToHouseholdUpdateRequestInput(input *types.Household) *types.HouseholdUpdateRequestInput {
	x := &types.HouseholdUpdateRequestInput{
		Name:          &input.Name,
		AddressLine1:  &input.AddressLine1,
		AddressLine2:  &input.AddressLine2,
		City:          &input.City,
		State:         &input.State,
		ZipCode:       &input.ZipCode,
		Country:       &input.Country,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		ContactPhone:  &input.ContactPhone,
		BelongsToUser: input.BelongsToUser,
	}

	return x
}

// ConvertHouseholdToHouseholdCreationRequestInput builds a faked HouseholdCreationRequestInput from a household.
func ConvertHouseholdToHouseholdCreationRequestInput(household *types.Household) *types.HouseholdCreationRequestInput {
	return &types.HouseholdCreationRequestInput{
		Name:         household.Name,
		AddressLine1: household.AddressLine1,
		AddressLine2: household.AddressLine2,
		City:         household.City,
		State:        household.State,
		ZipCode:      household.ZipCode,
		Country:      household.Country,
		Latitude:     household.Latitude,
		Longitude:    household.Longitude,
		ContactPhone: household.ContactPhone,
	}
}

// ConvertHouseholdToHouseholdDatabaseCreationInput builds a faked HouseholdCreationRequestInput.
func ConvertHouseholdToHouseholdDatabaseCreationInput(household *types.Household) *types.HouseholdDatabaseCreationInput {
	return &types.HouseholdDatabaseCreationInput{
		ID:                   household.ID,
		Name:                 household.Name,
		AddressLine1:         household.AddressLine1,
		AddressLine2:         household.AddressLine2,
		City:                 household.City,
		State:                household.State,
		ZipCode:              household.ZipCode,
		Country:              household.Country,
		Latitude:             household.Latitude,
		Longitude:            household.Longitude,
		ContactPhone:         household.ContactPhone,
		BelongsToUser:        household.BelongsToUser,
		WebhookEncryptionKey: household.WebhookEncryptionKey,
	}
}

// ConvertHouseholdUserMembershipCreationRequestInputToHouseholdUserMembershipDatabaseCreationInput builds a HouseholdUserMembershipDatabaseCreationInput from a HouseholdUserMembershipCreationRequestInput.
func ConvertHouseholdUserMembershipCreationRequestInputToHouseholdUserMembershipDatabaseCreationInput(input *types.HouseholdUserMembershipCreationRequestInput) *types.HouseholdUserMembershipDatabaseCreationInput {
	return &types.HouseholdUserMembershipDatabaseCreationInput{
		Reason: input.Reason,
		UserID: input.UserID,
	}
}
