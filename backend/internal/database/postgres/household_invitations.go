package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.HouseholdInvitationDataManager = (*Querier)(nil)
)

// HouseholdInvitationExists fetches whether a household invitation exists from the database.
func (q *Querier) HouseholdInvitationExists(ctx context.Context, householdInvitationID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdInvitationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)

	result, err := q.generatedQuerier.CheckHouseholdInvitationExistence(ctx, q.db, householdInvitationID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing household invitation existence check")
	}

	return result, nil
}

// GetHouseholdInvitationByHouseholdAndID fetches an invitation from the database.
func (q *Querier) GetHouseholdInvitationByHouseholdAndID(ctx context.Context, householdID, householdInvitationID string) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	if householdInvitationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)

	result, err := q.generatedQuerier.GetHouseholdInvitationByHouseholdAndID(ctx, q.db, &generated.GetHouseholdInvitationByHouseholdAndIDParams{
		DestinationHousehold: householdID,
		ID:                   householdInvitationID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching household invitation")
	}

	householdInvitation := &types.HouseholdInvitation{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		ToUser:        database.StringPointerFromNullString(result.ToUser),
		Status:        string(result.Status),
		ToEmail:       result.ToEmail,
		StatusNote:    result.StatusNote,
		Token:         result.Token,
		ID:            result.ID,
		Note:          result.Note,
		ToName:        result.ToName,
		ExpiresAt:     result.ExpiresAt,
		DestinationHousehold: types.Household{
			CreatedAt:                  result.HouseholdCreatedAt,
			SubscriptionPlanID:         database.StringPointerFromNullString(result.HouseholdSubscriptionPlanID),
			LastUpdatedAt:              database.TimePointerFromNullTime(result.HouseholdLastUpdatedAt),
			ArchivedAt:                 database.TimePointerFromNullTime(result.HouseholdArchivedAt),
			ContactPhone:               result.HouseholdContactPhone,
			BillingStatus:              result.HouseholdBillingStatus,
			AddressLine1:               result.HouseholdAddressLine1,
			AddressLine2:               result.HouseholdAddressLine2,
			City:                       result.HouseholdCity,
			State:                      result.HouseholdState,
			ZipCode:                    result.HouseholdZipCode,
			Country:                    result.HouseholdCountry,
			Latitude:                   database.Float64PointerFromNullString(result.HouseholdLatitude),
			Longitude:                  database.Float64PointerFromNullString(result.HouseholdLongitude),
			PaymentProcessorCustomerID: result.HouseholdPaymentProcessorCustomerID,
			BelongsToUser:              result.HouseholdBelongsToUser,
			ID:                         result.HouseholdID,
			Name:                       result.HouseholdName,
			Members:                    nil,
		},
		FromUser: types.User{
			CreatedAt:                  result.UserCreatedAt,
			PasswordLastChangedAt:      database.TimePointerFromNullTime(result.UserPasswordLastChangedAt),
			LastUpdatedAt:              database.TimePointerFromNullTime(result.UserLastUpdatedAt),
			LastAcceptedTermsOfService: database.TimePointerFromNullTime(result.UserLastAcceptedTermsOfService),
			LastAcceptedPrivacyPolicy:  database.TimePointerFromNullTime(result.UserLastAcceptedPrivacyPolicy),
			TwoFactorSecretVerifiedAt:  database.TimePointerFromNullTime(result.UserTwoFactorSecretVerifiedAt),
			AvatarSrc:                  database.StringPointerFromNullString(result.UserAvatarSrc),
			Birthday:                   database.TimePointerFromNullTime(result.UserBirthday),
			ArchivedAt:                 database.TimePointerFromNullTime(result.UserArchivedAt),
			AccountStatusExplanation:   result.UserUserAccountStatusExplanation,
			TwoFactorSecret:            result.UserTwoFactorSecret,
			HashedPassword:             result.UserHashedPassword,
			ID:                         result.UserID,
			AccountStatus:              result.UserUserAccountStatus,
			Username:                   result.UserUsername,
			FirstName:                  result.UserFirstName,
			LastName:                   result.UserLastName,
			EmailAddress:               result.UserEmailAddress,
			EmailAddressVerifiedAt:     database.TimePointerFromNullTime(result.UserEmailAddressVerifiedAt),
			ServiceRole:                result.UserServiceRole,
			RequiresPasswordChange:     result.UserRequiresPasswordChange,
		},
	}

	return householdInvitation, nil
}

