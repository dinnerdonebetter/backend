-- Begin enumerated types

CREATE TYPE component_type AS ENUM (
    'unspecified',
    'amuse-bouche',
    'appetizer',
    'soup',
    'main',
    'salad',
    'beverage',
    'side',
    'dessert'
);

CREATE TYPE grocery_list_item_status AS ENUM (
    'unknown',
    'already owned',
    'needs',
    'unavailable',
    'acquired'
);

CREATE TYPE ingredient_attribute_type AS ENUM (
    'TEXTure',
    'consistency',
    'color',
    'appearance',
    'odor',
    'taste',
    'sound',
    'temperature',
    'other'
);

CREATE TYPE meal_name AS ENUM (
    'breakfast',
    'second_breakfast',
    'brunch',
    'lunch',
    'supper',
    'dinner'
);

CREATE TYPE meal_plan_status AS ENUM (
    'awaiting_votes',
    'finalized'
);

CREATE TYPE prep_step_status AS ENUM (
    'unfinished',
    'postponed',
    'ignored',
    'canceled',
    'finished'
);

CREATE TYPE recipe_ingredient_scale AS ENUM (
    'linear',
    'logarithmic',
    'exponential'
);

CREATE TYPE recipe_step_product_type AS ENUM (
    'ingredient',
    'instrument',
    'vessel'
);

CREATE TYPE storage_container_type AS ENUM (
    'uncovered',
    'covered',
    'on a wire rack',
    'in an airtight container'
);

CREATE TYPE valid_election_method AS ENUM (
    'schulze',
    'instant-runoff'
);

CREATE TYPE vessel_shape AS ENUM (
    'hemisphere',
    'rectangle',
    'cone',
    'pyramid',
    'cylinder',
    'sphere',
    'cube',
    'other'
);

