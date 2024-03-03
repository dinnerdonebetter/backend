package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetMealPlanGroceryListItemRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanGroceryListItem.ID)

		actual, err := helper.builder.BuildGetMealPlanGroceryListItemRequest(helper.ctx, exampleMealPlanID, exampleMealPlanGroceryListItem.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan task ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetMealPlanGroceryListItemRequest(helper.ctx, exampleMealPlanID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		actual, err := helper.builder.BuildGetMealPlanGroceryListItemRequest(helper.ctx, exampleMealPlanID, exampleMealPlanGroceryListItem.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateMealPlanGroceryListItemRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanGroceryListItemCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, exampleMealPlanID)

		actual, err := helper.builder.BuildCreateMealPlanGroceryListItemRequest(helper.ctx, exampleMealPlanID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan task ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateMealPlanGroceryListItemRequest(helper.ctx, exampleMealPlanID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanGroceryListItemCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealPlanGroceryListItemRequest(helper.ctx, exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetMealPlanGroceryListItemsForMealPlanRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealPlanID)

		actual, err := helper.builder.BuildGetMealPlanGroceryListItemsForMealPlanRequest(helper.ctx, exampleMealPlanID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetMealPlanGroceryListItemsForMealPlanRequest(helper.ctx, exampleMealPlanID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateMealPlanGroceryListItemRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanGroceryListItemID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanGroceryListItemUpdateRequestInput()
		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanGroceryListItemID)

		actual, err := helper.builder.BuildUpdateMealPlanGroceryListItemRequest(helper.ctx, exampleMealPlanID, exampleMealPlanGroceryListItemID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan task ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanGroceryListItemID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildUpdateMealPlanGroceryListItemRequest(helper.ctx, exampleMealPlanID, exampleMealPlanGroceryListItemID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanGroceryListItemID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanGroceryListItemUpdateRequestInput()

		actual, err := helper.builder.BuildUpdateMealPlanGroceryListItemRequest(helper.ctx, exampleMealPlanID, exampleMealPlanGroceryListItemID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
