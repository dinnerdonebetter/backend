package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertValidIngredientToValidIngredientUpdateRequestInput creates a ValidIngredientUpdateRequestInput from a ValidIngredient.
func ConvertValidIngredientToValidIngredientUpdateRequestInput(x *types.ValidIngredient) *types.ValidIngredientUpdateRequestInput {
	out := &types.ValidIngredientUpdateRequestInput{
		Name:                                    &x.Name,
		Description:                             &x.Description,
		Warning:                                 &x.Warning,
		IconPath:                                &x.IconPath,
		ContainsDairy:                           &x.ContainsDairy,
		ContainsPeanut:                          &x.ContainsPeanut,
		ContainsTreeNut:                         &x.ContainsTreeNut,
		ContainsEgg:                             &x.ContainsEgg,
		ContainsWheat:                           &x.ContainsWheat,
		ContainsShellfish:                       &x.ContainsShellfish,
		ContainsSesame:                          &x.ContainsSesame,
		ContainsFish:                            &x.ContainsFish,
		ContainsGluten:                          &x.ContainsGluten,
		AnimalFlesh:                             &x.AnimalFlesh,
		IsMeasuredVolumetrically:                &x.IsMeasuredVolumetrically,
		IsLiquid:                                &x.IsLiquid,
		ContainsSoy:                             &x.ContainsSoy,
		PluralName:                              &x.PluralName,
		AnimalDerived:                           &x.AnimalDerived,
		RestrictToPreparations:                  &x.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: x.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: x.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     &x.StorageInstructions,
		Slug:                                    &x.Slug,
		ContainsAlcohol:                         &x.ContainsAlcohol,
		ShoppingSuggestions:                     &x.ShoppingSuggestions,
		IsStarch:                                &x.IsStarch,
		IsProtein:                               &x.IsProtein,
		IsGrain:                                 &x.IsGrain,
		IsFruit:                                 &x.IsFruit,
		IsSalt:                                  &x.IsSalt,
		IsFat:                                   &x.IsFat,
		IsAcid:                                  &x.IsAcid,
		IsHeat:                                  &x.IsHeat,
	}

	return out
}

// ConvertValidIngredientCreationRequestInputToValidIngredientDatabaseCreationInput creates a DatabaseCreationInput from a ValidIngredientCreationRequestInput.
func ConvertValidIngredientCreationRequestInputToValidIngredientDatabaseCreationInput(x *types.ValidIngredientCreationRequestInput) *types.ValidIngredientDatabaseCreationInput {
	out := &types.ValidIngredientDatabaseCreationInput{
		ID:                                      identifiers.New(),
		Name:                                    x.Name,
		Description:                             x.Description,
		Warning:                                 x.Warning,
		ContainsEgg:                             x.ContainsEgg,
		ContainsDairy:                           x.ContainsDairy,
		ContainsPeanut:                          x.ContainsPeanut,
		ContainsTreeNut:                         x.ContainsTreeNut,
		ContainsSoy:                             x.ContainsSoy,
		ContainsWheat:                           x.ContainsWheat,
		ContainsShellfish:                       x.ContainsShellfish,
		ContainsSesame:                          x.ContainsSesame,
		ContainsFish:                            x.ContainsFish,
		ContainsGluten:                          x.ContainsGluten,
		AnimalFlesh:                             x.AnimalFlesh,
		IsMeasuredVolumetrically:                x.IsMeasuredVolumetrically,
		IsLiquid:                                x.IsLiquid,
		IconPath:                                x.IconPath,
		PluralName:                              x.PluralName,
		AnimalDerived:                           x.AnimalDerived,
		RestrictToPreparations:                  x.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: x.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: x.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     x.StorageInstructions,
		Slug:                                    x.Slug,
		ContainsAlcohol:                         x.ContainsAlcohol,
		ShoppingSuggestions:                     x.ShoppingSuggestions,
		IsStarch:                                x.IsStarch,
		IsProtein:                               x.IsProtein,
		IsGrain:                                 x.IsGrain,
		IsFruit:                                 x.IsFruit,
		IsSalt:                                  x.IsSalt,
		IsFat:                                   x.IsFat,
		IsAcid:                                  x.IsAcid,
		IsHeat:                                  x.IsHeat,
	}

	return out
}

