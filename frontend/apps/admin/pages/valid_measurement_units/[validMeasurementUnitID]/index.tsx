import { AxiosError } from 'axios';
import { useForm, zodResolver } from '@mantine/form';
import {
  TextInput,
  Button,
  Group,
  Container,
  Switch,
  Autocomplete,
  Divider,
  Space,
  Title,
  ActionIcon,
  Text,
  AutocompleteItem,
  Center,
  Grid,
  Pagination,
  Table,
  ThemeIcon,
  NumberInput,
} from '@mantine/core';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';
import { IconTrash } from '@tabler/icons';
import { z } from 'zod';

import {
  ValidIngredient,
  ValidIngredientMeasurementUnit,
  ValidMeasurementUnit,
  ValidMeasurementUnitUpdateRequestInput,
  ValidIngredientMeasurementUnitCreationRequestInput,
  ValidMeasurementUnitConversion,
  ValidMeasurementUnitConversionCreationRequestInput,
  QueryFilteredResult,
} from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { AppLayout } from '../../../src/layouts';
import { buildLocalClient, buildServerSideClient } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';
import { inputSlug } from '../../../src/schemas';

declare interface ValidMeasurementUnitPageProps {
  pageLoadValidMeasurementUnit: ValidMeasurementUnit;
  pageLoadIngredientsForMeasurementUnit: QueryFilteredResult<ValidIngredientMeasurementUnit>;
  pageLoadMeasurementUnitConversionsFromUnit: ValidMeasurementUnitConversion[];
  pageLoadMeasurementUnitConversionsToUnit: ValidMeasurementUnitConversion[];
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidMeasurementUnitPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidMeasurementUnitPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  const { validMeasurementUnitID } = context.query;
  if (!validMeasurementUnitID) {
    throw new Error('valid measurement unit ID is somehow missing!');
  }

  const fetchValidMeasurementUnitTimer = timing.addEvent('fetch ');
  const pageLoadValidMeasurementUnitPromise = apiClient
    .getValidMeasurementUnit(validMeasurementUnitID.toString())
    .then((result: ValidMeasurementUnit) => {
      span.addEvent('valid measurement unit retrieved');
      return result;
    })
    .finally(() => {
      fetchValidMeasurementUnitTimer.end();
    });

  const fetchIngredientsForMeasurementUnitTimer = timing.addEvent('fetch valid ingredient measurement units');
  const pageLoadIngredientsForMeasurementUnitPromise = apiClient
    .validIngredientMeasurementUnitsForMeasurementUnitID(validMeasurementUnitID.toString())
    .then((res: QueryFilteredResult<ValidIngredientMeasurementUnit>) => {
      span.addEvent('valid ingredient measurement units retrieved');
      return res;
    })
    .finally(() => {
      fetchIngredientsForMeasurementUnitTimer.end();
    });

  const fetchMeasurementUnitConversionsFromUnitTimer = timing.addEvent('fetch measurement unit conversions from unit');
  const pageLoadMeasurementUnitConversionsFromUnitPromise = apiClient
    .getValidMeasurementUnitConversionsFromUnit(validMeasurementUnitID.toString())
    .then((res: ValidMeasurementUnitConversion[]) => {
      span.addEvent('valid ingredient measurement units retrieved');
      return res;
    })
    .finally(() => {
      fetchMeasurementUnitConversionsFromUnitTimer.end();
    });

  const fetchMeasurementUnitConversionsToUnitTimer = timing.addEvent('fetch measurement unit conversions to unit');
  const pageLoadMeasurementUnitConversionsToUnitPromise = apiClient
    .getValidMeasurementUnitConversionsToUnit(validMeasurementUnitID.toString())
    .then((res: ValidMeasurementUnitConversion[]) => {
      span.addEvent('valid ingredient measurement units retrieved');
      return res;
    })
    .finally(() => {
      fetchMeasurementUnitConversionsToUnitTimer.end();
    });

