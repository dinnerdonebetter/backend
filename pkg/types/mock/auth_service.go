package mocktypes

import (
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.AuthService = (*AuthService)(nil)

// AuthService is a mock types.AuthService.
type AuthService struct {
	mock.Mock
}

// StatusHandler satisfies our interface contract.
func (m *AuthService) StatusHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(req, res)
}

// PermissionFilterMiddleware satisfies our interface contract.
func (m *AuthService) PermissionFilterMiddleware(perms ...authorization.Permission) func(next http.Handler) http.Handler {
	return m.Called(perms).Get(0).(func(http.Handler) http.Handler)
}

// BeginSessionHandler satisfies our interface contract.
func (m *AuthService) BeginSessionHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(req, res)
}

// EndSessionHandler satisfies our interface contract.
func (m *AuthService) EndSessionHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(req, res)
}

// CycleCookieSecretHandler satisfies our interface contract.
func (m *AuthService) CycleCookieSecretHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(req, res)
}

// PASETOHandler satisfies our interface contract.
func (m *AuthService) PASETOHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(req, res)
}

// ChangeActiveHouseholdHandler satisfies our interface contract.
func (m *AuthService) ChangeActiveHouseholdHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(req, res)
}

// CookieRequirementMiddleware satisfies our interface contract.
func (m *AuthService) CookieRequirementMiddleware(next http.Handler) http.Handler {
	return m.Called(next).Get(0).(http.Handler)
}

// UserAttributionMiddleware satisfies our interface contract.
func (m *AuthService) UserAttributionMiddleware(next http.Handler) http.Handler {
	return m.Called(next).Get(0).(http.Handler)
}

// AuthorizationMiddleware satisfies our interface contract.
func (m *AuthService) AuthorizationMiddleware(next http.Handler) http.Handler {
	return m.Called(next).Get(0).(http.Handler)
}

// ServiceAdminMiddleware satisfies our interface contract.
func (m *AuthService) ServiceAdminMiddleware(next http.Handler) http.Handler {
	return m.Called(next).Get(0).(http.Handler)
}

// UserLoginInputMiddleware satisfies our interface contract.
func (m *AuthService) UserLoginInputMiddleware(next http.Handler) http.Handler {
	return m.Called(next).Get(0).(http.Handler)
}

// PASETOCreationInputMiddleware satisfies our interface contract.
func (m *AuthService) PASETOCreationInputMiddleware(next http.Handler) http.Handler {
	return m.Called(next).Get(0).(http.Handler)
}

// ChangeActiveHouseholdInputMiddleware satisfies our interface contract.
func (m *AuthService) ChangeActiveHouseholdInputMiddleware(next http.Handler) http.Handler {
	return m.Called(next).Get(0).(http.Handler)
}
