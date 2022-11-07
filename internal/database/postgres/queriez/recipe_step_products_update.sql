-- name: UpdateRecipeStepProduct :exec
UPDATE recipe_step_products
SET
	name = $1,
	type = $2,
	measurement_unit = $3,
	minimum_quantity_value = $4,
	maximum_quantity_value = $5,
	quantity_notes = $6,
	compostable = $7,
	maximum_storage_duration_in_seconds = $8,
	minimum_storage_temperature_in_celsius = $9,
	maximum_storage_temperature_in_celsius = $10,
	storage_instructions = $11,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $12
	AND id = $13;