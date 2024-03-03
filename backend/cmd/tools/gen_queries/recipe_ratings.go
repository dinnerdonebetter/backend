package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeRatingsTableName = "recipe_ratings"
)

var recipeRatingsColumns = []string{
	idColumn,
	"recipe_id",
	"taste",
	"difficulty",
	"cleanup",
	"instructions",
	"overall",
	notesColumn,
	"by_user",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeRatingsQueries() []*Query {
	insertColumns := filterForInsert(recipeRatingsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveRecipeRating",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				recipeRatingsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateRecipeRating",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipeRatingsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckRecipeRatingExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				recipeRatingsTableName, idColumn,
				recipeRatingsTableName,
				recipeRatingsTableName, archivedAtColumn,
				recipeRatingsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeRatings",
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
				strings.Join(applyToEach(recipeRatingsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", recipeRatingsTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(
					recipeRatingsTableName,
					true,
					true,
				),
				buildTotalCountSelect(recipeRatingsTableName, true),
				recipeRatingsTableName,
				recipeRatingsTableName, archivedAtColumn,
				buildFilterConditions(
					recipeRatingsTableName,
					true,
				),
				recipeRatingsTableName, idColumn,
				recipeRatingsTableName, idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeRating",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(recipeRatingsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", recipeRatingsTableName, s)
				}), ",\n\t"),
				recipeRatingsTableName,
				recipeRatingsTableName, archivedAtColumn,
				recipeRatingsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipeRating",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				recipeRatingsTableName,
				strings.Join(applyToEach(filterForUpdate(recipeRatingsColumns, "by_user"), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				idColumn, idColumn,
			)),
		},
	}
}
