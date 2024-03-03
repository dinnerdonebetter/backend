package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeRecipeRating builds a faked valid ingredient.
func BuildFakeRecipeRating() *types.RecipeRating {
	return &types.RecipeRating{
		CreatedAt:    BuildFakeTime(),
		Notes:        buildUniqueString(),
		ID:           buildUniqueString(),
		RecipeID:     buildUniqueString(),
		ByUser:       buildUniqueString(),
		Taste:        float32(buildFakeNumber()),
		Instructions: float32(buildFakeNumber()),
		Overall:      float32(buildFakeNumber()),
		Cleanup:      float32(buildFakeNumber()),
		Difficulty:   float32(buildFakeNumber()),
	}
}

// BuildFakeRecipeRatingList builds a faked RecipeRatingList.
func BuildFakeRecipeRatingList() *types.QueryFilteredResult[types.RecipeRating] {
	var examples []*types.RecipeRating
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeRating())
	}

	return &types.QueryFilteredResult[types.RecipeRating]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeRatingUpdateRequestInput builds a faked RecipeRatingUpdateRequestInput from a valid ingredient.
func BuildFakeRecipeRatingUpdateRequestInput() *types.RecipeRatingUpdateRequestInput {
	recipeRating := BuildFakeRecipeRating()
	return converters.ConvertRecipeRatingToRecipeRatingUpdateRequestInput(recipeRating)
}

// BuildFakeRecipeRatingCreationRequestInput builds a faked RecipeRatingCreationRequestInput.
func BuildFakeRecipeRatingCreationRequestInput() *types.RecipeRatingCreationRequestInput {
	recipeRating := BuildFakeRecipeRating()
	return converters.ConvertRecipeRatingToRecipeRatingCreationRequestInput(recipeRating)
}
