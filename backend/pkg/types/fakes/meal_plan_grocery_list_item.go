package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

func BuildFakeMealPlanGroceryListItem() *types.MealPlanGroceryListItem {
	minQty := buildFakeNumber()

	return &types.MealPlanGroceryListItem{
		ID:                       BuildFakeID(),
		BelongsToMealPlan:        BuildFakeID(),
		Ingredient:               *BuildFakeValidIngredient(),
		MeasurementUnit:          *BuildFakeValidMeasurementUnit(),
		MinimumQuantityNeeded:    float32(minQty),
		MaximumQuantityNeeded:    pointer.To(float32(minQty + 1)),
		QuantityPurchased:        nil,
		PurchasedMeasurementUnit: nil,
		PurchasedUPC:             nil,
		PurchasePrice:            nil,
		StatusExplanation:        buildUniqueString(),
		Status:                   types.MealPlanGroceryListItemStatusUnknown,
		CreatedAt:                BuildFakeTime(),
	}
}

func BuildFakeMealPlanGroceryListItemList() *types.QueryFilteredResult[types.MealPlanGroceryListItem] {
	var examples []*types.MealPlanGroceryListItem
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanGroceryListItem())
	}

	return &types.QueryFilteredResult[types.MealPlanGroceryListItem]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

func BuildFakeMealPlanGroceryListItemCreationRequestInput() *types.MealPlanGroceryListItemCreationRequestInput {
	mealPlanGroceryListItem := BuildFakeMealPlanGroceryListItem()
	return converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemCreationRequestInput(mealPlanGroceryListItem)
}

func BuildFakeMealPlanGroceryListItemUpdateRequestInput() *types.MealPlanGroceryListItemUpdateRequestInput {
	mealPlanGroceryListItem := BuildFakeMealPlanGroceryListItem()
	return converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemUpdateRequestInput(mealPlanGroceryListItem)
}
