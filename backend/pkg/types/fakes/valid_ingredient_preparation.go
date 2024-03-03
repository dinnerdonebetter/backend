package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeValidIngredientPreparation builds a faked valid ingredient preparation.
func BuildFakeValidIngredientPreparation() *types.ValidIngredientPreparation {
	return &types.ValidIngredientPreparation{
		ID:          BuildFakeID(),
		Notes:       buildUniqueString(),
		Preparation: *BuildFakeValidPreparation(),
		Ingredient:  *BuildFakeValidIngredient(),
		CreatedAt:   BuildFakeTime(),
	}
}

// BuildFakeValidIngredientPreparationList builds a faked ValidIngredientPreparationList.
func BuildFakeValidIngredientPreparationList() *types.QueryFilteredResult[types.ValidIngredientPreparation] {
	var examples []*types.ValidIngredientPreparation
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredientPreparation())
	}

	return &types.QueryFilteredResult[types.ValidIngredientPreparation]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidIngredientPreparationUpdateRequestInput builds a faked ValidIngredientPreparationUpdateRequestInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationUpdateRequestInput() *types.ValidIngredientPreparationUpdateRequestInput {
	validIngredientPreparation := BuildFakeValidIngredientPreparation()
	return &types.ValidIngredientPreparationUpdateRequestInput{
		Notes:              &validIngredientPreparation.Notes,
		ValidPreparationID: &validIngredientPreparation.Preparation.ID,
		ValidIngredientID:  &validIngredientPreparation.Ingredient.ID,
	}
}

// BuildFakeValidIngredientPreparationCreationRequestInput builds a faked ValidIngredientPreparationCreationRequestInput.
func BuildFakeValidIngredientPreparationCreationRequestInput() *types.ValidIngredientPreparationCreationRequestInput {
	validIngredientPreparation := BuildFakeValidIngredientPreparation()
	return converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(validIngredientPreparation)
}
