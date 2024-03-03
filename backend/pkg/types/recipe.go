package types

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

const (
	// RecipeCreatedCustomerEventType indicates a recipe was created.
	RecipeCreatedCustomerEventType ServiceEventType = "recipe_created"
	// RecipeUpdatedCustomerEventType indicates a recipe was updated.
	RecipeUpdatedCustomerEventType ServiceEventType = "recipe_updated"
	// RecipeArchivedCustomerEventType indicates a recipe was archived.
	RecipeArchivedCustomerEventType ServiceEventType = "recipe_archived"
	// RecipeClonedCustomerEventType indicates a recipe was cloned.
	RecipeClonedCustomerEventType ServiceEventType = "recipe_cloned"
)

func init() {
	gob.Register(new(Recipe))
	gob.Register(new(RecipeCreationRequestInput))
	gob.Register(new(RecipeUpdateRequestInput))
}

type (
	// Recipe represents a recipe.
	Recipe struct {
		_ struct{} `json:"-"`

		CreatedAt                time.Time         `json:"createdAt"`
		InspiredByRecipeID       *string           `json:"inspiredByRecipeID"`
		LastUpdatedAt            *time.Time        `json:"lastUpdatedAt"`
		ArchivedAt               *time.Time        `json:"archivedAt"`
		MaximumEstimatedPortions *float32          `json:"maximumEstimatedPortions"`
		PluralPortionName        string            `json:"pluralPortionName"`
		Description              string            `json:"description"`
		Name                     string            `json:"name"`
		PortionName              string            `json:"portionName"`
		ID                       string            `json:"id"`
		CreatedByUser            string            `json:"createdByUser"`
		Source                   string            `json:"source"`
		Slug                     string            `json:"slug"`
		YieldsComponentType      string            `json:"yieldsComponentType"`
		PrepTasks                []*RecipePrepTask `json:"prepTasks"`
		Steps                    []*RecipeStep     `json:"steps"`
		Media                    []*RecipeMedia    `json:"media"`
		SupportingRecipes        []*Recipe         `json:"supportingRecipes"`
		MinimumEstimatedPortions float32           `json:"minimumEstimatedPortions"`
		SealOfApproval           bool              `json:"sealOfApproval"`
		EligibleForMeals         bool              `json:"eligibleForMeals"`
	}

	// RecipeCreationRequestInput represents what a user could set as input for creating recipes.
	RecipeCreationRequestInput struct {
		_ struct{} `json:"-"`

		InspiredByRecipeID       *string                                           `json:"inspiredByRecipeID"`
		MaximumEstimatedPortions *float32                                          `json:"maximumEstimatedPortions"`
		PortionName              string                                            `json:"portionName"`
		Source                   string                                            `json:"source"`
		Description              string                                            `json:"description"`
		PluralPortionName        string                                            `json:"pluralPortionName"`
		Name                     string                                            `json:"name"`
		Slug                     string                                            `json:"slug"`
		YieldsComponentType      string                                            `json:"yieldsComponentType"`
		PrepTasks                []*RecipePrepTaskWithinRecipeCreationRequestInput `json:"prepTasks"`
		Steps                    []*RecipeStepCreationRequestInput                 `json:"steps"`
		MinimumEstimatedPortions float32                                           `json:"minimumEstimatedPortions"`
		AlsoCreateMeal           bool                                              `json:"alsoCreateMeal"`
		SealOfApproval           bool                                              `json:"sealOfApproval"`
		EligibleForMeals         bool                                              `json:"eligibleForMeals"`
	}

	// RecipeDatabaseCreationInput represents what a user could set as input for creating recipes.
	RecipeDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		InspiredByRecipeID       *string
		MaximumEstimatedPortions *float32
		PortionName              string
		Description              string
		Name                     string
		Slug                     string
		Source                   string
		PluralPortionName        string
		CreatedByUser            string
		ID                       string
		YieldsComponentType      string
		Steps                    []*RecipeStepDatabaseCreationInput
		PrepTasks                []*RecipePrepTaskDatabaseCreationInput
		MinimumEstimatedPortions float32
		AlsoCreateMeal           bool
		SealOfApproval           bool
		EligibleForMeals         bool
	}

	// RecipeUpdateRequestInput represents what a user could set as input for updating recipes.
	RecipeUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name                     *string  `json:"name,omitempty"`
		Slug                     *string  `json:"slug"`
		Source                   *string  `json:"source,omitempty"`
		Description              *string  `json:"description,omitempty"`
		InspiredByRecipeID       *string  `json:"inspiredByRecipeID,omitempty"`
		SealOfApproval           *bool    `json:"sealOfApproval,omitempty"`
		MinimumEstimatedPortions *float32 `json:"minimumEstimatedPortions,omitempty"`
		MaximumEstimatedPortions *float32 `json:"maximumEstimatedPortions,omitempty"`
		PortionName              *string  `json:"portionName"`
		PluralPortionName        *string  `json:"pluralPortionName"`
		EligibleForMeals         *bool    `json:"eligibleForMeals"`
		YieldsComponentType      *string  `json:"yieldsComponentType"`
	}

	// RecipeSearchSubset represents the subset of values suitable to index for search.
	RecipeSearchSubset struct {
		_ struct{} `json:"-"`

		ID          string                    `json:"id,omitempty"`
		Name        string                    `json:"name,omitempty"`
		Description string                    `json:"description,omitempty"`
		Steps       []*RecipeStepSearchSubset `json:"steps,omitempty"`
	}

	// RecipeDataManager describes a structure capable of storing recipes permanently.
	RecipeDataManager interface {
		RecipeExists(ctx context.Context, recipeID string) (bool, error)
		GetRecipe(ctx context.Context, recipeID string) (*Recipe, error)
		GetRecipes(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[Recipe], error)
		SearchForRecipes(ctx context.Context, query string, filter *QueryFilter) (*QueryFilteredResult[Recipe], error)
		CreateRecipe(ctx context.Context, input *RecipeDatabaseCreationInput) (*Recipe, error)
		UpdateRecipe(ctx context.Context, updated *Recipe) error
		MarkRecipeAsIndexed(ctx context.Context, recipeID string) error
		ArchiveRecipe(ctx context.Context, recipeID, userID string) error
		GetRecipeIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error)
		GetRecipesWithIDs(ctx context.Context, ids []string) ([]*Recipe, error)
	}

	// RecipeDataService describes a structure capable of serving traffic related to recipes.
	RecipeDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		SearchHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
		EstimatedPrepStepsHandler(http.ResponseWriter, *http.Request)
		ImageUploadHandler(http.ResponseWriter, *http.Request)
		DAGHandler(http.ResponseWriter, *http.Request)
		MermaidHandler(http.ResponseWriter, *http.Request)
		CloneHandler(http.ResponseWriter, *http.Request)
	}
)

