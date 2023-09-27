// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: meal_plan_grocery_list_items.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveMealPlanGroceryListItem = `-- name: ArchiveMealPlanGroceryListItem :execrows

UPDATE meal_plan_grocery_list_items SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1
`

func (q *Queries) ArchiveMealPlanGroceryListItem(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, archiveMealPlanGroceryListItem, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkMealPlanGroceryListItemExistence = `-- name: CheckMealPlanGroceryListItemExistence :one

SELECT EXISTS ( SELECT meal_plan_grocery_list_items.id FROM meal_plan_grocery_list_items WHERE meal_plan_grocery_list_items.archived_at IS NULL AND meal_plan_grocery_list_items.id = $1 AND meal_plan_grocery_list_items.belongs_to_meal_plan = $2 )
`

type CheckMealPlanGroceryListItemExistenceParams struct {
	MealPlanGroceryListItemID string
	MealPlanID                string
}

func (q *Queries) CheckMealPlanGroceryListItemExistence(ctx context.Context, db DBTX, arg *CheckMealPlanGroceryListItemExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanGroceryListItemExistence, arg.MealPlanGroceryListItemID, arg.MealPlanID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createMealPlanGroceryListItem = `-- name: CreateMealPlanGroceryListItem :exec

INSERT INTO meal_plan_grocery_list_items (
    id,
    belongs_to_meal_plan,
    valid_ingredient,
    valid_measurement_unit,
    minimum_quantity_needed,
    maximum_quantity_needed,
    quantity_purchased,
    purchased_measurement_unit,
    purchased_upc,
    purchase_price,
    status_explanation,
    status
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12
)
`

type CreateMealPlanGroceryListItemParams struct {
	ID                       string
	BelongsToMealPlan        string
	ValidIngredient          string
	ValidMeasurementUnit     string
	MinimumQuantityNeeded    string
	StatusExplanation        string
	Status                   GroceryListItemStatus
	MaximumQuantityNeeded    sql.NullString
	QuantityPurchased        sql.NullString
	PurchasedMeasurementUnit sql.NullString
	PurchasedUpc             sql.NullString
	PurchasePrice            sql.NullString
}

func (q *Queries) CreateMealPlanGroceryListItem(ctx context.Context, db DBTX, arg *CreateMealPlanGroceryListItemParams) error {
	_, err := db.ExecContext(ctx, createMealPlanGroceryListItem,
		arg.ID,
		arg.BelongsToMealPlan,
		arg.ValidIngredient,
		arg.ValidMeasurementUnit,
		arg.MinimumQuantityNeeded,
		arg.MaximumQuantityNeeded,
		arg.QuantityPurchased,
		arg.PurchasedMeasurementUnit,
		arg.PurchasedUpc,
		arg.PurchasePrice,
		arg.StatusExplanation,
		arg.Status,
	)
	return err
}

const getMealPlanGroceryListItem = `-- name: GetMealPlanGroceryListItem :one

