package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	mealPlanOptionVotesTableName = "meal_plan_option_votes"

	mealPlanOptionVoteIDColumn    = "meal_plan_option_vote_id"
	belongsToMealPlanOptionColumn = "belongs_to_meal_plan_option"
)

var mealPlanOptionVotesColumns = []string{
	idColumn,
	"rank",
	"abstain",
	notesColumn,
	"by_user",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	belongsToMealPlanOptionColumn,
}

func buildMealPlanOptionVotesQueries() []*Query {
	insertColumns := filterForInsert(mealPlanOptionVotesColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveMealPlanOptionVote",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				mealPlanOptionVotesTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToMealPlanOptionColumn,
				belongsToMealPlanOptionColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateMealPlanOptionVote",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				mealPlanOptionVotesTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckMealPlanOptionVoteExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
		JOIN %s ON %s.%s=%s.%s
		JOIN %s ON %s.%s=%s.%s
		JOIN %s ON %s.%s=%s.%s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				mealPlanOptionVotesTableName, idColumn,
				mealPlanOptionVotesTableName,
				mealPlanOptionsTableName, mealPlanOptionVotesTableName, belongsToMealPlanOptionColumn, mealPlanOptionsTableName, idColumn,
				mealPlanEventsTableName, mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventsTableName, idColumn,
				mealPlansTableName, mealPlanEventsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
				mealPlanOptionVotesTableName, archivedAtColumn,
				mealPlanOptionVotesTableName, belongsToMealPlanOptionColumn, mealPlanOptionIDColumn,
				mealPlanOptionVotesTableName, idColumn, mealPlanOptionVoteIDColumn,
				mealPlanOptionsTableName, archivedAtColumn,
				mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventIDColumn,
				mealPlanEventsTableName, archivedAtColumn,
				mealPlanEventsTableName, belongsToMealPlanColumn, mealPlanIDColumn,
				mealPlanOptionsTableName, idColumn, mealPlanOptionIDColumn,
				mealPlansTableName, archivedAtColumn,
				mealPlansTableName, idColumn, mealPlanIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetMealPlanOptionVotesForMealPlanOption",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(mealPlanOptionVotesColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", mealPlanOptionVotesTableName, s)
				}), ",\n\t"),
				mealPlanOptionVotesTableName,
				mealPlanOptionsTableName, mealPlanOptionVotesTableName, belongsToMealPlanOptionColumn, mealPlanOptionsTableName, idColumn,
				mealPlanEventsTableName, mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventsTableName, idColumn,
				mealPlansTableName, mealPlanEventsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
				mealPlanOptionVotesTableName, archivedAtColumn,
				mealPlanOptionVotesTableName, belongsToMealPlanOptionColumn, mealPlanOptionIDColumn,
				mealPlanOptionsTableName, archivedAtColumn,
				mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventIDColumn,
				mealPlanOptionsTableName, idColumn, mealPlanOptionIDColumn,
				mealPlanEventsTableName, archivedAtColumn,
				mealPlanEventsTableName, belongsToMealPlanColumn, mealPlanIDColumn,
				mealPlanEventsTableName, idColumn, mealPlanEventIDColumn,
				mealPlansTableName, archivedAtColumn,
				mealPlansTableName, idColumn, mealPlanIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetMealPlanOptionVotes",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT 
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
GROUP BY
	%s.%s,
	%s.%s,
	%s.%s,
	%s.%s
%s;`,
				strings.Join(applyToEach(mealPlanOptionVotesColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", mealPlanOptionVotesTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(mealPlanOptionVotesTableName, true, true, "meal_plan_option_votes.belongs_to_meal_plan_option = sqlc.arg(meal_plan_option_id)"),
				buildTotalCountSelect(mealPlanOptionVotesTableName, true),
				mealPlanOptionVotesTableName,
				mealPlanOptionsTableName, mealPlanOptionVotesTableName, belongsToMealPlanOptionColumn, mealPlanOptionsTableName, idColumn,
				mealPlanEventsTableName, mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventsTableName, idColumn,
				mealPlansTableName, mealPlanEventsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
				mealPlanOptionVotesTableName, archivedAtColumn,
				mealPlanOptionVotesTableName, belongsToMealPlanOptionColumn, mealPlanOptionIDColumn,
				mealPlanOptionsTableName, archivedAtColumn,
				mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventIDColumn,
				mealPlanOptionsTableName, idColumn, mealPlanOptionIDColumn,
				mealPlanEventsTableName, archivedAtColumn,
				mealPlanEventsTableName, belongsToMealPlanColumn, mealPlanIDColumn,
				mealPlanEventsTableName, idColumn, mealPlanEventIDColumn,
				mealPlansTableName, archivedAtColumn,
				mealPlansTableName, idColumn, mealPlanIDColumn,
				buildFilterConditions(mealPlanOptionVotesTableName, true),
				mealPlanOptionVotesTableName, idColumn,
				mealPlanOptionsTableName, idColumn,
				mealPlanEventsTableName, idColumn,
				mealPlansTableName, idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetMealPlanOptionVote",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(mealPlanOptionVotesColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", mealPlanOptionVotesTableName, s)
				}), ",\n\t"),
				mealPlanOptionVotesTableName,
				mealPlanOptionsTableName, mealPlanOptionVotesTableName, belongsToMealPlanOptionColumn, mealPlanOptionsTableName, idColumn,
				mealPlanEventsTableName, mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventsTableName, idColumn,
				mealPlansTableName, mealPlanEventsTableName, belongsToMealPlanColumn, mealPlansTableName, idColumn,
				mealPlanOptionVotesTableName, archivedAtColumn,
				mealPlanOptionVotesTableName, belongsToMealPlanOptionColumn, mealPlanOptionIDColumn,
				mealPlanOptionVotesTableName, idColumn, mealPlanOptionVoteIDColumn,
				mealPlanOptionsTableName, archivedAtColumn,
				mealPlanOptionsTableName, belongsToMealPlanEventColumn, mealPlanEventIDColumn,
				mealPlanEventsTableName, archivedAtColumn,
				mealPlanEventsTableName, belongsToMealPlanColumn, mealPlanIDColumn,
				mealPlanOptionsTableName, idColumn, mealPlanOptionIDColumn,
				mealPlansTableName, archivedAtColumn,
				mealPlansTableName, idColumn, mealPlanIDColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateMealPlanOptionVote",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				mealPlanOptionVotesTableName,
				strings.Join(applyToEach(filterForUpdate(mealPlanOptionVotesColumns, belongsToMealPlanOptionColumn), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				belongsToMealPlanOptionColumn, belongsToMealPlanOptionColumn,
				idColumn,
				idColumn,
			)),
		},
	}
}
