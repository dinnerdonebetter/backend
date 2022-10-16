package converters

import (
	"fmt"

	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput creates a RecipePrepTaskUpdateRequestInput from a RecipePrepTask.
func ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput(input *types.RecipePrepTask) *types.RecipePrepTaskUpdateRequestInput {
	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, &types.RecipePrepTaskStepUpdateRequestInput{
			ID:                      x.ID,
			BelongsToRecipeStep:     &x.BelongsToRecipeStep,
			BelongsToRecipePrepTask: &x.BelongsToRecipePrepTask,
			SatisfiesRecipeStep:     &x.SatisfiesRecipeStep,
		})
	}
	x := &types.RecipePrepTaskUpdateRequestInput{
		Notes:                                  &input.Notes,
		ExplicitStorageInstructions:            &input.ExplicitStorageInstructions,
		MinimumTimeBufferBeforeRecipeInSeconds: &input.MinimumTimeBufferBeforeRecipeInSeconds,
		MaximumTimeBufferBeforeRecipeInSeconds: &input.MaximumTimeBufferBeforeRecipeInSeconds,
		StorageType:                            &input.StorageType,
		MinimumStorageTemperatureInCelsius:     &input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     &input.MaximumStorageTemperatureInCelsius,
		BelongsToRecipe:                        &input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
	}

	return x
}

// ConvertRecipePrepTaskCreationRequestInputToRecipePrepTaskDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipePrepTaskCreationRequestInputToRecipePrepTaskDatabaseCreationInput(input *types.RecipePrepTaskCreationRequestInput) *types.RecipePrepTaskDatabaseCreationInput {
	taskSteps := []*types.RecipePrepTaskStepDatabaseCreationInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, &types.RecipePrepTaskStepDatabaseCreationInput{
			BelongsToRecipeStep:     x.BelongsToRecipeStep,
			BelongsToRecipePrepTask: x.BelongsToRecipePrepTask,
			SatisfiesRecipeStep:     x.SatisfiesRecipeStep,
		})
	}

	x := &types.RecipePrepTaskDatabaseCreationInput{
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		StorageType:                            input.StorageType,
		BelongsToRecipe:                        input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     uint32(input.MinimumStorageTemperatureInCelsius * types.RecipePrepTaskStorageTemperatureModifier),
		MaximumStorageTemperatureInCelsius:     uint32(input.MaximumStorageTemperatureInCelsius * types.RecipePrepTaskStorageTemperatureModifier),
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
	}

	return x
}

// ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput(recipe *types.RecipeDatabaseCreationInput, input *types.RecipePrepTaskWithinRecipeCreationRequestInput) (*types.RecipePrepTaskDatabaseCreationInput, error) {
	taskSteps := []*types.RecipePrepTaskStepDatabaseCreationInput{}
	for i, x := range input.TaskSteps {
		if y := recipe.FindStepByIndex(x.BelongsToRecipeStepIndex); y != nil {
			taskSteps = append(taskSteps, &types.RecipePrepTaskStepDatabaseCreationInput{
				BelongsToRecipeStep:     y.ID,
				BelongsToRecipePrepTask: x.BelongsToRecipePrepTask,
				SatisfiesRecipeStep:     x.SatisfiesRecipeStep,
			})
		} else {
			return nil, fmt.Errorf("task step #%d has an invalid recipe step index", i+1)
		}
	}

	x := &types.RecipePrepTaskDatabaseCreationInput{
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		StorageType:                            input.StorageType,
		BelongsToRecipe:                        input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     uint32(input.MinimumStorageTemperatureInCelsius * types.RecipePrepTaskStorageTemperatureModifier),
		MaximumStorageTemperatureInCelsius:     uint32(input.MaximumStorageTemperatureInCelsius * types.RecipePrepTaskStorageTemperatureModifier),
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
	}

	return x, nil
}

// ConvertRecipePrepTaskToRecipePrepTaskDatabaseCreationInput builds a RecipePrepTaskDatabaseCreationInput from a recipe prep task.
func ConvertRecipePrepTaskToRecipePrepTaskDatabaseCreationInput(input *types.RecipePrepTask) *types.RecipePrepTaskDatabaseCreationInput {
	taskSteps := []*types.RecipePrepTaskStepDatabaseCreationInput{}
	for _, step := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepDatabaseCreationInput(step))
	}

	return &types.RecipePrepTaskDatabaseCreationInput{
		ID:                                     input.ID,
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		StorageType:                            input.StorageType,
		TaskSteps:                              taskSteps,
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     uint32(input.MinimumStorageTemperatureInCelsius * types.RecipePrepTaskStorageTemperatureModifier),
		MaximumStorageTemperatureInCelsius:     uint32(input.MaximumStorageTemperatureInCelsius * types.RecipePrepTaskStorageTemperatureModifier),
		BelongsToRecipe:                        input.BelongsToRecipe,
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepCreationRequestInput(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepCreationRequestInput {
	return &types.RecipePrepTaskStepCreationRequestInput{
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}
}

func ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(input *types.RecipePrepTask) *types.RecipePrepTaskCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepCreationRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepCreationRequestInput(x))
	}

	return &types.RecipePrepTaskCreationRequestInput{
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		StorageType:                            input.StorageType,
		BelongsToRecipe:                        input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     input.MaximumStorageTemperatureInCelsius,
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
	}
}

func ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(recipe *types.Recipe, input *types.RecipePrepTask) *types.RecipePrepTaskWithinRecipeCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepWithinRecipeCreationRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepWithinRecipeCreationRequestInput(recipe, x))
	}

	return &types.RecipePrepTaskWithinRecipeCreationRequestInput{
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		StorageType:                            input.StorageType,
		BelongsToRecipe:                        input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     input.MaximumStorageTemperatureInCelsius,
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepWithinRecipeCreationRequestInput(recipe *types.Recipe, input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepWithinRecipeCreationRequestInput {
	var belongsToIndex uint32
	if x := recipe.FindStepByID(input.BelongsToRecipeStep); x != nil {
		belongsToIndex = x.Index
	}

	return &types.RecipePrepTaskStepWithinRecipeCreationRequestInput{
		BelongsToRecipeStepIndex: belongsToIndex,
		BelongsToRecipePrepTask:  input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:      input.SatisfiesRecipeStep,
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepUpdateRequestInput(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepUpdateRequestInput {
	return &types.RecipePrepTaskStepUpdateRequestInput{
		ID:                      input.ID,
		BelongsToRecipeStep:     pointers.StringPointer(input.BelongsToRecipeStep),
		BelongsToRecipePrepTask: pointers.StringPointer(input.BelongsToRecipePrepTask),
		SatisfiesRecipeStep:     pointers.BoolPointer(input.SatisfiesRecipeStep),
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepDatabaseCreationInput(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepDatabaseCreationInput {
	return &types.RecipePrepTaskStepDatabaseCreationInput{
		ID:                      input.ID,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}
}
