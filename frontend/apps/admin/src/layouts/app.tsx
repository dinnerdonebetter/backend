import {
  ActionIcon,
  AppShell,
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
  IconCheese,
  IconCooker,
  IconFlame,
  IconList,
  IconLogout,
  IconPyramid,
  IconRuler2,
  IconSun,
  IconToolsKitchen,
  IconUsers,
} from '@tabler/icons';
import axios from 'axios';
import Head from 'next/head';
import { useRouter } from 'next/router';
import React, { useState } from 'react';

class AppLayoutProps {
  title: string = 'NO TITLE';
  containerSize?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' = 'xl';
  disableTitlePrefix?: boolean = false;
  children: React.ReactNode;
}

export function AppLayout(props: AppLayoutProps) {
  // TODO: how do I know if I'm authed here?

  const { title, containerSize, disableTitlePrefix, children } = props;
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
        console.error(err);
      });
  };

  const pageTitle = `${disableTitlePrefix ? '' : 'Dinner Done Better'}${title ? ` - ${title}` : ''}`;

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

        <Grid.Col span={3}>{/*  */}</Grid.Col>
      </Grid>
    </Header>
  );

  const pathName = router.pathname;

  const navBar = (
    <Navbar width={{ base: 200 }} fixed={true} hiddenBreakpoint="xl" hidden={!opened}>
      <Navbar.Section mx="-xs" px="xs" grow>
        <NavLink
          label="Enumerations"
          icon={<IconList size={16} />}
          childrenOffset={28}
          defaultOpened={pathName === '/' || (pathName.match(/^\/(valid_)/g) || []).length > 0}
        >
          <NavLink
            icon={<IconFlame size={16} />}
            label="Preparations"
            onClick={() => router.push('/valid_preparations')}
            active={pathName.startsWith('/valid_preparations')}
          />
          <NavLink
            icon={<IconCheese size={16} />}
            label="Ingredients"
            onClick={() => router.push('/valid_ingredients')}
            active={pathName.startsWith('/valid_ingredients')}
          />
          <NavLink
            icon={<IconRuler2 size={16} />}
            label="Measurement Units"
            onClick={() => router.push('/valid_measurement_units')}
            active={pathName.startsWith('/valid_measurement_units')}
          />
          <NavLink
            icon={<IconPyramid size={16} />}
            label="Ingredient States"
            onClick={() => router.push('/valid_ingredient_states')}
            active={pathName.startsWith('/valid_ingredient_states')}
          />
          <NavLink
            icon={<IconToolsKitchen size={16} />}
            label="Instruments"
            onClick={() => router.push('/valid_instruments')}
            active={pathName.startsWith('/valid_instruments')}
          />
          <NavLink
            icon={<IconCooker size={16} />}
            label="Vessels"
            onClick={() => router.push('/valid_vessels')}
            active={pathName.startsWith('/valid_vessels')}
          />
        </NavLink>

        <NavLink
          icon={<IconUsers size={16} />}
          label="Recipe Creator"
          onClick={() => router.push('/recipes/new')}
          active={pathName.startsWith('/recipes/new')}
        />

        <NavLink
          icon={<IconUsers size={16} />}
          label="Users"
          onClick={() => router.push('/users')}
          active={pathName.startsWith('/users')}
        />

        <NavLink
          icon={<IconUsers size={16} />}
          label="OAuth2 Clients"
          onClick={() => router.push('/oauth2_clients')}
          active={pathName.startsWith('/oauth2_clients')}
        />
      </Navbar.Section>
    </Navbar>
  );

  const footer = (
    <>
      <Space h="xl" />
      <Space h="xl" />
      <Box>
        <Space mt="xl" h="xl" />
        <Footer height={40} mt="xl" p="xs" pt={5} fixed>
          <ActionIcon onClick={() => toggleColorScheme()} sx={{ float: 'left' }} aria-label="toggle color scheme">
            <IconSun />
          </ActionIcon>

          {/* TODO: figure out when to show this, depending on auth status */}
          <Box sx={{ float: 'right' }}>
            <Group>
              <Text weight="300" size="xs" color="tomato" mr="-sm">
                Log off
              </Text>
              <ActionIcon onClick={() => logout()} aria-label="logout">
                <IconLogout color="tomato" />
              </ActionIcon>
            </Group>
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
      <Container size={containerSize}>{children}</Container>
    </AppShell>
  );
}
