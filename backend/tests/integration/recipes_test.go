package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeEquality(t *testing.T, expected, actual *types.Recipe) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for recipe %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.Equal(t, expected.Source, actual.Source, "expected Source for recipe %s to be %v, but it was %v", expected.ID, expected.Source, actual.Source)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for recipe %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.InspiredByRecipeID, actual.InspiredByRecipeID, "expected InspiredByRecipeID for recipe %s to be %v, but it was %v", expected.ID, expected.InspiredByRecipeID, actual.InspiredByRecipeID)
	assert.Equal(t, expected.MinimumEstimatedPortions, actual.MinimumEstimatedPortions, "expected MinimumEstimatedPortions for recipe %s to be %v, but it was %v", expected.ID, expected.MinimumEstimatedPortions, actual.MinimumEstimatedPortions)
	assert.Equal(t, expected.MaximumEstimatedPortions, actual.MaximumEstimatedPortions, "expected MaximumEstimatedPortions for recipe %s to be %v, but it was %v", expected.ID, expected.MaximumEstimatedPortions, actual.MaximumEstimatedPortions)
	assert.Equal(t, expected.YieldsComponentType, actual.YieldsComponentType, "expected YieldsComponentType for recipe %s to be %v, but it was %v", expected.ID, expected.YieldsComponentType, actual.YieldsComponentType)
	assert.Equal(t, expected.PortionName, actual.PortionName, "expected PortionName for recipe %s to be %v, but it was %v", expected.ID, expected.PortionName, actual.PortionName)
	assert.Equal(t, expected.PluralPortionName, actual.PluralPortionName, "expected PluralPortionName for recipe %s to be %v, but it was %v", expected.ID, expected.PluralPortionName, actual.PluralPortionName)
	assert.Equal(t, expected.SealOfApproval, actual.SealOfApproval, "expected SealOfApproval for recipe %s to be %v, but it was %v", expected.ID, expected.SealOfApproval, actual.SealOfApproval)
	assert.Equal(t, expected.EligibleForMeals, actual.EligibleForMeals, "expected EligibleForMeals for recipe %s to be %v, but it was %v", expected.ID, expected.EligibleForMeals, actual.EligibleForMeals)
	assert.NotZero(t, actual.CreatedAt)
}

func createRecipeForTest(ctx context.Context, t *testing.T, adminClient, client *apiclient.Client, recipe *types.Recipe, inputFilter ...func(input *types.RecipeCreationRequestInput)) ([]*types.ValidIngredient, *types.ValidPreparation, *types.Recipe) {
	t.Helper()

	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
	createdValidPreparation, err := adminClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
	require.NoError(t, err)

	createdValidMeasurementUnit := createValidMeasurementUnitForTest(t, ctx, adminClient)
	createdValidInstrument := createValidInstrumentForTest(t, ctx, adminClient)
	createdValidIngredientState := createValidIngredientStateForTest(t, ctx, adminClient)
	createdValidVessel := createValidVesselForTest(t, ctx, nil, adminClient)

	exampleRecipe := fakes.BuildFakeRecipe()
	if recipe != nil {
		exampleRecipe = recipe
	}

	createdValidIngredients := []*types.ValidIngredient{}
	for i, recipeStep := range exampleRecipe.Steps {
		for j := range recipeStep.Ingredients {
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, createdValidIngredientErr := adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, createdValidIngredientErr)

			createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

			exampleRecipe.Steps[i].Ingredients[j].Ingredient = createdValidIngredient
			exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit = *createdValidMeasurementUnit
		}

		for j := range recipeStep.Products {
			exampleRecipe.Steps[i].Products[j].MeasurementUnit = createdValidMeasurementUnit
		}

		for j := range recipeStep.Instruments {
			recipeStep.Instruments[j].Instrument = createdValidInstrument
		}

		for j := range recipeStep.Vessels {
			recipeStep.Vessels[j].Vessel = createdValidVessel
		}

		for j := range recipeStep.CompletionConditions {
			recipeStep.CompletionConditions[j].IngredientState = *createdValidIngredientState
			for k := range recipeStep.CompletionConditions[j].Ingredients {
				recipeStep.CompletionConditions[j].Ingredients[k].RecipeStepIngredient = createdValidIngredients[0].ID
			}
		}
	}

	exampleRecipeInput := converters.ConvertRecipeToRecipeCreationRequestInput(exampleRecipe)
	exampleRecipeInput.AlsoCreateMeal = true
	for i := range exampleRecipeInput.Steps {
		exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
	}

	examplePrepTask := fakes.BuildFakeRecipePrepTask()
	examplePrepTask.TaskSteps = []*types.RecipePrepTaskStep{
		{
			BelongsToRecipeStep: exampleRecipe.Steps[0].ID,
			SatisfiesRecipeStep: false,
		},
	}
	exampleRecipeInput.PrepTasks = []*types.RecipePrepTaskWithinRecipeCreationRequestInput{
		converters.ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(exampleRecipe, examplePrepTask),
	}

	for _, filter := range inputFilter {
		filter(exampleRecipeInput)
	}

	createdRecipe, err := adminClient.CreateRecipe(ctx, exampleRecipeInput)
	require.NoError(t, err)
	checkRecipeEquality(t, exampleRecipe, createdRecipe)

	createdRecipe, err = client.GetRecipe(ctx, createdRecipe.ID)
	requireNotNilAndNoProblems(t, createdRecipe, err)
	checkRecipeEquality(t, exampleRecipe, createdRecipe)

	require.NotEmpty(t, createdRecipe.Steps, "created recipe must have steps")

	return createdValidIngredients, createdValidPreparation, createdRecipe
}

