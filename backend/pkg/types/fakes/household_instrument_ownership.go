package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeHouseholdInstrumentOwnership builds a faked valid ingredient.
func BuildFakeHouseholdInstrumentOwnership() *types.HouseholdInstrumentOwnership {
	return &types.HouseholdInstrumentOwnership{
		CreatedAt:          BuildFakeTime(),
		ID:                 buildUniqueString(),
		Notes:              buildUniqueString(),
		BelongsToHousehold: buildUniqueString(),
		Instrument:         *BuildFakeValidInstrument(),
		Quantity:           uint16(buildFakeNumber()),
	}
}

// BuildFakeHouseholdInstrumentOwnershipList builds a faked HouseholdInstrumentOwnershipList.
func BuildFakeHouseholdInstrumentOwnershipList() *types.QueryFilteredResult[types.HouseholdInstrumentOwnership] {
	var examples []*types.HouseholdInstrumentOwnership
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeHouseholdInstrumentOwnership())
	}

	return &types.QueryFilteredResult[types.HouseholdInstrumentOwnership]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeHouseholdInstrumentOwnershipUpdateRequestInput builds a faked HouseholdInstrumentOwnershipUpdateRequestInput from a valid ingredient.
func BuildFakeHouseholdInstrumentOwnershipUpdateRequestInput() *types.HouseholdInstrumentOwnershipUpdateRequestInput {
	validIngredient := BuildFakeHouseholdInstrumentOwnership()
	return converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipUpdateRequestInput(validIngredient)
}

// BuildFakeHouseholdInstrumentOwnershipCreationRequestInput builds a faked HouseholdInstrumentOwnershipCreationRequestInput.
func BuildFakeHouseholdInstrumentOwnershipCreationRequestInput() *types.HouseholdInstrumentOwnershipCreationRequestInput {
	validIngredient := BuildFakeHouseholdInstrumentOwnership()
	return converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipCreationRequestInput(validIngredient)
}
