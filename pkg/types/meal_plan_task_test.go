package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/require"
)

func TestMealPlanTaskDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := MealPlanTaskDatabaseCreationInput{
			ID:               fake.LoremIpsumSentence(exampleQuantity),
			MealPlanOptionID: fake.LoremIpsumSentence(exampleQuantity),
			// RecipeStepID:         fake.LoremIpsumSentence(exampleQuantity),
			CannotCompleteBefore: fake.Date(),
			CannotCompleteAfter:  fake.Date(),
			CompletedAt:          nil,
		}

		require.NoError(t, x.ValidateWithContext(ctx))
	})
}