func (s *TestSuite) TestRecipes_Realistic() {
	s.runForEachClient("sopa de frijol", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			soakBase := fakes.BuildFakeValidPreparation()
			soakInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(soakBase)
			soak, err := testClients.admin.CreateValidPreparation(ctx, soakInput)
			require.NoError(t, err)

			mixBase := fakes.BuildFakeValidPreparation()
			mixInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(mixBase)
			mix, err := testClients.admin.CreateValidPreparation(ctx, mixInput)
			require.NoError(t, err)

			exampleGrams := fakes.BuildFakeValidMeasurementUnit()
			exampleGramsInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleGrams)
			grams, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleGramsInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			grams, err = testClients.admin.GetValidMeasurementUnit(ctx, grams.ID)
			requireNotNilAndNoProblems(t, grams, err)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			exampleCups := fakes.BuildFakeValidMeasurementUnit()
			exampleCupsInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleCups)
			cups, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleCupsInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, exampleCups, cups)

			cups, err = testClients.admin.GetValidMeasurementUnit(ctx, cups.ID)
			requireNotNilAndNoProblems(t, cups, err)
			checkValidMeasurementUnitEquality(t, exampleCups, cups)

			pintoBeanBase := fakes.BuildFakeValidIngredient()
			pintoBeanInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(pintoBeanBase)
			pintoBeans, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, pintoBeanInput)
			require.NoError(t, createdValidIngredientErr)

			waterBase := fakes.BuildFakeValidIngredient()
			waterInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(waterBase)
			water, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, waterInput)
			require.NoError(t, createdValidIngredientErr)

			garlicPaste := fakes.BuildFakeValidIngredient()
			garlicPasteInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(garlicPaste)
			garlicPaste, garlicPasteErr := testClients.admin.CreateValidIngredient(ctx, garlicPasteInput)
			require.NoError(t, garlicPasteErr)

			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			expected := &types.Recipe{
				Name:                     "sopa de frijol",
				Slug:                     "sopa-de-frijol-whatever-who-cares",
				YieldsComponentType:      types.MealComponentTypesMain,
				PortionName:              t.Name(),
				PluralPortionName:        t.Name(),
				MinimumEstimatedPortions: 1,
				Steps: []*types.RecipeStep{
					{
						Products: []*types.RecipeStepProduct{
							{
								Name:            "soaked pinto beans",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: grams,
								QuantityNotes:   "",
								MinimumQuantity: pointer.To(float32(1000)),
							},
						},
						Notes:       "first step",
						Preparation: *soak,
						Instruments: []*types.RecipeStepInstrument{
							{
								Name:       "whatever",
								Instrument: createdValidInstrument,
							},
						},
						Ingredients: []*types.RecipeStepIngredient{
							{
								Ingredient:      pintoBeans,
								Name:            "pinto beans",
								MeasurementUnit: *grams,
								MinimumQuantity: 500,
							},
							{
								Ingredient:      water,
								Name:            "water",
								MeasurementUnit: *cups,
								MinimumQuantity: 5,
							},
						},
						Index: 0,
					},
					{
						Products: []*types.RecipeStepProduct{
							{
								Name:            "final output",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: grams,
								QuantityNotes:   "",
								MinimumQuantity: pointer.To(float32(1010)),
							},
						},
						Notes:       "second step",
						Preparation: *mix,
						Instruments: []*types.RecipeStepInstrument{
							{
								Name:       "whatever",
								Instrument: createdValidInstrument,
							},
						},
						Ingredients: []*types.RecipeStepIngredient{
							{
								Name:            "soaked pinto beans",
								MeasurementUnit: *grams,
								MinimumQuantity: 1000,
							},
							{
								Ingredient:      garlicPaste,
								Name:            "garlic paste",
								MeasurementUnit: *grams,
								MinimumQuantity: 10,
							},
						},
						Index: 1,
					},
				},
			}

			expectedInput := &types.RecipeCreationRequestInput{
				Name:                     expected.Name,
				Description:              expected.Description,
				Slug:                     expected.Slug,
				YieldsComponentType:      expected.YieldsComponentType,
				PortionName:              expected.PortionName,
				PluralPortionName:        expected.PluralPortionName,
				MinimumEstimatedPortions: expected.MinimumEstimatedPortions,
				Steps: []*types.RecipeStepCreationRequestInput{
					{
						MinimumTemperatureInCelsius: expected.Steps[0].MinimumTemperatureInCelsius,
						Products: []*types.RecipeStepProductCreationRequestInput{
							{
								Name:              expected.Steps[0].Products[0].Name,
								Type:              expected.Steps[0].Products[0].Type,
								MeasurementUnitID: &expected.Steps[0].Products[0].MeasurementUnit.ID,
								QuantityNotes:     expected.Steps[0].Products[0].QuantityNotes,
								MinimumQuantity:   expected.Steps[0].Products[0].MinimumQuantity,
							},
						},
						Notes:         expected.Steps[0].Notes,
						PreparationID: expected.Steps[0].Preparation.ID,
						Instruments: []*types.RecipeStepInstrumentCreationRequestInput{
							{
								Name:         "whatever",
								InstrumentID: pointer.To(createdValidInstrument.ID),
							},
						},
						Ingredients: []*types.RecipeStepIngredientCreationRequestInput{
							{
								IngredientID:      &expected.Steps[0].Ingredients[0].Ingredient.ID,
								Name:              expected.Steps[0].Ingredients[0].Name,
								MeasurementUnitID: expected.Steps[0].Ingredients[0].MeasurementUnit.ID,
								MinimumQuantity:   expected.Steps[0].Ingredients[0].MinimumQuantity,
							},
							{
								IngredientID:      &expected.Steps[0].Ingredients[1].Ingredient.ID,
								Name:              expected.Steps[0].Ingredients[1].Name,
								MeasurementUnitID: expected.Steps[0].Ingredients[1].MeasurementUnit.ID,
								MinimumQuantity:   expected.Steps[0].Ingredients[1].MinimumQuantity,
							},
						},
						Index: expected.Steps[0].Index,
					},
					{
						MinimumTemperatureInCelsius: expected.Steps[1].MinimumTemperatureInCelsius,
						Products: []*types.RecipeStepProductCreationRequestInput{
							{
								Name:              expected.Steps[1].Products[0].Name,
								Type:              expected.Steps[1].Products[0].Type,
								MeasurementUnitID: &expected.Steps[1].Products[0].MeasurementUnit.ID,
								QuantityNotes:     expected.Steps[1].Products[0].QuantityNotes,
								MinimumQuantity:   expected.Steps[1].Products[0].MinimumQuantity,
							},
						},
						Notes:         expected.Steps[1].Notes,
						PreparationID: expected.Steps[1].Preparation.ID,
						Instruments: []*types.RecipeStepInstrumentCreationRequestInput{
							{
								Name:         "whatever",
								InstrumentID: pointer.To(createdValidInstrument.ID),
							},
						},
						Ingredients: []*types.RecipeStepIngredientCreationRequestInput{
							{
								Name:                            expected.Steps[1].Ingredients[0].Name,
								MeasurementUnitID:               expected.Steps[1].Ingredients[0].MeasurementUnit.ID,
								MinimumQuantity:                 expected.Steps[1].Ingredients[0].MinimumQuantity,
								ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
								ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
							},
							{
								IngredientID:      &expected.Steps[1].Ingredients[1].Ingredient.ID,
								Name:              expected.Steps[1].Ingredients[1].Name,
								MeasurementUnitID: expected.Steps[1].Ingredients[1].MeasurementUnit.ID,
								MinimumQuantity:   expected.Steps[1].Ingredients[1].MinimumQuantity,
							},
						},
						Index: expected.Steps[1].Index,
					},
				},
			}

			created, err := testClients.admin.CreateRecipe(ctx, expectedInput)
			require.NoError(t, err)
			checkRecipeEquality(t, expected, created)

			created, err = testClients.user.GetRecipe(ctx, created.ID)
			requireNotNilAndNoProblems(t, created, err)
			checkRecipeEquality(t, expected, created)

			recipeStepProductIndex := -1
			for i, ingredient := range created.Steps[1].Ingredients {
				if ingredient.RecipeStepProductID != nil {
					recipeStepProductIndex = i
				}
			}

			require.NotEqual(t, -1, recipeStepProductIndex)
			require.NotNil(t, created.Steps[1].Ingredients[recipeStepProductIndex].RecipeStepProductID)
			assert.Equal(t, created.Steps[0].Products[0].ID, *created.Steps[1].Ingredients[recipeStepProductIndex].RecipeStepProductID)
		}
	})
}

