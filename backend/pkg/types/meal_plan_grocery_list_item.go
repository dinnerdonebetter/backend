package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanGroceryListItemStatusUnknown represents the database-side enum member for grocery list item status.
	MealPlanGroceryListItemStatusUnknown = "unknown"
	// MealPlanGroceryListItemStatusAlreadyOwned represents the database-side enum member for grocery list item status.
	MealPlanGroceryListItemStatusAlreadyOwned = "already owned"
	// MealPlanGroceryListItemStatusNeeds represents the database-side enum member for grocery list item status.
	MealPlanGroceryListItemStatusNeeds = "needs"
	// MealPlanGroceryListItemStatusUnavailable represents the database-side enum member for grocery list item status.
	MealPlanGroceryListItemStatusUnavailable = "unavailable"
	// MealPlanGroceryListItemStatusAcquired represents the database-side enum member for grocery list item status.
	MealPlanGroceryListItemStatusAcquired = "acquired"

	// MealPlanGroceryListItemCreatedCustomerEventType indicates a meal plan grocery list item was created.
	MealPlanGroceryListItemCreatedCustomerEventType ServiceEventType = "meal_plan_grocery_list_item_created"
	// MealPlanGroceryListItemUpdatedCustomerEventType indicates a meal plan grocery list item was updated.
	MealPlanGroceryListItemUpdatedCustomerEventType ServiceEventType = "meal_plan_grocery_list_item_updated"
	// MealPlanGroceryListItemArchivedCustomerEventType indicates a meal plan grocery list item was archived.
	MealPlanGroceryListItemArchivedCustomerEventType ServiceEventType = "meal_plan_grocery_list_item_archived"
)

func init() {
	gob.Register(new(MealPlanGroceryListItem))
	gob.Register(new(MealPlanGroceryListItemCreationRequestInput))
	gob.Register(new(MealPlanGroceryListItemUpdateRequestInput))
}

