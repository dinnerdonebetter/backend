package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepProduct builds a faked recipe step product.
func BuildFakeRecipeStepProduct() *types.RecipeStepProduct {
	return &types.RecipeStepProduct{
		ID:                  uint64(fake.Uint32()),
		ExternalID:          fake.UUID(),
		Name:                fake.Word(),
		QuantityType:        fake.Word(),
		QuantityValue:       fake.Float32(),
		QuantityNotes:       fake.Word(),
		RecipeStepID:        uint64(fake.Uint32()),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.Uint64(),
	}
}

// BuildFakeRecipeStepProductList builds a faked RecipeStepProductList.
func BuildFakeRecipeStepProductList() *types.RecipeStepProductList {
	var examples []*types.RecipeStepProduct
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepProduct())
	}

	return &types.RecipeStepProductList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		RecipeStepProducts: examples,
	}
}

// BuildFakeRecipeStepProductUpdateInput builds a faked RecipeStepProductUpdateInput from a recipe step product.
func BuildFakeRecipeStepProductUpdateInput() *types.RecipeStepProductUpdateInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return &types.RecipeStepProductUpdateInput{
		Name:                recipeStepProduct.Name,
		QuantityType:        recipeStepProduct.QuantityType,
		QuantityValue:       recipeStepProduct.QuantityValue,
		QuantityNotes:       recipeStepProduct.QuantityNotes,
		RecipeStepID:        recipeStepProduct.RecipeStepID,
		BelongsToRecipeStep: recipeStepProduct.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct builds a faked RecipeStepProductUpdateInput from a recipe step product.
func BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductUpdateInput {
	return &types.RecipeStepProductUpdateInput{
		Name:                recipeStepProduct.Name,
		QuantityType:        recipeStepProduct.QuantityType,
		QuantityValue:       recipeStepProduct.QuantityValue,
		QuantityNotes:       recipeStepProduct.QuantityNotes,
		RecipeStepID:        recipeStepProduct.RecipeStepID,
		BelongsToRecipeStep: recipeStepProduct.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepProductCreationInput builds a faked RecipeStepProductCreationInput.
func BuildFakeRecipeStepProductCreationInput() *types.RecipeStepProductCreationInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(recipeStepProduct)
}

// BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct builds a faked RecipeStepProductCreationInput from a recipe step product.
func BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductCreationInput {
	return &types.RecipeStepProductCreationInput{
		Name:                recipeStepProduct.Name,
		QuantityType:        recipeStepProduct.QuantityType,
		QuantityValue:       recipeStepProduct.QuantityValue,
		QuantityNotes:       recipeStepProduct.QuantityNotes,
		RecipeStepID:        recipeStepProduct.RecipeStepID,
		BelongsToRecipeStep: recipeStepProduct.BelongsToRecipeStep,
	}
}