SELECT
	meal_plan_grocery_list_items.id,
	meal_plan_grocery_list_items.belongs_to_meal_plan,
    valid_ingredients.id as valid_ingredient_id,
    valid_ingredients.name as valid_ingredient_name,
    valid_ingredients.description as valid_ingredient_description,
    valid_ingredients.warning as valid_ingredient_warning,
    valid_ingredients.contains_egg as valid_ingredient_contains_egg,
    valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
    valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
    valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
    valid_ingredients.contains_soy as valid_ingredient_contains_soy,
    valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
    valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
    valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
    valid_ingredients.contains_fish as valid_ingredient_contains_fish,
    valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
    valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
    valid_ingredients.volumetric as valid_ingredient_volumetric,
    valid_ingredients.is_liquid as valid_ingredient_is_liquid,
    valid_ingredients.icon_path as valid_ingredient_icon_path,
    valid_ingredients.animal_derived as valid_ingredient_animal_derived,
    valid_ingredients.plural_name as valid_ingredient_plural_name,
    valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
    valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
    valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
    valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
    valid_ingredients.slug as valid_ingredient_slug,
    valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
    valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
    valid_ingredients.is_starch as valid_ingredient_is_starch,
    valid_ingredients.is_protein as valid_ingredient_is_protein,
    valid_ingredients.is_grain as valid_ingredient_is_grain,
    valid_ingredients.is_fruit as valid_ingredient_is_fruit,
    valid_ingredients.is_salt as valid_ingredient_is_salt,
    valid_ingredients.is_fat as valid_ingredient_is_fat,
    valid_ingredients.is_acid as valid_ingredient_is_acid,
    valid_ingredients.is_heat as valid_ingredient_is_heat,
    valid_ingredients.created_at as valid_ingredient_created_at,
    valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
    valid_ingredients.archived_at as valid_ingredient_archived_at,
    valid_measurement_units.id as valid_measurement_unit_id,
    valid_measurement_units.name as valid_measurement_unit_name,
    valid_measurement_units.description as valid_measurement_unit_description,
    valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
    valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
    valid_measurement_units.universal as valid_measurement_unit_universal,
    valid_measurement_units.metric as valid_measurement_unit_metric,
    valid_measurement_units.imperial as valid_measurement_unit_imperial,
    valid_measurement_units.slug as valid_measurement_unit_slug,
    valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
    valid_measurement_units.created_at as valid_measurement_unit_created_at,
    valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
    valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
	meal_plan_grocery_list_items.minimum_quantity_needed,
	meal_plan_grocery_list_items.maximum_quantity_needed,
	meal_plan_grocery_list_items.quantity_purchased,
	meal_plan_grocery_list_items.purchased_measurement_unit,
	meal_plan_grocery_list_items.purchased_upc,
	meal_plan_grocery_list_items.purchase_price,
	meal_plan_grocery_list_items.status_explanation,
	meal_plan_grocery_list_items.status,
	meal_plan_grocery_list_items.created_at,
	meal_plan_grocery_list_items.last_updated_at,
	meal_plan_grocery_list_items.archived_at
FROM meal_plan_grocery_list_items
	JOIN meal_plans ON meal_plan_grocery_list_items.belongs_to_meal_plan=meal_plans.id
    JOIN valid_ingredients ON meal_plan_grocery_list_items.valid_ingredient=valid_ingredients.id
    JOIN valid_measurement_units ON meal_plan_grocery_list_items.valid_measurement_unit=valid_measurement_units.id
WHERE meal_plan_grocery_list_items.archived_at IS NULL
    AND valid_measurement_units.archived_at IS NULL
    AND valid_ingredients.archived_at IS NULL
    AND meal_plan_grocery_list_items.id = $1
    AND meal_plan_grocery_list_items.belongs_to_meal_plan = $2
