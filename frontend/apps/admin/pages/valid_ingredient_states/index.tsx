import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Grid, Pagination, Stack, Table, TextInput } from '@mantine/core';
import { AxiosError } from 'axios';
import { formatRelative } from 'date-fns';
import router from 'next/router';
import { IconSearch } from '@tabler/icons';
import { useState, useEffect } from 'react';

import { QueryFilter, QueryFilteredResult, ValidIngredientState } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildLocalClient, buildServerSideClient } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';

declare interface ValidIngredientStatesPageProps {
  pageLoadValidIngredientStates: QueryFilteredResult<ValidIngredientState>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidIngredientStatesPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidIngredientStatesPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  // TODO: parse context.query as QueryFilter.
  let props!: GetServerSidePropsResult<ValidIngredientStatesPageProps>;

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const fetchValidIngredientStatesTimer = timing.addEvent('fetch valid ingredient states');
  await apiClient
    .getValidIngredientStates(qf)
    .then((res: QueryFilteredResult<ValidIngredientState>) => {
      span.addEvent('valid ingredientStates retrieved');
      props = { props: { pageLoadValidIngredientStates: res } };
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
      fetchValidIngredientStatesTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function ValidIngredientStatesPage(props: ValidIngredientStatesPageProps) {
  let { pageLoadValidIngredientStates } = props;

  const [validIngredientStates, setValidIngredientStates] =
    useState<QueryFilteredResult<ValidIngredientState>>(pageLoadValidIngredientStates);
  const [search, setSearch] = useState('');

  useEffect(() => {
    const apiClient = buildLocalClient();

    if (search.trim().length < 1) {
      const qf = QueryFilter.deriveFromGetServerSidePropsContext({ search });
      apiClient
        .getValidIngredientStates(qf)
        .then((res: QueryFilteredResult<ValidIngredientState>) => {
          setValidIngredientStates(res);
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    } else {
      apiClient
        .searchForValidIngredientStates(search)
        .then((res: ValidIngredientState[]) => {
          setValidIngredientStates({
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
    qf.page = validIngredientStates.page;

    apiClient
      .getValidIngredientStates(qf)
      .then((res: QueryFilteredResult<ValidIngredientState>) => {
        setValidIngredientStates(res);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [validIngredientStates.page]);

  const formatDate = (x: string | undefined): string => {
    return x ? formatRelative(new Date(x), new Date()) : 'never';
  };

  const rows = (validIngredientStates.data || []).map((ingredientState) => (
    <tr
      key={ingredientState.id}
      onClick={() => router.push(`/valid_ingredient_states/${ingredientState.id}`)}
      style={{ cursor: 'pointer' }}
    >
      <td>{ingredientState.id}</td>
      <td>{ingredientState.name}</td>
      <td>{ingredientState.pastTense}</td>
      <td>{ingredientState.slug}</td>
      <td>{formatDate(ingredientState.createdAt)}</td>
      <td>{formatDate(ingredientState.lastUpdatedAt)}</td>
    </tr>
  ));

  return (
    <AppLayout title="Valid IngredientStates">
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
                router.push('/valid_ingredient_states/new');
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
          page={validIngredientStates.page}
          total={Math.ceil(validIngredientStates.totalCount / validIngredientStates.limit)}
          onChange={(value: number) => {
            setValidIngredientStates({ ...validIngredientStates, page: value });
          }}
        />
      </Stack>
    </AppLayout>
  );
}

export default ValidIngredientStatesPage;
