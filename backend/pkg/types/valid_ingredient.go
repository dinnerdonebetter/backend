package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidIngredientCreatedCustomerEventType indicates a valid ingredient was created.
	ValidIngredientCreatedCustomerEventType ServiceEventType = "valid_ingredient_created"
	// ValidIngredientUpdatedCustomerEventType indicates a valid ingredient was updated.
	ValidIngredientUpdatedCustomerEventType ServiceEventType = "valid_ingredient_updated"
	// ValidIngredientArchivedCustomerEventType indicates a valid ingredient was archived.
	ValidIngredientArchivedCustomerEventType ServiceEventType = "valid_ingredient_archived"
)

func init() {
	gob.Register(new(ValidIngredient))
	gob.Register(new(ValidIngredientCreationRequestInput))
	gob.Register(new(ValidIngredientUpdateRequestInput))
}

type (
	// ValidIngredient represents a valid ingredient.
	ValidIngredient struct {
		_ struct{} `json:"-"`

		CreatedAt                               time.Time  `json:"createdAt"`
		LastUpdatedAt                           *time.Time `json:"lastUpdatedAt"`
		ArchivedAt                              *time.Time `json:"archivedAt"`
		MaximumIdealStorageTemperatureInCelsius *float32   `json:"maximumIdealStorageTemperatureInCelsius"`
		MinimumIdealStorageTemperatureInCelsius *float32   `json:"minimumIdealStorageTemperatureInCelsius"`
		IconPath                                string     `json:"iconPath"`
		Warning                                 string     `json:"warning"`
		PluralName                              string     `json:"pluralName"`
		StorageInstructions                     string     `json:"storageInstructions"`
		Name                                    string     `json:"name"`
		ID                                      string     `json:"id"`
		Description                             string     `json:"description"`
		Slug                                    string     `json:"slug"`
		ShoppingSuggestions                     string     `json:"shoppingSuggestions"`
		ContainsShellfish                       bool       `json:"containsShellfish"`
		IsMeasuredVolumetrically                bool       `json:"isMeasuredVolumetrically"`
		IsLiquid                                bool       `json:"isLiquid"`
		ContainsPeanut                          bool       `json:"containsPeanut"`
		ContainsTreeNut                         bool       `json:"containsTreeNut"`
		ContainsEgg                             bool       `json:"containsEgg"`
		ContainsWheat                           bool       `json:"containsWheat"`
		ContainsSoy                             bool       `json:"containsSoy"`
		AnimalDerived                           bool       `json:"animalDerived"`
		RestrictToPreparations                  bool       `json:"restrictToPreparations"`
		ContainsSesame                          bool       `json:"containsSesame"`
		ContainsFish                            bool       `json:"containsFish"`
		ContainsGluten                          bool       `json:"containsGluten"`
		ContainsDairy                           bool       `json:"containsDairy"`
		ContainsAlcohol                         bool       `json:"containsAlcohol"`
		AnimalFlesh                             bool       `json:"animalFlesh"`
		IsStarch                                bool       `json:"is_starch"`
		IsProtein                               bool       `json:"is_protein"`
		IsGrain                                 bool       `json:"is_grain"`
		IsFruit                                 bool       `json:"is_fruit"`
		IsSalt                                  bool       `json:"is_salt"`
		IsFat                                   bool       `json:"is_fat"`
		IsAcid                                  bool       `json:"is_acid"`
		IsHeat                                  bool       `json:"is_heat"`
	}

	// NullableValidIngredient represents a nullable valid ingredient.
	NullableValidIngredient struct {
		_ struct{} `json:"-"`

		CreatedAt                               *time.Time
		LastUpdatedAt                           *time.Time
		ArchivedAt                              *time.Time
		ID                                      *string
		Warning                                 *string
		Description                             *string
		IconPath                                *string
		PluralName                              *string
		StorageInstructions                     *string
		Name                                    *string
		MaximumIdealStorageTemperatureInCelsius *float32
		MinimumIdealStorageTemperatureInCelsius *float32
		ContainsShellfish                       *bool
		ContainsDairy                           *bool
		AnimalFlesh                             *bool
		IsMeasuredVolumetrically                *bool
		IsLiquid                                *bool
		ContainsPeanut                          *bool
		ContainsTreeNut                         *bool
		ContainsEgg                             *bool
		ContainsWheat                           *bool
		ContainsSoy                             *bool
		AnimalDerived                           *bool
		RestrictToPreparations                  *bool
		ContainsSesame                          *bool
		ContainsFish                            *bool
		ContainsGluten                          *bool
		Slug                                    *string
		ContainsAlcohol                         *bool
		ShoppingSuggestions                     *string
		IsStarch                                *bool
		IsProtein                               *bool
		IsGrain                                 *bool
		IsFruit                                 *bool
		IsSalt                                  *bool
		IsFat                                   *bool
		IsAcid                                  *bool
		IsHeat                                  *bool
	}

	// ValidIngredientCreationRequestInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientCreationRequestInput struct {
		_ struct{} `json:"-"`

		MinimumIdealStorageTemperatureInCelsius *float32 `json:"minimumIdealStorageTemperatureInCelsius"`
		MaximumIdealStorageTemperatureInCelsius *float32 `json:"maximumIdealStorageTemperatureInCelsius"`
		Warning                                 string   `json:"warning"`
		IconPath                                string   `json:"iconPath"`
		PluralName                              string   `json:"pluralName"`
		StorageInstructions                     string   `json:"storageInstructions"`
		Name                                    string   `json:"name"`
		Description                             string   `json:"description"`
		Slug                                    string   `json:"slug"`
		ShoppingSuggestions                     string   `json:"shoppingSuggestions"`
		IsMeasuredVolumetrically                bool     `json:"isMeasuredVolumetrically"`
		ContainsFish                            bool     `json:"containsFish"`
		ContainsShellfish                       bool     `json:"containsShellfish"`
		AnimalFlesh                             bool     `json:"animalFlesh"`
		ContainsEgg                             bool     `json:"containsEgg"`
		IsLiquid                                bool     `json:"isLiquid"`
		ContainsSoy                             bool     `json:"containsSoy"`
		ContainsPeanut                          bool     `json:"containsPeanut"`
		AnimalDerived                           bool     `json:"animalDerived"`
		RestrictToPreparations                  bool     `json:"restrictToPreparations"`
		ContainsDairy                           bool     `json:"containsDairy"`
		ContainsSesame                          bool     `json:"containsSesame"`
		ContainsTreeNut                         bool     `json:"containsTreeNut"`
		ContainsWheat                           bool     `json:"containsWheat"`
		ContainsAlcohol                         bool     `json:"containsAlcohol"`
		ContainsGluten                          bool     `json:"containsGluten"`
		IsStarch                                bool     `json:"is_starch"`
		IsProtein                               bool     `json:"is_protein"`
		IsGrain                                 bool     `json:"is_grain"`
		IsFruit                                 bool     `json:"is_fruit"`
		IsSalt                                  bool     `json:"is_salt"`
		IsFat                                   bool     `json:"is_fat"`
		IsAcid                                  bool     `json:"is_acid"`
		IsHeat                                  bool     `json:"is_heat"`
	}

	// ValidIngredientDatabaseCreationInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		MinimumIdealStorageTemperatureInCelsius *float32
		MaximumIdealStorageTemperatureInCelsius *float32
		ID                                      string
		Warning                                 string
		IconPath                                string
		PluralName                              string
		StorageInstructions                     string
		Name                                    string
		Description                             string
		Slug                                    string
		ShoppingSuggestions                     string
		IsMeasuredVolumetrically                bool
		ContainsFish                            bool
		ContainsShellfish                       bool
		AnimalFlesh                             bool
		ContainsEgg                             bool
		IsLiquid                                bool
		ContainsSoy                             bool
		ContainsPeanut                          bool
		AnimalDerived                           bool
		RestrictToPreparations                  bool
		ContainsDairy                           bool
		ContainsSesame                          bool
		ContainsTreeNut                         bool
		ContainsWheat                           bool
		ContainsAlcohol                         bool
		ContainsGluten                          bool
		IsStarch                                bool
		IsProtein                               bool
		IsGrain                                 bool
		IsFruit                                 bool
		IsSalt                                  bool
		IsFat                                   bool
		IsAcid                                  bool
		IsHeat                                  bool
	}

	// ValidIngredientUpdateRequestInput represents what a user could set as input for updating valid ingredients.
	ValidIngredientUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name                                    *string  `json:"name,omitempty"`
		Description                             *string  `json:"description,omitempty"`
		Warning                                 *string  `json:"warning,omitempty"`
		IconPath                                *string  `json:"iconPath,omitempty"`
		ContainsDairy                           *bool    `json:"containsDairy,omitempty"`
		ContainsPeanut                          *bool    `json:"containsPeanut,omitempty"`
		ContainsTreeNut                         *bool    `json:"containsTreeNut,omitempty"`
		ContainsEgg                             *bool    `json:"containsEgg,omitempty"`
		ContainsWheat                           *bool    `json:"containsWheat,omitempty"`
		ContainsShellfish                       *bool    `json:"containsShellfish,omitempty"`
		ContainsSesame                          *bool    `json:"containsSesame,omitempty"`
		ContainsFish                            *bool    `json:"containsFish,omitempty"`
		ContainsGluten                          *bool    `json:"containsGluten,omitempty"`
		AnimalFlesh                             *bool    `json:"animalFlesh,omitempty"`
		IsMeasuredVolumetrically                *bool    `json:"isMeasuredVolumetrically,omitempty"`
		IsLiquid                                *bool    `json:"isLiquid,omitempty"`
		ContainsSoy                             *bool    `json:"containsSoy,omitempty"`
		PluralName                              *string  `json:"pluralName,omitempty"`
		AnimalDerived                           *bool    `json:"animalDerived,omitempty"`
		RestrictToPreparations                  *bool    `json:"restrictToPreparations,omitempty"`
		MinimumIdealStorageTemperatureInCelsius *float32 `json:"minimumIdealStorageTemperatureInCelsius,omitempty"`
		MaximumIdealStorageTemperatureInCelsius *float32 `json:"maximumIdealStorageTemperatureInCelsius,omitempty"`
		StorageInstructions                     *string  `json:"storageInstructions,omitempty"`
		Slug                                    *string  `json:"slug,omitempty"`
		ContainsAlcohol                         *bool    `json:"containsAlcohol,omitempty"`
		ShoppingSuggestions                     *string  `json:"shoppingSuggestions,omitempty"`
		IsStarch                                *bool    `json:"is_starch"`
		IsProtein                               *bool    `json:"is_protein"`
		IsGrain                                 *bool    `json:"is_grain"`
		IsFruit                                 *bool    `json:"is_fruit"`
		IsSalt                                  *bool    `json:"is_salt"`
		IsFat                                   *bool    `json:"is_fat"`
		IsAcid                                  *bool    `json:"is_acid"`
		IsHeat                                  *bool    `json:"is_heat"`
	}

	// ValidIngredientSearchSubset represents the subset of values suitable to index for search.
	ValidIngredientSearchSubset struct {
		_ struct{} `json:"-"`

		PluralName          string `json:"pluralName,omitempty"`
		Name                string `json:"name,omitempty"`
		ID                  string `json:"id,omitempty"`
		Description         string `json:"description,omitempty"`
		ShoppingSuggestions string `json:"shoppingSuggestions,omitempty"`
	}

	// ValidIngredientDataManager describes a structure capable of storing valid ingredients permanently.
	ValidIngredientDataManager interface {
		ValidIngredientExists(ctx context.Context, validIngredientID string) (bool, error)
		GetValidIngredient(ctx context.Context, validIngredientID string) (*ValidIngredient, error)
		GetRandomValidIngredient(ctx context.Context) (*ValidIngredient, error)
		GetValidIngredients(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidIngredient], error)
		SearchForValidIngredients(ctx context.Context, query string, filter *QueryFilter) (*QueryFilteredResult[ValidIngredient], error)
		SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string, filter *QueryFilter) (x *QueryFilteredResult[ValidIngredient], err error)
		CreateValidIngredient(ctx context.Context, input *ValidIngredientDatabaseCreationInput) (*ValidIngredient, error)
		UpdateValidIngredient(ctx context.Context, updated *ValidIngredient) error
		MarkValidIngredientAsIndexed(ctx context.Context, validIngredientID string) error
		ArchiveValidIngredient(ctx context.Context, validIngredientID string) error
		GetValidIngredientIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error)
		GetValidIngredientsWithIDs(ctx context.Context, ids []string) ([]*ValidIngredient, error)
	}

	// ValidIngredientDataService describes a structure capable of serving traffic related to valid ingredients.
	ValidIngredientDataService interface {
		SearchHandler(http.ResponseWriter, *http.Request)
		ForValidIngredientStateHandler(http.ResponseWriter, *http.Request)
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		RandomHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
		SearchByPreparationAndIngredientNameHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ValidIngredientUpdateRequestInput with a valid ingredient.
func (x *ValidIngredient) Update(input *ValidIngredientUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.Warning != nil && *input.Warning != x.Warning {
		x.Warning = *input.Warning
	}

	if input.ContainsEgg != nil && *input.ContainsEgg != x.ContainsEgg {
		x.ContainsEgg = *input.ContainsEgg
	}

	if input.ContainsDairy != nil && *input.ContainsDairy != x.ContainsDairy {
		x.ContainsDairy = *input.ContainsDairy
	}

	if input.ContainsPeanut != nil && *input.ContainsPeanut != x.ContainsPeanut {
		x.ContainsPeanut = *input.ContainsPeanut
	}

	if input.ContainsTreeNut != nil && *input.ContainsTreeNut != x.ContainsTreeNut {
		x.ContainsTreeNut = *input.ContainsTreeNut
	}

	if input.ContainsSoy != nil && *input.ContainsSoy != x.ContainsSoy {
		x.ContainsSoy = *input.ContainsSoy
	}

	if input.ContainsWheat != nil && *input.ContainsWheat != x.ContainsWheat {
		x.ContainsWheat = *input.ContainsWheat
	}

	if input.ContainsShellfish != nil && *input.ContainsShellfish != x.ContainsShellfish {
		x.ContainsShellfish = *input.ContainsShellfish
	}

	if input.ContainsSesame != nil && *input.ContainsSesame != x.ContainsSesame {
		x.ContainsSesame = *input.ContainsSesame
	}

	if input.ContainsFish != nil && *input.ContainsFish != x.ContainsFish {
		x.ContainsFish = *input.ContainsFish
	}

	if input.ContainsGluten != nil && *input.ContainsGluten != x.ContainsGluten {
		x.ContainsGluten = *input.ContainsGluten
	}

	if input.AnimalFlesh != nil && *input.AnimalFlesh != x.AnimalFlesh {
		x.AnimalFlesh = *input.AnimalFlesh
	}

	if input.IsMeasuredVolumetrically != nil && *input.IsMeasuredVolumetrically != x.IsMeasuredVolumetrically {
		x.IsMeasuredVolumetrically = *input.IsMeasuredVolumetrically
	}

	if input.IsLiquid != nil && *input.IsLiquid != x.IsLiquid {
		x.IsLiquid = *input.IsLiquid
	}

	if input.IconPath != nil && *input.IconPath != x.IconPath {
		x.IconPath = *input.IconPath
	}

	if input.PluralName != nil && *input.PluralName != x.PluralName {
		x.PluralName = *input.PluralName
	}

	if input.AnimalDerived != nil && *input.AnimalDerived != x.AnimalDerived {
		x.AnimalDerived = *input.AnimalDerived
	}

	if input.RestrictToPreparations != nil && *input.RestrictToPreparations != x.RestrictToPreparations {
		x.RestrictToPreparations = *input.RestrictToPreparations
	}

	if input.MinimumIdealStorageTemperatureInCelsius != nil && input.MinimumIdealStorageTemperatureInCelsius != x.MinimumIdealStorageTemperatureInCelsius {
		x.MinimumIdealStorageTemperatureInCelsius = input.MinimumIdealStorageTemperatureInCelsius
	}

	if input.MaximumIdealStorageTemperatureInCelsius != nil && input.MaximumIdealStorageTemperatureInCelsius != x.MaximumIdealStorageTemperatureInCelsius {
		x.MaximumIdealStorageTemperatureInCelsius = input.MaximumIdealStorageTemperatureInCelsius
	}

	if input.StorageInstructions != nil && *input.StorageInstructions != x.StorageInstructions {
		x.StorageInstructions = *input.StorageInstructions
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
	}

	if input.ContainsAlcohol != nil && *input.ContainsAlcohol != x.ContainsAlcohol {
		x.ContainsAlcohol = *input.ContainsAlcohol
	}

	if input.ShoppingSuggestions != nil && *input.ShoppingSuggestions != x.ShoppingSuggestions {
		x.ShoppingSuggestions = *input.ShoppingSuggestions
	}

	if input.IsStarch != nil && *input.IsStarch != x.IsStarch {
		x.IsStarch = *input.IsStarch
	}

	if input.IsProtein != nil && *input.IsProtein != x.IsProtein {
		x.IsProtein = *input.IsProtein
	}

	if input.IsGrain != nil && *input.IsGrain != x.IsGrain {
		x.IsGrain = *input.IsGrain
	}

	if input.IsFruit != nil && *input.IsFruit != x.IsFruit {
		x.IsFruit = *input.IsFruit
	}

	if input.IsSalt != nil && *input.IsSalt != x.IsSalt {
		x.IsSalt = *input.IsSalt
	}

	if input.IsFat != nil && *input.IsFat != x.IsFat {
		x.IsFat = *input.IsFat
	}

	if input.IsAcid != nil && *input.IsAcid != x.IsAcid {
		x.IsAcid = *input.IsAcid
	}

	if input.IsHeat != nil && *input.IsHeat != x.IsHeat {
		x.IsHeat = *input.IsHeat
	}
}

var _ validation.ValidatableWithContext = (*ValidIngredientCreationRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientCreationRequestInput.
func (x *ValidIngredientCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidIngredientDatabaseCreationInput.
func (x *ValidIngredientDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientUpdateRequestInput.
func (x *ValidIngredientUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
