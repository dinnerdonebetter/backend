package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestValidVessel_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidVessel{
			CapacityUnit: &ValidMeasurementUnit{},
		}
		input := &ValidVesselUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.CapacityUnitID = pointer.To(t.Name())
		input.UsableForStorage = pointer.To(true)
		input.DisplayInSummaryLists = pointer.To(true)
		input.IncludeInGeneratedInstructions = pointer.To(true)

		x.Update(input)
	})
}

func TestValidVesselCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselCreationRequestInput{
			Name:           t.Name(),
			Description:    t.Name(),
			IconPath:       t.Name(),
			Capacity:       exampleQuantity,
			CapacityUnitID: pointer.To(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidVesselDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselDatabaseCreationInput{
			ID:             t.Name(),
			Name:           t.Name(),
			CapacityUnitID: pointer.To(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidVesselUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselUpdateRequestInput{
			Name:        pointer.To(t.Name()),
			Description: pointer.To(t.Name()),
			IconPath:    pointer.To(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidVesselUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
