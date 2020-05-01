package mock

import (
	"net/http"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RequiredPreparationInstrumentDataServer = (*RequiredPreparationInstrumentDataServer)(nil)

// RequiredPreparationInstrumentDataServer is a mocked models.RequiredPreparationInstrumentDataServer for testing.
type RequiredPreparationInstrumentDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements.
func (m *RequiredPreparationInstrumentDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements.
func (m *RequiredPreparationInstrumentDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// ListHandler implements our interface requirements.
func (m *RequiredPreparationInstrumentDataServer) ListHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}

// CreateHandler implements our interface requirements.
func (m *RequiredPreparationInstrumentDataServer) CreateHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}

// ExistenceHandler implements our interface requirements.
func (m *RequiredPreparationInstrumentDataServer) ExistenceHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}

// ReadHandler implements our interface requirements.
func (m *RequiredPreparationInstrumentDataServer) ReadHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}

// UpdateHandler implements our interface requirements.
func (m *RequiredPreparationInstrumentDataServer) UpdateHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}

// ArchiveHandler implements our interface requirements.
func (m *RequiredPreparationInstrumentDataServer) ArchiveHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}
