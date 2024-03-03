import { useRouter } from 'next/router';
import { useForm, zodResolver } from '@mantine/form';
import { TextInput, Button, Group, Container } from '@mantine/core';
import { z } from 'zod';

import { ValidIngredientGroup, ValidIngredientGroupCreationRequestInput } from '@dinnerdonebetter/models';

import { AppLayout } from '../../src/layouts';
import { buildLocalClient } from '../../src/client';
import { inputSlug } from '../../src/schemas';

const validIngredientGroupCreationFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  slug: inputSlug,
});

export default function ValidIngredientGroupCreator(): JSX.Element {
  const router = useRouter();

  const creationForm = useForm({
    initialValues: {
      name: '',
      description: '',
      slug: '',
    },
    validate: zodResolver(validIngredientGroupCreationFormSchema),
  });

  const submit = async () => {
    const validation = creationForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ValidIngredientGroupCreationRequestInput({
      name: creationForm.values.name,
      description: creationForm.values.description,
      slug: creationForm.values.slug,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .createValidIngredientGroup(submission)
      .then((result: ValidIngredientGroup) => {
        if (result) {
          router.push(`/valid_ingredient_groups/${result.id}`);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="New valid ingredient group">
      <Container size="sm">
        <form onSubmit={creationForm.onSubmit(submit)}>
          <TextInput label="Name" placeholder="thing" {...creationForm.getInputProps('name')} />
          <TextInput label="Slug" placeholder="thing" {...creationForm.getInputProps('slug')} />
          <TextInput
            label="Description"
            placeholder="stuff about things"
            {...creationForm.getInputProps('description')}
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
