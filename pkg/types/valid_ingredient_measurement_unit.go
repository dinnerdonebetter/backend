package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidIngredientMeasurementUnitDataType indicates an event is related to a valid ingredient measurement unit.
	ValidIngredientMeasurementUnitDataType dataType = "valid_ingredient_measurement_unit"

	// ValidIngredientMeasurementUnitCreatedCustomerEventType indicates a valid ingredient measurement unit was created.
	ValidIngredientMeasurementUnitCreatedCustomerEventType CustomerEventType = "valid_ingredient_measurement_unit_created"
	// ValidIngredientMeasurementUnitUpdatedCustomerEventType indicates a valid ingredient measurement unit was updated.
	ValidIngredientMeasurementUnitUpdatedCustomerEventType CustomerEventType = "valid_ingredient_measurement_unit_updated"
	// ValidIngredientMeasurementUnitArchivedCustomerEventType indicates a valid ingredient measurement unit was archived.
	ValidIngredientMeasurementUnitArchivedCustomerEventType CustomerEventType = "valid_ingredient_measurement_unit_archived"
)

func init() {
	gob.Register(new(ValidIngredientMeasurementUnit))
	gob.Register(new(ValidIngredientMeasurementUnitCreationRequestInput))
	gob.Register(new(ValidIngredientMeasurementUnitUpdateRequestInput))
}

type (
	// ValidIngredientMeasurementUnit represents a valid ingredient measurement unit.
	ValidIngredientMeasurementUnit struct {
		_                        struct{}
		CreatedAt                time.Time            `json:"createdAt"`
		LastUpdatedAt            *time.Time           `json:"lastUpdatedAt"`
		ArchivedAt               *time.Time           `json:"archivedAt"`
		Notes                    string               `json:"notes"`
		ID                       string               `json:"id"`
		MeasurementUnit          ValidMeasurementUnit `json:"measurementUnit"`
		Ingredient               ValidIngredient      `json:"ingredient"`
		MinimumAllowableQuantity float32              `json:"minimumAllowableQuantity"`
		MaximumAllowableQuantity float32              `json:"maximumAllowableQuantity"`
	}

	// ValidIngredientMeasurementUnitCreationRequestInput represents what a user could set as input for creating valid ingredient measurement units.
	ValidIngredientMeasurementUnitCreationRequestInput struct {
		_ struct{}

		Notes                    string  `json:"notes"`
		ValidMeasurementUnitID   string  `json:"validMeasurementUnitID"`
		ValidIngredientID        string  `json:"validIngredientID"`
		MinimumAllowableQuantity float32 `json:"minimumAllowableQuantity"`
		MaximumAllowableQuantity float32 `json:"maximumAllowableQuantity"`
	}

	// ValidIngredientMeasurementUnitDatabaseCreationInput represents what a user could set as input for creating valid ingredient measurement units.
	ValidIngredientMeasurementUnitDatabaseCreationInput struct {
		_ struct{}

		ID                       string
		Notes                    string
		ValidMeasurementUnitID   string
		ValidIngredientID        string
		MinimumAllowableQuantity float32
		MaximumAllowableQuantity float32
	}

	// ValidIngredientMeasurementUnitUpdateRequestInput represents what a user could set as input for updating valid ingredient measurement units.
	ValidIngredientMeasurementUnitUpdateRequestInput struct {
		_ struct{}

		Notes                    *string  `json:"notes"`
		ValidMeasurementUnitID   *string  `json:"validMeasurementUnitID"`
		ValidIngredientID        *string  `json:"validIngredientID"`
		MinimumAllowableQuantity *float32 `json:"minimumAllowableQuantity"`
		MaximumAllowableQuantity *float32 `json:"maximumAllowableQuantity"`
	}

	// ValidIngredientMeasurementUnitDataManager describes a structure capable of storing valid ingredient measurement units permanently.
	ValidIngredientMeasurementUnitDataManager interface {
		ValidIngredientMeasurementUnitExists(ctx context.Context, validIngredientMeasurementUnitID string) (bool, error)
		GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*ValidIngredientMeasurementUnit, error)
		GetValidIngredientMeasurementUnits(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidIngredientMeasurementUnit], error)
		GetValidIngredientMeasurementUnitsForIngredient(ctx context.Context, ingredientID string, filter *QueryFilter) (*QueryFilteredResult[ValidIngredientMeasurementUnit], error)
		GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx context.Context, ingredientID string, filter *QueryFilter) (*QueryFilteredResult[ValidIngredientMeasurementUnit], error)
		CreateValidIngredientMeasurementUnit(ctx context.Context, input *ValidIngredientMeasurementUnitDatabaseCreationInput) (*ValidIngredientMeasurementUnit, error)
		UpdateValidIngredientMeasurementUnit(ctx context.Context, updated *ValidIngredientMeasurementUnit) error
		ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error
	}

	// ValidIngredientMeasurementUnitDataService describes a structure capable of serving traffic related to valid ingredient measurement units.
	ValidIngredientMeasurementUnitDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
		SearchByIngredientHandler(res http.ResponseWriter, req *http.Request)
		SearchByMeasurementUnitHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidIngredientMeasurementUnitUpdateRequestInput with a valid ingredient measurement unit.
func (x *ValidIngredientMeasurementUnit) Update(input *ValidIngredientMeasurementUnitUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ValidMeasurementUnitID != nil && *input.ValidMeasurementUnitID != x.MeasurementUnit.ID {
		x.MeasurementUnit.ID = *input.ValidMeasurementUnitID
	}

	if input.ValidIngredientID != nil && *input.ValidIngredientID != x.Ingredient.ID {
		x.Ingredient.ID = *input.ValidIngredientID
	}

	if input.MinimumAllowableQuantity != nil && *input.MinimumAllowableQuantity != x.MinimumAllowableQuantity {
		x.MinimumAllowableQuantity = *input.MinimumAllowableQuantity
	}

	if input.MaximumAllowableQuantity != nil && *input.MaximumAllowableQuantity != x.MaximumAllowableQuantity {
		x.MaximumAllowableQuantity = *input.MaximumAllowableQuantity
	}
}

var _ validation.ValidatableWithContext = (*ValidIngredientMeasurementUnitCreationRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientMeasurementUnitCreationRequestInput.
func (x *ValidIngredientMeasurementUnitCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.MinimumAllowableQuantity, validation.Min(0.009)),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientMeasurementUnitDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidIngredientMeasurementUnitDatabaseCreationInput.
func (x *ValidIngredientMeasurementUnitDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.MinimumAllowableQuantity, validation.Min(0.009)),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientMeasurementUnitUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientMeasurementUnitUpdateRequestInput.
func (x *ValidIngredientMeasurementUnitUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
	)
}
