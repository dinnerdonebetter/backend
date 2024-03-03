package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeValidPreparationVessel builds a faked valid preparation instrument.
func BuildFakeValidPreparationVessel() *types.ValidPreparationVessel {
	return &types.ValidPreparationVessel{
		ID:          BuildFakeID(),
		Notes:       buildUniqueString(),
		Preparation: *BuildFakeValidPreparation(),
		Vessel:      *BuildFakeValidVessel(),
		CreatedAt:   BuildFakeTime(),
	}
}

// BuildFakeValidPreparationVesselList builds a faked ValidPreparationVesselList.
func BuildFakeValidPreparationVesselList() *types.QueryFilteredResult[types.ValidPreparationVessel] {
	var examples []*types.ValidPreparationVessel
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidPreparationVessel())
	}

	return &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidPreparationVesselUpdateRequestInput builds a faked ValidPreparationVesselUpdateRequestInput from a valid preparation instrument.
func BuildFakeValidPreparationVesselUpdateRequestInput() *types.ValidPreparationVesselUpdateRequestInput {
	validPreparationVessel := BuildFakeValidPreparationVessel()
	return &types.ValidPreparationVesselUpdateRequestInput{
		Notes:              &validPreparationVessel.Notes,
		ValidPreparationID: &validPreparationVessel.Preparation.ID,
		ValidVesselID:      &validPreparationVessel.Vessel.ID,
	}
}

// BuildFakeValidPreparationVesselCreationRequestInput builds a faked ValidPreparationVesselCreationRequestInput.
func BuildFakeValidPreparationVesselCreationRequestInput() *types.ValidPreparationVesselCreationRequestInput {
	validPreparationVessel := BuildFakeValidPreparationVessel()
	return converters.ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput(validPreparationVessel)
}
