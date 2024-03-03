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

func TestMealPlanGroceryListItems(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlanGroceryListItemsTestSuite))
}

type mealPlanGroceryListItemsBaseSuite struct {
	suite.Suite
	ctx                                        context.Context
	exampleMealPlanGroceryListItem             *types.MealPlanGroceryListItem
	exampleMealPlanGroceryListItemResponse     *types.APIResponse[*types.MealPlanGroceryListItem]
	exampleMealPlanGroceryListItemListResponse *types.APIResponse[[]*types.MealPlanGroceryListItem]
	exampleMealPlanID                          string
	exampleMealPlanGroceryListItemList         []*types.MealPlanGroceryListItem
}

var _ suite.SetupTestSuite = (*mealPlanGroceryListItemsBaseSuite)(nil)

func (s *mealPlanGroceryListItemsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlanID = fakes.BuildFakeID()
	s.exampleMealPlanGroceryListItem = fakes.BuildFakeMealPlanGroceryListItem()
	s.exampleMealPlanGroceryListItemResponse = &types.APIResponse[*types.MealPlanGroceryListItem]{
		Data: s.exampleMealPlanGroceryListItem,
	}
	exampleList := fakes.BuildFakeMealPlanGroceryListItemList()
	s.exampleMealPlanGroceryListItemList = exampleList.Data
	s.exampleMealPlanGroceryListItemListResponse = &types.APIResponse[[]*types.MealPlanGroceryListItem]{
		Data:       s.exampleMealPlanGroceryListItemList,
		Pagination: &exampleList.Pagination,
	}
}

type mealPlanGroceryListItemsTestSuite struct {
	suite.Suite
	mealPlanGroceryListItemsBaseSuite
}

func (s *mealPlanGroceryListItemsTestSuite) TestClient_GetMealPlanGroceryListItem() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanGroceryListItemResponse)
		actual, err := c.GetMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanGroceryListItem, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanGroceryListItem(s.ctx, "", s.exampleMealPlanGroceryListItem.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanGroceryListItemsTestSuite) TestClient_GetMealPlanGroceryListItems() {
	const expectedPath = "/api/v1/meal_plans/%s/grocery_list_items"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, s.exampleMealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanGroceryListItemListResponse)
		actual, err := c.GetMealPlanGroceryListItemsForMealPlan(s.ctx, s.exampleMealPlanID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanGroceryListItemList, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanGroceryListItemsForMealPlan(s.ctx, "")

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanGroceryListItemsForMealPlan(s.ctx, s.exampleMealPlanID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, s.exampleMealPlanID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanGroceryListItemsForMealPlan(s.ctx, s.exampleMealPlanID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanGroceryListItemsTestSuite) TestClient_CreateMealPlanGroceryListItem() {
	const expectedPath = "/api/v1/meal_plans/%s/grocery_list_items"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanGroceryListItemCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleMealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanGroceryListItemResponse)

		actual, err := c.CreateMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanGroceryListItem, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealPlanGroceryListItemCreationRequestInput{}

		actual, err := c.CreateMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemCreationRequestInput(s.exampleMealPlanGroceryListItem)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemCreationRequestInput(s.exampleMealPlanGroceryListItem)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanGroceryListItemsTestSuite) TestClient_UpdateMealPlanGroceryListItem() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items/%s"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemUpdateRequestInput(s.exampleMealPlanGroceryListItem)
		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanGroceryListItemResponse)

		err := c.UpdateMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID, exampleInput)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemUpdateRequestInput(s.exampleMealPlanGroceryListItem)
		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemUpdateRequestInput(s.exampleMealPlanGroceryListItem)
		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID, exampleInput)
		assert.Error(t, err)
	})
}

func (s *mealPlanGroceryListItemsTestSuite) TestClient_ArchiveMealPlanGroceryListItem() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanGroceryListItemResponse)

		err := c.ArchiveMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanGroceryListItem(s.ctx, "", s.exampleMealPlanGroceryListItem.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMealPlanGroceryListItem(s.ctx, s.exampleMealPlanID, s.exampleMealPlanGroceryListItem.ID)
		assert.Error(t, err)
	})
}
