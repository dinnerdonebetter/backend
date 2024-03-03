import { Reducer } from 'react';

import {
  ValidMeasurementUnit,
  ValidPreparation,
  RecipeStepIngredient,
  RecipeStepInstrument,
  ValidIngredientState,
  RecipeCreationRequestInput,
  RecipeStepCreationRequestInput,
  RecipeStepProductCreationRequestInput,
  RecipeStepIngredientCreationRequestInput,
  RecipeStepCompletionConditionCreationRequestInput,
  RecipeStepInstrumentCreationRequestInput,
  ValidRecipeStepProductType,
  RecipeStepVesselCreationRequestInput,
  RecipeStepVessel,
  ValidVessel,
} from '@dinnerdonebetter/models';

type RecipeCreationAction =
  | { type: 'SET_PAGE_STATE'; newState: RecipeCreationPageState }
  | { type: 'UPDATE_NAME'; newName: string }
  | { type: 'UPDATE_SLUG'; newSlug: string }
  | { type: 'UPDATE_DESCRIPTION'; newDescription: string }
  | { type: 'UPDATE_SOURCE'; newSource: string }
  | { type: 'UPDATE_PORTION_NAME'; newPortionName: string }
  | { type: 'UPDATE_PLURAL_PORTION_NAME'; newPluralPortionName: string }
  | { type: 'UPDATE_MINIMUM_ESTIMATED_PORTIONS'; newPortions?: number }
  | { type: 'TOGGLE_SHOW_ALL_INGREDIENTS' }
  | { type: 'TOGGLE_SHOW_ALL_INSTRUMENTS' }
  | { type: 'TOGGLE_SHOW_ADVANCED_PREP_STEPS' }
  | { type: 'TOGGLE_SHOW_DEBUG_MENU' }
  | { type: 'UPDATE_SUBMISSION_ERROR'; error: string }
  | { type: 'ADD_STEP' }
  | { type: 'TOGGLE_SHOW_STEP'; stepIndex: number }
  | { type: 'REMOVE_STEP'; stepIndex: number }
  | {
      type: 'SET_INGREDIENT_FOR_RECIPE_STEP_INGREDIENT';
      stepIndex: number;
      recipeStepIngredientIndex: number;
      selectedValidIngredient: RecipeStepIngredient;
      productOfRecipeStepIndex?: number;
      productOfRecipeStepProductIndex?: number;
    }
  | {
      type: 'SET_VESSEL_FOR_RECIPE_STEP_VESSEL';
      stepIndex: number;
      recipeStepIngredientIndex: number;
      selectedVessel: RecipeStepVessel;
      productOfRecipeStepIndex?: number;
      productOfRecipeStepProductIndex?: number;
    }
  | {
      type: 'ADD_INGREDIENT_TO_STEP';
      stepIndex: number;
    }
  | {
      type: 'ADD_PRODUCT_TO_STEP';
      stepIndex: number;
    }
  | {
      type: 'REMOVE_PRODUCT_FROM_STEP';
      stepIndex: number;
      productIndex: number;
    }
  | {
      type: 'UNSET_RECIPE_STEP_INGREDIENT';
      stepIndex: number;
      recipeStepIngredientIndex: number;
    }
  | {
      type: 'SET_VALID_INSTRUMENT_FOR_RECIPE_STEP_INSTRUMENT';
      stepIndex: number;
      recipeStepInstrumentIndex: number;
      selectedValidInstrument: RecipeStepInstrument;
      productOfRecipeStepIndex?: number;
      productOfRecipeStepProductIndex?: number;
    }
  | {
      type: 'SET_PRODUCT_FOR_RECIPE_STEP_VESSEL';
      stepIndex: number;
      recipeStepVesselIndex: number;
      selectedVessel: RecipeStepVessel;
    }
  | {
      type: 'SET_PRODUCT_INSTRUMENT_FOR_RECIPE_STEP_INSTRUMENT';
      stepIndex: number;
      recipeStepInstrumentIndex: number;
      selectedValidInstrument: RecipeStepInstrument;
      productOfRecipeStepIndex?: number;
      productOfRecipeStepProductIndex?: number;
    }
  | {
      type: 'SET_PRODUCT_INSTRUMENT_FOR_RECIPE_STEP_VESSEL';
      stepIndex: number;
      recipeStepVesselIndex: number;
      selectedValidInstrument: RecipeStepInstrument;
      productOfRecipeStepIndex?: number;
      productOfRecipeStepProductIndex?: number;
    }
  | {
      type: 'ADD_INSTRUMENT_TO_STEP';
      stepIndex: number;
    }
  | {
      type: 'ADD_VESSEL_TO_STEP';
      stepIndex: number;
    }
  | {
      type: 'SET_RECIPE_STEP_VESSEL_PREDICATE';
      stepIndex: number;
      recipeStepVesselIndex: number;
      vesselPreposition: string;
    }
  | {
      type: 'TOGGLE_RECIPE_STEP_VESSEL_PREDICATE';
      stepIndex: number;
      recipeStepVesselIndex: number;
    }
  | { type: 'REMOVE_INGREDIENT_FROM_STEP'; stepIndex: number; recipeStepIngredientIndex: number }
  | { type: 'REMOVE_INSTRUMENT_FROM_STEP'; stepIndex: number; recipeStepInstrumentIndex: number }
  | { type: 'REMOVE_VESSEL_FROM_STEP'; stepIndex: number; recipeStepVesselIndex: number }
  | { type: 'UPDATE_STEP_PREPARATION_QUERY'; stepIndex: number; newQuery: string }
  | { type: 'UPDATE_STEP_NOTES'; stepIndex: number; newNotes: string }
  | { type: 'UPDATE_STEP_EXPLICIT_INSTRUCTIONS'; stepIndex: number; newExplicitInstructions: string }
  | { type: 'UPDATE_STEP_MINIMUM_TIME_ESTIMATE'; stepIndex: number; newMinTimeEstimate: number }
  | { type: 'UPDATE_STEP_MAXIMUM_TIME_ESTIMATE'; stepIndex: number; newMaxTimeEstimate: number }
  | { type: 'UPDATE_STEP_MINIMUM_TEMPERATURE'; stepIndex: number; newMinTempInCelsius: number }
  | { type: 'UPDATE_STEP_MAXIMUM_TEMPERATURE'; stepIndex: number; newMaxTempInCelsius: number }
  | { type: 'UPDATE_STEP_INGREDIENT_QUERY'; stepIndex: number; recipeStepIngredientIndex: number; newQuery: string }
  | { type: 'TOGGLE_INGREDIENT_PRODUCT_STATE'; stepIndex: number; recipeStepIngredientIndex: number }
  | { type: 'TOGGLE_INSTRUMENT_PRODUCT_STATE'; stepIndex: number; recipeStepInstrumentIndex: number }
  | { type: 'TOGGLE_VESSEL_PRODUCT_STATE'; stepIndex: number; recipeStepVesselIndex: number }
  | {
      type: 'UPDATE_STEP_PRODUCT_MEASUREMENT_UNIT_QUERY';
      stepIndex: number;
      productIndex: number;
      newQuery: string;
    }
  | {
      type: 'UPDATE_STEP_VESSEL_INSTRUMENT_QUERY';
      stepIndex: number;
      vesselIndex: number;
      newQuery: string;
    }
  | {
      type: 'UPDATE_STEP_PRODUCT_MEASUREMENT_UNIT_SUGGESTIONS';
      stepIndex: number;
      productIndex: number;
      results: ValidMeasurementUnit[];
    }
  | {
      type: 'UPDATE_STEP_VESSEL_SUGGESTIONS';
      stepIndex: number;
      vesselIndex: number;
      results: ValidVessel[];
    }
  | {
      type: 'UNSET_STEP_PRODUCT_MEASUREMENT_UNIT';
      stepIndex: number;
      productIndex: number;
    }
  | {
      type: 'ADD_COMPLETION_CONDITION_TO_STEP';
      stepIndex: number;
    }
  | {
      type: 'UPDATE_COMPLETION_CONDITION_INGREDIENT_STATE_QUERY';
      stepIndex: number;
      conditionIndex: number;
      query: string;
    }
  | {
      type: 'UPDATE_COMPLETION_CONDITION_INGREDIENT_STATE_SUGGESTIONS';
      stepIndex: number;
      conditionIndex: number;
      results: ValidIngredientState[];
    }
  | {
      type: 'UPDATE_COMPLETION_CONDITION_INGREDIENT_STATE';
      stepIndex: number;
      conditionIndex: number;
      ingredientState: ValidIngredientState;
    }
  | {
      type: 'REMOVE_RECIPE_STEP_COMPLETION_CONDITION';
      stepIndex: number;
      conditionIndex: number;
    }
  | {
      type: 'UPDATE_STEP_PRODUCT_MEASUREMENT_UNIT';
      stepIndex: number;
      productIndex: number;
      measurementUnit: ValidMeasurementUnit;
    }
  | {
      type: 'UPDATE_STEP_VESSEL_INSTRUMENT';
      stepIndex: number;
      vesselIndex: number;
      selectedVessel: RecipeStepVessel;
    }
  | {
      type: 'UPDATE_STEP_PRODUCT_TYPE';
      stepIndex: number;
      productIndex: number;
      newType: ValidRecipeStepProductType;
    }
  | {
      type: 'UPDATE_STEP_PRODUCT_VESSEL';
      stepIndex: number;
      productIndex: number;
      vesselIndex: number;
    }
  | {
      type: 'UPDATE_STEP_INGREDIENT_MEASUREMENT_UNIT_SUGGESTIONS';
      stepIndex: number;
      recipeStepIngredientIndex: number;
      results: ValidMeasurementUnit[];
    }
  | {
      type: 'UPDATE_STEP_INGREDIENT_MEASUREMENT_UNIT';
      stepIndex: number;
      recipeStepIngredientIndex: number;
      measurementUnit?: ValidMeasurementUnit;
    }
  | {
      type: 'UPDATE_STEP_INGREDIENT_SUGGESTIONS';
      stepIndex: number;
      recipeStepIngredientIndex: number;
      results: RecipeStepIngredient[];
    }
  | {
      type: 'UPDATE_STEP_INSTRUMENT_SUGGESTIONS';
      stepIndex: number;
      recipeStepInstrumentIndex: number;
      results: RecipeStepInstrument[];
    }
  | {
      type: 'UPDATE_STEP_PREPARATION_SUGGESTIONS';
      stepIndex: number;
      results: ValidPreparation[];
    }
  | {
      type: 'UPDATE_STEP_INGREDIENT_MINIMUM_QUANTITY';
      stepIndex: number;
      recipeStepIngredientIndex: number;
      newAmount: number;
    }
  | {
      type: 'UPDATE_STEP_INGREDIENT_MAXIMUM_QUANTITY';
      stepIndex: number;
      recipeStepIngredientIndex: number;
      newAmount: number;
    }
  | {
      type: 'UPDATE_STEP_PRODUCT_MINIMUM_QUANTITY';
      stepIndex: number;
      productIndex: number;
      newAmount: number;
    }
  | {
      type: 'UPDATE_STEP_PRODUCT_MAXIMUM_QUANTITY';
      stepIndex: number;
      productIndex: number;
      newAmount: number;
    }
  | {
      type: 'UPDATE_STEP_INSTRUMENT_MINIMUM_QUANTITY';
      stepIndex: number;
      recipeStepInstrumentIndex: number;
      newAmount: number;
    }
  | {
      type: 'UPDATE_STEP_INSTRUMENT_MAXIMUM_QUANTITY';
      stepIndex: number;
      recipeStepInstrumentIndex: number;
      newAmount: number;
    }
  | {
      type: 'UPDATE_STEP_VESSEL_MINIMUM_QUANTITY';
      stepIndex: number;
      recipeStepVesselIndex: number;
      newAmount: number;
    }
  | {
      type: 'UPDATE_STEP_VESSEL_MAXIMUM_QUANTITY';
      stepIndex: number;
      recipeStepVesselIndex: number;
      newAmount: number;
    }
  | {
      type: 'UNSET_STEP_PREPARATION';
      stepIndex: number;
    }
  | {
      type: 'UPDATE_STEP_PREPARATION';
      stepIndex: number;
      selectedPreparation: ValidPreparation;
    }
  | {
      type: 'UPDATE_STEP_PRODUCT_NAME';
      newName: string;
      stepIndex: number;
      productIndex: number;
    }
  | { type: 'TOGGLE_INGREDIENT_RANGE'; stepIndex: number; recipeStepIngredientIndex: number }
  | { type: 'TOGGLE_INSTRUMENT_RANGE'; stepIndex: number; recipeStepInstrumentIndex: number }
  | { type: 'TOGGLE_VESSEL_RANGE'; stepIndex: number; recipeStepVesselIndex: number }
  | { type: 'TOGGLE_PRODUCT_RANGE'; stepIndex: number; productIndex: number }
  | {
      type: 'TOGGLE_MANUAL_PRODUCT_NAMING';
      stepIndex: number;
      productIndex: number;
    };

