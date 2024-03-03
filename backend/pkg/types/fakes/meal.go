package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeMeal builds a faked meal.
func BuildFakeMeal() *types.Meal {
	recipes := []*types.MealComponent{}
	for i := 0; i < exampleQuantity; i++ {
		recipes = append(recipes, BuildFakeMealComponent())
	}

	return &types.Meal{
		ID:                       BuildFakeID(),
		Name:                     buildUniqueString(),
		Description:              buildUniqueString(),
		MinimumEstimatedPortions: float32(buildFakeNumber()),
		CreatedAt:                BuildFakeTime(),
		CreatedByUser:            BuildFakeID(),
		Components:               recipes,
		EligibleForMealPlans:     true,
	}
}

// BuildFakeMealComponent builds a faked meal component.
func BuildFakeMealComponent() *types.MealComponent {
	return &types.MealComponent{
		Recipe:        *BuildFakeRecipe(),
		RecipeScale:   float32(1.0),
		ComponentType: types.MealComponentTypesMain,
	}
}

// BuildFakeMealList builds a faked MealList.
func BuildFakeMealList() *types.QueryFilteredResult[types.Meal] {
	var examples []*types.Meal
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMeal())
	}

	return &types.QueryFilteredResult[types.Meal]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeMealCreationRequestInput builds a faked MealCreationRequestInput.
func BuildFakeMealCreationRequestInput() *types.MealCreationRequestInput {
	recipe := BuildFakeMeal()
	return converters.ConvertMealToMealCreationRequestInput(recipe)
}
