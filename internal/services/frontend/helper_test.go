package frontend

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	mockrouting "gitlab.com/prixfixe/prixfixe/internal/routing/mock"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"
)

type serviceHTTPRoutesTestHelper struct {
	ctx              context.Context
	req              *http.Request
	res              *httptest.ResponseRecorder
	service          *service
	sessionCtxData   *types.SessionContextData
	exampleUser      *types.User
	exampleHousehold *types.Household
}

func buildTestHelper(t *testing.T) *serviceHTTPRoutesTestHelper {
	t.Helper()

	helper := &serviceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID

	cfg := &Config{}
	logger := logging.NewNoopLogger()
	authService := &mocktypes.AuthService{}

	rpm := mockrouting.NewRouteParamManager()

	rpm.On(
		"BuildRouteParamIDFetcher",
		mock.IsType(logging.NewNoopLogger()),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(func(*http.Request) uint64 { return 0 })

	var ok bool
	helper.service, ok = ProvideService(
		cfg,
		logger,
		authService,
	).(*service)
	require.True(t, ok)

	helper.sessionCtxData = &types.SessionContextData{
		Requester: types.RequesterInfo{
			UserID:                helper.exampleUser.ID,
			Reputation:            helper.exampleUser.ServiceHouseholdStatus,
			ReputationExplanation: helper.exampleUser.ReputationExplanation,
			ServicePermissions:    authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRoles...),
		},
		ActiveHouseholdID: helper.exampleHousehold.ID,
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
			helper.exampleHousehold.ID: authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdMemberRole.String()),
		},
	}

	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return helper.sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), types.SessionContextDataKey, helper.sessionCtxData))

	helper.res = httptest.NewRecorder()

	return helper
}
