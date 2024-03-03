package converters

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput creates a RecipePrepTaskUpdateRequestInput from a RecipePrepTask.
func ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput(input *types.RecipePrepTask) *types.RecipePrepTaskUpdateRequestInput {
	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for _, x := range input.TaskSteps {
		y := x
		taskSteps = append(taskSteps, &types.RecipePrepTaskStepUpdateRequestInput{
			BelongsToRecipeStep:     &y.BelongsToRecipeStep,
			BelongsToRecipePrepTask: &y.BelongsToRecipePrepTask,
			SatisfiesRecipeStep:     &y.SatisfiesRecipeStep,
		})
	}
	x := &types.RecipePrepTaskUpdateRequestInput{
		Name:                                   &input.Name,
		Description:                            &input.Description,
		Notes:                                  &input.Notes,
		ExplicitStorageInstructions:            &input.ExplicitStorageInstructions,
		Optional:                               &input.Optional,
		MinimumTimeBufferBeforeRecipeInSeconds: &input.MinimumTimeBufferBeforeRecipeInSeconds,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		StorageType:                            &input.StorageType,
		MinimumStorageTemperatureInCelsius:     input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     input.MaximumStorageTemperatureInCelsius,
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
			BelongsToRecipeStep: x.BelongsToRecipeStep,
			SatisfiesRecipeStep: x.SatisfiesRecipeStep,
		})
	}

	x := &types.RecipePrepTaskDatabaseCreationInput{
		ID:                                     identifiers.New(),
		Name:                                   input.Name,
		Description:                            input.Description,
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		Optional:                               input.Optional,
		StorageType:                            input.StorageType,
		BelongsToRecipe:                        input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     input.MaximumStorageTemperatureInCelsius,
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
	}

	return x
}

// ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput(recipe *types.RecipeDatabaseCreationInput, input *types.RecipePrepTaskWithinRecipeCreationRequestInput) (*types.RecipePrepTaskDatabaseCreationInput, error) {
	x := &types.RecipePrepTaskDatabaseCreationInput{
		ID:                                     identifiers.New(),
		Name:                                   input.Name,
		Description:                            input.Description,
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		Optional:                               input.Optional,
		StorageType:                            input.StorageType,
		BelongsToRecipe:                        input.BelongsToRecipe,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     input.MaximumStorageTemperatureInCelsius,
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
	}

	x.TaskSteps = []*types.RecipePrepTaskStepDatabaseCreationInput{}
	for i, ts := range input.TaskSteps {
		if rs := recipe.FindStepByIndex(ts.BelongsToRecipeStepIndex); rs != nil {
			x.TaskSteps = append(x.TaskSteps, &types.RecipePrepTaskStepDatabaseCreationInput{
				ID:                      identifiers.New(),
				BelongsToRecipeStep:     rs.ID,
				BelongsToRecipePrepTask: x.ID,
				SatisfiesRecipeStep:     ts.SatisfiesRecipeStep,
			})
		} else {
			return nil, fmt.Errorf("task step #%d has an invalid recipe step index", i+1)
		}
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
		Name:                                   input.Name,
		Description:                            input.Description,
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		Optional:                               input.Optional,
		StorageType:                            input.StorageType,
		TaskSteps:                              taskSteps,
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     input.MaximumStorageTemperatureInCelsius,
		BelongsToRecipe:                        input.BelongsToRecipe,
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepCreationRequestInput(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepCreationRequestInput {
	return &types.RecipePrepTaskStepCreationRequestInput{
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		SatisfiesRecipeStep: input.SatisfiesRecipeStep,
	}
}

func ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(input *types.RecipePrepTask) *types.RecipePrepTaskCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepCreationRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepCreationRequestInput(x))
	}

	return &types.RecipePrepTaskCreationRequestInput{
		Name:                                   input.Name,
		Description:                            input.Description,
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		Optional:                               input.Optional,
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
		Name:                                   input.Name,
		Description:                            input.Description,
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		Optional:                               input.Optional,
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
		SatisfiesRecipeStep:      input.SatisfiesRecipeStep,
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepUpdateRequestInput(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepUpdateRequestInput {
	return &types.RecipePrepTaskStepUpdateRequestInput{
		BelongsToRecipeStep:     pointer.To(input.BelongsToRecipeStep),
		BelongsToRecipePrepTask: pointer.To(input.BelongsToRecipePrepTask),
		SatisfiesRecipeStep:     pointer.To(input.SatisfiesRecipeStep),
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
