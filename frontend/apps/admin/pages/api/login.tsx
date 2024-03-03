import { AxiosError, AxiosResponse } from 'axios';
import type { NextApiRequest, NextApiResponse } from 'next';

import { IAPIError, UserLoginInput } from '@dinnerdonebetter/models';

import { buildCookielessServerSideClient } from '../../src/client';
import { serverSideTracer } from '../../src/tracer';
import { processWebappCookieHeader } from '../../src/auth';

async function LoginRoute(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === 'POST') {
    const span = serverSideTracer.startSpan('LoginRoute');
    const input = req.body as UserLoginInput;

    const apiClient = buildCookielessServerSideClient().withSpan(span);

    await apiClient
      .adminLogin(input)
      .then((result: AxiosResponse) => {
        span.addEvent('response received');
        if (result.status === 205) {
          res.status(205).send('');
          return;
        }

        res.setHeader('Set-Cookie', processWebappCookieHeader(result, result.data.userID, result.data.activeHousehold));

        res.status(202).send('');
      })
      .catch((err: AxiosError<IAPIError>) => {
        span.addEvent('error received');
        console.log('error from login route', err);
        res.status(err.response?.status || 500).send('');
        return;
      });

    span.end();
  } else {
    res.status(405).send(`Method ${req.method} Not Allowed`);
  }
}

export default LoginRoute;
