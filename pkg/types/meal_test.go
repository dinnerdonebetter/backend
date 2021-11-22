package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestMealCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealCreationRequestInput{
			Name:        fake.LoremIpsumSentence(exampleQuantity),
			Description: fake.LoremIpsumSentence(exampleQuantity),
			Recipes:     []string{fake.LoremIpsumSentence(exampleQuantity)},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestMealUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealUpdateRequestInput{
			Name:          fake.LoremIpsumSentence(exampleQuantity),
			Description:   fake.LoremIpsumSentence(exampleQuantity),
			CreatedByUser: fake.LoremIpsumSentence(exampleQuantity),
			Recipes:       []string{fake.LoremIpsumSentence(exampleQuantity)},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &MealUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
