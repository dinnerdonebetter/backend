package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// BuildFakeMealPlan builds a faked meal plan.
func BuildFakeMealPlan() *types.MealPlan {
	return &types.MealPlan{
		ID:                 ksuid.New().String(),
		State:              fake.LoremIpsumSentence(exampleQuantity),
		StartsAt:           uint64(fake.Uint32()),
		EndsAt:             uint64(fake.Uint32()),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
		BelongsToHousehold: fake.UUID(),
	}
}

// BuildFakeMealPlanList builds a faked MealPlanList.
func BuildFakeMealPlanList() *types.MealPlanList {
	var examples []*types.MealPlan
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlan())
	}

	return &types.MealPlanList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		MealPlans: examples,
	}
}

// BuildFakeMealPlanUpdateRequestInput builds a faked MealPlanUpdateRequestInput from a meal plan.
func BuildFakeMealPlanUpdateRequestInput() *types.MealPlanUpdateRequestInput {
	mealPlan := BuildFakeMealPlan()
	return &types.MealPlanUpdateRequestInput{
		State:              mealPlan.State,
		StartsAt:           mealPlan.StartsAt,
		EndsAt:             mealPlan.EndsAt,
		BelongsToHousehold: mealPlan.BelongsToHousehold,
	}
}

// BuildFakeMealPlanUpdateRequestInputFromMealPlan builds a faked MealPlanUpdateRequestInput from a meal plan.
func BuildFakeMealPlanUpdateRequestInputFromMealPlan(mealPlan *types.MealPlan) *types.MealPlanUpdateRequestInput {
	return &types.MealPlanUpdateRequestInput{
		State:              mealPlan.State,
		StartsAt:           mealPlan.StartsAt,
		EndsAt:             mealPlan.EndsAt,
		BelongsToHousehold: mealPlan.BelongsToHousehold,
	}
}

// BuildFakeMealPlanCreationRequestInput builds a faked MealPlanCreationRequestInput.
func BuildFakeMealPlanCreationRequestInput() *types.MealPlanCreationRequestInput {
	mealPlan := BuildFakeMealPlan()
	return BuildFakeMealPlanCreationRequestInputFromMealPlan(mealPlan)
}

// BuildFakeMealPlanCreationRequestInputFromMealPlan builds a faked MealPlanCreationRequestInput from a meal plan.
func BuildFakeMealPlanCreationRequestInputFromMealPlan(mealPlan *types.MealPlan) *types.MealPlanCreationRequestInput {
	return &types.MealPlanCreationRequestInput{
		ID:                 mealPlan.ID,
		State:              mealPlan.State,
		StartsAt:           mealPlan.StartsAt,
		EndsAt:             mealPlan.EndsAt,
		BelongsToHousehold: mealPlan.BelongsToHousehold,
	}
}

// BuildFakeMealPlanDatabaseCreationInput builds a faked MealPlanDatabaseCreationInput.
func BuildFakeMealPlanDatabaseCreationInput() *types.MealPlanDatabaseCreationInput {
	mealPlan := BuildFakeMealPlan()
	return BuildFakeMealPlanDatabaseCreationInputFromMealPlan(mealPlan)
}

// BuildFakeMealPlanDatabaseCreationInputFromMealPlan builds a faked MealPlanDatabaseCreationInput from a meal plan.
func BuildFakeMealPlanDatabaseCreationInputFromMealPlan(mealPlan *types.MealPlan) *types.MealPlanDatabaseCreationInput {
	return &types.MealPlanDatabaseCreationInput{
		ID:                 mealPlan.ID,
		State:              mealPlan.State,
		StartsAt:           mealPlan.StartsAt,
		EndsAt:             mealPlan.EndsAt,
		BelongsToHousehold: mealPlan.BelongsToHousehold,
	}
}
