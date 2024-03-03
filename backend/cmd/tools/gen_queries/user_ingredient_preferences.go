package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	userIngredientPreferencesTableName = "user_ingredient_preferences"

	userIngredientPreferencesIngredientColumn = "ingredient"
)

var userIngredientPreferencesColumns = []string{
	idColumn,
	userIngredientPreferencesIngredientColumn,
	"rating",
	notesColumn,
	"allergy",
	belongsToUserColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildUserIngredientPreferencesQueries() []*Query {
	insertColumns := filterForInsert(userIngredientPreferencesColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(filterFromSlice(userIngredientPreferencesColumns, userIngredientPreferencesIngredientColumn), func(i int, s string) string {
			return fmt.Sprintf("%s.%s", userIngredientPreferencesTableName, s)
		}),
		applyToEach(validIngredientsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
		}),
		1,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveUserIngredientPreference",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
				userIngredientPreferencesTableName,
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
				Name: "CreateUserIngredientPreference",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				userIngredientPreferencesTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckUserIngredientPreferenceExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
);`,
				userIngredientPreferencesTableName, idColumn,
				userIngredientPreferencesTableName,
				userIngredientPreferencesTableName, archivedAtColumn,
				userIngredientPreferencesTableName, idColumn, idColumn,
				userIngredientPreferencesTableName, belongsToUserColumn, belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserIngredientPreferencesForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(userIngredientPreferencesTableName, true, true),
				buildTotalCountSelect(userIngredientPreferencesTableName, true),
				userIngredientPreferencesTableName,
				validIngredientsTableName, validIngredientsTableName, idColumn, userIngredientPreferencesTableName, userIngredientPreferencesIngredientColumn,
				userIngredientPreferencesTableName, archivedAtColumn,
				userIngredientPreferencesTableName, belongsToUserColumn, belongsToUserColumn,
				validIngredientsTableName, archivedAtColumn,
				buildFilterConditions(userIngredientPreferencesTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserIngredientPreference",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				userIngredientPreferencesTableName,
				validIngredientsTableName, validIngredientsTableName, idColumn, userIngredientPreferencesTableName, userIngredientPreferencesIngredientColumn,
				userIngredientPreferencesTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				userIngredientPreferencesTableName, idColumn, idColumn,
				userIngredientPreferencesTableName, belongsToUserColumn, belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserIngredientPreference",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				userIngredientPreferencesTableName,
				strings.Join(applyToEach(filterForUpdate(userIngredientPreferencesColumns, belongsToUserColumn), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn, currentTimeExpression,
				archivedAtColumn,
				belongsToUserColumn, belongsToUserColumn,
				idColumn, idColumn,
			)),
		},
	}
}
