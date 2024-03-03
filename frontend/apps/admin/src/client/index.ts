import { GetServerSidePropsContext } from 'next';
import router from 'next/router';

import DinnerDoneBetterAPIClient from '@dinnerdonebetter/api-client';

import { apiCookieName } from '../constants';

export const buildServerSideClient = (context: GetServerSidePropsContext): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  const ddbClient = new DinnerDoneBetterAPIClient(apiEndpoint, context.req.cookies[apiCookieName]);

  return ddbClient;
};

export const buildServerSideClientWithRawCookie = (cookie: string): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  if (!cookie) {
    throw new Error('no cookie set!');
  }

  const ddbClient = new DinnerDoneBetterAPIClient(apiEndpoint, cookie);

  return ddbClient;
};

export const buildCookielessServerSideClient = (): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  const ddbClient = new DinnerDoneBetterAPIClient(apiEndpoint);

  return ddbClient;
};

export const buildBrowserSideClient = (): DinnerDoneBetterAPIClient => {
  const ddbClient = buildCookielessServerSideClient();

  ddbClient.configureRouterRejectionInterceptor((loc: Location) => {
    const destParam = new URLSearchParams(loc.search).get('dest') ?? encodeURIComponent(`${loc.pathname}${loc.search}`);
    router.push({ pathname: '/login', query: { dest: destParam } });
  });

  return ddbClient;
};

export const buildLocalClient = (): DinnerDoneBetterAPIClient => {
  return new DinnerDoneBetterAPIClient();
};
