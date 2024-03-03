package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepVesselCreatedCustomerEventType indicates a recipe step instrument was created.
	RecipeStepVesselCreatedCustomerEventType ServiceEventType = "recipe_step_vessel_created"
	// RecipeStepVesselUpdatedCustomerEventType indicates a recipe step instrument was updated.
	RecipeStepVesselUpdatedCustomerEventType ServiceEventType = "recipe_step_vessel_updated"
	// RecipeStepVesselArchivedCustomerEventType indicates a recipe step instrument was archived.
	RecipeStepVesselArchivedCustomerEventType ServiceEventType = "recipe_step_vessel_archived"
)

func init() {
	gob.Register(new(RecipeStepVessel))
	gob.Register(new(RecipeStepVesselCreationRequestInput))
	gob.Register(new(RecipeStepVesselUpdateRequestInput))
}

type (
	// RecipeStepVessel represents a recipe step instrument.
	RecipeStepVessel struct {
		_ struct{} `json:"-"`

		CreatedAt            time.Time    `json:"createdAt"`
		MaximumQuantity      *uint32      `json:"maximumQuantity"`
		LastUpdatedAt        *time.Time   `json:"lastUpdatedAt"`
		ArchivedAt           *time.Time   `json:"archivedAt"`
		RecipeStepProductID  *string      `json:"recipeStepProductID"`
		Vessel               *ValidVessel `json:"vessel"`
		ID                   string       `json:"id"`
		Notes                string       `json:"notes"`
		BelongsToRecipeStep  string       `json:"belongsToRecipeStep"`
		VesselPreposition    string       `json:"vesselPreposition"`
		Name                 string       `json:"name"`
		MinimumQuantity      uint32       `json:"minimumQuantity"`
		UnavailableAfterStep bool         `json:"unavailableAfterStep"`
	}

	// RecipeStepVesselCreationRequestInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepVesselCreationRequestInput struct {
		_ struct{} `json:"-"`

		RecipeStepProductID             *string `json:"recipeStepProductID"`
		ProductOfRecipeStepIndex        *uint64 `json:"productOfRecipeStepIndex"`
		ProductOfRecipeStepProductIndex *uint64 `json:"productOfRecipeStepProductIndex"`
		VesselID                        *string `json:"vesselID"`
		MaximumQuantity                 *uint32 `json:"maximumQuantity"`
		Name                            string  `json:"name"`
		Notes                           string  `json:"notes"`
		VesselPreposition               string  `json:"vesselPreposition"`
		MinimumQuantity                 uint32  `json:"minimumQuantity"`
		UnavailableAfterStep            bool    `json:"unavailableAfterStep"`
	}

	// RecipeStepVesselDatabaseCreationInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepVesselDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		VesselID                        *string
		RecipeStepProductID             *string
		ProductOfRecipeStepIndex        *uint64
		ProductOfRecipeStepProductIndex *uint64
		MaximumQuantity                 *uint32
		ID                              string
		Notes                           string
		BelongsToRecipeStep             string
		VesselPreposition               string
		Name                            string
		MinimumQuantity                 uint32
		UnavailableAfterStep            bool
	}

	// RecipeStepVesselUpdateRequestInput represents what a user could set as input for updating recipe step instruments.
	RecipeStepVesselUpdateRequestInput struct {
		_ struct{} `json:"-"`

		RecipeStepProductID  *string `json:"recipeStepProductID,omitempty"`
		Name                 *string `json:"name,omitempty"`
		Notes                *string `json:"notes,omitempty"`
		BelongsToRecipeStep  *string `json:"belongsToRecipeStep,omitempty"`
		VesselID             *string `json:"vesselID,omitempty"`
		MinimumQuantity      *uint32 `json:"minimumQuantity,omitempty"`
		MaximumQuantity      *uint32 `json:"maximumQuantity,omitempty"`
		VesselPreposition    *string `json:"vesselPreposition,omitempty"`
		UnavailableAfterStep *bool   `json:"unavailableAfterStep,omitempty"`
	}

	// RecipeStepVesselDataManager describes a structure capable of storing recipe step instruments permanently.
	RecipeStepVesselDataManager interface {
		RecipeStepVesselExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error)
		GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*RecipeStepVessel, error)
		GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*QueryFilteredResult[RecipeStepVessel], error)
		CreateRecipeStepVessel(ctx context.Context, input *RecipeStepVesselDatabaseCreationInput) (*RecipeStepVessel, error)
		UpdateRecipeStepVessel(ctx context.Context, updated *RecipeStepVessel) error
		ArchiveRecipeStepVessel(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error
	}

	// RecipeStepVesselDataService describes a structure capable of serving traffic related to recipe step instruments.
	RecipeStepVesselDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipeStepVesselUpdateRequestInput with a recipe step instrument.
func (x *RecipeStepVessel) Update(input *RecipeStepVesselUpdateRequestInput) {
	if input.RecipeStepProductID != nil && input.RecipeStepProductID != x.RecipeStepProductID {
		x.RecipeStepProductID = input.RecipeStepProductID
	}

	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.BelongsToRecipeStep != nil && *input.BelongsToRecipeStep != x.BelongsToRecipeStep {
		x.BelongsToRecipeStep = *input.BelongsToRecipeStep
	}

	if input.MinimumQuantity != nil && *input.MinimumQuantity != x.MinimumQuantity {
		x.MinimumQuantity = *input.MinimumQuantity
	}

	if input.MaximumQuantity != nil && x.MaximumQuantity != nil && *input.MaximumQuantity != *x.MaximumQuantity {
		x.MaximumQuantity = input.MaximumQuantity
	}

	if input.VesselPreposition != nil && *input.VesselPreposition != x.VesselPreposition {
		x.VesselPreposition = *input.VesselPreposition
	}

	if input.UnavailableAfterStep != nil && *input.UnavailableAfterStep != x.UnavailableAfterStep {
		x.UnavailableAfterStep = *input.UnavailableAfterStep
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.VesselID != nil && x.Vessel != nil && *input.VesselID != x.Vessel.ID {
		x.Vessel = &ValidVessel{ID: *input.VesselID}
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepVesselCreationRequestInput.
func (x *RecipeStepVesselCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MinimumQuantity, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepVesselDatabaseCreationInput.
func (x *RecipeStepVesselDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepVesselUpdateRequestInput.
func (x *RecipeStepVesselUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
	)
}
