import { useState } from 'react';
import axios, { AxiosError, AxiosResponse } from 'axios';
import { useRouter } from 'next/router';
import { useForm, zodResolver } from '@mantine/form';
import { Alert, TextInput, PasswordInput, Button, Group, Space, Container } from '@mantine/core';
import { z } from 'zod';

import { IAPIError, UserLoginInput, UserStatusResponse } from '@dinnerdonebetter/models';

import { AppLayout } from '../src/layouts';

const loginFormSchema = z.object({
  username: z.string().trim().min(1, 'username is required'),
  password: z.string().trim().min(8, 'password must have at least 8 characters'),
  totpToken: z.string().trim().regex(/\d{6}/, 'token must be 6 digits'),
});

export default function Login(): JSX.Element {
  const router = useRouter();

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
    const validation = loginForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const loginInput = new UserLoginInput({
      username: loginForm.values.username.trim(),
      password: loginForm.values.password.trim(),
      totpToken: loginForm.values.totpToken.trim(),
    });

    await axios
      .post('/api/login', loginInput)
      .then((_res: AxiosResponse<UserStatusResponse>) => {
        const redirect = decodeURIComponent(new URLSearchParams(window.location.search).get('dest') || '').trim();

        router.push(redirect || '/');
      })
      .catch((err: AxiosError<IAPIError>) => {
        setLoginError(err?.response?.data.message || 'unknown error occurred');
      });
  };

  return (
    <AppLayout title="Login">
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
          <TextInput
            data-qa="totp-input"
            label="TOTP Token"
            placeholder="123456"
            {...loginForm.getInputProps('totpToken')}
          />

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
                loginForm.values.totpToken.length === 0
              }
            >
              Login
            </Button>
          </Group>
        </form>

        <Space h="xl" mb="xl" />
        <Space h="xl" mb="xl" />
      </Container>
    </AppLayout>
  );
}
