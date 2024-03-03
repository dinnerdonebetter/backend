import { useRouter } from 'next/router';
import { useForm, zodResolver } from '@mantine/form';
import { TextInput, Button, Group, Container, Switch } from '@mantine/core';
import { z } from 'zod';

import { ValidMeasurementUnit, ValidMeasurementUnitCreationRequestInput } from '@dinnerdonebetter/models';

import { AppLayout } from '../../src/layouts';
import { buildLocalClient } from '../../src/client';
import { inputSlug } from '../../src/schemas';

const validMeasurementUnitCreationFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  pluralName: z.string().trim().min(1, 'plural name is required'),
  slug: inputSlug,
});

export default function ValidMeasurementUnitCreator(): JSX.Element {
  const router = useRouter();

  const creationForm = useForm({
    initialValues: {
      name: '',
      description: '',
      volumetric: false,
      universal: false,
      metric: false,
      imperial: false,
      pluralName: '',
      slug: '',
    },
    validate: zodResolver(validMeasurementUnitCreationFormSchema),
  });

  const submit = async () => {
    const validation = creationForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    if (creationForm.values.imperial && creationForm.values.metric) {
      creationForm.setErrors({
        metric: 'Cannot be imperial and metric at the same time',
        imperial: 'Cannot be imperial and metric at the same time',
      });
      return;
    }

    const submission = new ValidMeasurementUnitCreationRequestInput({
      name: creationForm.values.name,
      description: creationForm.values.description,
      volumetric: creationForm.values.volumetric,
      universal: creationForm.values.universal,
      metric: creationForm.values.metric,
      imperial: creationForm.values.imperial,
      pluralName: creationForm.values.pluralName,
      slug: creationForm.values.slug,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .createValidMeasurementUnit(submission)
      .then((result: ValidMeasurementUnit) => {
        if (result) {
          router.push(`/valid_measurement_units/${result.id}`);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Create New Valid Measurement Unit">
      <Container size="sm">
        <form onSubmit={creationForm.onSubmit(submit)}>
          <TextInput label="Name" placeholder="thing" {...creationForm.getInputProps('name')} />
          <TextInput label="Plural Name" placeholder="things" {...creationForm.getInputProps('pluralName')} />
          <TextInput label="Slug" placeholder="thing" {...creationForm.getInputProps('slug')} />
          <TextInput label="Description" placeholder="thing" {...creationForm.getInputProps('description')} />

          <Switch
            checked={creationForm.values.volumetric}
            label="Volumetric"
            {...creationForm.getInputProps('volumetric')}
          />
          <Switch
            checked={creationForm.values.universal}
            label="Universal"
            {...creationForm.getInputProps('universal')}
          />
          <Switch
            checked={creationForm.values.metric}
            disabled={creationForm.values.imperial}
            label="Metric"
            {...creationForm.getInputProps('metric')}
          />
          <Switch
            checked={creationForm.values.imperial}
            disabled={creationForm.values.metric}
            label="Imperial"
            {...creationForm.getInputProps('imperial')}
          />

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
