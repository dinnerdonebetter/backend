import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Grid, Pagination, Stack, Table, TextInput } from '@mantine/core';
import { AxiosError } from 'axios';
import { formatRelative } from 'date-fns';
import router from 'next/router';
import { IconSearch } from '@tabler/icons';
import { useState, useEffect } from 'react';

import { QueryFilter, ValidPreparation, QueryFilteredResult } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildLocalClient, buildServerSideClient } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';

declare interface ValidPreparationsPageProps {
  pageLoadValidPreparations: QueryFilteredResult<ValidPreparation>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidPreparationsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidPreparationsPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  // TODO: parse context.query as QueryFilter.
  let props!: GetServerSidePropsResult<ValidPreparationsPageProps>;

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const fetchValidPreparationsTimer = timing.addEvent('fetch valid preparations');
  await apiClient
    .getValidPreparations(qf)
    .then((res: QueryFilteredResult<ValidPreparation>) => {
      span.addEvent('valid preparations retrieved');
      props = { props: { pageLoadValidPreparations: res } };
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
      fetchValidPreparationsTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function ValidPreparationsPage(props: ValidPreparationsPageProps) {
  let { pageLoadValidPreparations } = props;

  const [validPreparations, setValidPreparations] =
    useState<QueryFilteredResult<ValidPreparation>>(pageLoadValidPreparations);
  const [search, setSearch] = useState('');

  useEffect(() => {
    const apiClient = buildLocalClient();

    if (search.trim().length < 1) {
      const qf = QueryFilter.deriveFromGetServerSidePropsContext({ search });
      apiClient
        .getValidPreparations(qf)
        .then((res: QueryFilteredResult<ValidPreparation>) => {
          setValidPreparations(res);
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    } else {
      apiClient
        .searchForValidPreparations(search)
        .then((res: ValidPreparation[]) => {
          setValidPreparations({
            ...QueryFilter.Default(),
            data: res || [],
            filteredCount: (res || []).length,
            totalCount: (res || []).length,
          });
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    }
  }, [search]);

  useEffect(() => {
    const apiClient = buildLocalClient();

    const qf = QueryFilter.deriveFromPage();
    qf.page = validPreparations.page;

    apiClient
      .getValidPreparations(qf)
      .then((res: QueryFilteredResult<ValidPreparation>) => {
        setValidPreparations(res);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [validPreparations.page]);

  const formatDate = (x: string | undefined): string => {
    return x ? formatRelative(new Date(x), new Date()) : 'never';
  };

  const rows = (validPreparations.data || []).map((preparation) => (
    <tr
      key={preparation.id}
      onClick={() => router.push(`/valid_preparations/${preparation.id}`)}
      style={{ cursor: 'pointer' }}
    >
      <td>{preparation.name}</td>
      <td>{preparation.pastTense}</td>
      <td>{preparation.slug}</td>
      <td>{formatDate(preparation.createdAt)}</td>
      <td>{formatDate(preparation.lastUpdatedAt)}</td>
    </tr>
  ));

  return (
    <AppLayout title="Valid Preparations">
      <Stack>
        <Grid justify="space-between">
          <Grid.Col md="auto" sm={12}>
            <TextInput
              placeholder="Search..."
              icon={<IconSearch size={14} />}
              onChange={(event) => setSearch(event.target.value || '')}
            />
          </Grid.Col>
          <Grid.Col md="content" sm={12}>
            <Button
              onClick={() => {
                router.push('/valid_preparations/new');
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
              <th>Past Tense</th>
              <th>Slug</th>
              <th>Created At</th>
              <th>Last Updated At</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </Table>

        <Pagination
          disabled={search.trim().length > 0}
          position="center"
          page={validPreparations.page}
          total={Math.ceil(validPreparations.totalCount / validPreparations.limit)}
          onChange={(value: number) => {
            setValidPreparations({ ...validPreparations, page: value });
          }}
        />
      </Stack>
    </AppLayout>
  );
}

export default ValidPreparationsPage;
