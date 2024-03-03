package validingredientstateingredients

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

type validIngredientStateIngredientsServiceHTTPRoutesTestHelper struct {
	ctx                                   context.Context
	req                                   *http.Request
	res                                   *httptest.ResponseRecorder
	service                               *service
	exampleUser                           *types.User
	exampleHousehold                      *types.Household
	exampleValidIngredient                *types.ValidIngredient
	exampleValidIngredientState           *types.ValidIngredientState
	exampleValidIngredientStateIngredient *types.ValidIngredientStateIngredient
	exampleCreationInput                  *types.ValidIngredientStateIngredientCreationRequestInput
	exampleUpdateInput                    *types.ValidIngredientStateIngredientUpdateRequestInput
}

func buildTestHelper(t *testing.T) *validIngredientStateIngredientsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &validIngredientStateIngredientsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleValidIngredient = fakes.BuildFakeValidIngredient()
	helper.exampleValidIngredientState = fakes.BuildFakeValidIngredientState()
	helper.exampleValidIngredientStateIngredient = fakes.BuildFakeValidIngredientStateIngredient()
	helper.exampleCreationInput = converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(helper.exampleValidIngredientStateIngredient)
	helper.exampleUpdateInput = converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientUpdateRequestInput(helper.exampleValidIngredientStateIngredient)

	helper.service.validIngredientStateIngredientIDFetcher = func(*http.Request) string {
		return helper.exampleValidIngredientStateIngredient.ID
	}

	helper.service.validIngredientIDFetcher = func(*http.Request) string {
		return helper.exampleValidIngredient.ID
	}

	helper.service.validIngredientStateIDFetcher = func(*http.Request) string {
		return helper.exampleValidIngredientState.ID
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
