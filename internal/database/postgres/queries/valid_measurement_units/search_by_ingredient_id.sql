-- name: SearchValidMeasurementUnitsByIngredientID :many

SELECT
	DISTINCT(valid_measurement_units.id),
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at,
	(
	SELECT
	  COUNT(valid_measurement_units.id)
	FROM
	  valid_measurement_units
	WHERE
	    valid_measurement_units.archived_at IS NULL
	    AND valid_measurement_units.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	    AND valid_measurement_units.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	    AND (
	        valid_measurement_units.last_updated_at IS NULL
	        OR valid_measurement_units.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years'))
	    )
	    AND (
	        valid_measurement_units.last_updated_at IS NULL
	        OR valid_measurement_units.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years'))
	    )
	    AND (
	        valid_ingredient_measurement_units.valid_ingredient_id = $5
	        OR valid_measurement_units.universal = true
	    )
	) as filtered_count,
	(
	    SELECT
	        COUNT(valid_measurement_units.id)
	    FROM
	        valid_measurement_units
	    WHERE
	        valid_measurement_units.archived_at IS NULL
	) as total_count
FROM valid_measurement_units
	FULL OUTER JOIN valid_ingredient_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
	FULL OUTER JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE
	(
	    valid_ingredient_measurement_units.valid_ingredient_id = $5
	    OR valid_measurement_units.universal = true
	)
	AND valid_measurement_units.archived_at IS NULL
	AND valid_ingredient_measurement_units.archived_at IS NULL
	LIMIT $6
	OFFSET $7;