`

type GetMealPlanGroceryListItemParams struct {
	MealPlanGroceryListItemID string
	MealPlanID                string
}

type GetMealPlanGroceryListItemRow struct {
	ValidIngredientCreatedAt                               time.Time
	ValidMeasurementUnitCreatedAt                          time.Time
	CreatedAt                                              time.Time
	ValidIngredientLastUpdatedAt                           sql.NullTime
	ValidIngredientArchivedAt                              sql.NullTime
	ValidMeasurementUnitLastUpdatedAt                      sql.NullTime
	ValidMeasurementUnitArchivedAt                         sql.NullTime
	ArchivedAt                                             sql.NullTime
	LastUpdatedAt                                          sql.NullTime
	Status                                                 GroceryListItemStatus
	ValidMeasurementUnitIconPath                           string
	BelongsToMealPlan                                      string
	ValidIngredientID                                      string
	ValidMeasurementUnitID                                 string
	ValidMeasurementUnitName                               string
	ValidMeasurementUnitDescription                        string
	ID                                                     string
	MinimumQuantityNeeded                                  string
	StatusExplanation                                      string
	ValidIngredientIconPath                                string
	ValidIngredientWarning                                 string
	ValidIngredientPluralName                              string
	ValidIngredientDescription                             string
	ValidMeasurementUnitSlug                               string
	ValidMeasurementUnitPluralName                         string
	ValidIngredientStorageInstructions                     string
	ValidIngredientSlug                                    string
	ValidIngredientName                                    string
	ValidIngredientShoppingSuggestions                     string
	ValidIngredientMaximumIdealStorageTemperatureInCelsius sql.NullString
	ValidIngredientMinimumIdealStorageTemperatureInCelsius sql.NullString
	MaximumQuantityNeeded                                  sql.NullString
	QuantityPurchased                                      sql.NullString
	PurchasedMeasurementUnit                               sql.NullString
	PurchasedUpc                                           sql.NullString
	PurchasePrice                                          sql.NullString
	ValidIngredientIsLiquid                                sql.NullBool
	ValidMeasurementUnitVolumetric                         sql.NullBool
	ValidIngredientAnimalFlesh                             bool
	ValidIngredientIsHeat                                  bool
	ValidIngredientIsAcid                                  bool
	ValidIngredientIsFat                                   bool
	ValidIngredientIsSalt                                  bool
	ValidIngredientIsFruit                                 bool
	ValidIngredientIsGrain                                 bool
	ValidMeasurementUnitUniversal                          bool
	ValidMeasurementUnitMetric                             bool
	ValidMeasurementUnitImperial                           bool
	ValidIngredientIsProtein                               bool
	ValidIngredientIsStarch                                bool
	ValidIngredientContainsAlcohol                         bool
	ValidIngredientRestrictToPreparations                  bool
	ValidIngredientAnimalDerived                           bool
	ValidIngredientVolumetric                              bool
	ValidIngredientContainsGluten                          bool
	ValidIngredientContainsFish                            bool
	ValidIngredientContainsSesame                          bool
	ValidIngredientContainsShellfish                       bool
	ValidIngredientContainsWheat                           bool
	ValidIngredientContainsSoy                             bool
	ValidIngredientContainsTreeNut                         bool
	ValidIngredientContainsPeanut                          bool
	ValidIngredientContainsDairy                           bool
	ValidIngredientContainsEgg                             bool
}

func (q *Queries) GetMealPlanGroceryListItem(ctx context.Context, db DBTX, arg *GetMealPlanGroceryListItemParams) (*GetMealPlanGroceryListItemRow, error) {
	row := db.QueryRowContext(ctx, getMealPlanGroceryListItem, arg.MealPlanGroceryListItemID, arg.MealPlanID)
	var i GetMealPlanGroceryListItemRow
	err := row.Scan(
		&i.ID,
		&i.BelongsToMealPlan,
		&i.ValidIngredientID,
		&i.ValidIngredientName,
		&i.ValidIngredientDescription,
		&i.ValidIngredientWarning,
		&i.ValidIngredientContainsEgg,
		&i.ValidIngredientContainsDairy,
		&i.ValidIngredientContainsPeanut,
		&i.ValidIngredientContainsTreeNut,
		&i.ValidIngredientContainsSoy,
		&i.ValidIngredientContainsWheat,
		&i.ValidIngredientContainsShellfish,
		&i.ValidIngredientContainsSesame,
		&i.ValidIngredientContainsFish,
		&i.ValidIngredientContainsGluten,
		&i.ValidIngredientAnimalFlesh,
		&i.ValidIngredientVolumetric,
		&i.ValidIngredientIsLiquid,
		&i.ValidIngredientIconPath,
		&i.ValidIngredientAnimalDerived,
		&i.ValidIngredientPluralName,
		&i.ValidIngredientRestrictToPreparations,
		&i.ValidIngredientMinimumIdealStorageTemperatureInCelsius,
		&i.ValidIngredientMaximumIdealStorageTemperatureInCelsius,
		&i.ValidIngredientStorageInstructions,
		&i.ValidIngredientSlug,
		&i.ValidIngredientContainsAlcohol,
		&i.ValidIngredientShoppingSuggestions,
		&i.ValidIngredientIsStarch,
		&i.ValidIngredientIsProtein,
		&i.ValidIngredientIsGrain,
		&i.ValidIngredientIsFruit,
		&i.ValidIngredientIsSalt,
		&i.ValidIngredientIsFat,
		&i.ValidIngredientIsAcid,
		&i.ValidIngredientIsHeat,
		&i.ValidIngredientCreatedAt,
		&i.ValidIngredientLastUpdatedAt,
		&i.ValidIngredientArchivedAt,
		&i.ValidMeasurementUnitID,
		&i.ValidMeasurementUnitName,
		&i.ValidMeasurementUnitDescription,
		&i.ValidMeasurementUnitVolumetric,
		&i.ValidMeasurementUnitIconPath,
		&i.ValidMeasurementUnitUniversal,
		&i.ValidMeasurementUnitMetric,
		&i.ValidMeasurementUnitImperial,
		&i.ValidMeasurementUnitSlug,
		&i.ValidMeasurementUnitPluralName,
		&i.ValidMeasurementUnitCreatedAt,
		&i.ValidMeasurementUnitLastUpdatedAt,
		&i.ValidMeasurementUnitArchivedAt,
		&i.MinimumQuantityNeeded,
		&i.MaximumQuantityNeeded,
		&i.QuantityPurchased,
		&i.PurchasedMeasurementUnit,
		&i.PurchasedUpc,
		&i.PurchasePrice,
		&i.StatusExplanation,
		&i.Status,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getMealPlanGroceryListItemsForMealPlan = `-- name: GetMealPlanGroceryListItemsForMealPlan :many

