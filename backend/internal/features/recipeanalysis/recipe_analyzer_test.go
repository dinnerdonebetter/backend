package recipeanalysis

import (
	"context"
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func saveDebugDAGOfRecipe(t *testing.T, g RecipeAnalyzer, recipe *types.Recipe) {
	t.Helper()

	ctx := context.Background()

	img, err := g.GenerateDAGDiagramForRecipe(ctx, recipe)
	require.NoError(t, err)

	f, err := os.CreateTemp("", "")
	require.NoError(t, err)
	require.NoError(t, png.Encode(f, img))
}

func newAnalyzerForTest(t *testing.T) *recipeAnalyzer {
	t.Helper()

	return &recipeAnalyzer{
		tracer: tracing.NewTracerForTest(t.Name()),
		logger: logging.NewNoopLogger(),
	}
}

func TestRecipeGrapher_makeGraphForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		g := newAnalyzerForTest(t)

		ctx := context.Background()
		r := &types.Recipe{
			Steps: []*types.RecipeStep{
				{},
			},
		}

		actual, err := g.MakeGraphForRecipe(ctx, r)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestRecipeGrapher_makeDAGForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		g := newAnalyzerForTest(t)

		ctx := context.Background()
		r := &types.Recipe{
			Steps: []*types.RecipeStep{
				{},
			},
		}

		actual, err := g.makeDAGForRecipe(ctx, r)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestRecipeAnalyzer_GenerateMealPlanTasksForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("creates frozen thawing steps", func(t *testing.T) {
		t.Parallel()

		g := newAnalyzerForTest(t)
		ctx := context.Background()

		exampleMeal := fakes.BuildFakeMeal()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		exampleMealPlanEvent.BelongsToMealPlan = exampleMealPlan.ID
		now := time.Now().Add(0).Truncate(time.Second).UTC()
		inThreeDays := now.Add((time.Hour * 24) * 3).Add(0).Truncate(time.Second).UTC()
		inOneWeek := now.Add((time.Hour * 24) * 7).Add(0).Truncate(time.Second).UTC()
		exampleMealPlanEvent.StartsAt = inThreeDays
		exampleMealPlanEvent.EndsAt = inOneWeek

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
		exampleMealPlanOption.BelongsToMealPlanEvent = exampleMealPlanEvent.ID
		exampleMealPlanOption.Meal = *exampleMeal

		recipeStepID := fakes.BuildFakeID()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipe := &types.Recipe{
			Name: "Recipe 1",
			ID:   exampleRecipeID,
			Steps: []*types.RecipeStep{
				{
					MaximumEstimatedTimeInSeconds: nil,
					MinimumTemperatureInCelsius:   nil,
					MinimumEstimatedTimeInSeconds: nil,
					MaximumTemperatureInCelsius:   nil,
					BelongsToRecipe:               exampleRecipeID,
					ID:                            recipeStepID,
					Preparation:                   types.ValidPreparation{Name: "dice"},
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient: &types.ValidIngredient{
								MaximumIdealStorageTemperatureInCelsius: nil,
								MinimumIdealStorageTemperatureInCelsius: pointer.To(float32(2.5)),
								PluralName:                              "chicken breasts",
								StorageInstructions:                     "keep frozen",
								Name:                                    "chicken breast",
								ID:                                      fakes.BuildFakeID(),
							},
							Name:                "chicken breast",
							ID:                  fakes.BuildFakeID(),
							BelongsToRecipeStep: recipeStepID,
							MeasurementUnit:     types.ValidMeasurementUnit{Name: "gram", PluralName: "grams"},
							MinimumQuantity:     900,
							MaximumQuantity:     pointer.To(float32(900)),
							Optional:            false,
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							MinimumStorageTemperatureInCelsius: nil,
							MaximumStorageTemperatureInCelsius: nil,
							StorageInstructions:                "",
							Name:                               "diced chicken breast",
							Type:                               types.RecipeStepProductIngredientType,
							BelongsToRecipeStep:                recipeStepID,
							ID:                                 fakes.BuildFakeID(),
							QuantityNotes:                      "",
							MeasurementUnit:                    &types.ValidMeasurementUnit{},
							MaximumStorageDurationInSeconds:    nil,
							MaximumQuantity:                    nil,
							MinimumQuantity:                    nil,
							Compostable:                        false,
						},
					},
					Instruments: nil,
				},
			},
		}

		expected := []*types.MealPlanTaskDatabaseCreationInput{
			{
				CreationExplanation: buildThawStepCreationExplanation(1, 0),
				MealPlanOptionID:    exampleMealPlanOption.ID,
			},
		}

		actual, err := g.GenerateMealPlanTasksForRecipe(ctx, exampleMealPlanOption.ID, exampleRecipe)
		assert.NoError(t, err)

		for i := range expected {
			expected[i].ID = actual[i].ID
		}

		assert.Equal(t, expected, actual)
	})
}

