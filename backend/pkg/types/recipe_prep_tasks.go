package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

const (
	// RecipePrepTaskStorageTypeUncovered is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeUncovered = "uncovered"
	// RecipePrepTaskStorageTypeCovered is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeCovered = "covered"
	// RecipePrepTaskStorageTypeWireRack is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeWireRack = "on a wire rack"
	// RecipePrepTaskStorageTypeAirtightContainer is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeAirtightContainer = "in an airtight container"

	// RecipePrepTaskCreatedCustomerEventType indicates a recipe prep task was created.
	RecipePrepTaskCreatedCustomerEventType ServiceEventType = "recipe_created"
	// RecipePrepTaskUpdatedCustomerEventType indicates a recipe prep task was updated.
	RecipePrepTaskUpdatedCustomerEventType ServiceEventType = "recipe_updated"
	// RecipePrepTaskArchivedCustomerEventType indicates a recipe prep task was archived.
	RecipePrepTaskArchivedCustomerEventType ServiceEventType = "recipe_archived"
)

func init() {
	gob.Register(new(RecipePrepTask))
	gob.Register(new(RecipePrepTaskCreationRequestInput))
	gob.Register(new(RecipePrepTaskUpdateRequestInput))
}

type (
	// RecipePrepTask represents a recipe prep task.
	RecipePrepTask struct {
		_ struct{} `json:"-"`

		CreatedAt                              time.Time             `json:"createdAt"`
		MaximumStorageTemperatureInCelsius     *float32              `json:"maximumStorageTemperatureInCelsius"`
		ArchivedAt                             *time.Time            `json:"archivedAt"`
		LastUpdatedAt                          *time.Time            `json:"lastUpdatedAt"`
		MinimumStorageTemperatureInCelsius     *float32              `json:"minimumStorageTemperatureInCelsius"`
		MaximumTimeBufferBeforeRecipeInSeconds *uint32               `json:"maximumTimeBufferBeforeRecipeInSeconds"`
		ID                                     string                `json:"id"`
		StorageType                            string                `json:"storageType"`
		BelongsToRecipe                        string                `json:"belongsToRecipe"`
		ExplicitStorageInstructions            string                `json:"explicitStorageInstructions"`
		Notes                                  string                `json:"notes"`
		Name                                   string                `json:"name"`
		Description                            string                `json:"description"`
		TaskSteps                              []*RecipePrepTaskStep `json:"recipeSteps"`
		MinimumTimeBufferBeforeRecipeInSeconds uint32                `json:"minimumTimeBufferBeforeRecipeInSeconds"`
		Optional                               bool                  `json:"optional"`
	}

	// RecipePrepTaskCreationRequestInput represents what a user could set as input for creating recipes.
	RecipePrepTaskCreationRequestInput struct {
		_ struct{} `json:"-"`

		MaximumTimeBufferBeforeRecipeInSeconds *uint32                                   `json:"maximumTimeBufferBeforeRecipeInSeconds"`
		MinimumStorageTemperatureInCelsius     *float32                                  `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius     *float32                                  `json:"maximumStorageTemperatureInCelsius"`
		StorageType                            string                                    `json:"storageType"`
		ExplicitStorageInstructions            string                                    `json:"explicitStorageInstructions"`
		Notes                                  string                                    `json:"notes"`
		Name                                   string                                    `json:"name"`
		Description                            string                                    `json:"description"`
		BelongsToRecipe                        string                                    `json:"belongsToRecipe"`
		TaskSteps                              []*RecipePrepTaskStepCreationRequestInput `json:"recipeSteps"`
		MinimumTimeBufferBeforeRecipeInSeconds uint32                                    `json:"minimumTimeBufferBeforeRecipeInSeconds"`
		Optional                               bool                                      `json:"optional"`
	}

	// RecipePrepTaskWithinRecipeCreationRequestInput represents what a user could set as input for creating recipes.
	RecipePrepTaskWithinRecipeCreationRequestInput struct {
		_ struct{} `json:"-"`

		MaximumTimeBufferBeforeRecipeInSeconds *uint32                                               `json:"maximumTimeBufferBeforeRecipeInSeconds"`
		MinimumStorageTemperatureInCelsius     *float32                                              `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius     *float32                                              `json:"maximumStorageTemperatureInCelsius"`
		StorageType                            string                                                `json:"storageType"`
		Name                                   string                                                `json:"name"`
		Description                            string                                                `json:"description"`
		ExplicitStorageInstructions            string                                                `json:"explicitStorageInstructions"`
		Notes                                  string                                                `json:"notes"`
		BelongsToRecipe                        string                                                `json:"belongsToRecipe"`
		TaskSteps                              []*RecipePrepTaskStepWithinRecipeCreationRequestInput `json:"recipeSteps"`
		MinimumTimeBufferBeforeRecipeInSeconds uint32                                                `json:"minimumTimeBufferBeforeRecipeInSeconds"`
		Optional                               bool                                                  `json:"optional"`
	}

	// RecipePrepTaskDatabaseCreationInput represents what a user could set as input for creating recipes.
	RecipePrepTaskDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		MaximumTimeBufferBeforeRecipeInSeconds *uint32
		MinimumStorageTemperatureInCelsius     *float32
		MaximumStorageTemperatureInCelsius     *float32
		ExplicitStorageInstructions            string
		Notes                                  string
		ID                                     string
		Name                                   string
		Description                            string
		StorageType                            string
		BelongsToRecipe                        string
		TaskSteps                              []*RecipePrepTaskStepDatabaseCreationInput
		MinimumTimeBufferBeforeRecipeInSeconds uint32
		Optional                               bool
	}

	// RecipePrepTaskUpdateRequestInput represents what a user could set as input for updating recipes.
	RecipePrepTaskUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes                                  *string                                 `json:"notes,omitempty"`
		ExplicitStorageInstructions            *string                                 `json:"explicitStorageInstructions,omitempty"`
		MinimumTimeBufferBeforeRecipeInSeconds *uint32                                 `json:"minimumTimeBufferBeforeRecipeInSeconds,omitempty"`
		MaximumTimeBufferBeforeRecipeInSeconds *uint32                                 `json:"maximumTimeBufferBeforeRecipeInSeconds,omitempty"`
		StorageType                            *string                                 `json:"storageType,omitempty"`
		Name                                   *string                                 `json:"name"`
		Optional                               *bool                                   `json:"optional"`
		Description                            *string                                 `json:"description"`
		MinimumStorageTemperatureInCelsius     *float32                                `json:"minimumStorageTemperatureInCelsius,omitempty"`
		MaximumStorageTemperatureInCelsius     *float32                                `json:"maximumStorageTemperatureInCelsius,omitempty"`
		BelongsToRecipe                        *string                                 `json:"belongsToRecipe,omitempty"`
		TaskSteps                              []*RecipePrepTaskStepUpdateRequestInput `json:"recipeSteps,omitempty"`
	}

	// RecipePrepTaskDataManager describes a structure capable of storing recipes permanently.
	RecipePrepTaskDataManager interface {
		RecipePrepTaskExists(ctx context.Context, recipeID, recipePrepTaskID string) (bool, error)
		GetRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*RecipePrepTask, error)
		GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) ([]*RecipePrepTask, error)
		CreateRecipePrepTask(ctx context.Context, input *RecipePrepTaskDatabaseCreationInput) (*RecipePrepTask, error)
		UpdateRecipePrepTask(ctx context.Context, updated *RecipePrepTask) error
		ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error
	}

	// RecipePrepTaskDataService describes a structure capable of serving traffic related to recipes.
	RecipePrepTaskDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipePrepTaskUpdateRequestInput with a recipe prep task.
