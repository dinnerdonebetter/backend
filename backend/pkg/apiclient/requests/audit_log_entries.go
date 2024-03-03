package requests

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	auditLogEntriesBasePath = "audit_log_entries"
)

// BuildGetAuditLogEntryRequest builds an HTTP request that fetches a given audit log entry.
func (b *Builder) BuildGetAuditLogEntryRequest(ctx context.Context, auditLogEntryID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if auditLogEntryID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.AuditLogEntryIDKey, auditLogEntryID)

	uri := b.buildAPIV1URL(ctx, nil, auditLogEntriesBasePath, auditLogEntryID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetAuditLogEntriesForUserRequest builds an HTTP request that fetches a given audit log entry.
func (b *Builder) BuildGetAuditLogEntriesForUserRequest(ctx context.Context, resourceTypes ...string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	qp := url.Values{}
	for _, rt := range resourceTypes {
		qp.Add(types.AuditLogResourceTypesQueryParamKey, rt)
	}
	tracing.AttachToSpan(span, keys.AuditLogEntryResourceTypesKey, resourceTypes)

	uri := b.buildAPIV1URL(ctx, qp, auditLogEntriesBasePath, "for_user")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetAuditLogEntriesForHouseholdRequest builds an HTTP request that fetches a given audit log entry.
func (b *Builder) BuildGetAuditLogEntriesForHouseholdRequest(ctx context.Context, resourceTypes ...string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	qp := url.Values{}
	for _, rt := range resourceTypes {
		qp.Add(types.AuditLogResourceTypesQueryParamKey, rt)
	}
	tracing.AttachToSpan(span, keys.AuditLogEntryResourceTypesKey, resourceTypes)

	uri := b.buildAPIV1URL(ctx, qp, auditLogEntriesBasePath, "for_household")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
