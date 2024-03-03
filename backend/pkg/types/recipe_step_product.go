package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepProductIngredientType represents one of the valid recipe step product type values.
	RecipeStepProductIngredientType = "ingredient"
	// RecipeStepProductInstrumentType represents one of the valid recipe step product type values.
	RecipeStepProductInstrumentType = "instrument"
	// RecipeStepProductVesselType represents one of the valid recipe step product type values.
	RecipeStepProductVesselType = "vessel"

	// RecipeStepProductCreatedCustomerEventType indicates a recipe step product was created.
	RecipeStepProductCreatedCustomerEventType ServiceEventType = "recipe_step_product_created"
	// RecipeStepProductUpdatedCustomerEventType indicates a recipe step product was updated.
	RecipeStepProductUpdatedCustomerEventType ServiceEventType = "recipe_step_product_updated"
	// RecipeStepProductArchivedCustomerEventType indicates a recipe step product was archived.
	RecipeStepProductArchivedCustomerEventType ServiceEventType = "recipe_step_product_archived"
)

func init() {
	gob.Register(new(RecipeStepProduct))
	gob.Register(new(RecipeStepProductCreationRequestInput))
	gob.Register(new(RecipeStepProductUpdateRequestInput))
}

type (
	// RecipeStepProduct represents a recipe step product.
	RecipeStepProduct struct {
		_ struct{} `json:"-"`

		CreatedAt                          time.Time             `json:"createdAt"`
		MaximumStorageTemperatureInCelsius *float32              `json:"maximumStorageTemperatureInCelsius"`
		MaximumStorageDurationInSeconds    *uint32               `json:"maximumStorageDurationInSeconds"`
		MinimumStorageTemperatureInCelsius *float32              `json:"minimumStorageTemperatureInCelsius"`
		ArchivedAt                         *time.Time            `json:"archivedAt"`
		LastUpdatedAt                      *time.Time            `json:"lastUpdatedAt"`
		MinimumQuantity                    *float32              `json:"minimumQuantity"`
		MeasurementUnit                    *ValidMeasurementUnit `json:"measurementUnit"`
		MaximumQuantity                    *float32              `json:"maximumQuantity"`
		ContainedInVesselIndex             *uint16               `json:"containedInVesselIndex"`
		Name                               string                `json:"name"`
		BelongsToRecipeStep                string                `json:"belongsToRecipeStep"`
		Type                               string                `json:"type"`
		ID                                 string                `json:"id"`
		StorageInstructions                string                `json:"storageInstructions"`
		QuantityNotes                      string                `json:"quantityNotes"`
		Index                              uint16                `json:"index"`
		IsWaste                            bool                  `json:"isWaste"`
		IsLiquid                           bool                  `json:"isLiquid"`
		Compostable                        bool                  `json:"compostable"`
	}

	// RecipeStepProductCreationRequestInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductCreationRequestInput struct {
		_ struct{} `json:"-"`

		MinimumQuantity                    *float32 `json:"minimumQuantity"`
		MinimumStorageTemperatureInCelsius *float32 `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageDurationInSeconds    *uint32  `json:"maximumStorageDurationInSeconds"`
		MeasurementUnitID                  *string  `json:"measurementUnitID"`
		MaximumStorageTemperatureInCelsius *float32 `json:"maximumStorageTemperatureInCelsius"`
		MaximumQuantity                    *float32 `json:"maximumQuantity"`
		ContainedInVesselIndex             *uint16  `json:"containedInVesselIndex"`
		QuantityNotes                      string   `json:"quantityNotes"`
		Name                               string   `json:"name"`
		StorageInstructions                string   `json:"storageInstructions"`
		Type                               string   `json:"type"`
		Index                              uint16   `json:"index"`
		Compostable                        bool     `json:"compostable"`
		IsLiquid                           bool     `json:"isLiquid"`
		IsWaste                            bool     `json:"isWaste"`
	}

	// RecipeStepProductDatabaseCreationInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		MinimumQuantity                    *float32
		MinimumStorageTemperatureInCelsius *float32
		MaximumStorageTemperatureInCelsius *float32
		MaximumQuantity                    *float32
		MeasurementUnitID                  *string
		MaximumStorageDurationInSeconds    *uint32
		ContainedInVesselIndex             *uint16
		Name                               string
		BelongsToRecipeStep                string
		StorageInstructions                string
		QuantityNotes                      string
		ID                                 string
		Type                               string
		Index                              uint16
		Compostable                        bool
		IsLiquid                           bool
		IsWaste                            bool
	}

	// RecipeStepProductUpdateRequestInput represents what a user could set as input for updating recipe step products.
	RecipeStepProductUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name                               *string  `json:"name,omitempty"`
		Type                               *string  `json:"type,omitempty"`
		MeasurementUnitID                  *string  `json:"measurementUnitID,omitempty"`
		QuantityNotes                      *string  `json:"quantityNotes,omitempty"`
		BelongsToRecipeStep                *string  `json:"belongsToRecipeStep,omitempty"`
		MinimumQuantity                    *float32 `json:"minimumQuantity,omitempty"`
		MaximumQuantity                    *float32 `json:"maximumQuantity,omitempty"`
		Compostable                        *bool    `json:"compostable,omitempty"`
		MaximumStorageDurationInSeconds    *uint32  `json:"maximumStorageDurationInSeconds,omitempty"`
		MinimumStorageTemperatureInCelsius *float32 `json:"minimumStorageTemperatureInCelsius,omitempty"`
		MaximumStorageTemperatureInCelsius *float32 `json:"maximumStorageTemperatureInCelsius,omitempty"`
		StorageInstructions                *string  `json:"storageInstructions,omitempty"`
		IsLiquid                           *bool    `json:"isLiquid,omitempty"`
		IsWaste                            *bool    `json:"isWaste,omitempty"`
		Index                              *uint16  `json:"index,omitempty"`
		ContainedInVesselIndex             *uint16  `json:"containedInVesselIndex,omitempty"`
	}

	// RecipeStepProductDataManager describes a structure capable of storing recipe step products permanently.
	RecipeStepProductDataManager interface {
		RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (bool, error)
		GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*RecipeStepProduct, error)
		GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*QueryFilteredResult[RecipeStepProduct], error)
		CreateRecipeStepProduct(ctx context.Context, input *RecipeStepProductDatabaseCreationInput) (*RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, updated *RecipeStepProduct) error
		ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error
	}

	// RecipeStepProductDataService describes a structure capable of serving traffic related to recipe step products.
	RecipeStepProductDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipeStepProductUpdateRequestInput with a recipe step product.
func (x *RecipeStepProduct) Update(input *RecipeStepProductUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Type != nil && *input.Type != x.Type {
		x.Type = *input.Type
	}

	if input.MeasurementUnitID != nil && x.MeasurementUnit != nil && *input.MeasurementUnitID != x.MeasurementUnit.ID {
		x.MeasurementUnit = &ValidMeasurementUnit{ID: *input.MeasurementUnitID}
	}

	if input.MinimumQuantity != nil && input.MinimumQuantity != x.MinimumQuantity {
		x.MinimumQuantity = input.MinimumQuantity
	}

	if input.MaximumQuantity != nil && input.MaximumQuantity != x.MaximumQuantity {
		x.MaximumQuantity = input.MaximumQuantity
	}

	if input.QuantityNotes != nil && *input.QuantityNotes != x.QuantityNotes {
		x.QuantityNotes = *input.QuantityNotes
	}

	if input.Compostable != nil && *input.Compostable != x.Compostable {
		x.Compostable = *input.Compostable
	}

	if input.MaximumStorageDurationInSeconds != nil && input.MaximumStorageDurationInSeconds != x.MaximumStorageDurationInSeconds {
		x.MaximumStorageDurationInSeconds = input.MaximumStorageDurationInSeconds
	}

	if input.MinimumStorageTemperatureInCelsius != nil && input.MinimumStorageTemperatureInCelsius != x.MinimumStorageTemperatureInCelsius {
		x.MinimumStorageTemperatureInCelsius = input.MinimumStorageTemperatureInCelsius
	}

	if input.MaximumStorageTemperatureInCelsius != nil && input.MaximumStorageTemperatureInCelsius != x.MaximumStorageTemperatureInCelsius {
		x.MaximumStorageTemperatureInCelsius = input.MaximumStorageTemperatureInCelsius
	}

	if input.StorageInstructions != nil && *input.StorageInstructions != x.StorageInstructions {
		x.StorageInstructions = *input.StorageInstructions
	}

	if input.IsLiquid != nil && *input.IsLiquid != x.IsLiquid {
		x.IsLiquid = *input.IsLiquid
	}

	if input.IsWaste != nil && *input.IsWaste != x.IsWaste {
		x.IsWaste = *input.IsWaste
	}

	if input.Index != nil && *input.Index != x.Index {
		x.Index = *input.Index
	}

	if input.ContainedInVesselIndex != nil && x.ContainedInVesselIndex != nil && *input.ContainedInVesselIndex != *x.ContainedInVesselIndex {
		x.ContainedInVesselIndex = input.ContainedInVesselIndex
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepProductCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepProductCreationRequestInput.
func (x *RecipeStepProductCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Type, validation.In(RecipeStepProductIngredientType, RecipeStepProductInstrumentType, RecipeStepProductVesselType)),
		validation.Field(&x.MinimumQuantity, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepProductDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepProductDatabaseCreationInput.
func (x *RecipeStepProductDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
		validation.Field(&x.Type, validation.In(RecipeStepProductIngredientType, RecipeStepProductInstrumentType, RecipeStepProductVesselType)),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepProductUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepProductUpdateRequestInput.
func (x *RecipeStepProductUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Type, validation.In(RecipeStepProductIngredientType, RecipeStepProductInstrumentType, RecipeStepProductVesselType)),
		validation.Field(&x.MeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantity, validation.Required),
		validation.Field(&x.MaximumQuantity, validation.Required),
	)
}
