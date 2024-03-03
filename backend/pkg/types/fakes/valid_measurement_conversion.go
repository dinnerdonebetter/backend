package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeValidMeasurementUnitConversion builds a faked valid preparation.
func BuildFakeValidMeasurementUnitConversion() *types.ValidMeasurementUnitConversion {
	return &types.ValidMeasurementUnitConversion{
		ID:                BuildFakeID(),
		From:              *BuildFakeValidMeasurementUnit(),
		To:                *BuildFakeValidMeasurementUnit(),
		OnlyForIngredient: nil,
		Modifier:          float32(buildFakeNumber()),
		Notes:             buildUniqueString(),
		CreatedAt:         BuildFakeTime(),
	}
}

// BuildFakeValidMeasurementUnitConversionList builds a faked ValidMeasurementUnitConversionList.
func BuildFakeValidMeasurementUnitConversionList() *types.QueryFilteredResult[types.ValidMeasurementUnitConversion] {
	var examples []*types.ValidMeasurementUnitConversion
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidMeasurementUnitConversion())
	}

	return &types.QueryFilteredResult[types.ValidMeasurementUnitConversion]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput builds a faked ValidMeasurementUnitConversionUpdateRequestInput from a valid preparation.
func BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput() *types.ValidMeasurementUnitConversionUpdateRequestInput {
	validMeasurementUnitConversion := BuildFakeValidMeasurementUnitConversion()

	x := &types.ValidMeasurementUnitConversionUpdateRequestInput{
		From:     &validMeasurementUnitConversion.From.ID,
		To:       &validMeasurementUnitConversion.To.ID,
		Modifier: &validMeasurementUnitConversion.Modifier,
		Notes:    &validMeasurementUnitConversion.Notes,
	}

	if validMeasurementUnitConversion.OnlyForIngredient != nil {
		x.OnlyForIngredient = &validMeasurementUnitConversion.OnlyForIngredient.ID
	}

	return x
}

// BuildFakeValidMeasurementUnitConversionCreationRequestInput builds a faked ValidMeasurementUnitConversionCreationRequestInput.
func BuildFakeValidMeasurementUnitConversionCreationRequestInput() *types.ValidMeasurementUnitConversionCreationRequestInput {
	validMeasurementUnitConversion := BuildFakeValidMeasurementUnitConversion()
	return converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionCreationRequestInput(validMeasurementUnitConversion)
}