CREATE TABLE IF NOT EXISTS valid_ingredient_states (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    past_tense TEXT DEFAULT ''::TEXT NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    icon_path TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    attribute_type ingredient_attribute_type DEFAULT 'other'::ingredient_attribute_type NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS valid_ingredients (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    warning TEXT NOT NULL,
    contains_egg BOOLEAN NOT NULL,
    contains_dairy BOOLEAN NOT NULL,
    contains_peanut BOOLEAN NOT NULL,
    contains_tree_nut BOOLEAN NOT NULL,
    contains_soy BOOLEAN NOT NULL,
    contains_wheat BOOLEAN NOT NULL,
    contains_shellfish BOOLEAN NOT NULL,
    contains_sesame BOOLEAN NOT NULL,
    contains_fish BOOLEAN NOT NULL,
    contains_gluten BOOLEAN NOT NULL,
    animal_flesh BOOLEAN NOT NULL,
    volumetric BOOLEAN NOT NULL,
    icon_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    is_liquid BOOLEAN DEFAULT FALSE,
    animal_derived BOOLEAN DEFAULT FALSE NOT NULL,
    plural_name TEXT DEFAULT ''::TEXT NOT NULL,
    restrict_to_preparations BOOLEAN DEFAULT FALSE NOT NULL,
    minimum_ideal_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_ideal_storage_temperature_in_celsius NUMERIC(14,2),
    storage_instructions TEXT DEFAULT ''::TEXT NOT NULL,
    contains_alcohol BOOLEAN DEFAULT FALSE NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    shopping_suggestions TEXT DEFAULT ''::TEXT NOT NULL,
    is_starch BOOLEAN DEFAULT FALSE NOT NULL,
    is_protein BOOLEAN DEFAULT FALSE NOT NULL,
    is_grain BOOLEAN DEFAULT FALSE NOT NULL,
    is_fruit BOOLEAN DEFAULT FALSE NOT NULL,
    is_salt BOOLEAN DEFAULT FALSE NOT NULL,
    is_fat BOOLEAN DEFAULT FALSE NOT NULL,
    is_acid BOOLEAN DEFAULT FALSE NOT NULL,
    is_heat BOOLEAN DEFAULT FALSE NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS valid_ingredient_groups (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT DEFAULT ''::TEXT NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS valid_ingredient_group_members (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_group TEXT NOT NULL REFERENCES valid_ingredient_groups("id") ON DELETE CASCADE,
    valid_ingredient TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS valid_measurement_units (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    icon_path TEXT DEFAULT ''::TEXT NOT NULL,
    volumetric BOOLEAN DEFAULT FALSE,
    universal BOOLEAN DEFAULT FALSE NOT NULL,
    metric BOOLEAN DEFAULT FALSE NOT NULL,
    imperial BOOLEAN DEFAULT FALSE NOT NULL,
    plural_name TEXT DEFAULT ''::TEXT NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_ingredient_measurement_units (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    valid_ingredient_id TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    valid_measurement_unit_id TEXT NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    minimum_allowable_quantity NUMERIC(14,2) DEFAULT 0 NOT NULL,
    maximum_allowable_quantity NUMERIC(14,2),
    UNIQUE(valid_ingredient_id, valid_measurement_unit_id, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_measurement_unit_conversions (
    id TEXT NOT NULL PRIMARY KEY,
    from_unit TEXT NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    to_unit TEXT NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    only_for_ingredient TEXT REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    modifier NUMERIC(14,5) NOT NULL,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(from_unit, to_unit, only_for_ingredient)
);

CREATE TABLE IF NOT EXISTS valid_ingredient_state_ingredients (
    id TEXT NOT NULL PRIMARY KEY,
    valid_ingredient TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    valid_ingredient_state TEXT NOT NULL REFERENCES valid_ingredient_states("id") ON DELETE CASCADE,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_ingredient, valid_ingredient_state, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_preparations (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    yields_nothing BOOLEAN DEFAULT FALSE NOT NULL,
    restrict_to_ingredients BOOLEAN DEFAULT FALSE NOT NULL,
    past_tense TEXT DEFAULT ''::TEXT NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    minimum_ingredient_count INTEGER DEFAULT 1 NOT NULL,
    maximum_ingredient_count INTEGER,
    minimum_instrument_count INTEGER DEFAULT 1 NOT NULL,
    maximum_instrument_count INTEGER,
    temperature_required BOOLEAN DEFAULT FALSE NOT NULL,
    time_estimate_required BOOLEAN DEFAULT FALSE NOT NULL,
    condition_expression_required BOOLEAN DEFAULT FALSE NOT NULL,
    consumes_vessel BOOLEAN DEFAULT FALSE NOT NULL,
    only_for_vessels BOOLEAN DEFAULT FALSE NOT NULL,
    minimum_vessel_count INTEGER DEFAULT 0 NOT NULL,
    maximum_vessel_count INTEGER,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_ingredient_preparations (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT NOT NULL,
    valid_preparation_id TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    valid_ingredient_id TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_preparation_id, valid_ingredient_id, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_instruments (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    plural_name TEXT DEFAULT ''::TEXT NOT NULL,
    usable_for_storage BOOLEAN DEFAULT FALSE NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    display_in_summary_lists BOOLEAN DEFAULT true NOT NULL,
    include_in_generated_instructions BOOLEAN DEFAULT true NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS valid_preparation_instruments (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    valid_preparation_id TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    valid_instrument_id TEXT NOT NULL REFERENCES valid_instruments("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_instrument_id, valid_preparation_id, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_vessels (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    plural_name TEXT DEFAULT ''::TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    icon_path TEXT NOT NULL,
    usable_for_storage BOOLEAN DEFAULT FALSE NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    display_in_summary_lists BOOLEAN DEFAULT true NOT NULL,
    include_in_generated_instructions BOOLEAN DEFAULT true NOT NULL,
    capacity NUMERIC(14,2) DEFAULT 0 NOT NULL,
    capacity_unit TEXT REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    width_in_millimeters NUMERIC(14,2),
    length_in_millimeters NUMERIC(14,2),
    height_in_millimeters NUMERIC(14,2),
    shape vessel_shape DEFAULT 'other'::vessel_shape NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name, archived_at),
    UNIQUE(slug, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_preparation_vessels (
    id TEXT NOT NULL PRIMARY KEY,
    valid_preparation_id TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    valid_vessel_id TEXT NOT NULL REFERENCES valid_vessels("id") ON DELETE CASCADE,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_preparation_id, valid_vessel_id, archived_at)
);

-- Begin non-enumerated types

CREATE TABLE IF NOT EXISTS user_ingredient_preferences (
    id TEXT NOT NULL PRIMARY KEY,
    ingredient TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    rating smallint DEFAULT 0 NOT NULL,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    allergy BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    UNIQUE(belongs_to_user, ingredient)
);

CREATE TABLE IF NOT EXISTS household_instrument_ownerships (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    quantity INTEGER DEFAULT 0 NOT NULL,
    valid_instrument_id TEXT NOT NULL,
    belongs_to_household TEXT NOT NULL REFERENCES households("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_instrument_id, belongs_to_household, archived_at)
);

CREATE TABLE IF NOT EXISTS recipes (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    source TEXT NOT NULL,
    description TEXT NOT NULL,
    inspired_by_recipe_id TEXT REFERENCES recipes("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    created_by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    min_estimated_portions NUMERIC(14,2) DEFAULT 1 NOT NULL,
    seal_of_approval BOOLEAN DEFAULT FALSE NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    portion_name TEXT DEFAULT 'portion'::TEXT NOT NULL,
    plural_portion_name TEXT DEFAULT 'portions'::TEXT NOT NULL,
    max_estimated_portions NUMERIC(14,2),
    eligible_for_meals BOOLEAN DEFAULT true NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    last_validated_at TIMESTAMP WITH TIME ZONE,
    yields_component_type component_type DEFAULT 'unspecified'::component_type NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_steps (
    id TEXT NOT NULL PRIMARY KEY,
    index INTEGER NOT NULL,
    preparation_id TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    minimum_estimated_time_in_seconds bigint,
    maximum_estimated_time_in_seconds bigint,
    minimum_temperature_in_celsius NUMERIC(14,2),
    notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    optional BOOLEAN DEFAULT FALSE NOT NULL,
    maximum_temperature_in_celsius NUMERIC(14,2),
    explicit_instructions TEXT DEFAULT ''::TEXT NOT NULL,
    condition_expression TEXT DEFAULT ''::TEXT NOT NULL,
    start_timer_automatically BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_step_products (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    quantity_notes TEXT NOT NULL,
    minimum_quantity_value NUMERIC(14,2),
    maximum_quantity_value NUMERIC(14,2),
    measurement_unit TEXT REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    type recipe_step_product_type NOT NULL,
    compostable BOOLEAN DEFAULT FALSE NOT NULL,
    maximum_storage_duration_in_seconds INTEGER,
    minimum_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_storage_temperature_in_celsius NUMERIC(14,2),
    storage_instructions TEXT DEFAULT ''::TEXT NOT NULL,
    is_liquid BOOLEAN DEFAULT FALSE NOT NULL,
    is_waste BOOLEAN DEFAULT FALSE NOT NULL,
    index INTEGER DEFAULT 0 NOT NULL,
    contained_in_vessel_index INTEGER
);

CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
    id TEXT NOT NULL PRIMARY KEY,
    ingredient_id TEXT REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    minimum_quantity_value NUMERIC(14,2) NOT NULL,
    quantity_notes TEXT NOT NULL,
    ingredient_notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    name TEXT NOT NULL,
    recipe_step_product_id TEXT REFERENCES recipe_step_products("id") ON DELETE CASCADE,
    maximum_quantity_value NUMERIC(14,2),
    measurement_unit TEXT REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    optional BOOLEAN DEFAULT FALSE NOT NULL,
    option_index INTEGER DEFAULT 0 NOT NULL,
    vessel_index INTEGER,
    to_taste BOOLEAN DEFAULT FALSE NOT NULL,
    product_percentage_to_use NUMERIC(14,2),
    recipe_step_product_recipe_id TEXT REFERENCES recipes("id") ON DELETE CASCADE,
    CONSTRAINT valid_instrument_or_product CHECK (((recipe_step_product_id IS NOT NULL) OR (ingredient_id IS NOT NULL))),
    UNIQUE(ingredient_id, belongs_to_recipe_step)
);

CREATE TABLE IF NOT EXISTS recipe_step_instruments (
    id TEXT NOT NULL PRIMARY KEY,
    instrument_id TEXT REFERENCES valid_instruments("id") ON DELETE CASCADE,
    notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    preference_rank INTEGER NOT NULL,
    recipe_step_product_id TEXT REFERENCES recipe_step_products("id") ON DELETE CASCADE,
    name TEXT DEFAULT ''::TEXT NOT NULL,
    optional BOOLEAN DEFAULT FALSE NOT NULL,
    minimum_quantity INTEGER DEFAULT 1 NOT NULL,
    maximum_quantity INTEGER,
    option_index INTEGER DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_step_vessels (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT DEFAULT ''::TEXT NOT NULL,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    recipe_step_product_id TEXT REFERENCES recipe_step_products("id") ON DELETE CASCADE,
    vessel_predicate TEXT DEFAULT ''::TEXT NOT NULL,
    minimum_quantity INTEGER DEFAULT 0 NOT NULL,
    maximum_quantity INTEGER,
    unavailable_after_step BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    valid_vessel_id TEXT REFERENCES valid_vessels("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_step_completion_conditions (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    ingredient_state TEXT NOT NULL REFERENCES valid_ingredient_states("id") ON DELETE CASCADE,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    optional BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(belongs_to_recipe_step, ingredient_state)
);

CREATE TABLE IF NOT EXISTS recipe_step_completion_condition_ingredients (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe_step_completion_condition TEXT NOT NULL REFERENCES recipe_step_completion_conditions("id") ON DELETE CASCADE,
    recipe_step_ingredient TEXT NOT NULL REFERENCES recipe_step_ingredients("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS recipe_media (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe TEXT REFERENCES recipes("id") ON DELETE CASCADE,
    belongs_to_recipe_step TEXT REFERENCES recipe_steps("id") ON DELETE CASCADE,
    mime_type TEXT NOT NULL,
    internal_path TEXT NOT NULL,
    external_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    index INTEGER DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_prep_tasks (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    explicit_storage_instructions TEXT DEFAULT ''::TEXT NOT NULL,
    minimum_time_buffer_before_recipe_in_seconds INTEGER NOT NULL,
    maximum_time_buffer_before_recipe_in_seconds INTEGER,
    storage_type storage_container_type,
    minimum_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_storage_temperature_in_celsius NUMERIC(14,2),
    belongs_to_recipe TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    name TEXT DEFAULT ''::TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    optional BOOLEAN DEFAULT true NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_prep_task_steps (
    id TEXT NOT NULL PRIMARY KEY,
    satisfies_recipe_step BOOLEAN DEFAULT FALSE NOT NULL,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    belongs_to_recipe_prep_task TEXT NOT NULL REFERENCES recipe_prep_tasks("id") ON DELETE CASCADE,
    UNIQUE(belongs_to_recipe_step, belongs_to_recipe_prep_task)
);

CREATE TABLE IF NOT EXISTS recipe_ratings (
    id TEXT NOT NULL PRIMARY KEY,
    recipe_id TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    taste NUMERIC(14,2),
    difficulty NUMERIC(14,2),
    cleanup NUMERIC(14,2),
    instructions NUMERIC(14,2),
    overall NUMERIC(14,2),
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(by_user, recipe_id)
);

CREATE TABLE IF NOT EXISTS meals (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    created_by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    min_estimated_portions NUMERIC(14,2) DEFAULT 1.0 NOT NULL,
    max_estimated_portions NUMERIC(14,2),
    eligible_for_meal_plans BOOLEAN DEFAULT true NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS meal_components (
    id TEXT NOT NULL PRIMARY KEY,
    meal_id TEXT NOT NULL REFERENCES meals("id") ON DELETE CASCADE,
    recipe_id TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    meal_component_type component_type DEFAULT 'unspecified'::component_type NOT NULL,
    recipe_scale NUMERIC(14,2) DEFAULT 1.0 NOT NULL
);

CREATE TABLE IF NOT EXISTS meal_plans (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT NOT NULL,
    status meal_plan_status DEFAULT 'awaiting_votes'::meal_plan_status NOT NULL,
    voting_deadline TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_household TEXT NOT NULL REFERENCES households("id") ON DELETE CASCADE,
    grocery_list_initialized BOOLEAN DEFAULT FALSE NOT NULL,
    tasks_created BOOLEAN DEFAULT FALSE NOT NULL,
    election_method valid_election_method DEFAULT 'schulze'::valid_election_method NOT NULL,
    created_by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS meal_plan_grocery_list_items (
    id TEXT NOT NULL PRIMARY KEY,
    valid_ingredient TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    valid_measurement_unit TEXT NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    minimum_quantity_needed NUMERIC(14,2) NOT NULL,
    maximum_quantity_needed NUMERIC(14,2),
    quantity_purchased NUMERIC(14,2),
    purchased_measurement_unit TEXT REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    purchased_upc TEXT,
    purchase_price NUMERIC(14,2),
    status_explanation TEXT DEFAULT ''::TEXT NOT NULL,
    status grocery_list_item_status DEFAULT 'unknown'::grocery_list_item_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_meal_plan TEXT NOT NULL REFERENCES meal_plans("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS meal_plan_events (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    starts_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    ends_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    meal_name meal_name NOT NULL,
    belongs_to_meal_plan TEXT NOT NULL REFERENCES meal_plans("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS meal_plan_options (
    id TEXT NOT NULL PRIMARY KEY,
    meal_id TEXT NOT NULL REFERENCES meals("id") ON DELETE CASCADE,
    notes TEXT NOT NULL,
    chosen BOOLEAN DEFAULT FALSE NOT NULL,
    tiebroken BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    assigned_cook TEXT REFERENCES users("id") ON DELETE CASCADE,
    assigned_dishwasher TEXT REFERENCES users("id") ON DELETE CASCADE,
    belongs_to_meal_plan_event TEXT REFERENCES meal_plan_events("id") ON DELETE CASCADE,
    meal_scale NUMERIC(14,2) DEFAULT 1.0 NOT NULL
);

CREATE TABLE IF NOT EXISTS meal_plan_option_votes (
    id TEXT NOT NULL PRIMARY KEY,
    rank INTEGER NOT NULL,
    abstain BOOLEAN NOT NULL,
    notes TEXT NOT NULL,
    by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_meal_plan_option TEXT NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    UNIQUE(by_user, belongs_to_meal_plan_option)
);

CREATE TABLE IF NOT EXISTS meal_plan_tasks (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_meal_plan_option TEXT NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    belongs_to_recipe_prep_task TEXT NOT NULL REFERENCES recipe_prep_tasks("id") ON DELETE CASCADE,
    creation_explanation TEXT DEFAULT ''::TEXT NOT NULL,
    status_explanation TEXT DEFAULT ''::TEXT NOT NULL,
    status prep_step_status DEFAULT 'unfinished'::prep_step_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    assigned_to_user TEXT REFERENCES users("id") ON DELETE CASCADE,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS meal_plan_events_belongs_to_meal_pla_index ON meal_plan_events USING btree (belongs_to_meal_plan);
CREATE INDEX IF NOT EXISTS meal_plan_grocery_list_items_belongs_to_meal_pla_index ON meal_plan_grocery_list_items USING btree (belongs_to_meal_plan);
CREATE INDEX IF NOT EXISTS meal_plan_grocery_list_items_purchased_measurement_unit_index ON meal_plan_grocery_list_items USING btree (purchased_measurement_unit);
CREATE INDEX IF NOT EXISTS meal_plan_grocery_list_items_valid_ingredient_index ON meal_plan_grocery_list_items USING btree (valid_ingredient);
CREATE INDEX IF NOT EXISTS meal_plan_grocery_list_items_valid_measurement_unit_index ON meal_plan_grocery_list_items USING btree (valid_measurement_unit);
CREATE INDEX IF NOT EXISTS meal_plan_options_assigned_cook_index ON meal_plan_options USING btree (assigned_cook);
CREATE INDEX IF NOT EXISTS meal_plan_options_assigned_dishwasher_index ON meal_plan_options USING btree (assigned_dishwasher);
CREATE INDEX IF NOT EXISTS meal_plan_options_belongs_to_meal_plan_even_index ON meal_plan_options USING btree (belongs_to_meal_plan_event);
CREATE INDEX IF NOT EXISTS meal_plan_options_belongs_to_meal_plan_option ON meal_plan_option_votes USING btree (belongs_to_meal_plan_option);
CREATE INDEX IF NOT EXISTS meal_plan_tasks_assigned_to_user_index ON meal_plan_tasks USING btree (assigned_to_user);
CREATE INDEX IF NOT EXISTS meal_plan_tasks_belongs_to_meal_plan_option_index ON meal_plan_tasks USING btree (belongs_to_meal_plan_option);
CREATE INDEX IF NOT EXISTS meal_plan_tasks_belongs_to_recipe_prep_task_index ON meal_plan_tasks USING btree (belongs_to_recipe_prep_task);
CREATE INDEX IF NOT EXISTS meal_plans_belongs_to_household ON meal_plans USING btree (belongs_to_household);
CREATE INDEX IF NOT EXISTS meal_recipes_meal_id ON meal_components USING btree (meal_id);
CREATE INDEX IF NOT EXISTS meal_recipes_recipe_id ON meal_components USING btree (recipe_id);
CREATE INDEX IF NOT EXISTS meals_created_by_user ON meals USING btree (created_by_user);
CREATE INDEX IF NOT EXISTS recipe_media_belongs_to_recipe_index ON recipe_media USING btree (belongs_to_recipe);
CREATE INDEX IF NOT EXISTS recipe_media_belongs_to_recipe_step_index ON recipe_media USING btree (belongs_to_recipe_step);
CREATE INDEX IF NOT EXISTS recipe_prep_task_steps_belongs_to_recipe_prep_task_index ON recipe_prep_task_steps USING btree (belongs_to_recipe_prep_task);
CREATE INDEX IF NOT EXISTS recipe_prep_task_steps_belongs_to_recipe_step_index ON recipe_prep_task_steps USING btree (belongs_to_recipe_step);
CREATE INDEX IF NOT EXISTS recipe_prep_tasks_belongs_to_recipe_index ON recipe_prep_tasks USING btree (belongs_to_recipe);
CREATE INDEX IF NOT EXISTS recipe_step_ingredients_measurement_unit_index ON recipe_step_ingredients USING btree (measurement_unit);
CREATE INDEX IF NOT EXISTS recipe_step_ingredients_product_of_recipe_step ON recipe_step_ingredients USING btree (recipe_step_product_id);
CREATE INDEX IF NOT EXISTS recipe_step_instruments_instrument_id_index ON recipe_step_instruments USING btree (instrument_id);
CREATE INDEX IF NOT EXISTS recipe_step_instruments_recipe_step_product_id_index ON recipe_step_instruments USING btree (recipe_step_product_id);
CREATE INDEX IF NOT EXISTS recipe_step_products_belongs_to_recipe_step ON recipe_step_products USING btree (belongs_to_recipe_step);
CREATE INDEX IF NOT EXISTS recipe_step_products_measurement_unit_index ON recipe_step_products USING btree (measurement_unit);
CREATE INDEX IF NOT EXISTS recipe_steps_belongs_to_recipe ON recipe_steps USING btree (belongs_to_recipe);
CREATE INDEX IF NOT EXISTS recipes_created_by_user ON recipes USING btree (created_by_user);
CREATE INDEX IF NOT EXISTS valid_ingredient_measurement_units_valid_ingredient_id_index ON valid_ingredient_measurement_units USING btree (valid_ingredient_id);
CREATE INDEX IF NOT EXISTS valid_ingredient_measurement_units_valid_measurement_unit_id_in ON valid_ingredient_measurement_units USING btree (valid_measurement_unit_id);
CREATE INDEX IF NOT EXISTS valid_ingredient_state_ingredients_referncing_valid_ingredient_ ON valid_ingredient_state_ingredients USING btree (valid_ingredient);
CREATE INDEX IF NOT EXISTS valid_measurement_conversions_from_unit_index ON valid_measurement_unit_conversions USING btree (from_unit);
CREATE INDEX IF NOT EXISTS valid_measurement_conversions_only_for_ingredient_index ON valid_measurement_unit_conversions USING btree (only_for_ingredient);
CREATE INDEX IF NOT EXISTS valid_measurement_conversions_to_unit_index ON valid_measurement_unit_conversions USING btree (to_unit);
CREATE INDEX IF NOT EXISTS valid_preparation_instruments_valid_instrument_index ON valid_preparation_instruments USING btree (valid_instrument_id);
CREATE INDEX IF NOT EXISTS valid_preparation_instruments_valid_preparation_index ON valid_preparation_instruments USING btree (valid_preparation_id);
CREATE INDEX IF NOT EXISTS valid_preparation_vessels_referencing_valid_preparations_idx ON valid_preparation_vessels USING btree (valid_preparation_id);
CREATE INDEX IF NOT EXISTS valid_preparation_vessels_referencing_valid_vessels_idx ON valid_preparation_vessels USING btree (valid_vessel_id);
