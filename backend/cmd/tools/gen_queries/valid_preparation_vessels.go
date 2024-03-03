package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validPreparationVesselsTableName = "valid_preparation_vessels"
)

var validPreparationVesselsColumns = []string{
	idColumn,
	notesColumn,
	validPreparationIDColumn,
	validVesselIDColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidPreparationVesselsQueries() []*Query {
	insertColumns := filterForInsert(validPreparationVesselsColumns)

	fullSelectColumns := append(
		mergeColumns(
			applyToEach(filterFromSlice(validPreparationVesselsColumns, "valid_preparation_id", "valid_vessel_id"), func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_preparation_vessel_%s", validPreparationVesselsTableName, s, s)
			}),
			applyToEach(validPreparationsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_preparation_%s", validPreparationsTableName, s, s)
			}),
			2,
		),
		mergeColumns(
			applyToEach(filterFromSlice(validVesselsColumns, "capacity_unit"), func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_vessel_%s", validVesselsTableName, s, s)
			}),
			applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_measurement_unit_%s", validMeasurementUnitsTableName, s, s)
			}),
			10,
		)...,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidPreparationVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validPreparationVesselsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidPreparationVessel",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				validPreparationVesselsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidPreparationVesselExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				validPreparationVesselsTableName, idColumn,
				validPreparationVesselsTableName,
				validPreparationVesselsTableName, archivedAtColumn,
				validPreparationVesselsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationVesselsForPreparation",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	LEFT JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(validPreparationVesselsTableName, true, true),
				buildTotalCountSelect(validPreparationVesselsTableName, true),
				validPreparationVesselsTableName,
				validVesselsTableName, validPreparationVesselsTableName, validVesselIDColumn, validVesselsTableName, idColumn,
				validPreparationsTableName, validPreparationVesselsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validMeasurementUnitsTableName, validVesselsTableName, capacityUnitColumn, validMeasurementUnitsTableName, idColumn,
				validPreparationVesselsTableName, archivedAtColumn,
				validVesselsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validPreparationVesselsTableName, validPreparationIDColumn, idColumn,
				buildFilterConditions(validPreparationVesselsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationVesselsForVessel",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	LEFT JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(validPreparationVesselsTableName, true, true),
				buildTotalCountSelect(validPreparationVesselsTableName, true),
				validPreparationVesselsTableName,

				validVesselsTableName, validPreparationVesselsTableName, validVesselIDColumn, validVesselsTableName, idColumn,
				validPreparationsTableName, validPreparationVesselsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validMeasurementUnitsTableName, validVesselsTableName, capacityUnitColumn, validMeasurementUnitsTableName, idColumn,
				validPreparationVesselsTableName, archivedAtColumn,
				validVesselsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validPreparationVesselsTableName, validVesselIDColumn, idColumn,
				buildFilterConditions(validPreparationVesselsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationVessels",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	LEFT JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(validPreparationVesselsTableName, true, true),
				buildTotalCountSelect(validPreparationVesselsTableName, true),
				validPreparationVesselsTableName,
				validVesselsTableName, validPreparationVesselsTableName, validVesselIDColumn, validVesselsTableName, idColumn,
				validPreparationsTableName, validPreparationVesselsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validMeasurementUnitsTableName, validVesselsTableName, capacityUnitColumn, validMeasurementUnitsTableName, idColumn,
				validPreparationVesselsTableName, archivedAtColumn,
				validVesselsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				buildFilterConditions(validPreparationVesselsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationVessel",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
	LEFT JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validPreparationVesselsTableName,
				validVesselsTableName, validPreparationVesselsTableName, validVesselIDColumn, validVesselsTableName, idColumn,
				validPreparationsTableName, validPreparationVesselsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validMeasurementUnitsTableName, validVesselsTableName, capacityUnitColumn, validMeasurementUnitsTableName, idColumn,
				validPreparationVesselsTableName, archivedAtColumn,
				validVesselsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validMeasurementUnitsTableName, archivedAtColumn,
				validPreparationVesselsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ValidPreparationVesselPairIsValid",
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
				validPreparationVesselsTableName,
				validVesselIDColumn,
				validVesselIDColumn,
				validPreparationIDColumn,
				validPreparationIDColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidPreparationVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validPreparationVesselsTableName,
				strings.Join(applyToEach(filterForUpdate(validPreparationVesselsColumns), func(i int, s string) string {
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