// ConvertValidIngredientToValidIngredientCreationRequestInput builds a ValidIngredientCreationRequestInput from a Ingredient.
func ConvertValidIngredientToValidIngredientCreationRequestInput(x *types.ValidIngredient) *types.ValidIngredientCreationRequestInput {
	return &types.ValidIngredientCreationRequestInput{
		Name:                                    x.Name,
		Description:                             x.Description,
		Warning:                                 x.Warning,
		ContainsEgg:                             x.ContainsEgg,
		ContainsDairy:                           x.ContainsDairy,
		ContainsPeanut:                          x.ContainsPeanut,
		ContainsTreeNut:                         x.ContainsTreeNut,
		ContainsSoy:                             x.ContainsSoy,
		ContainsWheat:                           x.ContainsWheat,
		ContainsShellfish:                       x.ContainsShellfish,
		ContainsSesame:                          x.ContainsSesame,
		ContainsFish:                            x.ContainsFish,
		ContainsGluten:                          x.ContainsGluten,
		AnimalFlesh:                             x.AnimalFlesh,
		IsMeasuredVolumetrically:                x.IsMeasuredVolumetrically,
		IsLiquid:                                x.IsLiquid,
		IconPath:                                x.IconPath,
		PluralName:                              x.PluralName,
		AnimalDerived:                           x.AnimalDerived,
		RestrictToPreparations:                  x.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: x.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: x.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     x.StorageInstructions,
		Slug:                                    x.Slug,
		ContainsAlcohol:                         x.ContainsAlcohol,
		ShoppingSuggestions:                     x.ShoppingSuggestions,
		IsStarch:                                x.IsStarch,
		IsProtein:                               x.IsProtein,
		IsGrain:                                 x.IsGrain,
		IsFruit:                                 x.IsFruit,
		IsSalt:                                  x.IsSalt,
		IsFat:                                   x.IsFat,
		IsAcid:                                  x.IsAcid,
		IsHeat:                                  x.IsHeat,
	}
}

// ConvertValidIngredientToValidIngredientDatabaseCreationInput builds a ValidIngredientDatabaseCreationInput from a ValidIngredient.
func ConvertValidIngredientToValidIngredientDatabaseCreationInput(x *types.ValidIngredient) *types.ValidIngredientDatabaseCreationInput {
	return &types.ValidIngredientDatabaseCreationInput{
		ID:                                      x.ID,
		Name:                                    x.Name,
		Description:                             x.Description,
		Warning:                                 x.Warning,
		ContainsEgg:                             x.ContainsEgg,
		ContainsDairy:                           x.ContainsDairy,
		ContainsPeanut:                          x.ContainsPeanut,
		ContainsTreeNut:                         x.ContainsTreeNut,
		ContainsSoy:                             x.ContainsSoy,
		ContainsWheat:                           x.ContainsWheat,
		ContainsShellfish:                       x.ContainsShellfish,
		ContainsSesame:                          x.ContainsSesame,
		ContainsFish:                            x.ContainsFish,
		ContainsGluten:                          x.ContainsGluten,
		AnimalFlesh:                             x.AnimalFlesh,
		IsMeasuredVolumetrically:                x.IsMeasuredVolumetrically,
		IsLiquid:                                x.IsLiquid,
		IconPath:                                x.IconPath,
		PluralName:                              x.PluralName,
		AnimalDerived:                           x.AnimalDerived,
		RestrictToPreparations:                  x.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: x.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: x.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     x.StorageInstructions,
		Slug:                                    x.Slug,
		ContainsAlcohol:                         x.ContainsAlcohol,
		ShoppingSuggestions:                     x.ShoppingSuggestions,
		IsStarch:                                x.IsStarch,
		IsProtein:                               x.IsProtein,
		IsGrain:                                 x.IsGrain,
		IsFruit:                                 x.IsFruit,
		IsSalt:                                  x.IsSalt,
		IsFat:                                   x.IsFat,
		IsAcid:                                  x.IsAcid,
		IsHeat:                                  x.IsHeat,
	}
}

