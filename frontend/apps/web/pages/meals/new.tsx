import {
  ActionIcon,
  Alert,
  Autocomplete,
  AutocompleteItem,
  Button,
  Container,
  Divider,
  Text,
  Group,
  List,
  NumberInput,
  Select,
  SimpleGrid,
  Space,
  TextInput,
  Title,
  Stack,
} from '@mantine/core';
import { IconAlertCircle, IconX } from '@tabler/icons';
import { useRouter } from 'next/router';
import { AxiosError } from 'axios';
import { ReactNode, Reducer, useReducer, useEffect } from 'react';

import {
  ALL_MEAL_COMPONENT_TYPES,
  Meal,
  MealComponent,
  MealComponentType,
  Recipe,
  QueryFilteredResult,
} from '@dinnerdonebetter/models';
import { ConvertMealToMealCreationRequestInput } from '@dinnerdonebetter/utils';

import { buildLocalClient } from '../../src/client';
import { AppLayout } from '../../src/layouts';

/* BEGIN Meal Creation Reducer */

type mealCreationReducerAction =
  | { type: 'UPDATE_SUBMISSION_ERROR'; error: string }
  | { type: 'UPDATE_RECIPE_SUGGESTIONS'; recipeSuggestions: Recipe[] }
  | { type: 'UPDATE_RECIPE_QUERY'; newQuery: string }
  | { type: 'UPDATE_RECIPE_COMPONENT_TYPE'; componentIndex: number; componentType: MealComponentType }
  | { type: 'UPDATE_NAME'; newName: string }
  | { type: 'UPDATE_MINIMUM_PORTION_ESTIMATE'; newValue: number }
  | { type: 'UPDATE_MAXIMUM_PORTION_ESTIMATE'; newValue: number }
  | { type: 'UPDATE_DESCRIPTION'; newDescription: string }
  | { type: 'UPDATE_MEAL_COMPONENT_SCALE'; componentIndex: number; newScale: number }
  | { type: 'ADD_RECIPE'; recipe: Recipe }
  | { type: 'REMOVE_RECIPE'; recipe: Recipe };

export class MealCreationPageState {
  meal: Meal = new Meal({ minimumEstimatedPortions: 1, maximumEstimatedPortions: undefined });
  submissionShouldBePrevented: boolean = true;
  recipeQuery: string = '';
  mealScales: number[] = [1.0];
  recipeSuggestions: Recipe[] = [];
  submissionError: string | null = null;
}

const mealSubmissionShouldBeDisabled = (pageState: MealCreationPageState): boolean => {
  const componentProblems: string[] = [];

  pageState.meal.components.forEach((component: MealComponent, index: number) => {
    if (!component.componentType || component.componentType === 'unspecified') {
      componentProblems.push(`Component ${index + 1} is missing a component type`);
    }
  });

  return !(pageState.meal.name.length > 0 && pageState.meal.components.length > 0 && componentProblems.length === 0);
};

