// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meal_plan_grocery_list_items_create.sql

package generated

import (
	"context"
	"database/sql"
)

const CreateMealPlanGroceryListItem = `-- name: CreateMealPlanGroceryListItem :exec
INSERT INTO meal_plan_grocery_list_items
(id,belongs_to_meal_plan,valid_ingredient,valid_measurement_unit,minimum_quantity_needed,maximum_quantity_needed,quantity_purchased,purchased_measurement_unit,purchased_upc,purchase_price,status_explanation,status)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
`

type CreateMealPlanGroceryListItemParams struct {
	ID                       string                `db:"id"`
	BelongsToMealPlan        string                `db:"belongs_to_meal_plan"`
	ValidIngredient          string                `db:"valid_ingredient"`
	ValidMeasurementUnit     string                `db:"valid_measurement_unit"`
	MinimumQuantityNeeded    string                `db:"minimum_quantity_needed"`
	MaximumQuantityNeeded    string                `db:"maximum_quantity_needed"`
	StatusExplanation        string                `db:"status_explanation"`
	Status                   GroceryListItemStatus `db:"status"`
	QuantityPurchased        sql.NullString        `db:"quantity_purchased"`
	PurchasedMeasurementUnit sql.NullString        `db:"purchased_measurement_unit"`
	PurchasedUpc             sql.NullString        `db:"purchased_upc"`
	PurchasePrice            sql.NullString        `db:"purchase_price"`
}

func (q *Queries) CreateMealPlanGroceryListItem(ctx context.Context, db DBTX, arg *CreateMealPlanGroceryListItemParams) error {
	_, err := db.ExecContext(ctx, CreateMealPlanGroceryListItem,
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