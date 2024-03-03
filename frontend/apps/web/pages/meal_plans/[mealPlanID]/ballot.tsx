import { AxiosError } from 'axios';
import { format } from 'date-fns';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Title, SimpleGrid, Grid, Center, Button, Card, Stack, ActionIcon, Indicator, Text } from '@mantine/core';
import Link from 'next/link';
import { Reducer, useEffect, useReducer, useState } from 'react';
import { IconArrowDown, IconArrowUp } from '@tabler/icons';
import router from 'next/router';

import {
  Household,
  HouseholdUserMembershipWithUser,
  MealPlan,
  MealPlanEvent,
  MealPlanOption,
  MealPlanOptionVote,
  MealPlanOptionVoteCreationRequestInput,
} from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildLocalClient, buildServerSideClient } from '../../../src/client';
import { AppLayout } from '../../../src/layouts';
import { serverSideTracer } from '../../../src/tracer';
import { serverSideAnalytics } from '../../../src/analytics';
import { extractUserInfoFromCookie } from '../../../src/auth';

declare interface MealPlanBallotPageProps {
  mealPlan: MealPlan;
  userID: string;
  household: Household;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<MealPlanBallotPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('MealPlanBallotPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const { mealPlanID: mealPlanIDParam } = context.query;
  if (!mealPlanIDParam) {
    throw new Error('meal plan ID is somehow missing!');
  }

  const mealPlanID = mealPlanIDParam.toString();

  const extractCookieTimer = timing.addEvent('extract cookie');
  const userSessionData = extractUserInfoFromCookie(context.req.cookies);
  if (userSessionData?.userID) {
    serverSideAnalytics.page(userSessionData.userID, 'MEAL_PLAN_BALLOT_PAGE', context, {
      mealPlanID,
      householdID: userSessionData.householdID,
    });
  } else {
    console.log(`no user session data found for ${context.req.url}`);
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

  let notFound = false;
  let notAuthorized = false;
  const retrievedData = await Promise.all([mealPlanPromise, householdPromise]).catch((error: AxiosError) => {
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
  const [mealPlan, household] = retrievedData;

  span.end();

  return {
    props: {
      mealPlan: mealPlan!,
      household: household!,
      userID: userSessionData?.userID || '',
    },
  };
};

const dateFormat = "h aa 'on' iiii',' M/d";

/* BEGIN Meal Plan Creation Reducer */

type mealPlanBallotPageAction =
  | { type: 'MOVE_OPTION'; eventIndex: number; optionIndex: number; direction: 'up' | 'down' }
  | { type: 'ADD_VOTES_TO_MEAL_PLAN'; eventIndex: number; votes: MealPlanOptionVote[] };

export class MealPlanBallotPageState {
  mealPlan: MealPlan = new MealPlan();

  constructor(mealPlan: MealPlan) {
    this.mealPlan = mealPlan;
  }
}

const useMealPlanReducer: Reducer<MealPlanBallotPageState, mealPlanBallotPageAction> = (
  state: MealPlanBallotPageState,
  action: mealPlanBallotPageAction,
): MealPlanBallotPageState => {
  switch (action.type) {
    case 'MOVE_OPTION':
      return {
        ...state,
        mealPlan: {
          ...state.mealPlan,
          events: (state.mealPlan.events || []).map((event: MealPlanEvent, eventIndex: number) => {
            if (
              (action.optionIndex === 0 && action.direction === 'up') ||
              (action.optionIndex === event.options.length - 1 && action.direction === 'down')
            ) {
              return event;
            }

            const options = [...event.options];
            [
              options[action.direction === 'up' ? action.optionIndex - 1 : action.optionIndex + 1],
              options[action.optionIndex],
            ] = [
              options[action.optionIndex],
              options[action.direction === 'up' ? action.optionIndex - 1 : action.optionIndex + 1],
            ];

            return eventIndex !== action.eventIndex
              ? event
              : {
                  ...event,
                  options: options,
                };
          }),
        },
      };

    case 'ADD_VOTES_TO_MEAL_PLAN':
      return {
        ...state,
        mealPlan: {
          ...state.mealPlan,
          events: (state.mealPlan.events || []).map((event: MealPlanEvent, eventIndex: number) => {
            return eventIndex !== action.eventIndex
              ? event
              : new MealPlanEvent({
                  ...event,
                  options: event.options.map((option: MealPlanOption) => {
                    const votes = (action.votes || []).filter(
                      (vote: MealPlanOptionVote) => vote.belongsToMealPlanOption === option.id,
                    );
                    return new MealPlanOption({
                      ...option,
                      votes: votes,
                    });
                  }),
                });
          }),
        },
      };

    default:
      console.error(`Unhandled action type`);
      return state;
  }
};

/* END Meal Plan Creation Reducer */

const getMissingVotersForMealPlanEvent = (
  mealPlanEvent: MealPlanEvent,
  household: Household,
  userID: string,
): Array<string> => {
  const missingVotes: Set<string> = new Set<string>();

  mealPlanEvent.options.forEach((option: MealPlanOption) => {
    household.members.forEach((member: HouseholdUserMembershipWithUser) => {
      if (
        (option.votes || []).find((vote: MealPlanOptionVote) => vote.byUser === member.belongsToUser!.id) === undefined
      ) {
        missingVotes.add(member.belongsToUser!.id !== userID ? member.belongsToUser!.username : 'you');
      }
    });
  });

  return Array.from(missingVotes.values());
};

const getUnvotedMealPlanEvents = (mealPlan: MealPlan, userID: string): Array<MealPlanEvent> => {
  return (mealPlan.events || []).filter((event: MealPlanEvent) => {
    return (
      event.options.find(
        (option: MealPlanOption) =>
          (option.votes || []).find((vote: MealPlanOptionVote) => vote.byUser === userID) === undefined,
      ) !== undefined
    );
  });
};

function MealPlanBallotPage({ mealPlan, userID, household }: MealPlanBallotPageProps) {
  const apiClient = buildLocalClient();
  const [pageState, dispatchPageEvent] = useReducer(useMealPlanReducer, new MealPlanBallotPageState(mealPlan));

  const [unvotedMealPlanEvents, setUnvotedMealPlanEvents] = useState<Array<MealPlanEvent>>([]);

  useEffect(() => {
    const x = getUnvotedMealPlanEvents(pageState.mealPlan, userID);
    setUnvotedMealPlanEvents(x);
    if (x.length === 0) {
      router.push(`/meal_plans/${mealPlan.id}`);
    }
  }, [pageState.mealPlan, mealPlan.id, userID]);

  const submitMealPlanVotes = (eventIndex: number): void => {
    const submission = new MealPlanOptionVoteCreationRequestInput({
      votes: pageState.mealPlan.events[eventIndex].options.map((option: MealPlanOption, rank: number) => {
        return {
          belongsToMealPlanOption: option.id,
          rank: rank,
          notes: '',
          abstain: false,
        };
      }),
    });

    apiClient
      .voteForMealPlan(mealPlan.id, pageState.mealPlan.events[eventIndex].id, submission)
      .then((votesResults: MealPlanOptionVote[]) => {
        dispatchPageEvent({
          type: 'ADD_VOTES_TO_MEAL_PLAN',
          eventIndex: eventIndex,
          votes: votesResults,
        });
      })
      .catch((error: Error) => {
        console.error(error);
      });
  };

  return (
    <AppLayout title="Meal Plan" containerSize="xl" userLoggedIn>
      <Center>
        <Stack>
          <Center>
            <Title order={3} p={5}>
              {`${format(
                new Date(
                  pageState.mealPlan.events.reduce((earliest: MealPlanEvent, event: MealPlanEvent) => {
                    return event.startsAt < earliest.startsAt ? event : earliest;
                  }).startsAt,
                ),
                dateFormat,
              )} - ${format(
                new Date(
                  pageState.mealPlan.events.reduce((earliest: MealPlanEvent, event: MealPlanEvent) => {
                    return event.startsAt > earliest.startsAt ? event : earliest;
                  }).startsAt,
                ),
                dateFormat,
              )}`}
            </Title>
          </Center>

          {unvotedMealPlanEvents.map((event: MealPlanEvent, eventIndex: number) => {
            return (
              <Card shadow="xs" radius="md" withBorder my="xl" key={eventIndex}>
                <Grid justify="center" align="center">
                  <Grid.Col span="auto">
                    <Text>
                      Rank choices for {event.mealName} at {format(new Date(event.startsAt), dateFormat)}
                    </Text>
                  </Grid.Col>
                  {mealPlan.status === 'awaiting_votes' && (
                    <Grid.Col span="content" sx={{ float: 'right' }}>
                      <Button onClick={() => submitMealPlanVotes(eventIndex)}>Submit Vote</Button>
                    </Grid.Col>
                  )}
                </Grid>
                {event.options.map((option: MealPlanOption, optionIndex: number) => {
                  return (
                    <Grid key={optionIndex}>
                      <Grid.Col span="auto">
                        <Indicator
                          position="top-start"
                          offset={2}
                          size={16}
                          disabled={optionIndex > 2}
                          color={
                            (optionIndex === 0 && 'yellow') ||
                            (optionIndex === 1 && 'gray') ||
                            (optionIndex === 2 && '#CD7F32') ||
                            'blue'
                          }
                          label={`#${optionIndex + 1}`}
                        >
                          <Card shadow="xs" radius="md" withBorder mt="xs">
                            <SimpleGrid>
                              <Link key={option.meal.id} href={`/meals/${option.meal.id}`}>
                                {option.meal.name}
                              </Link>
                            </SimpleGrid>
                          </Card>
                        </Indicator>
                      </Grid.Col>
                      {!event.options.find((opt: MealPlanOption) => opt.chosen) && (
                        <Grid.Col span="content">
                          <Stack align="center" spacing="xs" mt="sm">
                            <ActionIcon
                              variant="outline"
                              size="sm"
                              aria-label="remove recipe step vessel"
                              disabled={optionIndex === 0}
                              onClick={() => {
                                dispatchPageEvent({
                                  type: 'MOVE_OPTION',
                                  eventIndex: eventIndex,
                                  optionIndex: optionIndex,
                                  direction: 'up',
                                });
                              }}
                            >
                              <IconArrowUp size="md" />
                            </ActionIcon>
                            <ActionIcon
                              variant="outline"
                              size="sm"
                              aria-label="remove recipe step vessel"
                              disabled={optionIndex === event.options.length - 1}
                              onClick={() => {
                                dispatchPageEvent({
                                  type: 'MOVE_OPTION',
                                  eventIndex: eventIndex,
                                  optionIndex: optionIndex,
                                  direction: 'down',
                                });
                              }}
                            >
                              <IconArrowDown size="md" />
                            </ActionIcon>
                          </Stack>
                        </Grid.Col>
                      )}
                    </Grid>
                  );
                })}

                {getMissingVotersForMealPlanEvent(event, household, userID).length > 0 && (
                  <Grid justify="center" align="center">
                    <Grid.Col span="auto">
                      <sub>{`(awaiting votes from ${new Intl.ListFormat('en').format(
                        getMissingVotersForMealPlanEvent(event, household, userID),
                      )})`}</sub>
                    </Grid.Col>
                  </Grid>
                )}
              </Card>
            );
          })}
        </Stack>
      </Center>
    </AppLayout>
  );
}

export default MealPlanBallotPage;
