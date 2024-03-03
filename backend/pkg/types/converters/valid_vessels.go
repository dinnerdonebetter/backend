package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertValidVesselToValidVesselUpdateRequestInput creates a ValidVesselUpdateRequestInput from a ValidVessel.
func ConvertValidVesselToValidVesselUpdateRequestInput(x *types.ValidVessel) *types.ValidVesselUpdateRequestInput {
	return &types.ValidVesselUpdateRequestInput{
		Name:                           &x.Name,
		PluralName:                     &x.PluralName,
		Description:                    &x.Description,
		IconPath:                       &x.IconPath,
		UsableForStorage:               &x.UsableForStorage,
		Slug:                           &x.Slug,
		DisplayInSummaryLists:          &x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: &x.IncludeInGeneratedInstructions,
		Capacity:                       &x.Capacity,
		CapacityUnitID:                 &x.CapacityUnit.ID,
		WidthInMillimeters:             &x.WidthInMillimeters,
		LengthInMillimeters:            &x.LengthInMillimeters,
		HeightInMillimeters:            &x.HeightInMillimeters,
		Shape:                          &x.Shape,
	}
}

// ConvertValidVesselCreationRequestInputToValidVesselDatabaseCreationInput creates a ValidVesselDatabaseCreationInput from a ValidVesselCreationRequestInput.
func ConvertValidVesselCreationRequestInputToValidVesselDatabaseCreationInput(x *types.ValidVesselCreationRequestInput) *types.ValidVesselDatabaseCreationInput {
	return &types.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           x.Name,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		IconPath:                       x.IconPath,
		UsableForStorage:               x.UsableForStorage,
		Slug:                           x.Slug,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		Capacity:                       x.Capacity,
		CapacityUnitID:                 x.CapacityUnitID,
		WidthInMillimeters:             x.WidthInMillimeters,
		LengthInMillimeters:            x.LengthInMillimeters,
		HeightInMillimeters:            x.HeightInMillimeters,
		Shape:                          x.Shape,
	}
}

// ConvertNullableValidVesselToValidVessel produces a ValidVessel from a NullableValidVessel.
func ConvertNullableValidVesselToValidVessel(x *types.NullableValidVessel) *types.ValidVessel {
	v := &types.ValidVessel{
		CapacityUnit:  ConvertNullableValidMeasurementUnitToValidMeasurementUnit(x.CapacityUnit),
		LastUpdatedAt: x.LastUpdatedAt,
		ArchivedAt:    x.ArchivedAt,
	}

	if x.ID != nil {
		v.ID = *x.ID
	}
	if x.Name != nil {
		v.Name = *x.Name
	}
	if x.PluralName != nil {
		v.PluralName = *x.PluralName
	}
	if x.Description != nil {
		v.Description = *x.Description
	}
	if x.IconPath != nil {
		v.IconPath = *x.IconPath
	}
	if x.UsableForStorage != nil {
		v.UsableForStorage = *x.UsableForStorage
	}
	if x.Slug != nil {
		v.Slug = *x.Slug
	}
	if x.DisplayInSummaryLists != nil {
		v.DisplayInSummaryLists = *x.DisplayInSummaryLists
	}
	if x.IncludeInGeneratedInstructions != nil {
		v.IncludeInGeneratedInstructions = *x.IncludeInGeneratedInstructions
	}
	if x.Capacity != nil {
		v.Capacity = *x.Capacity
	}
	if x.WidthInMillimeters != nil {
		v.WidthInMillimeters = *x.WidthInMillimeters
	}
	if x.LengthInMillimeters != nil {
		v.LengthInMillimeters = *x.LengthInMillimeters
	}
	if x.HeightInMillimeters != nil {
		v.HeightInMillimeters = *x.HeightInMillimeters
	}
	if x.Shape != nil {
		v.Shape = *x.Shape
	}
	if x.CreatedAt != nil {
		v.CreatedAt = *x.CreatedAt
	}

	return v
}

// ConvertValidVesselToValidVesselCreationRequestInput builds a ValidVesselCreationRequestInput from a ValidVessel.
func ConvertValidVesselToValidVesselCreationRequestInput(x *types.ValidVessel) *types.ValidVesselCreationRequestInput {
	v := &types.ValidVesselCreationRequestInput{
		Name:                           x.Name,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		IconPath:                       x.IconPath,
		UsableForStorage:               x.UsableForStorage,
		Slug:                           x.Slug,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		Capacity:                       x.Capacity,
		WidthInMillimeters:             x.WidthInMillimeters,
		LengthInMillimeters:            x.LengthInMillimeters,
		HeightInMillimeters:            x.HeightInMillimeters,
		Shape:                          x.Shape,
	}

	if x.CapacityUnit != nil {
		v.CapacityUnitID = &x.CapacityUnit.ID
	}

	return v
}

// ConvertValidVesselToValidVesselDatabaseCreationInput builds a ValidVesselDatabaseCreationInput from a ValidVessel.
func ConvertValidVesselToValidVesselDatabaseCreationInput(x *types.ValidVessel) *types.ValidVesselDatabaseCreationInput {
	v := &types.ValidVesselDatabaseCreationInput{
		ID:                             x.ID,
		Name:                           x.Name,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		IconPath:                       x.IconPath,
		UsableForStorage:               x.UsableForStorage,
		Slug:                           x.Slug,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		Capacity:                       x.Capacity,
		WidthInMillimeters:             x.WidthInMillimeters,
		LengthInMillimeters:            x.LengthInMillimeters,
		HeightInMillimeters:            x.HeightInMillimeters,
		Shape:                          x.Shape,
	}

	if x.CapacityUnit != nil {
		v.CapacityUnitID = &x.CapacityUnit.ID
	}

	return v
}

// ConvertValidVesselToValidVesselSearchSubset converts a ValidVessel to a ValidVesselSearchSubset.
func ConvertValidVesselToValidVesselSearchSubset(x *types.ValidVessel) *types.ValidVesselSearchSubset {
	return &types.ValidVesselSearchSubset{
		ID:          x.ID,
		Name:        x.Name,
		PluralName:  x.PluralName,
		Description: x.Description,
	}
}
