import React, { StrictMode } from 'react';

import { AppLayout } from '../src/layouts';

export default function Web(): JSX.Element {
  return (
    <StrictMode>
      <AppLayout title="" userLoggedIn={false}>
        <>{/* TODO: get a home page screen, lol */}</>
      </AppLayout>
    </StrictMode>
  );
}
