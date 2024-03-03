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

func TestMealPlanOptions(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlanOptionsTestSuite))
}

type mealPlanOptionsBaseSuite struct {
	suite.Suite
	ctx                               context.Context
	exampleMealPlanOption             *types.MealPlanOption
	exampleMealPlanOptionResponse     *types.APIResponse[*types.MealPlanOption]
	exampleMealPlanOptionListResponse *types.APIResponse[[]*types.MealPlanOption]
	exampleMealPlanEventID            string
	exampleMealPlanID                 string
	exampleMealPlanOptionList         []*types.MealPlanOption
}

var _ suite.SetupTestSuite = (*mealPlanOptionsBaseSuite)(nil)

func (s *mealPlanOptionsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlanID = fakes.BuildFakeID()
	s.exampleMealPlanEventID = fakes.BuildFakeID()
	s.exampleMealPlanOption = fakes.BuildFakeMealPlanOption()
	s.exampleMealPlanOption.BelongsToMealPlanEvent = s.exampleMealPlanID
	s.exampleMealPlanOptionResponse = &types.APIResponse[*types.MealPlanOption]{
		Data: s.exampleMealPlanOption,
	}
	exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()
	s.exampleMealPlanOptionList = exampleMealPlanOptionList.Data
	s.exampleMealPlanOptionListResponse = &types.APIResponse[[]*types.MealPlanOption]{
		Data:       s.exampleMealPlanOptionList,
		Pagination: &exampleMealPlanOptionList.Pagination,
	}
}

type mealPlanOptionsTestSuite struct {
	suite.Suite
	mealPlanOptionsBaseSuite
}

func (s *mealPlanOptionsTestSuite) TestClient_GetMealPlanOption() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionResponse)
		actual, err := c.GetMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanOption, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOption(s.ctx, "", s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionsTestSuite) TestClient_GetMealPlanOptions() {
	const expectedPath = "/api/v1/meal_plans/%s/events/%s/options"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleMealPlanID, s.exampleMealPlanEventID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionListResponse)
		actual, err := c.GetMealPlanOptions(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanOptionList, actual.Data)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptions(s.ctx, "", s.exampleMealPlanEventID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOptions(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleMealPlanID, s.exampleMealPlanEventID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOptions(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionsTestSuite) TestClient_CreateMealPlanOption() {
	const expectedPath = "/api/v1/meal_plans/%s/events/%s/options"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleMealPlanID, s.exampleMealPlanEventID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionResponse)

		actual, err := c.CreateMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanOption, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealPlanOptionCreationRequestInput{}

		actual, err := c.CreateMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanOptionToMealPlanOptionCreationRequestInput(s.exampleMealPlanOption)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanOptionToMealPlanOptionCreationRequestInput(s.exampleMealPlanOption)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionsTestSuite) TestClient_UpdateMealPlanOption() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanOption.BelongsToMealPlanEvent, s.exampleMealPlanOption.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionResponse)

		err := c.UpdateMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanOption)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealPlanOption(s.ctx, s.exampleMealPlanID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanOption)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanOption)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionsTestSuite) TestClient_ArchiveMealPlanOption() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionResponse)

		err := c.ArchiveMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanOption(s.ctx, "", s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEventID, s.exampleMealPlanOption.ID)
		assert.Error(t, err)
	})
}
