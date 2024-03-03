package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeFinalizeMealPlansRequest builds a faked Webhook.
func BuildFakeFinalizeMealPlansRequest() *types.FinalizeMealPlansRequest {
	return &types.FinalizeMealPlansRequest{
		ReturnCount: fake.Bool(),
	}
}

// BuildFakeFinalizeMealPlansResponse builds a faked FinalizeMealPlansResponse.
func BuildFakeFinalizeMealPlansResponse() *types.FinalizeMealPlansResponse {
	return &types.FinalizeMealPlansResponse{
		Count: int(buildFakeNumber()),
	}
}
