package recipes

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/uploads"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			PublicMediaURLPrefix: t.Name(),
			Uploads:              uploads.Config{},
			DataChangesTopicName: "blah",
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
