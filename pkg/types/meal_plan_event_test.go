package types

import (
	"context"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestMealPlanEventCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealPlanEventCreationRequestInput{
			MealName: SecondBreakfastMealName,
			StartsAt: time.Now(),
			EndsAt:   time.Now(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})

	T.Run("with invalid time", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		tt := time.Now()
		x := &MealPlanEventCreationRequestInput{
			MealName: SecondBreakfastMealName,
			StartsAt: tt,
			EndsAt:   tt,
		}

		assert.Error(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanEventDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealPlanEventDatabaseCreationInput{
			ID:                t.Name(),
			BelongsToMealPlan: t.Name(),
			MealName:          SecondBreakfastMealName,
			StartsAt:          time.Now(),
			EndsAt:            time.Now(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanEventUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealPlanEventUpdateRequestInput{
			MealName: pointers.Pointer(SecondBreakfastMealName),
			StartsAt: pointers.Pointer(time.Now()),
			EndsAt:   pointers.Pointer(time.Now()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanEvent_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanEvent{}
		input := &MealPlanEventUpdateRequestInput{}

		fake.Struct(&input)
		input.StartsAt = pointers.Pointer(time.Now())
		input.EndsAt = pointers.Pointer(time.Now())

		x.Update(input)
	})
}