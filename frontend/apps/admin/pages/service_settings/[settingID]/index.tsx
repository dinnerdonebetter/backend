import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { useForm, zodResolver } from '@mantine/form';
import { TextInput, Button, Group, Container, Divider, Space, Select } from '@mantine/core';
import { useState } from 'react';
import { useRouter } from 'next/router';
import { z } from 'zod';

import { ServiceSetting, ServiceSettingUpdateRequestInput } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { AppLayout } from '../../../src/layouts';
import { buildLocalClient, buildServerSideClient } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';

declare interface ServiceSettingPageProps {
  pageLoadServiceSetting: ServiceSetting;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ServiceSettingPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ServiceSettingPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const { settingID } = context.query;
  if (!settingID) {
    throw new Error('service setting ID is somehow missing!');
  }

  const fetchServiceSettingTimer = timing.addEvent('fetch service setting');
  const pageLoadServiceSettingPromise = apiClient
    .getServiceSetting(settingID.toString())
    .then((result: ServiceSetting) => {
      span.addEvent('service setting retrieved');
      return result;
    })
    .finally(() => {
      fetchServiceSettingTimer.end();
    });

  const [pageLoadServiceSetting] = await Promise.all([pageLoadServiceSettingPromise]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: { pageLoadServiceSetting },
  };
};

const serviceSettingUpdateFormSchema = z.object({
  name: z.string().min(1, 'name is required').trim(),
});

function ServiceSettingPage(props: ServiceSettingPageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const { pageLoadServiceSetting } = props;

  const [serviceSetting, setServiceSetting] = useState<ServiceSetting>(pageLoadServiceSetting);
  const [originalServiceSetting, setOriginalServiceSetting] = useState<ServiceSetting>(pageLoadServiceSetting);

  const updateForm = useForm({
    initialValues: serviceSetting,
    validate: zodResolver(serviceSettingUpdateFormSchema),
  });

  const dataHasChanged = (): boolean => {
    return (
      originalServiceSetting.name !== updateForm.values.name ||
      originalServiceSetting.description !== updateForm.values.description ||
      originalServiceSetting.type !== updateForm.values.type ||
      originalServiceSetting.defaultValue !== updateForm.values.defaultValue ||
      originalServiceSetting.enumeration.join(',') !== updateForm.values.enumeration.join(',')
    );
  };

  const submit = async () => {
    const validation = updateForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ServiceSettingUpdateRequestInput({
      name: updateForm.values.name,
      description: updateForm.values.description,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .updateServiceSetting(serviceSetting.id, submission)
      .then((result: ServiceSetting) => {
        if (result) {
          updateForm.setValues(result);
          setServiceSetting(result);
          setOriginalServiceSetting(result);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Service Setting">
      <Container size="sm">
        <form onSubmit={updateForm.onSubmit(submit)}>
          <TextInput label="Name" placeholder="thing" {...updateForm.getInputProps('name')} />
          <TextInput label="Description" placeholder="thing" {...updateForm.getInputProps('description')} />

          <Select
            label="Type"
            required
            onChange={async (item: string) => {
              updateForm.setFieldValue('type', item);
            }}
            {...updateForm.getInputProps('type')}
            data={['user', 'household', 'membership']}
          />

          <TextInput label="Default Value" placeholder="thing" {...updateForm.getInputProps('defaultValue')} />
          <TextInput label="Enumeration" placeholder="thing" {...updateForm.getInputProps('enumeration')} />

          <Group position="center">
            <Button type="submit" mt="sm" fullWidth disabled={!dataHasChanged()}>
              Submit
            </Button>
            <Button
              type="submit"
              color="red"
              fullWidth
              onClick={() => {
                if (confirm('Are you sure you want to delete this service setting?')) {
                  apiClient.deleteServiceSetting(serviceSetting.id).then(() => {
                    router.push('/service_settings');
                  });
                }
              }}
            >
              Delete
            </Button>
          </Group>
        </form>

        <Space h="xl" />
        <Divider />
        <Space h="xl" />
      </Container>
    </AppLayout>
  );
}

export default ServiceSettingPage;