// GetHouseholdInvitationByTokenAndID fetches an invitation from the database.
func (q *Querier) GetHouseholdInvitationByTokenAndID(ctx context.Context, token, invitationID string) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if token == "" {
		return nil, ErrInvalidIDProvided
	}

	if invitationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, invitationID)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, invitationID)

	logger.Debug("fetching household invitation")

	result, err := q.generatedQuerier.GetHouseholdInvitationByTokenAndID(ctx, q.db, &generated.GetHouseholdInvitationByTokenAndIDParams{
		Token: token,
		ID:    invitationID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching household invitation")
	}

	householdInvitation := &types.HouseholdInvitation{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		ToUser:        database.StringPointerFromNullString(result.ToUser),
		Status:        string(result.Status),
		ToEmail:       result.ToEmail,
		StatusNote:    result.StatusNote,
		Token:         result.Token,
		ID:            result.ID,
		Note:          result.Note,
		ToName:        result.ToName,
		ExpiresAt:     result.ExpiresAt,
		DestinationHousehold: types.Household{
			CreatedAt:                  result.HouseholdCreatedAt,
			SubscriptionPlanID:         database.StringPointerFromNullString(result.HouseholdSubscriptionPlanID),
			LastUpdatedAt:              database.TimePointerFromNullTime(result.HouseholdLastUpdatedAt),
			ArchivedAt:                 database.TimePointerFromNullTime(result.HouseholdArchivedAt),
			ContactPhone:               result.HouseholdContactPhone,
			BillingStatus:              result.HouseholdBillingStatus,
			AddressLine1:               result.HouseholdAddressLine1,
			AddressLine2:               result.HouseholdAddressLine2,
			City:                       result.HouseholdCity,
			State:                      result.HouseholdState,
			ZipCode:                    result.HouseholdZipCode,
			Country:                    result.HouseholdCountry,
			Latitude:                   database.Float64PointerFromNullString(result.HouseholdLatitude),
			Longitude:                  database.Float64PointerFromNullString(result.HouseholdLongitude),
			PaymentProcessorCustomerID: result.HouseholdPaymentProcessorCustomerID,
			BelongsToUser:              result.HouseholdBelongsToUser,
			ID:                         result.HouseholdID,
			Name:                       result.HouseholdName,
			Members:                    nil,
		},
		FromUser: types.User{
			CreatedAt:                  result.UserCreatedAt,
			PasswordLastChangedAt:      database.TimePointerFromNullTime(result.UserPasswordLastChangedAt),
			LastUpdatedAt:              database.TimePointerFromNullTime(result.UserLastUpdatedAt),
			LastAcceptedTermsOfService: database.TimePointerFromNullTime(result.UserLastAcceptedTermsOfService),
			LastAcceptedPrivacyPolicy:  database.TimePointerFromNullTime(result.UserLastAcceptedPrivacyPolicy),
			TwoFactorSecretVerifiedAt:  database.TimePointerFromNullTime(result.UserTwoFactorSecretVerifiedAt),
			AvatarSrc:                  database.StringPointerFromNullString(result.UserAvatarSrc),
			Birthday:                   database.TimePointerFromNullTime(result.UserBirthday),
			ArchivedAt:                 database.TimePointerFromNullTime(result.UserArchivedAt),
			AccountStatusExplanation:   result.UserUserAccountStatusExplanation,
			TwoFactorSecret:            result.UserTwoFactorSecret,
			HashedPassword:             result.UserHashedPassword,
			ID:                         result.UserID,
			AccountStatus:              result.UserUserAccountStatus,
			Username:                   result.UserUsername,
			FirstName:                  result.UserFirstName,
			LastName:                   result.UserLastName,
			EmailAddress:               result.UserEmailAddress,
			EmailAddressVerifiedAt:     database.TimePointerFromNullTime(result.UserEmailAddressVerifiedAt),
			ServiceRole:                result.UserServiceRole,
			RequiresPasswordChange:     result.UserRequiresPasswordChange,
		},
	}

	return householdInvitation, nil
}

