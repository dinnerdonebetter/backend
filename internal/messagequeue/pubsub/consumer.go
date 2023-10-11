package pubsub

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"cloud.google.com/go/pubsub"
)

type (
	messageConsumer interface {
		CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error)
	}

	consumer struct {
		tracer      tracing.Tracer
		encoder     encoding.ClientEncoder
		logger      logging.Logger
		consumer    messageConsumer
		handlerFunc func(context.Context, []byte) error
		topic       string
	}
)

// buildPubSubConsumer provides a Pub/Sub-backed consumer.
func buildPubSubConsumer(
	logger logging.Logger,
	pubsubClient *pubsub.Client,
	tracerProvider tracing.TracerProvider,
	topic string,
	handlerFunc func(context.Context, []byte) error,
) messagequeue.Consumer {
	return &consumer{
		topic:       topic,
		encoder:     encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		logger:      logging.EnsureLogger(logger),
		consumer:    pubsubClient,
		handlerFunc: handlerFunc,
		tracer:      tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(fmt.Sprintf("%s_consumer", topic))),
	}
}

func (p *consumer) Consume(stopChan chan bool, errors chan error) {
	if stopChan == nil {
		stopChan = make(chan bool, 1)
	}

	ctx := context.Background()
	sub, err := p.consumer.CreateSubscription(ctx, p.topic, pubsub.SubscriptionConfig{})
	if err != nil {
		errors <- err
		return
	}

	go func() {
		<-stopChan
		if err = sub.Delete(ctx); err != nil {
			errors <- err
		}
	}()

	if err = sub.Receive(ctx, func(receivedContext context.Context, m *pubsub.Message) {
		if handleErr := p.handlerFunc(receivedContext, m.Data); handleErr != nil {
			errors <- err
		} else {
			m.Ack()
		}
	}); err != nil {
		log.Println(err)
	}
}

type consumerProvider struct {
	logger           logging.Logger
	consumerCache    map[string]messagequeue.Consumer
	pubsubClient     *pubsub.Client
	tracerProvider   tracing.TracerProvider
	consumerCacheHat sync.RWMutex
}

// ProvidePubSubConsumerProvider returns a ConsumerProvider for a given address.
func ProvidePubSubConsumerProvider(logger logging.Logger, tracerProvider tracing.TracerProvider, client *pubsub.Client) messagequeue.ConsumerProvider {
	return &consumerProvider{
		logger:         logging.EnsureLogger(logger),
		pubsubClient:   client,
		consumerCache:  map[string]messagequeue.Consumer{},
		tracerProvider: tracerProvider,
	}
}

// Close closes the connection topic.
func (p *consumerProvider) Close() {
	if err := p.pubsubClient.Close(); err != nil {
		p.logger.Error(err, "closing pubsub connection")
	}
}

// ProvideConsumer returns a consumer for a given topic.
func (p *consumerProvider) ProvideConsumer(ctx context.Context, topic string, handlerFunc func(context.Context, []byte) error) (messagequeue.Consumer, error) {
	if topic == "" {
		return nil, messagequeue.ErrEmptyTopicName
	}

	logger := logging.EnsureLogger(p.logger.Clone())

	p.consumerCacheHat.Lock()
	defer p.consumerCacheHat.Unlock()
	if cachedPub, ok := p.consumerCache[topic]; ok {
		return cachedPub, nil
	}

	pub := buildPubSubConsumer(logger, p.pubsubClient, p.tracerProvider, topic, handlerFunc)
	p.consumerCache[topic] = pub

	return pub, nil
}