SELECT
    meal_plan_grocery_list_items.id,
    meal_plan_grocery_list_items.belongs_to_meal_plan,
    valid_ingredients.id as valid_ingredient_id,
    valid_ingredients.name as valid_ingredient_name,
    valid_ingredients.description as valid_ingredient_description,
    valid_ingredients.warning as valid_ingredient_warning,
    valid_ingredients.contains_egg as valid_ingredient_contains_egg,
    valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
    valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
    valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
    valid_ingredients.contains_soy as valid_ingredient_contains_soy,
    valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
    valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
    valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
    valid_ingredients.contains_fish as valid_ingredient_contains_fish,
    valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
    valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
    valid_ingredients.volumetric as valid_ingredient_volumetric,
    valid_ingredients.is_liquid as valid_ingredient_is_liquid,
    valid_ingredients.icon_path as valid_ingredient_icon_path,
    valid_ingredients.animal_derived as valid_ingredient_animal_derived,
    valid_ingredients.plural_name as valid_ingredient_plural_name,
    valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
    valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
    valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
    valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
    valid_ingredients.slug as valid_ingredient_slug,
    valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
    valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
    valid_ingredients.is_starch as valid_ingredient_is_starch,
    valid_ingredients.is_protein as valid_ingredient_is_protein,
    valid_ingredients.is_grain as valid_ingredient_is_grain,
    valid_ingredients.is_fruit as valid_ingredient_is_fruit,
    valid_ingredients.is_salt as valid_ingredient_is_salt,
    valid_ingredients.is_fat as valid_ingredient_is_fat,
    valid_ingredients.is_acid as valid_ingredient_is_acid,
    valid_ingredients.is_heat as valid_ingredient_is_heat,
    valid_ingredients.created_at as valid_ingredient_created_at,
    valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
    valid_ingredients.archived_at as valid_ingredient_archived_at,
    valid_measurement_units.id as valid_measurement_unit_id,
    valid_measurement_units.name as valid_measurement_unit_name,
    valid_measurement_units.description as valid_measurement_unit_description,
    valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
    valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
    valid_measurement_units.universal as valid_measurement_unit_universal,
    valid_measurement_units.metric as valid_measurement_unit_metric,
    valid_measurement_units.imperial as valid_measurement_unit_imperial,
    valid_measurement_units.slug as valid_measurement_unit_slug,
    valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
    valid_measurement_units.created_at as valid_measurement_unit_created_at,
    valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
    valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
    meal_plan_grocery_list_items.minimum_quantity_needed,
    meal_plan_grocery_list_items.maximum_quantity_needed,
    meal_plan_grocery_list_items.quantity_purchased,
    meal_plan_grocery_list_items.purchased_measurement_unit,
    meal_plan_grocery_list_items.purchased_upc,
    meal_plan_grocery_list_items.purchase_price,
    meal_plan_grocery_list_items.status_explanation,
    meal_plan_grocery_list_items.status,
    meal_plan_grocery_list_items.created_at,
    meal_plan_grocery_list_items.last_updated_at,
    meal_plan_grocery_list_items.archived_at
FROM meal_plan_grocery_list_items
    JOIN meal_plans ON meal_plan_grocery_list_items.belongs_to_meal_plan=meal_plans.id
    JOIN valid_ingredients ON meal_plan_grocery_list_items.valid_ingredient=valid_ingredients.id
    JOIN valid_measurement_units ON meal_plan_grocery_list_items.valid_measurement_unit=valid_measurement_units.id
WHERE meal_plan_grocery_list_items.archived_at IS NULL
    AND valid_measurement_units.archived_at IS NULL
    AND valid_ingredients.archived_at IS NULL
    AND meal_plan_grocery_list_items.belongs_to_meal_plan = $1
    AND meal_plans.archived_at IS NULL
    AND meal_plans.id = $1
