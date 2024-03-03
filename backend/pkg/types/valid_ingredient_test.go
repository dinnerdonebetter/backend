package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := &ValidIngredient{
			MaximumIdealStorageTemperatureInCelsius: nil,
			MinimumIdealStorageTemperatureInCelsius: nil,
			IconPath:                                "",
			Warning:                                 "",
			PluralName:                              "",
			StorageInstructions:                     "",
			Name:                                    "",
			Description:                             "",
			Slug:                                    "",
			ShoppingSuggestions:                     "",
			ContainsShellfish:                       false,
			IsMeasuredVolumetrically:                false,
			IsLiquid:                                false,
			ContainsPeanut:                          false,
			ContainsTreeNut:                         false,
			ContainsEgg:                             false,
			ContainsWheat:                           false,
			ContainsSoy:                             false,
			AnimalDerived:                           false,
			RestrictToPreparations:                  false,
			ContainsSesame:                          false,
			ContainsFish:                            false,
			ContainsGluten:                          false,
			ContainsDairy:                           false,
			ContainsAlcohol:                         false,
			AnimalFlesh:                             false,
		}

		input := &ValidIngredientUpdateRequestInput{}
		assert.NoError(t, fake.Struct(&input))
		input.ContainsEgg = pointer.To(true)
		input.ContainsDairy = pointer.To(true)
		input.ContainsPeanut = pointer.To(true)
		input.ContainsTreeNut = pointer.To(true)
		input.ContainsSoy = pointer.To(true)
		input.ContainsWheat = pointer.To(true)
		input.ContainsShellfish = pointer.To(true)
		input.ContainsSesame = pointer.To(true)
		input.ContainsFish = pointer.To(true)
		input.ContainsGluten = pointer.To(true)
		input.AnimalFlesh = pointer.To(true)
		input.IsMeasuredVolumetrically = pointer.To(true)
		input.AnimalDerived = pointer.To(true)
		input.RestrictToPreparations = pointer.To(true)
		input.ContainsAlcohol = pointer.To(true)
		input.IsLiquid = pointer.To(true)
		input.IsProtein = pointer.To(true)
		input.IsStarch = pointer.To(true)
		input.IsGrain = pointer.To(true)
		input.IsFruit = pointer.To(true)
		input.IsSalt = pointer.To(true)
		input.IsFat = pointer.To(true)
		input.IsAcid = pointer.To(true)
		input.IsHeat = pointer.To(true)

		actual.Update(input)
	})
}

func TestValidIngredientCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidIngredientCreationRequestInput{
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidIngredientDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidIngredientDatabaseCreationInput{
			ID:   t.Name(),
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidIngredientUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidIngredientUpdateRequestInput{
			Name: pointer.To(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
