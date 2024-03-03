package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeValidIngredientMeasurementUnit builds a faked valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnit() *types.ValidIngredientMeasurementUnit {
	minQty := buildFakeNumber()

	return &types.ValidIngredientMeasurementUnit{
		ID:                       BuildFakeID(),
		Notes:                    buildUniqueString(),
		MeasurementUnit:          *BuildFakeValidMeasurementUnit(),
		Ingredient:               *BuildFakeValidIngredient(),
		MinimumAllowableQuantity: float32(minQty),
		MaximumAllowableQuantity: pointer.To(float32(minQty + 1)),
		CreatedAt:                BuildFakeTime(),
	}
}

// BuildFakeValidIngredientMeasurementUnitList builds a faked ValidIngredientMeasurementUnitList.
func BuildFakeValidIngredientMeasurementUnitList() *types.QueryFilteredResult[types.ValidIngredientMeasurementUnit] {
	var examples []*types.ValidIngredientMeasurementUnit
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredientMeasurementUnit())
	}

	return &types.QueryFilteredResult[types.ValidIngredientMeasurementUnit]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidIngredientMeasurementUnitUpdateRequestInput builds a faked ValidIngredientMeasurementUnitUpdateRequestInput from a valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnitUpdateRequestInput() *types.ValidIngredientMeasurementUnitUpdateRequestInput {
	validIngredientMeasurementUnit := BuildFakeValidIngredientMeasurementUnit()
	return &types.ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                    &validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID:   &validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:        &validIngredientMeasurementUnit.Ingredient.ID,
		MinimumAllowableQuantity: &validIngredientMeasurementUnit.MinimumAllowableQuantity,
		MaximumAllowableQuantity: validIngredientMeasurementUnit.MaximumAllowableQuantity,
	}
}

// BuildFakeValidIngredientMeasurementUnitCreationRequestInput builds a faked ValidIngredientMeasurementUnitCreationRequestInput.
func BuildFakeValidIngredientMeasurementUnitCreationRequestInput() *types.ValidIngredientMeasurementUnitCreationRequestInput {
	validIngredientMeasurementUnit := BuildFakeValidIngredientMeasurementUnit()
	return converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(validIngredientMeasurementUnit)
}
