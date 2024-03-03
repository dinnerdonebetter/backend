package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeStepCompletionConditionIngredientsTableName = "recipe_step_completion_condition_ingredients"

	belongsToRecipeStepCompletionConditionColumn = "belongs_to_recipe_step_completion_condition"
)

var recipeStepCompletionConditionIngredientsColumns = []string{
	idColumn,
	belongsToRecipeStepCompletionConditionColumn,
	"recipe_step_ingredient",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepCompletionConditionIngredientsQueries() []*Query {
	insertColumns := filterForInsert(recipeStepCompletionConditionIngredientsColumns)

	fullSelectColumns := append(
		applyToEach(recipeStepCompletionConditionIngredientsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as %s_%s", recipeStepCompletionConditionIngredientsTableName, s, strings.TrimSuffix(recipeStepCompletionConditionIngredientsTableName, "s"), s)
		}),
		applyToEach(validIngredientStatesColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as %s_%s", validIngredientStatesTableName, s, strings.TrimSuffix(validIngredientStatesTableName, "s"), s)
		})...,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "CreateRecipeStepCompletionConditionIngredient",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipeStepCompletionConditionIngredientsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAllRecipeStepCompletionConditionIngredientsForRecipeCompletionIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[])
	AND %s.%s IS NULL;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				recipeStepCompletionConditionIngredientsTableName,
				recipeStepCompletionConditionsTableName, recipeStepCompletionConditionIngredientsTableName, belongsToRecipeStepCompletionConditionColumn, recipeStepCompletionConditionsTableName, idColumn,
				validIngredientStatesTableName, recipeStepCompletionConditionsTableName, ingredientStateColumn, validIngredientStatesTableName, idColumn,
				recipeStepCompletionConditionsTableName, archivedAtColumn,
				recipeStepCompletionConditionIngredientsTableName, archivedAtColumn,
				recipeStepCompletionConditionIngredientsTableName, belongsToRecipeStepCompletionConditionColumn,
				validIngredientStatesTableName, archivedAtColumn,
			)),
		},
	}
}
