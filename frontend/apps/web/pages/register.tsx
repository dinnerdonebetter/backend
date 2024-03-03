import { useState } from 'react';
import { AxiosError } from 'axios';
import { useRouter } from 'next/router';
import { Alert, TextInput, PasswordInput, Button, Group, Space, Grid, Text, Container, Divider } from '@mantine/core';
import { DatePicker } from '@mantine/dates';
import { useForm, zodResolver } from '@mantine/form';
import { z } from 'zod';
import Link from 'next/link';
import { formatISO, subYears } from 'date-fns';

import { IAPIError, UserRegistrationInput } from '@dinnerdonebetter/models';

import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { serverSideTracer } from '../src/tracer';
import { buildBrowserSideClient } from '../src/client';
import { AppLayout } from '../src/layouts';

const registrationFormSchema = z.object({
  emailAddress: z.string().trim().email({ message: 'invalid email' }),
  householdName: z.string().trim(),
  username: z.string().trim().min(1, 'username is required'),
  firstName: z.string().trim().optional(),
  lastName: z.string().trim().optional(),
  password: z.string().trim().min(8, 'password must have at least 8 characters'),
  repeatedPassword: z.string().trim().min(8, 'repeated password must have at least 8 characters'),
  birthday: z.date().nullable(),
});

declare interface RegistrationPageProps {
  invitationToken?: string;
  invitationID?: string;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<RegistrationPageProps>> => {
  const span = serverSideTracer.startSpan('RegistrationPage.getServerSideProps');

  if (process.env.DISABLE_REGISTRATION === 'true') {
    span.end();
    return {
      redirect: {
        destination: '/login',
        permanent: false,
      },
    };
  }

  let props: GetServerSidePropsResult<RegistrationPageProps> = {
    props: {
      invitationID: context.query['i']?.toString() || '',
      invitationToken: context.query['t']?.toString() || '',
    },
  };

  span.end();
  return props;
};

export default function Register(props: RegistrationPageProps): JSX.Element {
  const router = useRouter();
  const [registrationError, setRegistrationError] = useState('');

  const { invitationID, invitationToken } = props;

  const registrationForm = useForm({
    initialValues: {
      emailAddress: '',
      username: '',
      firstName: '',
      lastName: '',
      password: '',
      householdName: '',
      repeatedPassword: '',
      birthday: null,
    },

    validate: zodResolver(registrationFormSchema),
  });

  const register = async () => {
    const validation = registrationForm.validate();
    if (validation.hasErrors) {
      return;
    }

    if (registrationForm.values.password !== registrationForm.values.repeatedPassword) {
      registrationForm.setFieldError('password', 'passwords do not match');
      registrationForm.setFieldError('repeatedPassword', 'passwords do not match');
      return;
    }

    const registrationInput = new UserRegistrationInput({
      emailAddress: registrationForm.values.emailAddress.trim(),
      firstName: registrationForm.values.firstName.trim(),
      lastName: registrationForm.values.lastName.trim(),
      username: registrationForm.values.username.trim(),
      password: registrationForm.values.password.trim(),
      householdName: registrationForm.values.householdName.trim(),
      invitationToken: (invitationToken || '').trim(),
      invitationID: (invitationID || '').trim(),
    });

    if (registrationForm.values.birthday) {
      registrationInput.birthday = formatISO(registrationForm.values.birthday);
    }

    await buildBrowserSideClient()
      .register(registrationInput)
      .then(() => {
        router.push('/login');
      })
      .catch((err: AxiosError<IAPIError>) => {
        setRegistrationError(err?.response?.data.message || 'unknown error occurred');
      });
  };

  return (
    <AppLayout title="register" userLoggedIn={false}>
      {' '}
      {/* TODO: this is actually unknown, not false */}
      <Container size="xs">
        <form onSubmit={registrationForm.onSubmit(register)}>
          <TextInput
            data-qa="registration-email-address-input"
            label="Email Address"
            required
            placeholder="cool_person@emailprovider.website"
            {...registrationForm.getInputProps('emailAddress')}
          />
          <TextInput
            data-qa="registration-username-input"
            label="Username"
            required
            placeholder="username"
            {...registrationForm.getInputProps('username')}
          />
          <PasswordInput
            data-qa="registration-password-input"
            label="Password"
            required
            placeholder="hunter2"
            {...registrationForm.getInputProps('password')}
          />
          <PasswordInput
            data-qa="registration-password-confirm-input"
            label="Password (again)"
            placeholder="hunter2"
            required
            {...registrationForm.getInputProps('repeatedPassword')}
          />

          <Divider label="optional fields" labelPosition="center" m="sm" />

          <TextInput
            data-qa="registration-household-name-input"
            label="Household Name"
            placeholder="username's Beloved Family"
            {...registrationForm.getInputProps('householdName')}
          />

          <DatePicker
            data-qa="registration-birthday-input"
            placeholder="optional :)"
            initialLevel="date"
            label="Birthday"
            dropdownType="popover"
            dropdownPosition="bottom-start"
            initialMonth={subYears(new Date(), 13)} // new Date('1970-01-02')
            maxDate={subYears(new Date(), 13)} // COPPA
            {...registrationForm.getInputProps('birthday')}
          />
          <Grid>
            <Grid.Col span={6}>
              <TextInput
                data-qa="registration-first-name-input"
                label="First Name"
                placeholder="optional :)"
                {...registrationForm.getInputProps('firstName')}
              />
            </Grid.Col>
            <Grid.Col span={6}>
              <TextInput
                data-qa="registration-last-name-input"
                label="Last Name"
                placeholder="optional :)"
                {...registrationForm.getInputProps('lastName')}
              />
            </Grid.Col>
          </Grid>

          {registrationError && (
            <>
              <Space h="md" />
              <Alert title="Oh no!" color="tomato">
                {registrationError}
              </Alert>
            </>
          )}

          <Group position="center">
            <Button data-qa="registration-button" type="submit" mt="lg" fullWidth>
              Register
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
                <Link href="/login">Login instead</Link>
              </Text>
            </Grid.Col>
          </Grid>
        </form>
      </Container>
    </AppLayout>
  );
}
