package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildRecipeForTestCreation(t *testing.T, ctx context.Context, userID string, dbc *Querier) *types.Recipe {
	t.Helper()

	if userID == "" {
		user := createUserForTest(t, ctx, nil, dbc)
		userID = user.ID
	}

	recipeStepID := identifiers.New()
	exampleRecipe := fakes.BuildFakeRecipe()

	exampleRecipeMedia := fakes.BuildFakeRecipeMedia()
	exampleRecipeMedia.BelongsToRecipe = &exampleRecipe.ID
	exampleRecipe.Media = []*types.RecipeMedia{
		exampleRecipeMedia,
	}

	exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
	exampleRecipePrepTask.BelongsToRecipe = exampleRecipe.ID
	exampleRecipePrepTask.TaskSteps = []*types.RecipePrepTaskStep{
		{
			ID:                      identifiers.New(),
			BelongsToRecipeStep:     recipeStepID,
			BelongsToRecipePrepTask: exampleRecipePrepTask.ID,
			SatisfiesRecipeStep:     true,
		},
	}
	exampleRecipe.PrepTasks = []*types.RecipePrepTask{
		exampleRecipePrepTask,
	}

	exampleRecipe.Steps = []*types.RecipeStep{
		buildRecipeStepForTestCreation(t, ctx, exampleRecipe.ID, dbc),
	}
	exampleRecipePrepTask.TaskSteps[0].BelongsToRecipeStep = exampleRecipe.Steps[0].ID
	exampleRecipe.CreatedByUser = userID

	exampleRecipe.Media = []*types.RecipeMedia{}
	for i := range exampleRecipe.Steps {
		exampleRecipe.Steps[i].Media = []*types.RecipeMedia{}
		exampleRecipe.Steps[i].CompletionConditions = nil // []*types.RecipeStepCompletionCondition{}
	}

	return exampleRecipe
}

