package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestMealCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealCreationRequestInput{
			Name:        t.Name(),
			Description: t.Name(),
			Components: []*MealComponentCreationRequestInput{
				{
					RecipeID:      t.Name(),
					ComponentType: MealComponentTypesMain,
				},
			},
		}

		assert.NoError(t, x.ValidateWithContext(context.Background()))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealCreationRequestInput{}

		assert.Error(t, x.ValidateWithContext(context.Background()))
	})

	T.Run("with invalid component", func(t *testing.T) {
		t.Parallel()

		x := &MealCreationRequestInput{
			Name:        t.Name(),
			Description: t.Name(),
			Components: []*MealComponentCreationRequestInput{
				{},
			},
		}

		assert.Error(t, x.ValidateWithContext(context.Background()))
	})
}

func TestMealUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealUpdateRequestInput{
			Name:          pointers.Pointer(t.Name()),
			Description:   pointers.Pointer(t.Name()),
			CreatedByUser: pointers.Pointer(t.Name()),
			Components: []*MealComponentUpdateRequestInput{
				{
					RecipeID:      pointers.Pointer(t.Name()),
					RecipeScale:   pointers.Pointer(float32(exampleQuantity)),
					ComponentType: pointers.Pointer(MealComponentTypesAmuseBouche),
				},
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestMealComponentCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealComponentCreationRequestInput{
			RecipeID:      t.Name(),
			RecipeScale:   exampleQuantity,
			ComponentType: MealComponentTypesAmuseBouche,
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})
}

func TestMealDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealDatabaseCreationInput{
			Name: t.Name(),
			Components: []*MealComponentDatabaseCreationInput{
				{
					RecipeID: t.Name(),
				},
			},
			CreatedByUser: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealUpdateRequestInput{
			Name:        pointers.Pointer(t.Name()),
			Description: pointers.Pointer(t.Name()),
			Components: []*MealComponentUpdateRequestInput{
				{
					RecipeID: pointers.Pointer(t.Name()),
				},
			},
			CreatedByUser: pointers.Pointer(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealComponent_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealComponent{}
		input := &MealComponentUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))

		x.Update(input)
	})
}

func TestMeal_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Meal{}
		input := &MealUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.EligibleForMealPlans = pointers.Pointer(true)

		x.Update(input)
	})
}