// GetHouseholdInvitationByEmailAndToken fetches an invitation from the database.
func (q *Querier) GetHouseholdInvitationByEmailAndToken(ctx context.Context, emailAddress, token string) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if emailAddress == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserEmailAddressKey, emailAddress)
	tracing.AttachToSpan(span, keys.UserEmailAddressKey, emailAddress)

	if token == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationTokenKey, token)
	tracing.AttachToSpan(span, keys.HouseholdInvitationTokenKey, token)

	result, err := q.generatedQuerier.GetHouseholdInvitationByEmailAndToken(ctx, q.db, &generated.GetHouseholdInvitationByEmailAndTokenParams{
		ToEmail: emailAddress,
		Token:   token,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching household invitation")
	}

	invitation := &types.HouseholdInvitation{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		ToUser:        database.StringPointerFromNullString(result.ToUser),
		Status:        string(result.Status),
		ToEmail:       result.ToEmail,
		StatusNote:    result.StatusNote,
		Token:         result.Token,
		ID:            result.ID,
		Note:          result.Note,
		ToName:        result.ToName,
		ExpiresAt:     result.ExpiresAt,
		DestinationHousehold: types.Household{
			CreatedAt:                  result.HouseholdCreatedAt,
			SubscriptionPlanID:         database.StringPointerFromNullString(result.HouseholdSubscriptionPlanID),
			LastUpdatedAt:              database.TimePointerFromNullTime(result.HouseholdLastUpdatedAt),
			ArchivedAt:                 database.TimePointerFromNullTime(result.HouseholdArchivedAt),
			ContactPhone:               result.HouseholdContactPhone,
			BillingStatus:              result.HouseholdBillingStatus,
			AddressLine1:               result.HouseholdAddressLine1,
			AddressLine2:               result.HouseholdAddressLine2,
			City:                       result.HouseholdCity,
			State:                      result.HouseholdState,
			ZipCode:                    result.HouseholdZipCode,
			Country:                    result.HouseholdCountry,
			Latitude:                   database.Float64PointerFromNullString(result.HouseholdLatitude),
			Longitude:                  database.Float64PointerFromNullString(result.HouseholdLongitude),
			PaymentProcessorCustomerID: result.HouseholdPaymentProcessorCustomerID,
			BelongsToUser:              result.HouseholdBelongsToUser,
			ID:                         result.HouseholdID,
			Name:                       result.HouseholdName,
			Members:                    nil,
		},
		FromUser: types.User{
			CreatedAt:                  result.UserCreatedAt,
			PasswordLastChangedAt:      database.TimePointerFromNullTime(result.UserPasswordLastChangedAt),
			LastUpdatedAt:              database.TimePointerFromNullTime(result.UserLastUpdatedAt),
			LastAcceptedTermsOfService: database.TimePointerFromNullTime(result.UserLastAcceptedTermsOfService),
			LastAcceptedPrivacyPolicy:  database.TimePointerFromNullTime(result.UserLastAcceptedPrivacyPolicy),
			TwoFactorSecretVerifiedAt:  database.TimePointerFromNullTime(result.UserTwoFactorSecretVerifiedAt),
			AvatarSrc:                  database.StringPointerFromNullString(result.UserAvatarSrc),
			Birthday:                   database.TimePointerFromNullTime(result.UserBirthday),
			ArchivedAt:                 database.TimePointerFromNullTime(result.UserArchivedAt),
			AccountStatusExplanation:   result.UserUserAccountStatusExplanation,
			TwoFactorSecret:            result.UserTwoFactorSecret,
			HashedPassword:             result.UserHashedPassword,
			ID:                         result.UserID,
			AccountStatus:              result.UserUserAccountStatus,
			Username:                   result.UserUsername,
			FirstName:                  result.UserFirstName,
			LastName:                   result.UserLastName,
			EmailAddress:               result.UserEmailAddress,
			EmailAddressVerifiedAt:     database.TimePointerFromNullTime(result.UserEmailAddressVerifiedAt),
			ServiceRole:                result.UserServiceRole,
			RequiresPasswordChange:     result.UserRequiresPasswordChange,
		},
	}

	return invitation, nil
}