func Test_recipeAnalyzer_RenderMermaidDiagramForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("basic", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		g := newAnalyzerForTest(t)

		dice := fakes.BuildFakeValidPreparation()
		dice.Name = "dice"
		sautee := fakes.BuildFakeValidPreparation()
		sautee.Name = "sautee"

		step1ID := fakes.BuildFakeID()
		step2ID := fakes.BuildFakeID()
		step3ID := fakes.BuildFakeID()
		step4ID := fakes.BuildFakeID()
		dicedOnionRecipeStepProductID := fakes.BuildFakeID()
		dicedCarrotRecipeStepProductID := fakes.BuildFakeID()
		dicedCeleryRecipeStepProductID := fakes.BuildFakeID()

		recipe := &types.Recipe{
			Name: "example recipe",
			Steps: []*types.RecipeStep{
				{
					ID:          step1ID,
					Preparation: *dice,
					Ingredients: []*types.RecipeStepIngredient{
						{
							Ingredient: fakes.BuildFakeValidIngredient(),
							Name:       "onion",
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							ID:   dicedOnionRecipeStepProductID,
							Name: "diced onion",
							Type: types.RecipeStepProductIngredientType,
						},
					},
					Notes: "first step",
					Index: 0,
				},
				{
					ID:          step2ID,
					Preparation: *dice,
					Ingredients: []*types.RecipeStepIngredient{
						{
							Ingredient: fakes.BuildFakeValidIngredient(),
							Name:       "carrot",
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							ID:   dicedCarrotRecipeStepProductID,
							Name: "diced carrot",
							Type: types.RecipeStepProductIngredientType,
						},
					},
					Notes: "second step",
					Index: 1,
				},
				{
					ID:          step3ID,
					Preparation: *dice,
					Ingredients: []*types.RecipeStepIngredient{
						{
							Ingredient: fakes.BuildFakeValidIngredient(),
							Name:       "celery",
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							ID:   dicedCeleryRecipeStepProductID,
							Name: "diced celery",
							Type: types.RecipeStepProductIngredientType,
						},
					},
					Notes: "third step",
					Index: 2,
				},
				{
					ID:          step4ID,
					Preparation: *sautee,
					Ingredients: []*types.RecipeStepIngredient{
						{
							Name:                "diced onion",
							RecipeStepProductID: pointer.To(dicedOnionRecipeStepProductID),
						},
						{
							Name:                "diced carrot",
							RecipeStepProductID: pointer.To(dicedCarrotRecipeStepProductID),
						},
						{
							Name:                "diced celery",
							RecipeStepProductID: pointer.To(dicedCeleryRecipeStepProductID),
						},
					},
					Products: []*types.RecipeStepProduct{
						{
							ID:   dicedOnionRecipeStepProductID,
							Name: "sauteed mire poix",
							Type: types.RecipeStepProductIngredientType,
						},
					},
					Notes: "fourth step",
					Index: 3,
				},
			},
		}

		expected := `flowchart TD;
	Step1["Step #1 (dice)"];
	Step2["Step #2 (dice)"];
	Step3["Step #3 (dice)"];
	Step4["Step #4 (sautee)"];
	Step1 -->|ingredient| Step4;
	Step2 -->|ingredient| Step4;
	Step3 -->|ingredient| Step4;
`
		actual := g.RenderMermaidDiagramForRecipe(ctx, recipe)

		assert.Equal(t, expected, actual)
	})
}
