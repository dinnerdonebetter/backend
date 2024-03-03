package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput creates a RecipeStepIngredientDatabaseCreationInput from a RecipeStepIngredientCreationRequestInput.
func ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(input *types.RecipeStepIngredientCreationRequestInput) *types.RecipeStepIngredientDatabaseCreationInput {
	x := &types.RecipeStepIngredientDatabaseCreationInput{
		ID:                              identifiers.New(),
		IngredientID:                    input.IngredientID,
		Name:                            input.Name,
		MeasurementUnitID:               input.MeasurementUnitID,
		MinimumQuantity:                 input.MinimumQuantity,
		MaximumQuantity:                 input.MaximumQuantity,
		QuantityNotes:                   input.QuantityNotes,
		IngredientNotes:                 input.IngredientNotes,
		Optional:                        input.Optional,
		OptionIndex:                     input.OptionIndex,
		ProductOfRecipeStepIndex:        input.ProductOfRecipeStepIndex,
		ProductOfRecipeStepProductIndex: input.ProductOfRecipeStepProductIndex,
		VesselIndex:                     input.VesselIndex,
		ToTaste:                         input.ToTaste,
		ProductPercentageToUse:          input.ProductPercentageToUse,
	}

	return x
}

// ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput creates a RecipeStepIngredientUpdateRequestInput from a RecipeStepIngredient.
func ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput(input *types.RecipeStepIngredient) *types.RecipeStepIngredientUpdateRequestInput {
	x := &types.RecipeStepIngredientUpdateRequestInput{
		IngredientID:           &input.Ingredient.ID,
		RecipeStepProductID:    input.RecipeStepProductID,
		Name:                   &input.Name,
		MeasurementUnitID:      &input.MeasurementUnit.ID,
		QuantityNotes:          &input.QuantityNotes,
		IngredientNotes:        &input.IngredientNotes,
		BelongsToRecipeStep:    &input.BelongsToRecipeStep,
		MinimumQuantity:        &input.MinimumQuantity,
		MaximumQuantity:        input.MaximumQuantity,
		Optional:               &input.Optional,
		OptionIndex:            &input.OptionIndex,
		VesselIndex:            input.VesselIndex,
		ToTaste:                &input.ToTaste,
		ProductPercentageToUse: input.ProductPercentageToUse,
	}

	return x
}

// ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput builds a RecipeStepIngredientCreationRequestInput from a RecipeStepIngredient.
func ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(input *types.RecipeStepIngredient) *types.RecipeStepIngredientCreationRequestInput {
	return &types.RecipeStepIngredientCreationRequestInput{
		Name:                   input.Name,
		Optional:               input.Optional,
		IngredientID:           &input.Ingredient.ID,
		MeasurementUnitID:      input.MeasurementUnit.ID,
		MinimumQuantity:        input.MinimumQuantity,
		MaximumQuantity:        input.MaximumQuantity,
		QuantityNotes:          input.QuantityNotes,
		IngredientNotes:        input.IngredientNotes,
		OptionIndex:            input.OptionIndex,
		VesselIndex:            input.VesselIndex,
		ToTaste:                input.ToTaste,
		ProductPercentageToUse: input.ProductPercentageToUse,
	}
}

// ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput builds a RecipeStepIngredientDatabaseCreationInput from a RecipeStepIngredient.
func ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(input *types.RecipeStepIngredient) *types.RecipeStepIngredientDatabaseCreationInput {
	return &types.RecipeStepIngredientDatabaseCreationInput{
		ID:                     input.ID,
		Name:                   input.Name,
		Optional:               input.Optional,
		IngredientID:           &input.Ingredient.ID,
		MeasurementUnitID:      input.MeasurementUnit.ID,
		MinimumQuantity:        input.MinimumQuantity,
		MaximumQuantity:        input.MaximumQuantity,
		QuantityNotes:          input.QuantityNotes,
		IngredientNotes:        input.IngredientNotes,
		BelongsToRecipeStep:    input.BelongsToRecipeStep,
		OptionIndex:            input.OptionIndex,
		VesselIndex:            input.VesselIndex,
		ToTaste:                input.ToTaste,
		ProductPercentageToUse: input.ProductPercentageToUse,
	}
}
