package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientMeasurementUnitCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitCreationRequestInput{
			Notes:                    fake.LoremIpsumSentence(exampleQuantity),
			ValidMeasurementUnitID:   fake.LoremIpsumSentence(exampleQuantity),
			ValidIngredientID:        fake.LoremIpsumSentence(exampleQuantity),
			MinimumAllowableQuantity: fake.Float32(),
			MaximumAllowableQuantity: pointers.Pointer(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidIngredientMeasurementUnitUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitUpdateRequestInput{
			Notes:                    pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			ValidMeasurementUnitID:   pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			ValidIngredientID:        pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumAllowableQuantity: pointers.Pointer(fake.Float32()),
			MaximumAllowableQuantity: pointers.Pointer(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
