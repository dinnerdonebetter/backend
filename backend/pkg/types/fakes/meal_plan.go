package fakes

import (
	"time"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeMealPlan builds a faked meal plan.
func BuildFakeMealPlan() *types.MealPlan {
	mealPlanID := BuildFakeID()

	var events []*types.MealPlanEvent
	for i := 0; i < exampleQuantity; i++ {
		event := BuildFakeMealPlanEvent()
		event.BelongsToMealPlan = mealPlanID
		events = append(events, event)
	}

	now := time.Now().Add(30 * time.Minute).Truncate(time.Second).UTC()
	return &types.MealPlan{
		ID:                     mealPlanID,
		Notes:                  buildUniqueString(),
		Status:                 string(types.MealPlanStatusAwaitingVotes),
		VotingDeadline:         now,
		CreatedAt:              BuildFakeTime(),
		BelongsToHousehold:     fake.UUID(),
		TasksCreated:           false,
		GroceryListInitialized: false,
		ElectionMethod:         types.MealPlanElectionMethodSchulze,
		Events:                 events,
		CreatedByUser:          BuildFakeID(),
	}
}

// BuildFakeMealPlanList builds a faked MealPlanList.
func BuildFakeMealPlanList() *types.QueryFilteredResult[types.MealPlan] {
	var examples []*types.MealPlan
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlan())
	}

	return &types.QueryFilteredResult[types.MealPlan]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeMealPlanUpdateRequestInput builds a faked MealPlanUpdateRequestInput from a meal plan.
func BuildFakeMealPlanUpdateRequestInput() *types.MealPlanUpdateRequestInput {
	mealPlan := BuildFakeMealPlan()
	return converters.ConvertMealPlanToMealPlanUpdateRequestInput(mealPlan)
}

// BuildFakeMealPlanCreationRequestInput builds a faked MealPlanCreationRequestInput.
func BuildFakeMealPlanCreationRequestInput() *types.MealPlanCreationRequestInput {
	mealPlan := BuildFakeMealPlan()
	return converters.ConvertMealPlanToMealPlanCreationRequestInput(mealPlan)
}
