package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidPreparationInstrumentCreatedCustomerEventType indicates a valid preparation instrument was created.
	ValidPreparationInstrumentCreatedCustomerEventType ServiceEventType = "valid_preparation_instrument_created"
	// ValidPreparationInstrumentUpdatedCustomerEventType indicates a valid preparation instrument was updated.
	ValidPreparationInstrumentUpdatedCustomerEventType ServiceEventType = "valid_preparation_instrument_updated"
	// ValidPreparationInstrumentArchivedCustomerEventType indicates a valid preparation instrument was archived.
	ValidPreparationInstrumentArchivedCustomerEventType ServiceEventType = "valid_preparation_instrument_archived"
)

func init() {
	gob.Register(new(ValidPreparationInstrument))
	gob.Register(new(ValidPreparationInstrumentCreationRequestInput))
	gob.Register(new(ValidPreparationInstrumentUpdateRequestInput))
}

type (
	// ValidPreparationInstrument represents a valid preparation instrument.
	ValidPreparationInstrument struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time        `json:"createdAt"`
		LastUpdatedAt *time.Time       `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time       `json:"archivedAt"`
		ID            string           `json:"id"`
		Notes         string           `json:"notes"`
		Instrument    ValidInstrument  `json:"instrument"`
		Preparation   ValidPreparation `json:"preparation"`
	}

	// ValidPreparationInstrumentCreationRequestInput represents what a user could set as input for creating valid preparation instruments.
	ValidPreparationInstrumentCreationRequestInput struct {
		_ struct{} `json:"-"`

		Notes              string `json:"notes"`
		ValidPreparationID string `json:"validPreparationID"`
		ValidInstrumentID  string `json:"validInstrumentID"`
	}

	// ValidPreparationInstrumentDatabaseCreationInput represents what a user could set as input for creating valid preparation instruments.
	ValidPreparationInstrumentDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                 string
		Notes              string
		ValidPreparationID string
		ValidInstrumentID  string
	}

	// ValidPreparationInstrumentUpdateRequestInput represents what a user could set as input for updating valid preparation instruments.
	ValidPreparationInstrumentUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes              *string `json:"notes,omitempty"`
		ValidPreparationID *string `json:"validPreparationID,omitempty"`
		ValidInstrumentID  *string `json:"validInstrumentID,omitempty"`
	}

	// ValidPreparationInstrumentDataManager describes a structure capable of storing valid preparation instruments permanently.
	ValidPreparationInstrumentDataManager interface {
		ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID string) (bool, error)
		GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*ValidPreparationInstrument, error)
		GetValidPreparationInstruments(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidPreparationInstrument], error)
		GetValidPreparationInstrumentsForPreparation(ctx context.Context, preparationID string, filter *QueryFilter) (*QueryFilteredResult[ValidPreparationInstrument], error)
		GetValidPreparationInstrumentsForInstrument(ctx context.Context, instrumentID string, filter *QueryFilter) (*QueryFilteredResult[ValidPreparationInstrument], error)
		CreateValidPreparationInstrument(ctx context.Context, input *ValidPreparationInstrumentDatabaseCreationInput) (*ValidPreparationInstrument, error)
		UpdateValidPreparationInstrument(ctx context.Context, updated *ValidPreparationInstrument) error
		ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error
	}

	// ValidPreparationInstrumentDataService describes a structure capable of serving traffic related to valid preparation instruments.
	ValidPreparationInstrumentDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
		SearchByPreparationHandler(http.ResponseWriter, *http.Request)
		SearchByInstrumentHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ValidPreparationInstrumentUpdateRequestInput with a valid preparation instrument.
func (x *ValidPreparationInstrument) Update(input *ValidPreparationInstrumentUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ValidPreparationID != nil && *input.ValidPreparationID != x.Preparation.ID {
		x.Preparation.ID = *input.ValidPreparationID
	}

	if input.ValidInstrumentID != nil && *input.ValidInstrumentID != x.Instrument.ID {
		x.Instrument.ID = *input.ValidInstrumentID
	}
}

var _ validation.ValidatableWithContext = (*ValidPreparationInstrumentCreationRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationInstrumentCreationRequestInput.
func (x *ValidPreparationInstrumentCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationInstrumentDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidPreparationInstrumentDatabaseCreationInput.
func (x *ValidPreparationInstrumentDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationInstrumentUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationInstrumentUpdateRequestInput.
func (x *ValidPreparationInstrumentUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}
