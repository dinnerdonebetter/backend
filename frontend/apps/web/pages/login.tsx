import { useState } from 'react';
import axios, { AxiosError, AxiosResponse } from 'axios';
import { useRouter } from 'next/router';
import { useForm, zodResolver } from '@mantine/form';
import { Alert, TextInput, PasswordInput, Button, Group, Space, Grid, Text, Container } from '@mantine/core';
import { z } from 'zod';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import Link from 'next/link';

import { IAPIError, UserLoginInput, UserStatusResponse } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { AppLayout } from '../src/layouts';
import { serverSideAnalytics } from '../src/analytics';
import { extractUserInfoFromCookie } from '../src/auth';
import { serverSideTracer } from '../src/tracer';

const loginFormSchema = z.object({
  username: z.string().trim().min(1, 'username is required'),
  password: z.string().trim().min(8, 'password must have at least 8 characters'),
  totpToken: z.string().trim().optional().or(z.string().trim().regex(/\d{6}/, 'token must be 6 digits')),
});

declare interface LoginPageProps {}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<LoginPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('LoginPageProps.getServerSideProps');

  const extractCookieTimer = timing.addEvent('extract cookie');
  const userSessionData = extractUserInfoFromCookie(context.req.cookies);
  if (userSessionData?.userID) {
    serverSideAnalytics.page(userSessionData.userID, 'LOGIN_PAGE', context, {
      householdID: userSessionData.householdID,
    });

    span.end();
    return {
      redirect: {
        destination: '/',
        permanent: false,
      },
    };
  }
  extractCookieTimer.end();

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return { props: {} };
};

export default function Login(_props: LoginPageProps): JSX.Element {
  const router = useRouter();

  const [needsTOTPToken, setNeedsTOTPToken] = useState(false);
  const [loginError, setLoginError] = useState('');

  const loginForm = useForm({
    initialValues: {
      username: '',
      password: '',
      totpToken: '',
    },
    validate: zodResolver(loginFormSchema),
  });

  const login = async () => {
    if (needsTOTPToken && !loginForm.values.totpToken) {
      loginForm.setFieldError('totpToken', 'TOTP token is required');
      return;
    }

    const validation = loginForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const loginInput = new UserLoginInput({
      username: (loginForm.values.username || '').trim(),
      password: (loginForm.values.password || '').trim(),
      totpToken: (loginForm.values.totpToken || '').trim(),
    });

    await axios
      .post('/api/login', loginInput)
      .then((result: AxiosResponse<UserStatusResponse>) => {
        if (result.status === 205) {
          setNeedsTOTPToken(true);
          return;
        }

        const redirect = decodeURIComponent(new URLSearchParams(window.location.search || '').get('dest') || '').trim();

        router.push(redirect || '/');
      })
      .catch((err: AxiosError<IAPIError>) => {
        setLoginError(err?.response?.data.message || 'unknown error occurred');
      });
  };

  return (
    <AppLayout title="Login" userLoggedIn={false}>
      {' '}
      {/* TODO: this is actually unknown, not false */}
      <Container size="xs">
        <form onSubmit={loginForm.onSubmit(login)}>
          <TextInput
            data-qa="username-input"
            label="Username"
            placeholder="username"
            {...loginForm.getInputProps('username')}
          />
          <PasswordInput
            data-qa="password-input"
            label="Password"
            placeholder="hunter2"
            {...loginForm.getInputProps('password')}
          />
          {needsTOTPToken && (
            <TextInput
              data-qa="totp-input"
              mt="md"
              label="TOTP Token"
              placeholder="123456"
              {...loginForm.getInputProps('totpToken')}
            />
          )}

          {loginError && (
            <>
              <Space h="md" />
              <Alert title="Oh no!" color="tomato">
                {loginError}
              </Alert>
            </>
          )}

          <Group position="center">
            <Button
              data-qa="submit"
              type="submit"
              mt="sm"
              fullWidth
              disabled={
                loginForm.values.username.length === 0 ||
                loginForm.values.password.length === 0 ||
                (needsTOTPToken && loginForm.values.totpToken.length === 0)
              }
            >
              Login
            </Button>
          </Group>

          <Grid justify="space-between" mt={2}>
            <Grid.Col span={3}>
              <Text size="xs" align="left">
                <Link href="/passwords/forgotten">Forgot password?</Link>
              </Text>
            </Grid.Col>
            <Grid.Col span={3}>
              <Text size="xs" align="right">
                <Link href="/register">Register instead</Link>
              </Text>
            </Grid.Col>
          </Grid>
        </form>
      </Container>
    </AppLayout>
  );
}
