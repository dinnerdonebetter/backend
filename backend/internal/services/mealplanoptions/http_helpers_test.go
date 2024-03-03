package mealplanoptions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"
)

type mealPlanOptionsServiceHTTPRoutesTestHelper struct {
	ctx                   context.Context
	req                   *http.Request
	res                   *httptest.ResponseRecorder
	service               *service
	exampleUser           *types.User
	exampleHousehold      *types.Household
	exampleMealPlan       *types.MealPlan
	exampleMealPlanEvent  *types.MealPlanEvent
	exampleMealPlanOption *types.MealPlanOption
	exampleCreationInput  *types.MealPlanOptionCreationRequestInput
	exampleUpdateInput    *types.MealPlanOptionUpdateRequestInput
}

func buildTestHelper(t *testing.T) *mealPlanOptionsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &mealPlanOptionsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleMealPlan = fakes.BuildFakeMealPlan()
	helper.exampleMealPlan.BelongsToHousehold = helper.exampleHousehold.ID
	helper.exampleMealPlanEvent = fakes.BuildFakeMealPlanEvent()
	helper.exampleMealPlanEvent.BelongsToMealPlan = helper.exampleMealPlan.ID
	helper.exampleMealPlanOption = fakes.BuildFakeMealPlanOption()
	helper.exampleMealPlanOption.BelongsToMealPlanEvent = helper.exampleMealPlanEvent.ID
	helper.exampleCreationInput = converters.ConvertMealPlanOptionToMealPlanOptionCreationRequestInput(helper.exampleMealPlanOption)
	helper.exampleUpdateInput = converters.ConvertMealPlanOptionToMealPlanOptionUpdateRequestInput(helper.exampleMealPlanOption)

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanEventIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanEvent.ID
	}

	helper.service.mealPlanOptionIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanOption.ID
	}

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
			UserID:                   helper.exampleUser.ID,
			AccountStatus:            helper.exampleUser.AccountStatus,
			AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
			ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
		},
		ActiveHouseholdID: helper.exampleHousehold.ID,
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
			helper.exampleHousehold.ID: authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdMemberRole.String()),
		},
	}

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), types.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
