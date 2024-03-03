package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeStepCompletionConditionsTableName = "recipe_step_completion_conditions"

	ingredientStateColumn                 = "ingredient_state"
	recipeStepCompletionConditionIDColumn = "recipe_step_completion_condition_id"
)

var recipeStepCompletionConditionsColumns = []string{
	idColumn,
	"optional",
	notesColumn,
	belongsToRecipeStepColumn,
	ingredientStateColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepCompletionConditionQueries() []*Query {
	insertColumns := filterForInsert(recipeStepCompletionConditionsColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(recipeStepCompletionConditionIngredientsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as recipe_step_completion_condition_ingredient_%s", recipeStepCompletionConditionIngredientsTableName, s, s)
		}),
		mergeColumns(
			applyToEach(recipeStepCompletionConditionsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s", recipeStepCompletionConditionsTableName, s)
			}),
			applyToEach(validIngredientStatesColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_state_%s", validIngredientStatesTableName, s, s)
			}),
			2,
		),
		3,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveRecipeStepCompletionCondition",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				recipeStepCompletionConditionsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToRecipeStepColumn,
				belongsToRecipeStepColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateRecipeStepCompletionCondition",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipeStepCompletionConditionsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckRecipeStepCompletionConditionExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
		JOIN %s ON %s.%s=%s.%s
		JOIN %s ON %s.%s=%s.%s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				recipeStepCompletionConditionsTableName, idColumn,
				recipeStepCompletionConditionsTableName,
				recipeStepsTableName, recipeStepCompletionConditionsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepCompletionConditionsTableName, archivedAtColumn,
				recipeStepCompletionConditionsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeStepCompletionConditionsTableName, idColumn, recipeStepCompletionConditionIDColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAllRecipeStepCompletionConditionsForRecipe",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
GROUP BY
	%s.%s,
	%s.%s,
	%s.%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				recipeStepCompletionConditionIngredientsTableName,
				recipeStepCompletionConditionsTableName, recipeStepCompletionConditionIngredientsTableName, belongsToRecipeStepCompletionConditionColumn, recipeStepCompletionConditionsTableName, idColumn,
				recipeStepsTableName, recipeStepCompletionConditionsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				validIngredientStatesTableName, recipeStepCompletionConditionsTableName, ingredientStateColumn, validIngredientStatesTableName, idColumn,
				recipeStepCompletionConditionsTableName, archivedAtColumn,
				recipeStepCompletionConditionIngredientsTableName, archivedAtColumn,
				recipeStepsTableName, archivedAtColumn,
				recipesTableName, archivedAtColumn,
				validIngredientStatesTableName, archivedAtColumn,
				recipesTableName, idColumn, idColumn,
				recipeStepCompletionConditionsTableName, idColumn,
				recipeStepCompletionConditionIngredientsTableName, idColumn,
				validIngredientStatesTableName, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepCompletionConditions",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(recipeStepCompletionConditionIngredientsTableName, true, true, "recipe_step_completion_conditions.belongs_to_recipe_step = sqlc.arg(recipe_step_id)"),
				buildTotalCountSelect(recipeStepCompletionConditionIngredientsTableName, true),
				recipeStepCompletionConditionIngredientsTableName,
				recipeStepCompletionConditionsTableName, recipeStepCompletionConditionIngredientsTableName, belongsToRecipeStepCompletionConditionColumn, recipeStepCompletionConditionsTableName, idColumn,
				recipeStepsTableName, recipeStepCompletionConditionsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				validIngredientStatesTableName, recipeStepCompletionConditionsTableName, ingredientStateColumn, validIngredientStatesTableName, idColumn,
				recipeStepCompletionConditionsTableName, archivedAtColumn,
				buildFilterConditions(
					recipeStepCompletionConditionIngredientsTableName,
					true,
					"recipe_step_completion_conditions.belongs_to_recipe_step = sqlc.arg(recipe_step_id)",
				),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepCompletionConditionWithIngredients",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(recipe_step_completion_condition_id)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				recipeStepCompletionConditionIngredientsTableName,
				recipeStepCompletionConditionsTableName, recipeStepCompletionConditionIngredientsTableName, belongsToRecipeStepCompletionConditionColumn, recipeStepCompletionConditionsTableName, idColumn,
				recipeStepsTableName, recipeStepCompletionConditionsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				validIngredientStatesTableName, recipeStepCompletionConditionsTableName, ingredientStateColumn, validIngredientStatesTableName, idColumn,
				recipeStepCompletionConditionsTableName, archivedAtColumn,
				recipeStepCompletionConditionIngredientsTableName, archivedAtColumn,
				recipeStepCompletionConditionsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeStepCompletionConditionsTableName, idColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipeStepCompletionCondition",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				recipeStepCompletionConditionsTableName,
				strings.Join(applyToEach(filterForUpdate(recipeStepCompletionConditionsColumns), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				idColumn, idColumn,
			)),
		},
	}
}
