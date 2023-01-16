SELECT
	meal_plan_tasks.id,
	meal_plan_options.id,
	meal_plan_options.assigned_cook,
	meal_plan_options.assigned_dishwasher,
	meal_plan_options.chosen,
	meal_plan_options.tiebroken,
	meal_plan_options.meal_id,
	meal_plan_options.notes,
	meal_plan_options.created_at,
	meal_plan_options.last_updated_at,
	meal_plan_options.archived_at,
	meal_plan_options.belongs_to_meal_plan_event,
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count,
	valid_preparations.minimum_instrument_count,
	valid_preparations.maximum_instrument_count,
	valid_preparations.temperature_required,
	valid_preparations.time_estimate_required,
	valid_preparations.condition_expression_required,
    valid_preparations.consumes_vessel,
    valid_preparations.only_for_vessels,
    valid_preparations.minimum_vessel_count,
    valid_preparations.maximum_vessel_count,
	valid_preparations.slug,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.condition_expression,
	recipe_steps.optional,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe,
	meal_plan_tasks.assigned_to_user,
	meal_plan_tasks.status,
	meal_plan_tasks.status_explanation,
	meal_plan_tasks.creation_explanation,
	meal_plan_tasks.created_at,
	meal_plan_tasks.completed_at
FROM meal_plan_tasks
	 FULL OUTER JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
	 FULL OUTER JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id
	 FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
	 JOIN recipe_steps ON meal_plan_tasks.satisfies_recipe_step=recipe_steps.id
	 JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE meal_plan_tasks.belongs_to_meal_plan_option = $1
AND meal_plan_tasks.completed_at IS NULL;
