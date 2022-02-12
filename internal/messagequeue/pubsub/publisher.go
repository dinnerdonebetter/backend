package pubsub

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

type (
	messagePublisher interface {
		Publish(context.Context, *pubsub.Message) *pubsub.PublishResult
	}

	publisher struct {
		tracer    tracing.Tracer
		encoder   encoding.ClientEncoder
		logger    logging.Logger
		publisher messagePublisher
		topic     string
	}
)

func (r *publisher) Publish(ctx context.Context, data interface{}) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	logger.Debug("publishing message")

	var b bytes.Buffer
	if err := r.encoder.Encode(ctx, &b, data); err != nil {
		return observability.PrepareError(err, r.logger, span, "encoding topic message")
	}

	result := r.publisher.Publish(ctx, &pubsub.Message{Data: b.Bytes()})
	go func(res *pubsub.PublishResult) {
		// The Get method blocks until a server-generated ID or an error is returned for the published message.
		id, err := res.Get(ctx)
		if err != nil {
			// Error handling code can be added here.
			observability.AcknowledgeError(err, logger, span, "publishing pubsub message")
			return
		}

		_ = id
	}(result)

	return nil
}

// providePubSubPublisher provides a sqs-backed publisher.
func providePubSubPublisher(logger logging.Logger, pubsubClient *pubsub.Topic, tracerProvider tracing.TracerProvider, topic string) *publisher {
	return &publisher{
		topic:     topic,
		encoder:   encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		logger:    logging.EnsureLogger(logger),
		publisher: pubsubClient,
		tracer:    tracing.NewTracer(tracerProvider.Tracer(fmt.Sprintf("%s_publisher", topic))),
	}
}

type publisherProvider struct {
	logger            logging.Logger
	publisherCache    map[string]messagequeue.Publisher
	pubsubClient      *pubsub.Client
	tracerProvider    tracing.TracerProvider
	publisherCacheHat sync.RWMutex
}

// ProvidePubSubPublisherProvider returns a PublisherProvider for a given address.
func ProvidePubSubPublisherProvider(logger logging.Logger, tracerProvider tracing.TracerProvider, client *pubsub.Client) messagequeue.PublisherProvider {
	return &publisherProvider{
		logger:         logging.EnsureLogger(logger),
		pubsubClient:   client,
		publisherCache: map[string]messagequeue.Publisher{},
		tracerProvider: tracerProvider,
	}
}

// ProviderPublisher returns a publisher for a given topic.
func (p *publisherProvider) ProviderPublisher(topic string) (messagequeue.Publisher, error) {
	logger := logging.EnsureLogger(p.logger).WithValue("topic", topic)

	p.publisherCacheHat.Lock()
	defer p.publisherCacheHat.Unlock()
	if cachedPub, ok := p.publisherCache[topic]; ok {
		return cachedPub, nil
	}

	t := p.pubsubClient.Topic(topic)

	pub := providePubSubPublisher(logger, t, p.tracerProvider, topic)
	p.publisherCache[topic] = pub

	return pub, nil
}