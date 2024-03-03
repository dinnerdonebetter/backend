package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeStepIngredientsTableName = "recipe_step_ingredients"

	recipeStepIngredientIDColumn = "recipe_step_ingredient_id"
	ingredientIDColumn           = "ingredient_id"
	measurementUnitColumn        = "measurement_unit"
)

var recipeStepIngredientsColumns = []string{
	idColumn,
	nameColumn,
	"optional",
	ingredientIDColumn,
	measurementUnitColumn,
	"minimum_quantity_value",
	"maximum_quantity_value",
	"quantity_notes",
	"recipe_step_product_id",
	"ingredient_notes",
	"option_index",
	"to_taste",
	"product_percentage_to_use",
	"vessel_index",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	"recipe_step_product_recipe_id",
	belongsToRecipeStepColumn,
}

func buildRecipeStepIngredientsQueries() []*Query {
	insertColumns := filterForInsert(recipeStepIngredientsColumns)

	fullSelectColumn := mergeColumns(
		applyToEach(filterFromSlice(recipeStepIngredientsColumns, ingredientIDColumn, measurementUnitColumn), func(i int, s string) string {
			return fmt.Sprintf("%s.%s", recipeStepIngredientsTableName, s)
		}),
		append(
			applyToEach(validIngredientsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
			}),
			applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_measurement_unit_%s", validMeasurementUnitsTableName, s, s)
			})...,
		),
		3,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveRecipeStepIngredient",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				recipeStepIngredientsTableName,
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
				Name: "CreateRecipeStepIngredient",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipeStepIngredientsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckRecipeStepIngredientExistence",
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
				recipeStepIngredientsTableName, idColumn,
				recipeStepIngredientsTableName,
				recipeStepsTableName, recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepIngredientsTableName, archivedAtColumn,
				recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeStepIngredientsTableName, idColumn, recipeStepIngredientIDColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAllRecipeStepIngredientsForRecipe",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	LEFT JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumn, ",\n\t"),
				recipeStepIngredientsTableName,
				recipeStepsTableName, recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				validIngredientsTableName, recipeStepIngredientsTableName, ingredientIDColumn, validIngredientsTableName, idColumn,
				validMeasurementUnitsTableName, recipeStepIngredientsTableName, measurementUnitColumn, validMeasurementUnitsTableName, idColumn,
				recipeStepIngredientsTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepIngredients",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	(
		SELECT COUNT(%s.%s)
		FROM %s
			JOIN %s ON %s.%s = %s.%s
			JOIN %s ON %s.%s = %s.%s
		WHERE
			%s.%s IS NULL
			AND %s.%s = sqlc.arg(%s)
			AND %s.%s = sqlc.arg(%s)
			AND %s.%s = sqlc.arg(%s)
			AND %s.%s = sqlc.arg(%s)
			%s
	) as filtered_count,
	(
		SELECT COUNT(%s.%s)
		FROM %s
			JOIN %s ON %s.%s = %s.%s
			JOIN %s ON %s.%s = %s.%s
		WHERE %s.%s IS NULL
			AND %s.%s = sqlc.arg(%s)
			AND %s.%s = sqlc.arg(%s)
	) as total_count
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	LEFT JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumn, ",\n\t"),
				recipeStepIngredientsTableName, idColumn,
				recipeStepIngredientsTableName,
				recipeStepsTableName, recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepIngredientsTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				strings.Join(strings.Split(buildFilterConditions(recipeStepIngredientsTableName, true), "\n"), "\n\t\t"),

				//

				recipeStepIngredientsTableName, idColumn,
				recipeStepIngredientsTableName,
				recipeStepsTableName, recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepIngredientsTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
				recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				//

				recipeStepIngredientsTableName,
				recipeStepsTableName, recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				validIngredientsTableName, recipeStepIngredientsTableName, ingredientIDColumn, validIngredientsTableName, idColumn,
				validMeasurementUnitsTableName, recipeStepIngredientsTableName, measurementUnitColumn, validMeasurementUnitsTableName, idColumn,

				//

				//

				recipeStepIngredientsTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				buildFilterConditions(recipeStepIngredientsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepIngredient",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	LEFT JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumn, ",\n\t"),
				recipeStepIngredientsTableName,
				recipeStepsTableName, recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				validIngredientsTableName, recipeStepIngredientsTableName, ingredientIDColumn, validIngredientsTableName, idColumn,
				validMeasurementUnitsTableName, recipeStepIngredientsTableName, measurementUnitColumn, validMeasurementUnitsTableName, idColumn,
				recipeStepIngredientsTableName, archivedAtColumn,
				recipeStepIngredientsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeStepIngredientsTableName, idColumn, recipeStepIngredientIDColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipeStepIngredient",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				recipeStepIngredientsTableName,
				strings.Join(applyToEach(filterForUpdate(recipeStepIngredientsColumns, belongsToRecipeStepColumn), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				belongsToRecipeStepColumn, belongsToRecipeStepColumn,
				idColumn, idColumn,
			)),
		},
	}
}
