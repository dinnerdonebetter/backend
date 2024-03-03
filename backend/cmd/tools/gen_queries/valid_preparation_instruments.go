package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validPreparationInstrumentsTableName = "valid_preparation_instruments"
)

var validPreparationInstrumentsColumns = []string{
	idColumn,
	notesColumn,
	validPreparationIDColumn,
	validInstrumentIDColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidPreparationInstrumentsQueries() []*Query {
	insertColumns := filterForInsert(validPreparationInstrumentsColumns)

	fullSelectColumns := mergeColumns(
		mergeColumns(
			applyToEach(filterFromSlice(validPreparationInstrumentsColumns, "valid_preparation_id", "valid_instrument_id"), func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_preparation_instrument_%s", validPreparationInstrumentsTableName, s, s)
			}),
			applyToEach(validInstrumentsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_instrument_%s", validInstrumentsTableName, s, s)
			}),
			2,
		),
		applyToEach(validPreparationsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_preparation_%s", validPreparationsTableName, s, s)
		}),
		2,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidPreparationInstrument",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validPreparationInstrumentsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidPreparationInstrument",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				validPreparationInstrumentsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidPreparationInstrumentExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				validPreparationInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName,
				validPreparationInstrumentsTableName,
				archivedAtColumn,
				validPreparationInstrumentsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstrumentsForInstrument",
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
			AND %s.%s IS NULL
			AND %s.%s IS NULL
			AND %s.%s = sqlc.arg(%s)
			%s
	) as filtered_count,
	(
		SELECT COUNT(%s.%s)
		FROM %s
			JOIN %s ON %s.%s = %s.%s
			JOIN %s ON %s.%s = %s.%s
		WHERE
			%s.%s IS NULL
			AND %s.%s IS NULL
			AND %s.%s IS NULL
			AND %s.%s = sqlc.arg(%s)
	) as total_count
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	%s
GROUP BY
	%s.%s,
	%s.%s,
	%s.%s
ORDER BY %s.%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validPreparationInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validPreparationInstrumentsTableName, validInstrumentIDColumn, idColumn, ///
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				validPreparationInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validPreparationInstrumentsTableName, validInstrumentIDColumn, idColumn, ///
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validPreparationInstrumentsTableName, validInstrumentIDColumn, idColumn, ///
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				validPreparationInstrumentsTableName, idColumn,
				validPreparationsTableName, idColumn,
				validInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName, idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstrumentsForPreparation",
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
			AND %s.%s IS NULL
			AND %s.%s IS NULL
			AND %s.%s = sqlc.arg(%s)
			%s
	) as filtered_count,
	(
		SELECT COUNT(%s.%s)
		FROM %s
			JOIN %s ON %s.%s = %s.%s
			JOIN %s ON %s.%s = %s.%s
		WHERE
			%s.%s IS NULL
			AND %s.%s IS NULL
			AND %s.%s IS NULL
			AND %s.%s = sqlc.arg(%s)
	) as total_count
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	%s
GROUP BY
	%s.%s,
	%s.%s,
	%s.%s
ORDER BY %s.%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validPreparationInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validPreparationInstrumentsTableName, validPreparationIDColumn, idColumn, ///
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				validPreparationInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validPreparationInstrumentsTableName, validPreparationIDColumn, idColumn, ///
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validPreparationInstrumentsTableName, validPreparationIDColumn, idColumn, ///
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				validPreparationInstrumentsTableName, idColumn,
				validPreparationsTableName, idColumn,
				validInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName, idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstruments",
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
			AND %s.%s IS NULL
			AND %s.%s IS NULL
			%s
	) as filtered_count,
	(
		SELECT COUNT(%s.%s)
		FROM %s
			JOIN %s ON %s.%s = %s.%s
			JOIN %s ON %s.%s = %s.%s
		WHERE
			%s.%s IS NULL
			AND %s.%s IS NULL
			AND %s.%s IS NULL
	) as total_count
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	%s
GROUP BY
	%s.%s,
	%s.%s,
	%s.%s
ORDER BY %s.%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validPreparationInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				validPreparationInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				validPreparationInstrumentsTableName, idColumn,
				validPreparationsTableName, idColumn,
				validInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName, idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstrument",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM
	%s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validPreparationInstrumentsTableName,
				validInstrumentsTableName, validPreparationInstrumentsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				validPreparationsTableName, validPreparationInstrumentsTableName, validPreparationIDColumn, validPreparationsTableName, idColumn,
				validPreparationInstrumentsTableName, archivedAtColumn,
				validInstrumentsTableName, archivedAtColumn,
				validPreparationsTableName, archivedAtColumn,
				validPreparationInstrumentsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ValidPreparationInstrumentPairIsValid",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s.%s
	FROM %s
	WHERE %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s)
	AND %s IS NULL
);`,
				validPreparationInstrumentsTableName, idColumn,
				validPreparationInstrumentsTableName,
				validInstrumentIDColumn, validInstrumentIDColumn,
				validPreparationIDColumn, validPreparationIDColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidPreparationInstrument",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validPreparationInstrumentsTableName,
				strings.Join(applyToEach(filterForUpdate(validPreparationInstrumentsColumns), func(i int, s string) string {
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
