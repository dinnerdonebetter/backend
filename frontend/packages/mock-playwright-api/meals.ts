import { Page, Route } from '@playwright/test';
import { Meal, QueryFilteredResult } from '@dinnerdonebetter/models';
import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockMealResponseConfig extends ResponseConfig<Meal> {
  mealID: string;

  constructor(mealID: string, status: number = 200, body?: Meal) {
    super();

    this.mealID = mealID;
    this.status = status;
    this.body = body;
  }
}

export const mockMeal = (resCfg: MockMealResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meals/${resCfg.mealID}`,
      (route: Route) => {
        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockMealListResponseConfig extends ResponseConfig<QueryFilteredResult<Meal>> {}

export const mockMeals = (resCfg: MockMealListResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meals?${resCfg.filter.asURLSearchParams().toString()}`,
      (route: Route) => {
        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
