package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertRecipeToRecipeUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeToRecipeUpdateRequestInput(input *types.Recipe) *types.RecipeUpdateRequestInput {
	x := &types.RecipeUpdateRequestInput{
		Name:                     &input.Name,
		Slug:                     &input.Slug,
		Source:                   &input.Source,
		Description:              &input.Description,
		InspiredByRecipeID:       input.InspiredByRecipeID,
		YieldsComponentType:      &input.YieldsComponentType,
		SealOfApproval:           &input.SealOfApproval,
		MinimumEstimatedPortions: &input.MinimumEstimatedPortions,
		MaximumEstimatedPortions: input.MaximumEstimatedPortions,
		PortionName:              &input.PortionName,
		PluralPortionName:        &input.PluralPortionName,
		EligibleForMeals:         &input.EligibleForMeals,
	}

	return x
}

// ConvertRecipeCreationRequestInputToRecipeDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeCreationRequestInputToRecipeDatabaseCreationInput(input *types.RecipeCreationRequestInput) (*types.RecipeDatabaseCreationInput, error) {
	x := &types.RecipeDatabaseCreationInput{
		ID:                       identifiers.New(),
		AlsoCreateMeal:           input.AlsoCreateMeal,
		Name:                     input.Name,
		Slug:                     input.Slug,
		Source:                   input.Source,
		Description:              input.Description,
		InspiredByRecipeID:       input.InspiredByRecipeID,
		SealOfApproval:           input.SealOfApproval,
		MinimumEstimatedPortions: input.MinimumEstimatedPortions,
		MaximumEstimatedPortions: input.MaximumEstimatedPortions,
		PortionName:              input.PortionName,
		PluralPortionName:        input.PluralPortionName,
		EligibleForMeals:         input.EligibleForMeals,
		YieldsComponentType:      input.YieldsComponentType,
	}

	for _, step := range input.Steps {
		s := ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(step)
		s.BelongsToRecipe = x.ID
		x.Steps = append(x.Steps, s)
	}

	for _, task := range input.PrepTasks {
		prepTaskDatabaseCreationInput, err := ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput(x, task)
		if err != nil {
			return nil, err
		}
		prepTaskDatabaseCreationInput.BelongsToRecipe = x.ID
		x.PrepTasks = append(x.PrepTasks, prepTaskDatabaseCreationInput)
	}

	return x, nil
}

// ConvertRecipeToRecipeCreationRequestInput builds a RecipeCreationRequestInput from a recipe.
func ConvertRecipeToRecipeCreationRequestInput(input *types.Recipe) *types.RecipeCreationRequestInput {
	steps := []*types.RecipeStepCreationRequestInput{}
	for _, step := range input.Steps {
		steps = append(steps, ConvertRecipeStepToRecipeStepCreationRequestInput(step))
	}

	prepTasks := []*types.RecipePrepTaskWithinRecipeCreationRequestInput{}
	for _, prepTask := range input.PrepTasks {
		prepTasks = append(prepTasks, ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(input, prepTask))
	}

	return &types.RecipeCreationRequestInput{
		Name:                     input.Name,
		Slug:                     input.Slug,
		Source:                   input.Source,
		Description:              input.Description,
		InspiredByRecipeID:       input.InspiredByRecipeID,
		SealOfApproval:           input.SealOfApproval,
		MinimumEstimatedPortions: input.MinimumEstimatedPortions,
		MaximumEstimatedPortions: input.MaximumEstimatedPortions,
		PortionName:              input.PortionName,
		PluralPortionName:        input.PluralPortionName,
		Steps:                    steps,
		PrepTasks:                prepTasks,
		EligibleForMeals:         input.EligibleForMeals,
		YieldsComponentType:      input.YieldsComponentType,
	}
}

// ConvertRecipeToRecipeDatabaseCreationInput builds a RecipeDatabaseCreationInput from a recipe.
func ConvertRecipeToRecipeDatabaseCreationInput(input *types.Recipe) *types.RecipeDatabaseCreationInput {
	steps := []*types.RecipeStepDatabaseCreationInput{}
	for _, step := range input.Steps {
		steps = append(steps, ConvertRecipeStepToRecipeStepDatabaseCreationInput(step))
	}

	prepTasks := []*types.RecipePrepTaskDatabaseCreationInput{}
	for _, prepTask := range input.PrepTasks {
		prepTasks = append(prepTasks, ConvertRecipePrepTaskToRecipePrepTaskDatabaseCreationInput(prepTask))
	}

	return &types.RecipeDatabaseCreationInput{
		ID:                       input.ID,
		Name:                     input.Name,
		Slug:                     input.Slug,
		Source:                   input.Source,
		Description:              input.Description,
		InspiredByRecipeID:       input.InspiredByRecipeID,
		CreatedByUser:            input.CreatedByUser,
		SealOfApproval:           input.SealOfApproval,
		MinimumEstimatedPortions: input.MinimumEstimatedPortions,
		MaximumEstimatedPortions: input.MaximumEstimatedPortions,
		PortionName:              input.PortionName,
		PluralPortionName:        input.PluralPortionName,
		Steps:                    steps,
		PrepTasks:                prepTasks,
		EligibleForMeals:         input.EligibleForMeals,
		YieldsComponentType:      input.YieldsComponentType,
	}
}

func ConvertRecipeToRecipeSearchSubset(r *types.Recipe) *types.RecipeSearchSubset {
	x := &types.RecipeSearchSubset{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, step := range r.Steps {
		x.Steps = append(x.Steps, ConvertRecipeStepToRecipeStepSearchSubset(step))
	}

	return x
}
