package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	mealPlanEventsTableName = "meal_plan_events"

	belongsToMealPlanEventColumn = "belongs_to_meal_plan_event"
	mealPlanEventIDColumn        = "meal_plan_event_id"
	belongsToMealPlanColumn      = "belongs_to_meal_plan"
)

var mealPlanEventsColumns = []string{
	idColumn,
	notesColumn,
	"starts_at",
	"ends_at",
	"meal_name",
	belongsToMealPlanColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealPlanEventsQueries() []*Query {
	insertColumns := filterForInsert(mealPlanEventsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveMealPlanEvent",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				mealPlanEventsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
				belongsToMealPlanColumn,
				belongsToMealPlanColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateMealPlanEvent",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				mealPlanEventsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "MealPlanEventIsEligibleForVoting",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
		JOIN %s ON %s.%s = %s.%s
	WHERE
		%s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = 'awaiting_votes'
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
);`,
				mealPlanEventsTableName, idColumn,
				mealPlanEventsTableName,
				mealPlansTableName, mealPlanEventsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
				mealPlanEventsTableName, archivedAtColumn,
				mealPlansTableName, idColumn, mealPlanIDColumn,
				mealPlansTableName, mealPlanStatusColumn,
				mealPlansTableName, archivedAtColumn,
				mealPlanEventsTableName, idColumn, mealPlanEventIDColumn,
				mealPlanEventsTableName, archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckMealPlanEventExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
);`,
				mealPlanEventsTableName, idColumn,
				mealPlanEventsTableName,
				mealPlanEventsTableName, archivedAtColumn,
				mealPlanEventsTableName, idColumn, idColumn,
				mealPlanEventsTableName, belongsToMealPlanColumn, mealPlanIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetMealPlanEvents",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE
	%s.%s IS NULL
	%s
GROUP BY %s.%s
ORDER BY %s.%s
%s;`,
				strings.Join(applyToEach(mealPlanEventsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", mealPlanEventsTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(mealPlanEventsTableName, true, true, "meal_plan_events.belongs_to_meal_plan = sqlc.arg(meal_plan_id)"),
				buildTotalCountSelect(mealPlanEventsTableName, true, "meal_plan_events.belongs_to_meal_plan = sqlc.arg(meal_plan_id)"),
				mealPlanEventsTableName,
				mealPlanEventsTableName, archivedAtColumn,
				buildFilterConditions(
					mealPlanEventsTableName,
					true,
					"meal_plan_events.belongs_to_meal_plan = sqlc.arg(meal_plan_id)",
				),
				mealPlanEventsTableName, idColumn,
				mealPlanEventsTableName, idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAllMealPlanEventsForMealPlan",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE
	%s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(mealPlanEventsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", mealPlanEventsTableName, s)
				}), ",\n\t"),
				mealPlanEventsTableName,
				mealPlanEventsTableName, archivedAtColumn,
				mealPlanEventsTableName, belongsToMealPlanColumn, mealPlanIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetMealPlanEvent",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(mealPlanEventsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", mealPlanEventsTableName, s)
				}), ",\n\t"),
				mealPlanEventsTableName,
				mealPlanEventsTableName, archivedAtColumn,
				mealPlanEventsTableName, idColumn, idColumn,
				mealPlanEventsTableName, belongsToMealPlanColumn, belongsToMealPlanColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateMealPlanEvent",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				mealPlanEventsTableName,
				strings.Join(applyToEach(filterForUpdate(mealPlanEventsColumns), func(i int, s string) string {
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
