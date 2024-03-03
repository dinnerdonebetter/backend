import { useRouter } from 'next/router';
import { useForm, zodResolver } from '@mantine/form';
import { TextInput, Button, Group, Container, Switch, NumberInput } from '@mantine/core';
import { z } from 'zod';

import { ValidIngredient, ValidIngredientCreationRequestInput } from '@dinnerdonebetter/models';

import { AppLayout } from '../../src/layouts';
import { buildLocalClient } from '../../src/client';
import { inputSlug } from '../../src/schemas';

const validIngredientCreationFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  pluralName: z.string().trim().min(1, 'plural name is required'),
  slug: inputSlug,
});

export default function ValidIngredientCreator(): JSX.Element {
  const router = useRouter();

  const creationForm = useForm({
    initialValues: {
      name: '',
      pluralName: '',
      description: '',
      warning: '',
      iconPath: '',
      containsDairy: false,
      containsPeanut: false,
      containsTreeNut: false,
      containsEgg: false,
      containsWheat: false,
      containsShellfish: false,
      containsSesame: false,
      containsFish: false,
      containsGluten: false,
      animalFlesh: false,
      isMeasuredVolumetrically: false,
      isLiquid: false,
      containsSoy: false,
      animalDerived: false,
      restrictToPreparations: true,
      minimumIdealStorageTemperatureInCelsius: undefined,
      maximumIdealStorageTemperatureInCelsius: undefined,
      slug: '',
      shoppingSuggestions: '',
      containsAlcohol: false,
    },
    validate: zodResolver(validIngredientCreationFormSchema),
  });

  const submit = async () => {
    const validation = creationForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ValidIngredientCreationRequestInput({
      name: creationForm.values.name,
      pluralName: creationForm.values.pluralName,
      description: creationForm.values.description,
      warning: creationForm.values.warning,
      iconPath: creationForm.values.iconPath,
      containsDairy: creationForm.values.containsDairy,
      containsPeanut: creationForm.values.containsPeanut,
      containsTreeNut: creationForm.values.containsTreeNut,
      containsEgg: creationForm.values.containsEgg,
      containsWheat: creationForm.values.containsWheat,
      containsShellfish: creationForm.values.containsShellfish,
      containsSesame: creationForm.values.containsSesame,
      containsFish: creationForm.values.containsFish,
      containsGluten: creationForm.values.containsGluten,
      animalFlesh: creationForm.values.animalFlesh,
      isMeasuredVolumetrically: creationForm.values.isMeasuredVolumetrically,
      isLiquid: creationForm.values.isLiquid,
      containsSoy: creationForm.values.containsSoy,
      animalDerived: creationForm.values.animalDerived,
      restrictToPreparations: creationForm.values.restrictToPreparations,
      containsAlcohol: creationForm.values.containsAlcohol,
      minimumIdealStorageTemperatureInCelsius: creationForm.values.minimumIdealStorageTemperatureInCelsius,
      maximumIdealStorageTemperatureInCelsius: creationForm.values.maximumIdealStorageTemperatureInCelsius,
      slug: creationForm.values.slug,
      shoppingSuggestions: creationForm.values.shoppingSuggestions,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .createValidIngredient(submission)
      .then((result: ValidIngredient) => {
        if (result.id) {
          router.push(`/valid_ingredients/${result.id}`);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Create New Valid Ingredient">
      <Container size="sm">
        <form onSubmit={creationForm.onSubmit(submit)}>
          <TextInput label="Name" placeholder="thing" {...creationForm.getInputProps('name')} />
          <TextInput label="Slug" placeholder="thing" {...creationForm.getInputProps('slug')} />
          <TextInput label="Plural Name" placeholder="things" {...creationForm.getInputProps('pluralName')} />
          <TextInput
            label="Description"
            placeholder="stuff about things"
            {...creationForm.getInputProps('description')}
          />
          <TextInput label="Warning" placeholder="warning" {...creationForm.getInputProps('warning')} />
          <NumberInput
            label="Min Storage Temp (C°)"
            precision={5}
            {...creationForm.getInputProps('minimumIdealStorageTemperatureInCelsius')}
          />
          <NumberInput
            label="Max Storage Temp (C°)"
            precision={5}
            {...creationForm.getInputProps('maximumIdealStorageTemperatureInCelsius')}
          />
          <Switch
            checked={creationForm.values.containsDairy}
            label="Contains Dairy"
            {...creationForm.getInputProps('containsDairy')}
          />
          <Switch
            checked={creationForm.values.containsPeanut}
            label="Contains Peanut"
            {...creationForm.getInputProps('containsPeanut')}
          />
          <Switch
            checked={creationForm.values.containsTreeNut}
            label="Contains Tree Nut"
            {...creationForm.getInputProps('containsTreeNut')}
          />
          <Switch
            checked={creationForm.values.containsEgg}
            label="Contains Egg"
            {...creationForm.getInputProps('containsEgg')}
          />
          <Switch
            checked={creationForm.values.containsWheat}
            label="Contains Wheat"
            {...creationForm.getInputProps('containsWheat')}
          />
          <Switch
            checked={creationForm.values.containsShellfish}
            label="Contains Shellfish"
            {...creationForm.getInputProps('containsShellfish')}
          />
          <Switch
            checked={creationForm.values.containsSesame}
            label="Contains Sesame"
            {...creationForm.getInputProps('containsSesame')}
          />
          <Switch
            checked={creationForm.values.containsFish}
            label="Contains Fish"
            {...creationForm.getInputProps('containsFish')}
          />
          <Switch
            checked={creationForm.values.containsGluten}
            label="Contains Gluten"
            {...creationForm.getInputProps('containsGluten')}
          />
          <Switch
            checked={creationForm.values.containsSoy}
            label="Contains Soy"
            {...creationForm.getInputProps('containsSoy')}
          />
          <Switch
            checked={creationForm.values.containsAlcohol}
            label="Contains Alcohol"
            {...creationForm.getInputProps('containsAlcohol')}
          />
          <Switch
            checked={creationForm.values.animalFlesh}
            label="Animal Flesh"
            {...creationForm.getInputProps('animalFlesh')}
          />
          <Switch
            checked={creationForm.values.animalDerived}
            label="Animal Derived"
            {...creationForm.getInputProps('animalDerived')}
          />
          <Switch
            checked={creationForm.values.isMeasuredVolumetrically}
            label="Measured Volumetrically"
            {...creationForm.getInputProps('isMeasuredVolumetrically')}
          />
          <Switch checked={creationForm.values.isLiquid} label="Liquid" {...creationForm.getInputProps('isLiquid')} />
          <Switch
            checked={creationForm.values.restrictToPreparations}
            label="Restrict To Preparations"
            {...creationForm.getInputProps('restrictToPreparations')}
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
