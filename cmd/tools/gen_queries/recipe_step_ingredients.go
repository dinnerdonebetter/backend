package main

import (
	"github.com/cristalhq/builq"
)

const recipeStepIngredientsTableName = "recipe_step_ingredients"

var recipeStepIngredientsColumns = []string{
	idColumn,
	"ingredient_id",
	"minimum_quantity_value",
	"quantity_notes",
	"ingredient_notes",
	"belongs_to_recipe_step",
	"name",
	"recipe_step_product_id",
	"maximum_quantity_value",
	"measurement_unit",
	"optional",
	"option_index",
	"vessel_index",
	"to_taste",
	"product_percentage_to_use",
	"recipe_step_product_recipe_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeStepIngredientsQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}
