import { EmailAddressVerificationRequestInput } from '@dinnerdonebetter/models';
import { AxiosError } from 'axios';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';

import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildCookielessServerSideClient } from '../src/client';
import { serverSideTracer } from '../src/tracer';

declare interface VerifyEmailAddressPageProps {}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<VerifyEmailAddressPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('RegistrationPage.getServerSideProps');
  const apiClient = buildCookielessServerSideClient().withSpan(span);

  const emailAddressVerificationTimer = timing.addEvent('verify email address');
  const emailVerificationToken = context.query['t']?.toString() || '';
  await apiClient
    .verifyEmailAddress(new EmailAddressVerificationRequestInput({ emailVerificationToken }))
    .then(() => {
      span.addEvent('email address verified');
    })
    .catch((err: AxiosError) => {
      span.addEvent('email address verification failed');
      span.setStatus({
        code: err.response?.status || 500,
        message: err.message,
      });
    })
    .finally(() => {
      emailAddressVerificationTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();

  return {
    redirect: {
      destination: `/`,
      permanent: false,
    },
  };
};

export default function VerifyEmailAddressPage(_props: VerifyEmailAddressPageProps): JSX.Element {
  return <></>;
}