// FindStepForIndex finds a step for a given index.
func (x *Recipe) FindStepForIndex(index uint32) *RecipeStep {
	for _, step := range x.Steps {
		if step.Index == index {
			return step
		}
	}

	// we could return an error here, but that would make my life a little harder
	// so if you fuck up and submit a wrong value, and it's your fault here.
	return nil
}

// FindStepByID finds a step for a given ID.
func (x *Recipe) FindStepByID(id string) *RecipeStep {
	for _, step := range x.Steps {
		if step.ID == id {
			return step
		}
	}

	// we could return an error here, but that would make my life a little harder
	// so if you fuck up and submit a wrong value, and it's your fault here.
	return nil
}

// FindStepForRecipeStepProductID finds a step for a given ID.
func (x *Recipe) FindStepForRecipeStepProductID(recipeStepProductID string) *RecipeStep {
	for _, step := range x.Steps {
		for _, product := range step.Products {
			if product.ID == recipeStepProductID {
				return step
			}
		}
	}

	// we could return an error here, but that would make my life a little harder
	// so if you fuck up and submit a wrong value, and it's your fault here.
	return nil
}

// FindStepIndexByID finds a step for a given ID.
func (x *Recipe) FindStepIndexByID(id string) int {
	for i, step := range x.Steps {
		if step.ID == id {
			return i
		}
	}

	return -1
}

