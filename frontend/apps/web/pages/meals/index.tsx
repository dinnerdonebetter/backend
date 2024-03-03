import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { Button, Center, Container, List } from '@mantine/core';

import { Meal, QueryFilteredResult, QueryFilter } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { serverSideTracer } from '../../src/tracer';
import { buildServerSideClient } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideAnalytics } from '../../src/analytics';
import { extractUserInfoFromCookie } from '../../src/auth';

declare interface MealsPageProps {
  meals: Meal[];
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<MealsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('MealsPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const extractCookieTimer = timing.addEvent('extract cookie');
  const userSessionData = extractUserInfoFromCookie(context.req.cookies);
  if (userSessionData?.userID) {
    serverSideAnalytics.page(userSessionData.userID, 'MEALS_PAGE', context, {
      query: context.query,
      householdID: userSessionData.householdID,
    });
  }
  extractCookieTimer.end();

  const fetchMealsTimer = timing.addEvent('fetch meals');
  let props!: GetServerSidePropsResult<MealsPageProps>;
  await apiClient
    .getMeals(qf)
    .then((result: QueryFilteredResult<Meal>) => {
      span.addEvent('meals retrieved');
      props = { props: { meals: result.data } };
    })
    .catch((error) => {
      span.addEvent('error retrieving meals');
      if (error.response?.status === 401) {
        props = {
          redirect: {
            destination: `/login?dest=${encodeURIComponent(context.resolvedUrl)}`,
            permanent: false,
          },
        };
        return;
      }

      throw error;
    })
    .finally(() => {
      fetchMealsTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function MealsPage(props: MealsPageProps) {
  const router = useRouter();
  const { meals } = props;

  const mealItems = (meals || []).map((meal: Meal) => (
    <List.Item key={meal.id}>
      <Link href={`/meals/${meal.id}`}>{meal.name}</Link>
    </List.Item>
  ));

  return (
    <AppLayout title="Meals" userLoggedIn>
      <Container size="xs">
        <Center>
          <Button
            my="lg"
            onClick={() => {
              router.push('/meals/new');
            }}
          >
            New Meal
          </Button>
        </Center>
        <List>{mealItems}</List>
      </Container>
    </AppLayout>
  );
}

export default MealsPage;
