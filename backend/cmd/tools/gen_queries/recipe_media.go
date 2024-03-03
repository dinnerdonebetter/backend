package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeMediaTableName = "recipe_media"
)

var recipeMediaColumns = []string{
	idColumn,
	belongsToRecipeColumn,
	belongsToRecipeStepColumn,
	"mime_type",
	"internal_path",
	"external_path",
	"index",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeMediaQueries() []*Query {
	insertColumns := filterForInsert(recipeMediaColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveRecipeMedia",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				recipeMediaTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateRecipeMedia",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipeMediaTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckRecipeMediaExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				recipeMediaTableName, idColumn,
				recipeMediaTableName,
				recipeMediaTableName, archivedAtColumn,
				recipeMediaTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeMediaForRecipe",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s IS NULL
GROUP BY %s.%s
ORDER BY %s.%s;`,
				strings.Join(applyToEach(recipeMediaColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", recipeMediaTableName, s)
				}), ",\n\t"),
				recipeMediaTableName,
				recipeMediaTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeMediaTableName, belongsToRecipeStepColumn,
				recipeMediaTableName, archivedAtColumn,
				recipeMediaTableName, idColumn,
				recipeMediaTableName, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeMediaForRecipeStep",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
GROUP BY %s.%s
ORDER BY %s.%s;`,
				strings.Join(applyToEach(recipeMediaColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", recipeMediaTableName, s)
				}), ",\n\t"),
				recipeMediaTableName,
				recipeMediaTableName, belongsToRecipeColumn, recipeIDColumn,
				recipeMediaTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
				recipeMediaTableName, archivedAtColumn,
				recipeMediaTableName, idColumn,
				recipeMediaTableName, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRecipeMedia",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(recipeMediaColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", recipeMediaTableName, s)
				}), ",\n\t"),
				recipeMediaTableName,
				recipeMediaTableName, archivedAtColumn,
				recipeMediaTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateRecipeMedia",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				recipeMediaTableName,
				strings.Join(applyToEach(filterForUpdate(recipeMediaColumns), func(i int, s string) string {
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