func (s *TestSuite) TestRecipes_Updating() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			newRecipe := fakes.BuildFakeRecipe()
			createdRecipe.Update(converters.ConvertRecipeToRecipeUpdateRequestInput(newRecipe))
			assert.NoError(t, testClients.admin.UpdateRecipe(ctx, createdRecipe))

			actual, err := testClients.user.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, newRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_ContentUploading() {
	s.runForEachClient("should be able to upload content for a recipe", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			newRecipe := fakes.BuildFakeRecipe()
			createdRecipe.Update(converters.ConvertRecipeToRecipeUpdateRequestInput(newRecipe))
			assert.NoError(t, testClients.admin.UpdateRecipe(ctx, createdRecipe))

			actual, err := testClients.user.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, newRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			_, img1Bytes := testutils.BuildArbitraryImagePNGBytes(200)
			_, img2Bytes := testutils.BuildArbitraryImagePNGBytes(250)
			_, img3Bytes := testutils.BuildArbitraryImagePNGBytes(300)

			files := map[string][]byte{
				"image_1.png": img1Bytes,
				"image_2.png": img2Bytes,
				"image_3.png": img3Bytes,
			}

			require.NoError(t, testClients.user.UploadRecipeMedia(ctx, files, actual.ID))

			assert.NoError(t, testClients.admin.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_AlsoCreateMeal() {
	s.runForEachClient("should be able to create a meal and a recipe", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil, func(input *types.RecipeCreationRequestInput) {
				input.AlsoCreateMeal = true
			})

			mealResults, err := testClients.user.SearchForMeals(ctx, createdRecipe.Name, nil)
			requireNotNilAndNoProblems(t, mealResults, err)

			foundMealID := ""
			for _, m := range mealResults.Data {
				meal, mealFetchErr := testClients.user.GetMeal(ctx, m.ID)
				requireNotNilAndNoProblems(t, meal, mealFetchErr)

				for _, component := range meal.Components {
					if component.Recipe.ID == createdRecipe.ID {
						foundMealID = meal.ID
					}
				}
			}

			require.NotEmpty(t, foundMealID)

			assert.NoError(t, testClients.admin.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.user.GetRecipes(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdRecipe := range expected {
				assert.NoError(t, testClients.admin.ArchiveRecipe(ctx, createdRecipe.ID))
			}
		}
	})
}

func (s *TestSuite) TestRecipes_Searching() {
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

			exampleRecipe := fakes.BuildFakeRecipe()

			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				exampleRecipe.Name = fmt.Sprintf("example%d", i)
				_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, exampleRecipe)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.user.SearchForRecipes(ctx, "example", nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdRecipe := range expected {
				assert.NoError(t, testClients.admin.ArchiveRecipe(ctx, createdRecipe.ID))
			}
		}
	})
}

