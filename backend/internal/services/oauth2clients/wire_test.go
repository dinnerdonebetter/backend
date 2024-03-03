package oauth2clients

import (
	"testing"

	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"

	"github.com/stretchr/testify/assert"
)

func TestProvideConfig(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, ProvideConfig(&authservice.Config{}))
	})
}
