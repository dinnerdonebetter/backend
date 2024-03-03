package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeRecipeStep builds a faked recipe step.
func BuildFakeRecipeStep() *types.RecipeStep {
	recipeStepID := BuildFakeID()

	minTemp := buildFakeNumber()

	var ingredients []*types.RecipeStepIngredient
	for i := 0; i < exampleQuantity; i++ {
		ing := BuildFakeRecipeStepIngredient()
		ing.BelongsToRecipeStep = recipeStepID
		ingredients = append(ingredients, ing)
	}

	var instruments []*types.RecipeStepInstrument
	for i := 0; i < exampleQuantity; i++ {
		ing := BuildFakeRecipeStepInstrument()
		ing.BelongsToRecipeStep = recipeStepID
		instruments = append(instruments, ing)
	}

	var vessels []*types.RecipeStepVessel
	for i := 0; i < exampleQuantity; i++ {
		ing := BuildFakeRecipeStepVessel()
		ing.BelongsToRecipeStep = recipeStepID
		vessels = append(vessels, ing)
	}

	var products []*types.RecipeStepProduct
	for i := 0; i < exampleQuantity; i++ {
		p := BuildFakeRecipeStepProduct()
		p.BelongsToRecipeStep = recipeStepID
		products = append(products, p)
	}

	completionConditionID := BuildFakeID()
	completionConditions := []*types.RecipeStepCompletionCondition{
		{
			ID:                  completionConditionID,
			BelongsToRecipeStep: recipeStepID,
			IngredientState:     types.ValidIngredientState{},
			Notes:               buildUniqueString(),
			Ingredients: []*types.RecipeStepCompletionConditionIngredient{
				{
					ID:                                     BuildFakeID(),
					BelongsToRecipeStepCompletionCondition: completionConditionID,
					RecipeStepIngredient:                   ingredients[0].ID,
				},
			},
			Optional: false,
		},
	}

	return &types.RecipeStep{
		ID:                            recipeStepID,
		Index:                         fake.Uint32(),
		Preparation:                   *BuildFakeValidPreparation(),
		MinimumEstimatedTimeInSeconds: func(x uint32) *uint32 { return &x }(fake.Uint32()),
		MaximumEstimatedTimeInSeconds: func(x uint32) *uint32 { return &x }(fake.Uint32()),
		MinimumTemperatureInCelsius:   pointer.To(float32(minTemp)),
		MaximumTemperatureInCelsius:   pointer.To(float32(minTemp + 1)),
		Notes:                         buildUniqueString(),
		Products:                      products,
		Optional:                      false,
		CreatedAt:                     BuildFakeTime(),
		BelongsToRecipe:               BuildFakeID(),
		Ingredients:                   ingredients,
		ExplicitInstructions:          buildUniqueString(),
		ConditionExpression:           buildUniqueString(),
		Instruments:                   instruments,
		Vessels:                       vessels,
		CompletionConditions:          completionConditions,
		StartTimerAutomatically:       fake.Bool(),
	}
}

// BuildFakeRecipeStepList builds a faked RecipeStepList.
func BuildFakeRecipeStepList() *types.QueryFilteredResult[types.RecipeStep] {
	var examples []*types.RecipeStep
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStep())
	}

	return &types.QueryFilteredResult[types.RecipeStep]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeStepUpdateRequestInput builds a faked RecipeStepUpdateRequestInput from a recipe step.
func BuildFakeRecipeStepUpdateRequestInput() *types.RecipeStepUpdateRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return converters.ConvertRecipeStepToRecipeStepUpdateRequestInput(recipeStep)
}

// BuildFakeRecipeStepCreationRequestInput builds a faked RecipeStepCreationRequestInput.
func BuildFakeRecipeStepCreationRequestInput() *types.RecipeStepCreationRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return converters.ConvertRecipeStepToRecipeStepCreationRequestInput(recipeStep)
}
