package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	householdInstrumentOwnershipsTableName = "household_instrument_ownerships"
)

var householdInstrumentOwnershipsColumns = []string{
	idColumn,
	notesColumn,
	"quantity",
	validInstrumentIDColumn,
	belongsToHouseholdColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildHouseholdInstrumentOwnershipQueries() []*Query {
	insertColumns := filterForInsert(householdInstrumentOwnershipsColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(filterFromSlice(householdInstrumentOwnershipsColumns, validInstrumentIDColumn), func(i int, s string) string {
			return fmt.Sprintf("%s.%s", householdInstrumentOwnershipsTableName, s)
		}),
		applyToEach(validInstrumentsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_instrument_%s", validInstrumentsTableName, s, s)
		}),
		3,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveHouseholdInstrumentOwnership",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdInstrumentOwnershipsTableName,
				archivedAtColumn, currentTimeExpression,
				archivedAtColumn,
				idColumn, idColumn,
				belongsToHouseholdColumn, belongsToHouseholdColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateHouseholdInstrumentOwnership",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				householdInstrumentOwnershipsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckHouseholdInstrumentOwnershipExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.id
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
);`,
				householdInstrumentOwnershipsTableName,
				householdInstrumentOwnershipsTableName,
				householdInstrumentOwnershipsTableName,
				archivedAtColumn,
				householdInstrumentOwnershipsTableName, idColumn, idColumn,
				householdInstrumentOwnershipsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdInstrumentOwnerships",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
INNER JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	%s
GROUP BY
	%s.%s,
	%s.%s
ORDER BY
	%s.%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(householdInstrumentOwnershipsTableName, true, true, "household_instrument_ownerships.belongs_to_household = sqlc.arg(household_id)"),
				buildTotalCountSelect(householdInstrumentOwnershipsTableName, true, "household_instrument_ownerships.belongs_to_household = sqlc.arg(household_id)"),
				householdInstrumentOwnershipsTableName,
				validInstrumentsTableName, householdInstrumentOwnershipsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				householdInstrumentOwnershipsTableName, archivedAtColumn,
				buildFilterConditions(
					householdInstrumentOwnershipsTableName,
					true,
					"household_instrument_ownerships.belongs_to_household = sqlc.arg(household_id)",
				),
				householdInstrumentOwnershipsTableName, idColumn,
				validInstrumentsTableName, idColumn,
				householdInstrumentOwnershipsTableName, idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdInstrumentOwnership",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
INNER JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				householdInstrumentOwnershipsTableName,
				validInstrumentsTableName, householdInstrumentOwnershipsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
				householdInstrumentOwnershipsTableName,
				archivedAtColumn,
				householdInstrumentOwnershipsTableName, idColumn, idColumn,
				householdInstrumentOwnershipsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateHouseholdInstrumentOwnership",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
				householdInstrumentOwnershipsTableName,
				strings.Join(applyToEach(filterForUpdate(householdInstrumentOwnershipsColumns, belongsToHouseholdColumn), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn, idColumn,
				householdInstrumentOwnershipsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn,
			)),
		},
	}
}
