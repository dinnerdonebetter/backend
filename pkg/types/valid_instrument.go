package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidInstrumentDataType indicates an event is related to a valid instrument.
	ValidInstrumentDataType dataType = "valid_instrument"

	// ValidInstrumentCreatedCustomerEventType indicates a valid instrument was created.
	ValidInstrumentCreatedCustomerEventType CustomerEventType = "valid_instrument_created"
	// ValidInstrumentUpdatedCustomerEventType indicates a valid instrument was updated.
	ValidInstrumentUpdatedCustomerEventType CustomerEventType = "valid_instrument_updated"
	// ValidInstrumentArchivedCustomerEventType indicates a valid instrument was archived.
	ValidInstrumentArchivedCustomerEventType CustomerEventType = "valid_instrument_archived"
)

func init() {
	gob.Register(new(ValidInstrument))
	gob.Register(new(ValidInstrumentCreationRequestInput))
	gob.Register(new(ValidInstrumentUpdateRequestInput))
}

type (
	// ValidInstrument represents a valid instrument.
	ValidInstrument struct {
		_ struct{}

		CreatedAt                      time.Time  `json:"createdAt"`
		LastUpdatedAt                  *time.Time `json:"lastUpdatedAt"`
		ArchivedAt                     *time.Time `json:"archivedAt"`
		IconPath                       string     `json:"iconPath"`
		ID                             string     `json:"id"`
		Name                           string     `json:"name"`
		PluralName                     string     `json:"pluralName"`
		Description                    string     `json:"description"`
		Slug                           string     `json:"slug"`
		DisplayInSummaryLists          bool       `json:"displayInSummaryLists"`
		IncludeInGeneratedInstructions bool       `json:"includeInGeneratedInstructions"`
		UsableForStorage               bool       `json:"usableForStorage"`
		IsVessel                       bool       `json:"isVessel"`
		IsExclusivelyVessel            bool       `json:"isExclusivelyVessel"`
	}

	// NullableValidInstrument represents a fully nullable valid instrument.
	NullableValidInstrument struct {
		_ struct{}

		LastUpdatedAt                  *time.Time
		ArchivedAt                     *time.Time
		Description                    *string
		IconPath                       *string
		ID                             *string
		Name                           *string
		Slug                           *string
		DisplayInSummaryLists          *bool
		IncludeInGeneratedInstructions *bool
		PluralName                     *string
		UsableForStorage               *bool
		IsVessel                       *bool
		IsExclusivelyVessel            *bool
		CreatedAt                      *time.Time
	}

	// ValidInstrumentCreationRequestInput represents what a user could set as input for creating valid instruments.
	ValidInstrumentCreationRequestInput struct {
		_ struct{}

		Name                           string `json:"name"`
		PluralName                     string `json:"pluralName"`
		Description                    string `json:"description"`
		IconPath                       string `json:"iconPath"`
		Slug                           string `json:"slug"`
		DisplayInSummaryLists          bool   `json:"displayInSummaryLists"`
		IncludeInGeneratedInstructions bool   `json:"includeInGeneratedInstructions"`
		UsableForStorage               bool   `json:"usableForStorage"`
		IsExclusivelyVessel            bool   `json:"isExclusivelyVessel"`
		IsVessel                       bool   `json:"isVessel"`
	}

	// ValidInstrumentDatabaseCreationInput represents what a user could set as input for creating valid instruments.
	ValidInstrumentDatabaseCreationInput struct {
		_ struct{}

		ID                             string
		Name                           string
		PluralName                     string
		Description                    string
		IconPath                       string
		Slug                           string
		DisplayInSummaryLists          bool
		UsableForStorage               bool
		IsVessel                       bool
		IsExclusivelyVessel            bool
		IncludeInGeneratedInstructions bool
	}

	// ValidInstrumentUpdateRequestInput represents what a user could set as input for updating valid instruments.
	ValidInstrumentUpdateRequestInput struct {
		_ struct{}

		Name                           *string `json:"name,omitempty"`
		PluralName                     *string `json:"pluralName,omitempty"`
		Description                    *string `json:"description,omitempty"`
		IconPath                       *string `json:"iconPath,omitempty"`
		Slug                           *string `json:"slug,omitempty"`
		UsableForStorage               *bool   `json:"usableForStorage,omitempty"`
		DisplayInSummaryLists          *bool   `json:"displayInSummaryLists,omitempty"`
		IncludeInGeneratedInstructions *bool   `json:"includeInGeneratedInstructions,omitempty"`
		IsVessel                       *bool   `json:"isVessel,omitempty"`
		IsExclusivelyVessel            *bool   `json:"isExclusivelyVessel,omitempty"`
	}

	// ValidInstrumentDataManager describes a structure capable of storing valid instruments permanently.
	ValidInstrumentDataManager interface {
		ValidInstrumentExists(ctx context.Context, validInstrumentID string) (bool, error)
		GetValidInstrument(ctx context.Context, validInstrumentID string) (*ValidInstrument, error)
		GetRandomValidInstrument(ctx context.Context) (*ValidInstrument, error)
		GetValidInstruments(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidInstrument], error)
		SearchForValidInstruments(ctx context.Context, query string) ([]*ValidInstrument, error)
		CreateValidInstrument(ctx context.Context, input *ValidInstrumentDatabaseCreationInput) (*ValidInstrument, error)
		UpdateValidInstrument(ctx context.Context, updated *ValidInstrument) error
		ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error
	}

	// ValidInstrumentDataService describes a structure capable of serving traffic related to valid instruments.
	ValidInstrumentDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		RandomHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidInstrumentUpdateRequestInput with a valid instrument.
func (x *ValidInstrument) Update(input *ValidInstrumentUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.PluralName != nil && *input.PluralName != x.PluralName {
		x.PluralName = *input.PluralName
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.IconPath != nil && *input.IconPath != x.IconPath {
		x.IconPath = *input.IconPath
	}

	if input.UsableForStorage != nil && *input.UsableForStorage != x.UsableForStorage {
		x.UsableForStorage = *input.UsableForStorage
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
	}

	if input.DisplayInSummaryLists != nil && *input.DisplayInSummaryLists != x.DisplayInSummaryLists {
		x.DisplayInSummaryLists = *input.DisplayInSummaryLists
	}

	if input.IsVessel != nil && *input.IsVessel != x.IsVessel {
		x.IsVessel = *input.IsVessel
	}

	if input.IsExclusivelyVessel != nil && *input.IsExclusivelyVessel != x.IsExclusivelyVessel {
		x.IsExclusivelyVessel = *input.IsExclusivelyVessel
	}

	if input.IncludeInGeneratedInstructions != nil && *input.IncludeInGeneratedInstructions != x.IncludeInGeneratedInstructions {
		x.IncludeInGeneratedInstructions = *input.IncludeInGeneratedInstructions
	}
}

var _ validation.ValidatableWithContext = (*ValidInstrumentCreationRequestInput)(nil)

// ValidateWithContext validates a ValidInstrumentCreationRequestInput.
func (x *ValidInstrumentCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidInstrumentDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidInstrumentDatabaseCreationInput.
func (x *ValidInstrumentDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidInstrumentUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidInstrumentUpdateRequestInput.
func (x *ValidInstrumentUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