// Update merges an RecipeUpdateRequestInput with a recipe.
func (x *Recipe) Update(input *RecipeUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
	}

	if input.Source != nil && *input.Source != x.Source {
		x.Source = *input.Source
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.InspiredByRecipeID != nil && (x.InspiredByRecipeID == nil || (*input.InspiredByRecipeID != "" && *input.InspiredByRecipeID != *x.InspiredByRecipeID)) {
		x.InspiredByRecipeID = input.InspiredByRecipeID
	}

	if input.SealOfApproval != nil && *input.SealOfApproval != x.SealOfApproval {
		x.SealOfApproval = *input.SealOfApproval
	}

	if input.MinimumEstimatedPortions != nil && *input.MinimumEstimatedPortions != x.MinimumEstimatedPortions {
		x.MinimumEstimatedPortions = *input.MinimumEstimatedPortions
	}

	if input.MaximumEstimatedPortions != nil && input.MaximumEstimatedPortions != x.MaximumEstimatedPortions {
		x.MaximumEstimatedPortions = input.MaximumEstimatedPortions
	}

	if input.PortionName != nil && *input.PortionName != x.PortionName {
		x.PortionName = *input.PortionName
	}

	if input.PluralPortionName != nil && *input.PluralPortionName != x.PluralPortionName {
		x.PluralPortionName = *input.PluralPortionName
	}

	if input.EligibleForMeals != nil && *input.EligibleForMeals != x.EligibleForMeals {
		x.EligibleForMeals = *input.EligibleForMeals
	}

	if input.YieldsComponentType != nil && *input.YieldsComponentType != x.YieldsComponentType {
		x.YieldsComponentType = *input.YieldsComponentType
	}
}

var _ validation.ValidatableWithContext = (*RecipeCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeCreationRequestInput.
func (x *RecipeCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	var errResult *multierror.Error

	if len(x.Steps) < 2 {
		errResult = multierror.Append(fmt.Errorf("recipe must have at least 2 steps"), errResult)
	}

	stepIndicesMentionedInPrepTasks := map[uint32]bool{}
	for i, task := range x.PrepTasks {
		for j, step := range task.TaskSteps {
			if _, ok := stepIndicesMentionedInPrepTasks[step.BelongsToRecipeStepIndex]; ok {
				errResult = multierror.Append(fmt.Errorf("duplicate step mentioned in step %d for task %d", i+1, j+1), errResult)
			} else {
				stepIndicesMentionedInPrepTasks[step.BelongsToRecipeStepIndex] = true
			}
		}
	}

	if errResult != nil {
		return errResult
	}

	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.MinimumEstimatedPortions, validation.Required),
		validation.Field(&x.PluralPortionName, validation.Required),
		validation.Field(&x.PortionName, validation.Required),
		validation.Field(&x.Slug, validation.Required),
		validation.Field(&x.Steps, validation.Required),
		validation.Field(&x.YieldsComponentType,
			validation.Required,
			validation.In(
				MealComponentTypesUnspecified,
				MealComponentTypesAmuseBouche,
				MealComponentTypesAppetizer,
				MealComponentTypesSoup,
				MealComponentTypesMain,
				MealComponentTypesSalad,
				MealComponentTypesBeverage,
				MealComponentTypesSide,
				MealComponentTypesDessert,
			),
		),
	)
}

// FindStepByIndex finds a step for a given index.
func (x *RecipeDatabaseCreationInput) FindStepByIndex(index uint32) *RecipeStepDatabaseCreationInput {
	for _, step := range x.Steps {
		if step.Index == index {
			return step
		}
	}

	// we could return an error here, but that would make my life a little harder
	// so if you fuck up and submit a wrong value, and it's your fault here.
	return nil
}

// FindStepByID finds a step for a given ID.
func (x *RecipeDatabaseCreationInput) FindStepByID(id string) *RecipeStepDatabaseCreationInput {
	for _, step := range x.Steps {
		if step.ID == id {
			return step
		}
	}

	// we could return an error here, but that would make my life a little harder
	// so if you fuck up and submit a wrong value, and it's your fault here.
	return nil
}

var _ validation.ValidatableWithContext = (*RecipeDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeDatabaseCreationInput.
func (x *RecipeDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.CreatedByUser, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeUpdateRequestInput.
func (x *RecipeUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Source, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.MinimumEstimatedPortions, validation.Required),
	)
}
