package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeMealPlanTask builds a faked meal plan task.
func BuildFakeMealPlanTask() *types.MealPlanTask {
	return &types.MealPlanTask{
		ID:                  BuildFakeID(),
		CreatedAt:           BuildFakeTime(),
		Status:              "unfinished",
		StatusExplanation:   buildUniqueString(),
		CreationExplanation: buildUniqueString(),
		CompletedAt:         nil,
		RecipePrepTask:      *BuildFakeRecipePrepTask(),
	}
}

// BuildFakeMealPlanTaskCreationRequestInput builds a faked meal plan task.
func BuildFakeMealPlanTaskCreationRequestInput() *types.MealPlanTaskCreationRequestInput {
	x := BuildFakeMealPlanTask()

	return converters.ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(x)
}

// BuildFakeMealPlanTaskList builds a faked MealPlanTaskList.
func BuildFakeMealPlanTaskList() *types.QueryFilteredResult[types.MealPlanTask] {
	var examples []*types.MealPlanTask
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanTask())
	}

	return &types.QueryFilteredResult[types.MealPlanTask]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeMealPlanTaskDatabaseCreationInputs builds a faked MealPlanTaskList.
func BuildFakeMealPlanTaskDatabaseCreationInputs() []*types.MealPlanTaskDatabaseCreationInput {
	var examples []*types.MealPlanTaskDatabaseCreationInput
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, &types.MealPlanTaskDatabaseCreationInput{
			MealPlanOptionID:    "",
			ID:                  BuildFakeID(),
			StatusExplanation:   buildUniqueString(),
			CreationExplanation: buildUniqueString(),
		})
	}

	return examples
}

// BuildFakeMealPlanTaskStatusChangeRequestInput builds a faked meal plan task.
func BuildFakeMealPlanTaskStatusChangeRequestInput() *types.MealPlanTaskStatusChangeRequestInput {
	return &types.MealPlanTaskStatusChangeRequestInput{
		ID:                BuildFakeID(),
		Status:            pointer.To("unfinished"),
		StatusExplanation: buildUniqueString(),
	}
}
