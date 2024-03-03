import {
  ActionIcon,
  AppShell,
  Avatar,
  Box,
  Burger,
  Center,
  Container,
  Footer,
  Grid,
  Group,
  Header,
  Navbar,
  NavLink,
  Space,
  Text,
  useMantineColorScheme,
} from '@mantine/core';
import {
  IconCalendarEvent,
  IconFlame,
  IconHome,
  IconLogin,
  IconLogout,
  IconNotebook,
  IconSettings,
  IconSoup,
  IconSun,
  IconToolsKitchen,
  IconUser,
} from '@tabler/icons';
import axios from 'axios';
import Head from 'next/head';
import { useRouter } from 'next/router';
import React, { useState } from 'react';

interface AppLayoutProps {
  title: string;
  titlePosition?: 'right' | 'left';
  containerSize?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  disableTitlePrefix?: boolean;
  children?: React.ReactNode;
  userLoggedIn: boolean;
  userAvatar?: string;
}

export function AppLayout(
  props: AppLayoutProps = {
    title: 'NO TITLE',
    titlePosition: 'right',
    containerSize: 'xl',
    disableTitlePrefix: false,
    userLoggedIn: false,
    userAvatar: '',
  },
) {
  // TODO: how do I know if I'm authed here?

  const { title, titlePosition, containerSize, disableTitlePrefix, children, userLoggedIn, userAvatar } = props;
  const router = useRouter();
  const [opened, setOpened] = useState(false);
  const navVerb = opened ? 'Close' : 'Open';

  const { toggleColorScheme } = useMantineColorScheme();

  const logout = async () => {
    await axios
      .post('/api/logout')
      .then(() => {
        router.push('/login');
      })
      .catch((err) => {
        router.push('/login');
        console.error(err);
      });
  };

  const pageTitle =
    (titlePosition || 'right') === 'right'
      ? `${disableTitlePrefix ? '' : 'Dinner Done Better'}${title ? ` - ${title}` : ''}`
      : `${title ?? ''}${disableTitlePrefix ? '' : ' - Dinner Done Better'}`;

  const header = (
    <Header height={50} p="xs">
      <Grid>
        <Grid.Col span={3}>
          <Burger
            size="sm"
            opened={opened}
            onClick={() => setOpened((o) => !o)}
            title={`${navVerb} navigation`}
            aria-label={`${navVerb} navigation`}
          />
        </Grid.Col>

        <Grid.Col span="auto">
          <Center>
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
          {userAvatar && <Avatar radius={100} component="a" src={userAvatar} alt="user avatar" />}
        </Grid.Col>
      </Grid>
    </Header>
  );

  const pathName = router.pathname;
  const navBar = (
    <Navbar width={{ base: 200 }} fixed={true} hiddenBreakpoint="xl" hidden={!opened}>
      <Navbar.Section mx="-xs" px="xs" grow>
        <NavLink
          label="Cooking"
          icon={<IconFlame size={16} />}
          childrenOffset={28}
          defaultOpened={pathName === '/' || (pathName.match(/^\/(recipes)/g) || []).length > 0}
        >
          <NavLink
            icon={<IconNotebook size={16} />}
            label="Recipes"
            onClick={() => router.push('/recipes')}
            active={(pathName.match(/^\/(recipes)\/(\w{20,})/g) || []).length > 0}
          />
        </NavLink>

        <NavLink
          label="Eating"
          icon={<IconToolsKitchen size={16} />}
          childrenOffset={28}
          defaultOpened={pathName === '/' || (pathName.match(/^\/(meal_plans|meals)/g) || []).length > 0}
        >
          <NavLink
            icon={<IconCalendarEvent size={16} />}
            label="Meal Plans"
            onClick={() => router.push('/meal_plans')}
            active={pathName.startsWith('/meal_plans')}
          />

          <NavLink
            icon={<IconSoup size={16} />}
            label="Meals"
            onClick={() => router.push('/meals')}
            active={pathName.startsWith('/meals')}
          />
        </NavLink>

        <NavLink
          label="Settings"
          icon={<IconSettings size={16} />}
          childrenOffset={28}
          defaultOpened={pathName === '/' || (pathName.match(/^\/(settings)/g) || []).length > 0}
        >
          <NavLink
            icon={<IconHome size={16} />}
            label="Household"
            onClick={() => router.push('/settings/household')}
            active={pathName.startsWith('/settings/household')}
          />
          <NavLink
            icon={<IconUser size={16} />}
            label="User"
            onClick={() => router.push('/settings/user')}
            active={pathName.startsWith('/settings/user')}
          />
        </NavLink>
      </Navbar.Section>
    </Navbar>
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

          {/* TODO: figure out when to show this, depending on auth status */}
          <Box sx={{ float: 'right' }}>
            {userLoggedIn ? (
              <Group>
                <Text weight="300" size="xs" color="tomato" mr="-sm">
                  Log off
                </Text>
                <ActionIcon onClick={() => logout()} aria-label="logout">
                  <IconLogout color="tomato" />
                </ActionIcon>
              </Group>
            ) : (
              !pathName.startsWith('/login') && (
                <Group>
                  <Text weight="300" size="xs" mr="-sm">
                    Login
                  </Text>
                  <ActionIcon onClick={() => router.push('/login')} aria-label="logout">
                    <IconLogin />
                  </ActionIcon>
                </Group>
              )
            )}
          </Box>
        </Footer>
      </Box>
    </>
  );

  return (
    <AppShell fixed={false} padding="md" header={header} navbar={navBar} footer={footer}>
      <Head>
        <title>{pageTitle}</title>
      </Head>
      <Container size={containerSize ?? 'md'}>{children}</Container>
    </AppShell>
  );
}
