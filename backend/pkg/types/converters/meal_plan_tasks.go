package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertMealPlanTaskCreationRequestInputToMealPlanTaskDatabaseCreationInput creates a MealPlanTaskDatabaseCreationInput from a MealPlanTaskCreationRequestInput.
func ConvertMealPlanTaskCreationRequestInputToMealPlanTaskDatabaseCreationInput(input *types.MealPlanTaskCreationRequestInput) *types.MealPlanTaskDatabaseCreationInput {
	x := &types.MealPlanTaskDatabaseCreationInput{
		ID:                  identifiers.New(),
		AssignedToUser:      input.AssignedToUser,
		CreationExplanation: input.CreationExplanation,
		StatusExplanation:   input.StatusExplanation,
		MealPlanOptionID:    input.MealPlanOptionID,
		RecipePrepTaskID:    input.RecipePrepTaskID,
	}

	return x
}

// ConvertMealPlanTaskToMealPlanTaskCreationRequestInput builds a MealPlanTaskCreationRequestInput.
func ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(x *types.MealPlanTask) *types.MealPlanTaskCreationRequestInput {
	return &types.MealPlanTaskCreationRequestInput{
		Status:              x.Status,
		StatusExplanation:   x.StatusExplanation,
		CreationExplanation: x.CreationExplanation,
	}
}

// ConvertMealPlanTaskToMealPlanTaskDatabaseCreationInput builds a MealPlanTaskDatabaseCreationInput.
func ConvertMealPlanTaskToMealPlanTaskDatabaseCreationInput(x *types.MealPlanTask) *types.MealPlanTaskDatabaseCreationInput {
	return &types.MealPlanTaskDatabaseCreationInput{
		ID:                  x.ID,
		AssignedToUser:      x.AssignedToUser,
		CreationExplanation: x.CreationExplanation,
		StatusExplanation:   x.StatusExplanation,
		MealPlanOptionID:    x.MealPlanOption.ID,
		RecipePrepTaskID:    x.RecipePrepTask.ID,
	}
}
