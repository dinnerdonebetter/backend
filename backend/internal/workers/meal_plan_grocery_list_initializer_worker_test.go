package workers

import (
	"context"
	"testing"

	analyticsmock "github.com/dinnerdonebetter/backend/internal/analytics/mock"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/features/grocerylistpreparation"
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProvideMealPlanGroceryListInitializer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideMealPlanGroceryListInitializer(
			logging.NewNoopLogger(),
			&database.MockDatabase{},
			&recipeanalysis.MockRecipeAnalyzer{},
			&mockpublishers.Publisher{},
			&analyticsmock.EventReporter{},
			tracing.NewNoopTracerProvider(),
			&grocerylistpreparation.MockGroceryListCreator{},
		)
		assert.NotNil(t, actual)
	})
}

func TestMealPlanGroceryListInitializer_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := ProvideMealPlanGroceryListInitializer(
			logging.NewNoopLogger(),
			&database.MockDatabase{},
			&recipeanalysis.MockRecipeAnalyzer{},
			&mockpublishers.Publisher{},
			&analyticsmock.EventReporter{},
			tracing.NewNoopTracerProvider(),
			&grocerylistpreparation.MockGroceryListCreator{},
		).(*mealPlanGroceryListInitializer)
		assert.NotNil(t, w)

		onion := fakes.BuildFakeValidIngredient()
		carrot := fakes.BuildFakeValidIngredient()
		celery := fakes.BuildFakeValidIngredient()
		salt := fakes.BuildFakeValidIngredient()

		grams := fakes.BuildFakeValidMeasurementUnit()

		expectedMealPlans := []*types.MealPlan{
			{
				ID: fakes.BuildFakeID(),
				Events: []*types.MealPlanEvent{
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      onion,
																MinimumQuantity: 100,
																MaximumQuantity: pointer.To(float32(100)),
																MeasurementUnit: *grams,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      carrot,
																MinimumQuantity: 100,
																MaximumQuantity: pointer.To(float32(100)),
																MeasurementUnit: *grams,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      celery,
																MinimumQuantity: 100,
																MaximumQuantity: pointer.To(float32(100)),
																MeasurementUnit: *grams,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      salt,
																MinimumQuantity: 100,
																MaximumQuantity: pointer.To(float32(100)),
																MeasurementUnit: *grams,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Options: []*types.MealPlanOption{
							{
								Chosen: true,
								Meal: types.Meal{
									Components: []*types.MealComponent{
										{
											Recipe: types.Recipe{
												Steps: []*types.RecipeStep{
													{
														Ingredients: []*types.RecipeStepIngredient{
															{
																Ingredient:      onion,
																MinimumQuantity: 100,
																MaximumQuantity: pointer.To(float32(100)),
																MeasurementUnit: *grams,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		ctx := context.Background()
		mdm := database.NewMockDatabase()

		mdm.MealPlanDataManagerMock.On("GetFinalizedMealPlansWithUninitializedGroceryLists", testutils.ContextMatcher).Return(expectedMealPlans, nil)

		firstMealPlanExpectedGroceryListItemInputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      onion.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  200,
				MaximumQuantityNeeded:  pointer.To(float32(200)),
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      carrot.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  100,
				MaximumQuantityNeeded:  pointer.To(float32(100)),
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      celery.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  100,
				MaximumQuantityNeeded:  pointer.To(float32(100)),
			},
			{
				Status:                 types.MealPlanGroceryListItemStatusUnknown,
				ValidMeasurementUnitID: grams.ID,
				ValidIngredientID:      salt.ID,
				BelongsToMealPlan:      expectedMealPlans[0].ID,
				MinimumQuantityNeeded:  100,
				MaximumQuantityNeeded:  pointer.To(float32(100)),
			},
		}

		expectedInputSets := map[string][]*types.MealPlanGroceryListItemDatabaseCreationInput{
			expectedMealPlans[0].ID: firstMealPlanExpectedGroceryListItemInputs,
		}

		mglm := &grocerylistpreparation.MockGroceryListCreator{}
		mglm.On(
			"GenerateGroceryListInputs",
			testutils.ContextMatcher,
			expectedMealPlans[0],
		).Return(firstMealPlanExpectedGroceryListItemInputs, nil)
		w.groceryListCreator = mglm

		for _, inputs := range expectedInputSets {
			for _, input := range inputs {
				mdm.MealPlanGroceryListItemDataManagerMock.On("CreateMealPlanGroceryListItem", testutils.ContextMatcher, input).Return((*types.MealPlanGroceryListItem)(nil), nil)
			}
		}
		w.dataManager = mdm

		assert.NoError(t, w.InitializeGroceryListsForFinalizedMealPlans(ctx, []byte("{}")))
		mock.AssertExpectationsForObjects(t, mdm)
	})
}
