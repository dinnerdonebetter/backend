import type { Page, Route } from '@playwright/test';

import { QueryFilteredResult, ValidPreparation, ValidPreparationUpdateRequestInput } from '@dinnerdonebetter/models';
import { spellWord } from './utils';
import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockValidPreparationResponseConfig extends ResponseConfig<ValidPreparation> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200, body?: ValidPreparation) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
    this.body = body;
  }
}

export const mockValidPreparation = (resCfg: MockValidPreparationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparations/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockValidPreparationListResponseConfig extends ResponseConfig<QueryFilteredResult<ValidPreparation>> {}

export const mockValidPreparationsList = (resCfg: MockValidPreparationListResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparations?${resCfg.filter.asURLSearchParams().toString()}`,
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

export class MockValidPreparationSearchResponseConfig extends ResponseConfig<ValidPreparation[]> {}

export const mockValidPreparationsSearch = (resCfg: MockValidPreparationSearchResponseConfig) => {
  return (page: Page) => {
    for (const word of spellWord(resCfg.query)) {
      page.route(`**/api/v1/valid_preparations/search?q=${encodeURIComponent(word)}`, (route: Route) => {
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

export class MockValidPreparationUpdateResponseConfig extends ResponseConfig<ValidPreparationUpdateRequestInput> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200, body?: ValidPreparationUpdateRequestInput) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
    this.body = body;
  }
}

export const mockUpdateValidPreparation = (resCfg: MockValidPreparationUpdateResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparations/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockValidPreparationDeleteResponseConfig extends ResponseConfig<void> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
  }
}

export const mockDeleteValidPreparation = (resCfg: MockValidPreparationDeleteResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparations/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
