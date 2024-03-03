package grocerylistpreparation

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/shopspring/decimal"
)

// GroceryListCreator creates meal plan grocery lists for a given meal plan.
type GroceryListCreator interface {
	GenerateGroceryListInputs(ctx context.Context, mealPlan *types.MealPlan) ([]*types.MealPlanGroceryListItemDatabaseCreationInput, error)
}

type groceryListCreator struct {
	logger logging.Logger
	tracer tracing.Tracer
}

func NewGroceryListCreator(logger logging.Logger, tracerProvider tracing.TracerProvider) GroceryListCreator {
	return &groceryListCreator{
		logger: logging.EnsureLogger(logger).WithName("grocery_list_creator"),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("grocery_list_creator")),
	}
}

func (g *groceryListCreator) GenerateGroceryListInputs(ctx context.Context, mealPlan *types.MealPlan) ([]*types.MealPlanGroceryListItemDatabaseCreationInput, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	inputs := map[string]*types.MealPlanGroceryListItemDatabaseCreationInput{}
	logger := g.logger.Clone().WithValue(keys.MealPlanIDKey, mealPlan.ID)

	for _, event := range mealPlan.Events {
		logger = logger.WithValue(keys.MealPlanEventIDKey, event.ID)
		for _, option := range event.Options {
			if option.Chosen {
				mealScale := decimal.NewFromFloat32(option.MealScale)
				logger = logger.WithValue(keys.MealPlanOptionIDKey, option.ID)
				for _, component := range option.Meal.Components {
					recipeScale := decimal.NewFromFloat32(component.RecipeScale).Mul(mealScale)
					logger = logger.WithValue(keys.RecipeIDKey, component.Recipe.ID)
					for _, step := range component.Recipe.Steps {
						logger = logger.WithValue(keys.RecipeStepIDKey, step.ID)
						for _, ingredient := range step.Ingredients {
							if ingredient.Ingredient != nil {
								logger = logger.WithValue(keys.RecipeStepIngredientIDKey, ingredient.ID)
								if _, ok := inputs[ingredient.Ingredient.ID]; !ok {
									minQty := float32(recipeScale.Mul(decimal.NewFromFloat32(ingredient.MinimumQuantity)).Truncate(2).InexactFloat64())
									var maxQty *float32
									if ingredient.MaximumQuantity != nil {
										max := float32(recipeScale.Mul(decimal.NewFromFloat32(*ingredient.MaximumQuantity)).Truncate(2).InexactFloat64())
										maxQty = &max
									}

									inputs[ingredient.Ingredient.ID] = &types.MealPlanGroceryListItemDatabaseCreationInput{
										Status:                 types.MealPlanGroceryListItemStatusNeeds,
										ValidMeasurementUnitID: ingredient.MeasurementUnit.ID,
										ValidIngredientID:      ingredient.Ingredient.ID,
										BelongsToMealPlan:      mealPlan.ID,
										ID:                     identifiers.New(),
										MinimumQuantityNeeded:  minQty,
										MaximumQuantityNeeded:  maxQty,
									}
								} else {
									if inputs[ingredient.Ingredient.ID].ValidMeasurementUnitID == ingredient.MeasurementUnit.ID {
										inputs[ingredient.Ingredient.ID].MinimumQuantityNeeded += ingredient.MinimumQuantity

										if inputs[ingredient.Ingredient.ID].MaximumQuantityNeeded != nil {
											if ingredient.MaximumQuantity != nil {
												*inputs[ingredient.Ingredient.ID].MaximumQuantityNeeded += *ingredient.MaximumQuantity
											}
										} else if ingredient.MaximumQuantity != nil {
											inputs[ingredient.Ingredient.ID].MaximumQuantityNeeded = ingredient.MaximumQuantity
										}
									} else {
										logger.Error(fmt.Errorf("mismatched measurement units: %s and %s", inputs[ingredient.Ingredient.ID].ValidMeasurementUnitID, ingredient.MeasurementUnit.ID), "creating grocery list")
									}
								}
							}
						}
					}
				}
			}
		}
	}

	dbInputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{}
	for _, i := range inputs {
		dbInputs = append(dbInputs, i)
	}

	return dbInputs, nil
}
