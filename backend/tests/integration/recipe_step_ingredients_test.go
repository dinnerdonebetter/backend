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

func checkRecipeStepIngredientEquality(t *testing.T, expected, actual *types.RecipeStepIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.RecipeStepProductID, actual.RecipeStepProductID, "expected RecipeStepProductID for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.RecipeStepProductID, actual.RecipeStepProductID)
	assert.Equal(t, expected.Ingredient, actual.Ingredient, "expected Ingredient for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.Ingredient, actual.Ingredient)
	assert.Equal(t, expected.MaximumQuantity, actual.MaximumQuantity, "expected MaximumQuantity for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.MaximumQuantity, actual.MaximumQuantity)
	assert.Equal(t, expected.VesselIndex, actual.VesselIndex, "expected VesselIndex for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.VesselIndex, actual.VesselIndex)
	assert.Equal(t, expected.RecipeStepProductRecipeID, actual.RecipeStepProductRecipeID, "expected RecipeStepProductRecipeID for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.RecipeStepProductRecipeID, actual.RecipeStepProductRecipeID)
	assert.Equal(t, expected.ProductPercentageToUse, actual.ProductPercentageToUse, "expected ProductPercentageToUse for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.ProductPercentageToUse, actual.ProductPercentageToUse)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.BelongsToRecipeStep, actual.BelongsToRecipeStep, "expected BelongsToRecipeStep for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.BelongsToRecipeStep, actual.BelongsToRecipeStep)
	assert.Equal(t, expected.IngredientNotes, actual.IngredientNotes, "expected IngredientNotes for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.IngredientNotes, actual.IngredientNotes)
	assert.Equal(t, expected.MeasurementUnit, actual.MeasurementUnit, "expected MeasurementUnit for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.MeasurementUnit, actual.MeasurementUnit)
	assert.Equal(t, expected.MinimumQuantity, actual.MinimumQuantity, "expected MinimumQuantity for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.MinimumQuantity, actual.MinimumQuantity)
	assert.Equal(t, expected.OptionIndex, actual.OptionIndex, "expected OptionIndex for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.OptionIndex, actual.OptionIndex)
	assert.Equal(t, expected.Optional, actual.Optional, "expected Optional for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.Optional, actual.Optional)
	assert.Equal(t, expected.ToTaste, actual.ToTaste, "expected ToTaste for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.ToTaste, actual.ToTaste)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestRecipeStepIngredients_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			firstStep := createdRecipe.Steps[0]
			createdRecipeStepID := firstStep.ID
			createdRecipeStepIngredientID := firstStep.Ingredients[0].ID

			require.NotEmpty(t, createdRecipeStepID, "created recipe step ID must not be empty")
			require.NotEmpty(t, createdRecipeStepIngredientID, "created recipe step ingredient ID must not be empty")

			createdRecipeStepIngredient, err := testClients.user.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
			requireNotNilAndNoProblems(t, createdRecipeStepIngredient, err)

			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.admin.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			newRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			newRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
			newRecipeStepIngredient.ID = createdRecipeStepIngredientID
			newRecipeStepIngredient.Ingredient = createdValidIngredient
			newRecipeStepIngredient.MeasurementUnit = createdRecipeStepIngredient.MeasurementUnit

			createdRecipeStepIngredient.Update(converters.ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput(newRecipeStepIngredient))

			require.NoError(t, testClients.admin.UpdateRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepIngredient))

			actual, err := testClients.user.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, newRecipeStepIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.user.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID))

			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			assert.NoError(t, testClients.admin.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var (
				createdRecipeStepID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			var expected []*types.RecipeStepIngredient
			for i := 0; i < 5; i++ {
				x, _, _ := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

				exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
				exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepIngredient.Ingredient = &types.ValidIngredient{ID: x[0].ID}
				exampleRecipeStepIngredient.MeasurementUnit = types.ValidMeasurementUnit{ID: createdValidMeasurementUnit.ID}

				exampleRecipeStepIngredientInput := converters.ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(exampleRecipeStepIngredient)
				createdRecipeStepIngredient, createdRecipeStepIngredientErr := testClients.admin.CreateRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, exampleRecipeStepIngredientInput)
				require.NoError(t, createdRecipeStepIngredientErr)

				checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, createdRecipeStepIngredient)

				createdRecipeStepIngredient, createdRecipeStepIngredientErr = testClients.user.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredient.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepIngredient, createdRecipeStepIngredientErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepIngredient.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepIngredient)
			}

			// assert recipe step ingredient list equality
			actual, err := testClients.user.GetRecipeStepIngredients(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdRecipeStepIngredient := range expected {
				assert.NoError(t, testClients.user.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredient.ID))
			}

			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			assert.NoError(t, testClients.admin.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
