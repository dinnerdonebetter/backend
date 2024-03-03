import { AxiosError } from 'axios';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import Link from 'next/link';
import { Container, Divider, Grid, List, NumberInput, Space, Title } from '@mantine/core';
import { ReactNode, useState } from 'react';

import { ALL_MEAL_COMPONENT_TYPES, Meal, MealComponent } from '@dinnerdonebetter/models';
import { determineAllIngredientsForRecipes, determineAllInstrumentsForRecipes } from '@dinnerdonebetter/utils';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildServerSideClient } from '../../../src/client';
import { AppLayout } from '../../../src/layouts';
import { serverSideTracer } from '../../../src/tracer';
import { serverSideAnalytics } from '../../../src/analytics';
import { extractUserInfoFromCookie } from '../../../src/auth';
import { RecipeInstrumentListComponent } from '../../../src/components';
import { RecipeIngredientListComponent } from '../../../src/components/IngredientList';

declare interface MealPageProps {
  meal: Meal;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<MealPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('MealPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const { mealID } = context.query;
  if (!mealID) {
    throw new Error('meal ID is somehow missing!');
  }

  const extractCookieTimer = timing.addEvent('extract cookie');
  const userSessionData = extractUserInfoFromCookie(context.req.cookies);
  if (userSessionData?.userID) {
    serverSideAnalytics.page(userSessionData.userID, 'MEAL_PAGE', context, {
      mealID,
      householdID: userSessionData.householdID,
    });
  }
  extractCookieTimer.end();

  const fetchRecipeTimer = timing.addEvent('fetch meal');
  let props!: GetServerSidePropsResult<MealPageProps>;
  await apiClient
    .getMeal(mealID.toString())
    .then((result: Meal) => {
      span.addEvent(`recipe retrieved`);
      props = { props: { meal: result } };
    })
    .catch((error: AxiosError) => {
      if (error.response?.status === 404) {
        props = {
          redirect: {
            destination: '/meals',
            permanent: false,
          },
        };
      }
    })
    .finally(() => {
      fetchRecipeTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

// https://stackoverflow.com/a/14872766
const ordering: Record<string, number> = {};
for (let i = 0; i < ALL_MEAL_COMPONENT_TYPES.length; i++) {
  ordering[ALL_MEAL_COMPONENT_TYPES[i]] = i;
}

const formatRecipeList = (meal: Meal): ReactNode => {
  const sorted = (meal.components || []).sort(function (a: MealComponent, b: MealComponent) {
    return ordering[a.componentType] - ordering[b.componentType] || a.componentType.localeCompare(b.componentType);
  });

  return sorted.map((c: MealComponent, index: number) => {
    return (
      <List.Item key={index} mb="md">
        <Link href={`/recipes/${c.recipe.id}`}>{c.recipe.name}</Link>
        <em>{c.recipe.description}</em>
      </List.Item>
    );
  });
};

function MealPage({ meal }: MealPageProps) {
  const [mealScale, setMealScale] = useState(1.0);

  return (
    <AppLayout title={meal.name} userLoggedIn>
      <Container size="xs">
        <Title order={3}>{meal.name}</Title>

        <NumberInput
          mt="sm"
          mb="lg"
          value={mealScale}
          precision={2}
          step={0.25}
          removeTrailingZeros={true}
          description={`this meal normally yields about ${meal.minimumEstimatedPortions} ${
            meal.minimumEstimatedPortions === 1 ? 'portion' : 'portions'
          }${
            mealScale === 1.0
              ? ''
              : `, but is now set up to yield ${meal.minimumEstimatedPortions * mealScale}  ${
                  meal.minimumEstimatedPortions === 1 ? 'portion' : 'portions'
                }`
          }`}
          onChange={(value: number | undefined) => {
            if (!value) return;

            setMealScale(value);
          }}
        />

        <Divider label="recipes" labelPosition="center" mb="xs" />

        <Grid grow gutter="md">
          <List mt="sm">{formatRecipeList(meal)}</List>
        </Grid>

        <Divider label="resources" labelPosition="center" mb="md" />

        <Grid>
          <Grid.Col span={6}>
            <Title order={6}>Tools:</Title>
            {determineAllInstrumentsForRecipes((meal.components || []).map((x) => x.recipe)).length > 0 && (
              <RecipeInstrumentListComponent recipes={(meal.components || []).map((x) => x.recipe)} />
            )}
          </Grid.Col>

          <Grid.Col span={6}>
            <Title order={6}>Ingredients:</Title>
            {determineAllIngredientsForRecipes(
              (meal.components || []).map((x) => {
                return { scale: mealScale, recipe: x.recipe };
              }),
            ).length > 0 && (
              <RecipeIngredientListComponent recipes={(meal.components || []).map((x) => x.recipe)} scale={mealScale} />
            )}
          </Grid.Col>
        </Grid>
      </Container>
      <Space h="xl" my="xl" />
    </AppLayout>
  );
}

export default MealPage;