// CreateHouseholdInvitation creates an invitation in a database.
func (q *Querier) CreateHouseholdInvitation(ctx context.Context, input *types.HouseholdInvitationDatabaseCreationInput) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.HouseholdInvitationIDKey, input.ID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, input.DestinationHouseholdID)

	if err := q.generatedQuerier.CreateHouseholdInvitation(ctx, q.db, &generated.CreateHouseholdInvitationParams{
		ExpiresAt:            input.ExpiresAt,
		ID:                   input.ID,
		FromUser:             input.FromUser,
		ToName:               input.ToName,
		Note:                 input.Note,
		ToEmail:              input.ToEmail,
		Token:                input.Token,
		DestinationHousehold: input.DestinationHouseholdID,
		ToUser:               database.NullStringFromStringPointer(input.ToUser),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing household invitation creation query")
	}

	x := &types.HouseholdInvitation{
		ID:                   input.ID,
		FromUser:             types.User{ID: input.FromUser},
		ToUser:               input.ToUser,
		Note:                 input.Note,
		ToName:               input.ToName,
		ToEmail:              input.ToEmail,
		Token:                input.Token,
		StatusNote:           "",
		Status:               string(types.PendingHouseholdInvitationStatus),
		DestinationHousehold: types.Household{ID: input.DestinationHouseholdID},
		ExpiresAt:            input.ExpiresAt,
		CreatedAt:            q.currentTime(),
	}

	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, x.ID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, x.ID)

	logger.Info("household invitation created")

	return x, nil
}