type (
	// MealPlanGroceryListItem represents a meal plan grocery list item.
	MealPlanGroceryListItem struct {
		_ struct{} `json:"-"`

		CreatedAt                time.Time             `json:"createdAt"`
		MaximumQuantityNeeded    *float32              `json:"maximumQuantityNeeded"`
		LastUpdatedAt            *time.Time            `json:"lastUpdatedAt"`
		PurchasePrice            *float32              `json:"purchasePrice"`
		PurchasedUPC             *string               `json:"purchasedUPC"`
		ArchivedAt               *time.Time            `json:"archivedAt"`
		QuantityPurchased        *float32              `json:"quantityPurchased"`
		PurchasedMeasurementUnit *ValidMeasurementUnit `json:"purchasedMeasurementUnit"`
		BelongsToMealPlan        string                `json:"belongsToMealPlan"`
		Status                   string                `json:"status"`
		StatusExplanation        string                `json:"statusExplanation"`
		ID                       string                `json:"id"`
		MeasurementUnit          ValidMeasurementUnit  `json:"measurementUnit"`
		Ingredient               ValidIngredient       `json:"ingredient"`
		MinimumQuantityNeeded    float32               `json:"minimumQuantityNeeded"`
	}

	// MealPlanGroceryListItemCreationRequestInput represents what a user could set as input for creating meal plan grocery list items.
	MealPlanGroceryListItemCreationRequestInput struct {
		_ struct{} `json:"-"`

		PurchasedMeasurementUnitID *string  `json:"purchasedMeasurementUnitID"`
		PurchasedUPC               *string  `json:"purchasedUPC"`
		PurchasePrice              *float32 `json:"purchasePrice"`
		QuantityPurchased          *float32 `json:"quantityPurchased"`
		MaximumQuantityNeeded      *float32 `json:"maximumQuantityNeeded"`
		Status                     string   `json:"status"`
		BelongsToMealPlan          string   `json:"belongsToMealPlan"`
		ValidIngredientID          string   `json:"validIngredientID"`
		ValidMeasurementUnitID     string   `json:"validMeasurementUnitID"`
		StatusExplanation          string   `json:"statusExplanation"`
		MinimumQuantityNeeded      float32  `json:"minimumQuantityNeeded"`
	}

	// MealPlanGroceryListItemDatabaseCreationInput represents what a user could set as input for creating meal plan grocery list items.
	MealPlanGroceryListItemDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		PurchasePrice              *float32
		PurchasedUPC               *string
		PurchasedMeasurementUnitID *string
		QuantityPurchased          *float32
		MaximumQuantityNeeded      *float32
		Status                     string
		ValidMeasurementUnitID     string
		ValidIngredientID          string
		BelongsToMealPlan          string
		ID                         string
		StatusExplanation          string
		MinimumQuantityNeeded      float32
	}

	// MealPlanGroceryListItemUpdateRequestInput represents what a user could set as input for updating meal plan grocery list items.
	MealPlanGroceryListItemUpdateRequestInput struct {
		_ struct{} `json:"-"`

		MaximumQuantityNeeded      *float32 `json:"maximumQuantityNeeded,omitempty"`
		BelongsToMealPlan          *string  `json:"belongsToMealPlan,omitempty"`
		ValidIngredientID          *string  `json:"validIngredientID,omitempty"`
		ValidMeasurementUnitID     *string  `json:"validMeasurementUnitID,omitempty"`
		MinimumQuantityNeeded      *float32 `json:"minimumQuantityNeeded,omitempty"`
		StatusExplanation          *string  `json:"statusExplanation,omitempty"`
		QuantityPurchased          *float32 `json:"quantityPurchased,omitempty"`
		PurchasedMeasurementUnitID *string  `json:"purchasedMeasurementUnitID,omitempty"`
		PurchasedUPC               *string  `json:"purchasedUPC,omitempty"`
		PurchasePrice              *float32 `json:"purchasePrice,omitempty"`
		Status                     *string  `json:"status,omitempty"`
	}

	// MealPlanGroceryListItemDataManager describes a structure capable of storing meal plan grocery list items permanently.
	MealPlanGroceryListItemDataManager interface {
		MealPlanGroceryListItemExists(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (bool, error)
		GetMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*MealPlanGroceryListItem, error)
		GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string) ([]*MealPlanGroceryListItem, error)
		CreateMealPlanGroceryListItem(ctx context.Context, input *MealPlanGroceryListItemDatabaseCreationInput) (*MealPlanGroceryListItem, error)
		UpdateMealPlanGroceryListItem(ctx context.Context, updated *MealPlanGroceryListItem) error
		ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItemID string) error
	}

	// MealPlanGroceryListItemDataService describes a structure capable of serving traffic related to meal plan grocery list items.
	MealPlanGroceryListItemDataService interface {
		ListByMealPlanHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an MealPlanGroceryListItemUpdateRequestInput with a meal plan grocery list item.
func (x *MealPlanGroceryListItem) Update(input *MealPlanGroceryListItemUpdateRequestInput) {
	if input.BelongsToMealPlan != nil && *input.BelongsToMealPlan != x.BelongsToMealPlan {
		x.BelongsToMealPlan = *input.BelongsToMealPlan
	}

	if input.ValidIngredientID != nil && *input.ValidIngredientID != x.Ingredient.ID {
		x.Ingredient = ValidIngredient{ID: *input.ValidIngredientID}
	}

	if input.ValidMeasurementUnitID != nil && *input.ValidMeasurementUnitID != x.MeasurementUnit.ID {
		x.MeasurementUnit = ValidMeasurementUnit{ID: *input.ValidMeasurementUnitID}
	}

	if input.MinimumQuantityNeeded != nil && *input.MinimumQuantityNeeded != x.MinimumQuantityNeeded {
		x.MinimumQuantityNeeded = *input.MinimumQuantityNeeded
	}

	if input.MaximumQuantityNeeded != nil && x.MaximumQuantityNeeded != nil && *input.MaximumQuantityNeeded != *x.MaximumQuantityNeeded {
		x.MaximumQuantityNeeded = input.MaximumQuantityNeeded
	}

	// was nil and now will not be nil
	if input.QuantityPurchased != nil && ((x.QuantityPurchased == nil) || (x.QuantityPurchased != nil && (*input.QuantityPurchased != *x.QuantityPurchased))) {
		x.QuantityPurchased = input.QuantityPurchased
	}

	if input.PurchasedMeasurementUnitID != nil && *input.PurchasedMeasurementUnitID != x.PurchasedMeasurementUnit.ID {
		x.PurchasedMeasurementUnit = &ValidMeasurementUnit{ID: *input.PurchasedMeasurementUnitID}
	}

	if input.PurchasedUPC != nil && ((x.PurchasedUPC == nil) || (x.PurchasedUPC != nil && (*input.PurchasedUPC != *x.PurchasedUPC))) {
		x.PurchasedUPC = input.PurchasedUPC
	}

	if input.PurchasePrice != nil && ((x.PurchasePrice == nil) || (x.PurchasePrice != nil && (*input.PurchasePrice != *x.PurchasePrice))) {
		x.PurchasePrice = input.PurchasePrice
	}

	if input.StatusExplanation != nil && *input.StatusExplanation != x.StatusExplanation {
		x.StatusExplanation = *input.StatusExplanation
	}
	if input.Status != nil && *input.Status != x.Status {
		x.Status = *input.Status
	}
}

var _ validation.ValidatableWithContext = (*MealPlanGroceryListItemCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanGroceryListItemCreationRequestInput.
func (x *MealPlanGroceryListItemCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToMealPlan, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantityNeeded, validation.Required),
		validation.Field(&x.Status, validation.In(
			MealPlanGroceryListItemStatusUnknown,
			MealPlanGroceryListItemStatusAlreadyOwned,
			MealPlanGroceryListItemStatusNeeds,
			MealPlanGroceryListItemStatusUnavailable,
			MealPlanGroceryListItemStatusAcquired,
		)),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanGroceryListItemDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanGroceryListItemDatabaseCreationInput.
func (x *MealPlanGroceryListItemDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToMealPlan, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantityNeeded, validation.Required),
		validation.Field(&x.Status, validation.In(
			MealPlanGroceryListItemStatusUnknown,
			MealPlanGroceryListItemStatusAlreadyOwned,
			MealPlanGroceryListItemStatusNeeds,
			MealPlanGroceryListItemStatusUnavailable,
			MealPlanGroceryListItemStatusAcquired,
		)),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanGroceryListItemUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanGroceryListItemUpdateRequestInput.
func (x *MealPlanGroceryListItemUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Status, validation.In(
			MealPlanGroceryListItemStatusUnknown,
			MealPlanGroceryListItemStatusAlreadyOwned,
			MealPlanGroceryListItemStatusNeeds,
			MealPlanGroceryListItemStatusUnavailable,
			MealPlanGroceryListItemStatusAcquired,
		)),
	)
}
