package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeRecipe builds a faked recipe.
func BuildFakeRecipe() *types.Recipe {
	recipeID := BuildFakeID()

	var steps []*types.RecipeStep
	for i := 0; i < exampleQuantity; i++ {
		step := BuildFakeRecipeStep()
		step.Index = uint32(i)
		step.BelongsToRecipe = recipeID
		steps = append(steps, step)
	}

	prepTasks := BuildFakeRecipePrepTaskList().Data
	for i := range prepTasks {
		prepTasks[i].BelongsToRecipe = recipeID
	}

	recipeMedia := BuildFakeRecipeMediaList().Data
	for i := range recipeMedia {
		recipeMedia[i].BelongsToRecipe = &recipeID
	}

	return &types.Recipe{
		ID:                       recipeID,
		Name:                     buildUniqueString(),
		Slug:                     buildUniqueString(),
		Source:                   buildUniqueString(),
		Description:              buildUniqueString(),
		InspiredByRecipeID:       nil,
		CreatedAt:                BuildFakeTime(),
		CreatedByUser:            BuildFakeID(),
		Steps:                    steps,
		PrepTasks:                prepTasks,
		SealOfApproval:           false,
		Media:                    recipeMedia,
		MinimumEstimatedPortions: float32(buildFakeNumber()),
		MaximumEstimatedPortions: pointer.To(float32(buildFakeNumber())),
		PortionName:              buildUniqueString(),
		PluralPortionName:        buildUniqueString(),
		EligibleForMeals:         true,
		YieldsComponentType:      "main",
		SupportingRecipes:        []*types.Recipe{},
	}
}

// BuildFakeRecipeList builds a faked RecipeList.
func BuildFakeRecipeList() *types.QueryFilteredResult[types.Recipe] {
	var examples []*types.Recipe
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipe())
	}

	return &types.QueryFilteredResult[types.Recipe]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeUpdateRequestInput builds a faked RecipeUpdateRequestInput from a recipe.
func BuildFakeRecipeUpdateRequestInput() *types.RecipeUpdateRequestInput {
	recipe := BuildFakeRecipe()
	return converters.ConvertRecipeToRecipeUpdateRequestInput(recipe)
}

// BuildFakeRecipeCreationRequestInput builds a faked RecipeCreationRequestInput.
func BuildFakeRecipeCreationRequestInput() *types.RecipeCreationRequestInput {
	exampleRecipe := BuildFakeRecipe()
	exampleCreationInput := converters.ConvertRecipeToRecipeCreationRequestInput(exampleRecipe)
	examplePrepTask := BuildFakeRecipePrepTask()
	examplePrepTaskInput := converters.ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(exampleRecipe, examplePrepTask)
	examplePrepTaskInput.TaskSteps = []*types.RecipePrepTaskStepWithinRecipeCreationRequestInput{
		{
			BelongsToRecipeStepIndex: exampleCreationInput.Steps[0].Index,
			SatisfiesRecipeStep:      false,
		},
	}
	exampleCreationInput.PrepTasks = []*types.RecipePrepTaskWithinRecipeCreationRequestInput{
		examplePrepTaskInput,
	}

	return exampleCreationInput
}
