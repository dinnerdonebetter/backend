package types

import (
	"encoding/gob"
)

const (
	// RecipePrepTaskStepCreatedCustomerEventType indicates a recipe prep task step was created.
	RecipePrepTaskStepCreatedCustomerEventType ServiceEventType = "recipe_prep_task_step_created"
	// RecipePrepTaskStepUpdatedCustomerEventType indicates a recipe prep task step was updated.
	RecipePrepTaskStepUpdatedCustomerEventType ServiceEventType = "recipe_prep_task_step_updated"
	// RecipePrepTaskStepArchivedCustomerEventType indicates a recipe prep task step was archived.
	RecipePrepTaskStepArchivedCustomerEventType ServiceEventType = "recipe_prep_task_step_archived"
)

func init() {
	gob.Register(new(RecipePrepTask))
	gob.Register(new(RecipePrepTaskCreationRequestInput))
	gob.Register(new(RecipePrepTaskUpdateRequestInput))
}

type (
	// RecipePrepTaskStep represents a recipe prep task step.
	RecipePrepTaskStep struct {
		_ struct{} `json:"-"`

		ID                      string `json:"id"`
		BelongsToRecipeStep     string `json:"belongsToRecipeStep"`
		BelongsToRecipePrepTask string `json:"belongsToRecipeStepTask"`
		SatisfiesRecipeStep     bool   `json:"satisfiesRecipeStep"`
	}

	// RecipePrepTaskStepWithinRecipeCreationRequestInput represents a recipe prep task step.
	RecipePrepTaskStepWithinRecipeCreationRequestInput struct {
		_ struct{} `json:"-"`

		BelongsToRecipeStepIndex uint32 `json:"belongsToRecipeStepIndex"`
		SatisfiesRecipeStep      bool   `json:"satisfiesRecipeStep"`
	}

	// RecipePrepTaskStepCreationRequestInput represents a recipe prep task step.
	RecipePrepTaskStepCreationRequestInput struct {
		_ struct{} `json:"-"`

		BelongsToRecipeStep string `json:"belongsToRecipeStep"`
		SatisfiesRecipeStep bool   `json:"satisfiesRecipeStep"`
	}

	// RecipePrepTaskStepDatabaseCreationInput represents a recipe prep task step.
	RecipePrepTaskStepDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                      string
		BelongsToRecipeStep     string
		BelongsToRecipePrepTask string
		SatisfiesRecipeStep     bool
	}

	// RecipePrepTaskStepUpdateRequestInput represents a recipe prep task step.
	RecipePrepTaskStepUpdateRequestInput struct {
		_ struct{} `json:"-"`

		SatisfiesRecipeStep     *bool   `json:"satisfiesRecipeStep,omitempty"`
		BelongsToRecipeStep     *string `json:"belongsToRecipeStep,omitempty"`
		BelongsToRecipePrepTask *string `json:"belongsToRecipeStepTask,omitempty"`
	}
)
