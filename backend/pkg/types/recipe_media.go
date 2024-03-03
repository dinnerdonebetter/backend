package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeMediaCreatedCustomerEventType indicates recipe media was created.
	RecipeMediaCreatedCustomerEventType ServiceEventType = "recipe_media_created"
)

func init() {
	gob.Register(new(RecipeMedia))
	gob.Register(new(RecipeMediaCreationRequestInput))
	gob.Register(new(RecipeMediaUpdateRequestInput))
}

type (
	// RecipeMedia represents recipe media.
	RecipeMedia struct {
		_ struct{} `json:"-"`

		CreatedAt           time.Time  `json:"createdAt"`
		ArchivedAt          *time.Time `json:"archivedAt"`
		LastUpdatedAt       *time.Time `json:"lastUpdatedAt"`
		ID                  string     `json:"id"`
		BelongsToRecipe     *string    `json:"belongsToRecipe"`
		BelongsToRecipeStep *string    `json:"belongsToRecipeStep"`
		MimeType            string     `json:"mimeType"`
		InternalPath        string     `json:"internalPath"`
		ExternalPath        string     `json:"externalPath"`
		Index               uint16     `json:"index"`
	}

	// RecipeMediaCreationRequestInput represents what a user could set as input for creating valid preparations.
	RecipeMediaCreationRequestInput struct {
		_ struct{} `json:"-"`

		BelongsToRecipe     *string `json:"belongsToRecipe"`
		BelongsToRecipeStep *string `json:"belongsToRecipeStep"`
		MimeType            string  `json:"mimeType"`
		InternalPath        string  `json:"internalPath"`
		ExternalPath        string  `json:"externalPath"`
		Index               uint16  `json:"index"`
	}

	// RecipeMediaDatabaseCreationInput represents what a user could set as input for creating valid preparations.
	RecipeMediaDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                  string
		BelongsToRecipe     *string
		BelongsToRecipeStep *string
		MimeType            string
		InternalPath        string
		ExternalPath        string
		Index               uint16
	}

	// RecipeMediaUpdateRequestInput represents what a user could set as input for updating valid preparations.
	RecipeMediaUpdateRequestInput struct {
		_ struct{} `json:"-"`

		BelongsToRecipe     *string `json:"belongsToRecipe,omitempty"`
		BelongsToRecipeStep *string `json:"belongsToRecipeStep,omitempty"`
		MimeType            *string `json:"mimeType,omitempty"`
		InternalPath        *string `json:"internalPath,omitempty"`
		ExternalPath        *string `json:"externalPath,omitempty"`
		Index               *uint16 `json:"index,omitempty"`
	}

	// RecipeMediaDataManager describes a structure capable of storing valid preparations permanently.
	RecipeMediaDataManager interface {
		RecipeMediaExists(ctx context.Context, recipeMediaID string) (bool, error)
		GetRecipeMedia(ctx context.Context, recipeMediaID string) (*RecipeMedia, error)
		CreateRecipeMedia(ctx context.Context, input *RecipeMediaDatabaseCreationInput) (*RecipeMedia, error)
		UpdateRecipeMedia(ctx context.Context, updated *RecipeMedia) error
		ArchiveRecipeMedia(ctx context.Context, recipeMediaID string) error
	}

	// RecipeMediaDataService describes a structure capable of serving traffic related to valid preparations.
	RecipeMediaDataService interface {
		SearchHandler(http.ResponseWriter, *http.Request)
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		RandomHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipeMediaUpdateRequestInput with recipe media.
func (x *RecipeMedia) Update(input *RecipeMediaUpdateRequestInput) {
	if input.BelongsToRecipe != nil && input.BelongsToRecipe != x.BelongsToRecipe {
		x.BelongsToRecipe = input.BelongsToRecipe
	}

	if input.BelongsToRecipeStep != nil && input.BelongsToRecipeStep != x.BelongsToRecipeStep {
		x.BelongsToRecipeStep = input.BelongsToRecipeStep
	}

	if input.MimeType != nil && *input.MimeType != x.MimeType {
		x.MimeType = *input.MimeType
	}

	if input.InternalPath != nil && *input.InternalPath != x.InternalPath {
		x.InternalPath = *input.InternalPath
	}

	if input.ExternalPath != nil && *input.ExternalPath != x.ExternalPath {
		x.ExternalPath = *input.ExternalPath
	}

	if input.Index != nil && *input.Index != x.Index {
		x.Index = *input.Index
	}
}

var _ validation.ValidatableWithContext = (*RecipeMediaCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeMediaCreationRequestInput.
func (x *RecipeMediaCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipe, validation.Required),
		validation.Field(&x.MimeType, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.NilOrNotEmpty),
	)
}

var _ validation.ValidatableWithContext = (*RecipeMediaDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeMediaDatabaseCreationInput.
func (x *RecipeMediaDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToRecipe, validation.Required),
		validation.Field(&x.MimeType, validation.Required),
		validation.Field(&x.InternalPath, validation.Required),
		validation.Field(&x.ExternalPath, validation.Required),
		validation.Field(&x.Index, validation.Required),
	)
}
