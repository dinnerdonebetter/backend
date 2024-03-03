package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkMealEquality(t *testing.T, expected, actual *types.Meal) {
	t.Helper()

	assert.NotZero(t, actual.ID)

	assert.Equal(t, expected.Name, actual.Name, "expected Name for meal %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for meal %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.MinimumEstimatedPortions, actual.MinimumEstimatedPortions, "expected MinimumEstimatedPortions for meal %s to be %v, but it was %v", expected.ID, expected.MinimumEstimatedPortions, actual.MinimumEstimatedPortions)
	assert.Equal(t, expected.MaximumEstimatedPortions, actual.MaximumEstimatedPortions, "expected MaximumEstimatedPortions for meal %s to be %v, but it was %v", expected.ID, expected.MaximumEstimatedPortions, actual.MaximumEstimatedPortions)
	assert.Equal(t, expected.EligibleForMealPlans, actual.EligibleForMealPlans, "expected EligibleForMealPlans for meal %s to be %v, but it was %v", expected.ID, expected.EligibleForMealPlans, actual.EligibleForMealPlans)

	assert.NotZero(t, actual.CreatedAt)
}

func createMealForTest(ctx context.Context, t *testing.T, adminClient, client *apiclient.Client, mealInput *types.Meal) *types.Meal {
	t.Helper()

	createdRecipes := []*types.Recipe{}
	createdRecipeIDs := []*types.MealComponentCreationRequestInput{}
	for i := 0; i < 3; i++ {
		_, _, recipe := createRecipeForTest(ctx, t, adminClient, client, nil)
		createdRecipes = append(createdRecipes, recipe)
		createdRecipeIDs = append(createdRecipeIDs, &types.MealComponentCreationRequestInput{
			RecipeID:      recipe.ID,
			ComponentType: types.MealComponentTypesMain,
		})
	}

	exampleMeal := mealInput
	if exampleMeal == nil {
		exampleMeal = fakes.BuildFakeMeal()
	}

	exampleMealInput := converters.ConvertMealToMealCreationRequestInput(exampleMeal)
	exampleMealInput.Components = createdRecipeIDs

	createdMeal, err := client.CreateMeal(ctx, exampleMealInput)
	require.NoError(t, err)

	createdMeal, err = client.GetMeal(ctx, createdMeal.ID)
	requireNotNilAndNoProblems(t, createdMeal, err)
	checkMealEquality(t, exampleMeal, createdMeal)

	return createdMeal
}

func (s *TestSuite) TestMeals_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)

			assert.NoError(t, testClients.user.ArchiveMeal(ctx, createdMeal.ID))
		}
	})
}

func (s *TestSuite) TestMeals_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)

			require.NoError(t, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.user.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			var expected []*types.Meal
			for i := 0; i < 5; i++ {
				createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)

				expected = append(expected, createdMeal)
			}

			// assert meal list equality
			actual, err := testClients.user.GetMeals(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdMeal := range expected {
				assert.NoError(t, testClients.user.ArchiveMeal(ctx, createdMeal.ID))
			}
		}
	})
}

func (s *TestSuite) TestMeals_Searching() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.user.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			exampleMeal := fakes.BuildFakeMeal()
			var expected []*types.Meal
			for i := 0; i < 5; i++ {
				exampleMeal.Name = fmt.Sprintf("example%d", i)
				createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, exampleMeal)

				expected = append(expected, createdMeal)
			}

			// assert meal list equality
			actual, err := testClients.user.SearchForMeals(ctx, "example", nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdMeal := range expected {
				assert.NoError(t, testClients.user.ArchiveMeal(ctx, createdMeal.ID))
			}
		}
	})
}
