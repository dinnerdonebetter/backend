package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// UserIngredientPreferenceCreatedCustomerEventType indicates a user ingredient preference was created.
	UserIngredientPreferenceCreatedCustomerEventType ServiceEventType = "user_ingredient_preference_created"
	// UserIngredientPreferenceUpdatedCustomerEventType indicates a user ingredient preference was updated.
	UserIngredientPreferenceUpdatedCustomerEventType ServiceEventType = "user_ingredient_preference_updated"
	// UserIngredientPreferenceArchivedCustomerEventType indicates a user ingredient preference was archived.
	UserIngredientPreferenceArchivedCustomerEventType ServiceEventType = "user_ingredient_preference_archived"

	minRating int8 = -10
	maxRating int8 = 10
)

func init() {
	gob.Register(new(UserIngredientPreference))
	gob.Register(new(UserIngredientPreferenceCreationRequestInput))
	gob.Register(new(UserIngredientPreferenceUpdateRequestInput))
}

type (
	// UserIngredientPreference represents a user ingredient preference.
	UserIngredientPreference struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time       `json:"createdAt"`
		LastUpdatedAt *time.Time      `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time      `json:"archivedAt"`
		ID            string          `json:"id"`
		Notes         string          `json:"notes"`
		BelongsToUser string          `json:"belongsToUser"`
		Ingredient    ValidIngredient `json:"ingredient"`
		Rating        int8            `json:"rating"`
		Allergy       bool            `json:"allergy"`
	}

	// UserIngredientPreferenceCreationRequestInput represents what a user could set as input for creating user ingredient preferences.
	UserIngredientPreferenceCreationRequestInput struct {
		_ struct{} `json:"-"`

		ValidIngredientGroupID string `json:"validIngredientGroupID"`
		ValidIngredientID      string `json:"validIngredientID"`
		Notes                  string `json:"notes"`
		Rating                 int8   `json:"rating"`
		Allergy                bool   `json:"allergy"`
	}

	// UserIngredientPreferenceDatabaseCreationInput represents what a user could set as input for creating user ingredient preferences.
	UserIngredientPreferenceDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ValidIngredientGroupID string
		ValidIngredientID      string
		Notes                  string
		BelongsToUser          string
		Rating                 int8
		Allergy                bool
	}

	// UserIngredientPreferenceUpdateRequestInput represents what a user could set as input for updating user ingredient preferences.
	UserIngredientPreferenceUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes        *string `json:"notes,omitempty"`
		IngredientID *string `json:"ingredientID"`
		Rating       *int8   `json:"rating"`
		Allergy      *bool   `json:"allergy"`
	}

	// UserIngredientPreferenceDataManager describes a structure capable of storing user ingredient preferences permanently.
	UserIngredientPreferenceDataManager interface {
		UserIngredientPreferenceExists(ctx context.Context, userIngredientPreferenceID, userID string) (bool, error)
		GetUserIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) (*UserIngredientPreference, error)
		GetUserIngredientPreferences(ctx context.Context, userID string, filter *QueryFilter) (*QueryFilteredResult[UserIngredientPreference], error)
		CreateUserIngredientPreference(ctx context.Context, input *UserIngredientPreferenceDatabaseCreationInput) ([]*UserIngredientPreference, error)
		UpdateUserIngredientPreference(ctx context.Context, updated *UserIngredientPreference) error
		ArchiveUserIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) error
	}

	// UserIngredientPreferenceDataService describes a structure capable of serving traffic related to user ingredient preferences.
	UserIngredientPreferenceDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an UserIngredientPreferenceUpdateRequestInput with a user ingredient preference.
func (x *UserIngredientPreference) Update(input *UserIngredientPreferenceUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.IngredientID != nil && *input.IngredientID != x.Ingredient.ID {
		x.Ingredient = ValidIngredient{ID: *input.IngredientID}
	}

	if input.Rating != nil && *input.Rating != x.Rating {
		x.Rating = *input.Rating
	}

	if input.Allergy != nil && *input.Allergy != x.Allergy {
		x.Allergy = *input.Allergy
	}
}

var _ validation.ValidatableWithContext = (*UserIngredientPreferenceCreationRequestInput)(nil)

// ValidateWithContext validates a UserIngredientPreferenceCreationRequestInput.
func (x *UserIngredientPreferenceCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidIngredientGroupID, validation.When(x.ValidIngredientID == "", validation.Required)),
		validation.Field(&x.ValidIngredientID, validation.When(x.ValidIngredientGroupID == "", validation.Required)),
		validation.Field(&x.Rating, validation.Min(minRating), validation.Max(maxRating)),
	)
}

var _ validation.ValidatableWithContext = (*UserIngredientPreferenceDatabaseCreationInput)(nil)

// ValidateWithContext validates a UserIngredientPreferenceDatabaseCreationInput.
func (x *UserIngredientPreferenceDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.Rating, validation.Min(minRating), validation.Max(maxRating)),
		validation.Field(&x.BelongsToUser, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*UserIngredientPreferenceUpdateRequestInput)(nil)

// ValidateWithContext validates a UserIngredientPreferenceUpdateRequestInput.
func (x *UserIngredientPreferenceUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientID, validation.Required),
		validation.Field(&x.Rating, validation.Min(minRating), validation.Max(maxRating)),
	)
}
