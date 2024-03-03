package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput creates a RecipeStepInstrumentUpdateRequestInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput(input *types.RecipeStepInstrument) *types.RecipeStepInstrumentUpdateRequestInput {
	x := &types.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID:        &input.Instrument.ID,
		Notes:               &input.Notes,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                &input.Name,
		PreferenceRank:      &input.PreferenceRank,
		BelongsToRecipeStep: &input.BelongsToRecipeStep,
		Optional:            &input.Optional,
		OptionIndex:         &input.OptionIndex,
		MinimumQuantity:     &input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
	}

	return x
}

// ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput creates a RecipeStepInstrumentDatabaseCreationInput from a RecipeStepInstrumentCreationRequestInput.
func ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput(input *types.RecipeStepInstrumentCreationRequestInput) *types.RecipeStepInstrumentDatabaseCreationInput {
	x := &types.RecipeStepInstrumentDatabaseCreationInput{
		ID:                              identifiers.New(),
		InstrumentID:                    input.InstrumentID,
		RecipeStepProductID:             input.RecipeStepProductID,
		Name:                            input.Name,
		Notes:                           input.Notes,
		PreferenceRank:                  input.PreferenceRank,
		Optional:                        input.Optional,
		OptionIndex:                     input.OptionIndex,
		MinimumQuantity:                 input.MinimumQuantity,
		MaximumQuantity:                 input.MaximumQuantity,
		ProductOfRecipeStepIndex:        input.ProductOfRecipeStepIndex,
		ProductOfRecipeStepProductIndex: input.ProductOfRecipeStepProductIndex,
	}

	return x
}

// ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput builds a RecipeStepInstrumentCreationRequestInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(input *types.RecipeStepInstrument) *types.RecipeStepInstrumentCreationRequestInput {
	var instrumentID *string
	if input.Instrument != nil {
		instrumentID = &input.Instrument.ID
	}

	return &types.RecipeStepInstrumentCreationRequestInput{
		InstrumentID:        instrumentID,
		Name:                input.Name,
		RecipeStepProductID: input.RecipeStepProductID,
		Notes:               input.Notes,
		PreferenceRank:      input.PreferenceRank,
		Optional:            input.Optional,
		OptionIndex:         input.OptionIndex,
		MinimumQuantity:     input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
	}
}

// ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput builds a RecipeStepInstrumentDatabaseCreationInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput(input *types.RecipeStepInstrument) *types.RecipeStepInstrumentDatabaseCreationInput {
	var instrumentID *string
	if input.Instrument != nil {
		instrumentID = &input.Instrument.ID
	}

	return &types.RecipeStepInstrumentDatabaseCreationInput{
		ID:                  input.ID,
		InstrumentID:        instrumentID,
		Name:                input.Name,
		RecipeStepProductID: input.RecipeStepProductID,
		Notes:               input.Notes,
		PreferenceRank:      input.PreferenceRank,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Optional:            input.Optional,
		OptionIndex:         input.OptionIndex,
		MinimumQuantity:     input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
	}
}
