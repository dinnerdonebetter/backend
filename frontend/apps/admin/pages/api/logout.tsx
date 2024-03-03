import { AxiosError, AxiosResponse } from 'axios';
import type { NextApiRequest, NextApiResponse } from 'next';

import { buildServerSideLogger } from '@dinnerdonebetter/logger';

import { buildServerSideClientWithRawCookie } from '../../src/client';
import { apiCookieName } from '../../src/constants';
import { processWebappCookieHeader } from '../../src/auth';
import { serverSideTracer } from '../../src/tracer';

const logger = buildServerSideLogger('logout_route');

async function LogoutRoute(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === 'POST') {
    const span = serverSideTracer.startSpan('LogoutRoute');
    const spanContext = span.spanContext();
    const spanLogDetails = { spanID: spanContext.spanId, traceID: spanContext.traceId };

    const cookie = (req.headers['cookie'] || '').replace(`${apiCookieName}=`, '');
    if (!cookie) {
      logger.debug('cookie missing from logout request', spanLogDetails);
      res.status(401).send('no cookie attached');
      return;
    }

    logger.info('logging user out', spanLogDetails);

    const apiClient = buildServerSideClientWithRawCookie(cookie).withSpan(span);
    await apiClient
      .logOut()
      .then((result: AxiosResponse) => {
        span.addEvent('response received');
        const responseCookie = processWebappCookieHeader(result, '', '');
        res.setHeader('Set-Cookie', responseCookie).status(result.status).send('logged out');
        return;
      })
      .catch((err: AxiosError) => {
        span.addEvent('error received');
        logger.debug('error response received from logout', { status: err.response?.status, ...spanLogDetails });
        res.status(207).send('error logging out');
        return;
      });

    span.end();
  } else {
    res.status(405).send(`Method ${req.method} Not Allowed`);
  }
}

export default LogoutRoute;
