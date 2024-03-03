import {
  ActionIcon,
  AppShell,
  Box,
  Center,
  Container,
  Footer,
  Grid,
  Header,
  Space,
  useMantineColorScheme,
} from '@mantine/core';
import { IconSun } from '@tabler/icons';
import Head from 'next/head';
import Link from 'next/link';
import React from 'react';

class AppLayoutProps {
  title: string = 'NO TITLE';
  containerSize?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' = 'xl';
  disableTitlePrefix?: boolean = false;
  children: React.ReactNode;
}

export function AppLayout(props: AppLayoutProps) {
  const { title, containerSize, disableTitlePrefix, children } = props;
  const { toggleColorScheme } = useMantineColorScheme();

  const pageTitle = `${disableTitlePrefix ? '' : 'Dinner Done Better'}${title ? ` - ${title}` : ''}`;

  let loginPath = '';
  switch (process.env.NEXT_PUBLIC_ENVIRONMENT) {
    case 'localdev':
      loginPath = 'http://localhost:7000/login';
      break;
    case 'prod':
      loginPath = 'https://app.dinnerdonebetter.com/login';
      break;
    case 'dev':
      loginPath = 'https://app.dinnerdonebetter.dev/login';
      break;
    default:
      loginPath = 'https://app.dinnerdonebetter.com/login';
      break;
  }

  const header = (
    <Header height={50} p="xs">
      <Grid>
        <Grid.Col span={3}></Grid.Col>

        <Grid.Col span="auto">
          <Center>
            Dinner Done Better
            {/*
            <Image
              src={logo}
              alt="logo"
              priority
              style={{
                height: 'auto',
                width: '6rem',
              }}
            />
            */}
          </Center>
        </Grid.Col>

        <Grid.Col span={3}>
          <Link href={loginPath}>Login</Link>
        </Grid.Col>
      </Grid>
    </Header>
  );

  const footer = (
    <>
      <Space h="xl" />
      <Space h="xl" />
      <Box>
        <Footer height={40} mt="xl" p="xs" pt={5} fixed>
          <ActionIcon onClick={() => toggleColorScheme()} sx={{ float: 'left' }} aria-label="toggle color scheme">
            <IconSun />
          </ActionIcon>
        </Footer>
      </Box>
    </>
  );

  return (
    <AppShell fixed={false} padding="md" header={header} footer={footer}>
      <Head>
        <title>{pageTitle}</title>
      </Head>
      <Container size={containerSize}>{children}</Container>
    </AppShell>
  );
}
