package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepCompletionCondition_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionCondition{}
		input := &RecipeStepCompletionConditionUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.Optional = pointer.To(true)

		x.Update(input)
	})
}

func TestRecipeStepCompletionConditionCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionCreationRequestInput{
			IngredientStateID:   t.Name(),
			BelongsToRecipeStep: t.Name(),
			Optional:            fake.Bool(),
			Ingredients:         []uint64{123},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepCompletionConditionIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionIngredientCreationRequestInput{
			RecipeStepIngredient: t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionIngredientCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepCompletionConditionForExistingRecipeCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{
			IngredientStateID: t.Name(),
			Ingredients: []*RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput{
				{
					RecipeStepIngredient: t.Name(),
				},
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput{
			RecipeStepIngredient: t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepCompletionConditionDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionDatabaseCreationInput{
			ID:                  t.Name(),
			IngredientStateID:   t.Name(),
			BelongsToRecipeStep: t.Name(),
			Ingredients: []*RecipeStepCompletionConditionIngredientDatabaseCreationInput{
				{
					RecipeStepIngredient:                   t.Name(),
					BelongsToRecipeStepCompletionCondition: t.Name(),
				},
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepCompletionConditionIngredientDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionIngredientDatabaseCreationInput{
			RecipeStepIngredient:                   t.Name(),
			BelongsToRecipeStepCompletionCondition: t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionIngredientDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepCompletionConditionUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionUpdateRequestInput{
			IngredientStateID:   pointer.To(t.Name()),
			BelongsToRecipeStep: pointer.To(t.Name()),
			Optional:            pointer.To(fake.Bool()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
