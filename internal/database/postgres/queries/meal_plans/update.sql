-- name: UpdateMealPlan :exec

UPDATE meal_plans SET notes = $1, status = $2, voting_deadline = $3, last_updated_at = NOW() WHERE archived_at IS NULL AND belongs_to_household = $4 AND id = $5;