func (s *TestSuite) TestRecipes_GetMealPlanTasksForRecipe() {
	s.runForEachClient("meal plan tasks with frozen chicken breast", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			diceBase := fakes.BuildFakeValidPreparation()
			diceInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(diceBase)
			dice, err := testClients.admin.CreateValidPreparation(ctx, diceInput)
			require.NoError(t, err)

			exampleGrams := fakes.BuildFakeValidMeasurementUnit()
			exampleGramsInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleGrams)
			grams, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleGramsInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			grams, err = testClients.admin.GetValidMeasurementUnit(ctx, grams.ID)
			requireNotNilAndNoProblems(t, grams, err)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			chickenBreastBase := fakes.BuildFakeValidIngredient()
			chickenBreastInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(chickenBreastBase)
			chickenBreastInput.MinimumIdealStorageTemperatureInCelsius = pointer.To(float32(2.5))
			chickenBreast, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, chickenBreastInput)
			require.NoError(t, createdValidIngredientErr)

			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			sauteeBase := fakes.BuildFakeValidPreparation()
			sauteeInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(sauteeBase)
			sautee, err := testClients.admin.CreateValidPreparation(ctx, sauteeInput)
			require.NoError(t, err)

			expected := &types.Recipe{
				Name:                     "sopa de frijol",
				Slug:                     "whatever-who-cares-sopa-de-frijol",
				YieldsComponentType:      types.MealComponentTypesMain,
				PortionName:              t.Name(),
				PluralPortionName:        t.Name(),
				MinimumEstimatedPortions: 1,
				Steps: []*types.RecipeStep{
					{
						MinimumTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProduct{
							{
								Name:            "diced chicken breast",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: grams,
								QuantityNotes:   "",
								MinimumQuantity: pointer.To(float32(1000)),
							},
						},
						Notes:       "first step",
						Preparation: *dice,
						Instruments: []*types.RecipeStepInstrument{
							{
								Name:       "whatever",
								Instrument: createdValidInstrument,
							},
						},
						Ingredients: []*types.RecipeStepIngredient{
							{
								RecipeStepProductID: nil,
								Ingredient:          chickenBreast,
								Name:                "pinto beans",
								MeasurementUnit:     *grams,
								MinimumQuantity:     500,
							},
						},
						Index: 0,
					},
					{
						MinimumTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProduct{
							{
								Name:            "final output",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: grams,
								QuantityNotes:   "",
								MinimumQuantity: pointer.To(float32(1010)),
							},
						},
						Notes:       "second step",
						Preparation: *sautee,
						Instruments: []*types.RecipeStepInstrument{
							{
								Name:       "whatever",
								Instrument: createdValidInstrument,
							},
						},
						Ingredients: []*types.RecipeStepIngredient{
							{
								Name:            "diced chicken breast",
								MeasurementUnit: *grams,
								MinimumQuantity: 1000,
							},
						},
						Index: 1,
					},
				},
			}

			expectedInput := &types.RecipeCreationRequestInput{
				Name:                     expected.Name,
				Slug:                     expected.Slug,
				YieldsComponentType:      expected.YieldsComponentType,
				PortionName:              expected.PortionName,
				PluralPortionName:        expected.PluralPortionName,
				MinimumEstimatedPortions: expected.MinimumEstimatedPortions,
				Steps: []*types.RecipeStepCreationRequestInput{
					{
						MinimumTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProductCreationRequestInput{
							{
								Name:              "diced chicken breast",
								Type:              types.RecipeStepProductIngredientType,
								MeasurementUnitID: &grams.ID,
								QuantityNotes:     "",
								MinimumQuantity:   pointer.To(float32(1000)),
							},
						},
						Notes:         "first step",
						PreparationID: dice.ID,
						Instruments: []*types.RecipeStepInstrumentCreationRequestInput{
							{
								Name:         "whatever",
								InstrumentID: pointer.To(createdValidInstrument.ID),
							},
						},
						Ingredients: []*types.RecipeStepIngredientCreationRequestInput{
							{
								IngredientID:      &chickenBreast.ID,
								Name:              "pinto beans",
								MeasurementUnitID: grams.ID,
								MinimumQuantity:   500,
							},
						},
						Index: 0,
					},
					{
						MinimumTemperatureInCelsius: nil,
						Products: []*types.RecipeStepProductCreationRequestInput{
							{
								Name:              "final output",
								Type:              types.RecipeStepProductIngredientType,
								MeasurementUnitID: &grams.ID,
								QuantityNotes:     "",
								MinimumQuantity:   pointer.To(float32(1010)),
							},
						},
						Notes:         "second step",
						PreparationID: sautee.ID,
						Instruments: []*types.RecipeStepInstrumentCreationRequestInput{
							{
								Name:         "whatever",
								InstrumentID: pointer.To(createdValidInstrument.ID),
							},
						},
						Ingredients: []*types.RecipeStepIngredientCreationRequestInput{
							{
								Name:                            "diced chicken breast",
								MeasurementUnitID:               grams.ID,
								MinimumQuantity:                 1000,
								ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
								ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
							},
						},
						Index: 1,
					},
				},
			}

			created, err := testClients.admin.CreateRecipe(ctx, expectedInput)
			require.NoError(t, err)
			checkRecipeEquality(t, expected, created)

			steps, err := testClients.user.GetMealPlanTasksForRecipe(ctx, created.ID)
			requireNotNilAndNoProblems(t, created, err)

			require.NotEmpty(t, steps)
		}
	})
}

func (s *TestSuite) TestRecipes_Cloning() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			actual, err := testClients.user.CloneRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			require.Equal(t, createdRecipe.Name, actual.Name)
			require.Equal(t, len(createdRecipe.Steps), len(actual.Steps))

			assert.NoError(t, testClients.admin.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_DAGGeneration() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			actual, err := testClients.user.GetRecipeDAG(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			assert.NoError(t, testClients.admin.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