// GetPendingHouseholdInvitationsFromUser fetches pending household invitations sent from a given user.
func (q *Querier) GetPendingHouseholdInvitationsFromUser(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInvitation], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger := q.logger.WithValue(keys.UserIDKey, userID)
	filter.AttachToLogger(logger)

	x := &types.QueryFilteredResult[types.HouseholdInvitation]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetPendingInvitesFromUser(ctx, q.db, &generated.GetPendingInvitesFromUserParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
		Status:        generated.InvitationState(types.PendingHouseholdInvitationStatus),
		FromUser:      userID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing household invitation query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.HouseholdInvitation{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			ToUser:        database.StringPointerFromNullString(result.ToUser),
			Status:        string(result.Status),
			ToEmail:       result.ToEmail,
			StatusNote:    result.StatusNote,
			Token:         result.Token,
			ID:            result.ID,
			Note:          result.Note,
			ToName:        result.ToName,
			ExpiresAt:     result.ExpiresAt,
			DestinationHousehold: types.Household{
				CreatedAt:                  result.HouseholdCreatedAt,
				SubscriptionPlanID:         database.StringPointerFromNullString(result.HouseholdSubscriptionPlanID),
				LastUpdatedAt:              database.TimePointerFromNullTime(result.HouseholdLastUpdatedAt),
				ArchivedAt:                 database.TimePointerFromNullTime(result.HouseholdArchivedAt),
				ContactPhone:               result.HouseholdContactPhone,
				BillingStatus:              result.HouseholdBillingStatus,
				AddressLine1:               result.HouseholdAddressLine1,
				AddressLine2:               result.HouseholdAddressLine2,
				City:                       result.HouseholdCity,
				State:                      result.HouseholdState,
				ZipCode:                    result.HouseholdZipCode,
				Country:                    result.HouseholdCountry,
				Latitude:                   database.Float64PointerFromNullString(result.HouseholdLatitude),
				Longitude:                  database.Float64PointerFromNullString(result.HouseholdLongitude),
				PaymentProcessorCustomerID: result.HouseholdPaymentProcessorCustomerID,
				BelongsToUser:              result.HouseholdBelongsToUser,
				ID:                         result.HouseholdID,
				Name:                       result.HouseholdName,
				Members:                    nil,
			},
			FromUser: types.User{
				CreatedAt:                  result.UserCreatedAt,
				PasswordLastChangedAt:      database.TimePointerFromNullTime(result.UserPasswordLastChangedAt),
				LastUpdatedAt:              database.TimePointerFromNullTime(result.UserLastUpdatedAt),
				LastAcceptedTermsOfService: database.TimePointerFromNullTime(result.UserLastAcceptedTermsOfService),
				LastAcceptedPrivacyPolicy:  database.TimePointerFromNullTime(result.UserLastAcceptedPrivacyPolicy),
				TwoFactorSecretVerifiedAt:  database.TimePointerFromNullTime(result.UserTwoFactorSecretVerifiedAt),
				AvatarSrc:                  database.StringPointerFromNullString(result.UserAvatarSrc),
				Birthday:                   database.TimePointerFromNullTime(result.UserBirthday),
				ArchivedAt:                 database.TimePointerFromNullTime(result.UserArchivedAt),
				AccountStatusExplanation:   result.UserUserAccountStatusExplanation,
				TwoFactorSecret:            result.UserTwoFactorSecret,
				HashedPassword:             result.UserHashedPassword,
				ID:                         result.UserID,
				AccountStatus:              result.UserUserAccountStatus,
				Username:                   result.UserUsername,
				FirstName:                  result.UserFirstName,
				LastName:                   result.UserLastName,
				EmailAddress:               result.UserEmailAddress,
				EmailAddressVerifiedAt:     database.TimePointerFromNullTime(result.UserEmailAddressVerifiedAt),
				ServiceRole:                result.UserServiceRole,
				RequiresPasswordChange:     result.UserRequiresPasswordChange,
			},
		})

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetPendingHouseholdInvitationsForUser fetches pending household invitations sent to a given user.
func (q *Querier) GetPendingHouseholdInvitationsForUser(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInvitation], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger := q.logger.WithValue(keys.UserIDKey, userID)
	filter.AttachToLogger(logger)

	x := &types.QueryFilteredResult[types.HouseholdInvitation]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetPendingInvitesForUser(ctx, q.db, &generated.GetPendingInvitesForUserParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
		Status:        generated.InvitationState(types.PendingHouseholdInvitationStatus),
		ToUser:        database.NullStringFromString(userID),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing household invitation query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.HouseholdInvitation{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			ToUser:        database.StringPointerFromNullString(result.ToUser),
			Status:        string(result.Status),
			ToEmail:       result.ToEmail,
			StatusNote:    result.StatusNote,
			Token:         result.Token,
			ID:            result.ID,
			Note:          result.Note,
			ToName:        result.ToName,
			ExpiresAt:     result.ExpiresAt,
			DestinationHousehold: types.Household{
				CreatedAt:                  result.HouseholdCreatedAt,
				SubscriptionPlanID:         database.StringPointerFromNullString(result.HouseholdSubscriptionPlanID),
				LastUpdatedAt:              database.TimePointerFromNullTime(result.HouseholdLastUpdatedAt),
				ArchivedAt:                 database.TimePointerFromNullTime(result.HouseholdArchivedAt),
				ContactPhone:               result.HouseholdContactPhone,
				BillingStatus:              result.HouseholdBillingStatus,
				AddressLine1:               result.HouseholdAddressLine1,
				AddressLine2:               result.HouseholdAddressLine2,
				City:                       result.HouseholdCity,
				State:                      result.HouseholdState,
				ZipCode:                    result.HouseholdZipCode,
				Country:                    result.HouseholdCountry,
				Latitude:                   database.Float64PointerFromNullString(result.HouseholdLatitude),
				Longitude:                  database.Float64PointerFromNullString(result.HouseholdLongitude),
				PaymentProcessorCustomerID: result.HouseholdPaymentProcessorCustomerID,
				BelongsToUser:              result.HouseholdBelongsToUser,
				ID:                         result.HouseholdID,
				Name:                       result.HouseholdName,
				Members:                    nil,
			},
			FromUser: types.User{
				CreatedAt:                  result.UserCreatedAt,
				PasswordLastChangedAt:      database.TimePointerFromNullTime(result.UserPasswordLastChangedAt),
				LastUpdatedAt:              database.TimePointerFromNullTime(result.UserLastUpdatedAt),
				LastAcceptedTermsOfService: database.TimePointerFromNullTime(result.UserLastAcceptedTermsOfService),
				LastAcceptedPrivacyPolicy:  database.TimePointerFromNullTime(result.UserLastAcceptedPrivacyPolicy),
				TwoFactorSecretVerifiedAt:  database.TimePointerFromNullTime(result.UserTwoFactorSecretVerifiedAt),
				AvatarSrc:                  database.StringPointerFromNullString(result.UserAvatarSrc),
				Birthday:                   database.TimePointerFromNullTime(result.UserBirthday),
				ArchivedAt:                 database.TimePointerFromNullTime(result.UserArchivedAt),
				AccountStatusExplanation:   result.UserUserAccountStatusExplanation,
				TwoFactorSecret:            result.UserTwoFactorSecret,
				HashedPassword:             result.UserHashedPassword,
				ID:                         result.UserID,
				AccountStatus:              result.UserUserAccountStatus,
				Username:                   result.UserUsername,
				FirstName:                  result.UserFirstName,
				LastName:                   result.UserLastName,
				EmailAddress:               result.UserEmailAddress,
				EmailAddressVerifiedAt:     database.TimePointerFromNullTime(result.UserEmailAddressVerifiedAt),
				ServiceRole:                result.UserServiceRole,
				RequiresPasswordChange:     result.UserRequiresPasswordChange,
			},
		})

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

