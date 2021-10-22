package postgres

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var _ types.AdminUserDataManager = (*SQLQuerier)(nil)

const setUserReputationQuery = `
	UPDATE users SET reputation = $1, reputation_explanation = $2 WHERE archived_on IS NULL AND id = $3
`

// UpdateUserReputation updates a user's household status.
func (q *SQLQuerier) UpdateUserReputation(ctx context.Context, userID string, input *types.UserReputationUpdateInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	args := []interface{}{
		input.NewReputation,
		input.Reason,
		input.TargetUserID,
	}

	if err := q.performWriteQuery(ctx, q.db, "user status update query", setUserReputationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "user status update")
	}

	logger.Info("user reputation updated")

	return nil
}
