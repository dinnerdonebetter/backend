package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	mealComponentsTableName = "meal_components"
)

var mealComponentsColumns = []string{
	idColumn,
	"meal_id",
	"recipe_id",
	"meal_component_type",
	"recipe_scale",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealComponentsQueries() []*Query {
	insertColumns := filterForInsert(mealComponentsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "CreateMealComponent",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				mealComponentsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
	}
}
