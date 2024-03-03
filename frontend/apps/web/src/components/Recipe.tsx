import { Card, List, Title, Grid, ActionIcon, Collapse, NumberInput } from '@mantine/core';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { IconCaretDown, IconCaretUp, IconRotate } from '@tabler/icons';

import { Recipe } from '@dinnerdonebetter/models';
import { renderMermaidDiagramForRecipe, toDAG } from '@dinnerdonebetter/utils';

import {
  Mermaid,
  RecipeStepComponent,
  RecipeIngredientListComponent,
  RecipeInstrumentListComponent,
} from '../components';

declare interface RecipeComponentProps {
  recipe: Recipe;
  scale?: number;
}

export const RecipeComponent = ({ recipe, scale = 1.0 }: RecipeComponentProps): JSX.Element => {
  let recipeGraph = toDAG(recipe);

  const [stepsNeedingCompletion, setStepsNeedingCompletion] = useState(
    Array((recipe.steps || []).length).fill(true) as boolean[],
  );
  const [flowChartVisible, setFlowChartVisibility] = useState(true);
  const [allIngredientListVisible, setIngredientListVisibility] = useState(false);
  const [allInstrumentListVisible, setInstrumentListVisibility] = useState(false);

  const [recipeScale, setRecipeScale] = useState(scale);

  const [graphDirection, setGraphDirection] = useState<'TB' | 'LR' | 'BT' | 'RL'>('TB');
  const [recipeGraphDiagram, setRecipeGraphDiagram] = useState(renderMermaidDiagramForRecipe(recipe, graphDirection));

  useEffect(() => {
    setRecipeGraphDiagram(renderMermaidDiagramForRecipe(recipe, graphDirection));
  }, [recipe, graphDirection]);

  return (
    <>
      <Grid grow justify="space-between" mx="xs">
        <Grid.Col span="content">
          <Title order={3} mr={'-xs'}>
            {recipe.name}
          </Title>
        </Grid.Col>
        <Grid.Col span="auto">
          {recipe.source && (
            <Link style={{ float: 'right' }} href={recipe.source}>
              (source)
            </Link>
          )}
        </Grid.Col>
      </Grid>

      <Grid grow>
        <Grid.Col>
          {/* flow chart */}
          <Card shadow="sm" p="sm" radius="md" withBorder sx={{ width: '100%', margin: '1rem' }}>
            <Card.Section px="xs" sx={{ cursor: 'pointer' }}>
              <Grid justify="space-between" align="center">
                <Grid.Col span="content">
                  <Title order={5} sx={{ display: 'inline-block' }} mt="xs">
                    Flow Chart
                  </Title>
                  <ActionIcon
                    sx={{ float: 'right' }}
                    pt="sm"
                    variant="transparent"
                    aria-label="rotate recipe flow chart orientation"
                    onClick={() => setGraphDirection(graphDirection === 'TB' ? 'LR' : 'TB')}
                  >
                    <IconRotate size={15} color="green" />
                  </ActionIcon>
                </Grid.Col>
                <Grid.Col span="auto" onClick={() => setFlowChartVisibility((x: boolean) => !x)}>
                  <ActionIcon sx={{ float: 'right' }} aria-label="toggle recipe flow chart">
                    {flowChartVisible && <IconCaretUp />}
                    {!flowChartVisible && <IconCaretDown />}
                  </ActionIcon>
                </Grid.Col>
              </Grid>
            </Card.Section>

            {flowChartVisible && <Mermaid chartDefinition={recipeGraphDiagram} />}
          </Card>
          <Grid>
            <Grid.Col span={6}>
              {/* Ingredients */}
              <Card shadow="sm" radius="md" withBorder style={{ width: '100%', margin: '1rem' }}>
                <Card.Section px="xs" sx={{ cursor: 'pointer' }}>
                  <Grid justify="space-between" align="center">
                    <Grid.Col span="content">
                      <Title order={5} style={{ display: 'inline-block' }} mt="xs">
                        All Ingredients
                      </Title>
                    </Grid.Col>
                    <Grid.Col span="auto" onClick={() => setIngredientListVisibility((x: boolean) => !x)}>
                      <ActionIcon sx={{ float: 'right' }} aria-label="toggle recipe flow chart">
                        {allIngredientListVisible && <IconCaretUp />}
                        {!allIngredientListVisible && <IconCaretDown />}
                      </ActionIcon>
                    </Grid.Col>
                  </Grid>
                </Card.Section>

                <Collapse in={allIngredientListVisible}>
                  <List icon={<></>} spacing={-15}>
                    <RecipeIngredientListComponent recipes={[recipe]} scale={recipeScale} />
                  </List>
                </Collapse>
              </Card>
            </Grid.Col>
            <Grid.Col span={6}>
              {/* Instruments */}
              <Card shadow="sm" radius="md" withBorder style={{ width: '100%', margin: '1rem' }}>
                <Card.Section px="xs" sx={{ cursor: 'pointer' }}>
                  <Grid justify="space-between" align="center">
                    <Grid.Col span="content">
                      <Title order={5} style={{ display: 'inline-block' }} mt="xs">
                        All Instruments
                      </Title>
                    </Grid.Col>
                    <Grid.Col span="auto" onClick={() => setInstrumentListVisibility((x: boolean) => !x)}>
                      <ActionIcon sx={{ float: 'right' }} aria-label="toggle recipe flow chart">
                        {allInstrumentListVisible && <IconCaretUp />}
                        {!allInstrumentListVisible && <IconCaretDown />}
                      </ActionIcon>
                    </Grid.Col>
                  </Grid>
                </Card.Section>

                <Collapse in={allInstrumentListVisible}>
                  <List icon={<></>} spacing={-15}>
                    <RecipeInstrumentListComponent recipes={[recipe]} />
                  </List>
                </Collapse>
              </Card>
            </Grid.Col>
          </Grid>

          {/* Scale */}
          <Card shadow="sm" radius="md" withBorder style={{ width: '100%', margin: '1rem' }}>
            <Card.Section px="xs" sx={{ cursor: 'pointer' }}>
              <Grid justify="space-between" align="center">
                <Grid.Col span="content">
                  <Title order={5} style={{ display: 'inline-block' }} mt="xs">
                    Scale
                  </Title>
                </Grid.Col>
              </Grid>

              <NumberInput
                mt="sm"
                mb="lg"
                value={recipeScale}
                precision={2}
                step={0.25}
                removeTrailingZeros={true}
                description={`this recipe normally yields about ${recipe.minimumEstimatedPortions} ${
                  recipe.minimumEstimatedPortions === 1 ? recipe.portionName : recipe.pluralPortionName
                }${
                  recipeScale === 1.0
                    ? ''
                    : `, but is now set up to yield ${recipe.minimumEstimatedPortions * recipeScale}  ${
                        recipe.minimumEstimatedPortions === 1 ? recipe.portionName : recipe.pluralPortionName
                      }`
                }`}
                onChange={(value: number | undefined) => {
                  if (!value) return;

                  setRecipeScale(value);
                }}
              />
            </Card.Section>
          </Card>

          {(recipe.steps || []).map((recipeStep, stepIndex) => {
            return (
              <RecipeStepComponent
                recipe={recipe}
                recipeStep={recipeStep}
                stepIndex={stepIndex}
                recipeGraph={recipeGraph}
                stepsNeedingCompletion={stepsNeedingCompletion}
                stepCheckboxClicked={(index: number) => {
                  setStepsNeedingCompletion(
                    stepsNeedingCompletion.map((x: boolean, i: number) => {
                      return i === index ? !x : x;
                    }),
                  );
                }}
                scale={scale}
              />
            );
          })}

          {/* Steps */}
          {/* {recipeSteps} */}
        </Grid.Col>
      </Grid>
    </>
  );
};
