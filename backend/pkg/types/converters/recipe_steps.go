package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertRecipeStepToRecipeStepUpdateRequestInput creates a RecipeStepUpdateRequestInput from a RecipeStep.
func ConvertRecipeStepToRecipeStepUpdateRequestInput(input *types.RecipeStep) *types.RecipeStepUpdateRequestInput {
	x := &types.RecipeStepUpdateRequestInput{
		MinimumTemperatureInCelsius:   input.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   input.MaximumTemperatureInCelsius,
		Notes:                         &input.Notes,
		BelongsToRecipe:               input.BelongsToRecipe,
		Preparation:                   &input.Preparation,
		Index:                         &input.Index,
		MinimumEstimatedTimeInSeconds: input.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: input.MaximumEstimatedTimeInSeconds,
		Optional:                      &input.Optional,
		ExplicitInstructions:          &input.ExplicitInstructions,
		ConditionExpression:           &input.ConditionExpression,
		StartTimerAutomatically:       &input.StartTimerAutomatically,
	}

	return x
}

// ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput creates a RecipeStepDatabaseCreationInput from a RecipeStepCreationRequestInput.
func ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(input *types.RecipeStepCreationRequestInput) *types.RecipeStepDatabaseCreationInput {
	x := &types.RecipeStepDatabaseCreationInput{
		ID:                            identifiers.New(),
		Index:                         input.Index,
		PreparationID:                 input.PreparationID,
		MinimumEstimatedTimeInSeconds: input.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: input.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   input.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   input.MaximumTemperatureInCelsius,
		Notes:                         input.Notes,
		Optional:                      input.Optional,
		ExplicitInstructions:          input.ExplicitInstructions,
		ConditionExpression:           input.ConditionExpression,
		StartTimerAutomatically:       input.StartTimerAutomatically,
	}

	x.Ingredients = []*types.RecipeStepIngredientDatabaseCreationInput{}
	for _, ingredient := range input.Ingredients {
		convertedIngredient := ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(ingredient)
		convertedIngredient.ID = identifiers.New()
		convertedIngredient.BelongsToRecipeStep = x.ID
		x.Ingredients = append(x.Ingredients, convertedIngredient)
	}

	x.Instruments = []*types.RecipeStepInstrumentDatabaseCreationInput{}
	for _, instrument := range input.Instruments {
		convertedInstrument := ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput(instrument)
		convertedInstrument.ID = identifiers.New()
		convertedInstrument.BelongsToRecipeStep = x.ID
		x.Instruments = append(x.Instruments, convertedInstrument)
	}

	x.Vessels = []*types.RecipeStepVesselDatabaseCreationInput{}
	for _, vessel := range input.Vessels {
		convertedVessel := ConvertRecipeStepVesselCreationRequestInputToRecipeStepVesselDatabaseCreationInput(vessel)
		convertedVessel.ID = identifiers.New()
		convertedVessel.BelongsToRecipeStep = x.ID
		x.Vessels = append(x.Vessels, convertedVessel)
	}

	x.Products = []*types.RecipeStepProductDatabaseCreationInput{}
	for _, product := range input.Products {
		convertedProduct := ConvertRecipeStepProductCreationInputToRecipeStepProductDatabaseCreationInput(product)
		convertedProduct.ID = identifiers.New()
		convertedProduct.BelongsToRecipeStep = x.ID
		x.Products = append(x.Products, convertedProduct)
	}

	x.CompletionConditions = []*types.RecipeStepCompletionConditionDatabaseCreationInput{}
	for _, product := range input.CompletionConditions {
		convertedCompletionCondition := ConvertRecipeStepCompletionConditionCreationRequestInputToRecipeStepCompletionConditionDatabaseCreationInput(x, product)
		convertedCompletionCondition.ID = identifiers.New()
		convertedCompletionCondition.BelongsToRecipeStep = x.ID
		x.CompletionConditions = append(x.CompletionConditions, convertedCompletionCondition)
	}

	return x
}

