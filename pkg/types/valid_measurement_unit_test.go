package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidMeasurementUnit_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnit{
			Imperial: true,
		}
		input := &ValidMeasurementUnitUpdateRequestInput{}

		fake.Struct(&input)
		input.Volumetric = pointers.Pointer(true)
		input.Universal = pointers.Pointer(true)
		input.Imperial = pointers.Pointer(false)
		input.Metric = pointers.Pointer(true)

		x.Update(input)
	})
}

func TestValidMeasurementUnitCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnitCreationRequestInput{
			Name:        t.Name(),
			Description: t.Name(),
			Volumetric:  fake.Bool(),
			IconPath:    t.Name(),
			Universal:   fake.Bool(),
			Metric:      true,
			Imperial:    false,
			PluralName:  t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnitCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidMeasurementUnitUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnitUpdateRequestInput{
			Name:        pointers.Pointer(t.Name()),
			Description: pointers.Pointer(t.Name()),
			Volumetric:  pointers.Pointer(fake.Bool()),
			IconPath:    pointers.Pointer(t.Name()),
			Universal:   pointers.Pointer(fake.Bool()),
			Metric:      pointers.Pointer(fake.Bool()),
			Imperial:    pointers.Pointer(fake.Bool()),
			PluralName:  pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnitUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
