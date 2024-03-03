package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertRecipeRatingToRecipeRatingUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeRatingToRecipeRatingUpdateRequestInput(x *types.RecipeRating) *types.RecipeRatingUpdateRequestInput {
	out := &types.RecipeRatingUpdateRequestInput{
		RecipeID:     &x.RecipeID,
		Taste:        &x.Taste,
		Difficulty:   &x.Difficulty,
		Cleanup:      &x.Cleanup,
		Instructions: &x.Instructions,
		Overall:      &x.Overall,
		Notes:        &x.Notes,
		ByUser:       &x.ByUser,
	}

	return out
}

// ConvertRecipeRatingCreationRequestInputToRecipeRatingDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeRatingCreationRequestInputToRecipeRatingDatabaseCreationInput(x *types.RecipeRatingCreationRequestInput) *types.RecipeRatingDatabaseCreationInput {
	out := &types.RecipeRatingDatabaseCreationInput{
		ID:           identifiers.New(),
		RecipeID:     x.RecipeID,
		Notes:        x.Notes,
		ByUser:       x.ByUser,
		Taste:        x.Taste,
		Difficulty:   x.Difficulty,
		Cleanup:      x.Cleanup,
		Instructions: x.Instructions,
		Overall:      x.Overall,
	}

	return out
}

// ConvertRecipeRatingToRecipeRatingCreationRequestInput builds a RecipeRatingCreationRequestInput from a Ingredient.
func ConvertRecipeRatingToRecipeRatingCreationRequestInput(x *types.RecipeRating) *types.RecipeRatingCreationRequestInput {
	return &types.RecipeRatingCreationRequestInput{
		RecipeID:     x.RecipeID,
		Notes:        x.Notes,
		ByUser:       x.ByUser,
		Taste:        x.Taste,
		Difficulty:   x.Difficulty,
		Cleanup:      x.Cleanup,
		Instructions: x.Instructions,
		Overall:      x.Overall,
	}
}

// ConvertRecipeRatingToRecipeRatingDatabaseCreationInput builds a RecipeRatingDatabaseCreationInput from a RecipeRating.
func ConvertRecipeRatingToRecipeRatingDatabaseCreationInput(x *types.RecipeRating) *types.RecipeRatingDatabaseCreationInput {
	return &types.RecipeRatingDatabaseCreationInput{
		ID:           x.ID,
		RecipeID:     x.RecipeID,
		Notes:        x.Notes,
		ByUser:       x.ByUser,
		Taste:        x.Taste,
		Difficulty:   x.Difficulty,
		Cleanup:      x.Cleanup,
		Instructions: x.Instructions,
		Overall:      x.Overall,
	}
}
