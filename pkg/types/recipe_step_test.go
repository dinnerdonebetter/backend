package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/pointers"
)

func buildValidRecipeStepCreationRequestInput() *RecipeStepCreationRequestInput {
	return &RecipeStepCreationRequestInput{
		Index:                         fake.Uint32(),
		PreparationID:                 fake.LoremIpsumSentence(exampleQuantity),
		MinimumEstimatedTimeInSeconds: uint32Pointer(fake.Uint32()),
		MaximumEstimatedTimeInSeconds: uint32Pointer(fake.Uint32()),
		MinimumTemperatureInCelsius:   pointers.Float32Pointer(float32(123.45)),
		Notes:                         fake.LoremIpsumSentence(exampleQuantity),
		ExplicitInstructions:          fake.LoremIpsumSentence(exampleQuantity),
		Products: []*RecipeStepProductCreationRequestInput{
			{
				Name:              fake.LoremIpsumSentence(exampleQuantity),
				Type:              RecipeStepProductIngredientType,
				MeasurementUnitID: fake.LoremIpsumSentence(exampleQuantity),
				MinimumQuantity:   fake.Float32(),
				QuantityNotes:     fake.LoremIpsumSentence(exampleQuantity),
			},
		},
		Ingredients: []*RecipeStepIngredientCreationRequestInput{
			{
				IngredientID:        func(s string) *string { return &s }(fake.LoremIpsumSentence(exampleQuantity)),
				ID:                  fake.LoremIpsumSentence(exampleQuantity),
				MeasurementUnitID:   fake.LoremIpsumSentence(exampleQuantity),
				QuantityNotes:       fake.LoremIpsumSentence(exampleQuantity),
				IngredientNotes:     fake.LoremIpsumSentence(exampleQuantity),
				BelongsToRecipeStep: fake.LoremIpsumSentence(exampleQuantity),
				MinimumQuantity:     1,
			},
		},
	}
}

func TestRecipeStepCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := buildValidRecipeStepCreationRequestInput()

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})

	T.Run("with too many ingredients", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCreationRequestInput{
			Index:                         fake.Uint32(),
			PreparationID:                 fake.LoremIpsumSentence(exampleQuantity),
			MinimumEstimatedTimeInSeconds: uint32Pointer(fake.Uint32()),
			MaximumEstimatedTimeInSeconds: uint32Pointer(fake.Uint32()),
			MinimumTemperatureInCelsius:   pointers.Float32Pointer(float32(123.45)),
			Notes:                         fake.LoremIpsumSentence(exampleQuantity),
			ExplicitInstructions:          fake.LoremIpsumSentence(exampleQuantity),
			Products: []*RecipeStepProductCreationRequestInput{
				{
					Name: fake.LoremIpsumSentence(exampleQuantity),
				},
			},
			Ingredients: []*RecipeStepIngredientCreationRequestInput{},
		}

		for i := 0; i < maxIngredientsPerStep*2; i++ {
			x.Ingredients = append(x.Ingredients, &RecipeStepIngredientCreationRequestInput{
				IngredientID:        func(s string) *string { return &s }(fake.LoremIpsumSentence(exampleQuantity)),
				ID:                  fake.LoremIpsumSentence(exampleQuantity),
				MeasurementUnitID:   fake.LoremIpsumSentence(exampleQuantity),
				QuantityNotes:       fake.LoremIpsumSentence(exampleQuantity),
				IngredientNotes:     fake.LoremIpsumSentence(exampleQuantity),
				BelongsToRecipeStep: fake.LoremIpsumSentence(exampleQuantity),
				MinimumQuantity:     1,
			})
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateRequestInput{
			Index:                         uint32Pointer(fake.Uint32()),
			Preparation:                   &ValidPreparation{},
			MinimumEstimatedTimeInSeconds: uint32Pointer(fake.Uint32()),
			MaximumEstimatedTimeInSeconds: uint32Pointer(fake.Uint32()),
			MinimumTemperatureInCelsius:   pointers.Float32Pointer(float32(123.45)),
			Notes:                         stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			ExplicitInstructions:          stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
