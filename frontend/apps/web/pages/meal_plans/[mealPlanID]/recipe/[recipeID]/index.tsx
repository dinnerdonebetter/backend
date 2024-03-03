import { AxiosError } from 'axios';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';

import { Household, MealPlan, MealPlanGroceryListItem, MealPlanTask, Recipe } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildServerSideClient } from '../../../../../src/client';
import { AppLayout } from '../../../../../src/layouts';
import { RecipeComponent } from '../../../../../src/components';
import { serverSideTracer } from '../../../../../src/tracer';
import { serverSideAnalytics } from '../../../../../src/analytics';
import { extractUserInfoFromCookie } from '../../../../../src/auth';

declare interface MealPlanRecipePageProps {
  recipe: Recipe;
  mealPlan: MealPlan;
  userID: string;
  household: Household;
  groceryList: MealPlanGroceryListItem[];
  tasks: MealPlanTask[];
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<MealPlanRecipePageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('RecipePage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const { mealPlanID: mealPlanIDParam, recipeID: recipeIDParam } = context.query;
  if (!mealPlanIDParam) {
    throw new Error('meal plan ID is somehow missing!');
  }
  const mealPlanID = mealPlanIDParam.toString();

  if (!recipeIDParam) {
    throw new Error('recipe ID is somehow missing!');
  }
  const recipeID = recipeIDParam.toString();

  const extractCookieTimer = timing.addEvent('extract cookie');
  const userSessionData = extractUserInfoFromCookie(context.req.cookies);
  if (userSessionData?.userID) {
    serverSideAnalytics.page(userSessionData.userID, 'RECIPE_PAGE', context, {
      recipeID,
      householdID: userSessionData.householdID,
    });
  }
  extractCookieTimer.end();

  const fetchMealPlanTimer = timing.addEvent('fetch meal plan');
  const mealPlanPromise = apiClient
    .getMealPlan(mealPlanID)
    .then((result: MealPlan) => {
      span.addEvent(`meal plan retrieved`);
      return result;
    })
    .finally(() => {
      fetchMealPlanTimer.end();
    });

  const fetchHouseholdTimer = timing.addEvent('fetch household');
  const householdPromise = apiClient
    .getCurrentHouseholdInfo()
    .then((result: Household) => {
      span.addEvent(`household retrieved`);
      return result;
    })
    .finally(() => {
      fetchHouseholdTimer.end();
    });

  const fetchMealPlanTasksTimer = timing.addEvent('fetch meal plan tasks');
  const tasksPromise = apiClient
    .getMealPlanTasks(mealPlanID)
    .then((result: MealPlanTask[]) => {
      span.addEvent('meal plan grocery list items retrieved');
      return result;
    })
    .finally(() => {
      fetchMealPlanTasksTimer.end();
    });

  const fetchMealPlanGroceryListItemsTimer = timing.addEvent('fetch meal plan grocery list items');
  const groceryListPromise = apiClient
    .getMealPlanGroceryListItems(mealPlanID)
    .then((result: MealPlanGroceryListItem[]) => {
      span.addEvent('meal plan grocery list items retrieved');
      return result;
    })
    .finally(() => {
      fetchMealPlanGroceryListItemsTimer.end();
    });

  const fetchRecipeTimer = timing.addEvent('fetch recipe');
  const recipePromise = apiClient
    .getRecipe(recipeID.toString())
    .then((result: Recipe) => {
      span.addEvent(`recipe retrieved`);
      return result;
    })
    .finally(() => {
      fetchRecipeTimer.end();
    });

  let notFound = false;
  let notAuthorized = false;
  const retrievedData = await Promise.all([
    mealPlanPromise,
    householdPromise,
    groceryListPromise,
    tasksPromise,
    recipePromise,
  ]).catch((error: AxiosError) => {
    if (error.response?.status === 404) {
      notFound = true;
    } else if (error.response?.status === 401) {
      notAuthorized = true;
    } else {
      console.error(`${error.response?.status} ${error.response?.config?.url}}`);
    }
  });

  if (notFound || !retrievedData) {
    return {
      redirect: {
        destination: '/meal_plans',
        permanent: false,
      },
    };
  }

  if (notAuthorized) {
    return {
      redirect: {
        destination: '/login',
        permanent: false,
      },
    };
  }

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  const [mealPlan, household, groceryList, tasks, recipe] = retrievedData;

  span.end();
  return {
    props: {
      recipe: recipe,
      mealPlan: mealPlan!,
      household: household!,
      userID: userSessionData?.userID || '',
      tasks: tasks,
      groceryList: groceryList || [],
    },
  };
};

export default function MealPlanRecipePage({ recipe }: MealPlanRecipePageProps) {
  return (
    <AppLayout title={recipe.name} userLoggedIn>
      <RecipeComponent recipe={recipe} />
    </AppLayout>
  );
}