export class RecipeCreationPageState {
  submissionError: string | null = null;
  showIngredientsSummary: boolean = false;
  showInstrumentsSummary: boolean = false;
  showAdvancedPrepStepInputs: boolean = false;
  showDebugMenu: boolean = false;

  stepHelpers: StepHelper[] = [new StepHelper()];

  recipe: RecipeCreationRequestInput = new RecipeCreationRequestInput({
    minimumEstimatedPortions: 1,
    portionName: 'portion',
    pluralPortionName: 'portions',
    steps: [
      new RecipeStepCreationRequestInput({
        instruments: [new RecipeStepInstrumentCreationRequestInput()],
        ingredients: [new RecipeStepIngredientCreationRequestInput()],
        products: [
          new RecipeStepProductCreationRequestInput({
            type: 'ingredient',
          }),
        ],
      }),
    ],
  });
}

export class StepHelper {
  show: boolean = true;
  locked: boolean = false;

  // preparations
  preparationQuery: string = '';
  preparationSuggestions: ValidPreparation[] = [];
  selectedPreparation: ValidPreparation | null = null;

  // instruments
  instrumentIsRanged: boolean[] = [];
  instrumentSuggestions: RecipeStepInstrument[] = [];
  selectedInstruments: (RecipeStepInstrument | undefined)[] = [undefined];
  instrumentIsProduct: boolean[] = [false];

