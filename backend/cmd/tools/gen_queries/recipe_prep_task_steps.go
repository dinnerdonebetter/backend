package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipePrepTaskStepsTableName = "recipe_prep_task_steps"

	satisfiesRecipeStepColumn = "satisfies_recipe_step"
)

var recipePrepTaskStepsColumns = []string{
	idColumn,
	belongsToRecipeStepColumn,
	"belongs_to_recipe_prep_task",
	satisfiesRecipeStepColumn,
}

func buildRecipePrepTaskStepsQueries() []*Query {
	insertColumns := filterForInsert(recipePrepTaskStepsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "CreateRecipePrepTaskStep",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				recipePrepTaskStepsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
	}
}
