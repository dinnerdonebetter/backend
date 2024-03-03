import { Card, List, Title, Text, Grid, Collapse, Checkbox, Group } from '@mantine/core';
import { ReactNode } from 'react';
import dagre from 'dagre';

import {
  Recipe,
  RecipeStep,
  RecipeStepIngredient,
  RecipeStepInstrument,
  RecipeStepVessel,
} from '@dinnerdonebetter/models';
import {
  buildRecipeStepText,
  cleanFloat,
  englishListFormatter,
  getRecipeStepIndexByProductID,
  recipeStepCanBePerformed,
  stepElementIsProduct,
} from '@dinnerdonebetter/utils';

import { TimerComponent } from '../components';
import { browserSideAnalytics } from '../../src/analytics';

const formatInstrumentList = (
  instruments: (RecipeStepInstrument | RecipeStepVessel)[],
  recipe: Recipe,
  stepIndex: number,
  recipeGraph: dagre.graphlib.Graph<string>,
  stepsNeedingCompletion: boolean[],
): ReactNode => {
  return (instruments || []).map((instrument: RecipeStepInstrument | RecipeStepVessel) => {
    const elementIsProduct = stepElementIsProduct(instrument);
    const checkboxDisabled = recipeStepCanBePerformed(stepIndex, recipeGraph, stepsNeedingCompletion);

    const dump = JSON.parse(JSON.stringify(instrument));
    const displayInSummaryLists = dump.hasOwnProperty('vessel')
      ? (instrument as RecipeStepVessel).vessel?.displayInSummaryLists
      : (instrument as RecipeStepInstrument).instrument?.displayInSummaryLists;

    return (
      (displayInSummaryLists || instrument.recipeStepProductID) && (
        <List.Item key={instrument.id}>
          {(elementIsProduct && (
            <>
              <Text size="sm" italic>
                {instrument.name}
              </Text>{' '}
              <Text size="sm">
                &nbsp;{` from step #${getRecipeStepIndexByProductID(recipe, instrument.recipeStepProductID!)}`}
              </Text>
            </>
          )) || <Checkbox size="sm" label={instrument.name} disabled={checkboxDisabled} />}
        </List.Item>
      )
    );
  });
};

declare interface RecipeStepComponentProps {
  recipe: Recipe;
  recipeStep: RecipeStep;
  recipeGraph: dagre.graphlib.Graph<string>;
  stepsNeedingCompletion: boolean[];
  stepIndex: number;
  stepCheckboxClicked: (_stepIndex: number) => void;
  scale?: number;
}