  // vessels
  vesselQueries: string[] = [''];
  vesselIsRanged: boolean[] = [];
  vesselSuggestions: ValidVessel[][] = [[]];
  selectedVessels: (RecipeStepVessel | undefined)[] = [undefined];
  vesselIsProduct: boolean[] = [false];

  // ingredients
  ingredientIsRanged: boolean[] = [false];
  ingredientQueries: string[] = [''];
  ingredientSuggestions: RecipeStepIngredient[][] = [[]];
  ingredientIsProduct: boolean[] = [false];
  selectedIngredients: (RecipeStepIngredient | undefined)[] = [undefined];
  ingredientMeasurementUnitSuggestions: ValidMeasurementUnit[][] = [[]];
  selectedMeasurementUnits: (ValidMeasurementUnit | undefined)[] = [undefined];

  // products
  productIsRanged: boolean[] = [false];
  productIsNamedManually: boolean[] = [false];
  productMeasurementUnitQueries: string[] = [''];
  productMeasurementUnitSuggestions: ValidMeasurementUnit[][] = [[]];
  selectedProductMeasurementUnits: (ValidMeasurementUnit | undefined)[] = [undefined];

  // completion condition ingredient states
  completionConditionIngredientStateQueries: string[] = [];
  completionConditionIngredientStateSuggestions: ValidIngredientState[][] = [];
}

