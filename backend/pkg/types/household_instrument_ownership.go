package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// HouseholdInstrumentOwnershipCreatedCustomerEventType indicates a household instrument ownership was created.
	HouseholdInstrumentOwnershipCreatedCustomerEventType ServiceEventType = "household_instrument_ownership_created"
	// HouseholdInstrumentOwnershipUpdatedCustomerEventType indicates a household instrument ownership was updated.
	HouseholdInstrumentOwnershipUpdatedCustomerEventType ServiceEventType = "household_instrument_ownership_updated"
	// HouseholdInstrumentOwnershipArchivedCustomerEventType indicates a household instrument ownership was archived.
	HouseholdInstrumentOwnershipArchivedCustomerEventType ServiceEventType = "household_instrument_ownership_archived"
)

func init() {
	gob.Register(new(HouseholdInstrumentOwnership))
	gob.Register(new(HouseholdInstrumentOwnershipCreationRequestInput))
	gob.Register(new(HouseholdInstrumentOwnershipUpdateRequestInput))
}

type (
	// HouseholdInstrumentOwnership represents a household instrument ownership.
	HouseholdInstrumentOwnership struct {
		_ struct{} `json:"-"`

		CreatedAt          time.Time       `json:"createdAt"`
		ArchivedAt         *time.Time      `json:"archivedAt"`
		LastUpdatedAt      *time.Time      `json:"lastUpdatedAt"`
		ID                 string          `json:"id"`
		Notes              string          `json:"notes"`
		BelongsToHousehold string          `json:"belongsToHousehold"`
		Instrument         ValidInstrument `json:"instrument"`
		Quantity           uint16          `json:"quantity"`
	}

	// HouseholdInstrumentOwnershipCreationRequestInput represents what a user could set as input for creating household instrument ownerships.
	HouseholdInstrumentOwnershipCreationRequestInput struct {
		_ struct{} `json:"-"`

		Notes              string `json:"notes"`
		ValidInstrumentID  string `json:"validInstrumentID"`
		BelongsToHousehold string `json:"belongsToHousehold"`
		Quantity           uint16 `json:"quantity"`
	}

	// HouseholdInstrumentOwnershipDatabaseCreationInput represents what a user could set as input for creating household instrument ownerships.
	HouseholdInstrumentOwnershipDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                 string
		Notes              string
		ValidInstrumentID  string
		BelongsToHousehold string
		Quantity           uint16
	}

	// HouseholdInstrumentOwnershipUpdateRequestInput represents what a user could set as input for updating household instrument ownerships.
	HouseholdInstrumentOwnershipUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes             *string `json:"notes"`
		Quantity          *uint16 `json:"quantity"`
		ValidInstrumentID *string `json:"validInstrumentID"`
	}

	// HouseholdInstrumentOwnershipDataManager describes a structure capable of storing household instrument ownerships permanently.
	HouseholdInstrumentOwnershipDataManager interface {
		HouseholdInstrumentOwnershipExists(ctx context.Context, householdInstrumentOwnershipID, householdID string) (bool, error)
		GetHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) (*HouseholdInstrumentOwnership, error)
		GetHouseholdInstrumentOwnerships(ctx context.Context, householdID string, filter *QueryFilter) (*QueryFilteredResult[HouseholdInstrumentOwnership], error)
		CreateHouseholdInstrumentOwnership(ctx context.Context, input *HouseholdInstrumentOwnershipDatabaseCreationInput) (*HouseholdInstrumentOwnership, error)
		UpdateHouseholdInstrumentOwnership(ctx context.Context, updated *HouseholdInstrumentOwnership) error
		ArchiveHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) error
	}

	// HouseholdInstrumentOwnershipDataService describes a structure capable of serving traffic related to household instrument ownerships.
	HouseholdInstrumentOwnershipDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an HouseholdInstrumentOwnershipUpdateRequestInput with a household instrument ownership.
func (x *HouseholdInstrumentOwnership) Update(input *HouseholdInstrumentOwnershipUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.Quantity != nil && *input.Quantity != x.Quantity {
		x.Quantity = *input.Quantity
	}

	if input.ValidInstrumentID != nil && *input.ValidInstrumentID != x.Instrument.ID {
		x.Instrument = ValidInstrument{ID: *input.ValidInstrumentID}
	}
}

var _ validation.ValidatableWithContext = (*HouseholdInstrumentOwnershipCreationRequestInput)(nil)

// ValidateWithContext validates a HouseholdInstrumentOwnershipCreationRequestInput.
func (x *HouseholdInstrumentOwnershipCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Quantity, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*HouseholdInstrumentOwnershipDatabaseCreationInput)(nil)

// ValidateWithContext validates a HouseholdInstrumentOwnershipDatabaseCreationInput.
func (x *HouseholdInstrumentOwnershipDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Quantity, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*HouseholdInstrumentOwnershipUpdateRequestInput)(nil)

// ValidateWithContext validates a HouseholdInstrumentOwnershipUpdateRequestInput.
func (x *HouseholdInstrumentOwnershipUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Quantity, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}
