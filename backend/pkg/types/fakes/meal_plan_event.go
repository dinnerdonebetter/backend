package fakes

import (
	"time"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeMealPlanEvent builds a faked meal plan event.
func BuildFakeMealPlanEvent() *types.MealPlanEvent {
	mealPlanEventID := BuildFakeID()

	now := time.Now().Add(0).Truncate(time.Second).UTC()
	inTenMinutes := now.Add(time.Minute * 10).Add(0).Truncate(time.Second).UTC()
	inOneWeek := now.Add((time.Hour * 24) * 7).Add(0).Truncate(time.Second).UTC()

	options := []*types.MealPlanOption{}
	for _, opt := range BuildFakeMealPlanOptionList().Data {
		opt.BelongsToMealPlanEvent = mealPlanEventID
		options = append(options, opt)
	}

	return &types.MealPlanEvent{
		ID:       mealPlanEventID,
		Notes:    buildUniqueString(),
		StartsAt: inTenMinutes,
		EndsAt:   inOneWeek,
		MealName: fake.RandomString([]string{
			types.BreakfastMealName,
			types.SecondBreakfastMealName,
			types.BrunchMealName,
			types.LunchMealName,
			types.SupperMealName,
			types.DinnerMealName,
		}),
		CreatedAt:         BuildFakeTime(),
		BelongsToMealPlan: BuildFakeID(),
		Options:           options,
	}
}

// BuildFakeMealPlanEventList builds a faked MealPlanEventList.
func BuildFakeMealPlanEventList() *types.QueryFilteredResult[types.MealPlanEvent] {
	var examples []*types.MealPlanEvent
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanEvent())
	}

	return &types.QueryFilteredResult[types.MealPlanEvent]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeMealPlanEventUpdateRequestInput builds a faked MealPlanEventUpdateRequestInput from a meal plan.
func BuildFakeMealPlanEventUpdateRequestInput() *types.MealPlanEventUpdateRequestInput {
	mealPlanEvent := BuildFakeMealPlanEvent()
	return &types.MealPlanEventUpdateRequestInput{
		Notes:             &mealPlanEvent.Notes,
		StartsAt:          &mealPlanEvent.StartsAt,
		EndsAt:            &mealPlanEvent.EndsAt,
		MealName:          &mealPlanEvent.MealName,
		BelongsToMealPlan: mealPlanEvent.BelongsToMealPlan,
	}
}

// BuildFakeMealPlanEventCreationRequestInput builds a faked MealPlanEventCreationRequestInput.
func BuildFakeMealPlanEventCreationRequestInput() *types.MealPlanEventCreationRequestInput {
	mealPlan := BuildFakeMealPlanEvent()
	return converters.ConvertMealPlanEventToMealPlanEventCreationRequestInput(mealPlan)
}
