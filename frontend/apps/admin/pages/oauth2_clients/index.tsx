import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Grid, Pagination, Stack, Table } from '@mantine/core';
import { AxiosError } from 'axios';
import { formatRelative } from 'date-fns';
import router from 'next/router';
import { useState, useEffect } from 'react';

import { QueryFilter, OAuth2Client, QueryFilteredResult } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildLocalClient, buildServerSideClient } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';

declare interface OAuth2ClientsPageProps {
  pageLoadOAuth2Clients: QueryFilteredResult<OAuth2Client>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<OAuth2ClientsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('OAuth2ClientsPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  // TODO: parse context.query as QueryFilter.
  let props!: GetServerSidePropsResult<OAuth2ClientsPageProps>;

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const fetchOAuth2ClientsTimer = timing.addEvent('fetch oauth2 clients');
  await apiClient
    .getOAuth2Clients(qf)
    .then((res: QueryFilteredResult<OAuth2Client>) => {
      span.addEvent('valid preparations retrieved');
      props = { props: { pageLoadOAuth2Clients: res } };
    })
    .catch((error: AxiosError) => {
      span.addEvent('error occurred');
      if (error.response?.status === 401) {
        props = {
          redirect: {
            destination: `/login?dest=${encodeURIComponent(context.resolvedUrl)}`,
            permanent: false,
          },
        };
      }
    })
    .finally(() => {
      fetchOAuth2ClientsTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function OAuth2ClientsPage(props: OAuth2ClientsPageProps) {
  let { pageLoadOAuth2Clients } = props;

  const [oauth2Clients, setOAuth2Clients] = useState<QueryFilteredResult<OAuth2Client>>(pageLoadOAuth2Clients);

  useEffect(() => {
    const apiClient = buildLocalClient();

    const qf = QueryFilter.deriveFromPage();
    qf.page = oauth2Clients.page;

    apiClient
      .getOAuth2Clients(qf)
      .then((res: QueryFilteredResult<OAuth2Client>) => {
        setOAuth2Clients(res || []);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [oauth2Clients.page]);

  const formatDate = (x: string | undefined): string => {
    return x ? formatRelative(new Date(x), new Date()) : 'never';
  };

  const rows = (oauth2Clients.data || []).map((preparation) => (
    <tr
      key={preparation.id}
      onClick={() => router.push(`/oauth2_clients/${preparation.id}`)}
      style={{ cursor: 'pointer' }}
    >
      <td>{preparation.name}</td>
      <td>{formatDate(preparation.createdAt)}</td>
    </tr>
  ));

  return (
    <AppLayout title="Valid Preparations">
      <Stack>
        <Grid justify="space-between">
          <Grid.Col md="content" sm={12}>
            <Button
              onClick={() => {
                router.push('/oauth2_clients/new');
              }}
            >
              Create New
            </Button>
          </Grid.Col>
        </Grid>

        <Table mt="xl" striped highlightOnHover withBorder withColumnBorders>
          <thead>
            <tr>
              <th>Name</th>
              <th>Created At</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </Table>

        <Pagination
          position="center"
          page={oauth2Clients.page}
          total={Math.ceil(oauth2Clients.totalCount / oauth2Clients.limit)}
          onChange={(value: number) => {
            setOAuth2Clients({ ...oauth2Clients, page: value });
          }}
        />
      </Stack>
    </AppLayout>
  );
}

export default OAuth2ClientsPage;
