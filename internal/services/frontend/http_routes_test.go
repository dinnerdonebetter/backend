package frontend

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/routing/chi"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"

	"github.com/stretchr/testify/mock"
)

func TestService_SetupRoutes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		obligatoryHandler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

		authService := &mocktypes.AuthService{}
		authService.On(
			"ServiceAdminMiddleware",
			mock.IsType(obligatoryHandler),
		).Return(http.Handler(obligatoryHandler))

		authService.On(
			"PermissionFilterMiddleware",
			mock.IsType([]authorization.Permission{}),
		).Return(func(next http.Handler) http.Handler { return http.Handler(obligatoryHandler) })

		authService.On(
			"UserAttributionMiddleware",
			mock.IsType(obligatoryHandler),
		).Return(http.Handler(obligatoryHandler))
		s.service.authService = authService

		router := chi.NewRouter(logging.NewNoopLogger())

		s.service.SetupRoutes(router)
	})
}