  const [
    pageLoadValidMeasurementUnit,
    pageLoadIngredientsForMeasurementUnit,
    pageLoadMeasurementUnitConversionsFromUnit,
    pageLoadMeasurementUnitConversionsToUnit,
  ] = await Promise.all([
    pageLoadValidMeasurementUnitPromise,
    pageLoadIngredientsForMeasurementUnitPromise,
    pageLoadMeasurementUnitConversionsFromUnitPromise,
    pageLoadMeasurementUnitConversionsToUnitPromise,
  ]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: {
      pageLoadValidMeasurementUnit,
      pageLoadIngredientsForMeasurementUnit,
      pageLoadMeasurementUnitConversionsFromUnit,
      pageLoadMeasurementUnitConversionsToUnit,
    },
  };
};

const validMeasurementUnitUpdateFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  pluralName: z.string().trim().min(1, 'plural name is required'),
  slug: inputSlug,
});

function ValidMeasurementUnitPage(props: ValidMeasurementUnitPageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const {
    pageLoadValidMeasurementUnit,
    pageLoadIngredientsForMeasurementUnit,
    pageLoadMeasurementUnitConversionsFromUnit,
    pageLoadMeasurementUnitConversionsToUnit,
  } = props;

  const [validMeasurementUnit, setValidMeasurementUnit] = useState<ValidMeasurementUnit>(pageLoadValidMeasurementUnit);
  const [originalValidMeasurementUnit, setOriginalValidMeasurementUnit] =
    useState<ValidMeasurementUnit>(pageLoadValidMeasurementUnit);

  const [newIngredientForMeasurementUnitInput, setNewIngredientForMeasurementUnitInput] =
    useState<ValidIngredientMeasurementUnitCreationRequestInput>(
      new ValidIngredientMeasurementUnitCreationRequestInput({
        validMeasurementUnitID: validMeasurementUnit.id,
        minimumAllowableQuantity: 1,
      }),
    );
  const [ingredientQuery, setIngredientQuery] = useState('');
  const [ingredientsForMeasurementUnit, setIngredientsForMeasurementUnit] = useState<
    QueryFilteredResult<ValidIngredientMeasurementUnit>
  >(pageLoadIngredientsForMeasurementUnit);
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
          return !(ingredientsForMeasurementUnit.data || []).some((vimu: ValidIngredientMeasurementUnit) => {
            return vimu.ingredient.id === mu.id;
          });
        });

        setSuggestedIngredients(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [ingredientQuery, ingredientsForMeasurementUnit.data]);

  const [newMeasurementUnitConversionFromMeasurementUnit, setNewMeasurementUnitConversionFromMeasurementUnit] =
    useState<ValidMeasurementUnitConversionCreationRequestInput>(
      new ValidMeasurementUnitConversionCreationRequestInput({
        from: validMeasurementUnit.id,
        modifier: 1,
      }),
    );
  const [conversionFromUnitQuery, setConversionFromUnitQuery] = useState('');
  const [measurementUnitsToConvertFrom, setMeasurementUnitsToConvertFrom] = useState<ValidMeasurementUnitConversion[]>(
    pageLoadMeasurementUnitConversionsFromUnit,
  );
  const [suggestedMeasurementUnitsToConvertFrom, setSuggestedMeasurementUnitsToConvertFrom] = useState<
    ValidMeasurementUnit[]
  >([]);
  const [conversionFromOnlyIngredientQuery, setConversionFromOnlyIngredientQuery] = useState('');
  const [suggestedIngredientsToRestrictConversionFrom, setSuggestedIngredientsToRestrictConversionFrom] = useState<
    ValidIngredient[]
  >([]);

  useEffect(() => {
    if (conversionFromUnitQuery.length <= 2) {
      setSuggestedMeasurementUnitsToConvertFrom([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidMeasurementUnits(conversionFromUnitQuery)
      .then((res: ValidMeasurementUnit[]) => {
        const newSuggestions = (res || []).filter((mu: ValidMeasurementUnit) => {
          return mu.id != validMeasurementUnit.id;
        });

        setSuggestedMeasurementUnitsToConvertFrom(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [conversionFromUnitQuery, validMeasurementUnit.id]);

  const [newMeasurementUnitConversionToMeasurementUnit, setNewMeasurementUnitConversionToMeasurementUnit] =
    useState<ValidMeasurementUnitConversionCreationRequestInput>(
      new ValidMeasurementUnitConversionCreationRequestInput({
        to: validMeasurementUnit.id,
        modifier: 1,
      }),
    );
  const [conversionToUnitQuery, setConversionToUnitQuery] = useState('');
  const [measurementUnitsToConvertTo, setMeasurementUnitsToConvertTo] = useState<ValidMeasurementUnitConversion[]>(
    pageLoadMeasurementUnitConversionsToUnit,
  );
  const [suggestedMeasurementUnitsToConvertTo, setSuggestedMeasurementUnitsToConvertTo] = useState<
    ValidMeasurementUnit[]
  >([]);
  const [conversionToOnlyIngredientQuery, setConversionToOnlyIngredientQuery] = useState('');
  const [suggestedIngredientsToRestrictConversionTo, setSuggestedIngredientsToRestrictConversionTo] = useState<
    ValidIngredient[]
  >([]);

  useEffect(() => {
    if (conversionToUnitQuery.length <= 2) {
      setSuggestedMeasurementUnitsToConvertTo([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidMeasurementUnits(conversionToUnitQuery)
      .then((res: ValidMeasurementUnit[]) => {
        const newSuggestions = (res || []).filter((mu: ValidMeasurementUnit) => {
          return mu.id != validMeasurementUnit.id;
        });

        setSuggestedMeasurementUnitsToConvertTo(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [conversionToUnitQuery, validMeasurementUnit.id]);

  useEffect(() => {
    if (conversionFromOnlyIngredientQuery.length <= 2) {
      setSuggestedIngredientsToRestrictConversionFrom([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidIngredients(conversionFromOnlyIngredientQuery)
      .then((res: QueryFilteredResult<ValidIngredient>) => {
        const newSuggestions = (res.data || []).filter((mu: ValidIngredient) => {
          return !(ingredientsForMeasurementUnit.data || []).some((vimu: ValidIngredientMeasurementUnit) => {
            return vimu.ingredient.id === mu.id;
          });
        });

        setSuggestedIngredientsToRestrictConversionFrom(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [conversionFromOnlyIngredientQuery, ingredientsForMeasurementUnit.data]);

  useEffect(() => {
    if (conversionToOnlyIngredientQuery.length <= 2) {
      setSuggestedIngredientsToRestrictConversionTo([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidIngredients(conversionToOnlyIngredientQuery)
      .then((res: QueryFilteredResult<ValidIngredient>) => {
        const newSuggestions = (res.data || []).filter((mu: ValidIngredient) => {
          return !(ingredientsForMeasurementUnit.data || []).some((vimu: ValidIngredientMeasurementUnit) => {
            return vimu.ingredient.id === mu.id;
          });
        });

        setSuggestedIngredientsToRestrictConversionTo(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [conversionToOnlyIngredientQuery, ingredientsForMeasurementUnit.data]);

  const updateForm = useForm({
    initialValues: validMeasurementUnit,
    validate: zodResolver(validMeasurementUnitUpdateFormSchema),
  });

  const dataHasChanged = (): boolean => {
    return (
      originalValidMeasurementUnit.name !== updateForm.values.name ||
      originalValidMeasurementUnit.description !== updateForm.values.description ||
      originalValidMeasurementUnit.pluralName !== updateForm.values.pluralName ||
      originalValidMeasurementUnit.universal !== updateForm.values.universal ||
      originalValidMeasurementUnit.volumetric !== updateForm.values.volumetric ||
      originalValidMeasurementUnit.metric !== updateForm.values.metric ||
      originalValidMeasurementUnit.imperial !== updateForm.values.imperial ||
      originalValidMeasurementUnit.slug !== updateForm.values.slug
    );
  };

  const submit = async () => {
    const validation = updateForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ValidMeasurementUnitUpdateRequestInput({
      name: updateForm.values.name,
      description: updateForm.values.description,
      universal: updateForm.values.universal,
      metric: updateForm.values.metric,
      volumetric: updateForm.values.volumetric,
      imperial: updateForm.values.imperial,
      pluralName: updateForm.values.pluralName,
      slug: updateForm.values.slug,
    });

    if (updateForm.values.imperial && updateForm.values.metric) {
      updateForm.setErrors({
        metric: 'Cannot be imperial and metric at the same time',
        imperial: 'Cannot be imperial and metric at the same time',
      });
      return;
    }

    const apiClient = buildLocalClient();

    await apiClient
      .updateValidMeasurementUnit(validMeasurementUnit.id, submission)
      .then((result: ValidMeasurementUnit) => {
        if (result) {
          updateForm.setValues(result);
          setValidMeasurementUnit(result);
          setOriginalValidMeasurementUnit(result);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Valid Measurement Unit">
      <Container size="sm">
        <form onSubmit={updateForm.onSubmit(submit)}>
          <TextInput label="Name" placeholder="thing" {...updateForm.getInputProps('name')} />
          <TextInput label="Plural Name" placeholder="things" {...updateForm.getInputProps('pluralName')} />
          <TextInput label="Slug" placeholder="thing" {...updateForm.getInputProps('slug')} />
          <TextInput label="Description" placeholder="thing" {...updateForm.getInputProps('description')} />

          <Switch
            checked={updateForm.values.volumetric}
            label="Volumetric"
            {...updateForm.getInputProps('volumetric')}
          />
          <Switch checked={updateForm.values.universal} label="Universal" {...updateForm.getInputProps('universal')} />
          <Switch
            checked={updateForm.values.metric}
            disabled={updateForm.values.imperial}
            label="Metric"
            {...updateForm.getInputProps('metric')}
          />
          <Switch
            checked={updateForm.values.imperial}
            disabled={updateForm.values.metric}
            label="Imperial"
            {...updateForm.getInputProps('imperial')}
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
                if (confirm('Are you sure you want to delete this valid measurement unit?')) {
                  apiClient.deleteValidMeasurementUnit(validMeasurementUnit.id).then(() => {
                    router.push('/valid_measurement_units');
                  });
                }
              }}
            >
              Delete
            </Button>
          </Group>
        </form>

        {!validMeasurementUnit.universal && (
          <>
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

              {ingredientsForMeasurementUnit.data && (ingredientsForMeasurementUnit.data || []).length !== 0 && (
                <>
                  <Table mt="xl" withColumnBorders>
                    <thead>
                      <tr>
                        <th>Name</th>
                        <th>Min Qty</th>
                        <th>Max Qty</th>
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
                      {(ingredientsForMeasurementUnit.data || []).map(
                        (validIngredientMeasurementUnit: ValidIngredientMeasurementUnit) => {
                          return (
                            <tr key={validIngredientMeasurementUnit.id}>
                              <td>
                                <Link href={`/valid_ingredients/${validIngredientMeasurementUnit.ingredient.id}`}>
                                  {validIngredientMeasurementUnit.ingredient.name}
                                </Link>
                              </td>
                              <td>
                                <Text>{validIngredientMeasurementUnit.minimumAllowableQuantity}</Text>
                              </td>
                              <td>
                                <Text>{validIngredientMeasurementUnit.maximumAllowableQuantity}</Text>
                              </td>
                              <td>
                                <Text>{validIngredientMeasurementUnit.notes}</Text>
                              </td>
                              <td>
                                <Center>
                                  <ActionIcon
                                    variant="outline"
                                    aria-label="remove valid ingredient measurement unit"
                                    onClick={async () => {
                                      await apiClient
                                        .deleteValidIngredientMeasurementUnit(validIngredientMeasurementUnit.id)
                                        .then(() => {
                                          setIngredientsForMeasurementUnit({
                                            ...ingredientsForMeasurementUnit,
                                            data: ingredientsForMeasurementUnit.data.filter(
                                              (x: ValidIngredientMeasurementUnit) =>
                                                x.id !== validIngredientMeasurementUnit.id,
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
                      Math.ceil(ingredientsForMeasurementUnit.totalCount / ingredientsForMeasurementUnit.limit) <=
                      ingredientsForMeasurementUnit.page
                    }
                    position="center"
                    page={ingredientsForMeasurementUnit.page}
                    total={Math.ceil(ingredientsForMeasurementUnit.totalCount / ingredientsForMeasurementUnit.limit)}
                    onChange={(value: number) => {
                      setIngredientsForMeasurementUnit({ ...ingredientsForMeasurementUnit, page: value });
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
                      const selectedIngredient = suggestedIngredients.find(
                        (x: ValidIngredient) => x.name === item.value,
                      );

                      if (!selectedIngredient) {
                        console.error(`selectedIngredient not found for item ${item.value}}`);
                        return;
                      }

                      setNewIngredientForMeasurementUnitInput({
                        ...newIngredientForMeasurementUnitInput,
                        validIngredientID: selectedIngredient.id,
                      });
                    }}
                    data={suggestedIngredients.map((x: ValidIngredient) => {
                      return { value: x.name, label: x.name };
                    })}
                  />
                </Grid.Col>
                <Grid.Col span={2}>
                  <NumberInput
                    value={newIngredientForMeasurementUnitInput.minimumAllowableQuantity}
                    label="Min. Qty"
                    onChange={(value: number) =>
                      setNewIngredientForMeasurementUnitInput({
                        ...newIngredientForMeasurementUnitInput,
                        minimumAllowableQuantity: value,
                      })
                    }
                  />
                </Grid.Col>
                <Grid.Col span={2}>
                  <NumberInput
                    value={newIngredientForMeasurementUnitInput.maximumAllowableQuantity}
                    label="Max. Qty"
                    onChange={(value: number) =>
                      setNewIngredientForMeasurementUnitInput({
                        ...newIngredientForMeasurementUnitInput,
                        maximumAllowableQuantity: value,
                      })
                    }
                  />
                </Grid.Col>
                <Grid.Col span="auto">
                  <TextInput
                    label="Notes"
                    value={newIngredientForMeasurementUnitInput.notes}
                    onChange={(event) =>
                      setNewIngredientForMeasurementUnitInput({
                        ...newIngredientForMeasurementUnitInput,
                        notes: event.target.value,
                      })
                    }
                  />
                </Grid.Col>
                <Grid.Col span={2}>
                  <Button
                    mt="xl"
                    disabled={
                      newIngredientForMeasurementUnitInput.validMeasurementUnitID === '' ||
                      newIngredientForMeasurementUnitInput.validIngredientID === ''
                    }
                    onClick={async () => {
                      await apiClient
                        .createValidIngredientMeasurementUnit(newIngredientForMeasurementUnitInput)
                        .then((res: ValidIngredientMeasurementUnit) => {
                          // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                          apiClient
                            .getValidIngredientMeasurementUnit(res.id)
                            .then((res: ValidIngredientMeasurementUnit) => {
                              setIngredientsForMeasurementUnit({
                                ...ingredientsForMeasurementUnit,
                                data: [...(ingredientsForMeasurementUnit.data || []), res],
                              });

                              setNewIngredientForMeasurementUnitInput(
                                new ValidIngredientMeasurementUnitCreationRequestInput({
                                  validMeasurementUnitID: validMeasurementUnit.id,
                                  minimumAllowableQuantity: 1,
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
          </>
        )}

        {/*

        CONVERSIONS FROM THIS MEASUREMENT UNIT

        */}

        <Space h="xl" />
        <Divider />
        <Space h="xl" />

        <form>
          <Center>
            <Title order={4}>Conversions From</Title>
          </Center>

          {measurementUnitsToConvertFrom && (measurementUnitsToConvertFrom || []).length !== 0 && (
            <>
              <Table mt="xl" withColumnBorders>
                <thead>
                  <tr>
                    <th>From</th>
                    <th>To</th>
                    <th>Modifier</th>
                    <th>Only For Ingredient</th>
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
                  {(measurementUnitsToConvertFrom || []).map(
                    (validMeasurementUnitConversion: ValidMeasurementUnitConversion) => {
                      return (
                        <tr key={validMeasurementUnitConversion.id}>
                          <td>
                            <Link href={`/valid_measurement_units/${validMeasurementUnitConversion.from.id}`}>
                              {validMeasurementUnitConversion.from.pluralName}
                            </Link>
                          </td>
                          <td>
                            <Link href={`/valid_measurement_units/${validMeasurementUnitConversion.to.id}`}>
                              {validMeasurementUnitConversion.to.pluralName}
                            </Link>
                          </td>
                          <td>
                            <Text>{validMeasurementUnitConversion.modifier}</Text>
                          </td>
                          <td>
                            {(validMeasurementUnitConversion.onlyForIngredient && (
                              <Link href={`/valid_ingredients/${validMeasurementUnitConversion.onlyForIngredient.id}`}>
                                {validMeasurementUnitConversion.to.pluralName}
                              </Link>
                            )) || <Text> - </Text>}
                          </td>
                          <td>
                            <Text>{validMeasurementUnitConversion.notes}</Text>
                          </td>
                          <td>
                            <Center>
                              <ActionIcon
                                variant="outline"
                                aria-label="remove valid ingredient measurement unit"
                                onClick={async () => {
                                  await apiClient
                                    .deleteValidMeasurementUnitConversion(validMeasurementUnitConversion.id)
                                    .then(() => {
                                      setMeasurementUnitsToConvertFrom([
                                        ...measurementUnitsToConvertFrom.filter(
                                          (x: ValidMeasurementUnitConversion) =>
                                            x.id !== validMeasurementUnitConversion.id,
                                        ),
                                      ]);
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

              {/*
              <Pagination
                disabled={
                  Math.ceil(measurementUnitsToConvertFrom.totalCount / measurementUnitsToConvertFrom.limit) <=
                  measurementUnitsToConvertFrom.page
                }
                position="center"
                page={measurementUnitsToConvertFrom.page}
                total={Math.ceil(measurementUnitsToConvertFrom.totalCount / measurementUnitsToConvertFrom.limit)}
                onChange={(value: number) => {
                  setMeasurementUnitsToConvertFrom({ ...measurementUnitsToConvertFrom, page: value });
                }}
              />
              */}
            </>
          )}

          <Grid>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="grams"
                label="Measurement Unit"
                value={conversionFromUnitQuery}
                onChange={setConversionFromUnitQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedMeasurementUnit = suggestedMeasurementUnitsToConvertFrom.find(
                    (x: ValidMeasurementUnit) => x.name === item.value,
                  );

                  if (!selectedMeasurementUnit) {
                    console.error(`selectedMeasurementUnit not found for item ${item.value}}`);
                    return;
                  }

                  setNewMeasurementUnitConversionFromMeasurementUnit({
                    ...newMeasurementUnitConversionFromMeasurementUnit,
                    to: selectedMeasurementUnit.id,
                  });
                }}
                data={suggestedMeasurementUnitsToConvertFrom.map((x: ValidMeasurementUnit) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <NumberInput
                label="Modifier"
                step={0.00001}
                precision={5}
                value={newMeasurementUnitConversionFromMeasurementUnit.modifier}
                onChange={(value: number) =>
                  setNewMeasurementUnitConversionFromMeasurementUnit({
                    ...newMeasurementUnitConversionFromMeasurementUnit,
                    modifier: value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="garlic"
                label="Only For Ingredient"
                value={conversionFromOnlyIngredientQuery}
                onChange={setConversionFromOnlyIngredientQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedValidIngredient = suggestedIngredientsToRestrictConversionFrom.find(
                    (x: ValidIngredient) => x.name === item.value,
                  );

                  if (!selectedValidIngredient) {
                    console.error(`selectedValidIngredient not found for item ${item.value}}`);
                    return;
                  }

                  setNewMeasurementUnitConversionFromMeasurementUnit({
                    ...newMeasurementUnitConversionFromMeasurementUnit,
                    onlyForIngredient: selectedValidIngredient.id,
                  });
                }}
                data={suggestedIngredientsToRestrictConversionFrom.map((x: ValidIngredient) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <TextInput
                label="Notes"
                value={newMeasurementUnitConversionFromMeasurementUnit.notes}
                onChange={(event) =>
                  setNewMeasurementUnitConversionFromMeasurementUnit({
                    ...newMeasurementUnitConversionFromMeasurementUnit,
                    notes: event.target.value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span={2}>
              <Button
                mt="xl"
                disabled={
                  newMeasurementUnitConversionFromMeasurementUnit.from === '' ||
                  newMeasurementUnitConversionFromMeasurementUnit.to === ''
                }
                onClick={async () => {
                  await apiClient
                    .createValidMeasurementUnitConversion(newMeasurementUnitConversionFromMeasurementUnit)
                    .then((res: ValidMeasurementUnitConversion) => {
                      // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                      apiClient
                        .getValidMeasurementUnitConversion(res.id)
                        .then((res: ValidMeasurementUnitConversion) => {
                          setMeasurementUnitsToConvertFrom([...(measurementUnitsToConvertFrom || []), res]);

                          setNewMeasurementUnitConversionFromMeasurementUnit(
                            new ValidMeasurementUnitConversionCreationRequestInput({
                              from: validMeasurementUnit.id,
                              to: '',
                              notes: '',
                              modifier: 1,
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

        {/*

        CONVERSIONS FROM THIS MEASUREMENT UNIT

        */}

        <Space h="xl" />
        <Divider />
        <Space h="xl" />

        <form>
          <Center>
            <Title order={4}>Conversions To</Title>
          </Center>

          {measurementUnitsToConvertTo && (measurementUnitsToConvertTo || []).length !== 0 && (
            <>
              <Table mt="xl" withColumnBorders>
                <thead>
                  <tr>
                    <th>From</th>
                    <th>To</th>
                    <th>Modifier</th>
                    <th>Only for Ingredient</th>
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
                  {(measurementUnitsToConvertTo || []).map(
                    (validMeasurementUnitConversion: ValidMeasurementUnitConversion) => {
                      return (
                        <tr key={validMeasurementUnitConversion.id}>
                          <td>
                            <Link href={`/valid_measurement_units/${validMeasurementUnitConversion.from.id}`}>
                              {validMeasurementUnitConversion.from.pluralName}
                            </Link>
                          </td>
                          <td>
                            <Link href={`/valid_measurement_units/${validMeasurementUnitConversion.to.id}`}>
                              {validMeasurementUnitConversion.to.pluralName}
                            </Link>
                          </td>
                          <td>
                            <Text>{validMeasurementUnitConversion.modifier}</Text>
                          </td>
                          <td>
                            {(validMeasurementUnitConversion.onlyForIngredient && (
                              <Link href={`/valid_ingredients/${validMeasurementUnitConversion.onlyForIngredient.id}`}>
                                {validMeasurementUnitConversion.to.pluralName}
                              </Link>
                            )) || <Text> - </Text>}
                          </td>
                          <td>
                            <Text>{validMeasurementUnitConversion.notes}</Text>
                          </td>
                          <td>
                            <Center>
                              <ActionIcon
                                variant="outline"
                                aria-label="remove valid ingredient measurement unit"
                                onClick={async () => {
                                  await apiClient
                                    .deleteValidMeasurementUnitConversion(validMeasurementUnitConversion.id)
                                    .then(() => {
                                      setMeasurementUnitsToConvertTo([
                                        ...measurementUnitsToConvertTo.filter(
                                          (x: ValidMeasurementUnitConversion) =>
                                            x.id !== validMeasurementUnitConversion.id,
                                        ),
                                      ]);
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

              {/*
              <Pagination
                disabled={
                  Math.ceil(measurementUnitsToConvertTo.totalCount / measurementUnitsToConvertTo.limit) <=
                  measurementUnitsToConvertTo.page
                }
                position="center"
                page={measurementUnitsToConvertTo.page}
                total={Math.ceil(measurementUnitsToConvertTo.totalCount / measurementUnitsToConvertTo.limit)}
                onChange={(value: number) => {
                  setMeasurementUnitsToConvertTo({ ...measurementUnitsToConvertTo, page: value });
                }}
              />
              */}
            </>
          )}

          <Grid>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="grams"
                label="Measurement Unit"
                value={conversionToUnitQuery}
                onChange={setConversionToUnitQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedMeasurementUnit = suggestedMeasurementUnitsToConvertTo.find(
                    (x: ValidMeasurementUnit) => x.name === item.value,
                  );

                  if (!selectedMeasurementUnit) {
                    console.error(`selectedMeasurementUnit not found for item ${item.value}}`);
                    return;
                  }

                  setNewMeasurementUnitConversionToMeasurementUnit({
                    ...newMeasurementUnitConversionToMeasurementUnit,
                    from: selectedMeasurementUnit.id,
                  });
                }}
                data={suggestedMeasurementUnitsToConvertTo.map((x: ValidMeasurementUnit) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <NumberInput
                label="Modifier"
                step={0.00001}
                precision={5}
                value={newMeasurementUnitConversionToMeasurementUnit.modifier}
                onChange={(value: number) =>
                  setNewMeasurementUnitConversionToMeasurementUnit({
                    ...newMeasurementUnitConversionToMeasurementUnit,
                    modifier: value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="garlic"
                label="Only For Ingredient"
                value={conversionToOnlyIngredientQuery}
                onChange={setConversionToOnlyIngredientQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedValidIngredient = suggestedIngredientsToRestrictConversionTo.find(
                    (x: ValidIngredient) => x.name === item.value,
                  );

                  if (!selectedValidIngredient) {
                    console.error(`selectedValidIngredient not found for item ${item.value}}`);
                    return;
                  }

                  setNewMeasurementUnitConversionToMeasurementUnit({
                    ...newMeasurementUnitConversionToMeasurementUnit,
                    onlyForIngredient: selectedValidIngredient.id,
                  });
                }}
                data={suggestedIngredientsToRestrictConversionTo.map((x: ValidIngredient) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <TextInput
                label="Notes"
                value={newMeasurementUnitConversionToMeasurementUnit.notes}
                onChange={(event) =>
                  setNewMeasurementUnitConversionToMeasurementUnit({
                    ...newMeasurementUnitConversionToMeasurementUnit,
                    notes: event.target.value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span={2}>
              <Button
                mt="xl"
                disabled={
                  newMeasurementUnitConversionToMeasurementUnit.from === '' ||
                  newMeasurementUnitConversionToMeasurementUnit.to === ''
                }
                onClick={async () => {
                  await apiClient
                    .createValidMeasurementUnitConversion(newMeasurementUnitConversionToMeasurementUnit)
                    .then((res: ValidMeasurementUnitConversion) => {
                      // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                      apiClient
                        .getValidMeasurementUnitConversion(res.id)
                        .then((res: ValidMeasurementUnitConversion) => {
                          setMeasurementUnitsToConvertTo([...(measurementUnitsToConvertTo || []), res]);

                          setNewMeasurementUnitConversionToMeasurementUnit(
                            new ValidMeasurementUnitConversionCreationRequestInput({
                              to: validMeasurementUnit.id,
                              from: '',
                              notes: '',
                              modifier: 1,
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

export default ValidMeasurementUnitPage;
