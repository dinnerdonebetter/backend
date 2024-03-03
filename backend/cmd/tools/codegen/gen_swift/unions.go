package main

import (
	"fmt"
	"strings"
)

type unionDefinition struct {
	Name,
	ConstName,
	Comment string
	Values []string
}

var (
	unions = []*unionDefinition{
		{
			Comment:   "ingredient state attribute types",
			Name:      "ValidIngredientStateAttributeType",
			ConstName: "VALID_INGREDIENT_STATE_ATTRIBUTE_TYPES",
			Values: []string{
				"texture",
				"consistency",
				"temperature",
				"color",
				"appearance",
				"odor",
				"taste",
				"sound",
				"other",
			},
		},
		{
			Comment:   "meal plan statuses",
			Name:      "ValidMealPlanStatus",
			ConstName: "VALID_MEAL_PLAN_STATUSES",
			Values: []string{
				"awaiting_votes",
				"finalized",
			},
		},
		{
			Comment:   "meal plan grocery list item statuses",
			Name:      "ValidMealPlanGroceryListItemStatus",
			ConstName: "VALID_MEAL_PLAN_GROCERY_LIST_ITEM_STATUSES",
			Values: []string{
				"unknown",
				"already owned",
				"needs",
				"unavailable",
				"acquired",
			},
		},
		{
			Comment:   "meal plan election methods",
			Name:      "ValidMealPlanElectionMethod",
			ConstName: "VALID_MEAL_PLAN_ELECTION_METHODS",
			Values: []string{
				"schulze",
				"instant-runoff",
			},
		},
		{
			Comment:   "recipe step product types",
			Name:      "ValidRecipeStepProductType",
			ConstName: "RECIPE_STEP_PRODUCT_TYPES",
			Values: []string{
				"ingredient",
				"instrument",
				"vessel",
			},
		},
		{
			Comment:   "meal component types",
			Name:      "MealComponentType",
			ConstName: "MEAL_COMPONENT_TYPES",
			Values: []string{
				"main",
				"side",
				"appetizer",
				"beverage",
				"dessert",
				"soup",
				"salad",
				"amuse-bouche",
				"unspecified",
			},
		},
		{
			Comment:   "meal plan task statuses",
			Name:      "MealPlanTaskStatus",
			ConstName: "MEAL_PLAN_TASK_STATUSES",
			Values: []string{
				"unfinished",
				"postponed",
				"ignored",
				"canceled",
				"finished",
			},
		},
	}
)

func buildUnionDeclaration(def *unionDefinition) string {
	if def == nil {
		return ""
	}

	output := ""

	output += fmt.Sprintf("/**\n * %s\n */\n", def.Comment)
	output += fmt.Sprintf("enum %s: String{\n", def.Name)
	for _, value := range def.Values {
		output += fmt.Sprintf("case %s = %q\n", strings.ReplaceAll(value, " ", "_"), value)
	}
	output += "}\n"

	return output
}

func buildUnionsFile() string {
	output := copyString(generatedDisclaimer)

	for _, def := range unions {
		output += buildUnionDeclaration(def)
	}

	return output
}
