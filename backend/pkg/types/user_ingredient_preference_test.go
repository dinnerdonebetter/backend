package types

import (
	"context"
	"math"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestUserIngredientPreference_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &UserIngredientPreference{}
		input := &UserIngredientPreferenceUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))

		x.Update(input)
	})
}

func TestUserIngredientPreferenceCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		input := &UserIngredientPreferenceCreationRequestInput{
			ValidIngredientID: t.Name(),
			Rating:            1,
		}

		assert.NoError(t, input.ValidateWithContext(ctx))
	})

	T.Run("invalid range", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		input := &UserIngredientPreferenceCreationRequestInput{
			ValidIngredientID: t.Name(),
			Rating:            math.MaxInt8,
		}

		assert.Error(t, input.ValidateWithContext(ctx))
	})

	T.Run("with valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		input := &UserIngredientPreferenceCreationRequestInput{
			ValidIngredientGroupID: t.Name(),
			Rating:                 1,
		}

		assert.NoError(t, input.ValidateWithContext(ctx))
	})
}

func TestUserIngredientPreferenceDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &UserIngredientPreferenceDatabaseCreationInput{
			ValidIngredientID: t.Name(),
			Rating:            minRating,
			BelongsToUser:     t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &UserIngredientPreferenceDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestUserIngredientPreferenceUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &UserIngredientPreferenceUpdateRequestInput{
			IngredientID: pointer.To(t.Name()),
			Rating:       pointer.To(minRating),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &UserIngredientPreferenceUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
