package main

import (
	"github.com/cristalhq/builq"
)

const validIngredientsTableName = "valid_ingredients"

var validIngredientsColumns = []string{
	idColumn,
	"name",
	"description",
	"warning",
	"contains_egg",
	"contains_dairy",
	"contains_peanut",
	"contains_tree_nut",
	"contains_soy",
	"contains_wheat",
	"contains_shellfish",
	"contains_sesame",
	"contains_fish",
	"contains_gluten",
	"animal_flesh",
	"volumetric",
	"icon_path",
	"is_liquid",
	"animal_derived",
	"plural_name",
	"restrict_to_preparations",
	"minimum_ideal_storage_temperature_in_celsius",
	"maximum_ideal_storage_temperature_in_celsius",
	"storage_instructions",
	"contains_alcohol",
	"slug",
	"shopping_suggestions",
	"is_starch",
	"is_protein",
	"is_grain",
	"is_fruit",
	"is_salt",
	"is_fat",
	"is_acid",
	"is_heat",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientsQueries() []*Query {
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