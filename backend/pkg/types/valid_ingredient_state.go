package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidIngredientStateCreatedCustomerEventType indicates a valid ingredient state was created.
	ValidIngredientStateCreatedCustomerEventType ServiceEventType = "valid_ingredient_state_created"
	// ValidIngredientStateUpdatedCustomerEventType indicates a valid ingredient state was updated.
	ValidIngredientStateUpdatedCustomerEventType ServiceEventType = "valid_ingredient_state_updated"
	// ValidIngredientStateArchivedCustomerEventType indicates a valid ingredient state was archived.
	ValidIngredientStateArchivedCustomerEventType ServiceEventType = "valid_ingredient_state_archived"

	// ValidIngredientStateAttributeTypeTexture represents the ingredient attribute type for texture.
	ValidIngredientStateAttributeTypeTexture = "texture"
	// ValidIngredientStateAttributeTypeConsistency represents the ingredient attribute type for consistency.
	ValidIngredientStateAttributeTypeConsistency = "consistency"
	// ValidIngredientStateAttributeTypeTemperature represents the ingredient attribute type for temperature.
	ValidIngredientStateAttributeTypeTemperature = "temperature"
	// ValidIngredientStateAttributeTypeColor represents the ingredient attribute type for color.
	ValidIngredientStateAttributeTypeColor = "color"
	// ValidIngredientStateAttributeTypeAppearance represents the ingredient attribute type for appearance.
	ValidIngredientStateAttributeTypeAppearance = "appearance"
	// ValidIngredientStateAttributeTypeOdor represents the ingredient attribute type for odor.
	ValidIngredientStateAttributeTypeOdor = "odor"
	// ValidIngredientStateAttributeTypeTaste represents the ingredient attribute type for taste.
	ValidIngredientStateAttributeTypeTaste = "taste"
	// ValidIngredientStateAttributeTypeSound represents the ingredient attribute type for sound.
	ValidIngredientStateAttributeTypeSound = "sound"
	// ValidIngredientStateAttributeTypeOther represents the ingredient attribute type for other.
	ValidIngredientStateAttributeTypeOther = "other"
)

func init() {
	gob.Register(new(ValidIngredientState))
	gob.Register(new(ValidIngredientStateCreationRequestInput))
	gob.Register(new(ValidIngredientStateUpdateRequestInput))
}

