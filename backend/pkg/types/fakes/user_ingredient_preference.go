package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeUserIngredientPreference builds a faked valid preparation.
func BuildFakeUserIngredientPreference() *types.UserIngredientPreference {
	return &types.UserIngredientPreference{
		ID:         BuildFakeID(),
		Ingredient: *BuildFakeValidIngredient(),
		Rating:     1,
		Notes:      buildUniqueString(),
		Allergy:    fake.Bool(),
		CreatedAt:  BuildFakeTime(),
	}
}

// BuildFakeUserIngredientPreferenceList builds a faked UserIngredientPreferenceList.
func BuildFakeUserIngredientPreferenceList() *types.QueryFilteredResult[types.UserIngredientPreference] {
	var examples []*types.UserIngredientPreference
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeUserIngredientPreference())
	}

	return &types.QueryFilteredResult[types.UserIngredientPreference]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeUserIngredientPreferenceUpdateRequestInput builds a faked UserIngredientPreferenceUpdateRequestInput from a valid preparation.
func BuildFakeUserIngredientPreferenceUpdateRequestInput() *types.UserIngredientPreferenceUpdateRequestInput {
	validPreparation := BuildFakeUserIngredientPreference()
	return converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceUpdateRequestInput(validPreparation)
}

// BuildFakeUserIngredientPreferenceCreationRequestInput builds a faked UserIngredientPreferenceCreationRequestInput.
func BuildFakeUserIngredientPreferenceCreationRequestInput() *types.UserIngredientPreferenceCreationRequestInput {
	validPreparation := BuildFakeUserIngredientPreference()
	return converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(validPreparation)
}
