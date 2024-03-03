package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeStepVesselsTableName = "recipe_step_vessels"

	recipeStepVesselIDColumn = "recipe_step_vessel_id"
)

var recipeStepVesselsColumns = []string{
	idColumn,
	nameColumn,
	notesColumn,
	belongsToRecipeStepColumn,
	"recipe_step_product_id",
	"valid_vessel_id",
	"vessel_predicate",
	"minimum_quantity",
	"maximum_quantity",
	"unavailable_after_step",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepVesselsQueries() []*Query {
	insertColumns := filterForInsert(recipeStepVesselsColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(filterFromSlice(recipeStepVesselsColumns, validVesselIDColumn), func(i int, s string) string {
			return fmt.Sprintf("%s.%s", recipeStepVesselsTableName, s)
		}),
		mergeColumns(
			applyToEach(filterFromSlice(validVesselsColumns, capacityUnitColumn), func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_vessel_%s", validVesselsTableName, s, s)
			}),
			applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_measurement_unit_%s", validMeasurementUnitsTableName, s, s)
			}),
			10,
		),
		1,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveRecipeStepVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				recipeStepVesselsTableName,
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
				Name: "CreateRecipeStepVessel",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipeStepVesselsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckRecipeStepVesselExistence",
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
				recipeStepVesselsTableName, idColumn,
				recipeStepVesselsTableName,
				recipeStepsTableName, recipeStepVesselsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepVesselsTableName, archivedAtColumn,
				recipeStepVesselsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeStepVesselsTableName, idColumn, recipeStepVesselIDColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepVesselsForRecipe",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	LEFT JOIN %s ON %s.%s=%s.%s
	LEFT JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				recipeStepVesselsTableName,
				validVesselsTableName, recipeStepVesselsTableName, validVesselIDColumn, validVesselsTableName, idColumn,
				validMeasurementUnitsTableName, validVesselsTableName, capacityUnitColumn, validMeasurementUnitsTableName, idColumn,
				recipeStepsTableName, recipeStepVesselsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepVesselsTableName, archivedAtColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepVessel",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	LEFT JOIN %s ON %s.%s=%s.%s
	LEFT JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				recipeStepVesselsTableName,
				validVesselsTableName, recipeStepVesselsTableName, validVesselIDColumn, validVesselsTableName, idColumn,
				validMeasurementUnitsTableName, validVesselsTableName, capacityUnitColumn, validMeasurementUnitsTableName, idColumn,
				recipeStepsTableName, recipeStepVesselsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepVesselsTableName, archivedAtColumn,
				recipeStepVesselsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeStepVesselsTableName, idColumn, recipeStepVesselIDColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeStepVessels",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	 LEFT JOIN %s ON %s.%s=%s.%s
	 LEFT JOIN %s ON %s.%s=%s.%s
	 JOIN %s ON %s.%s=%s.%s
	 JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(recipeStepVesselsTableName, true, true),
				buildTotalCountSelect(recipeStepVesselsTableName, true),
				recipeStepVesselsTableName,
				validVesselsTableName, recipeStepVesselsTableName, validVesselIDColumn, validVesselsTableName, idColumn,
				validMeasurementUnitsTableName, validVesselsTableName, capacityUnitColumn, validMeasurementUnitsTableName, idColumn,
				recipeStepsTableName, recipeStepVesselsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
				recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
				recipeStepVesselsTableName, archivedAtColumn,
				recipeStepVesselsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeStepsTableName, archivedAtColumn,
				recipeStepsTableName, idColumn, recipeStepIDColumn,
				recipesTableName, archivedAtColumn,
				recipesTableName, idColumn, recipeIDColumn,
				buildFilterConditions(recipeStepVesselsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipeStepVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				recipeStepVesselsTableName,
				strings.Join(applyToEach(filterForUpdate(recipeStepVesselsColumns), func(i int, s string) string {
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
