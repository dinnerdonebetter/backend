package types

import (
	"context"
	"database/sql"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanOptionVoteCreatedCustomerEventType indicates a meal plan option vote was created.
	MealPlanOptionVoteCreatedCustomerEventType ServiceEventType = "meal_plan_option_vote_created"
	// MealPlanOptionVoteUpdatedCustomerEventType indicates a meal plan option vote was updated.
	MealPlanOptionVoteUpdatedCustomerEventType ServiceEventType = "meal_plan_option_vote_updated"
	// MealPlanOptionVoteArchivedCustomerEventType indicates a meal plan option vote was archived.
	MealPlanOptionVoteArchivedCustomerEventType ServiceEventType = "meal_plan_option_vote_archived"
)

func init() {
	gob.Register(new(MealPlanOptionVote))
	gob.Register(new(MealPlanOptionVoteCreationRequestInput))
	gob.Register(new(MealPlanOptionVoteUpdateRequestInput))
}

type (
	// MealPlanOptionVote represents a meal plan option vote.
	MealPlanOptionVote struct {
		_ struct{} `json:"-"`

		CreatedAt               time.Time  `json:"createdAt"`
		ArchivedAt              *time.Time `json:"archivedAt"`
		LastUpdatedAt           *time.Time `json:"lastUpdatedAt"`
		ID                      string     `json:"id"`
		Notes                   string     `json:"notes"`
		BelongsToMealPlanOption string     `json:"belongsToMealPlanOption"`
		ByUser                  string     `json:"byUser"`
		Rank                    uint8      `json:"rank"`
		Abstain                 bool       `json:"abstain"`
	}

	// NullableMealPlanOptionVote represents a fully nullable meal plan option vote.
	NullableMealPlanOptionVote struct {
		_ struct{} `json:"-"`

		Rank                    *uint8
		ID                      *string
		Notes                   *string
		BelongsToMealPlanOption *string
		ByUser                  *string
		Abstain                 *bool
		LastUpdatedAt           sql.NullTime
		CreatedAt               sql.NullTime
		ArchivedAt              sql.NullTime
	}

	// MealPlanOptionVoteCreationInput represents what a user could set as input for creating meal plan option votes.
	MealPlanOptionVoteCreationInput struct {
		_ struct{} `json:"-"`

		ID                      string `json:"-"`
		Notes                   string `json:"notes"`
		ByUser                  string `json:"-"`
		BelongsToMealPlanOption string `json:"belongsToMealPlanOption"`
		Rank                    uint8  `json:"rank"`
		Abstain                 bool   `json:"abstain"`
	}

	// MealPlanOptionVoteCreationRequestInput is a pending container for multiple votes.
	MealPlanOptionVoteCreationRequestInput struct {
		_ struct{} `json:"-"`

		Votes []*MealPlanOptionVoteCreationInput `json:"votes"`
	}

	// MealPlanOptionVotesDatabaseCreationInput represents what a user could set as input for creating meal plan option votes.
	MealPlanOptionVotesDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ByUser string
		Votes  []*MealPlanOptionVoteCreationInput
	}

	// MealPlanOptionVoteUpdateRequestInput represents what a user could set as input for updating meal plan option votes.
	MealPlanOptionVoteUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes                   *string `json:"notes,omitempty"`
		Rank                    *uint8  `json:"rank,omitempty"`
		Abstain                 *bool   `json:"abstain,omitempty"`
		BelongsToMealPlanOption string  `json:"belongsToMealPlanOption"`
	}

	// MissingVote represents missing votes in a meal plan.
	MissingVote struct {
		EventID  string `json:"eventID"`
		OptionID string `json:"optionID"`
		UserID   string `json:"userID"`
	}

	// MealPlanOptionVoteDataManager describes a structure capable of storing meal plan option votes permanently.
	MealPlanOptionVoteDataManager interface {
		MealPlanOptionVoteExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (bool, error)
		GetMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*MealPlanOptionVote, error)
		GetMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *QueryFilter) (*QueryFilteredResult[MealPlanOptionVote], error)
		GetMealPlanOptionVotesForMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (x []*MealPlanOptionVote, err error)
		CreateMealPlanOptionVote(ctx context.Context, input *MealPlanOptionVotesDatabaseCreationInput) ([]*MealPlanOptionVote, error)
		UpdateMealPlanOptionVote(ctx context.Context, updated *MealPlanOptionVote) error
		ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error
	}

	// MealPlanOptionVoteDataService describes a structure capable of serving traffic related to meal plan option votes.
	MealPlanOptionVoteDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an MealPlanOptionVoteUpdateRequestInput with a meal plan option vote.
func (x *MealPlanOptionVote) Update(input *MealPlanOptionVoteUpdateRequestInput) {
	if input.Rank != nil && *input.Rank != x.Rank {
		x.Rank = *input.Rank
	}

	if input.Abstain != nil && *input.Abstain != x.Abstain {
		x.Abstain = *input.Abstain
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}
}

var _ validation.ValidatableWithContext = (*MealPlanOptionVoteCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionVoteCreationRequestInput.
func (x *MealPlanOptionVoteCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	// LATER: we should validate the contents of each individual vote
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Votes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanOptionVotesDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanOptionVotesDatabaseCreationInput.
func (x *MealPlanOptionVotesDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Votes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanOptionVoteUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionVoteUpdateRequestInput.
func (x *MealPlanOptionVoteUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		// we cannot validate rank because zero is a valid value.
		validation.Field(&x.BelongsToMealPlanOption, validation.Required),
	)
}
