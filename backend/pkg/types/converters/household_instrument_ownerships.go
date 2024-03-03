package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipUpdateRequestInput(x *types.HouseholdInstrumentOwnership) *types.HouseholdInstrumentOwnershipUpdateRequestInput {
	out := &types.HouseholdInstrumentOwnershipUpdateRequestInput{
		Notes:             &x.Notes,
		Quantity:          &x.Quantity,
		ValidInstrumentID: &x.Instrument.ID,
	}

	return out
}

// ConvertHouseholdInstrumentOwnershipCreationRequestInputToHouseholdInstrumentOwnershipDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertHouseholdInstrumentOwnershipCreationRequestInputToHouseholdInstrumentOwnershipDatabaseCreationInput(x *types.HouseholdInstrumentOwnershipCreationRequestInput) *types.HouseholdInstrumentOwnershipDatabaseCreationInput {
	out := &types.HouseholdInstrumentOwnershipDatabaseCreationInput{
		ID:                identifiers.New(),
		Notes:             x.Notes,
		Quantity:          x.Quantity,
		ValidInstrumentID: x.ValidInstrumentID,
	}

	return out
}

// ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipCreationRequestInput builds a HouseholdInstrumentOwnershipCreationRequestInput from a Ingredient.
func ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipCreationRequestInput(x *types.HouseholdInstrumentOwnership) *types.HouseholdInstrumentOwnershipCreationRequestInput {
	return &types.HouseholdInstrumentOwnershipCreationRequestInput{
		Notes:              x.Notes,
		Quantity:           x.Quantity,
		ValidInstrumentID:  x.Instrument.ID,
		BelongsToHousehold: x.BelongsToHousehold,
	}
}

// ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipDatabaseCreationInput builds a HouseholdInstrumentOwnershipDatabaseCreationInput from a HouseholdInstrumentOwnership.
func ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipDatabaseCreationInput(x *types.HouseholdInstrumentOwnership) *types.HouseholdInstrumentOwnershipDatabaseCreationInput {
	return &types.HouseholdInstrumentOwnershipDatabaseCreationInput{
		ID:                 x.ID,
		Notes:              x.Notes,
		Quantity:           x.Quantity,
		ValidInstrumentID:  x.Instrument.ID,
		BelongsToHousehold: x.BelongsToHousehold,
	}
}
