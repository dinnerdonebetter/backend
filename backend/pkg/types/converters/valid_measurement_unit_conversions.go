package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionUpdateRequestInput creates a ValidMeasurementUnitConversionUpdateRequestInput from a ValidMeasurementUnitConversion.
func ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionUpdateRequestInput(input *types.ValidMeasurementUnitConversion) *types.ValidMeasurementUnitConversionUpdateRequestInput {
	x := &types.ValidMeasurementUnitConversionUpdateRequestInput{
		From:     &input.From.ID,
		To:       &input.To.ID,
		Modifier: &input.Modifier,
		Notes:    &input.Notes,
	}

	if input.OnlyForIngredient != nil {
		x.OnlyForIngredient = &input.OnlyForIngredient.ID
	}

	return x
}

// ConvertValidMeasurementUnitConversionCreationRequestInputToValidMeasurementUnitConversionDatabaseCreationInput creates a ValidMeasurementUnitConversionDatabaseCreationInput from a ValidMeasurementUnitConversionCreationRequestInput.
func ConvertValidMeasurementUnitConversionCreationRequestInputToValidMeasurementUnitConversionDatabaseCreationInput(input *types.ValidMeasurementUnitConversionCreationRequestInput) *types.ValidMeasurementUnitConversionDatabaseCreationInput {
	x := &types.ValidMeasurementUnitConversionDatabaseCreationInput{
		ID:                identifiers.New(),
		From:              input.From,
		To:                input.To,
		OnlyForIngredient: input.OnlyForIngredient,
		Modifier:          input.Modifier,
		Notes:             input.Notes,
	}

	return x
}

// ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionCreationRequestInput builds a ValidMeasurementUnitConversionCreationRequestInput from a ValidMeasurementUnitConversion.
func ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionCreationRequestInput(validMeasurementUnitConversion *types.ValidMeasurementUnitConversion) *types.ValidMeasurementUnitConversionCreationRequestInput {
	x := &types.ValidMeasurementUnitConversionCreationRequestInput{
		From:     validMeasurementUnitConversion.From.ID,
		To:       validMeasurementUnitConversion.To.ID,
		Modifier: validMeasurementUnitConversion.Modifier,
		Notes:    validMeasurementUnitConversion.Notes,
	}

	if validMeasurementUnitConversion.OnlyForIngredient != nil {
		x.OnlyForIngredient = &validMeasurementUnitConversion.OnlyForIngredient.ID
	}

	return x
}

// ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionDatabaseCreationInput builds a ValidMeasurementUnitConversionDatabaseCreationInput from a ValidMeasurementUnitConversion.
func ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionDatabaseCreationInput(validMeasurementUnitConversion *types.ValidMeasurementUnitConversion) *types.ValidMeasurementUnitConversionDatabaseCreationInput {
	x := &types.ValidMeasurementUnitConversionDatabaseCreationInput{
		ID:       validMeasurementUnitConversion.ID,
		From:     validMeasurementUnitConversion.From.ID,
		To:       validMeasurementUnitConversion.To.ID,
		Modifier: validMeasurementUnitConversion.Modifier,
		Notes:    validMeasurementUnitConversion.Notes,
	}

	if validMeasurementUnitConversion.OnlyForIngredient != nil {
		x.OnlyForIngredient = &validMeasurementUnitConversion.OnlyForIngredient.ID
	}

	return x
}
