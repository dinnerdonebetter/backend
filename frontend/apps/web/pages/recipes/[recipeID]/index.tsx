import { AxiosError } from 'axios';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';

import { Recipe } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildServerSideClient } from '../../../src/client';
import { AppLayout } from '../../../src/layouts';
import { RecipeComponent } from '../../../src/components';
import { serverSideTracer } from '../../../src/tracer';
import { serverSideAnalytics } from '../../../src/analytics';
import { extractUserInfoFromCookie } from '../../../src/auth';

declare interface RecipePageProps {
  recipe: Recipe;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<RecipePageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('RecipePage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const { recipeID } = context.query;
  if (!recipeID) {
    throw new Error('recipe ID is somehow missing!');
  }

  const extractCookieTimer = timing.addEvent('extract cookie');
  const userSessionData = extractUserInfoFromCookie(context.req.cookies);
  if (userSessionData?.userID) {
    serverSideAnalytics.page(userSessionData.userID, 'RECIPE_PAGE', context, {
      recipeID,
      householdID: userSessionData.householdID,
    });
  }
  extractCookieTimer.end();

  const fetchRecipeTimer = timing.addEvent('fetch recipe');
  let props!: GetServerSidePropsResult<RecipePageProps>;
  await apiClient
    .getRecipe(recipeID.toString())
    .then((result: Recipe) => {
      span.addEvent(`recipe retrieved`);
      props = { props: { recipe: result } };
    })
    .catch((error: AxiosError) => {
      if (error.response?.status === 404) {
        props = {
          redirect: {
            destination: '/recipes',
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

function RecipePage({ recipe }: RecipePageProps) {
  return (
    <AppLayout title={recipe.name} titlePosition="left" userLoggedIn>
      <RecipeComponent recipe={recipe} />
    </AppLayout>
  );
}

export default RecipePage;
