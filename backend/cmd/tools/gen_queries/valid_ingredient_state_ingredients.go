package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validIngredientStateIngredientsTableName = "valid_ingredient_state_ingredients"

	validIngredientStateColumn = "valid_ingredient_state"
	validIngredientColumn      = "valid_ingredient"
)

var validIngredientStateIngredientsColumns = []string{
	idColumn,
	notesColumn,
	validIngredientStateColumn,
	validIngredientColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientStateIngredientsQueries() []*Query {
	insertColumns := filterForInsert(validIngredientStateIngredientsColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(filterFromSlice(validIngredientStateIngredientsColumns, "valid_ingredient_id", "valid_measurement_unit_id"), func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_ingredient_state_ingredient_%s", validIngredientStateIngredientsTableName, s, s)
		}),
		append(
			applyToEach(validIngredientStatesColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_state_%s", validIngredientStatesTableName, s, s)
			}),
			applyToEach(validIngredientsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
			})...),
		2,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidIngredientStateIngredient",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validIngredientStateIngredientsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidIngredientStateIngredient",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				validIngredientStateIngredientsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidIngredientStateIngredientExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				validIngredientStateIngredientsTableName, idColumn,
				validIngredientStateIngredientsTableName,
				validIngredientStateIngredientsTableName, archivedAtColumn,
				validIngredientStateIngredientsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredientsForIngredient",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(validIngredientStateIngredientsTableName, true, true),
				buildTotalCountSelect(validIngredientStateIngredientsTableName, true),
				validIngredientStateIngredientsTableName,
				validIngredientsTableName, validIngredientStateIngredientsTableName, validIngredientColumn, validIngredientsTableName, idColumn,
				validIngredientStatesTableName, validIngredientStateIngredientsTableName, validIngredientStateColumn, validIngredientStatesTableName, idColumn,
				validIngredientStateIngredientsTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				validIngredientStatesTableName, archivedAtColumn,
				validIngredientStateIngredientsTableName, validIngredientColumn, validIngredientColumn,
				buildFilterConditions(validIngredientStateIngredientsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredientsForIngredientState",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(validIngredientStateIngredientsTableName, true, true),
				buildTotalCountSelect(validIngredientStateIngredientsTableName, true),
				validIngredientStateIngredientsTableName,
				validIngredientsTableName, validIngredientStateIngredientsTableName, validIngredientColumn, validIngredientsTableName, idColumn,
				validIngredientStatesTableName, validIngredientStateIngredientsTableName, validIngredientStateColumn, validIngredientStatesTableName, idColumn,
				validIngredientStateIngredientsTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				validIngredientStatesTableName, archivedAtColumn,
				validIngredientStateIngredientsTableName, validIngredientStateColumn, validIngredientStateColumn,
				buildFilterConditions(validIngredientStateIngredientsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredients",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(validIngredientStateIngredientsTableName, true, true),
				buildTotalCountSelect(validIngredientStateIngredientsTableName, true),
				validIngredientStateIngredientsTableName,
				validIngredientsTableName, validIngredientStateIngredientsTableName, validIngredientColumn, validIngredientsTableName, idColumn,
				validIngredientStatesTableName, validIngredientStateIngredientsTableName, validIngredientStateColumn, validIngredientStatesTableName, idColumn,
				validIngredientStateIngredientsTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				validIngredientStatesTableName, archivedAtColumn,
				buildFilterConditions(validIngredientStateIngredientsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredient",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validIngredientStateIngredientsTableName,
				validIngredientsTableName, validIngredientStateIngredientsTableName, validIngredientColumn, validIngredientsTableName, idColumn,
				validIngredientStatesTableName, validIngredientStateIngredientsTableName, validIngredientStateColumn, validIngredientStatesTableName, idColumn,
				validIngredientStateIngredientsTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				validIngredientStatesTableName, archivedAtColumn,
				validIngredientStateIngredientsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateIngredientsWithIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[]);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validIngredientStateIngredientsTableName,
				validIngredientsTableName, validIngredientStateIngredientsTableName, validIngredientColumn, validIngredientsTableName, idColumn,
				validIngredientStatesTableName, validIngredientStateIngredientsTableName, validIngredientStateColumn, validIngredientStatesTableName, idColumn,
				validIngredientStateIngredientsTableName, archivedAtColumn,
				validIngredientsTableName, archivedAtColumn,
				validIngredientStatesTableName, archivedAtColumn,
				validIngredientStateIngredientsTableName, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidityOfValidIngredientStateIngredientPair",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s.%s
	FROM %s
	WHERE %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s)
	AND %s IS NULL
);`,
				validIngredientStateIngredientsTableName, idColumn,
				validIngredientStateIngredientsTableName,
				validIngredientColumn, validIngredientColumn,
				validIngredientStateColumn, validIngredientStateColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidIngredientStateIngredient",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validIngredientStateIngredientsTableName,
				strings.Join(applyToEach(filterForUpdate(validIngredientStateIngredientsColumns), func(i int, s string) string {
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