func (x *RecipePrepTask) Update(input *RecipePrepTaskUpdateRequestInput) {
	if input.BelongsToRecipe != nil && *input.BelongsToRecipe != x.BelongsToRecipe {
		x.BelongsToRecipe = *input.BelongsToRecipe
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ExplicitStorageInstructions != nil && *input.ExplicitStorageInstructions != x.ExplicitStorageInstructions {
		x.ExplicitStorageInstructions = *input.ExplicitStorageInstructions
	}

	if input.MinimumTimeBufferBeforeRecipeInSeconds != nil && *input.MinimumTimeBufferBeforeRecipeInSeconds != x.MinimumTimeBufferBeforeRecipeInSeconds {
		x.MinimumTimeBufferBeforeRecipeInSeconds = *input.MinimumTimeBufferBeforeRecipeInSeconds
	}

	if input.MaximumTimeBufferBeforeRecipeInSeconds != nil && input.MaximumTimeBufferBeforeRecipeInSeconds != x.MaximumTimeBufferBeforeRecipeInSeconds {
		x.MaximumTimeBufferBeforeRecipeInSeconds = input.MaximumTimeBufferBeforeRecipeInSeconds
	}

	if input.StorageType != nil && *input.StorageType != x.StorageType {
		x.StorageType = *input.StorageType
	}

	if input.MinimumStorageTemperatureInCelsius != nil && input.MinimumStorageTemperatureInCelsius != x.MinimumStorageTemperatureInCelsius {
		x.MinimumStorageTemperatureInCelsius = input.MinimumStorageTemperatureInCelsius
	}

	if input.MaximumStorageTemperatureInCelsius != nil && input.MaximumStorageTemperatureInCelsius != x.MaximumStorageTemperatureInCelsius {
		x.MaximumStorageTemperatureInCelsius = input.MaximumStorageTemperatureInCelsius
	}

	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}
}

var _ validation.ValidatableWithContext = (*RecipePrepTaskCreationRequestInput)(nil)

// ValidateWithContext validates a RecipePrepTaskCreationRequestInput.
func (x *RecipePrepTaskCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	var result *multierror.Error

	if err := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipe, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.StorageType, validation.Required),
		validation.Field(&x.MaximumTimeBufferBeforeRecipeInSeconds, validation.Required),
		validation.Field(&x.MinimumStorageTemperatureInCelsius, validation.Required),
		validation.Field(&x.MaximumStorageTemperatureInCelsius, validation.Required),
		validation.Field(&x.MinimumTimeBufferBeforeRecipeInSeconds, validation.Required),
	); err != nil {
		result = multierror.Append(err, result)
	}

	// TODO: uncomment me
	// if x.MinimumStorageTemperatureInCelsius != nil && x.MaximumStorageTemperatureInCelsius != nil && *x.MinimumStorageTemperatureInCelsius > *x.MaximumStorageTemperatureInCelsius {
	// 	result = multierror.Append(fmt.Errorf("minimum storage temperature (%d) is greater than maximum storage temperature (%d)", x.MinimumStorageTemperatureInCelsius, x.MaximumStorageTemperatureInCelsius))
	// }

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*RecipePrepTaskDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipePrepTaskDatabaseCreationInput.
func (x *RecipePrepTaskDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToRecipe, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipePrepTaskUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipePrepTaskUpdateRequestInput.
func (x *RecipePrepTaskUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipe, validation.Required),
	)
}
