import { z } from 'zod';

export const ANALYTICS_ROUTE_ID = '/_auth/sites/$siteId/analytics' as const;

const PAGE_KEYS = [
  'eventsPage',
  'eventsCountsPage',
  'topPagesPage',
  'referrersPage',
  'devicesPage',
  'osPage',
  'countriesPage',
] as const;
const FILTER_KEYS = [
  'referrer',
  'browser',
  'device',
  'os',
  'page',
  'country',
  'eventName',
  'eventPath',
] as const;

export type AnalyticsPageKey = (typeof PAGE_KEYS)[number];
export type AnalyticsFilterKey = (typeof FILTER_KEYS)[number];

const optionalPage = z.preprocess((value) => {
  const raw = Array.isArray(value) ? value[0] : value;
  if (raw === undefined) return undefined;
  const page = Number(raw);
  return Number.isFinite(page) && page >= 1 ? Math.floor(page) : undefined;
}, z.number().int().min(1).optional());

const optionalFilter = z.preprocess((value) => {
  if (value === undefined || value === '') return undefined;
  const values = Array.isArray(value) ? value : [value];
  return values.filter((item): item is string => typeof item === 'string' && item !== '');
}, z.array(z.string()).optional());

const optionalSearchText = z.preprocess((value) => {
  const raw = Array.isArray(value) ? value[0] : value;
  if (typeof raw !== 'string') return undefined;
  const trimmed = raw.trim();
  return trimmed === '' ? undefined : trimmed;
}, z.string().optional());

const optionalDate = z
  .string()
  .regex(/^\d{4}-\d{2}-\d{2}$/v)
  .optional()
  .catch(undefined);
const optionalTime = z
  .string()
  .regex(/^\d{2}:\d{2}(?::\d{2})?$/v)
  .optional()
  .catch(undefined);

export const analyticsSearchSchema = z.object({
  preset: z.enum(['7d', '30d', '90d', 'custom', 'all']).optional().catch(undefined),
  from: optionalDate,
  to: optionalDate,
  fromTime: optionalTime,
  toTime: optionalTime,
  statsBucket: z.enum(['daily', 'hourly']).optional().catch(undefined),
  eventsPage: optionalPage,
  eventsCountsPage: optionalPage,
  topPagesPage: optionalPage,
  referrersPage: optionalPage,
  devicesPage: optionalPage,
  osPage: optionalPage,
  countriesPage: optionalPage,
  referrer: optionalFilter,
  browser: optionalFilter,
  device: optionalFilter,
  os: optionalFilter,
  page: optionalFilter,
  pagePathContains: optionalSearchText,
  country: optionalFilter,
  eventName: optionalFilter,
  eventPath: optionalFilter,
});

export type AnalyticsSearch = z.output<typeof analyticsSearchSchema>;

export function clearPagination(search: AnalyticsSearch): AnalyticsSearch {
  const next = { ...search };
  for (const key of PAGE_KEYS) delete next[key];
  return next;
}

export function setAnalyticsPage(
  search: AnalyticsSearch,
  key: AnalyticsPageKey,
  page: number
): AnalyticsSearch {
  const next = { ...search };
  if (page <= 1) delete next[key];
  else next[key] = page;
  return next;
}

export function addAnalyticsFilter(
  search: AnalyticsSearch,
  key: AnalyticsFilterKey,
  value: string
): AnalyticsSearch {
  const current = search[key] ?? [];
  return current.includes(value) ? search : { ...search, [key]: [...current, value] };
}

export function removeAnalyticsFilter(
  search: AnalyticsSearch,
  key: AnalyticsFilterKey,
  value: string
): AnalyticsSearch {
  const next = { ...search };
  const values = (search[key] ?? []).filter((item) => item !== value);
  if (values.length === 0) delete next[key];
  else next[key] = values;
  return next;
}

export function clearAnalyticsFilters(search: AnalyticsSearch): AnalyticsSearch {
  const next = { ...search };
  for (const key of FILTER_KEYS) delete next[key];
  delete next.pagePathContains;
  return next;
}

export function setPagePathSearch(search: AnalyticsSearch, value: string): AnalyticsSearch {
  const next = { ...search };
  const trimmed = value.trim();
  if (trimmed === '') delete next.pagePathContains;
  else next.pagePathContains = trimmed;
  return next;
}
