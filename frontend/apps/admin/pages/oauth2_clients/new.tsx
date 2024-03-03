import { useRouter } from 'next/router';
import { useForm, zodResolver } from '@mantine/form';
import { TextInput, Button, Group, Container } from '@mantine/core';
import { z } from 'zod';

import { OAuth2Client, OAuth2ClientCreationRequestInput } from '@dinnerdonebetter/models';

import { AppLayout } from '../../src/layouts';
import { buildLocalClient } from '../../src/client';

const oauth2ClientCreationFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
});

export default function OAuth2ClientCreator(): JSX.Element {
  const router = useRouter();

  const creationForm = useForm({
    initialValues: {
      name: '',
      description: '',
      yieldsNothing: false,
      restrictToIngredients: true,
      pastTense: '',
      slug: '',
      minimumIngredientCount: 1,
      maximumIngredientCount: undefined,
      minimumInstrumentCount: 1,
      maximumInstrumentCount: undefined,
      temperatureRequired: false,
      timeEstimateRequired: false,
      consumesVessel: false,
      onlyForVessels: false,
      minimumVesselCount: 1,
      maximumVesselCount: undefined,
    },
    validate: zodResolver(oauth2ClientCreationFormSchema),
  });

  const submit = async () => {
    const validation = creationForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new OAuth2ClientCreationRequestInput({
      name: creationForm.values.name,
      description: creationForm.values.description,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .createOAuth2Client(submission)
      .then((result: OAuth2Client) => {
        if (result) {
          router.push(`/oauth2_clients/${result.id}`);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Create New Valid Preparation">
      <Container size="sm">
        <form onSubmit={creationForm.onSubmit(submit)}>
          <TextInput label="Name" placeholder="thing" {...creationForm.getInputProps('name')} />
          <TextInput label="Description" placeholder="thing" {...creationForm.getInputProps('description')} />

          <Group position="center">
            <Button type="submit" mt="sm" fullWidth>
              Submit
            </Button>
          </Group>
        </form>
      </Container>
    </AppLayout>
  );
}
