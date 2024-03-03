package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidVesselCreatedCustomerEventType indicates a valid vessel was created.
	ValidVesselCreatedCustomerEventType ServiceEventType = "valid_vessel_created"
	// ValidVesselUpdatedCustomerEventType indicates a valid vessel was updated.
	ValidVesselUpdatedCustomerEventType ServiceEventType = "valid_vessel_updated"
	// ValidVesselArchivedCustomerEventType indicates a valid vessel was archived.
	ValidVesselArchivedCustomerEventType ServiceEventType = "valid_vessel_archived"
)

func init() {
	gob.Register(new(ValidVessel))
	gob.Register(new(ValidVesselCreationRequestInput))
	gob.Register(new(ValidVesselUpdateRequestInput))
}

type (
	// ValidVessel represents a valid vessel.
	ValidVessel struct {
		_ struct{} `json:"-"`

		CreatedAt                      time.Time             `json:"createdAt"`
		ArchivedAt                     *time.Time            `json:"archivedAt"`
		LastUpdatedAt                  *time.Time            `json:"lastUpdatedAt"`
		CapacityUnit                   *ValidMeasurementUnit `json:"capacityUnit"`
		IconPath                       string                `json:"iconPath"`
		PluralName                     string                `json:"pluralName"`
		Description                    string                `json:"description"`
		Name                           string                `json:"name"`
		Slug                           string                `json:"slug"`
		Shape                          string                `json:"shape"`
		ID                             string                `json:"id"`
		WidthInMillimeters             float32               `json:"widthInMillimeters"`
		LengthInMillimeters            float32               `json:"lengthInMillimeters"`
		HeightInMillimeters            float32               `json:"heightInMillimeters"`
		Capacity                       float32               `json:"capacity"`
		IncludeInGeneratedInstructions bool                  `json:"includeInGeneratedInstructions"`
		DisplayInSummaryLists          bool                  `json:"displayInSummaryLists"`
		UsableForStorage               bool                  `json:"usableForStorage"`
	}

	// NullableValidVessel represents a fully nullable valid vessel.
	NullableValidVessel struct {
		_ struct{} `json:"-"`

		ID                             *string                       `json:"id"`
		Name                           *string                       `json:"name"`
		PluralName                     *string                       `json:"pluralName"`
		Description                    *string                       `json:"description"`
		IconPath                       *string                       `json:"iconPath"`
		UsableForStorage               *bool                         `json:"usableForStorage"`
		Slug                           *string                       `json:"slug"`
		DisplayInSummaryLists          *bool                         `json:"displayInSummaryLists"`
		IncludeInGeneratedInstructions *bool                         `json:"includeInGeneratedInstructions"`
		Capacity                       *float32                      `json:"capacity"`
		CapacityUnit                   *NullableValidMeasurementUnit `json:"capacityUnit"`
		WidthInMillimeters             *float32                      `json:"widthInMillimeters"`
		LengthInMillimeters            *float32                      `json:"lengthInMillimeters"`
		HeightInMillimeters            *float32                      `json:"heightInMillimeters"`
		Shape                          *string                       `json:"shape"`
		CreatedAt                      *time.Time                    `json:"createdAt"`
		LastUpdatedAt                  *time.Time                    `json:"lastUpdatedAt"`
		ArchivedAt                     *time.Time                    `json:"archivedAt"`
	}

	// ValidVesselCreationRequestInput represents what a user could set as input for creating valid vessels.
	ValidVesselCreationRequestInput struct {
		_ struct{} `json:"-"`

		CapacityUnitID                 *string `json:"capacityUnitID"`
		Shape                          string  `json:"shape"`
		IconPath                       string  `json:"iconPath"`
		PluralName                     string  `json:"pluralName"`
		Name                           string  `json:"name"`
		Description                    string  `json:"description"`
		Slug                           string  `json:"slug"`
		LengthInMillimeters            float32 `json:"lengthInMillimeters"`
		HeightInMillimeters            float32 `json:"heightInMillimeters"`
		Capacity                       float32 `json:"capacity"`
		WidthInMillimeters             float32 `json:"widthInMillimeters"`
		UsableForStorage               bool    `json:"usableForStorage"`
		IncludeInGeneratedInstructions bool    `json:"includeInGeneratedInstructions"`
		DisplayInSummaryLists          bool    `json:"displayInSummaryLists"`
	}

	// ValidVesselDatabaseCreationInput represents what a user could set as input for creating valid vessels.
	ValidVesselDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		CapacityUnitID                 *string
		ID                             string
		Name                           string
		PluralName                     string
		Description                    string
		IconPath                       string
		Shape                          string
		Slug                           string
		WidthInMillimeters             float32
		Capacity                       float32
		LengthInMillimeters            float32
		HeightInMillimeters            float32
		IncludeInGeneratedInstructions bool
		DisplayInSummaryLists          bool
		UsableForStorage               bool
	}

	// ValidVesselUpdateRequestInput represents what a user could set as input for updating valid vessels.
	ValidVesselUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name                           *string  `json:"name"`
		PluralName                     *string  `json:"pluralName"`
		Description                    *string  `json:"description"`
		IconPath                       *string  `json:"iconPath"`
		UsableForStorage               *bool    `json:"usableForStorage"`
		Slug                           *string  `json:"slug"`
		DisplayInSummaryLists          *bool    `json:"displayInSummaryLists"`
		IncludeInGeneratedInstructions *bool    `json:"includeInGeneratedInstructions"`
		Capacity                       *float32 `json:"capacity"`
		CapacityUnitID                 *string  `json:"capacityUnitID"`
		WidthInMillimeters             *float32 `json:"widthInMillimeters"`
		LengthInMillimeters            *float32 `json:"lengthInMillimeters"`
		HeightInMillimeters            *float32 `json:"heightInMillimeters"`
		Shape                          *string  `json:"shape"`
	}

	// ValidVesselSearchSubset represents the subset of values suitable to index for search.
	ValidVesselSearchSubset struct {
		_ struct{} `json:"-"`

		ID               string  `json:"id,omitempty"`
		Name             string  `json:"name,omitempty"`
		PluralName       string  `json:"pluralName,omitempty"`
		Description      string  `json:"description,omitempty"`
		CapacityUnitName string  `json:"capacityUnitName"`
		Capacity         float32 `json:"capacity,omitempty"`
	}

	// ValidVesselDataManager describes a structure capable of storing valid vessels permanently.
	ValidVesselDataManager interface {
		ValidVesselExists(ctx context.Context, validVesselID string) (bool, error)
		GetValidVessel(ctx context.Context, validVesselID string) (*ValidVessel, error)
		GetRandomValidVessel(ctx context.Context) (*ValidVessel, error)
		GetValidVessels(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidVessel], error)
		SearchForValidVessels(ctx context.Context, query string) ([]*ValidVessel, error)
		CreateValidVessel(ctx context.Context, input *ValidVesselDatabaseCreationInput) (*ValidVessel, error)
		UpdateValidVessel(ctx context.Context, updated *ValidVessel) error
		MarkValidVesselAsIndexed(ctx context.Context, validVesselID string) error
		ArchiveValidVessel(ctx context.Context, validVesselID string) error
		GetValidVesselIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error)
		GetValidVesselsWithIDs(ctx context.Context, ids []string) ([]*ValidVessel, error)
	}

	// ValidVesselDataService describes a structure capable of serving traffic related to valid vessels.
	ValidVesselDataService interface {
		SearchHandler(http.ResponseWriter, *http.Request)
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		RandomHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ValidVesselUpdateRequestInput with a valid vessel.
func (x *ValidVessel) Update(input *ValidVesselUpdateRequestInput) {
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

	if input.IncludeInGeneratedInstructions != nil && *input.IncludeInGeneratedInstructions != x.IncludeInGeneratedInstructions {
		x.IncludeInGeneratedInstructions = *input.IncludeInGeneratedInstructions
	}

	if input.Capacity != nil && *input.Capacity != x.Capacity {
		x.Capacity = *input.Capacity
	}

	if input.CapacityUnitID != nil && x.CapacityUnit != nil && *input.CapacityUnitID != x.CapacityUnit.ID {
		x.CapacityUnit = &ValidMeasurementUnit{ID: *input.CapacityUnitID}
	}

	if input.WidthInMillimeters != nil && *input.WidthInMillimeters != x.WidthInMillimeters {
		x.WidthInMillimeters = *input.WidthInMillimeters
	}

	if input.LengthInMillimeters != nil && *input.LengthInMillimeters != x.LengthInMillimeters {
		x.LengthInMillimeters = *input.LengthInMillimeters
	}

	if input.HeightInMillimeters != nil && *input.HeightInMillimeters != x.HeightInMillimeters {
		x.HeightInMillimeters = *input.HeightInMillimeters
	}

	if input.Shape != nil && *input.Shape != x.Shape {
		x.Shape = *input.Shape
	}
}

var _ validation.ValidatableWithContext = (*ValidVesselCreationRequestInput)(nil)

// ValidateWithContext validates a ValidVesselCreationRequestInput.
func (x *ValidVesselCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Capacity, validation.When(x.CapacityUnitID != nil, validation.Required)),
		validation.Field(&x.Shape, validation.In(
			"hemisphere",
			"rectangle",
			"cone",
			"pyramid",
			"cylinder",
			"sphere",
			"cube",
			"other",
		)),
	)
}

var _ validation.ValidatableWithContext = (*ValidVesselDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidVesselDatabaseCreationInput.
func (x *ValidVesselDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.CapacityUnitID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidVesselUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidVesselUpdateRequestInput.
func (x *ValidVesselUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