func createRecipeForTest(t *testing.T, ctx context.Context, exampleRecipe *types.Recipe, dbc *Querier, alsoCreateMeal bool) *types.Recipe {
	t.Helper()

	// create
	if exampleRecipe == nil {
		exampleRecipe = buildRecipeForTestCreation(t, ctx, "", dbc)
	}
	dbInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)
	dbInput.AlsoCreateMeal = alsoCreateMeal

	created, err := dbc.CreateRecipe(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)

	recipe, err := dbc.GetRecipe(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, recipe)
	exampleRecipe.CreatedAt = recipe.CreatedAt

	for i := range recipe.PrepTasks {
		exampleRecipe.PrepTasks[i].CreatedAt = recipe.PrepTasks[i].CreatedAt
	}

	for i := range recipe.Steps {
		assert.Equal(t, exampleRecipe.Steps[i].Preparation.ID, recipe.Steps[i].Preparation.ID)
		exampleRecipe.Steps[i].Preparation = recipe.Steps[i].Preparation
		exampleRecipe.Steps[i].CreatedAt = recipe.Steps[i].CreatedAt

		for j := range recipe.Steps[i].Ingredients {
			exampleRecipe.Steps[i].Ingredients[j].CreatedAt = recipe.Steps[i].Ingredients[j].CreatedAt

			assert.Equal(t, exampleRecipe.Steps[i].Ingredients[j].Ingredient.ID, recipe.Steps[i].Ingredients[j].Ingredient.ID)
			exampleRecipe.Steps[i].Ingredients[j].Ingredient = recipe.Steps[i].Ingredients[j].Ingredient

			assert.Equal(t, exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit.ID, recipe.Steps[i].Ingredients[j].MeasurementUnit.ID)
			exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit = recipe.Steps[i].Ingredients[j].MeasurementUnit
		}
		for j := range recipe.Steps[i].Instruments {
			exampleRecipe.Steps[i].Instruments[j].CreatedAt = recipe.Steps[i].Instruments[j].CreatedAt

			assert.Equal(t, exampleRecipe.Steps[i].Instruments[j].Instrument.ID, recipe.Steps[i].Instruments[j].Instrument.ID)
			exampleRecipe.Steps[i].Instruments[j].Instrument = recipe.Steps[i].Instruments[j].Instrument
		}
		for j := range recipe.Steps[i].Vessels {
			exampleRecipe.Steps[i].Vessels[j].CreatedAt = recipe.Steps[i].Vessels[j].CreatedAt

			assert.Equal(t, exampleRecipe.Steps[i].Vessels[j].Vessel.ID, recipe.Steps[i].Vessels[j].Vessel.ID)
			exampleRecipe.Steps[i].Vessels[j].Vessel = recipe.Steps[i].Vessels[j].Vessel

			if recipe.Steps[i].Vessels[j].Vessel.CapacityUnit != nil {
				assert.Equal(t, exampleRecipe.Steps[i].Vessels[j].Vessel.CapacityUnit.ID, recipe.Steps[i].Vessels[j].Vessel.CapacityUnit.ID)
				exampleRecipe.Steps[i].Vessels[j].Vessel.CapacityUnit = recipe.Steps[i].Vessels[j].Vessel.CapacityUnit
			}
		}
		for j := range recipe.Steps[i].CompletionConditions {
			exampleRecipe.Steps[i].CompletionConditions[j].CreatedAt = recipe.Steps[i].CompletionConditions[j].CreatedAt

			assert.Equal(t, exampleRecipe.Steps[i].CompletionConditions[j].IngredientState.ID, recipe.Steps[i].CompletionConditions[j].IngredientState.ID)
			exampleRecipe.Steps[i].CompletionConditions[j].IngredientState = recipe.Steps[i].CompletionConditions[j].IngredientState

			for k := range recipe.Steps[i].CompletionConditions[j].Ingredients {
				assert.Equal(t, exampleRecipe.Steps[i].CompletionConditions[j].Ingredients[k].ID, recipe.Steps[i].CompletionConditions[j].Ingredients[k].ID)
				exampleRecipe.Steps[i].CompletionConditions[j].Ingredients[k] = recipe.Steps[i].CompletionConditions[j].Ingredients[k]
			}
		}
		for j := range recipe.Steps[i].Products {
			exampleRecipe.Steps[i].Products[j].CreatedAt = recipe.Steps[i].Products[j].CreatedAt
			assert.Equal(t, exampleRecipe.Steps[i].Products[j].MeasurementUnit.ID, recipe.Steps[i].Products[j].MeasurementUnit.ID)
			exampleRecipe.Steps[i].Products[j].MeasurementUnit = recipe.Steps[i].Products[j].MeasurementUnit
		}
		for j := range recipe.Steps[i].Media {
			assert.Equal(t, exampleRecipe.Steps[i].Media[j].ID, recipe.Steps[i].Media[j].ID)
			exampleRecipe.Steps[i].Media[j] = recipe.Steps[i].Media[j]
		}

		require.Equal(t, len(exampleRecipe.Steps[i].Products), len(recipe.Steps[i].Products))
		require.Equal(t, len(exampleRecipe.Steps[i].Instruments), len(recipe.Steps[i].Instruments))
		require.Equal(t, len(exampleRecipe.Steps[i].Vessels), len(recipe.Steps[i].Vessels))
		require.Equal(t, len(exampleRecipe.Steps[i].Ingredients), len(recipe.Steps[i].Ingredients))
		require.Equal(t, len(exampleRecipe.Steps[i].Media), len(recipe.Steps[i].Media))
		require.Equal(t, len(exampleRecipe.Steps[i].CompletionConditions), len(recipe.Steps[i].CompletionConditions))

		expectedStep := exampleRecipe.Steps[i]
		actualStep := recipe.Steps[i]
		require.Equal(t, expectedStep, actualStep)
	}

	require.Equal(t, exampleRecipe, recipe)

	return recipe
}

