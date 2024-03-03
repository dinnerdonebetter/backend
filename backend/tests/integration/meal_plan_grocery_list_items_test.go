package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkMealPlanGroceryListItemEquality(t *testing.T, expected, actual *types.MealPlanGroceryListItem) {
	t.Helper()

	assert.NotZero(t, actual.ID)

	assert.Equal(t, expected.QuantityPurchased, actual.QuantityPurchased, "expected QuantityPurchased for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.QuantityPurchased, actual.QuantityPurchased)
	assert.Equal(t, expected.PurchasePrice, actual.PurchasePrice, "expected PurchasePrice for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.PurchasePrice, actual.PurchasePrice)
	assert.Equal(t, expected.PurchasedUPC, actual.PurchasedUPC, "expected PurchasedUPC for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.PurchasedUPC, actual.PurchasedUPC)
	assert.Equal(t, expected.PurchasedMeasurementUnit, actual.PurchasedMeasurementUnit, "expected PurchasedMeasurementUnit for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.PurchasedMeasurementUnit, actual.PurchasedMeasurementUnit)
	assert.Equal(t, expected.StatusExplanation, actual.StatusExplanation, "expected StatusExplanation for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.StatusExplanation, actual.StatusExplanation)
	assert.Equal(t, expected.Status, actual.Status, "expected Status for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.Status, actual.Status)
	assert.Equal(t, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID, "expected MeasurementUnitID for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID)
	assert.Equal(t, expected.BelongsToMealPlan, actual.BelongsToMealPlan, "expected BelongsToMealPlan for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.BelongsToMealPlan, actual.BelongsToMealPlan)
	assert.Equal(t, expected.Ingredient.ID, actual.Ingredient.ID, "expected ValidIngredientID for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.Ingredient.ID, actual.Ingredient.ID)
	assert.Equal(t, expected.MaximumQuantityNeeded, actual.MaximumQuantityNeeded, "expected MaximumQuantityNeeded for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.MaximumQuantityNeeded, actual.MaximumQuantityNeeded)
	assert.Equal(t, expected.MinimumQuantityNeeded, actual.MinimumQuantityNeeded, "expected MinimumQuantityNeeded for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.MinimumQuantityNeeded, actual.MinimumQuantityNeeded)

	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestMealPlanGroceryListItems_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanForTest(ctx, t, nil, testClients.admin, testClients.user)

			exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
			exampleMealPlanGroceryListItemInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemCreationRequestInput(exampleMealPlanGroceryListItem)

			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.admin.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			exampleMealPlanGroceryListItem.BelongsToMealPlan = createdMealPlan.ID
			exampleMealPlanGroceryListItem.Ingredient = *createdValidIngredient
			exampleMealPlanGroceryListItem.MeasurementUnit = *createdValidMeasurementUnit

			exampleMealPlanGroceryListItemInput.BelongsToMealPlan = createdMealPlan.ID
			exampleMealPlanGroceryListItemInput.ValidIngredientID = createdValidIngredient.ID
			exampleMealPlanGroceryListItemInput.ValidMeasurementUnitID = createdValidMeasurementUnit.ID

			createdMealPlanGroceryListItem, err := testClients.admin.CreateMealPlanGroceryListItem(ctx, createdMealPlan.ID, exampleMealPlanGroceryListItemInput)
			require.NoError(t, err)
			checkMealPlanGroceryListItemEquality(t, exampleMealPlanGroceryListItem, createdMealPlanGroceryListItem)

			actual, err := testClients.admin.GetMealPlanGroceryListItem(ctx, createdMealPlan.ID, createdMealPlanGroceryListItem.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan task equality
			checkMealPlanGroceryListItemEquality(t, exampleMealPlanGroceryListItem, actual)

			actualList, err := testClients.admin.GetMealPlanGroceryListItemsForMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, actual, err)

			assert.NoError(t, testClients.user.ArchiveMealPlan(ctx, createdMealPlan.ID))

			assert.Len(t, actualList, 1)
			checkMealPlanGroceryListItemEquality(t, actualList[0], actual)
		}
	})
}
