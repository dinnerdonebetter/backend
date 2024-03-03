import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { useForm, zodResolver } from '@mantine/form';
import {
  TextInput,
  Button,
  Group,
  Container,
  Switch,
  NumberInput,
  Autocomplete,
  Divider,
  Space,
  Title,
  ActionIcon,
  AutocompleteItem,
  Center,
  Grid,
  Pagination,
  Table,
  Text,
  ThemeIcon,
} from '@mantine/core';
import { AxiosError } from 'axios';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import { IconTrash } from '@tabler/icons';
import { useRouter } from 'next/router';
import { z } from 'zod';

import {
  ValidIngredient,
  ValidIngredientPreparation,
  ValidPreparation,
  QueryFilteredResult,
  ValidPreparationInstrument,
  ValidPreparationUpdateRequestInput,
  ValidIngredientPreparationCreationRequestInput,
  ValidPreparationInstrumentCreationRequestInput,
  ValidInstrument,
} from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { AppLayout } from '../../../src/layouts';
import { buildLocalClient, buildServerSideClient } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';
import { inputSlug } from '../../../src/schemas';

declare interface ValidPreparationPageProps {
  pageLoadValidPreparation: ValidPreparation;
  pageLoadValidPreparationInstruments: QueryFilteredResult<ValidPreparationInstrument>;
  pageLoadValidIngredientPreparations: QueryFilteredResult<ValidIngredientPreparation>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidPreparationPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidPreparationPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const { validPreparationID } = context.query;
  if (!validPreparationID) {
    throw new Error('valid preparation ID is somehow missing!');
  }

  const fetchValidPreparationTimer = timing.addEvent('fetch valid preparation');
  const pageLoadValidPreparationPromise = apiClient
    .getValidPreparation(validPreparationID.toString())
    .then((result: ValidPreparation) => {
      span.addEvent('valid preparation retrieved');
      return result;
    })
    .finally(() => {
      fetchValidPreparationTimer.end();
    });

  const fetchValidPreparationInstrumentsTimer = timing.addEvent('fetch valid preparation instruments');
  const pageLoadValidPreparationInstrumentsPromise = apiClient
    .validPreparationInstrumentsForPreparationID(validPreparationID.toString())
    .then((result: QueryFilteredResult<ValidPreparationInstrument>) => {
      span.addEvent('valid preparation retrieved');
      return result;
    })
    .finally(() => {
      fetchValidPreparationInstrumentsTimer.end();
    });

  const fetchValidIngredientPreparationsTimer = timing.addEvent('fetch valid ingredient preparations');
  const pageLoadValidIngredientPreparationsPromise = apiClient
    .validIngredientPreparationsForPreparationID(validPreparationID.toString())
    .then((result: QueryFilteredResult<ValidIngredientPreparation>) => {
      span.addEvent('valid preparation retrieved');
      return result;
    })
    .finally(() => {
      fetchValidIngredientPreparationsTimer.end();
    });

  const [pageLoadValidPreparation, pageLoadValidPreparationInstruments, pageLoadValidIngredientPreparations] =
    await Promise.all([
      pageLoadValidPreparationPromise,
      pageLoadValidPreparationInstrumentsPromise,
      pageLoadValidIngredientPreparationsPromise,
    ]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: { pageLoadValidPreparation, pageLoadValidPreparationInstruments, pageLoadValidIngredientPreparations },
  };
};

const validPreparationUpdateFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  pastTense: z.string().trim().min(1, 'past tense is required'),
  slug: inputSlug,
});

