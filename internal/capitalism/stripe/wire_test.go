package stripe

import (
	"testing"

	"github.com/prixfixeco/api_server/internal/capitalism"

	"github.com/stretchr/testify/assert"
)

func TestProvideAPIKey(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		cfg := &capitalism.StripeConfig{
			APIKey: t.Name(),
		}

		assert.NotEmpty(t, ProvideAPIKey(cfg))
	})
}

func TestProvideWebhookSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		cfg := &capitalism.StripeConfig{
			WebhookSecret: t.Name(),
		}

		assert.NotEmpty(t, ProvideWebhookSecret(cfg))
	})
}
