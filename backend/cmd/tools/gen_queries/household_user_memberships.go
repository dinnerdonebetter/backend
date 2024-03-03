package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	householdUserMembershipsTableName = "household_user_memberships"

	defaultHouseholdColumn = "default_household"
	householdRoleColumn    = "household_role"
)

var householdUserMembershipsColumns = []string{
	idColumn,
	belongsToHouseholdColumn,
	belongsToUserColumn,
	defaultHouseholdColumn,
	householdRoleColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildHouseholdUserMembershipsQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "AddUserToHousehold",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				householdUserMembershipsTableName,
				strings.Join(filterForInsert(householdUserMembershipsColumns, "default_household"), ",\n\t"),
				strings.Join(applyToEach(filterForInsert(householdUserMembershipsColumns, "default_household"), func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateHouseholdUserMembershipForNewUser",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				householdUserMembershipsTableName,
				strings.Join(filterForInsert(householdUserMembershipsColumns), ",\n\t"),
				strings.Join(applyToEach(filterForInsert(householdUserMembershipsColumns), func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetDefaultHouseholdIDForUser",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT %s.%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s = sqlc.arg(%s)
	AND %s.%s = TRUE;`,
				householdsTableName, idColumn,
				householdsTableName,
				householdUserMembershipsTableName, householdUserMembershipsTableName, belongsToHouseholdColumn, householdsTableName, idColumn,
				householdUserMembershipsTableName, belongsToUserColumn, belongsToUserColumn,
				householdUserMembershipsTableName, defaultHouseholdColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdUserMembershipsForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(householdUserMembershipsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", householdUserMembershipsTableName, s)
				}), ",\n\t"),
				householdUserMembershipsTableName,
				householdsTableName, householdsTableName, idColumn, householdUserMembershipsTableName, belongsToHouseholdColumn,
				householdUserMembershipsTableName, archivedAtColumn,
				householdUserMembershipsTableName, belongsToUserColumn, belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "MarkHouseholdUserMembershipAsUserDefault",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = (%s = sqlc.arg(%s) AND %s = sqlc.arg(%s))
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				householdUserMembershipsTableName,
				defaultHouseholdColumn, belongsToUserColumn, belongsToUserColumn, belongsToHouseholdColumn, belongsToHouseholdColumn,
				archivedAtColumn,
				belongsToUserColumn, belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ModifyHouseholdUserPermissions",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = sqlc.arg(%s)
WHERE %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdUserMembershipsTableName,
				householdRoleColumn, householdRoleColumn,
				belongsToHouseholdColumn, belongsToHouseholdColumn,
				belongsToUserColumn, belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "RemoveUserFromHousehold",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s,
	%s = 'false'
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
				householdUserMembershipsTableName,
				archivedAtColumn, currentTimeExpression,
				defaultHouseholdColumn,
				householdUserMembershipsTableName, archivedAtColumn,
				householdUserMembershipsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn,
				householdUserMembershipsTableName, belongsToUserColumn, belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "TransferHouseholdMembership",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = sqlc.arg(%s)
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdUserMembershipsTableName,
				belongsToUserColumn,
				belongsToUserColumn,
				archivedAtColumn,
				belongsToHouseholdColumn,
				belongsToHouseholdColumn,
				belongsToUserColumn,
				belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "TransferHouseholdOwnership",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = sqlc.arg(new_owner)
WHERE %s IS NULL
	AND %s = sqlc.arg(old_owner)
	AND %s = sqlc.arg(household_id);`,
				householdsTableName,
				belongsToUserColumn,
				archivedAtColumn,
				belongsToUserColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UserIsHouseholdMember",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
);`,
				householdUserMembershipsTableName, idColumn,
				householdUserMembershipsTableName,
				householdUserMembershipsTableName, archivedAtColumn,
				householdUserMembershipsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn,
				householdUserMembershipsTableName, belongsToUserColumn, belongsToUserColumn,
			)),
		},
	}
}
