import type { Page, Route } from '@playwright/test';

import { QueryFilteredResult, ValidInstrument, ValidInstrumentUpdateRequestInput } from '@dinnerdonebetter/models';
import { spellWord } from './utils';
import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockValidInstrumentResponseConfig extends ResponseConfig<ValidInstrument> {
  validInstrumentID: string;

  constructor(validInstrumentID: string, status: number = 200, body?: ValidInstrument) {
    super();

    this.validInstrumentID = validInstrumentID;
    this.status = status;
    this.body = body;
  }
}

export const mockValidInstrument = (resCfg: MockValidInstrumentResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_instruments/${resCfg.validInstrumentID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockValidInstrumentListResponseConfig extends ResponseConfig<QueryFilteredResult<ValidInstrument>> {}

export const mockValidInstrumentsList = (resCfg: MockValidInstrumentListResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_instruments?${resCfg.filter.asURLSearchParams().toString()}`,
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

export class MockValidInstrumentSearchResponseConfig extends ResponseConfig<ValidInstrument[]> {}

export const mockValidInstrumentsSearch = (resCfg: MockValidInstrumentSearchResponseConfig) => {
  return (page: Page) => {
    for (const word of spellWord(resCfg.query)) {
      page.route(`**/api/v1/valid_instruments/search?q=${encodeURIComponent(word)}`, (route: Route) => {
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

export class MockValidInstrumentUpdateResponseConfig extends ResponseConfig<ValidInstrumentUpdateRequestInput> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200, body?: ValidInstrumentUpdateRequestInput) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
    this.body = body;
  }
}

export const mockUpdateValidInstrument = (resCfg: MockValidInstrumentUpdateResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_instruments/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockValidInstrumentDeleteResponseConfig extends ResponseConfig<void> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
  }
}

export const mockDeleteValidInstrument = (resCfg: MockValidInstrumentDeleteResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_instruments/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