// ConvertRecipeStepToRecipeStepCreationRequestInput builds a RecipeStepCreationRequestInput from a RecipeStep.
func ConvertRecipeStepToRecipeStepCreationRequestInput(recipeStep *types.RecipeStep) *types.RecipeStepCreationRequestInput {
	ingredients := []*types.RecipeStepIngredientCreationRequestInput{}
	for _, ingredient := range recipeStep.Ingredients {
		ingredients = append(ingredients, ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(ingredient))
	}

	instruments := []*types.RecipeStepInstrumentCreationRequestInput{}
	for _, instrument := range recipeStep.Instruments {
		instruments = append(instruments, ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(instrument))
	}

	vessels := []*types.RecipeStepVesselCreationRequestInput{}
	for _, vessel := range recipeStep.Vessels {
		vessels = append(vessels, ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(vessel))
	}

	products := []*types.RecipeStepProductCreationRequestInput{}
	for _, product := range recipeStep.Products {
		products = append(products, ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(product))
	}

	completionConditions := []*types.RecipeStepCompletionConditionCreationRequestInput{}
	for _, completionCondition := range recipeStep.CompletionConditions {
		completionConditions = append(completionConditions, ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionCreationRequestInput(recipeStep, completionCondition))
	}

	return &types.RecipeStepCreationRequestInput{
		Optional:                      recipeStep.Optional,
		Index:                         recipeStep.Index,
		PreparationID:                 recipeStep.Preparation.ID,
		MinimumEstimatedTimeInSeconds: recipeStep.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: recipeStep.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   recipeStep.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   recipeStep.MaximumTemperatureInCelsius,
		Notes:                         recipeStep.Notes,
		ExplicitInstructions:          recipeStep.ExplicitInstructions,
		ConditionExpression:           recipeStep.ConditionExpression,
		StartTimerAutomatically:       recipeStep.StartTimerAutomatically,
		Products:                      products,
		Ingredients:                   ingredients,
		Instruments:                   instruments,
		Vessels:                       vessels,
		CompletionConditions:          completionConditions,
	}
}

// ConvertRecipeStepToRecipeStepDatabaseCreationInput builds a RecipeStepDatabaseCreationInput from a RecipeStep.
func ConvertRecipeStepToRecipeStepDatabaseCreationInput(recipeStep *types.RecipeStep) *types.RecipeStepDatabaseCreationInput {
	ingredients := []*types.RecipeStepIngredientDatabaseCreationInput{}
	for _, i := range recipeStep.Ingredients {
		ingredients = append(ingredients, ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(i))
	}

	instruments := []*types.RecipeStepInstrumentDatabaseCreationInput{}
	for _, i := range recipeStep.Instruments {
		instruments = append(instruments, ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput(i))
	}

	vessels := []*types.RecipeStepVesselDatabaseCreationInput{}
	for _, v := range recipeStep.Vessels {
		vessels = append(vessels, ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(v))
	}

	products := []*types.RecipeStepProductDatabaseCreationInput{}
	for _, p := range recipeStep.Products {
		products = append(products, ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput(p))
	}

	completionConditions := []*types.RecipeStepCompletionConditionDatabaseCreationInput{}
	for _, completionCondition := range recipeStep.CompletionConditions {
		completionConditions = append(completionConditions, ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput(completionCondition))
	}

	return &types.RecipeStepDatabaseCreationInput{
		ID:                            recipeStep.ID,
		Index:                         recipeStep.Index,
		PreparationID:                 recipeStep.Preparation.ID,
		Optional:                      recipeStep.Optional,
		MinimumEstimatedTimeInSeconds: recipeStep.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: recipeStep.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   recipeStep.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   recipeStep.MaximumTemperatureInCelsius,
		StartTimerAutomatically:       recipeStep.StartTimerAutomatically,
		Notes:                         recipeStep.Notes,
		ExplicitInstructions:          recipeStep.ExplicitInstructions,
		ConditionExpression:           recipeStep.ConditionExpression,
		Ingredients:                   ingredients,
		Instruments:                   instruments,
		Products:                      products,
		Vessels:                       vessels,
		BelongsToRecipe:               recipeStep.BelongsToRecipe,
		CompletionConditions:          completionConditions,
	}
}

func ConvertRecipeStepToRecipeStepSearchSubset(x *types.RecipeStep) *types.RecipeStepSearchSubset {
	stepSubset := &types.RecipeStepSearchSubset{
		Preparation: x.Preparation.Name,
	}

	for _, ingredient := range x.Ingredients {
		stepSubset.Ingredients = append(stepSubset.Ingredients, types.NamedID{ID: ingredient.ID, Name: ingredient.Name})
	}

	for _, instrument := range x.Instruments {
		stepSubset.Instruments = append(stepSubset.Instruments, types.NamedID{ID: instrument.ID, Name: instrument.Name})
	}

	for _, vessel := range x.Vessels {
		stepSubset.Vessels = append(stepSubset.Vessels, types.NamedID{ID: vessel.ID, Name: vessel.Name})
	}

	return stepSubset
}