const useMealCreationReducer: Reducer<MealCreationPageState, mealCreationReducerAction> = (
  state: MealCreationPageState,
  action: mealCreationReducerAction,
): MealCreationPageState => {
  switch (action.type) {
    case 'UPDATE_SUBMISSION_ERROR':
      return { ...state, submissionError: action.error };

    case 'UPDATE_RECIPE_QUERY':
      return { ...state, recipeQuery: action.newQuery };

    case 'UPDATE_RECIPE_SUGGESTIONS':
      return { ...state, recipeSuggestions: action.recipeSuggestions };

    case 'UPDATE_RECIPE_COMPONENT_TYPE':
      let newComponents = [...state.meal.components];
      newComponents[action.componentIndex].componentType = action.componentType;

      return {
        ...state,
        meal: { ...state.meal, components: newComponents },
      };

    case 'UPDATE_NAME':
      return {
        ...state,
        meal: { ...state.meal, name: action.newName },
      };

    case 'UPDATE_MINIMUM_PORTION_ESTIMATE':
      return {
        ...state,
        meal: { ...state.meal, minimumEstimatedPortions: action.newValue },
      };

    case 'UPDATE_MAXIMUM_PORTION_ESTIMATE':
      return {
        ...state,
        meal: { ...state.meal, maximumEstimatedPortions: action.newValue },
      };

    case 'UPDATE_DESCRIPTION':
      return { ...state, meal: { ...state.meal, description: action.newDescription } };

    case 'UPDATE_MEAL_COMPONENT_SCALE':
      return {
        ...state,
        meal: {
          ...state.meal,
          components: state.meal.components.map((mc: MealComponent, index: number) => {
            if (index === action.componentIndex) {
              return { ...mc, recipeScale: action.newScale };
            }
            return mc;
          }),
        },
      };

    case 'ADD_RECIPE':
      const mealName = state.meal.name || action.recipe.name;

      return {
        ...state,
        recipeQuery: '',
        recipeSuggestions: [],
        meal: {
          ...state.meal,
          name: mealName,
          components: [...state.meal.components, new MealComponent({ recipe: action.recipe, recipeScale: 1 })],
        },
      };

    case 'REMOVE_RECIPE':
      return {
        ...state,
        meal: {
          ...state.meal,
          components: state.meal.components.filter((mc: MealComponent) => mc.recipe.id !== action.recipe.id),
        },
      };

    default:
      return state;
  }
};

/* END Meal Creation Reducer */