func TestQuerier_Integration_Recipes(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)

	createdRecipes := []*types.Recipe{}

	// create
	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	createdRecipes = append(createdRecipes, createRecipeForTest(t, ctx, exampleRecipe, dbc, true))

	// update
	updatedRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	updatedRecipe.ID = createdRecipes[0].ID
	assert.NoError(t, dbc.UpdateRecipe(ctx, updatedRecipe))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		exampleRecipe = buildRecipeForTestCreation(t, ctx, user.ID, dbc)
		exampleRecipe.Name = fmt.Sprintf("%s %d", exampleRecipe.Name, i)
		createdRecipes = append(createdRecipes, createRecipeForTest(t, ctx, exampleRecipe, dbc, true))
	}

	// fetch as list
	recipes, err := dbc.GetRecipes(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipes.Data)
	assert.Equal(t, len(createdRecipes), len(recipes.Data))

	needingIndexing, err := dbc.GetRecipeIDsThatNeedSearchIndexing(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, needingIndexing)

	// search
	searchResults, err := dbc.SearchForRecipes(ctx, createdRecipes[0].Name, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, searchResults)

	byIDs, err := dbc.GetRecipesWithIDs(ctx, []string{createdRecipes[0].ID})
	assert.NoError(t, err)
	assert.NotEmpty(t, byIDs)
	assert.Equal(t, createdRecipes[0].ID, byIDs[0].ID)

	// delete
	for _, recipe := range createdRecipes {
		assert.NoError(t, dbc.MarkRecipeAsIndexed(ctx, recipe.ID))
		assert.NoError(t, dbc.ArchiveRecipe(ctx, recipe.ID, user.ID))

		var exists bool
		exists, err = dbc.RecipeExists(ctx, recipe.ID)
		assert.NoError(t, err)
		assert.False(t, exists)
	}
}

func TestQuerier_RecipeExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_CreateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipe(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipe(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipe(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipe(ctx, exampleRecipe.ID, ""))
	})
}

func Test_findCreatedRecipeStepProductsForIngredients(T *testing.T) {
	T.Parallel()

	T.Run("sopa de frijol", func(t *testing.T) {
		t.Parallel()

		soak := fakes.BuildFakeValidPreparation()
		water := fakes.BuildFakeValidIngredient()
		pintoBeans := fakes.BuildFakeValidIngredient()
		garlicPaste := fakes.BuildFakeValidIngredient()
		productName := "soaked pinto beans"

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        "sopa de frijol",
			Description: "",
			Steps: []*types.RecipeStepDatabaseCreationInput{
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:                fakes.BuildFakeID(),
							Name:              productName,
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "first step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							IngredientID:      &pintoBeans.ID,
							Name:              "pinto beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
						{
							IngredientID:      &water.ID,
							Name:              "water",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
					},
					Index: 0,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							Name:              "final output",
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "second step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							Name:                            productName,
							MeasurementUnitID:               fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:                 1000,
							ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
							ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
						},
						{
							IngredientID:      &garlicPaste.ID,
							Name:              "garlic paste",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   10,
						},
					},
					Index: 1,
				},
			},
		}

		findCreatedRecipeStepProductsForIngredients(exampleRecipeInput)

		require.NotNil(t, exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
	})

	T.Run("slightly more complicated recipe", func(t *testing.T) {
		t.Parallel()

		soak := fakes.BuildFakeValidPreparation()
		water := fakes.BuildFakeValidIngredient()
		pintoBeans := fakes.BuildFakeValidIngredient()
		garlicPaste := fakes.BuildFakeValidIngredient()
		productName := "soaked pinto beans"

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        "sopa de frijol",
			Description: "",
			Steps: []*types.RecipeStepDatabaseCreationInput{
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:                fakes.BuildFakeID(),
							Name:              productName,
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "first step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							IngredientID:      &pintoBeans.ID,
							Name:              "pinto beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
						{
							IngredientID:      &water.ID,
							Name:              "water",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   5,
						},
					},
					Index: 0,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							Name:              "pressure cooked beans",
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "second step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							Name:                            productName,
							MeasurementUnitID:               fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:                 1000,
							ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
							ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
						},
						{
							IngredientID:      &garlicPaste.ID,
							Name:              "garlic paste",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   10,
						},
					},
					Index: 1,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:                fakes.BuildFakeID(),
							Name:              productName,
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "third step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							IngredientID:      &pintoBeans.ID,
							Name:              "pinto beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
						{
							IngredientID:      &water.ID,
							Name:              "water",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   5,
						},
					},
					Index: 2,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							Name:              "final output",
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "fourth step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							Name:                            productName,
							MeasurementUnitID:               fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:                 1000,
							ProductOfRecipeStepIndex:        pointer.To(uint64(2)),
							ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
						},
						{
							Name:              "pressure cooked beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   10,
						},
					},
					Index: 3,
				},
			},
		}

		findCreatedRecipeStepProductsForIngredients(exampleRecipeInput)

		require.NotNil(t, exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
		require.NotNil(t, exampleRecipeInput.Steps[3].Ingredients[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[2].Products[0].ID, *exampleRecipeInput.Steps[3].Ingredients[0].RecipeStepProductID)
	})
}

