package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.AuditLogEntryDataManager = (*Querier)(nil)
)

// GetAuditLogEntry fetches an audit log entry from the database.
func (q *Querier) GetAuditLogEntry(ctx context.Context, auditLogEntryID string) (*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if auditLogEntryID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AuditLogEntryIDKey, auditLogEntryID)
	tracing.AttachToSpan(span, keys.AuditLogEntryIDKey, auditLogEntryID)

	result, err := q.generatedQuerier.GetAuditLogEntry(ctx, q.db, auditLogEntryID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entry")
	}

	auditLogEntry := &types.AuditLogEntry{
		CreatedAt:          result.CreatedAt,
		BelongsToHousehold: database.StringPointerFromNullString(result.BelongsToHousehold),
		ID:                 result.ID,
		ResourceType:       result.ResourceType,
		RelevantID:         result.RelevantID,
		EventType:          types.AuditLogEntryEventType(result.EventType),
		BelongsToUser:      database.StringFromNullString(result.BelongsToUser),
	}

	if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
	}

	return auditLogEntry, nil
}

// GetAuditLogEntriesForUser fetches a list of audit log entries from the database that meet a particular filter.
func (q *Querier) GetAuditLogEntriesForUser(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &types.QueryFilteredResult[types.AuditLogEntry]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetAuditLogEntriesForUser(ctx, q.db, &generated.GetAuditLogEntriesForUserParams{
		BelongsToUser: database.NullStringFromString(userID),
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	for _, result := range results {
		auditLogEntry := &types.AuditLogEntry{
			CreatedAt:          result.CreatedAt,
			BelongsToHousehold: database.StringPointerFromNullString(result.BelongsToHousehold),
			ID:                 result.ID,
			ResourceType:       result.ResourceType,
			RelevantID:         result.RelevantID,
			EventType:          types.AuditLogEntryEventType(result.EventType),
			BelongsToUser:      database.StringFromNullString(result.BelongsToUser),
		}

		if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		x.Data = append(x.Data, auditLogEntry)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetAuditLogEntriesForUserAndResourceType fetches a list of audit log entries from the database that meet a particular filter.
func (q *Querier) GetAuditLogEntriesForUserAndResourceType(ctx context.Context, userID string, resourceTypes []string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if len(resourceTypes) == 0 {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.AuditLogEntryResourceTypesKey, resourceTypes)
	tracing.AttachToSpan(span, keys.AuditLogEntryResourceTypesKey, resourceTypes)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &types.QueryFilteredResult[types.AuditLogEntry]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetAuditLogEntriesForUserAndResourceType(ctx, q.db, &generated.GetAuditLogEntriesForUserAndResourceTypeParams{
		BelongsToUser: database.NullStringFromString(userID),
		Resources:     resourceTypes,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	for _, result := range results {
		auditLogEntry := &types.AuditLogEntry{
			CreatedAt:          result.CreatedAt,
			BelongsToHousehold: database.StringPointerFromNullString(result.BelongsToHousehold),
			ID:                 result.ID,
			ResourceType:       result.ResourceType,
			RelevantID:         result.RelevantID,
			EventType:          types.AuditLogEntryEventType(result.EventType),
			BelongsToUser:      database.StringFromNullString(result.BelongsToUser),
		}

		if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		x.Data = append(x.Data, auditLogEntry)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetAuditLogEntriesForHousehold fetches a list of audit log entries from the database that meet a particular filter.
func (q *Querier) GetAuditLogEntriesForHousehold(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &types.QueryFilteredResult[types.AuditLogEntry]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetAuditLogEntriesForHousehold(ctx, q.db, &generated.GetAuditLogEntriesForHouseholdParams{
		BelongsToHousehold: database.NullStringFromString(householdID),
		CreatedBefore:      database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:       database.NullTimeFromTimePointer(filter.CreatedAfter),
		QueryOffset:        database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:         database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	for _, result := range results {
		auditLogEntry := &types.AuditLogEntry{
			CreatedAt:          result.CreatedAt,
			BelongsToHousehold: database.StringPointerFromNullString(result.BelongsToHousehold),
			ID:                 result.ID,
			ResourceType:       result.ResourceType,
			RelevantID:         result.RelevantID,
			EventType:          types.AuditLogEntryEventType(result.EventType),
			BelongsToUser:      database.StringFromNullString(result.BelongsToUser),
		}

		if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		x.Data = append(x.Data, auditLogEntry)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetAuditLogEntriesForHouseholdAndResourceType fetches a list of audit log entries from the database that meet a particular filter.
func (q *Querier) GetAuditLogEntriesForHouseholdAndResourceType(ctx context.Context, householdID string, resourceTypes []string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	if len(resourceTypes) == 0 {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.AuditLogEntryResourceTypesKey, resourceTypes)
	tracing.AttachToSpan(span, keys.AuditLogEntryResourceTypesKey, resourceTypes)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &types.QueryFilteredResult[types.AuditLogEntry]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetAuditLogEntriesForHouseholdAndResourceType(ctx, q.db, &generated.GetAuditLogEntriesForHouseholdAndResourceTypeParams{
		BelongsToHousehold: database.NullStringFromString(householdID),
		Resources:          resourceTypes,
		CreatedBefore:      database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:       database.NullTimeFromTimePointer(filter.CreatedAfter),
		QueryOffset:        database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:         database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	for _, result := range results {
		auditLogEntry := &types.AuditLogEntry{
			CreatedAt:          result.CreatedAt,
			BelongsToHousehold: database.StringPointerFromNullString(result.BelongsToHousehold),
			ID:                 result.ID,
			ResourceType:       result.ResourceType,
			RelevantID:         result.RelevantID,
			EventType:          types.AuditLogEntryEventType(result.EventType),
			BelongsToUser:      database.StringFromNullString(result.BelongsToUser),
		}

		if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		x.Data = append(x.Data, auditLogEntry)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// createAuditLogEntry creates an audit log entry in a database.
func (q *Querier) createAuditLogEntry(ctx context.Context, querier database.SQLQueryExecutor, input *types.AuditLogEntryDatabaseCreationInput) (*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, input.BelongsToHousehold)
	logger = logger.WithValue(keys.HouseholdIDKey, input.BelongsToHousehold)

	tracing.AttachToSpan(span, keys.UserIDKey, input.BelongsToUser)
	logger = logger.WithValue(keys.UserIDKey, input.BelongsToUser)

	logger.Debug("createAuditLogEntry invoked")

	marshaledChanges, err := json.Marshal(input.Changes)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "serializing audit log change list")
	}

	if err = q.generatedQuerier.CreateAuditLogEntry(ctx, querier, &generated.CreateAuditLogEntryParams{
		ID:                 input.ID,
		ResourceType:       input.ResourceType,
		RelevantID:         input.RelevantID,
		EventType:          generated.AuditLogEventType(input.EventType),
		Changes:            marshaledChanges,
		BelongsToUser:      sql.NullString{String: input.BelongsToUser, Valid: strings.TrimSpace(input.BelongsToUser) != ""},
		BelongsToHousehold: database.NullStringFromStringPointer(input.BelongsToHousehold),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing audit log creation query")
	}

	x := &types.AuditLogEntry{
		ID:                 input.ID,
		Changes:            input.Changes,
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedAt:          q.currentTime(),
		ResourceType:       input.ResourceType,
		RelevantID:         input.RelevantID,
		EventType:          input.EventType,
		BelongsToUser:      input.BelongsToUser,
	}

	tracing.AttachToSpan(span, keys.AuditLogEntryIDKey, x.ID)

	return x, nil
}
