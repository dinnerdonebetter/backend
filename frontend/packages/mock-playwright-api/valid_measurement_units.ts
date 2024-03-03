import type { Page, Route } from '@playwright/test';

import {
  QueryFilteredResult,
  ValidMeasurementUnit,
  ValidMeasurementUnitUpdateRequestInput,
} from '@dinnerdonebetter/models';
import { spellWord } from './utils';
import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockValidMeasurementUnitResponseConfig extends ResponseConfig<ValidMeasurementUnit> {
  validMeasurementUnitID: string;

  constructor(validMeasurementUnitID: string, status: number = 200, body?: ValidMeasurementUnit) {
    super();

    this.validMeasurementUnitID = validMeasurementUnitID;
    this.status = status;
    this.body = body;
  }
}

export const mockValidMeasurementUnit = (resCfg: MockValidMeasurementUnitResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredients/${resCfg.validMeasurementUnitID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockValidMeasurementUnitListResponseConfig extends ResponseConfig<
  QueryFilteredResult<ValidMeasurementUnit>
> {}

export const mockValidMeasurementUnitsList = (resCfg: MockValidMeasurementUnitListResponseConfig) => {
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

export class MockValidMeasurementUnitSearchResponseConfig extends ResponseConfig<ValidMeasurementUnit[]> {}

export const mockValidMeasurementUnitsSearch = (resCfg: MockValidMeasurementUnitSearchResponseConfig) => {
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

export class MockValidMeasurementUnitUpdateResponseConfig extends ResponseConfig<ValidMeasurementUnitUpdateRequestInput> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200, body?: ValidMeasurementUnitUpdateRequestInput) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
    this.body = body;
  }
}

export const mockUpdateValidMeasurementUnit = (resCfg: MockValidMeasurementUnitUpdateResponseConfig) => {
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

export class MockValidMeasurementUnitDeleteResponseConfig extends ResponseConfig<void> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 200) {
    super();

    this.validPreparationID = validPreparationID;
    this.status = status;
  }
}

export const mockDeleteValidMeasurementUnit = (resCfg: MockValidMeasurementUnitDeleteResponseConfig) => {
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
