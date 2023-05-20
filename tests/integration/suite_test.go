package integration

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	cookieAuthType = "cookie"
	pasetoAuthType = "PASETO"
)

var (
	globalClientExceptions []string
)

type testClientWrapper struct {
	user     *apiclient.Client
	admin    *apiclient.Client
	authType string
}

func TestIntegration(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(TestSuite))
}

type TestSuite struct {
	suite.Suite

	ctx    context.Context
	user   *types.User
	cookie *http.Cookie
	cookieClient,
	pasetoClient,
	adminCookieClient,
	adminPASETOClient *apiclient.Client
}

var _ suite.SetupTestSuite = (*TestSuite)(nil)

func (s *TestSuite) SetupTest() {
	t := s.T()
	testName := t.Name()

	ctx, span := tracing.StartCustomSpan(context.Background(), testName)
	defer span.End()

	s.ctx, _ = tracing.StartCustomSpan(ctx, testName)
	s.user, s.cookie, s.cookieClient, s.pasetoClient = createUserAndClientForTest(s.ctx, t, nil)
	s.adminCookieClient, s.adminPASETOClient = buildAdminCookieAndPASETOClients(s.ctx, t)
}

func (s *TestSuite) runForCookieClient(name string, subtestBuilder func(*testClientWrapper) func()) {
	for a, c := range s.eachClientExcept(pasetoAuthType) {
		authType, testClients := a, c
		s.Run(fmt.Sprintf("%s via %s", name, authType), subtestBuilder(testClients))
	}
}

func (s *TestSuite) runForPASETOClient(name string, subtestBuilder func(*testClientWrapper) func()) {
	if x, _ := strconv.ParseBool(os.Getenv("SKIP_PASETO_TESTS")); x {
		return
	}

	for a, c := range s.eachClientExcept(cookieAuthType) {
		authType, testClients := a, c
		s.Run(fmt.Sprintf("%s via %s", name, authType), subtestBuilder(testClients))
	}
}

func (s *TestSuite) runForEachClient(name string, subtestBuilder func(*testClientWrapper) func()) {
	for a, c := range s.eachClientExcept() {
		authType, testClients := a, c
		s.Run(fmt.Sprintf("%s via %s", name, authType), subtestBuilder(testClients))
	}
}

func (s *TestSuite) runForEachClientExcept(name string, subtestBuilder func(*testClientWrapper) func(), exceptions ...string) {
	for a, c := range s.eachClientExcept(exceptions...) {
		authType, testClients := a, c
		s.Run(fmt.Sprintf("%s via %s", name, authType), subtestBuilder(testClients))
	}
}

func (s *TestSuite) eachClientExcept(exceptions ...string) map[string]*testClientWrapper {
	t := s.T()

	clients := map[string]*testClientWrapper{
		cookieAuthType: {authType: cookieAuthType, user: s.cookieClient, admin: s.adminCookieClient},
		pasetoAuthType: {authType: pasetoAuthType, user: s.pasetoClient, admin: s.adminPASETOClient},
	}

	for _, name := range exceptions {
		delete(clients, name)
	}

	for _, name := range globalClientExceptions {
		delete(clients, name)
	}

	require.NotEmpty(t, clients)

	return clients
}