GROUP BY meal_plan_grocery_list_items.id,
    valid_ingredients.id,
    valid_measurement_units.id,
    meal_plans.id
ORDER BY meal_plan_grocery_list_items.id
`

type GetMealPlanGroceryListItemsForMealPlanRow struct {
	ValidIngredientCreatedAt                               time.Time
	ValidMeasurementUnitCreatedAt                          time.Time
	CreatedAt                                              time.Time
	ValidIngredientLastUpdatedAt                           sql.NullTime
	ValidIngredientArchivedAt                              sql.NullTime
	ValidMeasurementUnitLastUpdatedAt                      sql.NullTime
	ValidMeasurementUnitArchivedAt                         sql.NullTime
	ArchivedAt                                             sql.NullTime
	LastUpdatedAt                                          sql.NullTime
	Status                                                 GroceryListItemStatus
	ValidMeasurementUnitIconPath                           string
	BelongsToMealPlan                                      string
	ValidIngredientID                                      string
	ValidMeasurementUnitID                                 string
	ValidMeasurementUnitName                               string
	ValidMeasurementUnitDescription                        string
	ID                                                     string
	MinimumQuantityNeeded                                  string
	StatusExplanation                                      string
	ValidIngredientIconPath                                string
	ValidIngredientWarning                                 string
	ValidIngredientPluralName                              string
	ValidIngredientDescription                             string
	ValidMeasurementUnitSlug                               string
	ValidMeasurementUnitPluralName                         string
	ValidIngredientStorageInstructions                     string
	ValidIngredientSlug                                    string
	ValidIngredientName                                    string
	ValidIngredientShoppingSuggestions                     string
	ValidIngredientMaximumIdealStorageTemperatureInCelsius sql.NullString
	ValidIngredientMinimumIdealStorageTemperatureInCelsius sql.NullString
	MaximumQuantityNeeded                                  sql.NullString
	QuantityPurchased                                      sql.NullString
	PurchasedMeasurementUnit                               sql.NullString
	PurchasedUpc                                           sql.NullString
	PurchasePrice                                          sql.NullString
	ValidIngredientIsLiquid                                sql.NullBool
	ValidMeasurementUnitVolumetric                         sql.NullBool
	ValidIngredientAnimalFlesh                             bool
	ValidIngredientIsHeat                                  bool
	ValidIngredientIsAcid                                  bool
	ValidIngredientIsFat                                   bool
	ValidIngredientIsSalt                                  bool
	ValidIngredientIsFruit                                 bool
	ValidIngredientIsGrain                                 bool
	ValidMeasurementUnitUniversal                          bool
	ValidMeasurementUnitMetric                             bool
	ValidMeasurementUnitImperial                           bool
	ValidIngredientIsProtein                               bool
	ValidIngredientIsStarch                                bool
	ValidIngredientContainsAlcohol                         bool
	ValidIngredientRestrictToPreparations                  bool
	ValidIngredientAnimalDerived                           bool
	ValidIngredientVolumetric                              bool
	ValidIngredientContainsGluten                          bool
	ValidIngredientContainsFish                            bool
	ValidIngredientContainsSesame                          bool
	ValidIngredientContainsShellfish                       bool
	ValidIngredientContainsWheat                           bool
	ValidIngredientContainsSoy                             bool
	ValidIngredientContainsTreeNut                         bool
	ValidIngredientContainsPeanut                          bool
	ValidIngredientContainsDairy                           bool
	ValidIngredientContainsEgg                             bool
}

func (q *Queries) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, db DBTX, mealPlanID string) ([]*GetMealPlanGroceryListItemsForMealPlanRow, error) {
	rows, err := db.QueryContext(ctx, getMealPlanGroceryListItemsForMealPlan, mealPlanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetMealPlanGroceryListItemsForMealPlanRow{}
	for rows.Next() {
		var i GetMealPlanGroceryListItemsForMealPlanRow
		if err := rows.Scan(
			&i.ID,
			&i.BelongsToMealPlan,
			&i.ValidIngredientID,
			&i.ValidIngredientName,
			&i.ValidIngredientDescription,
			&i.ValidIngredientWarning,
			&i.ValidIngredientContainsEgg,
			&i.ValidIngredientContainsDairy,
			&i.ValidIngredientContainsPeanut,
			&i.ValidIngredientContainsTreeNut,
			&i.ValidIngredientContainsSoy,
			&i.ValidIngredientContainsWheat,
			&i.ValidIngredientContainsShellfish,
			&i.ValidIngredientContainsSesame,
			&i.ValidIngredientContainsFish,
			&i.ValidIngredientContainsGluten,
			&i.ValidIngredientAnimalFlesh,
			&i.ValidIngredientVolumetric,
			&i.ValidIngredientIsLiquid,
			&i.ValidIngredientIconPath,
			&i.ValidIngredientAnimalDerived,
			&i.ValidIngredientPluralName,
			&i.ValidIngredientRestrictToPreparations,
			&i.ValidIngredientMinimumIdealStorageTemperatureInCelsius,
			&i.ValidIngredientMaximumIdealStorageTemperatureInCelsius,
			&i.ValidIngredientStorageInstructions,
			&i.ValidIngredientSlug,
			&i.ValidIngredientContainsAlcohol,
			&i.ValidIngredientShoppingSuggestions,
			&i.ValidIngredientIsStarch,
			&i.ValidIngredientIsProtein,
			&i.ValidIngredientIsGrain,
			&i.ValidIngredientIsFruit,
			&i.ValidIngredientIsSalt,
			&i.ValidIngredientIsFat,
			&i.ValidIngredientIsAcid,
			&i.ValidIngredientIsHeat,
			&i.ValidIngredientCreatedAt,
			&i.ValidIngredientLastUpdatedAt,
			&i.ValidIngredientArchivedAt,
			&i.ValidMeasurementUnitID,
			&i.ValidMeasurementUnitName,
			&i.ValidMeasurementUnitDescription,
			&i.ValidMeasurementUnitVolumetric,
			&i.ValidMeasurementUnitIconPath,
			&i.ValidMeasurementUnitUniversal,
			&i.ValidMeasurementUnitMetric,
			&i.ValidMeasurementUnitImperial,
			&i.ValidMeasurementUnitSlug,
			&i.ValidMeasurementUnitPluralName,
			&i.ValidMeasurementUnitCreatedAt,
			&i.ValidMeasurementUnitLastUpdatedAt,
			&i.ValidMeasurementUnitArchivedAt,
			&i.MinimumQuantityNeeded,
			&i.MaximumQuantityNeeded,
			&i.QuantityPurchased,
			&i.PurchasedMeasurementUnit,
			&i.PurchasedUpc,
			&i.PurchasePrice,
			&i.StatusExplanation,
			&i.Status,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMealPlanGroceryListItem = `-- name: UpdateMealPlanGroceryListItem :execrows

