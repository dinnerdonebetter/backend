package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestMealPlans(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlansTestSuite))
}

type mealPlansBaseSuite struct {
	suite.Suite
	ctx                         context.Context
	exampleMealPlan             *types.MealPlan
	exampleMealPlanResponse     *types.APIResponse[*types.MealPlan]
	exampleMealPlanListResponse *types.APIResponse[[]*types.MealPlan]
	exampleMealPlanList         []*types.MealPlan
}

var _ suite.SetupTestSuite = (*mealPlansBaseSuite)(nil)

func (s *mealPlansBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlan = fakes.BuildFakeMealPlan()
	s.exampleMealPlanResponse = &types.APIResponse[*types.MealPlan]{
		Data: s.exampleMealPlan,
	}
	exampleMealPlanList := fakes.BuildFakeMealPlanList()
	s.exampleMealPlanList = exampleMealPlanList.Data
	s.exampleMealPlanListResponse = &types.APIResponse[[]*types.MealPlan]{
		Data:       exampleMealPlanList.Data,
		Pagination: &exampleMealPlanList.Pagination,
	}
}

type mealPlansTestSuite struct {
	suite.Suite
	mealPlansBaseSuite
}

func (s *mealPlansTestSuite) TestClient_GetMealPlan() {
	const expectedPathFormat = "/api/v1/meal_plans/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlan.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanResponse)
		actual, err := c.GetMealPlan(s.ctx, s.exampleMealPlan.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)

		require.WithinDuration(t, s.exampleMealPlan.VotingDeadline, actual.VotingDeadline, 0)

		actual.VotingDeadline = s.exampleMealPlan.VotingDeadline

		assert.Equal(t, s.exampleMealPlan, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlan(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlan(s.ctx, s.exampleMealPlan.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlan.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlan(s.ctx, s.exampleMealPlan.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlansTestSuite) TestClient_GetMealPlans() {
	const expectedPath = "/api/v1/meal_plans"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanListResponse)
		actual, err := c.GetMealPlans(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)

		for i := range actual.Data {
			require.WithinDuration(t, s.exampleMealPlanList[i].VotingDeadline, actual.Data[i].VotingDeadline, 0)
			actual.Data[i].VotingDeadline = s.exampleMealPlanList[i].VotingDeadline
		}

		assert.Equal(t, s.exampleMealPlanList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlans(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlans(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlansTestSuite) TestClient_CreateMealPlan() {
	const expectedPath = "/api/v1/meal_plans"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanResponse)

		actual, err := c.CreateMealPlan(s.ctx, exampleInput)
		assert.NoError(t, err)

		require.WithinDuration(t, s.exampleMealPlan.VotingDeadline, actual.VotingDeadline, 0)

		actual.VotingDeadline = s.exampleMealPlan.VotingDeadline

		assert.Equal(t, s.exampleMealPlan, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlan(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealPlanCreationRequestInput{}

		actual, err := c.CreateMealPlan(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanToMealPlanCreationRequestInput(s.exampleMealPlan)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealPlan(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanToMealPlanCreationRequestInput(s.exampleMealPlan)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealPlan(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlansTestSuite) TestClient_UpdateMealPlan() {
	const expectedPathFormat = "/api/v1/meal_plans/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleMealPlan.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanResponse)

		err := c.UpdateMealPlan(s.ctx, s.exampleMealPlan)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealPlan(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateMealPlan(s.ctx, s.exampleMealPlan)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateMealPlan(s.ctx, s.exampleMealPlan)
		assert.Error(t, err)
	})
}

func (s *mealPlansTestSuite) TestClient_ArchiveMealPlan() {
	const expectedPathFormat = "/api/v1/meal_plans/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMealPlan.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanResponse)

		err := c.ArchiveMealPlan(s.ctx, s.exampleMealPlan.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlan(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMealPlan(s.ctx, s.exampleMealPlan.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMealPlan(s.ctx, s.exampleMealPlan.ID)
		assert.Error(t, err)
	})
}

func (s *mealPlansTestSuite) TestClient_FinalizeMealPlan() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/finalize"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodPost, "", expectedPathFormat, s.exampleMealPlan.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanResponse)

		err := c.FinalizeMealPlan(s.ctx, s.exampleMealPlan.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.FinalizeMealPlan(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.FinalizeMealPlan(s.ctx, s.exampleMealPlan.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.FinalizeMealPlan(s.ctx, s.exampleMealPlan.ID)
		assert.Error(t, err)
	})
}
