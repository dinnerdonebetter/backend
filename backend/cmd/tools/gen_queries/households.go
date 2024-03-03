package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	householdsTableName = "households"

	/* #nosec G101 */
	webhookHMACSecretColumn = "webhook_hmac_secret"
)

var householdsColumns = []string{
	idColumn,
	nameColumn,
	"billing_status",
	"contact_phone",
	"payment_processor_customer_id",
	"subscription_plan_id",
	belongsToUserColumn,
	"time_zone",
	"address_line_1",
	"address_line_2",
	"city",
	"state",
	"zip_code",
	"country",
	"latitude",
	"longitude",
	"last_payment_provider_sync_occurred_at",
	webhookHMACSecretColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildHouseholdsQueries() []*Query {
	insertColumns := filterForInsert(householdsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "AddToHouseholdDuringCreation",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO household_user_memberships (
	%s
) VALUES (
	%s
);`,
				strings.Join(filterForInsert(householdUserMembershipsColumns, "default_household"), ",\n\t"),
				strings.Join(applyToEach(filterForInsert(householdUserMembershipsColumns, "default_household"), func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveHousehold",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdsTableName,
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToUserColumn,
				belongsToUserColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateHousehold",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				householdsTableName,
				strings.Join(filterForInsert(
					insertColumns,
					"time_zone",
					"payment_processor_customer_id",
					"last_payment_provider_sync_occurred_at",
					"subscription_plan_id",
				), ",\n\t"),
				strings.Join(applyToEach(filterForInsert(
					insertColumns,
					"time_zone",
					"payment_processor_customer_id",
					"last_payment_provider_sync_occurred_at",
					"subscription_plan_id",
				), func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdByIDWithMemberships",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(append(
					append(
						applyToEach(householdsColumns, func(_ int, s string) string {
							return fmt.Sprintf("%s.%s", householdsTableName, s)
						}),
						applyToEach(usersColumns, func(_ int, s string) string {
							return fmt.Sprintf("%s.%s as user_%s", usersTableName, s, s)
						})...,
					),
					applyToEach(householdUserMembershipsColumns, func(_ int, s string) string {
						return fmt.Sprintf("%s.%s as membership_%s", householdUserMembershipsTableName, s, s)
					})...,
				), ",\n\t"),
				householdsTableName,
				householdUserMembershipsTableName, householdUserMembershipsTableName, belongsToHouseholdColumn, householdsTableName, idColumn,
				usersTableName, householdUserMembershipsTableName, belongsToUserColumn, usersTableName, idColumn,
				householdsTableName, archivedAtColumn,
				householdUserMembershipsTableName, archivedAtColumn,
				householdsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdsForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	(
		SELECT COUNT(%s.%s)
		FROM %s
			JOIN %s ON %s.%s = %s.%s
		WHERE %s.%s IS NULL
			AND %s.%s = sqlc.arg(%s)%s
	) as filtered_count,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(applyToEach(householdsColumns, func(_ int, s string) string {
					return fmt.Sprintf("%s.%s", householdsTableName, s)
				}), ",\n\t"),
				householdsTableName, idColumn,
				householdsTableName,
				householdUserMembershipsTableName, householdUserMembershipsTableName, belongsToHouseholdColumn, householdsTableName, idColumn,
				householdsTableName, archivedAtColumn,
				householdUserMembershipsTableName, belongsToUserColumn, belongsToUserColumn,
				strings.Join(applyToEach(strings.Split(buildFilterConditions(
					householdsTableName,
					true,
				), "\n"), func(i int, s string) string {
					if i == 0 {
						return fmt.Sprintf("\n\t\t\t%s", s)
					}
					return fmt.Sprintf("\n\t\t%s", s)
				}), ""),
				buildTotalCountSelect(householdsTableName, true),
				householdsTableName,
				householdUserMembershipsTableName, householdUserMembershipsTableName, belongsToHouseholdColumn, householdsTableName, idColumn,
				usersTableName, householdUserMembershipsTableName, belongsToUserColumn, usersTableName, idColumn,
				householdsTableName, archivedAtColumn,
				householdUserMembershipsTableName, archivedAtColumn,
				householdUserMembershipsTableName, belongsToUserColumn, belongsToUserColumn,
				buildFilterConditions(
					householdsTableName,
					true,
				),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateHousehold",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdsTableName,
				strings.Join(
					applyToEach(
						filterForUpdate(
							householdsColumns,
							"billing_status",
							"payment_processor_customer_id",
							"subscription_plan_id",
							belongsToUserColumn,
							"time_zone",
							"last_payment_provider_sync_occurred_at",
							"webhook_hmac_secret",
						),
						func(_ int, s string) string {
							return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
						},
					),
					",\n\t",
				),
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToUserColumn,
				belongsToUserColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateHouseholdWebhookEncryptionKey",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = sqlc.arg(%s),
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdsTableName,
				webhookHMACSecretColumn, webhookHMACSecretColumn,
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				belongsToUserColumn, belongsToUserColumn,
				idColumn, idColumn,
			)),
		},
	}
}
