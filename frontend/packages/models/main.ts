import { Span } from '@opentelemetry/api';
import { IAPIError } from './errors';

export class ResponseDetails {
  currentHouseholdID: string;
  traceID: string;

  constructor(
    input: {
      currentHouseholdID?: string;
      traceID?: string;
    } = {},
  ) {
    this.currentHouseholdID = input.currentHouseholdID || 'unknown';
    this.traceID = input.traceID || 'unknown';
  }
}

export class Pagination {
  page: number;
  limit: number;
  filteredCount: number;
  totalCount: number;

  constructor(
    input: {
      page?: number;
      limit?: number;
      filteredCount?: number;
      totalCount?: number;
    } = {},
  ) {
    this.page = input.page || 1;
    this.limit = input.limit || 50;
    this.filteredCount = input.filteredCount || 0;
    this.totalCount = input.totalCount || 0;
  }
}

export class APIResponse<T> {
  data: T;
  pagination?: Pagination;
  error?: IAPIError;
  details: ResponseDetails;

  constructor(
    input: {
      data?: T;
      pagination?: Pagination;
      error?: IAPIError;
      details: ResponseDetails;
    } = {
      details: new ResponseDetails(),
    },
  ) {
    this.data = input.data || ({} as T);
    this.pagination = input.pagination;
    this.error = input.error;
    this.details = input.details;
  }
}

export class QueryFilteredResult<T> {
  data: T[];
  page: number;
  limit: number;
  filteredCount: number;
  totalCount: number;

  constructor(
    input: {
      data?: T[];
      page?: number;
      limit?: number;
      filteredCount?: number;
      totalCount?: number;
    } = {},
  ) {
    this.data = input.data || [];
    this.page = input.page || 1;
    this.limit = input.limit || 20;
    this.filteredCount = input.filteredCount || 0;
    this.totalCount = input.totalCount || 0;
  }
}

export type ValidSortType = 'asc' | 'desc';

export class QueryFilter {
  sortBy: ValidSortType = 'asc';
  page: number = 1;
  createdBefore?: number;
  createdAfter?: number;
  updatedBefore?: number;
  updatedAfter?: number;
  limit: number;
  includeArchived?: boolean;

  constructor(
    input: {
      sortBy?: ValidSortType;
      page?: number;
      createdBefore?: number;
      createdAfter?: number;
      updatedBefore?: number;
      updatedAfter?: number;
      limit?: number;
      includeArchived?: boolean;
    } = {},
  ) {
    this.sortBy = input.sortBy || 'asc';
    this.page = Math.max(input.page || 1, 1);
    this.createdBefore = input.createdBefore;
    this.createdAfter = input.createdAfter;
    this.updatedBefore = input.updatedBefore;
    this.updatedAfter = input.updatedAfter;
    this.limit = input.limit || 20;
    this.includeArchived = Boolean(input.includeArchived);
  }

  public asURLSearchParams(): URLSearchParams {
    const out = new URLSearchParams();

    if (this.sortBy) out.set('sortBy', this.sortBy);
    if (this.page) out.set('page', this.page.toString());
    if (this.createdBefore) out.set('createdBefore', this.createdBefore.toString());
    if (this.createdAfter) out.set('createdAfter', this.createdAfter.toString());
    if (this.updatedBefore) out.set('updatedBefore', this.updatedBefore.toString());
    if (this.updatedAfter) out.set('updatedAfter', this.updatedAfter.toString());
    if (this.limit) out.set('limit', this.limit.toString());
    if (this.includeArchived) out.set('includeArchived', this.includeArchived.toString());

    return out;
  }

  public asRecord(): Record<string, string | number | ValidSortType> {
    const params = {} as Record<string, string | number | ValidSortType>;

    if (this.sortBy) params['sortBy'] = this.sortBy;
    if (this.page) params['page'] = this.page.toString();
    if (this.createdBefore) params['createdBefore'] = this.createdBefore.toString();
    if (this.createdAfter) params['createdAfter'] = this.createdAfter.toString();
    if (this.updatedBefore) params['updatedBefore'] = this.updatedBefore.toString();
    if (this.updatedAfter) params['updatedAfter'] = this.updatedAfter.toString();
    if (this.limit) params['limit'] = this.limit.toString();
    if (this.includeArchived) params['includeArchived'] = this.includeArchived.toString();

    return params;
  }

  public attachToSpan(span: Span): void {
    span.setAttributes({
      'pagination.sortBy': this.sortBy,
      'pagination.page': this.page,
      'pagination.createdBefore': this.createdBefore,
      'pagination.createdAfter': this.createdAfter,
      'pagination.updatedBefore': this.updatedBefore,
      'pagination.updatedAfter': this.updatedAfter,
      'pagination.limit': this.limit,
      'pagination.includeArchived': this.includeArchived,
    });
  }

  public static deriveFromPage(): QueryFilter {
    const qf = QueryFilter.Default();
    const pageParams = new URLSearchParams(window.location.search);

    if (pageParams.has('sortBy')) qf.sortBy = pageParams.get('sortBy') as 'asc' | 'desc';
    if (pageParams.has('page')) qf.page = parseInt(pageParams.get('page')!);
    if (pageParams.has('createdBefore')) qf.createdBefore = parseInt(pageParams.get('createdBefore')!);
    if (pageParams.has('createdAfter')) qf.createdAfter = parseInt(pageParams.get('createdAfter')!);
    if (pageParams.has('updatedBefore')) qf.updatedBefore = parseInt(pageParams.get('updatedBefore')!);
    if (pageParams.has('updatedAfter')) qf.updatedAfter = parseInt(pageParams.get('updatedAfter')!);
    if (pageParams.has('limit')) qf.limit = parseInt(pageParams.get('limit')!);
    if (pageParams.has('includeArchived')) qf.includeArchived = pageParams.get('includeArchived') === 'true';

    return qf;
  }

  public static deriveFromGetServerSidePropsContext(x: { [key: string]: string | string[] | undefined }): QueryFilter {
    const qf = QueryFilter.Default();

    if (x['sortBy']) qf.sortBy = x['sortBy'] as 'asc' | 'desc';
    if (x['page']) qf.page = parseInt(x['page'].toString()!);
    if (x['createdBefore']) qf.createdBefore = parseInt(x['createdBefore'].toString()!);
    if (x['createdAfter']) qf.createdAfter = parseInt(x['createdAfter'].toString()!);
    if (x['updatedBefore']) qf.updatedBefore = parseInt(x['updatedBefore'].toString()!);
    if (x['updatedAfter']) qf.updatedAfter = parseInt(x['updatedAfter'].toString()!);
    if (x['limit']) qf.limit = parseInt(x['limit'].toString()!);
    if (x['includeArchived']) qf.includeArchived = x['includeArchived'] === 'true';

    return qf;
  }

  public static Default(): QueryFilter {
    return new QueryFilter({
      page: 1,
      limit: 20,
    });
  }
}
