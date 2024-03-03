package requests

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	adminBasePath = "admin"
)

// BuildUserAccountStatusUpdateInputRequest builds a request to change a user's account status.
func (b *Builder) BuildUserAccountStatusUpdateInputRequest(ctx context.Context, input *types.UserAccountStatusUpdateInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	uri := b.BuildURL(ctx, nil, adminBasePath, usersBasePath, "status")

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}
