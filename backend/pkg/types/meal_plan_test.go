package types

import (
	"context"
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestMealPlan_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlan{}
		input := &MealPlanUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))

		x.Update(input)
	})
}

func TestMealPlanDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanDatabaseCreationInput{
			ID:                 t.Name(),
			VotingDeadline:     time.Now().Add(24 * time.Hour),
			BelongsToHousehold: t.Name(),
			CreatedByUser:      t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(context.Background()))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanDatabaseCreationInput{}

		assert.Error(t, x.ValidateWithContext(context.Background()))
	})
}

func TestMealPlanCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanCreationRequestInput{
			VotingDeadline: time.Now().Add(24 * time.Hour),
			Events: []*MealPlanEventCreationRequestInput{
				{
					MealName: BreakfastMealName,
					Notes:    t.Name(),
					StartsAt: time.Now(),
					EndsAt:   time.Now().Add(24 * time.Hour),
				},
			},
		}

		assert.NoError(t, x.ValidateWithContext(context.Background()))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanCreationRequestInput{}

		assert.Error(t, x.ValidateWithContext(context.Background()))
	})
}

func TestMealPlanUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleTime := time.Now()

		x := &MealPlanUpdateRequestInput{
			VotingDeadline: &exampleTime,
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
