package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestValidPreparationCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationRequestInput{
			Name:                  t.Name(),
			Description:           t.Name(),
			IconPath:              t.Name(),
			PastTense:             t.Name(),
			YieldsNothing:         fake.Bool(),
			RestrictToIngredients: fake.Bool(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateRequestInput{
			Name:                  pointer.To(t.Name()),
			Description:           pointer.To(t.Name()),
			IconPath:              pointer.To(t.Name()),
			PastTense:             pointer.To(t.Name()),
			YieldsNothing:         pointer.To(fake.Bool()),
			RestrictToIngredients: pointer.To(fake.Bool()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparation_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparation{
			MaximumInstrumentCount: pointer.To(int32(0)),
			MaximumIngredientCount: pointer.To(int32(0)),
			MaximumVesselCount:     pointer.To(int32(0)),
		}
		input := &ValidPreparationUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.YieldsNothing = pointer.To(true)
		input.RestrictToIngredients = pointer.To(true)
		input.MaximumIngredientCount = pointer.To(int32(1))
		input.MaximumInstrumentCount = pointer.To(int32(1))
		input.MaximumVesselCount = pointer.To(int32(1))
		input.TemperatureRequired = pointer.To(true)
		input.TimeEstimateRequired = pointer.To(true)
		input.OnlyForVessels = pointer.To(true)
		input.ConsumesVessel = pointer.To(true)
		input.ConditionExpressionRequired = pointer.To(true)

		x.Update(input)
	})
}

func TestValidPreparationCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationCreationRequestInput{
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationDatabaseCreationInput{
			ID:   t.Name(),
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationUpdateRequestInput{
			Name: pointer.To(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
