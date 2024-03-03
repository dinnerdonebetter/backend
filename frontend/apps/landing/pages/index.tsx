import { StrictMode } from 'react';
import { Button, Center, Container } from '@mantine/core';
import { useRouter } from 'next/router';

import { AppLayout } from '../src/layouts';

export default function Web(): JSX.Element {
  const router = useRouter();

  return (
    <StrictMode>
      <AppLayout title="The best dang lil' cookin' app on the internet">
        <Container size="xl">
          <Center>
            <Button onClick={() => router.push('https://app.dinnerdonebetter.dev/')}>Try the app</Button>
          </Center>
        </Container>
      </AppLayout>
    </StrictMode>
  );
}
