package websockets

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/gorilla/websocket"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName = "websockets_service"
)

type (
	websocketConnection interface {
		SetWriteDeadline(t time.Time) error
		WriteMessage(messageType int, data []byte) error
		WriteControl(messageType int, data []byte, deadline time.Time) error
	}

	// service handles websockets.
	service struct {
		logger                      logging.Logger
		encoderDecoder              encoding.ServerEncoderDecoder
		tracer                      tracing.Tracer
		connections                 map[string][]websocketConnection
		sessionContextDataFetcher   func(*http.Request) (*types.SessionContextData, error)
		authConfig                  *authservice.Config
		websocketConnectionUpgrader websocket.Upgrader
		websocketDeadline           time.Duration
		pollDuration                time.Duration
		connectionsHat              sync.RWMutex
	}
)

// ProvideService builds a new websocket service.
func ProvideService(
	ctx context.Context,
	authCfg *authservice.Config,
	cfg Config,
	logger logging.Logger,
	encoder encoding.ServerEncoderDecoder,
	consumerProvider messagequeue.ConsumerProvider,
	tracerProvider tracing.TracerProvider,
) (types.WebsocketDataService, error) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 10 * time.Second,
		Error:            buildWebsocketErrorFunc(encoder),
	}

	svc := &service{
		logger:                      logging.EnsureLogger(logger).WithName(serviceName),
		sessionContextDataFetcher:   authservice.FetchContextFromRequest,
		encoderDecoder:              encoder,
		websocketConnectionUpgrader: upgrader,
		connections:                 map[string][]websocketConnection{},
		websocketDeadline:           5 * time.Second,
		pollDuration:                10 * time.Second,
		authConfig:                  authCfg,
		tracer:                      tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	svc.logger.WithValue("topic_name", cfg.DataChangesTopicName).Info("fetching data change thing")

	dataChangesConsumer, err := consumerProvider.ProvideConsumer(ctx, cfg.DataChangesTopicName, svc.handleDataChange)
	if err != nil {
		return nil, fmt.Errorf("setting up event publisher: %w", err)
	}

	go svc.pingConnections()
	go dataChangesConsumer.Consume(nil, nil)

	return svc, nil
}

func buildWebsocketErrorFunc(encoder encoding.ServerEncoderDecoder) func(http.ResponseWriter, *http.Request, int, error) {
	return func(res http.ResponseWriter, req *http.Request, status int, reason error) {
		encoder.EncodeErrorResponse(req.Context(), res, reason.Error(), status)
	}
}