func Test_findCreatedRecipeStepProductsForInstruments(T *testing.T) {
	T.Parallel()

	T.Run("example", func(t *testing.T) {
		t.Parallel()

		bake := fakes.BuildFakeValidPreparation()
		line := fakes.BuildFakeValidPreparation()
		bakingSheet := fakes.BuildFakeValidInstrument()
		aluminumFoil := fakes.BuildFakeValidIngredient()
		asparagus := fakes.BuildFakeValidIngredient()
		grams := fakes.BuildFakeValidMeasurementUnit()
		sheet := fakes.BuildFakeValidMeasurementUnit()

		productName := "lined baking sheet"

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        "example",
			Description: "",
			Steps: []*types.RecipeStepDatabaseCreationInput{
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:   fakes.BuildFakeID(),
							Name: productName,
							Type: types.RecipeStepProductInstrumentType,
						},
					},
					Instruments: []*types.RecipeStepInstrumentDatabaseCreationInput{
						{
							InstrumentID:        &bakingSheet.ID,
							RecipeStepProductID: nil,
							Name:                "baking sheet",
						},
					},
					Notes:         "first step",
					PreparationID: line.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							RecipeStepProductID: nil,
							IngredientID:        &aluminumFoil.ID,
							Name:                "aluminum foil",
							MeasurementUnitID:   sheet.ID,
							MinimumQuantity:     1,
						},
					},
					Index: 0,
				},
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:   fakes.BuildFakeID(),
							Name: "roasted asparagus",
							Type: types.RecipeStepProductInstrumentType,
						},
					},
					Instruments: []*types.RecipeStepInstrumentDatabaseCreationInput{
						{
							InstrumentID:                    &bakingSheet.ID,
							RecipeStepProductID:             nil,
							Name:                            productName,
							ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
							ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
						},
					},
					Notes:         "second step",
					PreparationID: bake.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							RecipeStepProductID: nil,
							IngredientID:        &asparagus.ID,
							Name:                "asparagus",
							MeasurementUnitID:   grams.ID,
							MinimumQuantity:     1000,
						},
					},
					Index: 1,
				},
			},
		}

		findCreatedRecipeStepProductsForInstruments(exampleRecipeInput)

		require.NotNil(t, exampleRecipeInput.Steps[1].Instruments[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Instruments[0].RecipeStepProductID)
	})
}

func TestQuerier_MarkRecipeAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkRecipeAsIndexed(ctx, ""))
	})
}
