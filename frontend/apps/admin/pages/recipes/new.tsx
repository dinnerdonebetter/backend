import {
  ActionIcon,
  Button,
  Card,
  Divider,
  Grid,
  Text,
  Stack,
  Textarea,
  TextInput,
  Autocomplete,
  AutocompleteItem,
  SelectItem,
  NumberInput,
  Select,
  Switch,
  Space,
  Collapse,
  Box,
  Title,
  Center,
} from '@mantine/core';
import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { IconChevronDown, IconChevronUp, IconCircleX, IconEdit, IconEditOff, IconTrash } from '@tabler/icons';
import { FormEvent, useReducer, useState } from 'react';
import { z } from 'zod';

import {
  Recipe,
  RecipeStepIngredient,
  RecipeStepInstrument,
  ValidIngredient,
  ValidIngredientState,
  ValidMeasurementUnit,
  ValidPreparation,
  ValidPreparationInstrument,
  ValidRecipeStepProductType,
  QueryFilteredResult,
  RecipeStepCreationRequestInput,
  RecipeStepInstrumentCreationRequestInput,
  RecipeStepIngredientCreationRequestInput,
  RecipeStepCompletionConditionCreationRequestInput,
  RecipeStepProductCreationRequestInput,
  RecipeCreationRequestInput,
  RecipeStepVesselCreationRequestInput,
  RecipeStepVessel,
  ValidVessel,
} from '@dinnerdonebetter/models';
import {
  determineAvailableRecipeStepProducts,
  determineAvailableRecipeStepVessels,
  determinePreparedInstrumentOptions,
  RecipeStepInstrumentSuggestion,
  RecipeStepProductSuggestion,
  RecipeStepVesselSuggestion,
} from '@dinnerdonebetter/utils';

import { AppLayout } from '../..//src/layouts';
import { buildLocalClient } from '../../src/client';
import { useRecipeCreationReducer, RecipeCreationPageState } from '../../src/reducers';

const validRecipeStepProductTypes = ['ingredient', 'instrument', 'vessel'];

const addingStepsShouldBeDisabled = (pageState: RecipeCreationPageState): boolean => {
  const anyPreparationMissing =
    pageState.recipe.steps.filter((x: RecipeStepCreationRequestInput) => {
      return x.preparationID === '';
    }).length !== 0;

  const latestStep = pageState.recipe.steps[(pageState.recipe.steps || []).length - 1];

  const noProducts = latestStep.products.length === 0;
  const noPreparationSet = anyPreparationMissing || latestStep.preparationID === '';
  const noInstrumentsOrVessels = latestStep.vessels.length > 0 && latestStep.instruments.length === 0;

  const invalidIngredientsPresent =
    latestStep.ingredients.filter((x: RecipeStepIngredientCreationRequestInput) => {
      return x.ingredientID === '' && !x.productOfRecipeStepIndex && !x.productOfRecipeStepProductIndex;
    }).length > 0;

  const validIngredientsMissingMeasurementUnits =
    latestStep.ingredients.filter((x: RecipeStepIngredientCreationRequestInput) => {
      return (
        x.measurementUnitID === '' &&
        x.ingredientID !== '' &&
        !x.productOfRecipeStepIndex &&
        !x.productOfRecipeStepProductIndex
      );
    }).length > 0;

  const invalidInstrumentsPresent =
    latestStep.instruments.filter((x: RecipeStepInstrumentCreationRequestInput) => {
      return x.instrumentID === '' && !x.productOfRecipeStepIndex && !x.productOfRecipeStepProductIndex;
    }).length > 0;

  const invalidVesselsPresent =
    latestStep.vessels.filter((x: RecipeStepVesselCreationRequestInput) => {
      return x.vesselID === '' && !x.productOfRecipeStepIndex && !x.productOfRecipeStepProductIndex;
    }).length > 0;

  const rv =
    noPreparationSet ||
    validIngredientsMissingMeasurementUnits ||
    invalidIngredientsPresent ||
    noInstrumentsOrVessels ||
    invalidInstrumentsPresent ||
    invalidVesselsPresent ||
    noProducts;

  return rv;
};

const recipeCreationFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  minimumEstimatedPortions: z.number().min(1),
  slug: z
    .string()
    .trim()
    .min(0)
    .regex(new RegExp(/^[a-zA-Z\-]{1,}$/gm), 'must match expected URL slug pattern'),
  steps: z
    .array(
      z.object({
        preparationID: z.string().trim().min(1, 'preparation ID is required'),
        instruments: z.array(
          z
            .object({
              name: z.string().trim().min(1, 'instrument name is required'),
              instrumentID: z.string().trim().optional(),
              productOfRecipeStepIndex: z.number().optional(),
              productOfRecipeStepProductIndex: z.number().optional(),
              minimumQuantity: z.number().min(1),
            })
            .refine((instrument) => {
              return (
                Boolean(instrument.instrumentID) ||
                (instrument.productOfRecipeStepIndex !== undefined &&
                  instrument.productOfRecipeStepProductIndex !== undefined)
              );
            }),
        ),
        vessels: z.array(
          z
            .object({
              name: z.string().trim().min(1, 'vessel name is required'),
              instrumentID: z.string().trim().optional(),
              productOfRecipeStepIndex: z.number().optional(),
              productOfRecipeStepProductIndex: z.number().optional(),
              minimumQuantity: z.number().min(1),
            })
            .refine((vessel) => {
              return (
                Boolean(vessel.instrumentID) ||
                (vessel.productOfRecipeStepIndex !== undefined && vessel.productOfRecipeStepProductIndex !== undefined)
              );
            }, 'vessel must either be a valid instrument or a product of a previous step'),
        ),
        ingredients: z.array(
          z
            .object({
              name: z.string().trim().min(1, 'ingredient name is required'),
              ingredientID: z.string().trim().optional(),
              productOfRecipeStepIndex: z.number().optional(),
              productOfRecipeStepProductIndex: z.number().optional(),
              minimumQuantity: z.number().min(0.01),
            })
            .refine((ingredient) => {
              return (
                Boolean(ingredient.ingredientID) ||
                (ingredient.productOfRecipeStepIndex !== undefined &&
                  ingredient.productOfRecipeStepProductIndex !== undefined)
              );
            }),
        ),
        products: z
          .array(
            z.object({
              type: z.enum(['ingredient', 'instrument', 'vessel']), // for some reason, ALL_RECIPE_STEP_PRODUCT_TYPES doesn't work here
              name: z.string().trim().min(1, 'product name is required'),
              minimumQuantity: z.number().min(0.01),
            }),
          )
          .min(1),
      }),
    )
    .min(2),
});

const recipeSubmissionShouldBeDisabled = (recipe: RecipeCreationRequestInput): boolean => {
  const evaluatedResult = recipeCreationFormSchema.safeParse(recipe);

  const rv = !Boolean(evaluatedResult.success);

  return rv;
};

