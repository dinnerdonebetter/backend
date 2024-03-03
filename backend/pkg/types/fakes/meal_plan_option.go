package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeMealPlanOption builds a faked meal plan option.
func BuildFakeMealPlanOption() *types.MealPlanOption {
	var examples []*types.MealPlanOptionVote
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanOptionVote())
	}

	meal := BuildFakeMeal()
	meal.Components = nil

	return &types.MealPlanOption{
		ID:                     BuildFakeID(),
		AssignedCook:           func(s string) *string { return &s }(BuildFakeID()),
		Meal:                   *meal,
		Votes:                  examples,
		Chosen:                 false,
		Notes:                  buildUniqueString(),
		CreatedAt:              BuildFakeTime(),
		BelongsToMealPlanEvent: fake.UUID(),
	}
}

// BuildFakeMealPlanOptionList builds a faked MealPlanOptionList.
func BuildFakeMealPlanOptionList() *types.QueryFilteredResult[types.MealPlanOption] {
	var examples []*types.MealPlanOption
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanOption())
	}

	return &types.QueryFilteredResult[types.MealPlanOption]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeMealPlanOptionUpdateRequestInput builds a faked MealPlanOptionUpdateRequestInput from a meal plan option.
func BuildFakeMealPlanOptionUpdateRequestInput() *types.MealPlanOptionUpdateRequestInput {
	mealPlanOption := BuildFakeMealPlanOption()
	return converters.ConvertMealPlanOptionToMealPlanOptionUpdateRequestInput(mealPlanOption)
}

// BuildFakeMealPlanOptionCreationRequestInput builds a faked MealPlanOptionCreationRequestInput.
func BuildFakeMealPlanOptionCreationRequestInput() *types.MealPlanOptionCreationRequestInput {
	mealPlanOption := BuildFakeMealPlanOption()
	return converters.ConvertMealPlanOptionToMealPlanOptionCreationRequestInput(mealPlanOption)
}
