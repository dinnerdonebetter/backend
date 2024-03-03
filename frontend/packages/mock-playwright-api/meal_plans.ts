import { Page, Route } from '@playwright/test';
import { MealPlan, QueryFilteredResult } from '@dinnerdonebetter/models';
import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockMealPlanResponseConfig extends ResponseConfig<MealPlan> {
  mealPlanID: string;

  constructor(mealPlanID: string, status: number = 200, body?: MealPlan) {
    super();

    this.mealPlanID = mealPlanID;
    this.status = status;
    this.body = body;
  }
}

export const mockMealPlan = (resCfg: MockMealPlanResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}`,
      (route: Route) => {
        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockMealPlanListResponseConfig extends ResponseConfig<QueryFilteredResult<MealPlan>> {}

export const mockMealPlans = (resCfg: MockMealPlanListResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans?${resCfg.filter.asURLSearchParams().toString()}`,
      (route: Route) => {
        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
