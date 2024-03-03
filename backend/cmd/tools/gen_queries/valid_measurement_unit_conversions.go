package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validMeasurementUnitConversionsTableName = "valid_measurement_unit_conversions"

	validMeasurementUnitConversionsFromUnitColumn          = "from_unit"
	validMeasurementUnitConversionsToUnitColumn            = "to_unit"
	validMeasurementUnitConversionsOnlyForIngredientColumn = "only_for_ingredient"
)

var validMeasurementUnitConversionsColumns = []string{
	idColumn,
	"from_unit",
	"to_unit",
	validMeasurementUnitConversionsOnlyForIngredientColumn,
	"modifier",
	notesColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidMeasurementUnitConversionsQueries() []*Query {
	insertColumns := filterForInsert(validMeasurementUnitConversionsColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(validMeasurementUnitConversionsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_measurement_unit_conversion_%s", validMeasurementUnitConversionsTableName, s, s)
		}),
		append(
			append(
				applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s_from.%s as from_unit_%s", validMeasurementUnitsTableName, s, s)
				}),
				applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s_to.%s as to_unit_%s", validMeasurementUnitsTableName, s, s)
				})...,
			),
			applyToEach(validIngredientsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
			})...,
		),
		1,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidMeasurementUnitConversion",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validMeasurementUnitConversionsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidMeasurementUnitConversion",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				validMeasurementUnitConversionsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidMeasurementUnitConversionExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				validMeasurementUnitConversionsTableName, idColumn,
				validMeasurementUnitConversionsTableName,
				validMeasurementUnitConversionsTableName, archivedAtColumn,
				validMeasurementUnitConversionsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAllValidMeasurementUnitConversionsFromMeasurementUnit",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s AS %s_from ON %s.%s = %s_from.%s
	JOIN %s AS %s_to ON %s.%s = %s_to.%s
	LEFT JOIN %s ON %s.%s = %s.%s
WHERE
	%s_from.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s_from.%s IS NULL
	AND %s_to.%s IS NULL;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validMeasurementUnitConversionsTableName,
				validMeasurementUnitsTableName, validMeasurementUnitsTableName, validMeasurementUnitConversionsTableName, validMeasurementUnitConversionsFromUnitColumn, validMeasurementUnitsTableName, idColumn,
				validMeasurementUnitsTableName, validMeasurementUnitsTableName, validMeasurementUnitConversionsTableName, validMeasurementUnitConversionsToUnitColumn, validMeasurementUnitsTableName, idColumn,
				validIngredientsTableName, validMeasurementUnitConversionsTableName, validMeasurementUnitConversionsOnlyForIngredientColumn, validIngredientsTableName, idColumn,
				validMeasurementUnitsTableName, idColumn, idColumn,
				validMeasurementUnitConversionsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAllValidMeasurementUnitConversionsToMeasurementUnit",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s AS %s_from ON %s.%s = %s_from.%s
	JOIN %s AS %s_to ON %s.%s = %s_to.%s
	LEFT JOIN %s ON %s.%s = %s.%s
WHERE
	%s_to.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s_from.%s IS NULL
	AND %s_to.%s IS NULL;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validMeasurementUnitConversionsTableName,
				validMeasurementUnitsTableName, validMeasurementUnitsTableName, validMeasurementUnitConversionsTableName, validMeasurementUnitConversionsFromUnitColumn, validMeasurementUnitsTableName, idColumn,
				validMeasurementUnitsTableName, validMeasurementUnitsTableName, validMeasurementUnitConversionsTableName, validMeasurementUnitConversionsToUnitColumn, validMeasurementUnitsTableName, idColumn,
				validIngredientsTableName, validMeasurementUnitConversionsTableName, validMeasurementUnitConversionsOnlyForIngredientColumn, validIngredientsTableName, idColumn,
				validMeasurementUnitsTableName, idColumn, idColumn,
				validMeasurementUnitConversionsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidMeasurementUnitConversion",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s AS %s_from ON %s.%s = %s_from.%s
	JOIN %s AS %s_to ON %s.%s = %s_to.%s
	LEFT JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s_from.%s IS NULL
	AND %s_to.%s IS NULL;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validMeasurementUnitConversionsTableName,
				validMeasurementUnitsTableName, validMeasurementUnitsTableName, validMeasurementUnitConversionsTableName, validMeasurementUnitConversionsFromUnitColumn, validMeasurementUnitsTableName, idColumn,
				validMeasurementUnitsTableName, validMeasurementUnitsTableName, validMeasurementUnitConversionsTableName, validMeasurementUnitConversionsToUnitColumn, validMeasurementUnitsTableName, idColumn,
				validIngredientsTableName, validMeasurementUnitConversionsTableName, validMeasurementUnitConversionsOnlyForIngredientColumn, validIngredientsTableName, idColumn,
				validMeasurementUnitConversionsTableName, idColumn, idColumn,
				validMeasurementUnitConversionsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidMeasurementUnitConversion",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validMeasurementUnitConversionsTableName,
				strings.Join(applyToEach(filterForUpdate(validMeasurementUnitConversionsColumns), func(i int, s string) string {
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
