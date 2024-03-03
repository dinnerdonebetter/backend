import { QueryFilteredResult } from './main';
import { ValidIngredient } from './validIngredients';
import { ValidMeasurementUnit } from './validMeasurementUnits';

export class ValidMeasurementConversion {
  createdAt: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  onlyForIngredient?: ValidIngredient;
  notes: string;
  id: string;
  from: ValidMeasurementUnit;
  to: ValidMeasurementUnit;
  modifier: number;

  constructor(
    input: {
      createdAt?: string;
      lastUpdatedAt?: string;
      archivedAt?: string;
      onlyForIngredient?: ValidIngredient;
      notes?: string;
      id?: string;
      from?: ValidMeasurementUnit;
      to?: ValidMeasurementUnit;
      modifier?: number;
    } = {},
  ) {
    this.createdAt = input.createdAt || '1970-01-01T00:00:00Z';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
    this.onlyForIngredient = input.onlyForIngredient;
    this.notes = input.notes || '';
    this.id = input.id || '';
    this.from = input.from || new ValidMeasurementUnit();
    this.to = input.to || new ValidMeasurementUnit();
    this.modifier = input.modifier || 0;
  }
}

export class ValidMeasurementConversionList extends QueryFilteredResult<ValidMeasurementConversion> {
  constructor(
    input: {
      data?: ValidMeasurementConversion[];
      page?: number;
      limit?: number;
      filteredCount?: number;
      totalCount?: number;
    } = {},
  ) {
    super(input);

    this.data = input.data || [];
    this.page = input.page || 1;
    this.limit = input.limit || 20;
    this.filteredCount = input.filteredCount || 0;
    this.totalCount = input.totalCount || 0;
  }
}

export class ValidMeasurementConversionCreationRequestInput {
  onlyForIngredient?: string;
  from: string;
  to: string;
  notes: string;
  modifier: number;

  constructor(
    input: {
      onlyForIngredient?: string;
      from?: string;
      to?: string;
      notes?: string;
      modifier?: number;
    } = {},
  ) {
    this.onlyForIngredient = input.onlyForIngredient;
    this.notes = input.notes || '';
    this.from = input.from || '';
    this.to = input.to || '';
    this.modifier = input.modifier || 0;
  }

  static fromValidMeasurementConversion(
    input: ValidMeasurementConversion,
  ): ValidMeasurementConversionCreationRequestInput {
    const x = new ValidMeasurementConversionCreationRequestInput();

    x.onlyForIngredient = input.onlyForIngredient?.id;
    x.notes = input.notes;
    x.from = input.from.id;
    x.to = input.to.id;
    x.modifier = input.modifier;

    return x;
  }
}

export class ValidMeasurementConversionUpdateRequestInput {
  onlyForIngredient?: string;
  from: string;
  to: string;
  notes: string;
  modifier: number;

  constructor(
    input: {
      onlyForIngredient?: string;
      from?: string;
      to?: string;
      notes?: string;
      modifier?: number;
    } = {},
  ) {
    this.onlyForIngredient = input.onlyForIngredient;
    this.notes = input.notes || '';
    this.from = input.from || '';
    this.to = input.to || '';
    this.modifier = input.modifier || 0;
  }

  static fromValidMeasurementConversion(
    input: ValidMeasurementConversion,
  ): ValidMeasurementConversionUpdateRequestInput {
    const x = new ValidMeasurementConversionUpdateRequestInput();

    x.onlyForIngredient = input.onlyForIngredient?.id;
    x.notes = input.notes;
    x.from = input.from.id;
    x.to = input.to.id;
    x.modifier = input.modifier;

    return x;
  }
}
