package types

import (
	"context"
	"fmt"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// HouseholdCreatedCustomerEventType indicates a household was created.
	HouseholdCreatedCustomerEventType ServiceEventType = "household_created"
	// HouseholdUpdatedCustomerEventType indicates a household was updated.
	HouseholdUpdatedCustomerEventType ServiceEventType = "household_updated"
	// HouseholdArchivedCustomerEventType indicates a household was archived.
	HouseholdArchivedCustomerEventType ServiceEventType = "household_archived"
	// HouseholdMemberRemovedCustomerEventType indicates a household member was removed.
	HouseholdMemberRemovedCustomerEventType ServiceEventType = "household_member_removed"
	// HouseholdMembershipPermissionsUpdatedCustomerEventType indicates a household member's permissions were modified.
	HouseholdMembershipPermissionsUpdatedCustomerEventType ServiceEventType = "household_membership_permissions_updated"
	// HouseholdOwnershipTransferredCustomerEventType indicates a household was transferred to another owner.
	HouseholdOwnershipTransferredCustomerEventType ServiceEventType = "household_ownership_transferred"

	// UnpaidHouseholdBillingStatus indicates a household is not paid.
	UnpaidHouseholdBillingStatus = "unpaid"
)

type (
	// Household represents a household.
	Household struct {
		_ struct{} `json:"-"`

		CreatedAt                  time.Time                          `json:"createdAt"`
		SubscriptionPlanID         *string                            `json:"subscriptionPlanID"`
		LastUpdatedAt              *time.Time                         `json:"lastUpdatedAt"`
		ArchivedAt                 *time.Time                         `json:"archivedAt"`
		Longitude                  *float64                           `json:"longitude"`
		Latitude                   *float64                           `json:"latitude"`
		State                      string                             `json:"state"`
		ContactPhone               string                             `json:"contactPhone"`
		City                       string                             `json:"city"`
		AddressLine1               string                             `json:"addressLine1"`
		ZipCode                    string                             `json:"zipCode"`
		Country                    string                             `json:"country"`
		BillingStatus              string                             `json:"billingStatus"`
		AddressLine2               string                             `json:"addressLine2"`
		PaymentProcessorCustomerID string                             `json:"paymentProcessorCustomer"`
		BelongsToUser              string                             `json:"belongsToUser"`
		ID                         string                             `json:"id"`
		Name                       string                             `json:"name"`
		WebhookEncryptionKey       string                             `json:"-"`
		Members                    []*HouseholdUserMembershipWithUser `json:"members"`
	}

	// HouseholdCreationRequestInput represents what a User could set as input for creating households.
	HouseholdCreationRequestInput struct {
		_ struct{} `json:"-"`

		Latitude     *float64 `json:"latitude"`
		Longitude    *float64 `json:"longitude"`
		Name         string   `json:"name"`
		ContactPhone string   `json:"contactPhone"`
		AddressLine1 string   `json:"addressLine1"`
		AddressLine2 string   `json:"addressLine2"`
		City         string   `json:"city"`
		State        string   `json:"state"`
		ZipCode      string   `json:"zipCode"`
		Country      string   `json:"country"`
	}

	// HouseholdDatabaseCreationInput represents what a User could set as input for creating households.
	HouseholdDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                   string
		Name                 string
		AddressLine1         string
		AddressLine2         string
		City                 string
		State                string
		ZipCode              string
		Country              string
		Latitude             *float64
		Longitude            *float64
		ContactPhone         string
		BelongsToUser        string
		WebhookEncryptionKey string
	}

	// HouseholdUpdateRequestInput represents what a User could set as input for updating households.
	HouseholdUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name          *string  `json:"name,omitempty"`
		ContactPhone  *string  `json:"contactPhone,omitempty"`
		AddressLine1  *string  `json:"addressLine1"`
		AddressLine2  *string  `json:"addressLine2"`
		City          *string  `json:"city"`
		State         *string  `json:"state"`
		ZipCode       *string  `json:"zipCode"`
		Country       *string  `json:"country"`
		Latitude      *float64 `json:"latitude"`
		Longitude     *float64 `json:"longitude"`
		BelongsToUser string   `json:"-"`
	}

	// HouseholdDataManager describes a structure capable of storing households permanently.
	HouseholdDataManager interface {
		GetHousehold(ctx context.Context, householdID string) (*Household, error)
		GetHouseholds(ctx context.Context, userID string, filter *QueryFilter) (*QueryFilteredResult[Household], error)
		CreateHousehold(ctx context.Context, input *HouseholdDatabaseCreationInput) (*Household, error)
		UpdateHousehold(ctx context.Context, updated *Household) error
		ArchiveHousehold(ctx context.Context, householdID string, userID string) error
	}

	// HouseholdDataService describes a structure capable of serving traffic related to households.
	HouseholdDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		CurrentInfoHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
		RemoveMemberHandler(http.ResponseWriter, *http.Request)
		MarkAsDefaultHouseholdHandler(http.ResponseWriter, *http.Request)
		ModifyMemberPermissionsHandler(http.ResponseWriter, *http.Request)
		TransferHouseholdOwnershipHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges a HouseholdUpdateRequestInput with a household.
func (x *Household) Update(input *HouseholdUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.ContactPhone != nil && *input.ContactPhone != x.ContactPhone {
		x.ContactPhone = *input.ContactPhone
	}

	if input.AddressLine1 != nil && *input.AddressLine1 != x.AddressLine1 {
		x.AddressLine1 = *input.AddressLine1
	}

	if input.AddressLine2 != nil && *input.AddressLine2 != x.AddressLine2 {
		x.AddressLine2 = *input.AddressLine2
	}

	if input.City != nil && *input.City != x.City {
		x.City = *input.City
	}

	if input.State != nil && *input.State != x.State {
		x.State = *input.State
	}

	if input.ZipCode != nil && *input.ZipCode != x.ZipCode {
		x.ZipCode = *input.ZipCode
	}

	if input.Country != nil && *input.Country != x.Country {
		x.Country = *input.Country
	}

	if input.Latitude != nil && input.Latitude != x.Latitude {
		x.Latitude = input.Latitude
	}

	if input.Longitude != nil && input.Longitude != x.Longitude {
		x.Longitude = input.Longitude
	}
}

var _ validation.ValidatableWithContext = (*HouseholdCreationRequestInput)(nil)

// ValidateWithContext validates a HouseholdCreationRequestInput.
func (x *HouseholdCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Latitude, validation.NilOrNotEmpty),
		validation.Field(&x.Longitude, validation.NilOrNotEmpty),
	)
}

var _ validation.ValidatableWithContext = (*HouseholdUpdateRequestInput)(nil)

// ValidateWithContext validates a HouseholdUpdateRequestInput.
func (x *HouseholdUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Latitude, validation.NilOrNotEmpty),
		validation.Field(&x.Longitude, validation.NilOrNotEmpty),
	)
}

// HouseholdCreationInputForNewUser creates a new HouseholdInputCreation struct for a given user.
func HouseholdCreationInputForNewUser(u *User) *HouseholdDatabaseCreationInput {
	return &HouseholdDatabaseCreationInput{
		Name:          fmt.Sprintf("%s's cool household", u.Username),
		BelongsToUser: u.ID,
	}
}
