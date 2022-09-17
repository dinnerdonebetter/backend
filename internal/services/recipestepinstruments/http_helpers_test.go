package recipestepinstruments

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

type recipeStepInstrumentsServiceHTTPRoutesTestHelper struct {
	ctx                         context.Context
	req                         *http.Request
	res                         *httptest.ResponseRecorder
	service                     *service
	exampleUser                 *types.User
	exampleHousehold            *types.Household
	exampleRecipe               *types.Recipe
	exampleRecipeStep           *types.RecipeStep
	exampleRecipeStepInstrument *types.RecipeStepInstrument
	exampleCreationInput        *types.RecipeStepInstrumentCreationRequestInput
	exampleUpdateInput          *types.RecipeStepInstrumentUpdateRequestInput
}

func buildTestHelper(t *testing.T) *recipeStepInstrumentsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &recipeStepInstrumentsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleRecipe = fakes.BuildFakeRecipe()
	helper.exampleRecipe.CreatedByUser = helper.exampleHousehold.ID
	helper.exampleRecipeStep = fakes.BuildFakeRecipeStep()
	helper.exampleRecipeStep.BelongsToRecipe = helper.exampleRecipe.ID
	helper.exampleRecipeStepInstrument = fakes.BuildFakeRecipeStepInstrument()
	helper.exampleRecipeStepInstrument.BelongsToRecipeStep = helper.exampleRecipeStep.ID
	helper.exampleCreationInput = fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(helper.exampleRecipeStepInstrument)
	helper.exampleUpdateInput = fakes.BuildFakeRecipeStepInstrumentUpdateRequestInputFromRecipeStepInstrument(helper.exampleRecipeStepInstrument)

	helper.service.recipeIDFetcher = func(*http.Request) string {
		return helper.exampleRecipe.ID
	}

	helper.service.recipeStepIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStep.ID
	}

	helper.service.recipeStepInstrumentIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStepInstrument.ID
	}

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
			UserID:                   helper.exampleUser.ID,
			AccountStatus:            helper.exampleUser.AccountStatus,
			AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
			ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRoles...),
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