// ConvertNullableValidIngredientToValidIngredient converts a NullableValidIngredient to a ValidIngredient.
func ConvertNullableValidIngredientToValidIngredient(x *types.NullableValidIngredient) *types.ValidIngredient {
	return &types.ValidIngredient{
		CreatedAt:                               *x.CreatedAt,
		LastUpdatedAt:                           x.LastUpdatedAt,
		ArchivedAt:                              x.ArchivedAt,
		ID:                                      *x.ID,
		Warning:                                 *x.Warning,
		Description:                             *x.Description,
		IconPath:                                *x.IconPath,
		PluralName:                              *x.PluralName,
		StorageInstructions:                     *x.StorageInstructions,
		Name:                                    *x.Name,
		MaximumIdealStorageTemperatureInCelsius: x.MaximumIdealStorageTemperatureInCelsius,
		MinimumIdealStorageTemperatureInCelsius: x.MinimumIdealStorageTemperatureInCelsius,
		ContainsShellfish:                       *x.ContainsShellfish,
		ContainsDairy:                           *x.ContainsDairy,
		AnimalFlesh:                             *x.AnimalFlesh,
		IsMeasuredVolumetrically:                *x.IsMeasuredVolumetrically,
		IsLiquid:                                *x.IsLiquid,
		ContainsPeanut:                          *x.ContainsPeanut,
		ContainsTreeNut:                         *x.ContainsTreeNut,
		ContainsEgg:                             *x.ContainsEgg,
		ContainsWheat:                           *x.ContainsWheat,
		ContainsSoy:                             *x.ContainsSoy,
		AnimalDerived:                           *x.AnimalDerived,
		RestrictToPreparations:                  *x.RestrictToPreparations,
		ContainsSesame:                          *x.ContainsSesame,
		ContainsFish:                            *x.ContainsFish,
		ContainsGluten:                          *x.ContainsGluten,
		Slug:                                    *x.Slug,
		ContainsAlcohol:                         *x.ContainsAlcohol,
		ShoppingSuggestions:                     *x.ShoppingSuggestions,
		IsStarch:                                *x.IsStarch,
		IsProtein:                               *x.IsProtein,
		IsGrain:                                 *x.IsGrain,
		IsFruit:                                 *x.IsFruit,
		IsSalt:                                  *x.IsSalt,
		IsFat:                                   *x.IsFat,
		IsAcid:                                  *x.IsAcid,
		IsHeat:                                  *x.IsHeat,
	}
}

// ConvertValidIngredientToNullableValidIngredient converts a NullableValidIngredient to a ValidIngredient.
func ConvertValidIngredientToNullableValidIngredient(x *types.ValidIngredient) *types.NullableValidIngredient {
	return &types.NullableValidIngredient{
		CreatedAt:                               &x.CreatedAt,
		LastUpdatedAt:                           x.LastUpdatedAt,
		ArchivedAt:                              x.ArchivedAt,
		ID:                                      &x.ID,
		Warning:                                 &x.Warning,
		Description:                             &x.Description,
		IconPath:                                &x.IconPath,
		PluralName:                              &x.PluralName,
		StorageInstructions:                     &x.StorageInstructions,
		Name:                                    &x.Name,
		MaximumIdealStorageTemperatureInCelsius: x.MaximumIdealStorageTemperatureInCelsius,
		MinimumIdealStorageTemperatureInCelsius: x.MinimumIdealStorageTemperatureInCelsius,
		ContainsShellfish:                       &x.ContainsShellfish,
		ContainsDairy:                           &x.ContainsDairy,
		AnimalFlesh:                             &x.AnimalFlesh,
		IsMeasuredVolumetrically:                &x.IsMeasuredVolumetrically,
		IsLiquid:                                &x.IsLiquid,
		ContainsPeanut:                          &x.ContainsPeanut,
		ContainsTreeNut:                         &x.ContainsTreeNut,
		ContainsEgg:                             &x.ContainsEgg,
		ContainsWheat:                           &x.ContainsWheat,
		ContainsSoy:                             &x.ContainsSoy,
		AnimalDerived:                           &x.AnimalDerived,
		RestrictToPreparations:                  &x.RestrictToPreparations,
		ContainsSesame:                          &x.ContainsSesame,
		ContainsFish:                            &x.ContainsFish,
		ContainsGluten:                          &x.ContainsGluten,
		Slug:                                    &x.Slug,
		ContainsAlcohol:                         &x.ContainsAlcohol,
		ShoppingSuggestions:                     &x.ShoppingSuggestions,
		IsStarch:                                &x.IsStarch,
		IsProtein:                               &x.IsProtein,
		IsGrain:                                 &x.IsGrain,
		IsFruit:                                 &x.IsFruit,
		IsSalt:                                  &x.IsSalt,
		IsFat:                                   &x.IsFat,
		IsAcid:                                  &x.IsAcid,
		IsHeat:                                  &x.IsHeat,
	}
}

// ConvertValidIngredientToValidIngredientSearchSubset converts a ValidIngredient to a ValidIngredientSearchSubset.
func ConvertValidIngredientToValidIngredientSearchSubset(x *types.ValidIngredient) *types.ValidIngredientSearchSubset {
	return &types.ValidIngredientSearchSubset{
		ID:                  x.ID,
		Name:                x.Name,
		PluralName:          x.PluralName,
		Description:         x.Description,
		ShoppingSuggestions: x.ShoppingSuggestions,
	}
}
