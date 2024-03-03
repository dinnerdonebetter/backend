package types

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvert(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		original := ValidPreparation{
			Name:           "Name",
			Description:    "Description",
			OnlyForVessels: true,
		}
		updateInput := ValidPreparationUpdateRequestInput{
			Slug:      pointer.To("Slug"),
			PastTense: pointer.To("PastTense"),
		}
		updated := ValidPreparation{
			Name:           "Name",
			Description:    "Description",
			Slug:           "Slug",
			PastTense:      "PastTense",
			OnlyForVessels: true,
		}

		require.NoError(t, Merge(&original, &updateInput))

		assert.Equal(t, updated, original)
	})
}