export const useRecipeCreationReducer: Reducer<RecipeCreationPageState, RecipeCreationAction> = (
  state: RecipeCreationPageState,
  action: RecipeCreationAction,
): RecipeCreationPageState => {
  let newState: RecipeCreationPageState = structuredClone(state);

  switch (action.type) {
    case 'SET_PAGE_STATE': {
      newState = action.newState;
      break;
    }

    case 'TOGGLE_SHOW_ALL_INGREDIENTS': {
      newState.showIngredientsSummary = !state.showIngredientsSummary;
      break;
    }

    case 'TOGGLE_SHOW_ALL_INSTRUMENTS': {
      newState.showInstrumentsSummary = !state.showInstrumentsSummary;
      break;
    }

    case 'TOGGLE_SHOW_ADVANCED_PREP_STEPS': {
      newState.showAdvancedPrepStepInputs = !state.showAdvancedPrepStepInputs;
      break;
    }

    case 'TOGGLE_SHOW_DEBUG_MENU': {
      newState.showDebugMenu = !state.showDebugMenu;
      break;
    }

    case 'UPDATE_SUBMISSION_ERROR': {
      newState.submissionError = action.error;
      break;
    }

    case 'UPDATE_NAME': {
      newState.recipe.name = action.newName;
      break;
    }

    case 'UPDATE_SLUG': {
      newState.recipe.slug = action.newSlug;
      break;
    }

    case 'UPDATE_DESCRIPTION': {
      newState.recipe.description = action.newDescription;
      break;
    }

    case 'UPDATE_SOURCE': {
      newState.recipe.source = action.newSource;
      break;
    }

    case 'UPDATE_PORTION_NAME': {
      newState.recipe.portionName = action.newPortionName;
      break;
    }

    case 'UPDATE_PLURAL_PORTION_NAME': {
      newState.recipe.pluralPortionName = action.newPluralPortionName;
      break;
    }

    case 'UPDATE_MINIMUM_ESTIMATED_PORTIONS': {
      if ((action.newPortions || -1) > 0) {
        newState = { ...state, recipe: { ...state.recipe, minimumEstimatedPortions: action.newPortions! } };
      }
      break;
    }

    case 'ADD_STEP': {
      const newStepHelper = new StepHelper();

      newState.stepHelpers = [...state.stepHelpers, newStepHelper];
      newState.recipe.steps.push(
        new RecipeStepCreationRequestInput({
          instruments: [new RecipeStepInstrumentCreationRequestInput()],
          ingredients: [new RecipeStepIngredientCreationRequestInput()],
          products: [
            new RecipeStepProductCreationRequestInput({
              type: 'ingredient',
            }),
          ],
          completionConditions: [],
        }),
      );

      break;
    }

    case 'REMOVE_STEP': {
      newState.stepHelpers = newState.stepHelpers.filter(
        (_stepHelper: StepHelper, index: number) => index !== action.stepIndex,
      );
      newState.recipe.steps = newState.recipe.steps.filter(
        (_step: RecipeStepCreationRequestInput, index: number) => index !== action.stepIndex,
      );
      break;
    }

    case 'TOGGLE_SHOW_STEP': {
      newState.stepHelpers[action.stepIndex].show = !newState.stepHelpers[action.stepIndex].show;
      break;
    }

    case 'SET_INGREDIENT_FOR_RECIPE_STEP_INGREDIENT': {
      newState.stepHelpers[action.stepIndex].ingredientQueries[action.recipeStepIngredientIndex] =
        action.selectedValidIngredient.name;
      newState.stepHelpers[action.stepIndex].ingredientSuggestions[action.recipeStepIngredientIndex] = [];
      newState.stepHelpers[action.stepIndex].ingredientMeasurementUnitSuggestions[action.recipeStepIngredientIndex] =
        [];
      newState.stepHelpers[action.stepIndex].selectedIngredients[action.recipeStepIngredientIndex] =
        action.selectedValidIngredient;

      newState.recipe.steps[action.stepIndex].products[0].name = `${
        newState.stepHelpers[action.stepIndex].selectedPreparation?.pastTense
      } ${new Intl.ListFormat('en').format(
        newState.recipe.steps[action.stepIndex].ingredients.map(
          (x: RecipeStepIngredientCreationRequestInput, i: number) =>
            newState.stepHelpers[action.stepIndex]?.selectedIngredients[i]?.name || x.name,
        ),
      )}`;

      const newIngredient = new RecipeStepIngredientCreationRequestInput({
        name: action.selectedValidIngredient.name,
        ingredientID: action.selectedValidIngredient.ingredient?.id,
        measurementUnitID: action.selectedValidIngredient.measurementUnit.id,
        minimumQuantity: action.selectedValidIngredient.minimumQuantity,
        maximumQuantity: action.selectedValidIngredient.maximumQuantity,
        productOfRecipeStepIndex: action.productOfRecipeStepIndex,
        productOfRecipeStepProductIndex: action.productOfRecipeStepProductIndex,
      });

      if (action.productOfRecipeStepIndex && action.productOfRecipeStepProductIndex) {
        newIngredient.measurementUnitID =
          newState.recipe.steps[action.productOfRecipeStepIndex].products[action.productOfRecipeStepProductIndex]
            .measurementUnitID || '';
      }

      newState.recipe.steps[action.stepIndex].ingredients[action.recipeStepIngredientIndex] =
        new RecipeStepIngredientCreationRequestInput({
          name: action.selectedValidIngredient.name,
          ingredientID: action.selectedValidIngredient.ingredient?.id,
          measurementUnitID: action.selectedValidIngredient.measurementUnit.id,
          minimumQuantity: action.selectedValidIngredient.minimumQuantity,
          maximumQuantity: action.selectedValidIngredient.maximumQuantity,
          productOfRecipeStepIndex: action.productOfRecipeStepIndex,
          productOfRecipeStepProductIndex: action.productOfRecipeStepProductIndex,
        });

      break;
    }

    case 'SET_VESSEL_FOR_RECIPE_STEP_VESSEL': {
      newState.stepHelpers[action.stepIndex].vesselQueries[action.recipeStepIngredientIndex] =
        action.selectedVessel.name;
      newState.stepHelpers[action.stepIndex].vesselSuggestions[action.recipeStepIngredientIndex] = [];
      newState.stepHelpers[action.stepIndex].selectedVessels[action.recipeStepIngredientIndex] = action.selectedVessel;

      newState.recipe.steps[action.stepIndex].products.push(
        new RecipeStepProductCreationRequestInput({
          type: 'vessel',
          name: action.selectedVessel.name,
        }),
      );

      newState.recipe.steps[action.stepIndex].vessels[action.recipeStepIngredientIndex] =
        new RecipeStepVesselCreationRequestInput({
          name: action.selectedVessel.name,
          vesselID: action.selectedVessel.vessel?.id,
          minimumQuantity: action.selectedVessel.minimumQuantity,
          maximumQuantity: action.selectedVessel.maximumQuantity,
          productOfRecipeStepIndex: action.productOfRecipeStepIndex,
          productOfRecipeStepProductIndex: action.productOfRecipeStepProductIndex,
          // TODO: vesselPreposition:
        });

      break;
    }

    case 'ADD_INGREDIENT_TO_STEP': {
      newState.stepHelpers[action.stepIndex].ingredientIsRanged.push(false);
      newState.stepHelpers[action.stepIndex].ingredientQueries.push('');
      newState.stepHelpers[action.stepIndex].ingredientSuggestions.push([]);
      newState.stepHelpers[action.stepIndex].ingredientIsProduct.push(false);
      newState.stepHelpers[action.stepIndex].ingredientMeasurementUnitSuggestions.push([]);
      newState.stepHelpers[action.stepIndex].selectedIngredients.push(new RecipeStepIngredient());

      newState.recipe.steps[action.stepIndex].ingredients.push(new RecipeStepIngredientCreationRequestInput());
      break;
    }

    case 'UNSET_RECIPE_STEP_INGREDIENT': {
      newState.stepHelpers[action.stepIndex].ingredientQueries[action.recipeStepIngredientIndex] = '';
      newState.stepHelpers[action.stepIndex].ingredientSuggestions[action.recipeStepIngredientIndex] = [];
      newState.stepHelpers[action.stepIndex].ingredientMeasurementUnitSuggestions[action.recipeStepIngredientIndex] =
        [];
      newState.stepHelpers[action.stepIndex].selectedIngredients[action.recipeStepIngredientIndex] = undefined;
      break;
    }

    case 'REMOVE_INGREDIENT_FROM_STEP': {
      newState.stepHelpers[action.stepIndex].ingredientIsRanged = newState.stepHelpers[
        action.stepIndex
      ].ingredientIsRanged.filter(
        (_x: boolean, recipeStepIngredientIndex: number) =>
          recipeStepIngredientIndex !== action.recipeStepIngredientIndex,
      );
      newState.stepHelpers[action.stepIndex].ingredientQueries = newState.stepHelpers[
        action.stepIndex
      ].ingredientQueries.filter(
        (_x: string, recipeStepIngredientIndex: number) =>
          recipeStepIngredientIndex !== action.recipeStepIngredientIndex,
      );
      newState.stepHelpers[action.stepIndex].ingredientSuggestions = newState.stepHelpers[
        action.stepIndex
      ].ingredientSuggestions.filter(
        (_x: RecipeStepIngredient[], recipeStepIngredientIndex: number) =>
          recipeStepIngredientIndex !== action.recipeStepIngredientIndex,
      );
      newState.stepHelpers[action.stepIndex].ingredientIsProduct = newState.stepHelpers[
        action.stepIndex
      ].ingredientIsProduct.filter(
        (_x: boolean, recipeStepIngredientIndex: number) =>
          recipeStepIngredientIndex !== action.recipeStepIngredientIndex,
      );
      newState.stepHelpers[action.stepIndex].selectedIngredients = newState.stepHelpers[
        action.stepIndex
      ].selectedIngredients.filter(
        (_x: RecipeStepIngredient | undefined, recipeStepIngredientIndex: number) =>
          recipeStepIngredientIndex !== action.recipeStepIngredientIndex,
      );
      newState.stepHelpers[action.stepIndex].ingredientMeasurementUnitSuggestions = newState.stepHelpers[
        action.stepIndex
      ].ingredientMeasurementUnitSuggestions.filter(
        (_x: ValidMeasurementUnit[], ingredientIndex: number) => ingredientIndex !== action.recipeStepIngredientIndex,
      );
      newState.stepHelpers[action.stepIndex].selectedMeasurementUnits = newState.stepHelpers[
        action.stepIndex
      ].selectedMeasurementUnits.filter(
        (_x: ValidMeasurementUnit | undefined, recipeStepIngredientIndex: number) =>
          recipeStepIngredientIndex !== action.recipeStepIngredientIndex,
      );

      newState.recipe.steps[action.stepIndex].ingredients = newState.recipe.steps[action.stepIndex].ingredients.filter(
        (_x: RecipeStepIngredientCreationRequestInput, recipeStepIngredientIndex: number) =>
          recipeStepIngredientIndex !== action.recipeStepIngredientIndex,
      );

      break;
    }

    case 'ADD_PRODUCT_TO_STEP': {
      newState.stepHelpers[action.stepIndex].productIsRanged.push(false);
      newState.stepHelpers[action.stepIndex].productIsNamedManually.push(false);
      newState.stepHelpers[action.stepIndex].productMeasurementUnitQueries.push('');
      newState.stepHelpers[action.stepIndex].productMeasurementUnitSuggestions.push([]);
      newState.stepHelpers[action.stepIndex].selectedProductMeasurementUnits.push(undefined);

      newState.recipe.steps[action.stepIndex].products.push(new RecipeStepProductCreationRequestInput());
      break;
    }

    case 'REMOVE_PRODUCT_FROM_STEP': {
      newState.stepHelpers[action.stepIndex].productIsRanged = newState.stepHelpers[
        action.stepIndex
      ].productIsRanged.filter((_x: boolean, productIndex: number) => productIndex !== action.productIndex);
      newState.stepHelpers[action.stepIndex].productIsNamedManually = newState.stepHelpers[
        action.stepIndex
      ].productIsNamedManually.filter((_x: boolean, productIndex: number) => productIndex !== action.productIndex);
      newState.stepHelpers[action.stepIndex].productMeasurementUnitQueries = newState.stepHelpers[
        action.stepIndex
      ].productMeasurementUnitQueries.filter(
        (_x: string, productIndex: number) => productIndex !== action.productIndex,
      );
      newState.stepHelpers[action.stepIndex].productMeasurementUnitSuggestions = newState.stepHelpers[
        action.stepIndex
      ].productMeasurementUnitSuggestions.filter(
        (_x: ValidMeasurementUnit[], productIndex: number) => productIndex !== action.productIndex,
      );
      newState.stepHelpers[action.stepIndex].selectedProductMeasurementUnits = newState.stepHelpers[
        action.stepIndex
      ].selectedProductMeasurementUnits.filter(
        (_x: ValidMeasurementUnit | undefined, productIndex: number) => productIndex !== action.productIndex,
      );

      newState.recipe.steps[action.stepIndex].products = newState.recipe.steps[action.stepIndex].products.filter(
        (_product: RecipeStepProductCreationRequestInput, productIndex: number) => productIndex !== action.productIndex,
      );

      break;
    }

    case 'SET_VALID_INSTRUMENT_FOR_RECIPE_STEP_INSTRUMENT': {
      newState.stepHelpers[action.stepIndex].instrumentIsRanged[action.recipeStepInstrumentIndex] = false;
      newState.stepHelpers[action.stepIndex].selectedInstruments[action.recipeStepInstrumentIndex] =
        action.selectedValidInstrument;
      newState.stepHelpers[action.stepIndex].instrumentIsProduct[action.recipeStepInstrumentIndex] = false;
      newState.recipe.steps[action.stepIndex].instruments[action.recipeStepInstrumentIndex] =
        new RecipeStepInstrumentCreationRequestInput({
          name: action.selectedValidInstrument.name,
          instrumentID: action.selectedValidInstrument.instrument?.id,
          minimumQuantity: 1,
        });

      break;
    }

    case 'SET_PRODUCT_FOR_RECIPE_STEP_VESSEL': {
      newState.stepHelpers[action.stepIndex].vesselIsRanged[action.recipeStepVesselIndex] = false;
      newState.stepHelpers[action.stepIndex].selectedVessels[action.recipeStepVesselIndex] = action.selectedVessel;
      newState.stepHelpers[action.stepIndex].vesselIsProduct[action.recipeStepVesselIndex] = true;
      newState.recipe.steps[action.stepIndex].vessels[action.recipeStepVesselIndex] =
        new RecipeStepVesselCreationRequestInput({
          name: action.selectedVessel.name,
          vesselID: action.selectedVessel.vessel?.id,
          minimumQuantity: 1,
        });

      // TODO: upsert product instead of always pushing
      newState.recipe.steps[action.stepIndex].products.push(
        new RecipeStepProductCreationRequestInput({
          name: action.selectedVessel.name,
          minimumQuantity:
            newState.recipe.steps[action.stepIndex].vessels[action.recipeStepVesselIndex].minimumQuantity,
          maximumQuantity:
            newState.recipe.steps[action.stepIndex].vessels[action.recipeStepVesselIndex].maximumQuantity,
          type: 'vessel',
        }),
      );

      break;
    }

    case 'SET_PRODUCT_INSTRUMENT_FOR_RECIPE_STEP_INSTRUMENT': {
      newState.stepHelpers[action.stepIndex].instrumentIsRanged[action.recipeStepInstrumentIndex] = false;
      newState.stepHelpers[action.stepIndex].selectedInstruments[action.recipeStepInstrumentIndex] =
        action.selectedValidInstrument;
      newState.recipe.steps[action.stepIndex].instruments[action.recipeStepInstrumentIndex] =
        new RecipeStepInstrumentCreationRequestInput({
          name: action.selectedValidInstrument.name,
          instrumentID: action.selectedValidInstrument.instrument?.id,
          minimumQuantity: 1,
          productOfRecipeStepIndex: action.productOfRecipeStepIndex,
          productOfRecipeStepProductIndex: action.productOfRecipeStepProductIndex,
        });

      break;
    }

    case 'SET_PRODUCT_INSTRUMENT_FOR_RECIPE_STEP_VESSEL': {
      newState.stepHelpers[action.stepIndex].vesselIsRanged[action.recipeStepVesselIndex] = false;
      newState.stepHelpers[action.stepIndex].selectedInstruments[action.recipeStepVesselIndex] =
        action.selectedValidInstrument;
      newState.recipe.steps[action.stepIndex].vessels[action.recipeStepVesselIndex] =
        new RecipeStepVesselCreationRequestInput({
          name: action.selectedValidInstrument.name,
          vesselID: action.selectedValidInstrument.instrument?.id,
          minimumQuantity: 1,
          productOfRecipeStepIndex: action.productOfRecipeStepIndex,
          productOfRecipeStepProductIndex: action.productOfRecipeStepProductIndex,
          vesselPreposition: 'in',
        });

      break;
    }

    case 'ADD_INSTRUMENT_TO_STEP': {
      newState.stepHelpers[action.stepIndex].instrumentIsRanged.push(false);
      newState.stepHelpers[action.stepIndex].selectedInstruments.push(undefined);
      newState.stepHelpers[action.stepIndex].instrumentIsProduct.push(false);
      newState.recipe.steps[action.stepIndex].instruments.push(new RecipeStepInstrumentCreationRequestInput());
      break;
    }

    case 'ADD_VESSEL_TO_STEP': {
      newState.stepHelpers[action.stepIndex].vesselQueries.push('');
      newState.stepHelpers[action.stepIndex].vesselIsRanged.push(false);
      newState.stepHelpers[action.stepIndex].selectedVessels.push(undefined);
      newState.stepHelpers[action.stepIndex].vesselIsProduct.push(false);
      newState.recipe.steps[action.stepIndex].vessels.push(
        new RecipeStepVesselCreationRequestInput({ vesselPreposition: 'in' }),
      );
      break;
    }

    case 'SET_RECIPE_STEP_VESSEL_PREDICATE': {
      newState.recipe.steps[action.stepIndex].vessels[action.recipeStepVesselIndex].vesselPreposition =
        action.vesselPreposition;
      break;
    }

    case 'REMOVE_INSTRUMENT_FROM_STEP': {
      newState.stepHelpers[action.stepIndex].instrumentIsRanged = newState.stepHelpers[
        action.stepIndex
      ].instrumentIsRanged.filter(
        (_isRanged: boolean, instrumentIndex: number) => instrumentIndex !== action.recipeStepInstrumentIndex,
      );
      newState.stepHelpers[action.stepIndex].selectedInstruments = newState.stepHelpers[
        action.stepIndex
      ].selectedInstruments.filter(
        (_instrument: RecipeStepInstrument | undefined, instrumentIndex: number) =>
          instrumentIndex !== action.recipeStepInstrumentIndex,
      );
      newState.recipe.steps[action.stepIndex].instruments = newState.recipe.steps[action.stepIndex].instruments.filter(
        (_instrument: RecipeStepInstrumentCreationRequestInput, instrumentIndex: number) =>
          instrumentIndex !== action.recipeStepInstrumentIndex,
      );
      break;
    }

    case 'REMOVE_VESSEL_FROM_STEP': {
      newState.stepHelpers[action.stepIndex].vesselQueries = newState.stepHelpers[
        action.stepIndex
      ].vesselQueries.filter((_query: string, vesselIndex: number) => vesselIndex !== action.recipeStepVesselIndex);
      newState.stepHelpers[action.stepIndex].vesselIsRanged = newState.stepHelpers[
        action.stepIndex
      ].vesselIsRanged.filter(
        (_isRanged: boolean, vesselIndex: number) => vesselIndex !== action.recipeStepVesselIndex,
      );
      newState.stepHelpers[action.stepIndex].selectedInstruments = newState.stepHelpers[
        action.stepIndex
      ].selectedInstruments.filter(
        (_vessel: RecipeStepInstrument | undefined, vesselIndex: number) =>
          vesselIndex !== action.recipeStepVesselIndex,
      );
      newState.recipe.steps[action.stepIndex].vessels = newState.recipe.steps[action.stepIndex].vessels.filter(
        (_vessel: RecipeStepVesselCreationRequestInput, vesselIndex: number) =>
          vesselIndex !== action.recipeStepVesselIndex,
      );
      break;
    }

    case 'UPDATE_STEP_INGREDIENT_SUGGESTIONS': {
      newState.stepHelpers[action.stepIndex].ingredientSuggestions[action.recipeStepIngredientIndex] =
        action.results || [];
      break;
    }

    case 'UPDATE_STEP_INSTRUMENT_SUGGESTIONS': {
      newState.stepHelpers[action.stepIndex].instrumentSuggestions = action.results || [];
      break;
    }

    case 'UPDATE_STEP_PREPARATION_QUERY': {
      newState.stepHelpers[action.stepIndex].preparationQuery = action.newQuery;
      break;
    }

    case 'UPDATE_STEP_PREPARATION_SUGGESTIONS': {
      newState.stepHelpers[action.stepIndex].preparationSuggestions = action.results || [];
      break;
    }

    case 'UPDATE_STEP_INGREDIENT_MEASUREMENT_UNIT_SUGGESTIONS': {
      newState.stepHelpers[action.stepIndex].ingredientMeasurementUnitSuggestions[action.recipeStepIngredientIndex] =
        action.results || [];
      break;
    }

    case 'UPDATE_STEP_PRODUCT_MEASUREMENT_UNIT_QUERY': {
      newState.stepHelpers[action.stepIndex].productMeasurementUnitQueries[action.productIndex] = action.newQuery;
      break;
    }

    case 'UPDATE_STEP_VESSEL_INSTRUMENT_QUERY': {
      newState.stepHelpers[action.stepIndex].vesselQueries[action.vesselIndex] = action.newQuery;
      break;
    }

    case 'UPDATE_STEP_PRODUCT_MEASUREMENT_UNIT_SUGGESTIONS': {
      newState.stepHelpers[action.stepIndex].productMeasurementUnitSuggestions[action.productIndex] =
        action.results || [];
      break;
    }

    case 'UPDATE_STEP_VESSEL_SUGGESTIONS': {
      newState.stepHelpers[action.stepIndex].vesselSuggestions[action.vesselIndex] = action.results || [];
      break;
    }

    case 'UNSET_STEP_PRODUCT_MEASUREMENT_UNIT': {
      newState.stepHelpers[action.stepIndex].productMeasurementUnitQueries[action.productIndex] = '';
      newState.stepHelpers[action.stepIndex].productMeasurementUnitSuggestions[action.productIndex] = [];
      newState.stepHelpers[action.stepIndex].selectedProductMeasurementUnits[action.productIndex] = undefined;
      newState.recipe.steps[action.stepIndex].products[action.productIndex].measurementUnitID = '';
      break;
    }

    case 'ADD_COMPLETION_CONDITION_TO_STEP': {
      newState.stepHelpers[action.stepIndex].completionConditionIngredientStateQueries.push('');
      newState.stepHelpers[action.stepIndex].completionConditionIngredientStateSuggestions.push([]);
      newState.recipe.steps[action.stepIndex].completionConditions.push(
        new RecipeStepCompletionConditionCreationRequestInput(),
      );
      break;
    }

    case 'UPDATE_COMPLETION_CONDITION_INGREDIENT_STATE_QUERY': {
      newState.stepHelpers[action.stepIndex].completionConditionIngredientStateQueries[action.conditionIndex] =
        action.query;
      break;
    }

    case 'UPDATE_COMPLETION_CONDITION_INGREDIENT_STATE_SUGGESTIONS': {
      newState.stepHelpers[action.stepIndex].completionConditionIngredientStateSuggestions[action.conditionIndex] =
        action.results || [];
      break;
    }

    case 'UPDATE_COMPLETION_CONDITION_INGREDIENT_STATE': {
      if (!action.ingredientState) {
        console.error("couldn't find ingredient state to add");
        break;
      }

      newState.stepHelpers[action.stepIndex].completionConditionIngredientStateQueries[action.conditionIndex] =
        action.ingredientState.name;
      newState.stepHelpers[action.stepIndex].completionConditionIngredientStateSuggestions[action.conditionIndex] = [];

      newState.recipe.steps[action.stepIndex].completionConditions[action.conditionIndex].ingredientState =
        action.ingredientState!.id;
      break;
    }

    case 'REMOVE_RECIPE_STEP_COMPLETION_CONDITION': {
      newState.stepHelpers[action.stepIndex].completionConditionIngredientStateQueries = newState.stepHelpers[
        action.stepIndex
      ].completionConditionIngredientStateQueries.filter(
        (_: string, conditionIndex: number) => conditionIndex !== action.conditionIndex,
      );

      newState.stepHelpers[action.stepIndex].completionConditionIngredientStateSuggestions = newState.stepHelpers[
        action.stepIndex
      ].completionConditionIngredientStateSuggestions.filter(
        (_: ValidIngredientState[], conditionIndex: number) => conditionIndex !== action.conditionIndex,
      );

      newState.recipe.steps[action.stepIndex].completionConditions = newState.recipe.steps[
        action.stepIndex
      ].completionConditions.filter(
        (_: RecipeStepCompletionConditionCreationRequestInput, conditionIndex: number) =>
          conditionIndex !== action.conditionIndex,
      );
      break;
    }

    case 'UPDATE_STEP_INGREDIENT_MEASUREMENT_UNIT': {
      if (!action.measurementUnit) {
        console.error("couldn't find measurement unit to add for step ingredient");
        break;
      }

      newState.stepHelpers[action.stepIndex].selectedMeasurementUnits[action.recipeStepIngredientIndex] =
        action.measurementUnit!;
      newState.recipe.steps[action.stepIndex].ingredients[action.recipeStepIngredientIndex].measurementUnitID =
        action.measurementUnit!.id;
      break;
    }

    case 'UPDATE_STEP_PRODUCT_MEASUREMENT_UNIT': {
      if (!action.measurementUnit) {
        console.error("couldn't find measurement unit to add");
        break;
      }

      newState.stepHelpers[action.stepIndex].selectedProductMeasurementUnits[action.productIndex] =
        action.measurementUnit;
      newState.recipe.steps[action.stepIndex].products[action.productIndex].measurementUnitID =
        action.measurementUnit!.id;
      break;
    }

    case 'UPDATE_STEP_VESSEL_INSTRUMENT': {
      if (!action.selectedVessel) {
        console.error("couldn't find measurement unit to add");
        break;
      }

      newState.stepHelpers[action.stepIndex].selectedVessels[action.vesselIndex] = action.selectedVessel;
      newState.recipe.steps[action.stepIndex].vessels[action.vesselIndex] = new RecipeStepVesselCreationRequestInput({
        ...action.selectedVessel,
      });

      // TODO: upsert product instead of always pushing
      newState.recipe.steps[action.stepIndex].products.push(
        new RecipeStepProductCreationRequestInput({
          name: action.selectedVessel.name,
          minimumQuantity: newState.recipe.steps[action.stepIndex].vessels[action.vesselIndex].minimumQuantity,
          maximumQuantity: newState.recipe.steps[action.stepIndex].vessels[action.vesselIndex].maximumQuantity,
          type: 'vessel',
        }),
      );

      break;
    }

    case 'UPDATE_STEP_PRODUCT_TYPE': {
      newState.stepHelpers[action.stepIndex].productMeasurementUnitSuggestions[action.productIndex] = [];
      newState.recipe.steps[action.stepIndex].products[action.productIndex].type = action.newType;
      newState.recipe.steps[action.stepIndex].products[action.productIndex].minimumQuantity = 1;
      break;
    }

    case 'UPDATE_STEP_PRODUCT_VESSEL': {
      newState.recipe.steps[action.stepIndex].products[action.productIndex].containedInVesselIndex = action.vesselIndex;
      break;
    }

    case 'UPDATE_STEP_INGREDIENT_MINIMUM_QUANTITY': {
      newState.recipe.steps[action.stepIndex].ingredients[action.recipeStepIngredientIndex].minimumQuantity =
        action.newAmount;
      break;
    }

    case 'UPDATE_STEP_INGREDIENT_MAXIMUM_QUANTITY': {
      newState.recipe.steps[action.stepIndex].ingredients[action.recipeStepIngredientIndex].maximumQuantity =
        action.newAmount;
      break;
    }

    case 'UPDATE_STEP_PRODUCT_MINIMUM_QUANTITY': {
      newState.recipe.steps[action.stepIndex].products[action.productIndex].minimumQuantity = action.newAmount;
      break;
    }

    case 'UPDATE_STEP_PRODUCT_MAXIMUM_QUANTITY': {
      newState.recipe.steps[action.stepIndex].products[action.productIndex].maximumQuantity = action.newAmount;
      break;
    }

    case 'UPDATE_STEP_INSTRUMENT_MINIMUM_QUANTITY': {
      newState.recipe.steps[action.stepIndex].instruments[action.recipeStepInstrumentIndex].minimumQuantity =
        action.newAmount;
      break;
    }

    case 'UPDATE_STEP_INSTRUMENT_MAXIMUM_QUANTITY': {
      newState.recipe.steps[action.stepIndex].instruments[action.recipeStepInstrumentIndex].maximumQuantity =
        action.newAmount;
      break;
    }

    case 'UPDATE_STEP_VESSEL_MINIMUM_QUANTITY': {
      newState.recipe.steps[action.stepIndex].vessels[action.recipeStepVesselIndex].minimumQuantity = action.newAmount;
      break;
    }

    case 'UPDATE_STEP_VESSEL_MAXIMUM_QUANTITY': {
      newState.recipe.steps[action.stepIndex].vessels[action.recipeStepVesselIndex].maximumQuantity = action.newAmount;
      break;
    }

    case 'UPDATE_STEP_PRODUCT_NAME': {
      newState.recipe.steps[action.stepIndex].products[action.productIndex].name = action.newName;
      break;
    }

    case 'UNSET_STEP_PREPARATION': {
      // we need to effectively reset the step, since the preparation is the root.
      newState.stepHelpers[action.stepIndex].selectedPreparation = null;
      newState.stepHelpers[action.stepIndex].preparationQuery = '';
      newState.stepHelpers[action.stepIndex].ingredientSuggestions = [];
      newState.stepHelpers[action.stepIndex].preparationSuggestions = [];
      newState.stepHelpers[action.stepIndex].instrumentSuggestions = [];

      newState.recipe.steps[action.stepIndex].instruments = [new RecipeStepInstrumentCreationRequestInput()];

      newState.recipe.steps[action.stepIndex].ingredients = [new RecipeStepIngredientCreationRequestInput()];

      newState.recipe.steps[action.stepIndex].products = [
        new RecipeStepProductCreationRequestInput({
          type: 'ingredient',
        }),
      ];

      break;
    }

    case 'UPDATE_STEP_PREPARATION': {
      // we need to effectively reset the step, since the preparation is the root.
      newState.stepHelpers[action.stepIndex].selectedPreparation = action.selectedPreparation;
      newState.stepHelpers[action.stepIndex].preparationQuery = action.selectedPreparation.name;
      newState.stepHelpers[action.stepIndex].preparationSuggestions = [];

      newState.recipe.steps[action.stepIndex].preparationID = action.selectedPreparation.id;
      newState.recipe.steps[action.stepIndex].instruments = [new RecipeStepInstrumentCreationRequestInput()];
      newState.recipe.steps[action.stepIndex].products = [
        new RecipeStepProductCreationRequestInput({
          type: 'ingredient',
        }),
      ];
      newState.recipe.steps[action.stepIndex].ingredients = [new RecipeStepIngredientCreationRequestInput()];
      newState.recipe.steps[action.stepIndex].completionConditions = [];
      break;
    }

    case 'UPDATE_STEP_NOTES': {
      newState.recipe.steps[action.stepIndex].notes = action.newNotes;
      break;
    }

    case 'UPDATE_STEP_EXPLICIT_INSTRUCTIONS': {
      newState.recipe.steps[action.stepIndex].explicitInstructions = action.newExplicitInstructions;
      break;
    }

    case 'UPDATE_STEP_MINIMUM_TIME_ESTIMATE': {
      newState.recipe.steps[action.stepIndex].minimumEstimatedTimeInSeconds = action.newMinTimeEstimate;
      break;
    }

    case 'UPDATE_STEP_MAXIMUM_TIME_ESTIMATE': {
      newState.recipe.steps[action.stepIndex].maximumEstimatedTimeInSeconds = action.newMaxTimeEstimate;
      break;
    }

    case 'UPDATE_STEP_MINIMUM_TEMPERATURE': {
      newState.recipe.steps[action.stepIndex].minimumTemperatureInCelsius = action.newMinTempInCelsius;
      break;
    }

    case 'UPDATE_STEP_MAXIMUM_TEMPERATURE': {
      newState.recipe.steps[action.stepIndex].maximumTemperatureInCelsius = action.newMaxTempInCelsius;
      break;
    }

    case 'UPDATE_STEP_INGREDIENT_QUERY': {
      newState.stepHelpers[action.stepIndex].ingredientQueries[action.recipeStepIngredientIndex] = action.newQuery;
      break;
    }

    case 'TOGGLE_INGREDIENT_PRODUCT_STATE': {
      newState.stepHelpers[action.stepIndex].ingredientIsProduct[action.recipeStepIngredientIndex] =
        !newState.stepHelpers[action.stepIndex].ingredientIsProduct[action.recipeStepIngredientIndex];
      break;
    }

    case 'TOGGLE_INSTRUMENT_PRODUCT_STATE': {
      newState.stepHelpers[action.stepIndex].instrumentIsProduct[action.recipeStepInstrumentIndex] =
        !newState.stepHelpers[action.stepIndex].instrumentIsProduct[action.recipeStepInstrumentIndex];
      newState.recipe.steps[action.stepIndex].instruments[action.recipeStepInstrumentIndex] =
        new RecipeStepInstrumentCreationRequestInput();
      break;
    }

    case 'TOGGLE_VESSEL_PRODUCT_STATE': {
      newState.stepHelpers[action.stepIndex].vesselIsProduct[action.recipeStepVesselIndex] =
        !newState.stepHelpers[action.stepIndex].vesselIsProduct[action.recipeStepVesselIndex];
      newState.recipe.steps[action.stepIndex].vessels[action.recipeStepVesselIndex] =
        new RecipeStepVesselCreationRequestInput();
      break;
    }

    case 'TOGGLE_INGREDIENT_RANGE': {
      newState.stepHelpers[action.stepIndex].ingredientIsRanged[action.recipeStepIngredientIndex] =
        !newState.stepHelpers[action.stepIndex].ingredientIsRanged[action.recipeStepIngredientIndex];
      break;
    }

    case 'TOGGLE_INSTRUMENT_RANGE': {
      newState.stepHelpers[action.stepIndex].instrumentIsRanged[action.recipeStepInstrumentIndex] =
        !newState.stepHelpers[action.stepIndex].instrumentIsRanged[action.recipeStepInstrumentIndex];
      break;
    }

    case 'TOGGLE_VESSEL_RANGE': {
      newState.stepHelpers[action.stepIndex].vesselIsRanged[action.recipeStepVesselIndex] =
        !newState.stepHelpers[action.stepIndex].vesselIsRanged[action.recipeStepVesselIndex];
      break;
    }

    case 'TOGGLE_PRODUCT_RANGE': {
      newState.stepHelpers[action.stepIndex].productIsRanged[action.productIndex] =
        !newState.stepHelpers[action.stepIndex].productIsRanged[action.productIndex];
      break;
    }

    case 'TOGGLE_MANUAL_PRODUCT_NAMING': {
      newState.stepHelpers[action.stepIndex].productIsNamedManually[action.productIndex] =
        !newState.stepHelpers[action.stepIndex].productIsNamedManually[action.productIndex];

      const productIsNamedManually = state.stepHelpers[action.stepIndex].productIsNamedManually[action.productIndex];
      const productIsIngredient =
        state.recipe.steps[action.stepIndex].products[action.productIndex].type === 'ingredient';
      const productIsInstrument =
        state.recipe.steps[action.stepIndex].products[action.productIndex].type === 'instrument';

      let name = `product ${action.productIndex + 1} of step ${action.stepIndex + 1}`;
      if (!productIsNamedManually && productIsIngredient) {
        name = `${newState.stepHelpers[action.stepIndex].selectedPreparation?.pastTense} ${new Intl.ListFormat(
          'en',
        ).format(
          newState.recipe.steps[action.stepIndex].ingredients.map(
            (x: RecipeStepIngredientCreationRequestInput, i: number) =>
              newState.stepHelpers[action.stepIndex]?.selectedIngredients[i]?.name || x.name,
          ),
        )}`;
      } else if (!productIsNamedManually && productIsInstrument) {
        name = `${newState.stepHelpers[action.stepIndex].selectedPreparation?.pastTense} ${new Intl.ListFormat(
          'en',
        ).format(
          newState.recipe.steps[action.stepIndex].instruments.map(
            (x: RecipeStepInstrumentCreationRequestInput, i: number) =>
              newState.stepHelpers[action.stepIndex]?.selectedInstruments[i]?.name || x.name,
          ),
        )}`;
      }

      newState.recipe.steps[action.stepIndex].products[action.productIndex].name = name;
      break;
    }

    default:
      console.error(`Unhandled action type`);
  }

  return newState;
};
