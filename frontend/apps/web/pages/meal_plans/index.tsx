import { format } from 'date-fns';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Center, Container, Table } from '@mantine/core';
import { useRouter } from 'next/router';
import { useState } from 'react';

import { MealPlan, QueryFilter } from '@dinnerdonebetter/models';
import { getEarliestEvent, getLatestEvent } from '@dinnerdonebetter/utils';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildServerSideClient } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';
import { serverSideAnalytics } from '../../src/analytics';
import { extractUserInfoFromCookie } from '../../src/auth';

declare interface MealPlansPageProps {
  userID: string;
  mealPlans: MealPlan[];
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<MealPlansPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('MealPlansPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const extractCookieTimer = timing.addEvent('extract cookie');
  const userSessionData = extractUserInfoFromCookie(context.req.cookies);
  if (userSessionData?.userID) {
    serverSideAnalytics.page(userSessionData.userID, 'MEAL_PLANS_PAGE', context, {
      householdID: userSessionData.householdID,
    });
  } else {
    return {
      redirect: {
        destination: '/login',
        permanent: false,
      },
    };
  }
  extractCookieTimer.end();

  const fetchMealPlansTimer = timing.addEvent('fetch meal plans');
  const { data: mealPlans } = await apiClient
    .getMealPlans(qf)
    .then((result) => {
      span.addEvent('meal plan list retrieved');
      return result;
    })
    .finally(() => {
      fetchMealPlansTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return { props: { userID: userSessionData?.userID, mealPlans: mealPlans } };
};

const dateFormat = 'h aa M/d/yy';

function MealPlansPage(props: MealPlansPageProps) {
  const router = useRouter();

  const { mealPlans: pageLoadMealPlans } = props;
  const [mealPlans, updateMealPlans] = useState(pageLoadMealPlans);

  return (
    <AppLayout title="Meal Plans" userLoggedIn>
      <Container size="xs">
        {mealPlans.length > 0 && (
          <Table>
            <thead>
              <tr>
                <th
                  onClick={() =>
                    updateMealPlans(mealPlans.sort((a, b) => (b.status > a.status ? 1 : a.status > b.status ? -1 : 0)))
                  }
                >
                  Status
                </th>
                <th
                  onClick={() =>
                    updateMealPlans(
                      mealPlans.sort((a, b) =>
                        new Date(b.votingDeadline) > new Date(a.votingDeadline)
                          ? 1
                          : new Date(a.votingDeadline) > new Date(b.votingDeadline)
                          ? -1
                          : 0,
                      ),
                    )
                  }
                >
                  Voting Deadline
                </th>
                <th
                  onClick={() =>
                    updateMealPlans(
                      mealPlans.sort((a, b) =>
                        b.events.length > a.events.length ? 1 : a.events.length > b.events.length ? -1 : 0,
                      ),
                    )
                  }
                >
                  Events
                </th>
                <th
                  onClick={() =>
                    updateMealPlans(
                      mealPlans.sort((a, b) =>
                        new Date(getEarliestEvent(b).startsAt) > new Date(getEarliestEvent(a).startsAt)
                          ? 1
                          : new Date(getEarliestEvent(a).startsAt) > new Date(getEarliestEvent(b).startsAt)
                          ? -1
                          : 0,
                      ),
                    )
                  }
                >
                  Starts At
                </th>
                <th
                  onClick={() =>
                    updateMealPlans(
                      mealPlans.sort((a, b) =>
                        new Date(getEarliestEvent(b).endsAt) > new Date(getEarliestEvent(a).endsAt)
                          ? 1
                          : new Date(getEarliestEvent(a).endsAt) > new Date(getEarliestEvent(b).endsAt)
                          ? -1
                          : 0,
                      ),
                    )
                  }
                >
                  Ends At
                </th>
              </tr>
            </thead>
            <tbody>
              {mealPlans.map((mealPlan: MealPlan, mealPlanIndex: number) => {
                return (
                  <tr
                    style={{ cursor: 'pointer' }}
                    key={mealPlanIndex}
                    onClick={() => router.push(`/meal_plans/${mealPlan.id}`)}
                  >
                    <td>{mealPlan.status}</td>
                    <td>{format(new Date(mealPlan.votingDeadline), dateFormat)}</td>
                    <td>{mealPlan.events.length}</td>
                    <td>{format(new Date(getEarliestEvent(mealPlan).startsAt), dateFormat)}</td>
                    <td>{format(new Date(getLatestEvent(mealPlan).endsAt), dateFormat)}</td>
                  </tr>
                );
              })}
            </tbody>
          </Table>
        )}

        <Center>
          <Button my="lg" onClick={() => router.push('/meal_plans/new')}>
            New Meal Plan
          </Button>
        </Center>
      </Container>
    </AppLayout>
  );
}

export default MealPlansPage;
