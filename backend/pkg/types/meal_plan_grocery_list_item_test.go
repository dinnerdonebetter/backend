package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestMealPlanGroceryListItemCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealPlanGroceryListItemCreationRequestInput{
			BelongsToMealPlan:      t.Name(),
			ValidIngredientID:      t.Name(),
			ValidMeasurementUnitID: t.Name(),
			Status:                 MealPlanGroceryListItemStatusUnknown,
			MinimumQuantityNeeded:  1.23,
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanGroceryListItemDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealPlanGroceryListItemDatabaseCreationInput{
			ID:                     t.Name(),
			BelongsToMealPlan:      t.Name(),
			ValidIngredientID:      t.Name(),
			ValidMeasurementUnitID: t.Name(),
			MinimumQuantityNeeded:  1.23,
			Status:                 MealPlanGroceryListItemStatusUnknown,
			MaximumQuantityNeeded:  pointer.To(float32(1.23)),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanGroceryListItemUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealPlanGroceryListItemUpdateRequestInput{
			BelongsToMealPlan:      pointer.To(t.Name()),
			ValidIngredientID:      pointer.To(t.Name()),
			ValidMeasurementUnitID: pointer.To(t.Name()),
			MinimumQuantityNeeded:  pointer.To(float32(1.23)),
			MaximumQuantityNeeded:  pointer.To(float32(1.23)),
			Status:                 pointer.To(MealPlanGroceryListItemStatusUnknown),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanGroceryListItem_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanGroceryListItem{
			PurchasedMeasurementUnit: &ValidMeasurementUnit{},
			MaximumQuantityNeeded:    pointer.To(float32(1.23)),
		}
		input := &MealPlanGroceryListItemUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.PurchasedMeasurementUnitID = pointer.To(t.Name())
		input.MaximumQuantityNeeded = pointer.To(float32(3.21))

		x.Update(input)
	})
}
