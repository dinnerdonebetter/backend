package types

import (
	"context"
	"math"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrument{
			RecipeStepProductID: pointer.To(t.Name()),
			MaximumQuantity:     pointer.To(uint32(321)),
		}
		input := &RecipeStepInstrumentUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.Optional = pointer.To(true)
		input.RecipeStepProductID = pointer.To("whatever")
		input.MaximumQuantity = pointer.To(uint32(123))

		x.Update(input)
	})
}

func TestRecipeStepInstrumentCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentCreationRequestInput{
			InstrumentID:        pointer.To(t.Name()),
			Name:                t.Name(),
			RecipeStepProductID: pointer.To(t.Name()),
			Notes:               t.Name(),
			PreferenceRank:      uint8(fake.Number(1, math.MaxUint8)),
			Optional:            fake.Bool(),
			MinimumQuantity:     fake.Uint32(),
			MaximumQuantity:     pointer.To(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepInstrumentDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentDatabaseCreationInput{
			ID:                  t.Name(),
			InstrumentID:        pointer.To(t.Name()),
			BelongsToRecipeStep: t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepInstrumentUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentUpdateRequestInput{
			InstrumentID:        pointer.To(t.Name()),
			Name:                pointer.To(t.Name()),
			BelongsToRecipeStep: pointer.To(t.Name()),
			RecipeStepProductID: pointer.To(t.Name()),
			Notes:               pointer.To(t.Name()),
			PreferenceRank:      pointer.To(uint8(fake.Number(1, math.MaxUint8))),
			Optional:            pointer.To(fake.Bool()),
			MinimumQuantity:     pointer.To(fake.Uint32()),
			MaximumQuantity:     pointer.To(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
