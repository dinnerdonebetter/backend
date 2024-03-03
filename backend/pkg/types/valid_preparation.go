package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidPreparationCreatedCustomerEventType indicates a valid preparation was created.
	ValidPreparationCreatedCustomerEventType ServiceEventType = "valid_preparation_created"
	// ValidPreparationUpdatedCustomerEventType indicates a valid preparation was updated.
	ValidPreparationUpdatedCustomerEventType ServiceEventType = "valid_preparation_updated"
	// ValidPreparationArchivedCustomerEventType indicates a valid preparation was archived.
	ValidPreparationArchivedCustomerEventType ServiceEventType = "valid_preparation_archived"
)

func init() {
	gob.Register(new(ValidPreparation))
	gob.Register(new(ValidPreparationCreationRequestInput))
	gob.Register(new(ValidPreparationUpdateRequestInput))
}

type (
	// ValidPreparation represents a valid preparation.
	ValidPreparation struct {
		_ struct{} `json:"-"`

		CreatedAt                   time.Time  `json:"createdAt"`
		MaximumInstrumentCount      *int32     `json:"maximumInstrumentCount"` // TODO: make these uint16
		ArchivedAt                  *time.Time `json:"archivedAt"`
		MaximumIngredientCount      *int32     `json:"maximumIngredientCount"` // TODO: make these uint16
		LastUpdatedAt               *time.Time `json:"lastUpdatedAt"`
		MaximumVesselCount          *int32     `json:"maximumVesselCount"` // TODO: make these uint16
		IconPath                    string     `json:"iconPath"`
		PastTense                   string     `json:"pastTense"`
		ID                          string     `json:"id"`
		Name                        string     `json:"name"`
		Description                 string     `json:"description"`
		Slug                        string     `json:"slug"`
		MinimumIngredientCount      int32      `json:"minimumIngredientCount"`
		MinimumInstrumentCount      int32      `json:"minimumInstrumentCount"`
		MinimumVesselCount          int32      `json:"minimumVesselCount"`
		RestrictToIngredients       bool       `json:"restrictToIngredients"`
		TemperatureRequired         bool       `json:"temperatureRequired"`
		TimeEstimateRequired        bool       `json:"timeEstimateRequired"`
		ConditionExpressionRequired bool       `json:"conditionExpressionRequired"`
		ConsumesVessel              bool       `json:"consumesVessel"`
		OnlyForVessels              bool       `json:"onlyForVessels"`
		YieldsNothing               bool       `json:"yieldsNothing"`
	}

	// ValidPreparationCreationRequestInput represents what a user could set as input for creating valid preparations.
	ValidPreparationCreationRequestInput struct {
		_ struct{} `json:"-"`

		MaximumInstrumentCount      *int32 `json:"maximumInstrumentCount"`
		MaximumIngredientCount      *int32 `json:"maximumIngredientCount"`
		MaximumVesselCount          *int32 `json:"maximumVesselCount"`
		IconPath                    string `json:"iconPath"`
		PastTense                   string `json:"pastTense"`
		Slug                        string `json:"slug"`
		Name                        string `json:"name"`
		Description                 string `json:"description"`
		MinimumIngredientCount      int32  `json:"minimumIngredientCount"`
		MinimumVesselCount          int32  `json:"minimumVesselCount"`
		MinimumInstrumentCount      int32  `json:"minimumInstrumentCount"`
		TemperatureRequired         bool   `json:"temperatureRequired"`
		TimeEstimateRequired        bool   `json:"timeEstimateRequired"`
		ConditionExpressionRequired bool   `json:"conditionExpressionRequired"`
		ConsumesVessel              bool   `json:"consumesVessel"`
		OnlyForVessels              bool   `json:"onlyForVessels"`
		RestrictToIngredients       bool   `json:"restrictToIngredients"`
		YieldsNothing               bool   `json:"yieldsNothing"`
	}

	// ValidPreparationDatabaseCreationInput represents what a user could set as input for creating valid preparations.
	ValidPreparationDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		MaximumInstrumentCount      *int32
		MaximumIngredientCount      *int32
		MaximumVesselCount          *int32
		IconPath                    string
		PastTense                   string
		Slug                        string
		ID                          string
		Name                        string
		Description                 string
		MinimumIngredientCount      int32
		MinimumVesselCount          int32
		MinimumInstrumentCount      int32
		TemperatureRequired         bool
		TimeEstimateRequired        bool
		ConditionExpressionRequired bool
		ConsumesVessel              bool
		OnlyForVessels              bool
		RestrictToIngredients       bool
		YieldsNothing               bool
	}

	// ValidPreparationUpdateRequestInput represents what a user could set as input for updating valid preparations.
	ValidPreparationUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name                        *string `json:"name,omitempty"`
		Description                 *string `json:"description,omitempty"`
		IconPath                    *string `json:"iconPath,omitempty"`
		YieldsNothing               *bool   `json:"yieldsNothing,omitempty"`
		Slug                        *string `json:"slug,omitempty"`
		RestrictToIngredients       *bool   `json:"restrictToIngredients,omitempty"`
		PastTense                   *string `json:"pastTense,omitempty"`
		MinimumInstrumentCount      *int32  `json:"minimumInstrumentCount,omitempty"`
		MaximumInstrumentCount      *int32  `json:"maximumInstrumentCount,omitempty"`
		MinimumIngredientCount      *int32  `json:"minimumIngredientCount,omitempty"`
		MaximumIngredientCount      *int32  `json:"maximumIngredientCount,omitempty"`
		TemperatureRequired         *bool   `json:"temperatureRequired,omitempty"`
		TimeEstimateRequired        *bool   `json:"timeEstimateRequired,omitempty"`
		ConditionExpressionRequired *bool   `json:"conditionExpressionRequired,omitempty"`
		ConsumesVessel              *bool   `json:"consumesVessel,omitempty"`
		OnlyForVessels              *bool   `json:"onlyForVessels,omitempty"`
		MinimumVesselCount          *int32  `json:"minimumVesselCount,omitempty"`
		MaximumVesselCount          *int32  `json:"maximumVesselCount,omitempty"`
	}

	// ValidPreparationSearchSubset represents the subset of values suitable to index for search.
	ValidPreparationSearchSubset struct {
		_ struct{} `json:"-"`

		PastTense   string `json:"pastTense,omitempty"`
		ID          string `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}

	// ValidPreparationDataManager describes a structure capable of storing valid preparations permanently.
	ValidPreparationDataManager interface {
		ValidPreparationExists(ctx context.Context, validPreparationID string) (bool, error)
		GetValidPreparation(ctx context.Context, validPreparationID string) (*ValidPreparation, error)
		GetRandomValidPreparation(ctx context.Context) (*ValidPreparation, error)
		GetValidPreparations(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidPreparation], error)
		SearchForValidPreparations(ctx context.Context, query string) ([]*ValidPreparation, error)
		CreateValidPreparation(ctx context.Context, input *ValidPreparationDatabaseCreationInput) (*ValidPreparation, error)
		UpdateValidPreparation(ctx context.Context, updated *ValidPreparation) error
		MarkValidPreparationAsIndexed(ctx context.Context, validPreparationID string) error
		ArchiveValidPreparation(ctx context.Context, validPreparationID string) error
		GetValidPreparationIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error)
		GetValidPreparationsWithIDs(ctx context.Context, ids []string) ([]*ValidPreparation, error)
	}

	// ValidPreparationDataService describes a structure capable of serving traffic related to valid preparations.
	ValidPreparationDataService interface {
		SearchHandler(http.ResponseWriter, *http.Request)
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		RandomHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ValidPreparationUpdateRequestInput with a valid preparation.
func (x *ValidPreparation) Update(input *ValidPreparationUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.IconPath != nil && *input.IconPath != x.IconPath {
		x.IconPath = *input.IconPath
	}

	if input.YieldsNothing != nil && *input.YieldsNothing != x.YieldsNothing {
		x.YieldsNothing = *input.YieldsNothing
	}

	if input.RestrictToIngredients != nil && *input.RestrictToIngredients != x.RestrictToIngredients {
		x.RestrictToIngredients = *input.RestrictToIngredients
	}

	if input.MinimumIngredientCount != nil && *input.MinimumIngredientCount != x.MinimumIngredientCount {
		x.MinimumIngredientCount = *input.MinimumIngredientCount
	}

	if input.MaximumIngredientCount != nil && x.MaximumIngredientCount != nil && *input.MaximumIngredientCount != *x.MaximumIngredientCount {
		x.MaximumIngredientCount = input.MaximumIngredientCount
	}

	if input.MinimumInstrumentCount != nil && *input.MinimumInstrumentCount != x.MinimumInstrumentCount {
		x.MinimumInstrumentCount = *input.MinimumInstrumentCount
	}

	if input.MaximumInstrumentCount != nil && x.MaximumInstrumentCount != nil && *input.MaximumInstrumentCount != *x.MaximumInstrumentCount {
		x.MaximumInstrumentCount = input.MaximumInstrumentCount
	}

	if input.TemperatureRequired != nil && *input.TemperatureRequired != x.TemperatureRequired {
		x.TemperatureRequired = *input.TemperatureRequired
	}

	if input.TimeEstimateRequired != nil && *input.TimeEstimateRequired != x.TimeEstimateRequired {
		x.TimeEstimateRequired = *input.TimeEstimateRequired
	}

	if input.ConditionExpressionRequired != nil && *input.ConditionExpressionRequired != x.ConditionExpressionRequired {
		x.ConditionExpressionRequired = *input.ConditionExpressionRequired
	}

	if input.PastTense != nil && *input.PastTense != x.PastTense {
		x.PastTense = *input.PastTense
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
	}

	if input.ConsumesVessel != nil && *input.ConsumesVessel != x.ConsumesVessel {
		x.ConsumesVessel = *input.ConsumesVessel
	}

	if input.OnlyForVessels != nil && *input.OnlyForVessels != x.OnlyForVessels {
		x.OnlyForVessels = *input.OnlyForVessels
	}

	if input.MinimumVesselCount != nil && *input.MinimumVesselCount != x.MinimumVesselCount {
		x.MinimumVesselCount = *input.MinimumVesselCount
	}

	if input.MaximumVesselCount != nil && x.MaximumVesselCount != nil && *input.MaximumVesselCount != *x.MaximumVesselCount {
		x.MaximumVesselCount = input.MaximumVesselCount
	}
}

var _ validation.ValidatableWithContext = (*ValidPreparationCreationRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationCreationRequestInput.
func (x *ValidPreparationCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidPreparationDatabaseCreationInput.
func (x *ValidPreparationDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationUpdateRequestInput.
func (x *ValidPreparationUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
