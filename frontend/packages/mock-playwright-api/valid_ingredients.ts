import type { Page, Route } from '@playwright/test';

import { QueryFilteredResult, ValidIngredient, ValidIngredientUpdateRequestInput } from '@dinnerdonebetter/models';
import { spellWord } from './utils';
import { assertClient, assertMethod, methods, ResponseConfig } from './helpers';

export class MockValidIngredientResponseConfig extends ResponseConfig<ValidIngredient> {
  validIngredientID: string;

  constructor(validIngredientID: string, status: number = 200, body?: ValidIngredient) {
    super();

    this.validIngredientID = validIngredientID;
    this.status = status;
    this.body = body;
  }
}

export const mockValidIngredient = (resCfg: MockValidIngredientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredients/${resCfg.validIngredientID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockValidIngredientListResponseConfig extends ResponseConfig<QueryFilteredResult<ValidIngredient>> {}

export const mockValidIngredientsList = (resCfg: MockValidIngredientListResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredients?${resCfg.filter.asURLSearchParams().toString()}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        if (resCfg.body && resCfg.filter) resCfg.body.limit = resCfg.filter.limit;
        if (resCfg.body && resCfg.filter) resCfg.body.page = resCfg.filter.page;

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockValidIngredientSearchResponseConfig extends ResponseConfig<ValidIngredient[]> {}

export const mockValidIngredientsSearch = (resCfg: MockValidIngredientSearchResponseConfig) => {
  return (page: Page) => {
    for (const word of spellWord(resCfg.query)) {
      page.route(`**/api/v1/valid_ingredients/search?q=${encodeURIComponent(word)}`, (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        const rv = resCfg.fulfill();

        if (word !== resCfg.query) {
          rv.body = JSON.stringify([]);
        }

        route.fulfill(rv);
      });
    }
  };
};

export class MockValidIngredientUpdateResponseConfig extends ResponseConfig<ValidIngredientUpdateRequestInput> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200, body?: ValidIngredientUpdateRequestInput) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
    this.body = body;
  }
}

export const mockUpdateValidIngredient = (resCfg: MockValidIngredientUpdateResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredients/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockValidIngredientDeleteResponseConfig extends ResponseConfig<void> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
  }
}

export const mockDeleteValidIngredient = (resCfg: MockValidIngredientDeleteResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredients/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
