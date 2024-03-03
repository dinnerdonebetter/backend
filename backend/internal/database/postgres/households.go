package postgres

import (
	"context"
	"database/sql"
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

const (
	resourceTypeHouseholds = "households"
)

var (
	_ types.HouseholdDataManager = (*Querier)(nil)
)

// GetHousehold fetches a household from the database.
func (q *Querier) GetHousehold(ctx context.Context, householdID string) (*types.Household, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	results, err := q.generatedQuerier.GetHouseholdByIDWithMemberships(ctx, q.db, householdID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing households list retrieval query")
	}

	var household *types.Household
	for _, result := range results {
		if household == nil {
			household = &types.Household{
				CreatedAt:                  result.CreatedAt,
				SubscriptionPlanID:         database.StringPointerFromNullString(result.SubscriptionPlanID),
				LastUpdatedAt:              database.TimePointerFromNullTime(result.LastUpdatedAt),
				ArchivedAt:                 database.TimePointerFromNullTime(result.ArchivedAt),
				ContactPhone:               result.ContactPhone,
				BillingStatus:              result.BillingStatus,
				AddressLine1:               result.AddressLine1,
				AddressLine2:               result.AddressLine2,
				City:                       result.City,
				State:                      result.State,
				ZipCode:                    result.ZipCode,
				Country:                    result.Country,
				Latitude:                   database.Float64PointerFromNullString(result.Latitude),
				Longitude:                  database.Float64PointerFromNullString(result.Longitude),
				PaymentProcessorCustomerID: result.PaymentProcessorCustomerID,
				BelongsToUser:              result.BelongsToUser,
				ID:                         result.ID,
				Name:                       result.Name,
				WebhookEncryptionKey:       result.WebhookHmacSecret,
				Members:                    nil,
			}
		}

		household.Members = append(household.Members, &types.HouseholdUserMembershipWithUser{
			CreatedAt:     result.MembershipCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.MembershipLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.MembershipArchivedAt),
			ID:            result.MembershipID,
			BelongsToUser: &types.User{
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
			BelongsToHousehold: result.MembershipBelongsToHousehold,
			HouseholdRole:      result.MembershipHouseholdRole,
			DefaultHousehold:   result.MembershipDefaultHousehold,
		})
	}

	if household == nil {
		return nil, sql.ErrNoRows
	}

	return household, nil
}

// getHouseholdsForUser fetches a list of households from the database that meet a particular filter.
func (q *Querier) getHouseholdsForUser(ctx context.Context, querier database.SQLQueryExecutor, userID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Household], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.Household]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetHouseholdsForUser(ctx, querier, &generated.GetHouseholdsForUserParams{
		BelongsToUser: userID,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing households list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.Household{
			CreatedAt:                  result.CreatedAt,
			SubscriptionPlanID:         database.StringPointerFromNullString(result.SubscriptionPlanID),
			LastUpdatedAt:              database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                 database.TimePointerFromNullTime(result.ArchivedAt),
			ContactPhone:               result.ContactPhone,
			BillingStatus:              result.BillingStatus,
			AddressLine1:               result.AddressLine1,
			AddressLine2:               result.AddressLine2,
			City:                       result.City,
			State:                      result.State,
			ZipCode:                    result.ZipCode,
			Country:                    result.Country,
			Latitude:                   database.Float64PointerFromNullString(result.Latitude),
			Longitude:                  database.Float64PointerFromNullString(result.Longitude),
			PaymentProcessorCustomerID: result.PaymentProcessorCustomerID,
			BelongsToUser:              result.BelongsToUser,
			ID:                         result.ID,
			Name:                       result.Name,
			Members:                    nil,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetHouseholds fetches a list of households from the database that meet a particular filter.
func (q *Querier) GetHouseholds(ctx context.Context, userID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Household], err error) {
	return q.getHouseholdsForUser(ctx, q.db, userID, filter)
}

// CreateHousehold creates a household in the database.
func (q *Querier) CreateHousehold(ctx context.Context, input *types.HouseholdDatabaseCreationInput) (*types.Household, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, input.BelongsToUser)

	// begin household creation transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	// create the household.
	if writeErr := q.generatedQuerier.CreateHousehold(ctx, tx, &generated.CreateHouseholdParams{
		City:              input.City,
		Name:              input.Name,
		BillingStatus:     types.UnpaidHouseholdBillingStatus,
		ContactPhone:      input.ContactPhone,
		AddressLine1:      input.AddressLine1,
		AddressLine2:      input.AddressLine2,
		ID:                input.ID,
		State:             input.State,
		ZipCode:           input.ZipCode,
		Country:           input.Country,
		BelongsToUser:     input.BelongsToUser,
		WebhookHmacSecret: input.WebhookEncryptionKey,
		Latitude:          database.NullStringFromFloat64Pointer(input.Latitude),
		Longitude:         database.NullStringFromFloat64Pointer(input.Longitude),
	}); writeErr != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(writeErr, span, "creating household")
	}

	household := &types.Household{
		ID:            input.ID,
		Name:          input.Name,
		BelongsToUser: input.BelongsToUser,
		BillingStatus: types.UnpaidHouseholdBillingStatus,
		ContactPhone:  input.ContactPhone,
		AddressLine1:  input.AddressLine1,
		AddressLine2:  input.AddressLine2,
		City:          input.City,
		State:         input.State,
		ZipCode:       input.ZipCode,
		Country:       input.Country,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		CreatedAt:     q.currentTime(),
	}

	if _, err = q.createAuditLogEntry(ctx, tx, &types.AuditLogEntryDatabaseCreationInput{
		BelongsToHousehold: &household.ID,
		ID:                 identifiers.New(),
		ResourceType:       resourceTypeHouseholds,
		RelevantID:         household.ID,
		EventType:          types.AuditLogEventTypeCreated,
		BelongsToUser:      household.BelongsToUser,
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	householdMembershipID := identifiers.New()
	if err = q.generatedQuerier.AddUserToHousehold(ctx, tx, &generated.AddUserToHouseholdParams{
		ID:                 householdMembershipID,
		BelongsToUser:      household.BelongsToUser,
		BelongsToHousehold: household.ID,
		HouseholdRole:      authorization.HouseholdAdminRole.String(),
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing household membership creation query")
	}

	if _, err = q.createAuditLogEntry(ctx, tx, &types.AuditLogEntryDatabaseCreationInput{
		BelongsToHousehold: &household.ID,
		ID:                 identifiers.New(),
		ResourceType:       resourceTypeHouseholdUserMemberships,
		RelevantID:         householdMembershipID,
		EventType:          types.AuditLogEventTypeCreated,
		BelongsToUser:      household.BelongsToUser,
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, household.ID)
	logger.Info("household created")

	return household, nil
}

// UpdateHousehold updates a particular household. Note that UpdateHousehold expects the provided input to have a valid ID.
func (q *Querier) UpdateHousehold(ctx context.Context, updated *types.Household) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.HouseholdIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, updated.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	household, err := q.GetHousehold(ctx, updated.ID)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "fetching household")
	}

	if _, err = q.generatedQuerier.UpdateHousehold(ctx, q.db, &generated.UpdateHouseholdParams{
		Name:          updated.Name,
		ContactPhone:  updated.ContactPhone,
		AddressLine1:  updated.AddressLine1,
		AddressLine2:  updated.AddressLine2,
		City:          updated.City,
		State:         updated.State,
		ZipCode:       updated.ZipCode,
		Country:       updated.Country,
		BelongsToUser: updated.BelongsToUser,
		ID:            updated.ID,
		Latitude:      database.NullStringFromFloat64Pointer(updated.Latitude),
		Longitude:     database.NullStringFromFloat64Pointer(updated.Longitude),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating household")
	}

	if _, err = q.createAuditLogEntry(ctx, tx, &types.AuditLogEntryDatabaseCreationInput{
		BelongsToHousehold: &updated.ID,
		ID:                 identifiers.New(),
		ResourceType:       resourceTypeHouseholds,
		RelevantID:         updated.ID,
		EventType:          types.AuditLogEventTypeUpdated,
		BelongsToUser:      household.BelongsToUser,
		Changes:            buildChangesForHousehold(household, updated),
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("household updated")

	return nil
}

func buildChangesForHousehold(household, updated *types.Household) map[string]types.ChangeLog {
	changes := map[string]types.ChangeLog{}

	if household.Name != updated.Name {
		changes["name"] = types.ChangeLog{
			OldValue: household.Name,
			NewValue: updated.Name,
		}
	}

	if household.ContactPhone != updated.ContactPhone {
		changes["contact_phone"] = types.ChangeLog{
			OldValue: household.ContactPhone,
			NewValue: updated.ContactPhone,
		}
	}

	if household.AddressLine1 != updated.AddressLine1 {
		changes["address_line_1"] = types.ChangeLog{
			OldValue: household.AddressLine1,
			NewValue: updated.AddressLine1,
		}
	}

	if household.AddressLine2 != updated.AddressLine2 {
		changes["address_line_2"] = types.ChangeLog{
			OldValue: household.AddressLine2,
			NewValue: updated.AddressLine2,
		}
	}

	if household.City != updated.City {
		changes["city"] = types.ChangeLog{
			OldValue: household.City,
			NewValue: updated.City,
		}
	}

	if household.State != updated.State {
		changes["state"] = types.ChangeLog{
			OldValue: household.State,
			NewValue: updated.State,
		}
	}

	if household.ZipCode != updated.ZipCode {
		changes["zip_code"] = types.ChangeLog{
			OldValue: household.ZipCode,
			NewValue: updated.ZipCode,
		}
	}

	if household.Country != updated.Country {
		changes["country"] = types.ChangeLog{
			OldValue: household.Country,
			NewValue: updated.Country,
		}
	}

	if household.Latitude != updated.Latitude {
		changes["latitude"] = types.ChangeLog{
			OldValue: fmt.Sprintf("%v", household.Latitude),
			NewValue: fmt.Sprintf("%v", updated.Latitude),
		}
	}

	if household.Longitude != updated.Longitude {
		changes["longitude"] = types.ChangeLog{
			OldValue: fmt.Sprintf("%v", household.Longitude),
			NewValue: fmt.Sprintf("%v", updated.Longitude),
		}
	}

	return changes
}

// ArchiveHousehold archives a household from the database by its ID.
func (q *Querier) ArchiveHousehold(ctx context.Context, householdID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" || userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if _, err = q.generatedQuerier.ArchiveHousehold(ctx, q.db, &generated.ArchiveHouseholdParams{
		BelongsToUser: userID,
		ID:            householdID,
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving household")
	}

	if _, err = q.createAuditLogEntry(ctx, tx, &types.AuditLogEntryDatabaseCreationInput{
		BelongsToHousehold: &householdID,
		ID:                 identifiers.New(),
		ResourceType:       resourceTypeHouseholds,
		RelevantID:         householdID,
		EventType:          types.AuditLogEventTypeCreated,
		BelongsToUser:      userID,
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("household archived")

	return nil
}