export default function NewMealPage(): JSX.Element {
  const router = useRouter();
  const [pageState, dispatchMealUpdate] = useReducer(useMealCreationReducer, new MealCreationPageState());

  useEffect(() => {
    const apiClient = buildLocalClient();
    const recipeQuery = pageState.recipeQuery.trim();
    if (recipeQuery.length > 2) {
      apiClient.searchForRecipes(recipeQuery).then((res: QueryFilteredResult<Recipe>) => {
        dispatchMealUpdate({ type: 'UPDATE_RECIPE_SUGGESTIONS', recipeSuggestions: res.data || [] });
      });
    } else {
      dispatchMealUpdate({ type: 'UPDATE_RECIPE_SUGGESTIONS', recipeSuggestions: [] as Recipe[] });
    }
  }, [pageState.recipeQuery]);

  const selectRecipe = (item: AutocompleteItem) => {
    const selectedRecipe = pageState.recipeSuggestions.find(
      (x: Recipe) =>
        x.name === item.value && !pageState.meal.components.find((y: MealComponent) => y.recipe.id === x.id),
    );

    if (selectedRecipe) {
      dispatchMealUpdate({ type: 'ADD_RECIPE', recipe: selectedRecipe });
    }
  };

  const removeRecipe = (recipe: Recipe) => {
    dispatchMealUpdate({ type: 'REMOVE_RECIPE', recipe: recipe });
  };

  const submitMeal = async () => {
    const apiClient = buildLocalClient();
    apiClient
      .createMeal(ConvertMealToMealCreationRequestInput(pageState.meal))
      .then((res: Meal) => {
        router.push(`/meals/${res.id}`);
      })
      .catch((err: AxiosError) => {
        console.error(`Failed to create meal: ${err}`);
        dispatchMealUpdate({ type: 'UPDATE_SUBMISSION_ERROR', error: err.message });
      });
  };

  let chosenRecipes: ReactNode = (pageState.meal.components || []).map(
    (mealComponent: MealComponent, componentIndex: number) => (
      <List.Item key={mealComponent.recipe.id} icon={<></>} pt="xs">
        <SimpleGrid cols={4}>
          <div>
            <Select
              mt="-md"
              label="Type"
              placeholder="main, side, etc."
              value={mealComponent.componentType}
              onChange={(value: MealComponentType) =>
                dispatchMealUpdate({
                  type: 'UPDATE_RECIPE_COMPONENT_TYPE',
                  componentIndex: componentIndex,
                  componentType: value,
                })
              }
              data={ALL_MEAL_COMPONENT_TYPES.filter((x) => x != 'unspecified').map((x) => ({ label: x, value: x }))}
            />
          </div>

          <div>
            <Text mt="xs"> {mealComponent.recipe.name}</Text>
          </div>

          <div>
            <NumberInput
              precision={2}
              mt="-sm"
              step={0.25}
              descriptionProps={{ fontSize: 'sm' }}
              description={`This recipe will yield ${
                mealComponent.recipeScale * mealComponent.recipe.minimumEstimatedPortions
              }${
                mealComponent.recipe.maximumEstimatedPortions
                  ? `- ${mealComponent.recipeScale * mealComponent.recipe.maximumEstimatedPortions}`
                  : ''
              } ${
                mealComponent.recipeScale * mealComponent.recipe.minimumEstimatedPortions == 1
                  ? mealComponent.recipe.portionName
                  : mealComponent.recipe.pluralPortionName
              }`}
              value={mealComponent.recipeScale}
              onChange={(value: number) => {
                if (value <= 0) {
                  return;
                }

                dispatchMealUpdate({
                  type: 'UPDATE_MEAL_COMPONENT_SCALE',
                  componentIndex,
                  newScale: value,
                });
              }}
            />
          </div>

          <div>
            <ActionIcon
              mt="xs"
              onClick={() => removeRecipe(mealComponent.recipe)}
              sx={{ float: 'right' }}
              aria-label="remove recipe from meal"
            >
              <IconX color="tomato" />
            </ActionIcon>
          </div>
        </SimpleGrid>
      </List.Item>
    ),
  );

  return (
    <AppLayout title="New Meal" userLoggedIn>
      <Container size="md">
        <Title order={3}>Create Meal</Title>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            submitMeal();
          }}
        >
          <Stack>
            <TextInput
              withAsterisk
              label="Name"
              value={pageState.meal.name}
              onChange={(event) => dispatchMealUpdate({ type: 'UPDATE_NAME', newName: event.target.value })}
              mt="xs"
            />

            <TextInput
              label="Description"
              value={pageState.meal.description}
              onChange={(event) =>
                dispatchMealUpdate({ type: 'UPDATE_DESCRIPTION', newDescription: event.target.value })
              }
              mt="xs"
            />

            <Group position="center" spacing="xl" grow>
              <NumberInput
                label="Min. Portions"
                value={pageState.meal.minimumEstimatedPortions}
                onChange={(value: number) => {
                  if (value <= 0) {
                    return;
                  }

                  dispatchMealUpdate({
                    type: 'UPDATE_MINIMUM_PORTION_ESTIMATE',
                    newValue: value,
                  });
                }}
                mt="xs"
              />

              <NumberInput
                label="Max Portions"
                value={pageState.meal.maximumEstimatedPortions}
                onChange={(value: number) => {
                  if (value <= 0) {
                    return;
                  }

                  dispatchMealUpdate({
                    type: 'UPDATE_MAXIMUM_PORTION_ESTIMATE',
                    newValue: value,
                  });
                }}
                mt="xs"
              />
            </Group>
          </Stack>

          <Space h="lg" />
          <Divider />

          <Autocomplete
            value={pageState.recipeQuery}
            onChange={(value: string) => dispatchMealUpdate({ type: 'UPDATE_RECIPE_QUERY', newQuery: value })}
            limit={20}
            label="Recipe name"
            placeholder="baba ganoush"
            onItemSubmit={selectRecipe}
            data={pageState.recipeSuggestions.map((x: Recipe) => ({ value: x.name, label: x.name }))}
            mt="xs"
          />
          <Space h="md" />

          <List>{chosenRecipes}</List>

          <Space h="md" />
          {pageState.submissionError && (
            <Alert m="md" icon={<IconAlertCircle size={16} />} color="tomato">
              {pageState.submissionError}
            </Alert>
          )}

          <Group position="center">
            <Button type="submit" disabled={mealSubmissionShouldBeDisabled(pageState)}>
              Submit
            </Button>
          </Group>
        </form>
      </Container>
    </AppLayout>
  );
}
