package types

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/segmentio/ksuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidMeasurementUnitDataType indicates an event is related to a valid measurement unit.
	ValidMeasurementUnitDataType dataType = "valid_measurement_unit"

	// ValidMeasurementUnitCreatedCustomerEventType indicates a valid measurement unit was created.
	ValidMeasurementUnitCreatedCustomerEventType CustomerEventType = "valid_measurement_unit_created"
	// ValidMeasurementUnitUpdatedCustomerEventType indicates a valid measurement unit was updated.
	ValidMeasurementUnitUpdatedCustomerEventType CustomerEventType = "valid_measurement_unit_updated"
	// ValidMeasurementUnitArchivedCustomerEventType indicates a valid measurement unit was archived.
	ValidMeasurementUnitArchivedCustomerEventType CustomerEventType = "valid_measurement_unit_archived"
)

func init() {
	gob.Register(new(ValidMeasurementUnit))
	gob.Register(new(ValidMeasurementUnitList))
	gob.Register(new(ValidMeasurementUnitCreationRequestInput))
	gob.Register(new(ValidMeasurementUnitUpdateRequestInput))
}

type (
	// ValidMeasurementUnit represents a valid measurement unit.
	ValidMeasurementUnit struct {
		_             struct{}
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		ArchivedOn    *uint64 `json:"archivedOn"`
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		ID            string  `json:"id"`
		IconPath      string  `json:"iconPath"`
		CreatedOn     uint64  `json:"createdOn"`
		Volumetric    bool    `json:"volumetric"`
	}

	// ValidMeasurementUnitList represents a list of valid measurement units.
	ValidMeasurementUnitList struct {
		_                     struct{}
		ValidMeasurementUnits []*ValidMeasurementUnit `json:"data"`
		Pagination
	}

	// ValidMeasurementUnitCreationRequestInput represents what a user could set as input for creating valid measurement units.
	ValidMeasurementUnitCreationRequestInput struct {
		_           struct{}
		Name        string `json:"name"`
		Description string `json:"description"`
		IconPath    string `json:"iconPath"`
		CreatedOn   uint64 `json:"createdOn"`
		Volumetric  bool   `json:"volumetric"`
	}

	// ValidMeasurementUnitDatabaseCreationInput represents what a user could set as input for creating valid measurement units.
	ValidMeasurementUnitDatabaseCreationInput struct {
		_           struct{}
		Name        string `json:"name"`
		Description string `json:"description"`
		ID          string `json:"id"`
		IconPath    string `json:"iconPath"`
		CreatedOn   uint64 `json:"createdOn"`
		Volumetric  bool   `json:"volumetric"`
	}

	// ValidMeasurementUnitUpdateRequestInput represents what a user could set as input for updating valid measurement units.
	ValidMeasurementUnitUpdateRequestInput struct {
		_           struct{}
		Name        *string `json:"name"`
		Description *string `json:"description"`
		IconPath    *string `json:"iconPath"`
		CreatedOn   *uint64 `json:"createdOn"`
		Volumetric  *bool   `json:"volumetric"`
	}

	// ValidMeasurementUnitDataManager describes a structure capable of storing valid measurement units permanently.
	ValidMeasurementUnitDataManager interface {
		ValidMeasurementUnitExists(ctx context.Context, validMeasurementUnitID string) (bool, error)
		GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*ValidMeasurementUnit, error)
		GetValidMeasurementUnits(ctx context.Context, filter *QueryFilter) (*ValidMeasurementUnitList, error)
		SearchForValidMeasurementUnits(ctx context.Context, query string) ([]*ValidMeasurementUnit, error)
		CreateValidMeasurementUnit(ctx context.Context, input *ValidMeasurementUnitDatabaseCreationInput) (*ValidMeasurementUnit, error)
		UpdateValidMeasurementUnit(ctx context.Context, updated *ValidMeasurementUnit) error
		ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error
	}

	// ValidMeasurementUnitDataService describes a structure capable of serving traffic related to valid measurement units.
	ValidMeasurementUnitDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
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
}

var _ validation.ValidatableWithContext = (*ValidMeasurementUnitCreationRequestInput)(nil)

// ValidateWithContext validates a ValidMeasurementUnitCreationRequestInput.
func (x *ValidMeasurementUnitCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidMeasurementUnitDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidMeasurementUnitDatabaseCreationInput.
func (x *ValidMeasurementUnitDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
	)
}

// ValidMeasurementUnitUpdateRequestInputFromValidMeasurementUnit creates a DatabaseCreationInput from a CreationInput.
func ValidMeasurementUnitUpdateRequestInputFromValidMeasurementUnit(input *ValidMeasurementUnit) *ValidMeasurementUnitUpdateRequestInput {
	x := &ValidMeasurementUnitUpdateRequestInput{
		Name:        &input.Name,
		Description: &input.Description,
		IconPath:    &input.IconPath,
		Volumetric:  &input.Volumetric,
	}

	return x
}

// ValidMeasurementUnitDatabaseCreationInputFromValidMeasurementUnitCreationInput creates a DatabaseCreationInput from a CreationInput.
func ValidMeasurementUnitDatabaseCreationInputFromValidMeasurementUnitCreationInput(input *ValidMeasurementUnitCreationRequestInput) *ValidMeasurementUnitDatabaseCreationInput {
	x := &ValidMeasurementUnitDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		Volumetric:  input.Volumetric,
		IconPath:    input.IconPath,
	}

	return x
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