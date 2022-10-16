package converters

import (
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertValidMeasurementUnitToValidMeasurementUnitUpdateRequestInput creates a ValidMeasurementUnitUpdateRequestInput from a ValidMeasurementUnit.
func ConvertValidMeasurementUnitToValidMeasurementUnitUpdateRequestInput(input *types.ValidMeasurementUnit) *types.ValidMeasurementUnitUpdateRequestInput {
	x := &types.ValidMeasurementUnitUpdateRequestInput{
		Name:        &input.Name,
		Description: &input.Description,
		IconPath:    &input.IconPath,
		Volumetric:  &input.Volumetric,
		Universal:   &input.Universal,
		Metric:      &input.Metric,
		Imperial:    &input.Imperial,
		PluralName:  &input.PluralName,
	}

	return x
}

// ConvertValidMeasurementUnitCreationRequestInputToValidMeasurementUnitDatabaseCreationInput creates a ValidMeasurementUnitDatabaseCreationInput from a ValidMeasurementUnitCreationRequestInput.
func ConvertValidMeasurementUnitCreationRequestInputToValidMeasurementUnitDatabaseCreationInput(input *types.ValidMeasurementUnitCreationRequestInput) *types.ValidMeasurementUnitDatabaseCreationInput {
	x := &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		Volumetric:  input.Volumetric,
		IconPath:    input.IconPath,
		Universal:   input.Universal,
		Metric:      input.Metric,
		Imperial:    input.Imperial,
		PluralName:  input.PluralName,
	}

	return x
}

// ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput builds a ValidMeasurementUnitCreationRequestInput from a ValidMeasurementUnit.
func ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(validMeasurementUnit *types.ValidMeasurementUnit) *types.ValidMeasurementUnitCreationRequestInput {
	return &types.ValidMeasurementUnitCreationRequestInput{
		Name:        validMeasurementUnit.Name,
		Description: validMeasurementUnit.Description,
		Volumetric:  validMeasurementUnit.Volumetric,
		IconPath:    validMeasurementUnit.IconPath,
		Universal:   validMeasurementUnit.Universal,
		Metric:      validMeasurementUnit.Metric,
		Imperial:    validMeasurementUnit.Imperial,
		PluralName:  validMeasurementUnit.PluralName,
	}
}

// ConvertValidMeasurementUnitToValidMeasurementUnitDatabaseCreationInput builds a ValidMeasurementUnitDatabaseCreationInput from a ValidMeasurementUnit.
func ConvertValidMeasurementUnitToValidMeasurementUnitDatabaseCreationInput(validMeasurementUnit *types.ValidMeasurementUnit) *types.ValidMeasurementUnitDatabaseCreationInput {
	return &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:          validMeasurementUnit.ID,
		Name:        validMeasurementUnit.Name,
		Description: validMeasurementUnit.Description,
		Volumetric:  validMeasurementUnit.Volumetric,
		IconPath:    validMeasurementUnit.IconPath,
		Universal:   validMeasurementUnit.Universal,
		Metric:      validMeasurementUnit.Metric,
		Imperial:    validMeasurementUnit.Imperial,
		PluralName:  validMeasurementUnit.PluralName,
	}
}
