package validingredientmeasurementunits

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

type validIngredientMeasurementUnitsServiceHTTPRoutesTestHelper struct {
	ctx                                   context.Context
	req                                   *http.Request
	res                                   *httptest.ResponseRecorder
	service                               *service
	exampleUser                           *types.User
	exampleHousehold                      *types.Household
	exampleValidIngredient                *types.ValidIngredient
	exampleValidMeasurementUnit           *types.ValidMeasurementUnit
	exampleValidIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit
	exampleCreationInput                  *types.ValidIngredientMeasurementUnitCreationRequestInput
	exampleUpdateInput                    *types.ValidIngredientMeasurementUnitUpdateRequestInput
}

func buildTestHelper(t *testing.T) *validIngredientMeasurementUnitsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &validIngredientMeasurementUnitsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleValidIngredient = fakes.BuildFakeValidIngredient()
	helper.exampleValidMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
	helper.exampleValidIngredientMeasurementUnit = fakes.BuildFakeValidIngredientMeasurementUnit()
	helper.exampleCreationInput = converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(helper.exampleValidIngredientMeasurementUnit)
	helper.exampleUpdateInput = converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput(helper.exampleValidIngredientMeasurementUnit)

	helper.service.validIngredientMeasurementUnitIDFetcher = func(*http.Request) string {
		return helper.exampleValidIngredientMeasurementUnit.ID
	}

	helper.service.validIngredientIDFetcher = func(*http.Request) string {
		return helper.exampleValidIngredient.ID
	}

	helper.service.validMeasurementUnitIDFetcher = func(*http.Request) string {
		return helper.exampleValidMeasurementUnit.ID
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