function RecipeCreator() {
  const apiClient = buildLocalClient();
  const router = useRouter();
  const [pageState, dispatchPageEvent] = useReducer(useRecipeCreationReducer, new RecipeCreationPageState());

  const [debugOutput, setDebugOutput] = useState('');

  const submitRecipe = async () => {
    apiClient
      .createRecipe(pageState.recipe)
      .then((res: Recipe) => {
        router.push(`/recipes/${res.id}`);
      })
      .catch((err: AxiosError) => {
        console.error(`Failed to create recipe: ${err}`);
      });
  };

  const addingStepCompletionConditionsShouldBeDisabled = (step: RecipeStepCreationRequestInput): boolean => {
    return (
      step.completionConditions.length > 0 &&
      step.completionConditions[step.completionConditions.length - 1].ingredientState === '' &&
      step.completionConditions[step.completionConditions.length - 1].ingredients.length === 0
    );
  };

  const handlePreparationQueryChange = (stepIndex: number) => async (value: string) => {
    dispatchPageEvent({
      type: 'UPDATE_STEP_PREPARATION_QUERY',
      stepIndex: stepIndex,
      newQuery: value,
    });

    if (value.length > 2) {
      await apiClient
        .searchForValidPreparations(value)
        .then((res: ValidPreparation[]) => {
          dispatchPageEvent({
            type: 'UPDATE_STEP_PREPARATION_SUGGESTIONS',
            stepIndex: stepIndex,
            results: res || [],
          });
        })
        .catch((err: AxiosError) => {
          console.error(`Failed to get preparations: ${err}`);
        });
    }
  };

  const handlePreparationSelection = (stepIndex: number) => (value: AutocompleteItem) => {
    const selectedPreparation = (pageState.stepHelpers[stepIndex].preparationSuggestions || []).find(
      (preparationSuggestion: ValidPreparation) => preparationSuggestion.name === value.value,
    );

    if (!selectedPreparation) {
      console.error(
        `couldn't find preparation to add: ${value.value}, ${JSON.stringify(
          pageState.stepHelpers[stepIndex].preparationSuggestions.map((x: ValidPreparation) => x.name),
        )}`,
      );
      return;
    }

    dispatchPageEvent({
      type: 'UPDATE_STEP_PREPARATION',
      stepIndex: stepIndex,
      selectedPreparation: selectedPreparation,
    });

    apiClient
      .validPreparationInstrumentsForPreparationID(selectedPreparation.id)
      .then((res: QueryFilteredResult<ValidPreparationInstrument>) => {
        dispatchPageEvent({
          type: 'UPDATE_STEP_INSTRUMENT_SUGGESTIONS',
          stepIndex: stepIndex,
          recipeStepInstrumentIndex: pageState.recipe.steps[stepIndex].instruments.length - 1,
          results: (res.data || []).map((x: ValidPreparationInstrument) => {
            return new RecipeStepInstrument({
              instrument: x.instrument,
              name: x.instrument.name,
              notes: '',
              preferenceRank: 0,
              optional: false,
              optionIndex: 0,
            });
          }),
        });

        return res.data || [];
      })
      .catch((err: AxiosError) => {
        console.error(`Failed to get preparation instruments: ${err}`);
      });
  };

  const handleValidInstrumentSelection =
    (stepIndex: number, recipeStepInstrumentIndex: number) => (instrument: string) => {
      const rawSelectedInstrument = (pageState.stepHelpers[stepIndex].instrumentSuggestions || []).find(
        (instrumentSuggestion: RecipeStepInstrument) => {
          if (
            pageState.recipe.steps[stepIndex].instruments.find(
              (instrument: RecipeStepInstrumentCreationRequestInput) => {
                return (
                  instrument.instrumentID === instrumentSuggestion.instrument?.id ||
                  instrument.name === instrumentSuggestion.name
                );
              },
            )
          ) {
            return false;
          }
          return instrument === instrumentSuggestion.name;
        },
      );

      if (!rawSelectedInstrument) {
        console.error("couldn't find instrument to add");
        return;
      }

      const selectedInstrument = new RecipeStepInstrument({
        ...rawSelectedInstrument,
        minimumQuantity: 1,
      });

      if (instrument) {
        dispatchPageEvent({
          type: 'SET_VALID_INSTRUMENT_FOR_RECIPE_STEP_INSTRUMENT',
          stepIndex: stepIndex,
          recipeStepInstrumentIndex: recipeStepInstrumentIndex,
          selectedValidInstrument: selectedInstrument,
        });
      }
    };

  const determineInstrumentOptionsForInput = (stepIndex: number, filter: boolean = true) => {
    const base = pageState.stepHelpers[stepIndex].instrumentSuggestions || [];

    return (
      filter
        ? base.filter((x: RecipeStepInstrument) => {
            return !pageState.recipe.steps[stepIndex].instruments.find(
              (y: RecipeStepInstrumentCreationRequestInput) => y.name === x.name,
            );
          })
        : base
    ).map(
      (x: RecipeStepInstrument) =>
        ({
          value: x.instrument?.name || x.name || 'UNKNOWN',
          label: x.instrument?.name || x.name || 'UNKNOWN',
        } as SelectItem),
    );
  };

  const handleInstrumentProductSelection =
    (stepIndex: number, recipeStepInstrumentIndex: number) => (instrument: string) => {
      const base = determinePreparedInstrumentOptions(pageState.recipe, stepIndex) || [];
      const rawSelectedInstrument = base.find((instrumentSuggestion: RecipeStepInstrumentSuggestion) => {
        if (
          pageState.recipe.steps[stepIndex].instruments.find((instrument: RecipeStepInstrumentCreationRequestInput) => {
            return (
              instrument.instrumentID === instrumentSuggestion.product.id ||
              instrument.name === instrumentSuggestion.product.name
            );
          })
        ) {
          return false;
        }
        return instrument === instrumentSuggestion.product.name;
      });

      if (!rawSelectedInstrument) {
        console.error(`couldn't find instrument to add: ${JSON.stringify(base)} ${instrument}`);
        return;
      }

      const selectedInstrument = new RecipeStepInstrument({
        ...rawSelectedInstrument?.product,
        minimumQuantity: 1,
      });

      dispatchPageEvent({
        type: 'SET_PRODUCT_INSTRUMENT_FOR_RECIPE_STEP_INSTRUMENT',
        stepIndex: stepIndex,
        recipeStepInstrumentIndex: recipeStepInstrumentIndex,
        selectedValidInstrument: selectedInstrument,
        productOfRecipeStepIndex: rawSelectedInstrument.stepIndex,
        productOfRecipeStepProductIndex: rawSelectedInstrument.productIndex,
      });
    };

  const determineInstrumentProductOptionsForInput = (stepIndex: number, filter: boolean = true) => {
    const baseOptions = determinePreparedInstrumentOptions(pageState.recipe, stepIndex);

    const filteredOptions = filter
      ? baseOptions.filter((x: RecipeStepInstrumentSuggestion) => {
          return !pageState.recipe.steps[stepIndex].instruments.find(
            (y: RecipeStepInstrumentCreationRequestInput) => y.name !== '' && y.name === x.product.name,
          );
        })
      : baseOptions;

    return filteredOptions.map((x: RecipeStepInstrumentSuggestion) => ({
      value: x.product.name || 'UNKNOWN',
      label: x.product.name || 'UNKNOWN',
    }));
  };

  const handleProductVesselSelection =
    (stepIndex: number, recipeStepIngredientIndex: number) => async (item: string) => {
      const vessels = determineAvailableRecipeStepVessels(pageState.recipe, stepIndex) || [];
      const selectedVessel = vessels.find(
        (vesselSuggestion: RecipeStepVesselSuggestion) => vesselSuggestion.vessel.name === item,
      );

      if (!selectedVessel) {
        console.error("couldn't find vessel to add");
        return;
      }

      dispatchPageEvent({
        type: 'SET_VESSEL_FOR_RECIPE_STEP_VESSEL',
        stepIndex: stepIndex,
        recipeStepIngredientIndex: recipeStepIngredientIndex,
        selectedVessel: selectedVessel.vessel,
        productOfRecipeStepIndex: selectedVessel.stepIndex,
        productOfRecipeStepProductIndex: selectedVessel.productIndex,
      });
    };

  const handleRecipeStepVesselSelection = (stepIndex: number, vesselIndex: number) => (value: AutocompleteItem) => {
    const chosenVessel = (pageState.stepHelpers[stepIndex].vesselSuggestions[vesselIndex] || []).find(
      (vessel: ValidVessel) => {
        return vessel.name === value.value;
      },
    );

    if (!chosenVessel) {
      console.error('Could not find vessel', value);
      return;
    }

    const selectedVessel = new RecipeStepVessel({
      name: chosenVessel.name,
      vessel: chosenVessel,
      minimumQuantity: 1,
    });

    dispatchPageEvent({
      type: 'UPDATE_STEP_VESSEL_INSTRUMENT',
      stepIndex,
      vesselIndex,
      selectedVessel: selectedVessel,
    });
  };

  const handleIngredientQueryChange =
    (stepIndex: number, recipeStepIngredientIndex: number) => async (value: string) => {
      dispatchPageEvent({
        type: 'UPDATE_STEP_INGREDIENT_QUERY',
        newQuery: value,
        stepIndex: stepIndex,
        recipeStepIngredientIndex: recipeStepIngredientIndex,
      });

      // FIXME: if a user selects a choice from the dropdown, it updates the query value first, then
      // this code runs, which then updates the query value again.
      if (value.length > 2) {
        const chosenPreparationID = pageState.stepHelpers[stepIndex].selectedPreparation?.id || '';
        if (!chosenPreparationID) {
          return;
        }

        await apiClient
          .getValidIngredientsForPreparation(chosenPreparationID, value)
          .then((res: QueryFilteredResult<ValidIngredient>) => {
            dispatchPageEvent({
              type: 'UPDATE_STEP_INGREDIENT_SUGGESTIONS',
              stepIndex: stepIndex,
              recipeStepIngredientIndex: recipeStepIngredientIndex,
              results: (
                res.data.filter((validIngredient: ValidIngredient) => {
                  let found = false;

                  (pageState.recipe.steps[stepIndex]?.ingredients || []).forEach((ingredient) => {
                    if ((ingredient.ingredientID || '') === validIngredient.id) {
                      found = true;
                    }
                  });

                  // return true if the ingredient is not already used by another ingredient in the step
                  return !found;
                }) || []
              ).map((x: ValidIngredient) => {
                return new RecipeStepIngredient({
                  name: x.name,
                  ingredient: x,
                  minimumQuantity: 1,
                });
              }),
            });
          })
          .catch((err: AxiosError) => {
            console.error(`Failed to get ingredients: ${err}`);
          });
      } else {
        dispatchPageEvent({
          type: 'UPDATE_STEP_INGREDIENT_SUGGESTIONS',
          stepIndex: stepIndex,
          recipeStepIngredientIndex: recipeStepIngredientIndex,
          results: [],
        });
      }
    };

  const handleIngredientSelection =
    (stepIndex: number, recipeStepIngredientIndex: number) => async (item: AutocompleteItem) => {
      const selectedValidIngredient = (
        pageState.stepHelpers[stepIndex].ingredientSuggestions[recipeStepIngredientIndex] || []
      ).find((ingredientSuggestion: RecipeStepIngredient) => ingredientSuggestion.ingredient?.name === item.value);

      if (!selectedValidIngredient) {
        console.error("couldn't find ingredient to add");
        return;
      }

      if (
        pageState.recipe.steps[stepIndex].ingredients.find(
          (ingredient) => ingredient.ingredientID === selectedValidIngredient.ingredient?.id,
        )
      ) {
        console.error('ingredient already added');

        dispatchPageEvent({
          type: 'UPDATE_STEP_INGREDIENT_SUGGESTIONS',
          stepIndex: stepIndex,
          recipeStepIngredientIndex: recipeStepIngredientIndex,
          results: [],
        });

        return;
      }

      dispatchPageEvent({
        type: 'SET_INGREDIENT_FOR_RECIPE_STEP_INGREDIENT',
        stepIndex: stepIndex,
        recipeStepIngredientIndex: recipeStepIngredientIndex,
        selectedValidIngredient: selectedValidIngredient,
      });

      if ((selectedValidIngredient?.ingredient?.id || '').length > 2) {
        await apiClient
          .searchForValidMeasurementUnitsByIngredientID(selectedValidIngredient!.ingredient!.id)
          .then((res: QueryFilteredResult<ValidMeasurementUnit>) => {
            dispatchPageEvent({
              type: 'UPDATE_STEP_INGREDIENT_MEASUREMENT_UNIT_SUGGESTIONS',
              stepIndex: stepIndex,
              recipeStepIngredientIndex: recipeStepIngredientIndex,
              results: res.data.filter((vmu, index) => index === res.data.findIndex((other) => vmu.id === other.id)),
            });
          })
          .catch((err: AxiosError) => {
            console.error(`Failed to get ingredient measurement units: ${err}`);
          });
      }
    };

  const determineIngredientSuggestions = (stepIndex: number, recipeStepIngredientIndex: number) => {
    const products = pageState.stepHelpers[stepIndex].ingredientSuggestions[recipeStepIngredientIndex] || [];

    return products
      .filter((x?: RecipeStepIngredient) => Boolean(x))
      .map((x: RecipeStepIngredient) => ({
        value: x.ingredient?.name || x.name || 'UNKNOWN',
        label: x.ingredient?.name || x.name || 'UNKNOWN',
      }));
  };

  const handleRecipeStepProductSelection =
    (stepIndex: number, recipeStepIngredientIndex: number) => async (item: string) => {
      const products = determineAvailableRecipeStepProducts(pageState.recipe, stepIndex) || [];
      const selectedValidIngredient = products.find(
        (ingredientSuggestion: RecipeStepProductSuggestion) => ingredientSuggestion.product.name === item,
      );

      if (!selectedValidIngredient) {
        console.error("couldn't find ingredient to add");
        return;
      }

      dispatchPageEvent({
        type: 'SET_INGREDIENT_FOR_RECIPE_STEP_INGREDIENT',
        stepIndex: stepIndex,
        recipeStepIngredientIndex: recipeStepIngredientIndex,
        selectedValidIngredient: selectedValidIngredient.product,
        productOfRecipeStepIndex: selectedValidIngredient.stepIndex,
        productOfRecipeStepProductIndex: selectedValidIngredient.productIndex,
      });

      if (
        pageState.stepHelpers[selectedValidIngredient.stepIndex]?.selectedProductMeasurementUnits[
          selectedValidIngredient.productIndex
        ]
      ) {
        dispatchPageEvent({
          type: 'UPDATE_STEP_INGREDIENT_MEASUREMENT_UNIT_SUGGESTIONS',
          stepIndex: stepIndex,
          recipeStepIngredientIndex: recipeStepIngredientIndex,
          results: [
            pageState.stepHelpers[selectedValidIngredient.stepIndex].selectedProductMeasurementUnits[
              selectedValidIngredient.productIndex
            ]!,
          ],
        });

        dispatchPageEvent({
          type: 'UPDATE_STEP_INGREDIENT_MEASUREMENT_UNIT',
          stepIndex: stepIndex,
          recipeStepIngredientIndex: recipeStepIngredientIndex,
          measurementUnit:
            pageState.stepHelpers[selectedValidIngredient.stepIndex].selectedProductMeasurementUnits[
              selectedValidIngredient.productIndex
            ],
        });
      }
    };

  const determineRecipeStepProductSuggestions = (stepIndex: number) => {
    const products = determineAvailableRecipeStepProducts(pageState.recipe, stepIndex) || [];

    return products
      .filter((x?: RecipeStepProductSuggestion) => (x?.product.name || '') !== '')
      .map((x: RecipeStepProductSuggestion) => ({
        value: x.product.ingredient?.name || x.product.name || 'UNKNOWN',
        label: x.product.ingredient?.name || x.product.name || 'UNKNOWN',
      }));
  };

  const determineRecipeStepVesselSuggestions = (stepIndex: number) => {
    const vessels = determineAvailableRecipeStepVessels(pageState.recipe, stepIndex) || [];

    return vessels
      .filter((x?: RecipeStepVesselSuggestion) => (x?.vessel.name || '') !== '')
      .map((x: RecipeStepVesselSuggestion) => ({
        value: x.vessel.vessel?.name || x.vessel.name || 'UNKNOWN',
        label: x.vessel.vessel?.name || x.vessel.name || 'UNKNOWN',
      }));
  };

  const handleIngredientMeasurementUnitSelection =
    (stepIndex: number, recipeStepIngredientIndex: number) => (value: string) => {
      dispatchPageEvent({
        type: 'UPDATE_STEP_INGREDIENT_MEASUREMENT_UNIT',
        stepIndex: stepIndex,
        recipeStepIngredientIndex: recipeStepIngredientIndex,
        measurementUnit: (
          pageState.stepHelpers[stepIndex].ingredientMeasurementUnitSuggestions[recipeStepIngredientIndex] || []
        ).find((ingredientMeasurementUnitSuggestion: ValidMeasurementUnit) => {
          return ingredientMeasurementUnitSuggestion.pluralName === value;
        }),
      });
    };

  const handleCompletionConditionIngredientStateQueryChange =
    (stepIndex: number, conditionIndex: number) => async (value: string) => {
      dispatchPageEvent({
        type: 'UPDATE_COMPLETION_CONDITION_INGREDIENT_STATE_QUERY',
        stepIndex,
        conditionIndex,
        query: value,
      });

      if (value.length > 2 && !pageState.recipe.steps[stepIndex].completionConditions[conditionIndex].ingredientState) {
        await apiClient
          .searchForValidIngredientStates(value)
          .then((res: ValidIngredientState[]) => {
            dispatchPageEvent({
              type: 'UPDATE_COMPLETION_CONDITION_INGREDIENT_STATE_SUGGESTIONS',
              stepIndex: stepIndex,
              conditionIndex: conditionIndex,
              results: res,
            });
          })
          .catch((err: AxiosError) => {
            console.error(`Failed to get preparations: ${err}`);
          });
      }
    };

  const handleCompletionConditionIngredientStateSelection =
    (stepIndex: number, conditionIndex: number) => (value: AutocompleteItem) => {
      const selectedValidIngredientState = pageState.stepHelpers[
        stepIndex
      ].completionConditionIngredientStateSuggestions[conditionIndex].find(
        (x: ValidIngredientState) => x.name === value.value,
      );

      if (!selectedValidIngredientState) {
        console.error(`unknown ingredient state: ${value.value}`);
        return;
      }

      dispatchPageEvent({
        type: 'UPDATE_COMPLETION_CONDITION_INGREDIENT_STATE',
        stepIndex,
        conditionIndex,
        ingredientState: selectedValidIngredientState,
      });
    };

  const handleRecipeStepProductMeasurementUnitQueryUpdate =
    (stepIndex: number, productIndex: number) => async (value: string) => {
      dispatchPageEvent({
        type: 'UPDATE_STEP_PRODUCT_MEASUREMENT_UNIT_QUERY',
        stepIndex,
        productIndex,
        newQuery: value,
      });

      if (value.length > 2) {
        await apiClient
          .searchForValidMeasurementUnits(value)
          .then((res: ValidMeasurementUnit[]) => {
            dispatchPageEvent({
              type: 'UPDATE_STEP_PRODUCT_MEASUREMENT_UNIT_SUGGESTIONS',
              stepIndex: stepIndex,
              productIndex: productIndex,
              results: res || [],
            });
          })
          .catch((err: AxiosError) => {
            console.error(`Failed to get ingredient measurement units: ${err}`);
          });
      }
    };

  const handleRecipeStepVesselQueryUpdate = (stepIndex: number, vesselIndex: number) => async (value: string) => {
    dispatchPageEvent({
      type: 'UPDATE_STEP_VESSEL_INSTRUMENT_QUERY',
      stepIndex,
      vesselIndex,
      newQuery: value,
    });

    if (value.length > 2) {
      await apiClient
        .searchForValidVessels(value)
        .then((res: ValidVessel[]) => {
          dispatchPageEvent({
            type: 'UPDATE_STEP_VESSEL_SUGGESTIONS',
            stepIndex: stepIndex,
            vesselIndex: vesselIndex,
            results: res || [],
          });
        })
        .catch((err: AxiosError) => {
          console.error(`Failed to get ingredient measurement units: ${err}`);
        });
    }
  };

  const handleRecipeStepProductMeasurementUnitSelection =
    (stepIndex: number, productIndex: number) => (value: AutocompleteItem) => {
      const selectedMeasurementUnit = (
        pageState.stepHelpers[stepIndex].productMeasurementUnitSuggestions[productIndex] || []
      ).find((productMeasurementUnitSuggestion: ValidMeasurementUnit) => {
        return productMeasurementUnitSuggestion.pluralName === value.value;
      });

      if (!selectedMeasurementUnit) {
        console.error('Could not find measurement unit', value);
        return;
      }

      dispatchPageEvent({
        type: 'UPDATE_STEP_PRODUCT_MEASUREMENT_UNIT',
        stepIndex,
        productIndex,
        measurementUnit: selectedMeasurementUnit,
      });
    };

  const recipeStepProductIsUsedInLaterStep = (
    recipe: RecipeCreationRequestInput,
    stepIndex: number,
    productIndex: number,
  ): boolean => {
    return (
      (
        recipe.steps.filter((step: RecipeStepCreationRequestInput) => {
          return (
            (
              step.ingredients.filter((ingredient: RecipeStepIngredientCreationRequestInput) => {
                return (
                  ingredient.productOfRecipeStepIndex === stepIndex &&
                  ingredient.productOfRecipeStepProductIndex === productIndex
                );
              }) || []
            ).length > 0
          );
        }) || []
      ).length > 0
    );
  };

  const manualProductNamingShouldBeDisabled = (stepIndex: number, productIndex: number): boolean => {
    const step = pageState.recipe.steps[stepIndex];
    const product = step.products[productIndex];

    const productIsIngredientButStepHasNoIngredients = product.type === 'ingredient' && step.ingredients.length === 0;
    const productIsUsedElsewhere = recipeStepProductIsUsedInLaterStep(pageState.recipe, stepIndex, productIndex);
    const stepIsLocked = pageState.stepHelpers[stepIndex].locked;

    return productIsIngredientButStepHasNoIngredients || productIsUsedElsewhere || stepIsLocked;
  };

  const buildStep = (step: RecipeStepCreationRequestInput, stepIndex: number): JSX.Element => {
    return (
      <Card key={stepIndex} shadow="sm" radius="md" withBorder sx={{ width: '100%', marginBottom: '1rem' }}>
        {/* this is the top of the recipe step view, with the step index indicator and the delete step button */}
        <Card.Section px="xs" sx={{ cursor: 'pointer' }}>
          <Grid justify="space-between" align="center">
            <Grid.Col span="auto">
              <Text weight="bold">{`Step #${stepIndex + 1}`}</Text>
            </Grid.Col>
            <Grid.Col span="content">
              <ActionIcon
                data-qa={`toggle-step-${stepIndex}`}
                variant="outline"
                size="sm"
                style={{ float: 'right' }}
                aria-label="show step"
                disabled={
                  ((pageState.recipe.steps || []).length === 1 ||
                    pageState.stepHelpers.filter((x) => x.show).length === 1) &&
                  pageState.stepHelpers[stepIndex].show
                }
                onClick={() => dispatchPageEvent({ type: 'TOGGLE_SHOW_STEP', stepIndex: stepIndex })}
              >
                {pageState.stepHelpers[stepIndex].show ? (
                  <IconChevronUp size={16} color={(pageState.recipe.steps || []).length === 1 ? 'gray' : 'black'} />
                ) : (
                  <IconChevronDown size={16} color={(pageState.recipe.steps || []).length === 1 ? 'gray' : 'black'} />
                )}
              </ActionIcon>
            </Grid.Col>
            <Grid.Col span="content">
              <ActionIcon
                variant="outline"
                size="sm"
                style={{ float: 'right' }}
                aria-label="remove step"
                disabled={(pageState.recipe.steps || []).length === 1}
                onClick={() => dispatchPageEvent({ type: 'REMOVE_STEP', stepIndex: stepIndex })}
              >
                <IconTrash size={16} color={(pageState.recipe.steps || []).length === 1 ? 'gray' : 'tomato'} />
              </ActionIcon>
            </Grid.Col>
          </Grid>
        </Card.Section>

        <Collapse in={pageState.stepHelpers[stepIndex].show}>
          {/* this is the first input section */}
          <Card.Section px="xs" pb="xs">
            <Grid>
              <Grid.Col md="auto">
                <Stack>
                  <Autocomplete
                    label="Preparation"
                    required
                    tabIndex={0}
                    disabled={pageState.stepHelpers[stepIndex].locked}
                    value={pageState.stepHelpers[stepIndex].preparationQuery}
                    onChange={handlePreparationQueryChange(stepIndex)}
                    data={pageState.stepHelpers[stepIndex].preparationSuggestions
                      .filter((x: ValidPreparation) => {
                        return x.name !== pageState.stepHelpers[stepIndex].selectedPreparation?.name;
                      })
                      .map((x: ValidPreparation) => ({
                        value: x.name,
                        label: x.name,
                      }))}
                    onItemSubmit={handlePreparationSelection(stepIndex)}
                    rightSection={
                      pageState.stepHelpers[stepIndex].selectedPreparation && (
                        <IconCircleX
                          size={18}
                          color={pageState.stepHelpers[stepIndex].selectedPreparation ? 'tomato' : 'gray'}
                          onClick={() => {
                            if (!pageState.stepHelpers[stepIndex].selectedPreparation) {
                              return;
                            }

                            dispatchPageEvent({
                              type: 'UNSET_STEP_PREPARATION',
                              stepIndex: stepIndex,
                            });
                          }}
                        />
                      )
                    }
                  />

                  <Grid>
                    <Grid.Col span="auto">
                      <Textarea
                        label="Explicit Instructions"
                        value={step.explicitInstructions}
                        minRows={2}
                        disabled={pageState.stepHelpers[stepIndex].locked}
                        onChange={(event: React.ChangeEvent<HTMLTextAreaElement>) => {
                          dispatchPageEvent({
                            type: 'UPDATE_STEP_EXPLICIT_INSTRUCTIONS',
                            stepIndex: stepIndex,
                            newExplicitInstructions: event.target.value,
                          });
                        }}
                      />
                    </Grid.Col>

                    <Grid.Col span="auto">
                      <Textarea
                        label="Notes"
                        value={step.notes}
                        minRows={2}
                        disabled={pageState.stepHelpers[stepIndex].locked}
                        onChange={(event: React.ChangeEvent<HTMLTextAreaElement>) => {
                          dispatchPageEvent({
                            type: 'UPDATE_STEP_NOTES',
                            stepIndex: stepIndex,
                            newNotes: event.target.value,
                          });
                        }}
                      />
                    </Grid.Col>
                  </Grid>

                  <Grid>
                    <Grid.Col span="auto">
                      <NumberInput
                        data-qa={`recipe-step-${stepIndex}-min-estimated-time-in-seconds`}
                        label="Min. Time"
                        placeholder="seconds"
                        disabled={pageState.stepHelpers[stepIndex].locked}
                        onChange={(value: number) => {
                          if (value <= 0) {
                            return;
                          }

                          dispatchPageEvent({
                            type: 'UPDATE_STEP_MINIMUM_TIME_ESTIMATE',
                            stepIndex,
                            newMinTimeEstimate: value,
                          });
                        }}
                        value={step.minimumEstimatedTimeInSeconds}
                        maxLength={0}
                      />
                    </Grid.Col>

                    <Grid.Col span="auto">
                      <NumberInput
                        data-qa={`recipe-step-${stepIndex}-max-estimated-time-in-seconds`}
                        label="Max. Time"
                        placeholder="seconds"
                        disabled={pageState.stepHelpers[stepIndex].locked}
                        onChange={(value: number) => {
                          if (value <= 0) {
                            return;
                          }

                          dispatchPageEvent({
                            type: 'UPDATE_STEP_MAXIMUM_TIME_ESTIMATE',
                            stepIndex,
                            newMaxTimeEstimate: value,
                          });
                        }}
                        value={step.maximumEstimatedTimeInSeconds}
                        maxLength={0}
                      />
                    </Grid.Col>

                    <Grid.Col span="auto">
                      <NumberInput
                        data-qa={`recipe-step-${stepIndex}-min-temperature-in-celsius`}
                        label="Min. Temp"
                        placeholder="celsius"
                        disabled={pageState.stepHelpers[stepIndex].locked}
                        onChange={(value: number) => {
                          if (value <= 0) {
                            return;
                          }

                          dispatchPageEvent({
                            type: 'UPDATE_STEP_MINIMUM_TEMPERATURE',
                            stepIndex,
                            newMinTempInCelsius: value,
                          });
                        }}
                        value={step.minimumTemperatureInCelsius}
                        maxLength={0}
                      />
                    </Grid.Col>

                    <Grid.Col span="auto">
                      <NumberInput
                        data-qa={`recipe-step-${stepIndex}-max-temperature-in-celsius`}
                        label="Max. Temp"
                        placeholder="celsius"
                        disabled={pageState.stepHelpers[stepIndex].locked}
                        onChange={(value: number) => {
                          if (value <= 0) {
                            return;
                          }

                          dispatchPageEvent({
                            type: 'UPDATE_STEP_MAXIMUM_TEMPERATURE',
                            stepIndex,
                            newMaxTempInCelsius: value,
                          });
                        }}
                        value={step.maximumTemperatureInCelsius}
                        maxLength={0}
                      />
                    </Grid.Col>
                  </Grid>

                  <Divider label="with tools" labelPosition="center" />

                  {(step.instruments || []).map(
                    (instrument: RecipeStepInstrumentCreationRequestInput, recipeStepInstrumentIndex: number) => (
                      <Box key={recipeStepInstrumentIndex}>
                        <Grid>
                          <Grid.Col span="content">
                            {stepIndex !== 0 && (
                              <Switch
                                data-qa={`toggle-recipe-step-${stepIndex}-instrument-${recipeStepInstrumentIndex}-product-switch`}
                                mt="sm"
                                size="md"
                                onLabel="product"
                                offLabel="instrument"
                                disabled={
                                  pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                                  pageState.stepHelpers[stepIndex].selectedInstruments.length === 0 ||
                                  pageState.stepHelpers[stepIndex].locked
                                }
                                value={
                                  pageState.stepHelpers[stepIndex].instrumentIsProduct[recipeStepInstrumentIndex]
                                    ? 'product'
                                    : 'instrument'
                                }
                                onChange={() => {
                                  dispatchPageEvent({
                                    type: 'TOGGLE_INSTRUMENT_PRODUCT_STATE',
                                    stepIndex: stepIndex,
                                    recipeStepInstrumentIndex: recipeStepInstrumentIndex,
                                  });
                                }}
                              />
                            )}
                          </Grid.Col>

                          <Grid.Col span="auto">
                            {((stepIndex === 0 ||
                              !pageState.stepHelpers[stepIndex].instrumentIsProduct[recipeStepInstrumentIndex]) && (
                              <Select
                                label="Instrument"
                                required
                                disabled={
                                  pageState.stepHelpers[stepIndex].locked ||
                                  !pageState.stepHelpers[stepIndex].selectedPreparation ||
                                  determineInstrumentOptionsForInput(stepIndex, false).length == 0
                                }
                                onChange={handleValidInstrumentSelection(stepIndex, recipeStepInstrumentIndex)}
                                value={instrument.name}
                                data={determineInstrumentOptionsForInput(stepIndex, false)}
                              />
                            )) || (
                              <Select
                                label="Step Instrument"
                                required
                                disabled={
                                  pageState.stepHelpers[stepIndex].locked ||
                                  !pageState.stepHelpers[stepIndex].selectedPreparation ||
                                  determineInstrumentProductOptionsForInput(stepIndex).length == 0
                                }
                                onChange={handleInstrumentProductSelection(stepIndex, recipeStepInstrumentIndex)}
                                value={instrument.name}
                                data={determineInstrumentProductOptionsForInput(stepIndex, false)}
                              />
                            )}
                          </Grid.Col>

                          <Grid.Col span="auto">
                            <NumberInput
                              data-qa={`recipe-step-${stepIndex}-instrument-${recipeStepInstrumentIndex}-min-quantity-input`}
                              label={
                                pageState.stepHelpers[stepIndex].instrumentIsRanged[recipeStepInstrumentIndex]
                                  ? 'Min. Quantity'
                                  : 'Quantity'
                              }
                              required
                              disabled={pageState.stepHelpers[stepIndex].locked}
                              onChange={(value: number) => {
                                if (value <= 0) {
                                  return;
                                }

                                dispatchPageEvent({
                                  type: 'UPDATE_STEP_INSTRUMENT_MINIMUM_QUANTITY',
                                  stepIndex,
                                  recipeStepInstrumentIndex,
                                  newAmount: value,
                                });
                              }}
                              value={step.instruments[recipeStepInstrumentIndex].minimumQuantity}
                              maxLength={0}
                            />
                          </Grid.Col>

                          {pageState.stepHelpers[stepIndex].instrumentIsRanged[recipeStepInstrumentIndex] && (
                            <Grid.Col span="auto">
                              <NumberInput
                                data-qa={`recipe-step-${stepIndex}-instrument-${recipeStepInstrumentIndex}-max-quantity-input`}
                                label="Max Quantity"
                                maxLength={0}
                                disabled={pageState.stepHelpers[stepIndex].locked}
                                onChange={(value: number) => {
                                  if (value <= 0) {
                                    return;
                                  }

                                  dispatchPageEvent({
                                    type: 'UPDATE_STEP_INSTRUMENT_MAXIMUM_QUANTITY',
                                    stepIndex,
                                    recipeStepInstrumentIndex,
                                    newAmount: value,
                                  });
                                }}
                                value={step.instruments[recipeStepInstrumentIndex].maximumQuantity}
                              />
                            </Grid.Col>
                          )}

                          <Grid.Col span="content">
                            <Switch
                              data-qa={`toggle-recipe-step-${stepIndex}-instrument-${recipeStepInstrumentIndex}-ranged-status`}
                              mt="sm"
                              size="md"
                              onLabel="ranged"
                              offLabel="simple"
                              disabled={pageState.stepHelpers[stepIndex].locked}
                              value={
                                pageState.stepHelpers[stepIndex].instrumentIsRanged[recipeStepInstrumentIndex]
                                  ? 'ranged'
                                  : 'simple'
                              }
                              onChange={() => {
                                dispatchPageEvent({
                                  type: 'TOGGLE_INSTRUMENT_RANGE',
                                  stepIndex,
                                  recipeStepInstrumentIndex,
                                });
                              }}
                            />
                          </Grid.Col>

                          <Grid.Col span="content" mt="sm">
                            <ActionIcon
                              data-qa={`remove-recipe-step-${stepIndex}-instrument-${recipeStepInstrumentIndex}`}
                              mt="sm"
                              variant="outline"
                              size="sm"
                              aria-label="remove recipe step instrument"
                              disabled={
                                determineInstrumentOptionsForInput(stepIndex, true).length == 0 ||
                                pageState.stepHelpers[stepIndex].locked
                              }
                              onClick={() => {
                                dispatchPageEvent({
                                  type: 'REMOVE_INSTRUMENT_FROM_STEP',
                                  stepIndex,
                                  recipeStepInstrumentIndex,
                                });
                              }}
                            >
                              <IconTrash size="md" color="tomato" />
                            </ActionIcon>
                          </Grid.Col>
                        </Grid>
                      </Box>
                    ),
                  )}

                  {determineInstrumentOptionsForInput(stepIndex, true).length > 0 && (
                    <Center>
                      <Button
                        mt="sm"
                        disabled={
                          determineInstrumentOptionsForInput(stepIndex, true).length == 0 ||
                          pageState.stepHelpers[stepIndex].locked
                        }
                        style={{
                          cursor: addingStepCompletionConditionsShouldBeDisabled(step) ? 'not-allowed' : 'pointer',
                        }}
                        onClick={() => {
                          dispatchPageEvent({
                            type: 'ADD_INSTRUMENT_TO_STEP',
                            stepIndex: stepIndex,
                          });
                        }}
                      >
                        Add Instrument
                      </Button>
                    </Center>
                  )}
                </Stack>

                <Space h="sm" />

                <Divider label="consuming" labelPosition="center" mb="md" />

                {(step.ingredients || []).map(
                  (ingredient: RecipeStepIngredientCreationRequestInput, recipeStepIngredientIndex: number) => (
                    <Box key={recipeStepIngredientIndex}>
                      <Grid>
                        <Grid.Col span="content">
                          {stepIndex !== 0 && (
                            <Switch
                              data-qa={`toggle-recipe-step-${stepIndex}-ingredient-${recipeStepIngredientIndex}-product-switch`}
                              mt="sm"
                              size="md"
                              onLabel="product"
                              offLabel="ingredient"
                              disabled={
                                pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                                pageState.stepHelpers[stepIndex].selectedInstruments.length === 0 ||
                                pageState.stepHelpers[stepIndex].locked
                              }
                              value={
                                pageState.stepHelpers[stepIndex].ingredientIsProduct[recipeStepIngredientIndex]
                                  ? 'product'
                                  : 'ingredient'
                              }
                              onChange={() => {
                                dispatchPageEvent({
                                  type: 'TOGGLE_INGREDIENT_PRODUCT_STATE',
                                  stepIndex: stepIndex,
                                  recipeStepIngredientIndex: recipeStepIngredientIndex,
                                });
                              }}
                            />
                          )}
                        </Grid.Col>

                        <Grid.Col span="content">
                          {((stepIndex == 0 ||
                            !pageState.stepHelpers[stepIndex].ingredientIsProduct[recipeStepIngredientIndex]) && (
                            <Autocomplete
                              data-qa={`recipe-step-${stepIndex}-ingredient-input-${recipeStepIngredientIndex}`}
                              label="Ingredient"
                              limit={20}
                              required
                              value={pageState.stepHelpers[stepIndex].ingredientQueries[recipeStepIngredientIndex]}
                              disabled={
                                pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                                pageState.stepHelpers[stepIndex].selectedInstruments.length === 0 ||
                                pageState.stepHelpers[stepIndex].locked
                              }
                              onChange={handleIngredientQueryChange(stepIndex, recipeStepIngredientIndex)}
                              onItemSubmit={handleIngredientSelection(stepIndex, recipeStepIngredientIndex)}
                              data={determineIngredientSuggestions(stepIndex, recipeStepIngredientIndex)}
                              rightSection={
                                pageState.stepHelpers[stepIndex].selectedPreparation && (
                                  <IconCircleX
                                    size={18}
                                    color={
                                      pageState.stepHelpers[stepIndex].selectedIngredients[recipeStepIngredientIndex]
                                        ? 'tomato'
                                        : 'gray'
                                    }
                                    onClick={() => {
                                      if (
                                        !pageState.stepHelpers[stepIndex].selectedIngredients[recipeStepIngredientIndex]
                                      ) {
                                        return;
                                      }

                                      dispatchPageEvent({
                                        type: 'UNSET_RECIPE_STEP_INGREDIENT',
                                        stepIndex: stepIndex,
                                        recipeStepIngredientIndex: recipeStepIngredientIndex,
                                      });
                                    }}
                                  />
                                )
                              }
                            />
                          )) || (
                            <Select
                              data-qa={`recipe-step-${stepIndex}-ingredient-input-${recipeStepIngredientIndex}`}
                              label="Ingredient"
                              limit={20}
                              required
                              value={pageState.stepHelpers[stepIndex].ingredientQueries[recipeStepIngredientIndex]}
                              disabled={
                                pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                                pageState.stepHelpers[stepIndex].locked
                              }
                              onChange={handleRecipeStepProductSelection(stepIndex, recipeStepIngredientIndex)}
                              data={determineRecipeStepProductSuggestions(stepIndex)}
                            />
                          )}
                        </Grid.Col>

                        <Grid.Col span="auto">
                          <NumberInput
                            data-qa={`recipe-step-${stepIndex}-ingredient-${recipeStepIngredientIndex}-min-quantity-input`}
                            label={
                              pageState.stepHelpers[stepIndex].ingredientIsRanged[recipeStepIngredientIndex]
                                ? 'Min. Quantity'
                                : 'Quantity'
                            }
                            required={!pageState.stepHelpers[stepIndex].ingredientIsProduct[recipeStepIngredientIndex]}
                            disabled={
                              pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                              !pageState.stepHelpers[stepIndex].selectedIngredients[recipeStepIngredientIndex] ||
                              pageState.stepHelpers[stepIndex].locked
                            }
                            onChange={(value: number) => {
                              if (value <= 0) {
                                return;
                              }

                              dispatchPageEvent({
                                type: 'UPDATE_STEP_INGREDIENT_MINIMUM_QUANTITY',
                                stepIndex: stepIndex,
                                recipeStepIngredientIndex: recipeStepIngredientIndex,
                                newAmount: value,
                              });
                            }}
                            value={ingredient.minimumQuantity}
                          />
                        </Grid.Col>

                        {pageState.stepHelpers[stepIndex].ingredientIsRanged[recipeStepIngredientIndex] && (
                          <Grid.Col span="auto">
                            <NumberInput
                              data-qa={`recipe-step-${stepIndex}-ingredient-${recipeStepIngredientIndex}-max-quantity-input`}
                              label="Max Quantity"
                              disabled={
                                pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                                !pageState.stepHelpers[stepIndex].selectedIngredients[recipeStepIngredientIndex] ||
                                pageState.stepHelpers[stepIndex].locked
                              }
                              onChange={(value: number) => {
                                if (value <= 0) {
                                  return;
                                }

                                dispatchPageEvent({
                                  type: 'UPDATE_STEP_INGREDIENT_MAXIMUM_QUANTITY',
                                  stepIndex: stepIndex,
                                  recipeStepIngredientIndex: recipeStepIngredientIndex,
                                  newAmount: value,
                                });
                              }}
                              value={ingredient.maximumQuantity}
                            />
                          </Grid.Col>
                        )}

                        <Grid.Col span="content">
                          <Switch
                            data-qa={`toggle-recipe-step-${stepIndex}-ingredient-${recipeStepIngredientIndex}-range`}
                            mt="sm"
                            size="md"
                            onLabel="ranged"
                            offLabel="single"
                            disabled={
                              pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                              pageState.stepHelpers[stepIndex].locked
                            }
                            value={
                              pageState.stepHelpers[stepIndex].ingredientIsRanged[recipeStepIngredientIndex]
                                ? 'ranged'
                                : 'single'
                            }
                            onChange={() => {
                              dispatchPageEvent({
                                type: 'TOGGLE_INGREDIENT_RANGE',
                                stepIndex: stepIndex,
                                recipeStepIngredientIndex: recipeStepIngredientIndex,
                              });
                            }}
                          />
                        </Grid.Col>

                        <Grid.Col span="auto">
                          <Select
                            data-qa={`recipe-step-${stepIndex}-ingredient-${recipeStepIngredientIndex}-measurement-unit-input`}
                            label="Measurement"
                            required={!pageState.stepHelpers[stepIndex].ingredientIsProduct[recipeStepIngredientIndex]}
                            disabled={
                              pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                              !pageState.stepHelpers[stepIndex].selectedIngredients[recipeStepIngredientIndex] ||
                              pageState.stepHelpers[stepIndex].locked
                            }
                            value={
                              pageState.stepHelpers[stepIndex].selectedMeasurementUnits[recipeStepIngredientIndex]
                                ?.pluralName || ''
                            }
                            onChange={handleIngredientMeasurementUnitSelection(stepIndex, recipeStepIngredientIndex)}
                            data={(
                              pageState.stepHelpers[stepIndex].ingredientMeasurementUnitSuggestions[
                                recipeStepIngredientIndex
                              ] || []
                            ).map((y: ValidMeasurementUnit) => ({
                              value: y.pluralName,
                              label: y.pluralName,
                            }))}
                          />
                        </Grid.Col>

                        <Grid.Col span="content" mt="sm">
                          <ActionIcon
                            data-qa={`remove-recipe-step-${stepIndex}-ingredient-${recipeStepIngredientIndex}`}
                            mt="sm"
                            variant="outline"
                            size="sm"
                            disabled={
                              pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                              pageState.stepHelpers[stepIndex].locked
                            }
                            aria-label="remove recipe step ingredient"
                            onClick={() => {
                              dispatchPageEvent({
                                type: 'REMOVE_INGREDIENT_FROM_STEP',
                                stepIndex: stepIndex,
                                recipeStepIngredientIndex: recipeStepIngredientIndex,
                              });
                            }}
                          >
                            <IconTrash
                              size="md"
                              color={pageState.stepHelpers[stepIndex].selectedPreparation === null ? 'grey' : 'tomato'}
                            />
                          </ActionIcon>
                        </Grid.Col>
                      </Grid>
                    </Box>
                  ),
                )}

                <Grid>
                  <Grid.Col span="auto">
                    <Center>
                      <Button
                        mt="sm"
                        disabled={pageState.stepHelpers[stepIndex].locked}
                        style={{
                          cursor: addingStepCompletionConditionsShouldBeDisabled(step) ? 'not-allowed' : 'pointer',
                        }}
                        onClick={() => {
                          dispatchPageEvent({
                            type: 'ADD_INGREDIENT_TO_STEP',
                            stepIndex: stepIndex,
                          });
                        }}
                      >
                        Add Ingredient
                      </Button>
                    </Center>
                  </Grid.Col>
                </Grid>

                <Divider label="in vessels" labelPosition="center" mb="md" mt="md" />

                {(step.vessels || []).map(
                  (vessel: RecipeStepVesselCreationRequestInput, recipeStepVesselIndex: number) => (
                    <Box key={recipeStepVesselIndex}>
                      <Grid>
                        <Grid.Col span="content">
                          {stepIndex !== 0 && (
                            <Switch
                              data-qa={`toggle-recipe-step-${stepIndex}-vessel-${recipeStepVesselIndex}-product-switch`}
                              mt="sm"
                              size="md"
                              onLabel="product"
                              offLabel="vessel"
                              disabled={
                                pageState.stepHelpers[stepIndex].selectedPreparation === null ||
                                pageState.stepHelpers[stepIndex].selectedVessels.length === 0 ||
                                pageState.stepHelpers[stepIndex].locked
                              }
                              value={
                                pageState.stepHelpers[stepIndex].vesselIsProduct[recipeStepVesselIndex]
                                  ? 'product'
                                  : 'vessel'
                              }
                              onChange={() => {
                                dispatchPageEvent({
                                  type: 'TOGGLE_VESSEL_PRODUCT_STATE',
                                  stepIndex: stepIndex,
                                  recipeStepVesselIndex: recipeStepVesselIndex,
                                });
                              }}
                            />
                          )}
                        </Grid.Col>

                        <Grid.Col span={2}>
                          <TextInput
                            data-qa={`recipe-step-${stepIndex}-vessel-input-${recipeStepVesselIndex}`}
                            label="Preposition"
                            required
                            value={vessel.vesselPreposition}
                            disabled={
                              pageState.stepHelpers[stepIndex].locked ||
                              !pageState.stepHelpers[stepIndex].selectedPreparation
                            }
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                              dispatchPageEvent({
                                type: 'SET_RECIPE_STEP_VESSEL_PREDICATE',
                                stepIndex: stepIndex,
                                recipeStepVesselIndex: recipeStepVesselIndex,
                                vesselPreposition: event.target.value || '',
                              });
                            }}
                            placeholder={"'in', 'on', etc."}
                          />
                        </Grid.Col>

                        <Grid.Col span="auto">
                          {((stepIndex === 0 ||
                            !pageState.stepHelpers[stepIndex].vesselIsProduct[recipeStepVesselIndex]) && (
                            <Autocomplete
                              label="Vessel"
                              value={pageState.stepHelpers[stepIndex].vesselQueries[recipeStepVesselIndex]}
                              data={pageState.stepHelpers[stepIndex].vesselSuggestions[recipeStepVesselIndex].map(
                                (x) => {
                                  return {
                                    value: x.name,
                                    label: x.name,
                                  };
                                },
                              )}
                              onItemSubmit={handleRecipeStepVesselSelection(stepIndex, recipeStepVesselIndex)}
                              onChange={handleRecipeStepVesselQueryUpdate(stepIndex, recipeStepVesselIndex)}
                            />
                          )) || (
                            <Select
                              label="Step Vessel"
                              required
                              disabled={
                                pageState.stepHelpers[stepIndex].locked ||
                                !pageState.stepHelpers[stepIndex].selectedPreparation ||
                                determineRecipeStepVesselSuggestions(stepIndex).length == 0
                              }
                              onChange={handleProductVesselSelection(stepIndex, recipeStepVesselIndex)}
                              value={vessel.name}
                              data={determineRecipeStepVesselSuggestions(stepIndex)}
                            />
                          )}
                        </Grid.Col>

                        <Grid.Col span="auto">
                          <NumberInput
                            data-qa={`recipe-step-${stepIndex}-vessel-${recipeStepVesselIndex}-min-quantity-input`}
                            label={
                              pageState.stepHelpers[stepIndex].vesselIsRanged[recipeStepVesselIndex]
                                ? 'Min. Quantity'
                                : 'Quantity'
                            }
                            required
                            disabled={pageState.stepHelpers[stepIndex].locked}
                            onChange={(value: number) => {
                              if (value <= 0) {
                                return;
                              }

                              dispatchPageEvent({
                                type: 'UPDATE_STEP_VESSEL_MINIMUM_QUANTITY',
                                stepIndex,
                                recipeStepVesselIndex,
                                newAmount: value,
                              });
                            }}
                            value={step.vessels[recipeStepVesselIndex].minimumQuantity}
                            maxLength={0}
                          />
                        </Grid.Col>

                        <Grid.Col span="content">
                          <Switch
                            data-qa={`toggle-recipe-step-${stepIndex}-vessel-${recipeStepVesselIndex}-ranged-status`}
                            mt="sm"
                            size="md"
                            onLabel="ranged"
                            offLabel="simple"
                            disabled={pageState.stepHelpers[stepIndex].locked}
                            value={
                              pageState.stepHelpers[stepIndex].vesselIsRanged[recipeStepVesselIndex]
                                ? 'ranged'
                                : 'simple'
                            }
                            onChange={() => {
                              dispatchPageEvent({
                                type: 'TOGGLE_VESSEL_RANGE',
                                stepIndex,
                                recipeStepVesselIndex,
                              });
                            }}
                          />
                        </Grid.Col>

                        {pageState.stepHelpers[stepIndex].vesselIsRanged[recipeStepVesselIndex] && (
                          <Grid.Col span="auto">
                            <NumberInput
                              data-qa={`recipe-step-${stepIndex}-vessel-${recipeStepVesselIndex}-max-quantity-input`}
                              label="Max Quantity"
                              maxLength={0}
                              disabled={pageState.stepHelpers[stepIndex].locked}
                              onChange={(value: number) => {
                                if (value <= 0) {
                                  return;
                                }

                                dispatchPageEvent({
                                  type: 'UPDATE_STEP_VESSEL_MAXIMUM_QUANTITY',
                                  stepIndex,
                                  recipeStepVesselIndex,
                                  newAmount: value,
                                });
                              }}
                              value={step.vessels[recipeStepVesselIndex].maximumQuantity}
                            />
                          </Grid.Col>
                        )}

                        <Grid.Col span="content" mt="sm">
                          <ActionIcon
                            data-qa={`remove-recipe-step-${stepIndex}-vessel-${recipeStepVesselIndex}`}
                            mt="sm"
                            variant="outline"
                            size="sm"
                            aria-label="remove recipe step vessel"
                            disabled={pageState.stepHelpers[stepIndex].locked}
                            onClick={() => {
                              dispatchPageEvent({
                                type: 'REMOVE_VESSEL_FROM_STEP',
                                stepIndex,
                                recipeStepVesselIndex,
                              });
                            }}
                          >
                            <IconTrash size="md" color="tomato" />
                          </ActionIcon>
                        </Grid.Col>
                      </Grid>
                    </Box>
                  ),
                )}

                <Center>
                  <Button
                    mt="sm"
                    disabled={
                      addingStepCompletionConditionsShouldBeDisabled(step) || pageState.stepHelpers[stepIndex].locked
                    }
                    style={{
                      cursor: addingStepCompletionConditionsShouldBeDisabled(step) ? 'not-allowed' : 'pointer',
                    }}
                    onClick={() => {
                      dispatchPageEvent({
                        type: 'ADD_VESSEL_TO_STEP',
                        stepIndex,
                      });
                    }}
                  >
                    Add Vessel
                  </Button>
                </Center>

                <Space h="sm" />

                <Divider label="until" labelPosition="center" my="md" />

                {(step.completionConditions || []).map(
                  (completionCondition: RecipeStepCompletionConditionCreationRequestInput, conditionIndex: number) => {
                    return (
                      <Grid key={conditionIndex}>
                        <Grid.Col span="auto">
                          <Autocomplete
                            data-qa={`recipe-step-${stepIndex}-completion-condition-${conditionIndex}-ingredient-state-input`}
                            label="Ingredient State"
                            required
                            disabled={step.ingredients.length === 0 || pageState.stepHelpers[stepIndex].locked}
                            value={
                              pageState.stepHelpers[stepIndex].completionConditionIngredientStateQueries[conditionIndex]
                            }
                            data={pageState.stepHelpers[stepIndex].completionConditionIngredientStateSuggestions[
                              conditionIndex
                            ].map((x: ValidIngredientState) => {
                              return {
                                value: x.name,
                                label: x.name,
                              };
                            })}
                            onChange={handleCompletionConditionIngredientStateQueryChange(stepIndex, conditionIndex)}
                            onItemSubmit={handleCompletionConditionIngredientStateSelection(stepIndex, conditionIndex)}
                          />
                        </Grid.Col>

                        <Grid.Col span="auto">
                          <Select
                            data-qa={`recipe-step-${stepIndex}-completion-condition-${conditionIndex}-ingredient-selection-input`}
                            disabled={
                              step.ingredients.length === 0 ||
                              !completionCondition.ingredientState ||
                              pageState.stepHelpers[stepIndex].locked
                            }
                            label="Add Ingredient"
                            required
                            data={step.ingredients
                              .filter((x: RecipeStepIngredientCreationRequestInput) => x.ingredientID)
                              .map((x: RecipeStepIngredientCreationRequestInput) => {
                                console.log(`condition ingredient filter: ${JSON.stringify(x)}`);
                                return {
                                  value: x.name,
                                  label: x.name,
                                };
                              })}
                          />
                        </Grid.Col>

                        <Grid.Col span="content" mt="xl">
                          <ActionIcon
                            data-qa={`remove-recipe-step-${stepIndex}-completion-condition-${conditionIndex}`}
                            mt={5}
                            style={{ float: 'right' }}
                            variant="outline"
                            size="md"
                            disabled={pageState.stepHelpers[stepIndex].locked}
                            aria-label="remove condition"
                            onClick={() => {
                              dispatchPageEvent({
                                type: 'REMOVE_RECIPE_STEP_COMPLETION_CONDITION',
                                stepIndex,
                                conditionIndex,
                              });
                            }}
                          >
                            <IconTrash size="md" color="tomato" />
                          </ActionIcon>
                        </Grid.Col>
                      </Grid>
                    );
                  },
                )}

                <Grid>
                  <Grid.Col span="auto">
                    <Center>
                      <Button
                        mt="sm"
                        disabled={
                          addingStepCompletionConditionsShouldBeDisabled(step) ||
                          pageState.stepHelpers[stepIndex].locked
                        }
                        style={{
                          cursor: addingStepCompletionConditionsShouldBeDisabled(step) ? 'not-allowed' : 'pointer',
                        }}
                        onClick={() => {
                          dispatchPageEvent({
                            type: 'ADD_COMPLETION_CONDITION_TO_STEP',
                            stepIndex,
                          });
                        }}
                      >
                        Add Completion Condition
                      </Button>
                    </Center>
                  </Grid.Col>
                </Grid>

                <Divider label="producing" labelPosition="center" my="md" />

                {(step.products || []).map((product: RecipeStepProductCreationRequestInput, productIndex: number) => {
                  return (
                    <Grid key={productIndex}>
                      <Grid.Col md="auto">
                        <Select
                          label="Type"
                          value={product.type}
                          data={validRecipeStepProductTypes}
                          disabled={pageState.stepHelpers[stepIndex].locked}
                          onChange={(value: string) => {
                            dispatchPageEvent({
                              type: 'UPDATE_STEP_PRODUCT_TYPE',
                              stepIndex: stepIndex,
                              productIndex: productIndex,
                              newType: validRecipeStepProductTypes.includes(value)
                                ? (value as ValidRecipeStepProductType)
                                : 'ingredient',
                            });
                          }}
                        />
                      </Grid.Col>

                      <Grid.Col md="auto">
                        <NumberInput
                          label={
                            pageState.stepHelpers[stepIndex].productIsRanged[productIndex]
                              ? 'Min. Quantity'
                              : 'Quantity'
                          }
                          disabled={
                            (product.type === 'ingredient' && step.ingredients.length === 0) ||
                            pageState.stepHelpers[stepIndex].locked
                          }
                          onChange={(value: number) => {
                            if (value <= 0) {
                              return;
                            }

                            dispatchPageEvent({
                              type: 'UPDATE_STEP_PRODUCT_MINIMUM_QUANTITY',
                              stepIndex: stepIndex,
                              productIndex: productIndex,
                              newAmount: value,
                            });
                          }}
                          value={product.minimumQuantity}
                        />
                      </Grid.Col>

                      {pageState.stepHelpers[stepIndex].productIsRanged[productIndex] && (
                        <Grid.Col md="auto">
                          <NumberInput
                            label="Max Quantity"
                            disabled={
                              (product.type === 'ingredient' && step.ingredients.length === 0) ||
                              recipeStepProductIsUsedInLaterStep(pageState.recipe, stepIndex, productIndex) ||
                              pageState.stepHelpers[stepIndex].locked
                            }
                            onChange={(value: number) => {
                              if (value <= 0) {
                                return;
                              }

                              dispatchPageEvent({
                                type: 'UPDATE_STEP_PRODUCT_MAXIMUM_QUANTITY',
                                stepIndex: stepIndex,
                                productIndex: productIndex,
                                newAmount: value,
                              });
                            }}
                            value={product.maximumQuantity}
                          />
                        </Grid.Col>
                      )}

                      <Grid.Col span="content">
                        <Switch
                          mt="lg"
                          size="md"
                          onLabel="ranged"
                          offLabel="single"
                          disabled={
                            pageState.recipe.steps[stepIndex].products[productIndex].type === 'vessel' ||
                            pageState.stepHelpers[stepIndex].locked
                          }
                          value={pageState.stepHelpers[stepIndex].productIsRanged[productIndex] ? 'ranged' : 'single'}
                          onChange={() => {
                            dispatchPageEvent({
                              type: 'TOGGLE_PRODUCT_RANGE',
                              stepIndex: stepIndex,
                              productIndex: productIndex,
                            });
                          }}
                        />
                      </Grid.Col>

                      {product.type !== 'vessel' && (
                        <Grid.Col md="auto">
                          <Autocomplete
                            data-qa={`recipe-step-${stepIndex}-product-${productIndex}-measurement-unit-input`}
                            label="Measurement"
                            disabled={
                              (product.type === 'ingredient' && step.ingredients.length === 0) ||
                              pageState.stepHelpers[stepIndex].locked
                            }
                            value={pageState.stepHelpers[stepIndex].productMeasurementUnitQueries[productIndex]}
                            data={(
                              pageState.stepHelpers[stepIndex].productMeasurementUnitSuggestions[productIndex] || []
                            ).map((y: ValidMeasurementUnit) => ({
                              value: y.pluralName,
                              label: y.pluralName,
                            }))}
                            onItemSubmit={handleRecipeStepProductMeasurementUnitSelection(stepIndex, productIndex)}
                            onChange={handleRecipeStepProductMeasurementUnitQueryUpdate(stepIndex, productIndex)}
                            rightSection={
                              product.measurementUnitID &&
                              !recipeStepProductIsUsedInLaterStep(pageState.recipe, stepIndex, productIndex) && (
                                <IconCircleX
                                  size={18}
                                  color={product.measurementUnitID ? 'tomato' : 'gray'}
                                  onClick={() => {
                                    if (!product.measurementUnitID) {
                                      return;
                                    }

                                    dispatchPageEvent({
                                      type: 'UNSET_STEP_PRODUCT_MEASUREMENT_UNIT',
                                      stepIndex: stepIndex,
                                      productIndex: productIndex,
                                    });
                                  }}
                                />
                              )
                            }
                          />
                        </Grid.Col>
                      )}

                      <Grid.Col span="auto">
                        <TextInput
                          required
                          label="Name"
                          disabled={
                            !pageState.stepHelpers[stepIndex].productIsNamedManually[productIndex] ||
                            pageState.stepHelpers[stepIndex].locked
                          }
                          value={product.name}
                          onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                            dispatchPageEvent({
                              type: 'UPDATE_STEP_PRODUCT_NAME',
                              stepIndex: stepIndex,
                              productIndex: productIndex,
                              newName: event.target.value || '',
                            });
                          }}
                          rightSection={
                            <ActionIcon
                              style={{ float: 'right' }}
                              size="sm"
                              aria-label="add product"
                              disabled={manualProductNamingShouldBeDisabled(stepIndex, productIndex)}
                              onClick={() => {
                                if (manualProductNamingShouldBeDisabled(stepIndex, productIndex)) {
                                  return;
                                }

                                dispatchPageEvent({
                                  type: 'TOGGLE_MANUAL_PRODUCT_NAMING',
                                  stepIndex: stepIndex,
                                  productIndex: productIndex,
                                });
                              }}
                            >
                              {manualProductNamingShouldBeDisabled(stepIndex, productIndex) ? (
                                <IconEditOff />
                              ) : (
                                <IconEdit />
                              )}
                            </ActionIcon>
                          }
                        />
                      </Grid.Col>

                      {product.type !== 'vessel' && (
                        <Grid.Col span="auto">
                          <Select
                            label="In Vessel"
                            value={
                              pageState.recipe.steps[stepIndex].vessels[product.containedInVesselIndex ?? -1]?.name ||
                              ''
                            }
                            data={pageState.recipe.steps[stepIndex].vessels
                              .filter((vessel) => {
                                return vessel.name !== '';
                              })
                              .map((vessel) => {
                                return { value: vessel.name, label: vessel.name };
                              })}
                            disabled={pageState.stepHelpers[stepIndex].locked}
                            onChange={(value: string) => {
                              const selectedVessel = pageState.stepHelpers[stepIndex].selectedVessels.find(
                                (vessel: RecipeStepVessel | undefined) => {
                                  if (!vessel) {
                                    return false;
                                  }

                                  return vessel.name === value;
                                },
                              );

                              if (!selectedVessel) {
                                console.error('Could not find vessel with name', value);
                                return;
                              }

                              const selectedVesselIndex =
                                pageState.stepHelpers[stepIndex].selectedVessels.indexOf(selectedVessel);

                              dispatchPageEvent({
                                type: 'UPDATE_STEP_PRODUCT_VESSEL',
                                stepIndex: stepIndex,
                                productIndex: productIndex,
                                vesselIndex: selectedVesselIndex,
                              });
                            }}
                          />
                        </Grid.Col>
                      )}

                      <Grid.Col span="content" mt="xl">
                        <ActionIcon
                          mt={5}
                          variant="outline"
                          size="md"
                          style={{ float: 'right' }}
                          aria-label="remove step"
                          disabled={
                            pageState.stepHelpers[stepIndex].locked // ||
                            // (step.products || []).length === 1 ||
                            // (step.products || []).length - 1 !== productIndex ||
                            // (step.products[productIndex].type === 'vessel' &&
                            //   !pageState.stepHelpers[stepIndex].selectedPreparation?.consumesVessel)
                          }
                          onClick={() =>
                            dispatchPageEvent({
                              type: 'REMOVE_PRODUCT_FROM_STEP',
                              stepIndex: stepIndex,
                              productIndex: productIndex,
                            })
                          }
                        >
                          <IconTrash
                            size={16}
                            color={
                              (step.products || []).length === 1 ||
                              (step.products || []).length - 1 !== productIndex ||
                              step.products[productIndex].type === 'vessel'
                                ? 'gray'
                                : 'tomato'
                            }
                          />
                        </ActionIcon>
                      </Grid.Col>
                    </Grid>
                  );
                })}
              </Grid.Col>
            </Grid>

            <Grid>
              <Grid.Col span="auto">
                <Center>
                  <Button
                    mt="sm"
                    disabled={pageState.stepHelpers[stepIndex].locked}
                    onClick={() => {
                      dispatchPageEvent({
                        type: 'ADD_PRODUCT_TO_STEP',
                        stepIndex,
                      });
                    }}
                  >
                    Add Product
                  </Button>
                </Center>
              </Grid.Col>
            </Grid>
          </Card.Section>
        </Collapse>
      </Card>
    );
  };

  return (
    <AppLayout title="New Recipe" containerSize="xl">
      <form
        onSubmit={(e: FormEvent) => {
          e.preventDefault();
          submitRecipe();
        }}
      >
        <Grid>
          <Grid.Col md={4}>
            <Stack>
              <Stack>
                <TextInput
                  data-qa="recipe-name-input"
                  required
                  label="Name"
                  value={pageState.recipe.name}
                  onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                    dispatchPageEvent({ type: 'UPDATE_NAME', newName: event.target.value });
                  }}
                  mt="xs"
                />

                <TextInput
                  data-qa="recipe-slug-input"
                  label="Slug"
                  value={pageState.recipe.slug}
                  onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                    dispatchPageEvent({ type: 'UPDATE_SLUG', newSlug: event.target.value });
                  }}
                  mt="xs"
                />

                <NumberInput
                  data-qa="recipe-minimum-estimated-portions-input"
                  label="Estimated Portion Quantity"
                  required
                  value={pageState.recipe.minimumEstimatedPortions}
                  onChange={(value: number) => {
                    if (value <= 0) {
                      return;
                    }

                    dispatchPageEvent({ type: 'UPDATE_MINIMUM_ESTIMATED_PORTIONS', newPortions: value });
                  }}
                  mt="xs"
                />

                <Grid>
                  <Grid.Col span="auto">
                    <TextInput
                      data-qa="recipe-source-input"
                      label="Portion Name"
                      value={pageState.recipe.portionName}
                      onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                        dispatchPageEvent({ type: 'UPDATE_PORTION_NAME', newPortionName: event.target.value });
                      }}
                      mt="xs"
                    />
                  </Grid.Col>

                  <Grid.Col span="auto">
                    <TextInput
                      data-qa="recipe-source-input"
                      label="Plural Portion Name"
                      value={pageState.recipe.pluralPortionName}
                      onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                        dispatchPageEvent({
                          type: 'UPDATE_PLURAL_PORTION_NAME',
                          newPluralPortionName: event.target.value,
                        });
                      }}
                      mt="xs"
                    />
                  </Grid.Col>
                </Grid>

                <TextInput
                  data-qa="recipe-source-input"
                  label="Source"
                  value={pageState.recipe.source}
                  onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                    dispatchPageEvent({ type: 'UPDATE_SOURCE', newSource: event.target.value });
                  }}
                  mt="xs"
                />

                <Textarea
                  data-qa="recipe-description-input"
                  label="Description"
                  value={pageState.recipe.description}
                  onChange={(event: React.ChangeEvent<HTMLTextAreaElement>) => {
                    dispatchPageEvent({ type: 'UPDATE_DESCRIPTION', newDescription: event.target.value });
                  }}
                  minRows={4}
                  mt="xs"
                />

                <Button onClick={submitRecipe} disabled={recipeSubmissionShouldBeDisabled(pageState.recipe)}>
                  Save
                </Button>
              </Stack>

              <Divider />

              <Grid justify="space-between" align="center">
                <Grid.Col span="auto">
                  <Title order={4}>All Ingredients</Title>
                </Grid.Col>

                <Grid.Col span="auto">
                  <ActionIcon
                    data-qa="toggle-all-ingredients"
                    variant="outline"
                    size="sm"
                    style={{ float: 'right' }}
                    aria-label="show all ingredients"
                    onClick={() => {
                      dispatchPageEvent({ type: 'TOGGLE_SHOW_ALL_INGREDIENTS' });
                    }}
                  >
                    {(pageState.showIngredientsSummary && <IconChevronUp size={16} color="gray" />) || (
                      <IconChevronDown size={16} color="gray" />
                    )}
                  </ActionIcon>
                </Grid.Col>
              </Grid>

              <Divider />

              <Collapse sx={{ minHeight: '10rem' }} in={pageState.showIngredientsSummary}>
                <Box />
              </Collapse>

              <Grid justify="space-between" align="center">
                <Grid.Col span="auto">
                  <Title order={4}>All Instruments</Title>
                </Grid.Col>
                <Grid.Col span="auto">
                  <ActionIcon
                    data-qa="toggle-all-instruments"
                    variant="outline"
                    size="sm"
                    style={{ float: 'right' }}
                    aria-label="show all instruments"
                    onClick={() => {
                      dispatchPageEvent({ type: 'TOGGLE_SHOW_ALL_INSTRUMENTS' });
                    }}
                  >
                    {(pageState.showInstrumentsSummary && <IconChevronUp size={16} color="gray" />) || (
                      <IconChevronDown size={16} color="gray" />
                    )}
                  </ActionIcon>
                </Grid.Col>
              </Grid>

              <Divider />

              <Collapse sx={{ minHeight: '10rem' }} in={pageState.showInstrumentsSummary}>
                <Box />
              </Collapse>

              {pageState.recipe.steps.length > 1 && (
                <>
                  <Grid justify="space-between" align="center">
                    <Grid.Col span="auto">
                      <Title order={4}>Advanced Prep</Title>
                    </Grid.Col>
                    <Grid.Col span="auto">
                      <ActionIcon
                        data-qa="toggle-all-advanced-prep-steps"
                        variant="outline"
                        size="sm"
                        style={{ float: 'right' }}
                        aria-label="show advanced prep tasks"
                        onClick={() => {
                          dispatchPageEvent({ type: 'TOGGLE_SHOW_ADVANCED_PREP_STEPS' });
                        }}
                      >
                        {(pageState.showAdvancedPrepStepInputs && <IconChevronUp size={16} color="gray" />) || (
                          <IconChevronDown size={16} color="gray" />
                        )}
                      </ActionIcon>
                    </Grid.Col>
                  </Grid>

                  <Collapse sx={{ minHeight: '10rem' }} in={pageState.showAdvancedPrepStepInputs}>
                    <Box />
                  </Collapse>

                  <Divider />
                </>
              )}

              <>
                <Grid justify="space-between" align="center">
                  <Grid.Col span="auto">
                    <Title order={4}>Debug</Title>
                  </Grid.Col>
                  <Grid.Col span="auto">
                    <ActionIcon
                      data-qa="toggle-debug-menu"
                      variant="outline"
                      size="sm"
                      style={{ float: 'right' }}
                      aria-label="show debug menu"
                      onClick={() => {
                        dispatchPageEvent({ type: 'TOGGLE_SHOW_DEBUG_MENU' });
                      }}
                    >
                      {(pageState.showDebugMenu && <IconChevronUp size={16} color="gray" />) || (
                        <IconChevronDown size={16} color="gray" />
                      )}
                    </ActionIcon>
                  </Grid.Col>
                </Grid>

                <Collapse sx={{ minHeight: '10rem' }} in={pageState.showDebugMenu}>
                  <Grid>
                    <Grid.Col span={12}>
                      <Textarea
                        value={debugOutput}
                        autosize
                        minRows={2}
                        maxRows={10}
                        onChange={(event: React.ChangeEvent<HTMLTextAreaElement>) => {
                          setDebugOutput(event.target.value);
                        }}
                      />
                    </Grid.Col>
                    <Grid.Col span={6}>
                      <Button
                        fullWidth
                        onClick={() => {
                          setDebugOutput(JSON.stringify(pageState, null, 2));
                        }}
                      >
                        Dump State
                      </Button>
                    </Grid.Col>
                    <Grid.Col span={6}>
                      <Button
                        fullWidth
                        color="red"
                        onClick={() => {
                          dispatchPageEvent({
                            type: 'SET_PAGE_STATE',
                            newState: JSON.parse(debugOutput) as RecipeCreationPageState,
                          });
                          setDebugOutput('');
                        }}
                      >
                        Load State
                      </Button>
                    </Grid.Col>
                    <Grid.Col span={12}>
                      <Button
                        fullWidth
                        onClick={() => {
                          setDebugOutput(JSON.stringify(recipeCreationFormSchema.safeParse(pageState.recipe), null, 2));
                        }}
                      >
                        Recipe Submission Should Be Disabled
                      </Button>
                    </Grid.Col>
                  </Grid>
                </Collapse>
              </>
            </Stack>
          </Grid.Col>

          <Grid.Col span="auto" mt={'2.2rem'} mb="xl">
            {(pageState.recipe.steps || []).map(buildStep)}
            <Button
              fullWidth
              onClick={() => dispatchPageEvent({ type: 'ADD_STEP' })}
              mb="xl"
              disabled={addingStepsShouldBeDisabled(pageState)}
            >
              Add Step
            </Button>
          </Grid.Col>
        </Grid>
      </form>
    </AppLayout>
  );
}

export default RecipeCreator;
