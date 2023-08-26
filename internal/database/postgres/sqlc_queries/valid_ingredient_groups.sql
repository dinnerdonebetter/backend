-- name: ArchiveValidIngredientGroup :execrows

UPDATE valid_ingredient_groups SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: ArchiveValidIngredientGroupMember :execrows

UPDATE valid_ingredient_group_members SET archived_at = NOW() WHERE id = $1 AND belongs_to_group = $2;

-- name: CreateValidIngredientGroup :exec

INSERT INTO valid_ingredient_groups
(
	id,
	"name",
	description,
	slug
) VALUES ($1,$2,$3,$4);

-- name: CreateValidIngredientGroupMember :exec

INSERT INTO valid_ingredient_group_members
(
	id,
	belongs_to_group,
	valid_ingredient
) VALUES ($1,$2,$3);

-- name: CheckValidIngredientGroupExistence :one

SELECT EXISTS ( SELECT valid_ingredient_groups.id FROM valid_ingredient_groups WHERE valid_ingredient_groups.archived_at IS NULL AND valid_ingredient_groups.id = $1 );

-- name: GetValidIngredientGroups :many

SELECT
    valid_ingredient_groups.id,
    valid_ingredient_groups.name,
    valid_ingredient_groups.description,
    valid_ingredient_groups.slug,
    valid_ingredient_groups.created_at,
    valid_ingredient_groups.last_updated_at,
    valid_ingredient_groups.archived_at,
    (
        SELECT
            COUNT(valid_ingredient_groups.id)
        FROM
            valid_ingredient_groups
        WHERE
            valid_ingredient_groups.archived_at IS NULL
          AND valid_ingredient_groups.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND valid_ingredient_groups.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (
                valid_ingredient_groups.last_updated_at IS NULL
                OR valid_ingredient_groups.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
            )
          AND (
                valid_ingredient_groups.last_updated_at IS NULL
                OR valid_ingredient_groups.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
            )
        OFFSET sqlc.narg(query_offset)
    ) AS filtered_count,
    (
        SELECT
            COUNT(valid_ingredient_groups.id)
        FROM
            valid_ingredient_groups
        WHERE
            valid_ingredient_groups.archived_at IS NULL
    ) AS total_count
FROM valid_ingredient_groups
WHERE
	valid_ingredient_groups.archived_at IS NULL
	AND valid_ingredient_groups.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
	AND valid_ingredient_groups.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
	AND (
	    valid_ingredient_groups.last_updated_at IS NULL
	    OR valid_ingredient_groups.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
	)
	AND (
	    valid_ingredient_groups.last_updated_at IS NULL
	    OR valid_ingredient_groups.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
	)
	OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);

-- name: GetValidIngredientGroupMembers :many

SELECT
	valid_ingredient_group_members.id,
    valid_ingredient_group_members.belongs_to_group,
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
    valid_ingredient_group_members.created_at,
    valid_ingredient_group_members.archived_at
FROM valid_ingredient_group_members
  JOIN valid_ingredients ON valid_ingredients.id = valid_ingredient_group_members.valid_ingredient
WHERE valid_ingredient_group_members.archived_at IS NULL
AND valid_ingredient_group_members.belongs_to_group = $1;

-- name: GetValidIngredientGroup :one

SELECT
	valid_ingredient_groups.id,
	valid_ingredient_groups.name,
	valid_ingredient_groups.description,
	valid_ingredient_groups.slug,
	valid_ingredient_groups.created_at,
	valid_ingredient_groups.last_updated_at,
	valid_ingredient_groups.archived_at
FROM valid_ingredient_groups
WHERE valid_ingredient_groups.archived_at IS NULL
AND valid_ingredient_groups.id = $1;

-- name: SearchForValidIngredientGroups :many

SELECT
    valid_ingredient_groups.id,
    valid_ingredient_groups.name,
    valid_ingredient_groups.description,
    valid_ingredient_groups.slug,
    valid_ingredient_groups.created_at,
    valid_ingredient_groups.last_updated_at,
    valid_ingredient_groups.archived_at,
    (
        SELECT
            COUNT(valid_ingredient_groups.id)
        FROM
            valid_ingredient_groups
        WHERE
            valid_ingredient_groups.archived_at IS NULL
          AND valid_ingredient_groups.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND valid_ingredient_groups.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (
                valid_ingredient_groups.last_updated_at IS NULL
                OR valid_ingredient_groups.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
            )
          AND (
                valid_ingredient_groups.last_updated_at IS NULL
                OR valid_ingredient_groups.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
            )
        OFFSET sqlc.narg(query_offset)
    ) AS filtered_count,
    (
        SELECT
            COUNT(valid_ingredient_groups.id)
        FROM
            valid_ingredient_groups
        WHERE
            valid_ingredient_groups.archived_at IS NULL
    ) AS total_count
FROM valid_ingredient_groups
WHERE
    valid_ingredient_groups.archived_at IS NULL
  AND valid_ingredient_groups.name ILIKE '%' || sqlc.arg(name)::text || '%'
  AND valid_ingredient_groups.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
  AND valid_ingredient_groups.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
  AND (
    valid_ingredient_groups.last_updated_at IS NULL
   OR valid_ingredient_groups.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
    )
  AND (
    valid_ingredient_groups.last_updated_at IS NULL
   OR valid_ingredient_groups.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
    )
OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);

-- name: UpdateValidIngredientGroup :execrows

UPDATE valid_ingredient_groups SET
	name = $1,
	description = $2,
	slug = $3,
	last_updated_at = NOW()
WHERE archived_at IS NULL AND id = $4;
