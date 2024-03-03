import { StrictMode } from 'react';
import { Container, Title } from '@mantine/core';

import { AppLayout } from '../src/layouts';

export default function Web(): JSX.Element {
  return (
    <StrictMode>
      <AppLayout title="">
        <Container size="xs">
          <Title order={5}>Obligatory home page</Title>
        </Container>
      </AppLayout>
    </StrictMode>
  );
}
