package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validIngredientMeasurementUnitsTableName = "valid_ingredient_measurement_units"
	validMeasurementUnitColumn               = "valid_measurement_unit"
	validMeasurementUnitIDColumn             = "valid_measurement_unit_id"
)

var validIngredientMeasurementUnitsColumns = []string{
	idColumn,
	notesColumn,
	validMeasurementUnitIDColumn,
	validIngredientIDColumn,
	"minimum_allowable_quantity",
	"maximum_allowable_quantity",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientMeasurementUnitsQueries() []*Query {
	insertColumns := filterForInsert(validIngredientMeasurementUnitsColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(filterFromSlice(validIngredientMeasurementUnitsColumns, "valid_ingredient_id", "valid_measurement_unit_id"), func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_ingredient_measurement_unit_%s", validIngredientMeasurementUnitsTableName, s, s)
		}),
		append(
			applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_measurement_unit_%s", validMeasurementUnitsTableName, s, s)
			}),
			applyToEach(validIngredientsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
			})...),
		2,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidIngredientMeasurementUnit",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validIngredientMeasurementUnitsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidIngredientMeasurementUnit",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				validIngredientMeasurementUnitsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidIngredientMeasurementUnitExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				validIngredientMeasurementUnitsTableName, idColumn,
				validIngredientMeasurementUnitsTableName,
				validIngredientMeasurementUnitsTableName, archivedAtColumn,
				validIngredientMeasurementUnitsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientMeasurementUnitsForIngredient",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(validIngredientMeasurementUnitsTableName, true, true),
				buildTotalCountSelect(validIngredientMeasurementUnitsTableName, true),
				validIngredientMeasurementUnitsTableName,
				validMeasurementUnitsTableName, validIngredientMeasurementUnitsTableName, validMeasurementUnitIDColumn, validMeasurementUnitsTableName, idColumn,
				validIngredientsTableName, validIngredientMeasurementUnitsTableName, validIngredientIDColumn, validIngredientsTableName, idColumn,
				validIngredientMeasurementUnitsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				validIngredientMeasurementUnitsTableName, validIngredientIDColumn, validIngredientIDColumn,
				buildFilterConditions(validIngredientMeasurementUnitsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientMeasurementUnitsForMeasurementUnit",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(validIngredientMeasurementUnitsTableName, true, true),
				buildTotalCountSelect(validIngredientMeasurementUnitsTableName, true),
				validIngredientMeasurementUnitsTableName,
				validMeasurementUnitsTableName, validIngredientMeasurementUnitsTableName, validMeasurementUnitIDColumn, validMeasurementUnitsTableName, idColumn,
				validIngredientsTableName, validIngredientMeasurementUnitsTableName, validIngredientIDColumn, validIngredientsTableName, idColumn,
				validIngredientMeasurementUnitsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				validIngredientMeasurementUnitsTableName, validMeasurementUnitIDColumn, validMeasurementUnitIDColumn,
				buildFilterConditions(validIngredientMeasurementUnitsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientMeasurementUnits",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(validIngredientMeasurementUnitsTableName, true, true),
				buildTotalCountSelect(validIngredientMeasurementUnitsTableName, true),
				validIngredientMeasurementUnitsTableName,
				validMeasurementUnitsTableName, validIngredientMeasurementUnitsTableName, validMeasurementUnitIDColumn, validMeasurementUnitsTableName, idColumn,
				validIngredientsTableName, validIngredientMeasurementUnitsTableName, validIngredientIDColumn, validIngredientsTableName, idColumn,
				validIngredientMeasurementUnitsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				buildFilterConditions(validIngredientMeasurementUnitsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientMeasurementUnit",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validIngredientMeasurementUnitsTableName,
				validMeasurementUnitsTableName, validIngredientMeasurementUnitsTableName, validMeasurementUnitIDColumn, validMeasurementUnitsTableName, idColumn,
				validIngredientsTableName, validIngredientMeasurementUnitsTableName, validIngredientIDColumn, validIngredientsTableName, idColumn,
				validIngredientMeasurementUnitsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				validIngredientMeasurementUnitsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ValidIngredientMeasurementUnitPairIsValid",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s
	FROM %s
	WHERE %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s)
	AND %s IS NULL
);`,
				idColumn,
				validIngredientMeasurementUnitsTableName,
				validMeasurementUnitIDColumn, validMeasurementUnitIDColumn,
				validIngredientIDColumn, validIngredientIDColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidIngredientMeasurementUnit",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validIngredientMeasurementUnitsTableName,
				strings.Join(applyToEach(filterForUpdate(validIngredientMeasurementUnitsColumns), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
	}
}
