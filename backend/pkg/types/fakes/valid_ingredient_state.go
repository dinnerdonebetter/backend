package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeValidIngredientState builds a faked valid preparation.
func BuildFakeValidIngredientState() *types.ValidIngredientState {
	return &types.ValidIngredientState{
		ID:            BuildFakeID(),
		Name:          buildUniqueString(),
		Description:   buildUniqueString(),
		IconPath:      buildUniqueString(),
		Slug:          buildUniqueString(),
		PastTense:     buildUniqueString(),
		AttributeType: types.ValidIngredientStateAttributeTypeOther,
		CreatedAt:     BuildFakeTime(),
	}
}

// BuildFakeValidIngredientStateList builds a faked ValidIngredientStateList.
func BuildFakeValidIngredientStateList() *types.QueryFilteredResult[types.ValidIngredientState] {
	var examples []*types.ValidIngredientState
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredientState())
	}

	return &types.QueryFilteredResult[types.ValidIngredientState]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidIngredientStateUpdateRequestInput builds a faked ValidIngredientStateUpdateRequestInput from a valid preparation.
func BuildFakeValidIngredientStateUpdateRequestInput() *types.ValidIngredientStateUpdateRequestInput {
	validIngredientState := BuildFakeValidIngredientState()
	return converters.ConvertValidIngredientStateToValidIngredientStateUpdateRequestInput(validIngredientState)
}

// BuildFakeValidIngredientStateCreationRequestInput builds a faked ValidIngredientStateCreationRequestInput.
func BuildFakeValidIngredientStateCreationRequestInput() *types.ValidIngredientStateCreationRequestInput {
	validIngredientState := BuildFakeValidIngredientState()
	return converters.ConvertValidIngredientStateToValidIngredientStateCreationRequestInput(validIngredientState)
}