export const RecipeStepComponent = ({
  recipe,
  recipeStep,
  stepIndex,
  recipeGraph,
  stepsNeedingCompletion,
  stepCheckboxClicked,
  scale = 1.0,
}: RecipeStepComponentProps): JSX.Element => {
  const checkboxDisabled = recipeStepCanBePerformed(stepIndex, recipeGraph, stepsNeedingCompletion);

  const allInstruments = (recipeStep.instruments || []).filter(
    (instrument: RecipeStepInstrument) =>
      (instrument.instrument && instrument.instrument?.displayInSummaryLists) || instrument.recipeStepProductID,
  );
  const allVessels = (recipeStep.vessels || []).filter(
    (vessel: RecipeStepVessel) => (vessel.vessel && vessel.vessel?.displayInSummaryLists) || vessel.recipeStepProductID,
  );
  const allTools: (RecipeStepInstrument | RecipeStepVessel)[] = [...allInstruments, ...allVessels];

  const validIngredients = (recipeStep.ingredients || []).filter((ingredient) => ingredient.ingredient !== null);
  const productIngredients = (recipeStep.ingredients || []).filter(stepElementIsProduct);
  const allIngredients = validIngredients.concat(productIngredients);

  return (
    <>
      <Card key={recipeStep.id} shadow="sm" p="sm" radius="md" withBorder style={{ width: '100%', margin: '1rem' }}>
        <Card.Section px="sm">
          <Grid justify="space-between">
            <Grid.Col span="content">
              <Text italic={!stepsNeedingCompletion[stepIndex]}>
                {recipeStep.preparation.name}{' '}
                {englishListFormatter.format(
                  allIngredients.map((ingredient) =>
                    stepElementIsProduct(ingredient)
                      ? `${ingredient.name} from step #${getRecipeStepIndexByProductID(
                          recipe,
                          ingredient.recipeStepProductID!,
                        )}`
                      : ingredient.name,
                  ),
                )}
              </Text>
            </Grid.Col>
            <Grid.Col span="auto" />
            <Grid.Col span="content">
              <Group style={{ float: 'right' }}>
                <Checkbox
                  checked={!stepsNeedingCompletion[stepIndex]}
                  label={
                    checkboxDisabled ? 'Not Ready' : !stepsNeedingCompletion[stepIndex] ? 'Completed' : 'Not Completed'
                  }
                  labelPosition="left"
                  onChange={() => {}}
                  onClick={() => {
                    browserSideAnalytics.track('RECIPE_STEP_TOGGLED', {
                      recipeID: recipe.id,
                      recipeStepID: recipeStep.id,
                      checked: !stepsNeedingCompletion[stepIndex],
                    });

                    stepCheckboxClicked(stepIndex);
                  }}
                  disabled={checkboxDisabled}
                />
              </Group>
            </Grid.Col>
          </Grid>
        </Card.Section>

        <Collapse in={stepsNeedingCompletion[stepIndex]}>
          <Grid justify="center">
            <Grid.Col sm={12} md={8}>
              <Text strikethrough={!stepsNeedingCompletion[stepIndex]}>
                {buildRecipeStepText(recipe, recipeStep, scale)}
              </Text>

              <Text strikethrough={!recipeStepCanBePerformed(stepIndex, recipeGraph, stepsNeedingCompletion)} mt="md">
                {recipeStep.notes}
              </Text>

              {recipeStep.minimumEstimatedTimeInSeconds && (
                <TimerComponent durationInSeconds={recipeStep.minimumEstimatedTimeInSeconds} />
              )}
            </Grid.Col>

            <Grid.Col sm={12} md={4}>
              {allTools.length > 0 && (
                <Card.Section px="sm">
                  <Title order={6}>Tools:</Title>
                  <List icon={<></>} mt={-10}>
                    {formatInstrumentList(allTools, recipe, stepIndex, recipeGraph, stepsNeedingCompletion)}
                  </List>
                </Card.Section>
              )}

              {(recipeStep.ingredients || []).length > 0 && (
                <Card.Section px="sm" pt="sm">
                  <Title order={6}>Ingredients:</Title>
                  <List icon={<></>} mt={-10}>
                    {allIngredients.map((ingredient: RecipeStepIngredient): ReactNode => {
                      const checkboxDisabled = recipeStepCanBePerformed(stepIndex, recipeGraph, stepsNeedingCompletion);

                      const shouldDisplayMinQuantity = !stepElementIsProduct(ingredient);
                      const shouldDisplayMaxQuantity =
                        shouldDisplayMinQuantity &&
                        ingredient.maximumQuantity !== undefined &&
                        ingredient.maximumQuantity !== null &&
                        (ingredient.maximumQuantity ?? -1) > ingredient.minimumQuantity &&
                        ingredient.minimumQuantity != ingredient.maximumQuantity;
                      const elementIsProduct = stepElementIsProduct(ingredient);

                      let measurementName = shouldDisplayMinQuantity
                        ? cleanFloat(ingredient.minimumQuantity * scale) === 1
                          ? ingredient.measurementUnit.name
                          : ingredient.measurementUnit.pluralName
                        : '';
                      measurementName = ['unit', 'units'].includes(measurementName) ? '' : measurementName;

                      const ingredientName =
                        cleanFloat(ingredient.minimumQuantity * scale) === 1
                          ? ingredient.ingredient?.name || ingredient.name
                          : ingredient.ingredient?.pluralName || ingredient.name;

                      const lineText = (
                        <>
                          {`${shouldDisplayMinQuantity ? cleanFloat(ingredient.minimumQuantity * scale) : ''}${
                            shouldDisplayMaxQuantity ? `- ${cleanFloat((ingredient.maximumQuantity ?? 0) * scale)}` : ''
                          } ${measurementName}`}
                          {elementIsProduct ? <em>{ingredientName}</em> : <>{ingredientName}</>}
                          {`${
                            elementIsProduct
                              ? ` from step #${getRecipeStepIndexByProductID(recipe, ingredient.recipeStepProductID!)}`
                              : ''
                          }`}
                        </>
                      );

                      return (
                        <List.Item key={ingredient.id} mt="xs">
                          <Checkbox label={lineText} disabled={checkboxDisabled} mt="-sm" />
                        </List.Item>
                      );
                    })}
                  </List>
                </Card.Section>
              )}
            </Grid.Col>
          </Grid>
        </Collapse>
      </Card>
    </>
  );
};