type (
	// ValidIngredientState represents a valid ingredient state.
	ValidIngredientState struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time  `json:"createdAt"`
		ArchivedAt    *time.Time `json:"archivedAt"`
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		PastTense     string     `json:"pastTense"`
		Description   string     `json:"description"`
		IconPath      string     `json:"iconPath"`
		ID            string     `json:"id"`
		Name          string     `json:"name"`
		AttributeType string     `json:"attributeType"`
		Slug          string     `json:"slug"`
	}

	// ValidIngredientStateCreationRequestInput represents what a user could set as input for creating valid ingredient states.
	ValidIngredientStateCreationRequestInput struct {
		_ struct{} `json:"-"`

		Name          string `json:"name"`
		Slug          string `json:"slug"`
		PastTense     string `json:"pastTense"`
		Description   string `json:"description"`
		AttributeType string `json:"attributeType"`
		IconPath      string `json:"iconPath"`
	}

	// ValidIngredientStateDatabaseCreationInput represents what a user could set as input for creating valid ingredient states.
	ValidIngredientStateDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID            string
		Name          string
		Slug          string
		PastTense     string
		Description   string
		AttributeType string
		IconPath      string
	}

	// ValidIngredientStateUpdateRequestInput represents what a user could set as input for updating valid ingredient states.
	ValidIngredientStateUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name          *string `json:"name,omitempty"`
		Slug          *string `json:"slug,omitempty"`
		PastTense     *string `json:"pastTense,omitempty"`
		Description   *string `json:"description,omitempty"`
		AttributeType *string `json:"attributeType,omitempty"`
		IconPath      *string `json:"iconPath,omitempty"`
	}

	// ValidIngredientStateSearchSubset represents the subset of values suitable to index for search.
	ValidIngredientStateSearchSubset struct {
		_ struct{} `json:"-"`

		ID            string `json:"id,omitempty"`
		PastTense     string `json:"pastTense,omitempty"`
		Description   string `json:"description,omitempty"`
		Name          string `json:"name,omitempty"`
		AttributeType string `json:"attributeType,omitempty"`
	}

	// ValidIngredientStateDataManager describes a structure capable of storing valid ingredient states permanently.
	ValidIngredientStateDataManager interface {
		ValidIngredientStateExists(ctx context.Context, validIngredientState string) (bool, error)
		GetValidIngredientState(ctx context.Context, validIngredientState string) (*ValidIngredientState, error)
		GetValidIngredientStates(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidIngredientState], error)
		SearchForValidIngredientStates(ctx context.Context, query string) ([]*ValidIngredientState, error)
		CreateValidIngredientState(ctx context.Context, input *ValidIngredientStateDatabaseCreationInput) (*ValidIngredientState, error)
		UpdateValidIngredientState(ctx context.Context, updated *ValidIngredientState) error
		MarkValidIngredientStateAsIndexed(ctx context.Context, validIngredientState string) error
		ArchiveValidIngredientState(ctx context.Context, validIngredientState string) error
		GetValidIngredientStateIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error)
		GetValidIngredientStatesWithIDs(ctx context.Context, ids []string) ([]*ValidIngredientState, error)
	}

	// ValidIngredientStateDataService describes a structure capable of serving traffic related to valid ingredient states.
	ValidIngredientStateDataService interface {
		SearchHandler(http.ResponseWriter, *http.Request)
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ValidIngredientStateUpdateRequestInput with a valid ingredient state.
func (x *ValidIngredientState) Update(input *ValidIngredientStateUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.IconPath != nil && *input.IconPath != x.IconPath {
		x.IconPath = *input.IconPath
	}

	if input.PastTense != nil && *input.PastTense != x.PastTense {
		x.PastTense = *input.PastTense
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
	}

	if input.AttributeType != nil && *input.AttributeType != x.AttributeType {
		x.AttributeType = *input.AttributeType
	}
}

var _ validation.ValidatableWithContext = (*ValidIngredientStateCreationRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientStateCreationRequestInput.
func (x *ValidIngredientStateCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.AttributeType, validation.In(
			ValidIngredientStateAttributeTypeTexture,
			ValidIngredientStateAttributeTypeConsistency,
			ValidIngredientStateAttributeTypeTemperature,
			ValidIngredientStateAttributeTypeColor,
			ValidIngredientStateAttributeTypeAppearance,
			ValidIngredientStateAttributeTypeOdor,
			ValidIngredientStateAttributeTypeTaste,
			ValidIngredientStateAttributeTypeSound,
			ValidIngredientStateAttributeTypeOther,
		)),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientStateDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidIngredientStateDatabaseCreationInput.
func (x *ValidIngredientStateDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.AttributeType, validation.In(
			ValidIngredientStateAttributeTypeTexture,
			ValidIngredientStateAttributeTypeConsistency,
			ValidIngredientStateAttributeTypeTemperature,
			ValidIngredientStateAttributeTypeColor,
			ValidIngredientStateAttributeTypeAppearance,
			ValidIngredientStateAttributeTypeOdor,
			ValidIngredientStateAttributeTypeTaste,
			ValidIngredientStateAttributeTypeSound,
			ValidIngredientStateAttributeTypeOther,
		)),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientStateUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientStateUpdateRequestInput.
func (x *ValidIngredientStateUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.AttributeType, validation.In(
			ValidIngredientStateAttributeTypeTexture,
			ValidIngredientStateAttributeTypeConsistency,
			ValidIngredientStateAttributeTypeTemperature,
			ValidIngredientStateAttributeTypeColor,
			ValidIngredientStateAttributeTypeAppearance,
			ValidIngredientStateAttributeTypeOdor,
			ValidIngredientStateAttributeTypeTaste,
			ValidIngredientStateAttributeTypeSound,
			ValidIngredientStateAttributeTypeOther,
		)),
	)
}
