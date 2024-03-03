import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { useForm, zodResolver } from '@mantine/form';
import { TextInput, Button, Group, Container, Switch } from '@mantine/core';
import { z } from 'zod';
import { useState } from 'react';
import { useRouter } from 'next/router';

import { ValidVessel, ValidVesselUpdateRequestInput } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { AppLayout } from '../../../src/layouts';
import { buildLocalClient, buildServerSideClient } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';
import { inputSlug } from '../../../src/schemas';

declare interface ValidVesselPageProps {
  pageLoadValidVessel: ValidVessel;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidVesselPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidVesselPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const { validVesselID } = context.query;
  if (!validVesselID) {
    throw new Error('valid vessel ID is somehow missing!');
  }

  const fetchValidVesselTimer = timing.addEvent('fetch valid vessel');
  const pageLoadValidVesselPromise = apiClient
    .getValidVessel(validVesselID.toString())
    .then((result: ValidVessel) => {
      span.addEvent('valid vessel retrieved');
      return result;
    })
    .finally(() => {
      fetchValidVesselTimer.end();
    });

  const [pageLoadValidVessel] = await Promise.all([pageLoadValidVesselPromise]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return { props: { pageLoadValidVessel } };
};

const validVesselUpdateFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  pluralName: z.string().trim().min(1, 'plural name is required'),
  slug: inputSlug,
});

function ValidVesselPage(props: ValidVesselPageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const { pageLoadValidVessel } = props;

  const [validVessel, setValidVessel] = useState<ValidVessel>(pageLoadValidVessel);
  const [originalValidVessel, setOriginalValidVessel] = useState<ValidVessel>(pageLoadValidVessel);

  const updateForm = useForm({
    initialValues: validVessel,
    validate: zodResolver(validVesselUpdateFormSchema),
  });

  const dataHasChanged = (): boolean => {
    return (
      originalValidVessel.description !== updateForm.values.description ||
      originalValidVessel.iconPath !== updateForm.values.iconPath ||
      originalValidVessel.name !== updateForm.values.name ||
      originalValidVessel.pluralName !== updateForm.values.pluralName ||
      originalValidVessel.slug !== updateForm.values.slug ||
      originalValidVessel.displayInSummaryLists !== updateForm.values.displayInSummaryLists ||
      originalValidVessel.includeInGeneratedInstructions !== updateForm.values.includeInGeneratedInstructions
    );
  };

  const submit = async () => {
    const validation = updateForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ValidVesselUpdateRequestInput({
      name: updateForm.values.name,
      pluralName: updateForm.values.pluralName,
      description: updateForm.values.description,
      iconPath: updateForm.values.iconPath,
      slug: updateForm.values.slug,
      displayInSummaryLists: updateForm.values.displayInSummaryLists,
      includeInGeneratedInstructions: updateForm.values.includeInGeneratedInstructions,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .updateValidVessel(validVessel.id, submission)
      .then((result: ValidVessel) => {
        if (result) {
          updateForm.setValues(result);
          setValidVessel(result);
          setOriginalValidVessel(result);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Valid Vessel">
      <Container size="sm">
        <form onSubmit={updateForm.onSubmit(submit)}>
          <TextInput label="Name" placeholder="thing" {...updateForm.getInputProps('name')} />
          <TextInput label="Plural Name" placeholder="things" {...updateForm.getInputProps('pluralName')} />
          <TextInput label="Slug" placeholder="thing" {...updateForm.getInputProps('slug')} />
          <TextInput label="Description" placeholder="thing" {...updateForm.getInputProps('description')} />

          <Switch
            checked={updateForm.values.displayInSummaryLists}
            label="Display in summary lists"
            {...updateForm.getInputProps('displayInSummaryLists')}
          />

          <Switch
            checked={updateForm.values.includeInGeneratedInstructions}
            label="Include in generated instructions"
            {...updateForm.getInputProps('includeInGeneratedInstructions')}
          />

          <Group position="center">
            <Button type="submit" mt="sm" fullWidth disabled={!dataHasChanged()}>
              Submit
            </Button>
            <Button
              type="submit"
              color="red"
              fullWidth
              onClick={() => {
                if (confirm('Are you sure you want to delete this valid vessel?')) {
                  apiClient.deleteValidVessel(validVessel.id).then(() => {
                    router.push('/valid_vessels');
                  });
                }
              }}
            >
              Delete
            </Button>
          </Group>
        </form>
      </Container>
    </AppLayout>
  );
}

export default ValidVesselPage;