UPDATE meal_plan_grocery_list_items SET
	belongs_to_meal_plan = $1,
	valid_ingredient = $2,
	valid_measurement_unit = $3,
	minimum_quantity_needed = $4,
	maximum_quantity_needed = $5,
	quantity_purchased = $6,
	purchased_measurement_unit = $7,
	purchased_upc = $8,
	purchase_price = $9,
	status_explanation = $10,
	status = $11,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $12
`

type UpdateMealPlanGroceryListItemParams struct {
	BelongsToMealPlan        string
	ValidIngredient          string
	ValidMeasurementUnit     string
	MinimumQuantityNeeded    string
	StatusExplanation        string
	Status                   GroceryListItemStatus
	ID                       string
	MaximumQuantityNeeded    sql.NullString
	QuantityPurchased        sql.NullString
	PurchasedMeasurementUnit sql.NullString
	PurchasedUpc             sql.NullString
	PurchasePrice            sql.NullString
}

func (q *Queries) UpdateMealPlanGroceryListItem(ctx context.Context, db DBTX, arg *UpdateMealPlanGroceryListItemParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateMealPlanGroceryListItem,
		arg.BelongsToMealPlan,
		arg.ValidIngredient,
		arg.ValidMeasurementUnit,
		arg.MinimumQuantityNeeded,
		arg.MaximumQuantityNeeded,
		arg.QuantityPurchased,
		arg.PurchasedMeasurementUnit,
		arg.PurchasedUpc,
		arg.PurchasePrice,
		arg.StatusExplanation,
		arg.Status,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