func (q *Querier) setInvitationStatus(ctx context.Context, querier database.SQLQueryExecutor, householdInvitationID, note, status string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("new_status", status)

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)

	if err := q.generatedQuerier.SetHouseholdInvitationStatus(ctx, querier, &generated.SetHouseholdInvitationStatusParams{
		Status:     generated.InvitationState(status),
		StatusNote: note,
		ID:         householdInvitationID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "changing household invitation status")
	}

	logger.Debug("household invitation updated")

	return nil
}

// CancelHouseholdInvitation cancels a household invitation by its ID with a note.
func (q *Querier) CancelHouseholdInvitation(ctx context.Context, householdInvitationID, note string) error {
	return q.setInvitationStatus(ctx, q.db, householdInvitationID, note, string(types.CancelledHouseholdInvitationStatus))
}

// AcceptHouseholdInvitation accepts a household invitation by its ID with a note.
func (q *Querier) AcceptHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if err = q.setInvitationStatus(ctx, tx, householdInvitationID, note, string(types.AcceptedHouseholdInvitationStatus)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "accepting household invitation")
	}

	invitation, err := q.GetHouseholdInvitationByTokenAndID(ctx, token, householdInvitationID)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "fetching household invitation")
	}

	addUserInput := &types.HouseholdUserMembershipDatabaseCreationInput{
		ID:            identifiers.New(),
		Reason:        fmt.Sprintf("accepted household invitation %q", householdInvitationID),
		HouseholdID:   invitation.DestinationHousehold.ID,
		HouseholdRole: "household_member",
	}
	if invitation.ToUser != nil {
		addUserInput.UserID = *invitation.ToUser
	}

	if err = q.addUserToHousehold(ctx, tx, addUserInput); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "adding user to household")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return nil
}

// RejectHouseholdInvitation rejects a household invitation by its ID with a note.
func (q *Querier) RejectHouseholdInvitation(ctx context.Context, householdInvitationID, note string) error {
	return q.setInvitationStatus(ctx, q.db, householdInvitationID, note, string(types.RejectedHouseholdInvitationStatus))
}

func (q *Querier) attachInvitationsToUser(ctx context.Context, querier database.SQLQueryExecutor, userEmail, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if userEmail == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserEmailAddressKey, userEmail)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, userEmail)

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, userID)

	if err := q.generatedQuerier.AttachHouseholdInvitationsToUserID(ctx, querier, &generated.AttachHouseholdInvitationsToUserIDParams{
		ToEmail: userEmail,
		ToUser:  database.NullStringFromString(userID),
	}); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareAndLogError(err, logger, span, "attaching invitations to user")
	}

	logger.Info("invitations associated with user")

	return nil
}

func (q *Querier) acceptInvitationForUser(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.UserDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue(keys.UsernameKey, input.Username).WithValue(keys.UserEmailAddressKey, input.EmailAddress)

	invitation, tokenCheckErr := q.GetHouseholdInvitationByEmailAndToken(ctx, input.EmailAddress, input.InvitationToken)
	if tokenCheckErr != nil {
		q.rollbackTransaction(ctx, querier)
		return observability.PrepareError(tokenCheckErr, span, "fetching household invitation")
	}

	logger.Debug("fetched invitation to accept for user")

	if err := q.generatedQuerier.CreateHouseholdUserMembershipForNewUser(ctx, querier, &generated.CreateHouseholdUserMembershipForNewUserParams{
		ID:                 identifiers.New(),
		BelongsToUser:      input.ID,
		BelongsToHousehold: input.DestinationHouseholdID,
		HouseholdRole:      authorization.HouseholdMemberRole.String(),
		DefaultHousehold:   true,
	}); err != nil {
		q.rollbackTransaction(ctx, querier)
		return observability.PrepareAndLogError(err, logger, span, "writing destination household membership")
	}

	logger.Debug("created membership via invitation")

	if err := q.setInvitationStatus(ctx, querier, invitation.ID, "", string(types.AcceptedHouseholdInvitationStatus)); err != nil {
		q.rollbackTransaction(ctx, querier)
		return observability.PrepareAndLogError(err, logger, span, "accepting household invitation")
	}

	logger.Debug("marked invitation as accepted")

	return nil
}
