package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestWorkers(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(workersTestSuite))
}

type workersTestSuite struct {
	suite.Suite

	ctx context.Context
}

var _ suite.SetupTestSuite = (*workersTestSuite)(nil)

func (s *workersTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *workersTestSuite) TestClient_RunFinalizeMealPlansWorker() {
	const expectedPath = "/api/v1/workers/finalize_meal_plans"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeFinalizeMealPlansRequest()
		exampleResponse := fakes.BuildFakeFinalizeMealPlansResponse()
		exampleAPIResponse := &types.APIResponse[*types.FinalizeMealPlansResponse]{
			Data: exampleResponse,
		}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAPIResponse)

		actual, err := c.RunFinalizeMealPlansWorker(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.RunFinalizeMealPlansWorker(s.ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeFinalizeMealPlansRequest()
		c := buildTestClientWithInvalidURL(t)

		actual, err := c.RunFinalizeMealPlansWorker(s.ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeFinalizeMealPlansRequest()
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.RunFinalizeMealPlansWorker(s.ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *workersTestSuite) TestClient_RunMealPlanGroceryListInitializationWorker() {
	const expectedPath = "/api/v1/workers/meal_plan_grocery_list_init"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, &types.APIResponse[*types.FinalizeMealPlansResponse]{
			Data: fakes.BuildFakeFinalizeMealPlansResponse(),
		})

		assert.NoError(t, c.RunMealPlanGroceryListInitializationWorker(s.ctx))
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.RunMealPlanGroceryListInitializationWorker(s.ctx))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.RunMealPlanGroceryListInitializationWorker(s.ctx))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.RunMealPlanGroceryListInitializationWorker(s.ctx))
	})
}

func (s *workersTestSuite) TestClient_RunMealPlanTaskCreationWorker() {
	const expectedPath = "/api/v1/workers/meal_plan_tasks"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, &types.APIResponse[any]{})

		assert.NoError(t, c.RunMealPlanTaskCreationWorker(s.ctx))
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.RunMealPlanTaskCreationWorker(s.ctx))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.RunMealPlanTaskCreationWorker(s.ctx))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.RunMealPlanGroceryListInitializationWorker(s.ctx))
	})
}