function ValidPreparationPage(props: ValidPreparationPageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const { pageLoadValidPreparation, pageLoadValidPreparationInstruments, pageLoadValidIngredientPreparations } = props;

  const [validPreparation, setValidPreparation] = useState<ValidPreparation>(pageLoadValidPreparation);
  const [originalValidPreparation, setOriginalValidPreparation] = useState<ValidPreparation>(pageLoadValidPreparation);

  const [newIngredientForPreparationInput, setNewIngredientForPreparationInput] =
    useState<ValidIngredientPreparationCreationRequestInput>(
      new ValidIngredientPreparationCreationRequestInput({
        validPreparationID: validPreparation.id,
      }),
    );
  const [ingredientQuery, setIngredientQuery] = useState('');
  const [ingredientsForPreparation, setIngredientsForPreparation] = useState<
    QueryFilteredResult<ValidIngredientPreparation>
  >(pageLoadValidIngredientPreparations);
  const [suggestedIngredients, setSuggestedIngredients] = useState<ValidIngredient[]>([]);

  useEffect(() => {
    if (ingredientQuery.length <= 2) {
      setSuggestedIngredients([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidIngredients(ingredientQuery)
      .then((res: QueryFilteredResult<ValidIngredient>) => {
        const newSuggestions = (res.data || []).filter((mu: ValidIngredient) => {
          return !(ingredientsForPreparation.data || []).some((vimu: ValidIngredientPreparation) => {
            return vimu.ingredient.id === mu.id;
          });
        });

        setSuggestedIngredients(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [ingredientQuery, ingredientsForPreparation.data]);

  const [newInstrumentForPreparationInput, setNewInstrumentForPreparationInput] =
    useState<ValidPreparationInstrumentCreationRequestInput>(
      new ValidPreparationInstrumentCreationRequestInput({
        validPreparationID: validPreparation.id,
      }),
    );
  const [instrumentQuery, setInstrumentQuery] = useState('');
  const [instrumentsForPreparation, setInstrumentsForPreparation] = useState<
    QueryFilteredResult<ValidPreparationInstrument>
  >(pageLoadValidPreparationInstruments);
  const [suggestedInstruments, setSuggestedInstruments] = useState<ValidInstrument[]>([]);

  useEffect(() => {
    if (instrumentQuery.length <= 2) {
      setSuggestedInstruments([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidInstruments(instrumentQuery)
      .then((res: ValidInstrument[]) => {
        const newSuggestions = res.filter((mu: ValidInstrument) => {
          return !(instrumentsForPreparation.data || []).some((vimu: ValidPreparationInstrument) => {
            return vimu.preparation.id === mu.id;
          });
        });

        setSuggestedInstruments(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [instrumentQuery, instrumentsForPreparation.data]);

  const updateForm = useForm({
    initialValues: validPreparation,
    validate: zodResolver(validPreparationUpdateFormSchema),
  });

  const dataHasChanged = (): boolean => {
    return (
      originalValidPreparation.name !== updateForm.values.name ||
      originalValidPreparation.description !== updateForm.values.description ||
      originalValidPreparation.yieldsNothing !== updateForm.values.yieldsNothing ||
      originalValidPreparation.restrictToIngredients !== updateForm.values.restrictToIngredients ||
      originalValidPreparation.pastTense !== updateForm.values.pastTense ||
      originalValidPreparation.slug !== updateForm.values.slug ||
      originalValidPreparation.minimumIngredientCount !== updateForm.values.minimumIngredientCount ||
      originalValidPreparation.maximumIngredientCount !== updateForm.values.maximumIngredientCount ||
      originalValidPreparation.minimumInstrumentCount !== updateForm.values.minimumInstrumentCount ||
      originalValidPreparation.maximumInstrumentCount !== updateForm.values.maximumInstrumentCount ||
      originalValidPreparation.temperatureRequired !== updateForm.values.temperatureRequired ||
      originalValidPreparation.timeEstimateRequired !== updateForm.values.timeEstimateRequired ||
      originalValidPreparation.consumesVessel !== updateForm.values.consumesVessel ||
      originalValidPreparation.onlyForVessels !== updateForm.values.onlyForVessels ||
      originalValidPreparation.minimumVesselCount !== updateForm.values.minimumVesselCount ||
      originalValidPreparation.maximumVesselCount !== updateForm.values.maximumVesselCount
    );
  };

  const submit = async () => {
    const validation = updateForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ValidPreparationUpdateRequestInput({
      name: updateForm.values.name,
      description: updateForm.values.description,
      yieldsNothing: updateForm.values.yieldsNothing,
      restrictToIngredients: updateForm.values.restrictToIngredients,
      pastTense: updateForm.values.pastTense,
      slug: updateForm.values.slug,
      minimumIngredientCount: updateForm.values.minimumIngredientCount,
      maximumIngredientCount: updateForm.values.maximumIngredientCount,
      minimumInstrumentCount: updateForm.values.minimumInstrumentCount,
      maximumInstrumentCount: updateForm.values.maximumInstrumentCount,
      temperatureRequired: updateForm.values.temperatureRequired,
      timeEstimateRequired: updateForm.values.timeEstimateRequired,
      consumesVessel: updateForm.values.consumesVessel,
      onlyForVessels: updateForm.values.onlyForVessels,
      minimumVesselCount: updateForm.values.minimumVesselCount,
      maximumVesselCount: updateForm.values.maximumVesselCount,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .updateValidPreparation(validPreparation.id, submission)
      .then((result: ValidPreparation) => {
        if (result) {
          updateForm.setValues(result);
          setValidPreparation(result);
          setOriginalValidPreparation(result);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Valid Preparation">
      <Container size="sm">
        <form onSubmit={updateForm.onSubmit(submit)}>
          <TextInput label="Name" placeholder="thing" {...updateForm.getInputProps('name')} />
          <TextInput label="Past Tense" placeholder="thinged" {...updateForm.getInputProps('pastTense')} />
          <TextInput label="Slug" placeholder="thing" {...updateForm.getInputProps('slug')} />
          <TextInput label="Description" placeholder="thing" {...updateForm.getInputProps('description')} />

          <Switch
            checked={updateForm.values.yieldsNothing}
            label="Yields Nothing"
            {...updateForm.getInputProps('yieldsNothing')}
          />
          <Switch
            checked={updateForm.values.restrictToIngredients}
            label="Restrict To Ingredients"
            {...updateForm.getInputProps('restrictToIngredients')}
          />
          <Switch
            checked={updateForm.values.temperatureRequired}
            label="Temperature Required"
            {...updateForm.getInputProps('temperatureRequired')}
          />
          <Switch
            checked={updateForm.values.timeEstimateRequired}
            label="Time Estimate Required"
            {...updateForm.getInputProps('timeEstimateRequired')}
          />

          <NumberInput label="Minimum Ingredient Count" {...updateForm.getInputProps('minimumIngredientCount')} />
          <NumberInput label="Maximum Ingredient Count" {...updateForm.getInputProps('maximumIngredientCount')} />
          <NumberInput label="Minimum Instrument Count" {...updateForm.getInputProps('minimumInstrumentCount')} />
          <NumberInput label="Maximum Instrument Count" {...updateForm.getInputProps('maximumInstrumentCount')} />

          <Switch
            checked={updateForm.values.consumesVessel}
            label="Consumes Vessel"
            {...updateForm.getInputProps('consumesVessel')}
          />
          <Switch
            checked={updateForm.values.onlyForVessels}
            label="Only For Vessels"
            {...updateForm.getInputProps('onlyForVessels')}
          />
          <NumberInput label="Minimum Vessel Count" {...updateForm.getInputProps('minimumVesselCount')} />
          <NumberInput label="Maximum Vessel Count" {...updateForm.getInputProps('maximumVesselCount')} />

          <Group position="center">
            <Button type="submit" mt="sm" fullWidth disabled={!dataHasChanged()}>
              Submit
            </Button>
            <Button
              type="submit"
              color="red"
              fullWidth
              onClick={() => {
                if (confirm('Are you sure you want to delete this valid preparation?')) {
                  apiClient.deleteValidPreparation(validPreparation.id).then(() => {
                    router.push('/valid_preparations');
                  });
                }
              }}
            >
              Delete
            </Button>
          </Group>
        </form>

        {/*

        INSTRUMENTS

        */}

        <Space h="xl" />
        <Divider />
        <Space h="xl" />

        <form>
          <Center>
            <Title order={4}>Instruments</Title>
          </Center>

          {instrumentsForPreparation.data && (instrumentsForPreparation.data || []).length !== 0 && (
            <>
              <Table mt="xl" withColumnBorders>
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Notes</th>
                    <th>
                      <Center>
                        <ThemeIcon variant="outline" color="gray">
                          <IconTrash size="sm" color="gray" />
                        </ThemeIcon>
                      </Center>
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {(instrumentsForPreparation.data || []).map(
                    (validInstrumentsForPreparation: ValidPreparationInstrument) => {
                      return (
                        <tr key={validInstrumentsForPreparation.id}>
                          <td>
                            <Link href={`/valid_instruments/${validInstrumentsForPreparation.instrument.id}`}>
                              {validInstrumentsForPreparation.instrument.name}
                            </Link>
                          </td>
                          <td>
                            <Text>{validInstrumentsForPreparation.notes}</Text>
                          </td>
                          <td>
                            <Center>
                              <ActionIcon
                                variant="outline"
                                aria-label="remove valid instrument measurement unit"
                                onClick={async () => {
                                  await apiClient
                                    .deleteValidPreparationInstrument(validInstrumentsForPreparation.id)
                                    .then(() => {
                                      setInstrumentsForPreparation({
                                        ...instrumentsForPreparation,
                                        data: instrumentsForPreparation.data.filter(
                                          (x: ValidPreparationInstrument) => x.id !== validInstrumentsForPreparation.id,
                                        ),
                                      });
                                    })
                                    .catch((error) => {
                                      console.error(error);
                                    });
                                }}
                              >
                                <IconTrash size="md" color="tomato" />
                              </ActionIcon>
                            </Center>
                          </td>
                        </tr>
                      );
                    },
                  )}
                </tbody>
              </Table>

              <Space h="xs" />

              <Pagination
                disabled={
                  Math.ceil(instrumentsForPreparation.totalCount / instrumentsForPreparation.limit) <=
                  instrumentsForPreparation.page
                }
                position="center"
                page={instrumentsForPreparation.page}
                total={Math.ceil(instrumentsForPreparation.totalCount / instrumentsForPreparation.limit)}
                onChange={(value: number) => {
                  setInstrumentsForPreparation({ ...instrumentsForPreparation, page: value });
                }}
              />
            </>
          )}

          <Grid>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="spoon"
                label="Instrument"
                value={instrumentQuery}
                onChange={setInstrumentQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedInstrument = suggestedInstruments.find((x: ValidInstrument) => x.name === item.value);

                  if (!selectedInstrument) {
                    console.error(`selectedInstrument not found for item ${item.value}}`);
                    return;
                  }

                  setNewInstrumentForPreparationInput({
                    ...newInstrumentForPreparationInput,
                    validInstrumentID: selectedInstrument.id,
                  });
                }}
                data={suggestedInstruments.map((x: ValidInstrument) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <TextInput
                label="Notes"
                value={newInstrumentForPreparationInput.notes}
                onChange={(event) =>
                  setNewInstrumentForPreparationInput({
                    ...newInstrumentForPreparationInput,
                    notes: event.target.value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span={2}>
              <Button
                mt="xl"
                disabled={
                  newInstrumentForPreparationInput.validPreparationID === '' ||
                  newInstrumentForPreparationInput.validInstrumentID === ''
                }
                onClick={async () => {
                  await apiClient
                    .createValidPreparationInstrument(newInstrumentForPreparationInput)
                    .then((res: ValidPreparationInstrument) => {
                      // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                      apiClient.getValidPreparationInstrument(res.id).then((res: ValidPreparationInstrument) => {
                        setInstrumentsForPreparation({
                          ...instrumentsForPreparation,
                          data: [...(instrumentsForPreparation.data || []), res],
                        });

                        setNewInstrumentForPreparationInput(
                          new ValidPreparationInstrumentCreationRequestInput({
                            validPreparationID: validPreparation.id,
                            validInstrumentID: '',
                            notes: '',
                          }),
                        );

                        setInstrumentQuery('');
                      });
                    })
                    .catch((error) => {
                      console.error(error);
                    });
                }}
              >
                Save
              </Button>
            </Grid.Col>
          </Grid>
        </form>

        {/*

        INGREDIENTS

        */}

        <Space h="xl" />
        <Divider />
        <Space h="xl" />

        <form>
          <Center>
            <Title order={4}>Ingredients</Title>
          </Center>

          {ingredientsForPreparation.data && (ingredientsForPreparation.data || []).length !== 0 && (
            <>
              <Table mt="xl" withColumnBorders>
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Notes</th>
                    <th>
                      <Center>
                        <ThemeIcon variant="outline" color="gray">
                          <IconTrash size="sm" color="gray" />
                        </ThemeIcon>
                      </Center>
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {(ingredientsForPreparation.data || []).map(
                    (validIngredientsForPreparation: ValidIngredientPreparation) => {
                      return (
                        <tr key={validIngredientsForPreparation.id}>
                          <td>
                            <Link href={`/valid_ingredients/${validIngredientsForPreparation.ingredient.id}`}>
                              {validIngredientsForPreparation.ingredient.name}
                            </Link>
                          </td>
                          <td>
                            <Text>{validIngredientsForPreparation.notes}</Text>
                          </td>
                          <td>
                            <Center>
                              <ActionIcon
                                variant="outline"
                                aria-label="remove valid ingredient measurement unit"
                                onClick={async () => {
                                  await apiClient
                                    .deleteValidIngredientPreparation(validIngredientsForPreparation.id)
                                    .then(() => {
                                      setIngredientsForPreparation({
                                        ...ingredientsForPreparation,
                                        data: ingredientsForPreparation.data.filter(
                                          (x: ValidIngredientPreparation) => x.id !== validIngredientsForPreparation.id,
                                        ),
                                      });
                                    })
                                    .catch((error) => {
                                      console.error(error);
                                    });
                                }}
                              >
                                <IconTrash size="md" color="tomato" />
                              </ActionIcon>
                            </Center>
                          </td>
                        </tr>
                      );
                    },
                  )}
                </tbody>
              </Table>

              <Space h="xs" />

              <Pagination
                disabled={
                  Math.ceil(ingredientsForPreparation.totalCount / ingredientsForPreparation.limit) <=
                  ingredientsForPreparation.page
                }
                position="center"
                page={ingredientsForPreparation.page}
                total={Math.ceil(ingredientsForPreparation.totalCount / ingredientsForPreparation.limit)}
                onChange={(value: number) => {
                  setIngredientsForPreparation({ ...ingredientsForPreparation, page: value });
                }}
              />
            </>
          )}

          <Grid>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="garlic"
                label="Ingredient"
                value={ingredientQuery}
                onChange={setIngredientQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedIngredient = suggestedIngredients.find((x: ValidIngredient) => x.name === item.value);

                  if (!selectedIngredient) {
                    console.error(`selectedIngredient not found for item ${item.value}}`);
                    return;
                  }

                  setNewIngredientForPreparationInput({
                    ...newIngredientForPreparationInput,
                    validIngredientID: selectedIngredient.id,
                  });
                }}
                data={suggestedIngredients.map((x: ValidIngredient) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <TextInput
                label="Notes"
                value={newIngredientForPreparationInput.notes}
                onChange={(event) =>
                  setNewIngredientForPreparationInput({
                    ...newIngredientForPreparationInput,
                    notes: event.target.value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span={2}>
              <Button
                mt="xl"
                disabled={
                  newIngredientForPreparationInput.validPreparationID === '' ||
                  newIngredientForPreparationInput.validIngredientID === ''
                }
                onClick={async () => {
                  await apiClient
                    .createValidIngredientPreparation(newIngredientForPreparationInput)
                    .then((res: ValidIngredientPreparation) => {
                      // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                      apiClient.getValidIngredientPreparation(res.id).then((res: ValidIngredientPreparation) => {
                        setIngredientsForPreparation({
                          ...ingredientsForPreparation,
                          data: [...(ingredientsForPreparation.data || []), res],
                        });

                        setNewIngredientForPreparationInput(
                          new ValidIngredientPreparationCreationRequestInput({
                            validPreparationID: validPreparation.id,
                            validIngredientID: '',
                            notes: '',
                          }),
                        );

                        setIngredientQuery('');
                      });
                    })
                    .catch((error) => {
                      console.error(error);
                    });
                }}
              >
                Save
              </Button>
            </Grid.Col>
          </Grid>
        </form>
      </Container>
    </AppLayout>
  );
}

export default ValidPreparationPage;
