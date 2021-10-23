CREATE TABLE IF NOT EXISTS valid_ingredients (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"name" TEXT NOT NULL,
	"variant" TEXT NOT NULL,
	"description" TEXT NOT NULL,
	"warning" TEXT NOT NULL,
	"contains_egg" BOOLEAN NOT NULL,
	"contains_dairy" BOOLEAN NOT NULL,
	"contains_peanut" BOOLEAN NOT NULL,
	"contains_tree_nut" BOOLEAN NOT NULL,
	"contains_soy" BOOLEAN NOT NULL,
	"contains_wheat" BOOLEAN NOT NULL,
	"contains_shellfish" BOOLEAN NOT NULL,
	"contains_sesame" BOOLEAN NOT NULL,
	"contains_fish" BOOLEAN NOT NULL,
	"contains_gluten" BOOLEAN NOT NULL,
	"animal_flesh" BOOLEAN NOT NULL,
	"animal_derived" BOOLEAN NOT NULL,
	"volumetric" BOOLEAN NOT NULL,
	"icon_path" TEXT NOT NULL,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL
);