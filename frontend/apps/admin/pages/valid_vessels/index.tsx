import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Grid, Pagination, Stack, Table, TextInput } from '@mantine/core';
import { AxiosError } from 'axios';
import { formatRelative } from 'date-fns';
import router from 'next/router';
import { IconSearch } from '@tabler/icons';
import { useState, useEffect } from 'react';

import { QueryFilter, ValidVessel, QueryFilteredResult } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildLocalClient, buildServerSideClient } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';

declare interface ValidVesselsPageProps {
  pageLoadValidVessels: QueryFilteredResult<ValidVessel>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidVesselsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidVesselsPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  // TODO: parse context.query as QueryFilter.
  let props!: GetServerSidePropsResult<ValidVesselsPageProps>;

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const fetchValidVesselsTimer = timing.addEvent('fetch valid vessels');
  await apiClient
    .getValidVessels(qf)
    .then((res: QueryFilteredResult<ValidVessel>) => {
      span.addEvent('valid vessels retrieved');
      props = { props: { pageLoadValidVessels: res } };
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
      fetchValidVesselsTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function ValidVesselsPage(props: ValidVesselsPageProps) {
  let { pageLoadValidVessels } = props;

  const [validVessels, setValidVessels] = useState<QueryFilteredResult<ValidVessel>>(pageLoadValidVessels);
  const [search, setSearch] = useState('');

  useEffect(() => {
    const apiClient = buildLocalClient();

    if (search.trim().length < 1) {
      const qf = QueryFilter.deriveFromGetServerSidePropsContext({ search });
      apiClient
        .getValidVessels(qf)
        .then((res: QueryFilteredResult<ValidVessel>) => {
          setValidVessels(res);
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    } else {
      apiClient
        .searchForValidVessels(search)
        .then((res: ValidVessel[]) => {
          setValidVessels({
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
    qf.page = validVessels.page;

    apiClient
      .getValidVessels(qf)
      .then((res: QueryFilteredResult<ValidVessel>) => {
        setValidVessels(res);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [validVessels.page]);

  const formatDate = (x: string | undefined): string => {
    return x ? formatRelative(new Date(x), new Date()) : 'never';
  };

  const rows = (validVessels.data || []).map((vessel) => (
    <tr key={vessel.id} onClick={() => router.push(`/valid_vessels/${vessel.id}`)} style={{ cursor: 'pointer' }}>
      <td>{vessel.id}</td>
      <td>{vessel.name}</td>
      <td>{vessel.slug}</td>
      <td>{formatDate(vessel.createdAt)}</td>
      <td>{formatDate(vessel.lastUpdatedAt)}</td>
    </tr>
  ));

  return (
    <AppLayout title="Valid Vessels">
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
                router.push('/valid_vessels/new');
              }}
            >
              Create New
            </Button>
          </Grid.Col>
        </Grid>

        <Table mt="xl" striped highlightOnHover withBorder withColumnBorders>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
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
          page={validVessels.page}
          total={Math.ceil(validVessels.totalCount / validVessels.limit)}
          onChange={(value: number) => {
            setValidVessels({ ...validVessels, page: value });
          }}
        />
      </Stack>
    </AppLayout>
  );
}

export default ValidVesselsPage;
