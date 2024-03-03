package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepProduct_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProduct{
			MeasurementUnit:        &ValidMeasurementUnit{},
			ContainedInVesselIndex: pointer.To(uint16(3)),
		}
		input := &RecipeStepProductUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.Compostable = pointer.To(true)
		input.ContainedInVesselIndex = pointer.To(uint16(1))
		input.IsLiquid = pointer.To(true)
		input.IsWaste = pointer.To(true)

		x.Update(input)
	})
}

func TestRecipeStepProductCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{
			Name:                               t.Name(),
			Type:                               RecipeStepProductIngredientType,
			MeasurementUnitID:                  pointer.To(t.Name()),
			MinimumQuantity:                    pointer.To(fake.Float32()),
			QuantityNotes:                      t.Name(),
			Compostable:                        fake.Bool(),
			MaximumStorageDurationInSeconds:    pointer.To(fake.Uint32()),
			MinimumStorageTemperatureInCelsius: pointer.To(fake.Float32()),
			MaximumStorageTemperatureInCelsius: pointer.To(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepProductUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{
			Name:                               pointer.To(t.Name()),
			Type:                               pointer.To(RecipeStepProductIngredientType),
			MeasurementUnitID:                  pointer.To(t.Name()),
			MinimumQuantity:                    pointer.To(fake.Float32()),
			MaximumQuantity:                    pointer.To(fake.Float32()),
			QuantityNotes:                      pointer.To(t.Name()),
			Compostable:                        pointer.To(fake.Bool()),
			MaximumStorageDurationInSeconds:    pointer.To(fake.Uint32()),
			MinimumStorageTemperatureInCelsius: pointer.To(fake.Float32()),
			MaximumStorageTemperatureInCelsius: pointer.To(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
