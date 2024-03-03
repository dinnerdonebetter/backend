package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidPreparationVesselCreatedCustomerEventType indicates a valid preparation instrument was created.
	ValidPreparationVesselCreatedCustomerEventType ServiceEventType = "valid_preparation_instrument_created"
	// ValidPreparationVesselUpdatedCustomerEventType indicates a valid preparation instrument was updated.
	ValidPreparationVesselUpdatedCustomerEventType ServiceEventType = "valid_preparation_instrument_updated"
	// ValidPreparationVesselArchivedCustomerEventType indicates a valid preparation instrument was archived.
	ValidPreparationVesselArchivedCustomerEventType ServiceEventType = "valid_preparation_instrument_archived"
)

func init() {
	gob.Register(new(ValidPreparationVessel))
	gob.Register(new(ValidPreparationVesselCreationRequestInput))
	gob.Register(new(ValidPreparationVesselUpdateRequestInput))
}

type (
	// ValidPreparationVessel represents a valid preparation instrument.
	ValidPreparationVessel struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time        `json:"createdAt"`
		LastUpdatedAt *time.Time       `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time       `json:"archivedAt"`
		ID            string           `json:"id"`
		Notes         string           `json:"notes"`
		Vessel        ValidVessel      `json:"instrument"`
		Preparation   ValidPreparation `json:"preparation"`
	}

	// ValidPreparationVesselCreationRequestInput represents what a user could set as input for creating valid preparation instruments.
	ValidPreparationVesselCreationRequestInput struct {
		_ struct{} `json:"-"`

		Notes              string `json:"notes"`
		ValidPreparationID string `json:"validPreparationID"`
		ValidVesselID      string `json:"validVesselID"`
	}

	// ValidPreparationVesselDatabaseCreationInput represents what a user could set as input for creating valid preparation instruments.
	ValidPreparationVesselDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                 string
		Notes              string
		ValidPreparationID string
		ValidVesselID      string
	}

	// ValidPreparationVesselUpdateRequestInput represents what a user could set as input for updating valid preparation instruments.
	ValidPreparationVesselUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes              *string `json:"notes,omitempty"`
		ValidPreparationID *string `json:"validPreparationID,omitempty"`
		ValidVesselID      *string `json:"validVesselID,omitempty"`
	}

	// ValidPreparationVesselDataManager describes a structure capable of storing valid preparation instruments permanently.
	ValidPreparationVesselDataManager interface {
		ValidPreparationVesselExists(ctx context.Context, validPreparationVesselID string) (bool, error)
		GetValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*ValidPreparationVessel, error)
		GetValidPreparationVessels(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidPreparationVessel], error)
		GetValidPreparationVesselsForPreparation(ctx context.Context, preparationID string, filter *QueryFilter) (*QueryFilteredResult[ValidPreparationVessel], error)
		GetValidPreparationVesselsForVessel(ctx context.Context, instrumentID string, filter *QueryFilter) (*QueryFilteredResult[ValidPreparationVessel], error)
		CreateValidPreparationVessel(ctx context.Context, input *ValidPreparationVesselDatabaseCreationInput) (*ValidPreparationVessel, error)
		UpdateValidPreparationVessel(ctx context.Context, updated *ValidPreparationVessel) error
		ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error
	}

	// ValidPreparationVesselDataService describes a structure capable of serving traffic related to valid preparation instruments.
	ValidPreparationVesselDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
		SearchByPreparationHandler(http.ResponseWriter, *http.Request)
		SearchByVesselHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ValidPreparationVesselUpdateRequestInput with a valid preparation instrument.
func (x *ValidPreparationVessel) Update(input *ValidPreparationVesselUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ValidPreparationID != nil && *input.ValidPreparationID != x.Preparation.ID {
		x.Preparation.ID = *input.ValidPreparationID
	}

	if input.ValidVesselID != nil && *input.ValidVesselID != x.Vessel.ID {
		x.Vessel.ID = *input.ValidVesselID
	}
}

var _ validation.ValidatableWithContext = (*ValidPreparationVesselCreationRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationVesselCreationRequestInput.
func (x *ValidPreparationVesselCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidVesselID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationVesselDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidPreparationVesselDatabaseCreationInput.
func (x *ValidPreparationVesselDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidVesselID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationVesselUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationVesselUpdateRequestInput.
func (x *ValidPreparationVesselUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidVesselID, validation.Required),
	)
}
