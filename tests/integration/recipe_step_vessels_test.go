package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/pointers"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func checkRecipeStepVesselEquality(t *testing.T, expected, actual *types.RecipeStepVessel, checkInstrument bool) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	if checkInstrument {
		checkValidInstrumentEquality(t, expected.Instrument, actual.Instrument)
	} else {
		assert.Equal(t, expected.Instrument.ID, actual.Instrument.ID, "expected Instrument.ID for recipe step vessel %s to be %v, but it was %v", expected.ID, expected.Instrument.ID, actual.Instrument.ID)
	}

	assert.Equal(t, expected.RecipeStepProductID, actual.RecipeStepProductID, "expected RecipeStepProductID for recipe step vessel %s to be %v, but it was %v", expected.ID, expected.RecipeStepProductID, actual.RecipeStepProductID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step vessel %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe step vessel %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.BelongsToRecipeStep, actual.BelongsToRecipeStep, "expected BelongsToRecipeStep for recipe step vessel %s to be %v, but it was %v", expected.ID, expected.BelongsToRecipeStep, actual.BelongsToRecipeStep)
	assert.Equal(t, expected.VesselPredicate, actual.VesselPredicate, "expected VesselPredicate for recipe step vessel %s to be %v, but it was %v", expected.ID, expected.VesselPredicate, actual.VesselPredicate)
	assert.Equal(t, expected.MaximumQuantity, actual.MaximumQuantity, "expected MaximumQuantity for recipe step vessel %s to be %v, but it was %v", expected.ID, expected.MaximumQuantity, actual.MaximumQuantity)
	assert.Equal(t, expected.MinimumQuantity, actual.MinimumQuantity, "expected MinimumQuantity for recipe step vessel %s to be %v, but it was %v", expected.ID, expected.MinimumQuantity, actual.MinimumQuantity)
	assert.Equal(t, expected.UnavailableAfterStep, actual.UnavailableAfterStep, "expected UnavailableAfterStep for recipe step vessel %s to be %v, but it was %v", expected.ID, expected.UnavailableAfterStep, actual.UnavailableAfterStep)

	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestRecipeStepVessels_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrument.IsVessel = true
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			t.Log("creating recipe step vessel")
			exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
			exampleRecipeStepVessel.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepVessel.Instrument = &types.ValidInstrument{ID: createdValidInstrument.ID}
			exampleRecipeStepVesselInput := converters.ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(exampleRecipeStepVessel)
			createdRecipeStepVessel, err := testClients.user.CreateRecipeStepVessel(ctx, createdRecipe.ID, createdRecipeStepID, exampleRecipeStepVesselInput)
			require.NoError(t, err)
			t.Logf("recipe step vessel %q created", createdRecipeStepVessel.ID)

			checkRecipeStepVesselEquality(t, exampleRecipeStepVessel, createdRecipeStepVessel, false)

			createdRecipeStepVessel, err = testClients.user.GetRecipeStepVessel(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepVessel.ID)
			requireNotNilAndNoProblems(t, createdRecipeStepVessel, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepVessel.BelongsToRecipeStep)
			exampleRecipeStepVessel.Instrument = createdValidInstrument
			exampleRecipeStepVessel.Instrument.CreatedAt = createdRecipeStepVessel.Instrument.CreatedAt

			checkRecipeStepVesselEquality(t, exampleRecipeStepVessel, createdRecipeStepVessel, false)

			t.Log("creating valid instrument")
			newExampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrument.IsVessel = true
			newExampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(newExampleValidInstrument)
			newValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, newExampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, newExampleValidInstrument, newValidInstrument)

			t.Log("changing recipe step vessel")
			newRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
			newRecipeStepVessel.BelongsToRecipeStep = createdRecipeStepID
			newRecipeStepVessel.Instrument = newValidInstrument
			createdRecipeStepVessel.Update(converters.ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput(newRecipeStepVessel))
			assert.NoError(t, testClients.user.UpdateRecipeStepVessel(ctx, createdRecipe.ID, createdRecipeStepVessel))

			t.Log("fetching changed recipe step vessel")
			actual, err := testClients.user.GetRecipeStepVessel(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepVessel.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step vessel equality
			checkRecipeStepVesselEquality(t, newRecipeStepVessel, actual, false)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up recipe step vessel")
			assert.NoError(t, testClients.user.ArchiveRecipeStepVessel(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepVessel.ID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepVessels_AsRecipeStepProducts() {
	s.runForEachClient("should be able to use a recipe step vessel that was the product of a prior recipe step", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid preparation")
			lineBase := fakes.BuildFakeValidPreparation()
			lineInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(lineBase)
			line, err := testClients.admin.CreateValidPreparation(ctx, lineInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", line.ID)

			t.Log("creating prerequisite valid preparation")
			roastBase := fakes.BuildFakeValidPreparation()
			roastInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(roastBase)
			roast, err := testClients.admin.CreateValidPreparation(ctx, roastInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", roast.ID)

			t.Log("creating valid instrument")
			bakingSheetBase := fakes.BuildFakeValidInstrument()
			bakingSheetBase.IsVessel = true
			bakingSheetBaseInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(bakingSheetBase)
			bakingSheet, err := testClients.admin.CreateValidInstrument(ctx, bakingSheetBaseInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", bakingSheet.ID)
			checkValidInstrumentEquality(t, bakingSheetBase, bakingSheet)

			t.Log("creating valid measurement units")
			sheetsBase := fakes.BuildFakeValidMeasurementUnit()
			sheetsBaseInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(sheetsBase)
			sheets, err := testClients.admin.CreateValidMeasurementUnit(ctx, sheetsBaseInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", sheets.ID)
			checkValidMeasurementUnitEquality(t, sheetsBase, sheets)

			t.Log("creating valid measurement units")
			headsBase := fakes.BuildFakeValidMeasurementUnit()
			headsBaseInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(headsBase)
			head, err := testClients.admin.CreateValidMeasurementUnit(ctx, headsBaseInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", head.ID)
			checkValidMeasurementUnitEquality(t, headsBase, head)

			t.Log("creating valid measurement units")
			exampleUnits := fakes.BuildFakeValidMeasurementUnit()
			exampleUnitsInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleUnits)
			unit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleUnitsInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", unit.ID)
			checkValidMeasurementUnitEquality(t, exampleUnits, unit)

			t.Log("creating prerequisite valid ingredient")
			aluminumFoilBase := fakes.BuildFakeValidIngredient()
			aluminumFoilInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(aluminumFoilBase)
			aluminumFoil, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, aluminumFoilInput)
			require.NoError(t, createdValidIngredientErr)

			t.Log("creating prerequisite valid ingredient")
			garlic := fakes.BuildFakeValidIngredient()
			garlicInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(garlic)
			garlic, garlicErr := testClients.admin.CreateValidIngredient(ctx, garlicInput)
			require.NoError(t, garlicErr)

			linedBakingSheetName := "lined baking sheet"

			t.Log("creating recipe")
			expected := &types.Recipe{
				Name: t.Name(),
				Steps: []*types.RecipeStep{
					{
						Products: []*types.RecipeStepProduct{
							{
								Name:            linedBakingSheetName,
								Type:            types.RecipeStepProductVesselType,
								MeasurementUnit: unit,
								QuantityNotes:   "",
								MinimumQuantity: pointers.Float32(1),
							},
						},
						Notes:       "first step",
						Preparation: *line,
						Ingredients: []*types.RecipeStepIngredient{
							{
								RecipeStepProductID: nil,
								Ingredient:          aluminumFoil,
								Name:                "aluminum foil",
								MeasurementUnit:     *sheets,
								MinimumQuantity:     3,
							},
						},
						Vessels: []*types.RecipeStepVessel{
							{
								Instrument: bakingSheet,
							},
						},
						Index: 0,
					},
					{
						Preparation: *roast,
						Vessels: []*types.RecipeStepVessel{
							{
								Name:       linedBakingSheetName,
								Instrument: nil,
							},
						},
						Products: []*types.RecipeStepProduct{
							{
								Name:            "roasted garlic",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: head,
								QuantityNotes:   "",
								MinimumQuantity: pointers.Float32(1),
							},
						},
						Notes: "second step",
						Ingredients: []*types.RecipeStepIngredient{
							{
								Ingredient:      garlic,
								Name:            "garlic",
								MeasurementUnit: *head,
								MinimumQuantity: 1,
							},
						},
						Index: 1,
					},
				},
			}

			exampleRecipeInput := converters.ConvertRecipeToRecipeCreationRequestInput(expected)
			exampleRecipeInput.Steps[1].Vessels[0].ProductOfRecipeStepIndex = pointers.Pointer(uint64(0))
			exampleRecipeInput.Steps[1].Vessels[0].ProductOfRecipeStepProductIndex = pointers.Pointer(uint64(0))

			inputJSON, _ := json.Marshal(exampleRecipeInput)
			t.Log(string(inputJSON))

			created, err := testClients.user.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", created.ID)
			checkRecipeEquality(t, expected, created)

			created, err = testClients.user.GetRecipe(ctx, created.ID)
			requireNotNilAndNoProblems(t, created, err)
			checkRecipeEquality(t, expected, created)

			createdJSON, _ := json.Marshal(created)
			t.Log(string(createdJSON))

			recipeStepProductIndex := -1
			for i, ingredient := range created.Steps[1].Vessels {
				if ingredient.RecipeStepProductID != nil {
					recipeStepProductIndex = i
				}
			}

			require.NotEqual(t, -1, recipeStepProductIndex)
			require.NotNil(t, created.Steps[1].Vessels[recipeStepProductIndex].RecipeStepProductID)
			assert.Equal(t, created.Steps[0].Products[0].ID, *created.Steps[1].Vessels[recipeStepProductIndex].RecipeStepProductID)
		}
	})
}

func (s *TestSuite) TestRecipeStepVessels_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrument.IsVessel = true
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			t.Log("creating recipe step vessels")
			var expected []*types.RecipeStepVessel
			for i := 0; i < 5; i++ {
				exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
				exampleRecipeStepVessel.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepVessel.Instrument = &types.ValidInstrument{ID: createdValidInstrument.ID}
				exampleRecipeStepVesselInput := converters.ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(exampleRecipeStepVessel)
				createdRecipeStepVessel, createdRecipeStepVesselErr := testClients.user.CreateRecipeStepVessel(ctx, createdRecipe.ID, createdRecipeStepID, exampleRecipeStepVesselInput)
				require.NoError(t, createdRecipeStepVesselErr)
				t.Logf("recipe step vessel %q created", createdRecipeStepVessel.ID)
				checkRecipeStepVesselEquality(t, exampleRecipeStepVessel, createdRecipeStepVessel, false)

				createdRecipeStepVessel, createdRecipeStepVesselErr = testClients.user.GetRecipeStepVessel(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepVessel.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepVessel, createdRecipeStepVesselErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepVessel.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepVessel)
			}

			// assert recipe step vessel list equality
			actual, err := testClients.user.GetRecipeStepVessels(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			t.Log("cleaning up")
			for _, createdRecipeStepVessel := range expected {
				assert.NoError(t, testClients.user.ArchiveRecipeStepVessel(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepVessel.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}