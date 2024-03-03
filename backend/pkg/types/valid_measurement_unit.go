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
	// ValidMeasurementUnitCreatedCustomerEventType indicates a valid measurement unit was created.
	ValidMeasurementUnitCreatedCustomerEventType ServiceEventType = "valid_measurement_unit_created"
	// ValidMeasurementUnitUpdatedCustomerEventType indicates a valid measurement unit was updated.
	ValidMeasurementUnitUpdatedCustomerEventType ServiceEventType = "valid_measurement_unit_updated"
	// ValidMeasurementUnitArchivedCustomerEventType indicates a valid measurement unit was archived.
	ValidMeasurementUnitArchivedCustomerEventType ServiceEventType = "valid_measurement_unit_archived"
)

func init() {
	gob.Register(new(ValidMeasurementUnit))
	gob.Register(new(ValidMeasurementUnitCreationRequestInput))
	gob.Register(new(ValidMeasurementUnitUpdateRequestInput))
}

type (
	// ValidMeasurementUnit represents a valid measurement unit.
	ValidMeasurementUnit struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time  `json:"createdAt"`
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time `json:"archivedAt"`
		Name          string     `json:"name"`
		IconPath      string     `json:"iconPath"`
		ID            string     `json:"id"`
		Description   string     `json:"description"`
		PluralName    string     `json:"pluralName"`
		Slug          string     `json:"slug"`
		Volumetric    bool       `json:"volumetric"`
		Universal     bool       `json:"universal"`
		Metric        bool       `json:"metric"`
		Imperial      bool       `json:"imperial"`
	}

	// NullableValidMeasurementUnit represents a nullable valid measurement unit.
	NullableValidMeasurementUnit struct {
		_ struct{} `json:"-"`

		CreatedAt     *time.Time
		LastUpdatedAt *time.Time
		ArchivedAt    *time.Time
		Name          *string
		IconPath      *string
		ID            *string
		Description   *string
		PluralName    *string
		Slug          *string
		Volumetric    *bool
		Universal     *bool
		Metric        *bool
		Imperial      *bool
	}

	// ValidMeasurementUnitCreationRequestInput represents what a user could set as input for creating valid measurement units.
	ValidMeasurementUnitCreationRequestInput struct {
		_ struct{} `json:"-"`

		Name        string `json:"name"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
		PluralName  string `json:"pluralName"`
		Slug        string `json:"slug"`
		Volumetric  bool   `json:"volumetric"`
		Universal   bool   `json:"universal"`
		Metric      bool   `json:"metric"`
		Imperial    bool   `json:"imperial"`
	}

	// ValidMeasurementUnitDatabaseCreationInput represents what a user could set as input for creating valid measurement units.
	ValidMeasurementUnitDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		Name        string
		Description string
		ID          string
		IconPath    string
		PluralName  string
		Slug        string
		Volumetric  bool
		Universal   bool
		Metric      bool
		Imperial    bool
	}

	// ValidMeasurementUnitUpdateRequestInput represents what a user could set as input for updating valid measurement units.
	ValidMeasurementUnitUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
		IconPath    *string `json:"iconPath,omitempty"`
		Volumetric  *bool   `json:"volumetric,omitempty"`
		Universal   *bool   `json:"universal,omitempty"`
		Metric      *bool   `json:"metric,omitempty"`
		Imperial    *bool   `json:"imperial,omitempty"`
		PluralName  *string `json:"pluralName,omitempty"`
		Slug        *string `json:"slug,omitempty"`
	}

	// ValidMeasurementUnitSearchSubset represents the subset of values suitable to index for search.
	ValidMeasurementUnitSearchSubset struct {
		_ struct{} `json:"-"`

		Name        string `json:"name,omitempty"`
		ID          string `json:"id,omitempty"`
		Description string `json:"description,omitempty"`
		PluralName  string `json:"pluralName,omitempty"`
	}

	// ValidMeasurementUnitDataManager describes a structure capable of storing valid measurement units permanently.
	ValidMeasurementUnitDataManager interface {
		ValidMeasurementUnitExists(ctx context.Context, validMeasurementUnitID string) (bool, error)
		GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*ValidMeasurementUnit, error)
		GetValidMeasurementUnits(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidMeasurementUnit], error)
		SearchForValidMeasurementUnits(ctx context.Context, query string) ([]*ValidMeasurementUnit, error)
		ValidMeasurementUnitsForIngredientID(ctx context.Context, validIngredientID string, filter *QueryFilter) (*QueryFilteredResult[ValidMeasurementUnit], error)
		CreateValidMeasurementUnit(ctx context.Context, input *ValidMeasurementUnitDatabaseCreationInput) (*ValidMeasurementUnit, error)
		UpdateValidMeasurementUnit(ctx context.Context, updated *ValidMeasurementUnit) error
		MarkValidMeasurementUnitAsIndexed(ctx context.Context, validMeasurementUnitID string) error
		ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error
		GetValidMeasurementUnitIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error)
		GetValidMeasurementUnitsWithIDs(ctx context.Context, ids []string) ([]*ValidMeasurementUnit, error)
	}

	// ValidMeasurementUnitDataService describes a structure capable of serving traffic related to valid measurement units.
	ValidMeasurementUnitDataService interface {
		SearchHandler(http.ResponseWriter, *http.Request)
		SearchByIngredientIDHandler(http.ResponseWriter, *http.Request)
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ValidMeasurementUnitUpdateRequestInput with a valid measurement unit.
func (x *ValidMeasurementUnit) Update(input *ValidMeasurementUnitUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.Volumetric != nil && *input.Volumetric != x.Volumetric {
		x.Volumetric = *input.Volumetric
	}

	if input.IconPath != nil && *input.IconPath != x.IconPath {
		x.IconPath = *input.IconPath
	}

	if input.Universal != nil && *input.Universal != x.Universal {
		x.Universal = *input.Universal
	}

	if input.Metric != nil && *input.Metric != x.Metric {
		x.Metric = *input.Metric
	}

	if input.Imperial != nil && *input.Imperial != x.Imperial {
		x.Imperial = *input.Imperial
	}

	if input.PluralName != nil && *input.PluralName != x.PluralName {
		x.PluralName = *input.PluralName
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
	}
}

var _ validation.ValidatableWithContext = (*ValidMeasurementUnitCreationRequestInput)(nil)

// ValidateWithContext validates a ValidMeasurementUnitCreationRequestInput.
func (x *ValidMeasurementUnitCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	var result *multierror.Error

	if (x.Metric && x.Imperial) || (!x.Metric && !x.Imperial) {
		result = multierror.Append(result, errMustBeEitherMetricOrImperial)
	}

	if err := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	); err != nil {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*ValidMeasurementUnitDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidMeasurementUnitDatabaseCreationInput.
func (x *ValidMeasurementUnitDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	var result *multierror.Error

	if (x.Metric && x.Imperial) || (!x.Metric && !x.Imperial) {
		result = multierror.Append(result, errMustBeEitherMetricOrImperial)
	}

	if err := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
	); err != nil {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*ValidMeasurementUnitUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidMeasurementUnitUpdateRequestInput.
func (x *ValidMeasurementUnitUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
