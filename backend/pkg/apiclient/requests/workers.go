package requests

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	workersBasePath = "workers"
)

// BuildRunFinalizeMealPlansWorkerRequest builds an HTTP request for running a worker.
func (b *Builder) BuildRunFinalizeMealPlansWorkerRequest(ctx context.Context, input *types.FinalizeMealPlansRequest) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	uri := b.BuildURL(ctx, nil, workersBasePath, "finalize_meal_plans")

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// BuildRunMealPlanGroceryListInitializationWorkerRequest builds an HTTP request for running a worker.
func (b *Builder) BuildRunMealPlanGroceryListInitializationWorkerRequest(ctx context.Context) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(ctx, nil, workersBasePath, "meal_plan_grocery_list_init")

	return b.buildDataRequest(ctx, http.MethodPost, uri, nil)
}

// BuildRunMealPlanTaskCreationWorkerRequest builds an HTTP request for running a worker.
func (b *Builder) BuildRunMealPlanTaskCreationWorkerRequest(ctx context.Context) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(ctx, nil, workersBasePath, "meal_plan_tasks")

	return b.buildDataRequest(ctx, http.MethodPost, uri, nil)
}
