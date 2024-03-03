import type { Page, Route } from '@playwright/test';

import { QueryFilteredResult, Recipe, RecipeUpdateRequestInput } from '@dinnerdonebetter/models';
import { spellWord } from './utils';
import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockRecipeResponseConfig extends ResponseConfig<Recipe> {
  recipeID: string;

  constructor(recipeID: string, status: number = 200, body?: Recipe) {
    super();

    this.recipeID = recipeID;
    this.status = status;
    this.body = body;
  }
}

export const mockRecipe = (resCfg: MockRecipeResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockRecipeListResponseConfig extends ResponseConfig<QueryFilteredResult<Recipe>> {}

export const mockRecipesList = (resCfg: MockRecipeListResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes?${resCfg.filter.asURLSearchParams().toString()}`,
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

export const mockRecipesSearch = (resCfg: MockRecipeListResponseConfig) => {
  return (page: Page) => {
    for (const word of spellWord(resCfg.query)) {
      page.route(`**/api/v1/recipes/search?q=${encodeURIComponent(word)}`, (route: Route) => {
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

export class MockRecipeUpdateResponseConfig extends ResponseConfig<RecipeUpdateRequestInput> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200, body?: RecipeUpdateRequestInput) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
    this.body = body;
  }
}

export const mockUpdateRecipe = (resCfg: MockRecipeUpdateResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockRecipeDeleteResponseConfig extends ResponseConfig<void> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
  }
}

export const mockDeleteRecipe = (resCfg: MockRecipeDeleteResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
