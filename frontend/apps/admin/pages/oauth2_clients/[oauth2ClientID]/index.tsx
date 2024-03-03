import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { TextInput, Button, Group, Container, Divider } from '@mantine/core';
import { useState } from 'react';
import { useRouter } from 'next/router';

import { OAuth2Client } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { AppLayout } from '../../../src/layouts';
import { buildLocalClient, buildServerSideClient } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';

declare interface OAuth2ClientPageProps {
  pageLoadOAuth2Client: OAuth2Client;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<OAuth2ClientPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('OAuth2ClientPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const { oauth2ClientID } = context.query;
  if (!oauth2ClientID) {
    throw new Error('oauth2 client ID is somehow missing!');
  }

  const fetchOAuth2ClientTimer = timing.addEvent('fetch OAuth2 client');
  const pageLoadOAuth2ClientPromise = apiClient
    .getOAuth2Client(oauth2ClientID.toString())
    .then((result: OAuth2Client) => {
      span.addEvent('oauth2 client retrieved');
      return result;
    })
    .finally(() => {
      fetchOAuth2ClientTimer.end();
    });

  const [pageLoadOAuth2Client] = await Promise.all([pageLoadOAuth2ClientPromise]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: { pageLoadOAuth2Client },
  };
};

function OAuth2ClientPage(props: OAuth2ClientPageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const { pageLoadOAuth2Client } = props;

  const [oauth2Client, _setOAuth2Client] = useState<OAuth2Client>(pageLoadOAuth2Client);

  return (
    <AppLayout title="Valid Preparation">
      <Container size="sm">
        <TextInput label="Name" value={oauth2Client.name} />
        <TextInput label="Description" value={oauth2Client.description} />
        <TextInput label="Client ID" value={oauth2Client.clientID} />
        <TextInput label="Client Secret" value={oauth2Client.clientSecret} />

        <Divider my="xl" />

        <Group position="center">
          <Button
            type="submit"
            color="red"
            fullWidth
            onClick={() => {
              if (confirm('Are you sure you want to delete this OAuth2 client?')) {
                apiClient.deleteOAuth2Client(oauth2Client.id).then(() => {
                  router.push('/oauth_clients');
                });
              }
            }}
          >
            Delete
          </Button>
        </Group>
      </Container>
    </AppLayout>
  );
}

export default OAuth2ClientPage;
