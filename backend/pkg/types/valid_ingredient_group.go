package types

import (
	"context"
	"encoding/gob"
	"math"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidIngredientGroupCreatedCustomerEventType indicates a valid ingredient group was created.
	ValidIngredientGroupCreatedCustomerEventType ServiceEventType = "valid_ingredient_group_created"
	// ValidIngredientGroupUpdatedCustomerEventType indicates a valid ingredient group was updated.
	ValidIngredientGroupUpdatedCustomerEventType ServiceEventType = "valid_ingredient_group_updated"
	// ValidIngredientGroupArchivedCustomerEventType indicates a valid ingredient group was archived.
	ValidIngredientGroupArchivedCustomerEventType ServiceEventType = "valid_ingredient_group_archived"
)

func init() {
	gob.Register(new(ValidIngredientGroup))
	gob.Register(new(ValidIngredientGroupCreationRequestInput))
	gob.Register(new(ValidIngredientGroupUpdateRequestInput))
}

type (
	// ValidIngredientGroup represents a valid ingredient group.
	ValidIngredientGroup struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time                     `json:"createdAt"`
		LastUpdatedAt *time.Time                    `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time                    `json:"archivedAt"`
		ID            string                        `json:"id"`
		Name          string                        `json:"name"`
		Slug          string                        `json:"slug"`
		Description   string                        `json:"description"`
		Members       []*ValidIngredientGroupMember `json:"members"`
	}

	// ValidIngredientGroupMember represents a valid ingredient group member.
	ValidIngredientGroupMember struct {
		_ struct{} `json:"-"`

		CreatedAt       time.Time       `json:"createdAt"`
		ArchivedAt      *time.Time      `json:"archivedAt"`
		ID              string          `json:"id"`
		BelongsToGroup  string          `json:"belongsToGroup"`
		ValidIngredient ValidIngredient `json:"validIngredient"`
	}

	// ValidIngredientGroupCreationRequestInput represents what a user could set as input for creating valid ingredient groups.
	ValidIngredientGroupCreationRequestInput struct {
		_ struct{} `json:"-"`

		Name        string                                            `json:"name"`
		Slug        string                                            `json:"slug"`
		Description string                                            `json:"description"`
		Members     []*ValidIngredientGroupMemberCreationRequestInput `json:"members"`
	}

	// ValidIngredientGroupMemberCreationRequestInput represents what a user could set as input for creating valid ingredient group members.
	ValidIngredientGroupMemberCreationRequestInput struct {
		_ struct{} `json:"-"`

		ValidIngredientID string `json:"validIngredientID"`
	}

	// ValidIngredientGroupDatabaseCreationInput represents what a user could set as input for creating valid ingredient groups.
	ValidIngredientGroupDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID          string
		Name        string
		Slug        string
		Description string
		Members     []*ValidIngredientGroupMemberDatabaseCreationInput
	}

	// ValidIngredientGroupMemberDatabaseCreationInput represents what a user could set as input for creating valid ingredient groups.
	ValidIngredientGroupMemberDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                string
		ValidIngredientID string
	}

	// ValidIngredientGroupUpdateRequestInput represents what a user could set as input for updating valid ingredient groups.
	ValidIngredientGroupUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name        *string `json:"name,omitempty"`
		Slug        *string `json:"slug"`
		Description *string `json:"description,omitempty"`
	}

	// ValidIngredientGroupDataManager describes a structure capable of storing valid ingredient groups permanently.
	ValidIngredientGroupDataManager interface {
		ValidIngredientGroupExists(ctx context.Context, validIngredientID string) (bool, error)
		GetValidIngredientGroup(ctx context.Context, validIngredientID string) (*ValidIngredientGroup, error)
		GetValidIngredientGroups(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidIngredientGroup], error)
		SearchForValidIngredientGroups(ctx context.Context, query string, filter *QueryFilter) ([]*ValidIngredientGroup, error)
		CreateValidIngredientGroup(ctx context.Context, input *ValidIngredientGroupDatabaseCreationInput) (*ValidIngredientGroup, error)
		UpdateValidIngredientGroup(ctx context.Context, updated *ValidIngredientGroup) error
		ArchiveValidIngredientGroup(ctx context.Context, validIngredientID string) error
	}

	// ValidIngredientGroupDataService describes a structure capable of serving traffic related to valid ingredient groups.
	ValidIngredientGroupDataService interface {
		SearchHandler(http.ResponseWriter, *http.Request)
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ValidIngredientGroupUpdateRequestInput with a valid ingredient group.
func (x *ValidIngredientGroup) Update(input *ValidIngredientGroupUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}
}

var _ validation.ValidatableWithContext = (*ValidIngredientGroupCreationRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientGroupCreationRequestInput.
func (x *ValidIngredientGroupCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Members, validation.Length(1, math.MaxUint8)),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientGroupDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidIngredientGroupDatabaseCreationInput.
func (x *ValidIngredientGroupDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientGroupUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientGroupUpdateRequestInput.
func (x *ValidIngredientGroupUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
