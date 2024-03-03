package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetAuditLogEntry fetches an audit log entry.
func (c *Client) GetAuditLogEntry(ctx context.Context, auditLogEntryID string) (*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if auditLogEntryID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AuditLogEntryIDKey, auditLogEntryID)
	tracing.AttachToSpan(span, keys.AuditLogEntryIDKey, auditLogEntryID)

	req, err := c.requestBuilder.BuildGetAuditLogEntryRequest(ctx, auditLogEntryID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building audit log entry request")
	}

	var apiResponse *types.APIResponse[*types.AuditLogEntry]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving audit log entry")
	}

	return apiResponse.Data, nil
}

// GetAuditLogEntriesForUser fetches audit log entries for a user.
func (c *Client) GetAuditLogEntriesForUser(ctx context.Context, resourceTypes ...string) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.WithValue(keys.AuditLogEntryResourceTypesKey, resourceTypes)

	req, err := c.requestBuilder.BuildGetAuditLogEntriesForUserRequest(ctx, resourceTypes...)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building user audit log entries request")
	}

	var apiResponse *types.APIResponse[[]*types.AuditLogEntry]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving user audit log entries")
	}

	result := &types.QueryFilteredResult[types.AuditLogEntry]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

// GetAuditLogEntriesForHousehold fetches audit log entries for a user's household.
func (c *Client) GetAuditLogEntriesForHousehold(ctx context.Context, resourceTypes ...string) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.WithValue(keys.AuditLogEntryResourceTypesKey, resourceTypes)

	req, err := c.requestBuilder.BuildGetAuditLogEntriesForHouseholdRequest(ctx, resourceTypes...)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building household audit log entries request")
	}

	var apiResponse *types.APIResponse[[]*types.AuditLogEntry]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving household audit log entries")
	}

	result := &types.QueryFilteredResult[types.AuditLogEntry]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}
