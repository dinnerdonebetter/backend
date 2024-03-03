import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Container, Grid, List, Space, TextInput } from '@mantine/core';
import { AxiosError } from 'axios';
import Link from 'next/link';
import { IconSearch } from '@tabler/icons';

import { QueryFilter, Recipe, QueryFilteredResult } from '@dinnerdonebetter/models';
import { buildServerSideLogger } from '@dinnerdonebetter/logger';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildServerSideClient } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';
import { extractUserInfoFromCookie } from '../../src/auth';
import { serverSideAnalytics } from '../../src/analytics';

declare interface RecipesPageProps {
  recipes: Recipe[];
}

const logger = buildServerSideLogger('RecipesPage');

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<RecipesPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('RecipesPage.getServerSideProps');
  const spanContext = span.spanContext();
  const spanLogDetails = { spanID: spanContext.spanId, traceID: spanContext.traceId };
  const apiClient = buildServerSideClient(context).withSpan(span);

  const extractCookieTimer = timing.addEvent('extract cookie');
  const userSessionData = extractUserInfoFromCookie(context.req.cookies);
  if (userSessionData?.userID) {
    serverSideAnalytics.page(userSessionData.userID, 'RECIPES_PAGE', context, {
      householdID: userSessionData.householdID,
    });
  }
  extractCookieTimer.end();

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const fetchRecipesTimer = timing.addEvent('fetch recipes');
  let props!: GetServerSidePropsResult<RecipesPageProps>;
  await apiClient
    .getRecipes(qf)
    .then((res: QueryFilteredResult<Recipe>) => {
      span.addEvent('recipes retrieved');
      const recipes = res.data;
      props = { props: { recipes } };
    })
    .catch((error: AxiosError) => {
      span.addEvent('error occurred');
      if (error.response?.status === 401) {
        logger.error('unauthorized access to recipes page', spanLogDetails);
        props = {
          redirect: {
            destination: `/login?dest=${encodeURIComponent(context.resolvedUrl)}`,
            permanent: false,
          },
        };
      }
    })
    .finally(() => {
      fetchRecipesTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function RecipesPage(props: RecipesPageProps) {
  const { recipes } = props;

  const recipeItems = (recipes || []).map((recipe: Recipe) => (
    <List.Item key={recipe.id}>
      <Link href={`/recipes/${recipe.id}`}>{recipe.name}</Link>
    </List.Item>
  ));

  return (
    <AppLayout title="Recipes" userLoggedIn>
      <Container size="xs">
        <Grid justify="space-between">
          <Grid.Col md="auto" sm={12}>
            <TextInput
              placeholder="Search..."
              disabled
              icon={<IconSearch size={14} />}
              // onChange={(event) => setSearch(event.target.value || '')}
            />
          </Grid.Col>
        </Grid>

        <Space my="xl" />

        <List>{recipeItems}</List>
      </Container>
    </AppLayout>
  );
}

export default RecipesPage;
