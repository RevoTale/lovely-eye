/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = T | null | undefined;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Time: { input: string; output: string; }
};

export type ActivePageStats = {
  __typename: 'ActivePageStats';
  path: Scalars['String']['output'];
  /** Number of visitors currently viewing this page */
  visitors: Scalars['Int']['output'];
};

export type AuthPayload = {
  __typename: 'AuthPayload';
  user: User;
};

export type BrowserStats = {
  __typename: 'BrowserStats';
  browser: Scalars['String']['output'];
  visitors: Scalars['Int']['output'];
};

export type CountryStats = {
  __typename: 'CountryStats';
  country: Scalars['String']['output'];
  visitors: Scalars['Int']['output'];
};

export type CreateSiteInput = {
  domains: Array<Scalars['String']['input']>;
  name: Scalars['String']['input'];
};

export type DailyStats = {
  __typename: 'DailyStats';
  date: Scalars['Time']['output'];
  pageViews: Scalars['Int']['output'];
  sessions: Scalars['Int']['output'];
  visitors: Scalars['Int']['output'];
};

export type DashboardStats = {
  __typename: 'DashboardStats';
  /** Average session duration in seconds */
  avgDuration: Scalars['Float']['output'];
  bounceRate: Scalars['Float']['output'];
  browsers: Array<BrowserStats>;
  countries: PagedCountryStats;
  dailyStats: Array<DailyStats>;
  devices: PagedDeviceStats;
  pageViews: Scalars['Int']['output'];
  sessions: Scalars['Int']['output'];
  topPages: PagedPageStats;
  topReferrers: PagedReferrerStats;
  visitors: Scalars['Int']['output'];
};


export type DashboardStatsBrowsersArgs = {
  paging: PagingInput;
};


export type DashboardStatsCountriesArgs = {
  paging: PagingInput;
};


export type DashboardStatsDailyStatsArgs = {
  bucket?: InputMaybe<TimeBucket>;
  limit: InputMaybe<Scalars['Int']['input']>;
  offset: InputMaybe<Scalars['Int']['input']>;
};


export type DashboardStatsDevicesArgs = {
  paging: PagingInput;
};


export type DashboardStatsTopPagesArgs = {
  paging: PagingInput;
};


export type DashboardStatsTopReferrersArgs = {
  paging: PagingInput;
};

export type DateRangeInput = {
  from: InputMaybe<Scalars['Time']['input']>;
  to: InputMaybe<Scalars['Time']['input']>;
};

export type DeviceStats = {
  __typename: 'DeviceStats';
  device: Scalars['String']['output'];
  visitors: Scalars['Int']['output'];
};

export type Event = {
  __typename: 'Event';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  path: Scalars['String']['output'];
  /** Key-value properties associated with the event */
  properties: Array<EventProperty>;
};

export type EventCount = {
  __typename: 'EventCount';
  count: Scalars['Int']['output'];
  event: Event;
};

export type EventDefinition = {
  __typename: 'EventDefinition';
  createdAt: Scalars['Time']['output'];
  fields: Array<EventDefinitionField>;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type EventDefinitionField = {
  __typename: 'EventDefinitionField';
  id: Scalars['ID']['output'];
  key: Scalars['String']['output'];
  maxLength: Scalars['Int']['output'];
  required: Scalars['Boolean']['output'];
  type: EventFieldType;
};

export type EventDefinitionFieldInput = {
  key: Scalars['String']['input'];
  maxLength: InputMaybe<Scalars['Int']['input']>;
  required: Scalars['Boolean']['input'];
  type: EventFieldType;
};

export type EventDefinitionInput = {
  fields: Array<EventDefinitionFieldInput>;
  name: Scalars['String']['input'];
};

export const EventFieldType = {
  Boolean: 'BOOLEAN',
  Int: 'INT',
  String: 'STRING'
} as const;

export type EventFieldType = typeof EventFieldType[keyof typeof EventFieldType];
export type EventProperty = {
  __typename: 'EventProperty';
  key: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type EventsResult = {
  __typename: 'EventsResult';
  events: Array<Event>;
  total: Scalars['Int']['output'];
};

export type FilterInput = {
  /** Filter by country (stored country name) */
  country: InputMaybe<Array<Scalars['String']['input']>>;
  /** Filter by device type (desktop, mobile, tablet) */
  device: InputMaybe<Array<Scalars['String']['input']>>;
  /** Filter by event name */
  eventName: InputMaybe<Array<Scalars['String']['input']>>;
  /** Filter by event path */
  eventPath: InputMaybe<Array<Scalars['String']['input']>>;
  /** Filter by page path */
  page: InputMaybe<Array<Scalars['String']['input']>>;
  /** Filter by specific referrer */
  referrer: InputMaybe<Array<Scalars['String']['input']>>;
};

export type GeoIpCountry = {
  __typename: 'GeoIPCountry';
  code: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type GeoIpStatus = {
  __typename: 'GeoIPStatus';
  dbPath: Scalars['String']['output'];
  lastError: Maybe<Scalars['String']['output']>;
  source: Maybe<Scalars['String']['output']>;
  state: Scalars['String']['output'];
  updatedAt: Maybe<Scalars['Time']['output']>;
};

export type LoginInput = {
  password: Scalars['String']['input'];
  username: Scalars['String']['input'];
};

export type Mutation = {
  __typename: 'Mutation';
  createSite: Site;
  deleteEventDefinition: Scalars['Boolean']['output'];
  /** Deletes site and all analytics data */
  deleteSite: Scalars['Boolean']['output'];
  login: AuthPayload;
  /** Clears auth cookies */
  logout: Scalars['Boolean']['output'];
  refreshGeoIPDatabase: GeoIpStatus;
  refreshToken: TokenPayload;
  /** Invalidates old tracking scripts */
  regenerateSiteKey: Site;
  /** First user becomes admin */
  register: AuthPayload;
  updateSite: Site;
  upsertEventDefinition: EventDefinition;
};


export type MutationCreateSiteArgs = {
  input: CreateSiteInput;
};


export type MutationDeleteEventDefinitionArgs = {
  name: Scalars['String']['input'];
  siteId: Scalars['ID']['input'];
};


export type MutationDeleteSiteArgs = {
  id: Scalars['ID']['input'];
};


export type MutationLoginArgs = {
  input: LoginInput;
};


export type MutationRefreshTokenArgs = {
  refreshToken: Scalars['String']['input'];
};


export type MutationRegenerateSiteKeyArgs = {
  id: Scalars['ID']['input'];
};


export type MutationRegisterArgs = {
  input: RegisterInput;
};


export type MutationUpdateSiteArgs = {
  id: Scalars['ID']['input'];
  input: UpdateSiteInput;
};


export type MutationUpsertEventDefinitionArgs = {
  input: EventDefinitionInput;
  siteId: Scalars['ID']['input'];
};

export type PageStats = {
  __typename: 'PageStats';
  path: Scalars['String']['output'];
  views: Scalars['Int']['output'];
  visitors: Scalars['Int']['output'];
};

export type PagedCountryStats = {
  __typename: 'PagedCountryStats';
  items: Array<CountryStats>;
  total: Scalars['Int']['output'];
  totalVisitors: Scalars['Int']['output'];
};

export type PagedDeviceStats = {
  __typename: 'PagedDeviceStats';
  items: Array<DeviceStats>;
  total: Scalars['Int']['output'];
  totalVisitors: Scalars['Int']['output'];
};

export type PagedPageStats = {
  __typename: 'PagedPageStats';
  items: Array<PageStats>;
  total: Scalars['Int']['output'];
};

export type PagedReferrerStats = {
  __typename: 'PagedReferrerStats';
  items: Array<ReferrerStats>;
  total: Scalars['Int']['output'];
};

export type PagingInput = {
  limit: Scalars['Int']['input'];
  offset: Scalars['Int']['input'];
};

export type Query = {
  __typename: 'Query';
  dashboard: DashboardStats;
  /** Get event counts aggregated by event with the most recent occurrence */
  eventCounts: Array<EventCount>;
  /** Get event definitions for a site */
  eventDefinitions: Array<EventDefinition>;
  /** Get events for a site with pagination */
  events: EventsResult;
  geoIPCountries: Array<GeoIpCountry>;
  geoIPStatus: GeoIpStatus;
  me: Maybe<User>;
  realtime: RealtimeStats;
  site: Maybe<Site>;
  sites: Array<Site>;
};


export type QueryDashboardArgs = {
  dateRange: InputMaybe<DateRangeInput>;
  filter: InputMaybe<FilterInput>;
  siteId: Scalars['ID']['input'];
};


export type QueryEventCountsArgs = {
  dateRange: InputMaybe<DateRangeInput>;
  filter: InputMaybe<FilterInput>;
  paging: PagingInput;
  siteId: Scalars['ID']['input'];
};


export type QueryEventDefinitionsArgs = {
  paging: PagingInput;
  siteId: Scalars['ID']['input'];
};


export type QueryEventsArgs = {
  dateRange: InputMaybe<DateRangeInput>;
  filter: InputMaybe<FilterInput>;
  limit: InputMaybe<Scalars['Int']['input']>;
  offset: InputMaybe<Scalars['Int']['input']>;
  siteId: Scalars['ID']['input'];
};


export type QueryGeoIpCountriesArgs = {
  paging: PagingInput;
  search: InputMaybe<Scalars['String']['input']>;
};


export type QueryRealtimeArgs = {
  siteId: Scalars['ID']['input'];
};


export type QuerySiteArgs = {
  id: Scalars['ID']['input'];
};


export type QuerySitesArgs = {
  paging: PagingInput;
};

export type RealtimeStats = {
  __typename: 'RealtimeStats';
  /** Active pages with visitor count */
  activePages: Array<ActivePageStats>;
  /** Visitors active in last 5 minutes */
  visitors: Scalars['Int']['output'];
};


export type RealtimeStatsActivePagesArgs = {
  paging: PagingInput;
};

export type ReferrerStats = {
  __typename: 'ReferrerStats';
  referrer: Scalars['String']['output'];
  visitors: Scalars['Int']['output'];
};

export type RegisterInput = {
  password: Scalars['String']['input'];
  username: Scalars['String']['input'];
};

export type Session = {
  __typename: 'Session';
  /** True when created from an event without a page view; flipped to false after a page view arrives. */
  eventOnly: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
};

export type Site = {
  __typename: 'Site';
  /** ISO country codes blocked from tracking */
  blockedCountries: Array<Scalars['String']['output']>;
  /** IP addresses blocked from tracking */
  blockedIPs: Array<Scalars['String']['output']>;
  createdAt: Scalars['Time']['output'];
  /** All tracked domains (includes primary) */
  domains: Array<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  /** Used in tracking script */
  publicKey: Scalars['String']['output'];
  /** Enable country tracking (requires GeoIP database) */
  trackCountry: Scalars['Boolean']['output'];
};

export const TimeBucket = {
  Daily: 'DAILY',
  Hourly: 'HOURLY'
} as const;

export type TimeBucket = typeof TimeBucket[keyof typeof TimeBucket];
export type TokenPayload = {
  __typename: 'TokenPayload';
  accessToken: Scalars['String']['output'];
  refreshToken: Scalars['String']['output'];
};

export type UpdateSiteInput = {
  /** Full list of blocked country codes */
  blockedCountries: InputMaybe<Array<Scalars['String']['input']>>;
  /** Full list of blocked IPs */
  blockedIPs: InputMaybe<Array<Scalars['String']['input']>>;
  /** Full list of tracked domains (includes primary) */
  domains: InputMaybe<Array<Scalars['String']['input']>>;
  name: Scalars['String']['input'];
  trackCountry: InputMaybe<Scalars['Boolean']['input']>;
};

export type User = {
  __typename: 'User';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  role: Scalars['String']['output'];
  sites: Maybe<Array<Site>>;
  username: Scalars['String']['output'];
};


export type UserSitesArgs = {
  paging: PagingInput;
};

export type RefreshGeoIpDatabaseMutationVariables = Exact<{ [key: string]: never; }>;


export type RefreshGeoIpDatabaseMutation = { __typename: 'Mutation', refreshGeoIPDatabase: (
    { __typename: 'GeoIPStatus' }
    & { ' $fragmentRefs'?: { 'CountryTrackingGeoIpStatusFieldsFragment': CountryTrackingGeoIpStatusFieldsFragment } }
  ) };

export type CountryTrackingGeoIpStatusFieldsFragment = { __typename: 'GeoIPStatus', state: string, dbPath: string, source: string | null, lastError: string | null, updatedAt: string | null } & { ' $fragmentName'?: 'CountryTrackingGeoIpStatusFieldsFragment' };

export type DeleteSiteMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteSiteMutation = { __typename: 'Mutation', deleteSite: boolean };

export type EventDefinitionsQueryVariables = Exact<{
  siteId: Scalars['ID']['input'];
  paging: PagingInput;
}>;


export type EventDefinitionsQuery = { __typename: 'Query', eventDefinitions: Array<(
    { __typename: 'EventDefinition' }
    & { ' $fragmentRefs'?: { 'EventDefinitionFieldsFragment': EventDefinitionFieldsFragment } }
  )> };

export type UpsertEventDefinitionMutationVariables = Exact<{
  siteId: Scalars['ID']['input'];
  input: EventDefinitionInput;
}>;


export type UpsertEventDefinitionMutation = { __typename: 'Mutation', upsertEventDefinition: (
    { __typename: 'EventDefinition' }
    & { ' $fragmentRefs'?: { 'EventDefinitionFieldsFragment': EventDefinitionFieldsFragment } }
  ) };

export type DeleteEventDefinitionMutationVariables = Exact<{
  siteId: Scalars['ID']['input'];
  name: Scalars['String']['input'];
}>;


export type DeleteEventDefinitionMutation = { __typename: 'Mutation', deleteEventDefinition: boolean };

export type EventDefinitionFieldFieldsFragment = { __typename: 'EventDefinitionField', id: string, key: string, type: EventFieldType, required: boolean, maxLength: number } & { ' $fragmentName'?: 'EventDefinitionFieldFieldsFragment' };

export type EventDefinitionFieldsFragment = { __typename: 'EventDefinition', id: string, name: string, createdAt: string, updatedAt: string, fields: Array<(
    { __typename: 'EventDefinitionField' }
    & { ' $fragmentRefs'?: { 'EventDefinitionFieldFieldsFragment': EventDefinitionFieldFieldsFragment } }
  )> } & { ' $fragmentName'?: 'EventDefinitionFieldsFragment' };

export type RegenerateSiteKeyMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type RegenerateSiteKeyMutation = { __typename: 'Mutation', regenerateSiteKey: (
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'TrackingCodeSiteSummaryFieldsFragment': TrackingCodeSiteSummaryFieldsFragment } }
  ) };

export type TrackingCodeSiteSummaryFieldsFragment = { __typename: 'Site', id: string, domains: Array<string>, name: string, publicKey: string, createdAt: string } & { ' $fragmentName'?: 'TrackingCodeSiteSummaryFieldsFragment' };

export type GeoIpCountriesQueryVariables = Exact<{
  search: InputMaybe<Scalars['String']['input']>;
  paging: PagingInput;
}>;


export type GeoIpCountriesQuery = { __typename: 'Query', geoIPCountries: Array<(
    { __typename: 'GeoIPCountry' }
    & { ' $fragmentRefs'?: { 'TrafficBlockingCountryFieldsFragment': TrafficBlockingCountryFieldsFragment } }
  )> };

export type TrafficBlockingCountryFieldsFragment = { __typename: 'GeoIPCountry', code: string, name: string } & { ' $fragmentName'?: 'TrafficBlockingCountryFieldsFragment' };

export type MeQueryVariables = Exact<{
  sitesPaging: PagingInput;
}>;


export type MeQuery = { __typename: 'Query', me: (
    { __typename: 'User' }
    & { ' $fragmentRefs'?: { 'AuthUserDetailsFieldsFragment': AuthUserDetailsFieldsFragment } }
  ) | null };

export type LoginMutationVariables = Exact<{
  input: LoginInput;
}>;


export type LoginMutation = { __typename: 'Mutation', login: { __typename: 'AuthPayload', user: (
      { __typename: 'User' }
      & { ' $fragmentRefs'?: { 'AuthUserFieldsFragment': AuthUserFieldsFragment } }
    ) } };

export type RegisterMutationVariables = Exact<{
  input: RegisterInput;
}>;


export type RegisterMutation = { __typename: 'Mutation', register: { __typename: 'AuthPayload', user: (
      { __typename: 'User' }
      & { ' $fragmentRefs'?: { 'AuthUserFieldsFragment': AuthUserFieldsFragment } }
    ) } };

export type LogoutMutationVariables = Exact<{ [key: string]: never; }>;


export type LogoutMutation = { __typename: 'Mutation', logout: boolean };

export type AuthUserFieldsFragment = { __typename: 'User', id: string, username: string, role: string } & { ' $fragmentName'?: 'AuthUserFieldsFragment' };

export type AuthSiteSummaryFieldsFragment = { __typename: 'Site', id: string, domains: Array<string>, name: string, publicKey: string, createdAt: string } & { ' $fragmentName'?: 'AuthSiteSummaryFieldsFragment' };

export type AuthUserDetailsFieldsFragment = { __typename: 'User', id: string, username: string, role: string, createdAt: string, sites: Array<(
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'AuthSiteSummaryFieldsFragment': AuthSiteSummaryFieldsFragment } }
  )> | null } & { ' $fragmentName'?: 'AuthUserDetailsFieldsFragment' };

export type DashboardQueryVariables = Exact<{
  siteId: Scalars['ID']['input'];
  dateRange: InputMaybe<DateRangeInput>;
  filter: InputMaybe<FilterInput>;
  topPagesPaging: PagingInput;
  referrersPaging: PagingInput;
  browsersPaging: PagingInput;
  devicesPaging: PagingInput;
  countriesPaging: PagingInput;
}>;


export type DashboardQuery = { __typename: 'Query', dashboard: (
    { __typename: 'DashboardStats' }
    & { ' $fragmentRefs'?: { 'DashboardStatsFieldsFragment': DashboardStatsFieldsFragment } }
  ) };

export type RealtimeQueryVariables = Exact<{
  siteId: Scalars['ID']['input'];
  activePagesPaging: PagingInput;
}>;


export type RealtimeQuery = { __typename: 'Query', realtime: (
    { __typename: 'RealtimeStats' }
    & { ' $fragmentRefs'?: { 'RealtimeStatsFieldsFragment': RealtimeStatsFieldsFragment } }
  ) };

export type EventsQueryVariables = Exact<{
  siteId: Scalars['ID']['input'];
  dateRange: InputMaybe<DateRangeInput>;
  filter: InputMaybe<FilterInput>;
  limit: InputMaybe<Scalars['Int']['input']>;
  offset: InputMaybe<Scalars['Int']['input']>;
}>;


export type EventsQuery = { __typename: 'Query', events: { __typename: 'EventsResult', total: number, events: Array<(
      { __typename: 'Event' }
      & { ' $fragmentRefs'?: { 'EventFieldsFragment': EventFieldsFragment } }
    )> } };

export type EventCountsQueryVariables = Exact<{
  siteId: Scalars['ID']['input'];
  dateRange: InputMaybe<DateRangeInput>;
  filter: InputMaybe<FilterInput>;
  paging: PagingInput;
}>;


export type EventCountsQuery = { __typename: 'Query', eventCounts: Array<(
    { __typename: 'EventCount' }
    & { ' $fragmentRefs'?: { 'EventCountFieldsFragment': EventCountFieldsFragment } }
  )> };

export type ChartDataQueryVariables = Exact<{
  siteId: Scalars['ID']['input'];
  dateRange: InputMaybe<DateRangeInput>;
  filter: InputMaybe<FilterInput>;
  bucket: InputMaybe<TimeBucket>;
  limit: InputMaybe<Scalars['Int']['input']>;
  offset: InputMaybe<Scalars['Int']['input']>;
}>;


export type ChartDataQuery = { __typename: 'Query', dashboard: { __typename: 'DashboardStats', dailyStats: Array<(
      { __typename: 'DailyStats' }
      & { ' $fragmentRefs'?: { 'DailyStatsFieldsFragment': DailyStatsFieldsFragment } }
    )> } };

export type PageStatsFieldsFragment = { __typename: 'PageStats', path: string, views: number, visitors: number } & { ' $fragmentName'?: 'PageStatsFieldsFragment' };

export type ReferrerStatsFieldsFragment = { __typename: 'ReferrerStats', referrer: string, visitors: number } & { ' $fragmentName'?: 'ReferrerStatsFieldsFragment' };

export type BrowserStatsFieldsFragment = { __typename: 'BrowserStats', browser: string, visitors: number } & { ' $fragmentName'?: 'BrowserStatsFieldsFragment' };

export type DeviceStatsFieldsFragment = { __typename: 'DeviceStats', device: string, visitors: number } & { ' $fragmentName'?: 'DeviceStatsFieldsFragment' };

export type CountryStatsFieldsFragment = { __typename: 'CountryStats', country: string, visitors: number } & { ' $fragmentName'?: 'CountryStatsFieldsFragment' };

export type ActivePageStatsFieldsFragment = { __typename: 'ActivePageStats', path: string, visitors: number } & { ' $fragmentName'?: 'ActivePageStatsFieldsFragment' };

export type DailyStatsFieldsFragment = { __typename: 'DailyStats', date: string, visitors: number, pageViews: number, sessions: number } & { ' $fragmentName'?: 'DailyStatsFieldsFragment' };

export type EventPropertyFieldsFragment = { __typename: 'EventProperty', key: string, value: string } & { ' $fragmentName'?: 'EventPropertyFieldsFragment' };

export type EventFieldsFragment = { __typename: 'Event', id: string, name: string, path: string, createdAt: string, properties: Array<(
    { __typename: 'EventProperty' }
    & { ' $fragmentRefs'?: { 'EventPropertyFieldsFragment': EventPropertyFieldsFragment } }
  )> } & { ' $fragmentName'?: 'EventFieldsFragment' };

export type EventCountFieldsFragment = { __typename: 'EventCount', count: number, event: (
    { __typename: 'Event' }
    & { ' $fragmentRefs'?: { 'EventFieldsFragment': EventFieldsFragment } }
  ) } & { ' $fragmentName'?: 'EventCountFieldsFragment' };

export type RealtimeStatsFieldsFragment = { __typename: 'RealtimeStats', visitors: number, activePages: Array<(
    { __typename: 'ActivePageStats' }
    & { ' $fragmentRefs'?: { 'ActivePageStatsFieldsFragment': ActivePageStatsFieldsFragment } }
  )> } & { ' $fragmentName'?: 'RealtimeStatsFieldsFragment' };

export type DashboardStatsFieldsFragment = { __typename: 'DashboardStats', visitors: number, pageViews: number, sessions: number, bounceRate: number, avgDuration: number, topPages: { __typename: 'PagedPageStats', total: number, items: Array<(
      { __typename: 'PageStats' }
      & { ' $fragmentRefs'?: { 'PageStatsFieldsFragment': PageStatsFieldsFragment } }
    )> }, topReferrers: { __typename: 'PagedReferrerStats', total: number, items: Array<(
      { __typename: 'ReferrerStats' }
      & { ' $fragmentRefs'?: { 'ReferrerStatsFieldsFragment': ReferrerStatsFieldsFragment } }
    )> }, browsers: Array<(
    { __typename: 'BrowserStats' }
    & { ' $fragmentRefs'?: { 'BrowserStatsFieldsFragment': BrowserStatsFieldsFragment } }
  )>, devices: { __typename: 'PagedDeviceStats', total: number, totalVisitors: number, items: Array<(
      { __typename: 'DeviceStats' }
      & { ' $fragmentRefs'?: { 'DeviceStatsFieldsFragment': DeviceStatsFieldsFragment } }
    )> }, countries: { __typename: 'PagedCountryStats', total: number, totalVisitors: number, items: Array<(
      { __typename: 'CountryStats' }
      & { ' $fragmentRefs'?: { 'CountryStatsFieldsFragment': CountryStatsFieldsFragment } }
    )> } } & { ' $fragmentName'?: 'DashboardStatsFieldsFragment' };

export type CreateSiteMutationVariables = Exact<{
  input: CreateSiteInput;
}>;


export type CreateSiteMutation = { __typename: 'Mutation', createSite: (
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'SiteDetailsFieldsFragment': SiteDetailsFieldsFragment } }
  ) };

export type GeoIpStatusQueryVariables = Exact<{ [key: string]: never; }>;


export type GeoIpStatusQuery = { __typename: 'Query', geoIPStatus: (
    { __typename: 'GeoIPStatus' }
    & { ' $fragmentRefs'?: { 'GeoIpStatusFieldsFragment': GeoIpStatusFieldsFragment } }
  ) };

export type SiteQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type SiteQuery = { __typename: 'Query', site: (
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'SiteDetailsFieldsFragment': SiteDetailsFieldsFragment } }
  ) | null };

export type UpdateSiteMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  input: UpdateSiteInput;
}>;


export type UpdateSiteMutation = { __typename: 'Mutation', updateSite: (
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'SiteDetailsFieldsFragment': SiteDetailsFieldsFragment } }
  ) };

export type GeoIpStatusFieldsFragment = { __typename: 'GeoIPStatus', state: string, dbPath: string, source: string | null, lastError: string | null, updatedAt: string | null } & { ' $fragmentName'?: 'GeoIpStatusFieldsFragment' };

export type SiteDetailsFieldsFragment = { __typename: 'Site', id: string, domains: Array<string>, name: string, publicKey: string, trackCountry: boolean, blockedIPs: Array<string>, blockedCountries: Array<string>, createdAt: string } & { ' $fragmentName'?: 'SiteDetailsFieldsFragment' };

export type SitesQueryVariables = Exact<{
  paging: PagingInput;
}>;


export type SitesQuery = { __typename: 'Query', sites: Array<(
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'SiteSummaryFieldsFragment': SiteSummaryFieldsFragment } }
  )> };

export type SiteSummaryFieldsFragment = { __typename: 'Site', id: string, domains: Array<string>, name: string, publicKey: string, createdAt: string } & { ' $fragmentName'?: 'SiteSummaryFieldsFragment' };

export const CountryTrackingGeoIpStatusFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryTrackingGeoIPStatusFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPStatus"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"state"}},{"kind":"Field","name":{"kind":"Name","value":"dbPath"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"lastError"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]} as unknown as DocumentNode<CountryTrackingGeoIpStatusFieldsFragment, unknown>;
export const EventDefinitionFieldFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFieldFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionField"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"required"}},{"kind":"Field","name":{"kind":"Name","value":"maxLength"}}]}}]} as unknown as DocumentNode<EventDefinitionFieldFieldsFragment, unknown>;
export const EventDefinitionFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinition"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"fields"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFieldFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFieldFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionField"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"required"}},{"kind":"Field","name":{"kind":"Name","value":"maxLength"}}]}}]} as unknown as DocumentNode<EventDefinitionFieldsFragment, unknown>;
export const TrackingCodeSiteSummaryFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"TrackingCodeSiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<TrackingCodeSiteSummaryFieldsFragment, unknown>;
export const TrafficBlockingCountryFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"TrafficBlockingCountryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPCountry"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]} as unknown as DocumentNode<TrafficBlockingCountryFieldsFragment, unknown>;
export const AuthUserFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}}]}}]} as unknown as DocumentNode<AuthUserFieldsFragment, unknown>;
export const AuthSiteSummaryFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthSiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<AuthSiteSummaryFieldsFragment, unknown>;
export const AuthUserDetailsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"sites"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"sitesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"AuthSiteSummaryFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthSiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<AuthUserDetailsFieldsFragment, unknown>;
export const DailyStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DailyStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DailyStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"date"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"pageViews"}},{"kind":"Field","name":{"kind":"Name","value":"sessions"}}]}}]} as unknown as DocumentNode<DailyStatsFieldsFragment, unknown>;
export const EventPropertyFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}}]} as unknown as DocumentNode<EventPropertyFieldsFragment, unknown>;
export const EventFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Event"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"properties"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventPropertyFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}}]} as unknown as DocumentNode<EventFieldsFragment, unknown>;
export const EventCountFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventCountFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventCount"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"count"}},{"kind":"Field","name":{"kind":"Name","value":"event"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Event"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"properties"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventPropertyFields"}}]}}]}}]} as unknown as DocumentNode<EventCountFieldsFragment, unknown>;
export const ActivePageStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ActivePageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ActivePageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<ActivePageStatsFieldsFragment, unknown>;
export const RealtimeStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"RealtimeStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"RealtimeStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"activePages"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"activePagesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ActivePageStatsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ActivePageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ActivePageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<RealtimeStatsFieldsFragment, unknown>;
export const PageStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"PageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"views"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<PageStatsFieldsFragment, unknown>;
export const ReferrerStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ReferrerStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ReferrerStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referrer"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<ReferrerStatsFieldsFragment, unknown>;
export const BrowserStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"BrowserStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"BrowserStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"browser"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<BrowserStatsFieldsFragment, unknown>;
export const DeviceStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DeviceStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DeviceStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"device"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<DeviceStatsFieldsFragment, unknown>;
export const CountryStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"CountryStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"country"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<CountryStatsFieldsFragment, unknown>;
export const DashboardStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DashboardStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DashboardStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"pageViews"}},{"kind":"Field","name":{"kind":"Name","value":"sessions"}},{"kind":"Field","name":{"kind":"Name","value":"bounceRate"}},{"kind":"Field","name":{"kind":"Name","value":"avgDuration"}},{"kind":"Field","name":{"kind":"Name","value":"topPages"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"topPagesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PageStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"topReferrers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"referrersPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ReferrerStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"browsers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"browsersPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"BrowserStatsFields"}}]}},{"kind":"Field","name":{"kind":"Name","value":"devices"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"devicesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"DeviceStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"countries"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"countriesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryStatsFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"PageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"views"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ReferrerStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ReferrerStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referrer"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"BrowserStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"BrowserStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"browser"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DeviceStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DeviceStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"device"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"CountryStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"country"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<DashboardStatsFieldsFragment, unknown>;
export const GeoIpStatusFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"GeoIPStatusFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPStatus"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"state"}},{"kind":"Field","name":{"kind":"Name","value":"dbPath"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"lastError"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]} as unknown as DocumentNode<GeoIpStatusFieldsFragment, unknown>;
export const SiteDetailsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"trackCountry"}},{"kind":"Field","name":{"kind":"Name","value":"blockedIPs"}},{"kind":"Field","name":{"kind":"Name","value":"blockedCountries"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<SiteDetailsFieldsFragment, unknown>;
export const SiteSummaryFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<SiteSummaryFieldsFragment, unknown>;
export const RefreshGeoIpDatabaseDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"RefreshGeoIPDatabase"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"refreshGeoIPDatabase"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryTrackingGeoIPStatusFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryTrackingGeoIPStatusFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPStatus"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"state"}},{"kind":"Field","name":{"kind":"Name","value":"dbPath"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"lastError"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]} as unknown as DocumentNode<RefreshGeoIpDatabaseMutation, RefreshGeoIpDatabaseMutationVariables>;
export const DeleteSiteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteSite"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteSite"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}]}}]} as unknown as DocumentNode<DeleteSiteMutation, DeleteSiteMutationVariables>;
export const EventDefinitionsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"EventDefinitions"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"eventDefinitions"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFieldFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionField"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"required"}},{"kind":"Field","name":{"kind":"Name","value":"maxLength"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinition"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"fields"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFieldFields"}}]}}]}}]} as unknown as DocumentNode<EventDefinitionsQuery, EventDefinitionsQueryVariables>;
export const UpsertEventDefinitionDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpsertEventDefinition"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"upsertEventDefinition"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFieldFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionField"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"required"}},{"kind":"Field","name":{"kind":"Name","value":"maxLength"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinition"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"fields"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFieldFields"}}]}}]}}]} as unknown as DocumentNode<UpsertEventDefinitionMutation, UpsertEventDefinitionMutationVariables>;
export const DeleteEventDefinitionDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteEventDefinition"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteEventDefinition"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}}]}]}}]} as unknown as DocumentNode<DeleteEventDefinitionMutation, DeleteEventDefinitionMutationVariables>;
export const RegenerateSiteKeyDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"RegenerateSiteKey"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"regenerateSiteKey"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"TrackingCodeSiteSummaryFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"TrackingCodeSiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<RegenerateSiteKeyMutation, RegenerateSiteKeyMutationVariables>;
export const GeoIpCountriesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GeoIPCountries"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"search"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"geoIPCountries"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"search"},"value":{"kind":"Variable","name":{"kind":"Name","value":"search"}}},{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"TrafficBlockingCountryFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"TrafficBlockingCountryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPCountry"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]} as unknown as DocumentNode<GeoIpCountriesQuery, GeoIpCountriesQueryVariables>;
export const MeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Me"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"sitesPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"AuthUserDetailsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthSiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"sites"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"sitesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"AuthSiteSummaryFields"}}]}}]}}]} as unknown as DocumentNode<MeQuery, MeQueryVariables>;
export const LoginDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Login"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"LoginInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"AuthUserFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}}]}}]} as unknown as DocumentNode<LoginMutation, LoginMutationVariables>;
export const RegisterDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Register"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"RegisterInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"register"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"AuthUserFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}}]}}]} as unknown as DocumentNode<RegisterMutation, RegisterMutationVariables>;
export const LogoutDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Logout"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"logout"}}]}}]} as unknown as DocumentNode<LogoutMutation, LogoutMutationVariables>;
export const DashboardDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Dashboard"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"DateRangeInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"FilterInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"topPagesPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"referrersPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"browsersPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"devicesPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"countriesPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"dashboard"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"dateRange"},"value":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"DashboardStatsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"PageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"views"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ReferrerStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ReferrerStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referrer"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"BrowserStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"BrowserStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"browser"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DeviceStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DeviceStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"device"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"CountryStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"country"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DashboardStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DashboardStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"pageViews"}},{"kind":"Field","name":{"kind":"Name","value":"sessions"}},{"kind":"Field","name":{"kind":"Name","value":"bounceRate"}},{"kind":"Field","name":{"kind":"Name","value":"avgDuration"}},{"kind":"Field","name":{"kind":"Name","value":"topPages"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"topPagesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PageStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"topReferrers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"referrersPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ReferrerStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"browsers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"browsersPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"BrowserStatsFields"}}]}},{"kind":"Field","name":{"kind":"Name","value":"devices"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"devicesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"DeviceStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"countries"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"countriesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryStatsFields"}}]}}]}}]}}]} as unknown as DocumentNode<DashboardQuery, DashboardQueryVariables>;
export const RealtimeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Realtime"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"activePagesPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"realtime"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"RealtimeStatsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ActivePageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ActivePageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"RealtimeStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"RealtimeStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"activePages"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"activePagesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ActivePageStatsFields"}}]}}]}}]} as unknown as DocumentNode<RealtimeQuery, RealtimeQueryVariables>;
export const EventsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Events"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"DateRangeInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"FilterInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"limit"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"offset"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"events"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"dateRange"},"value":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}},{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"Variable","name":{"kind":"Name","value":"limit"}}},{"kind":"Argument","name":{"kind":"Name","value":"offset"},"value":{"kind":"Variable","name":{"kind":"Name","value":"offset"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"events"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Event"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"properties"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventPropertyFields"}}]}}]}}]} as unknown as DocumentNode<EventsQuery, EventsQueryVariables>;
export const EventCountsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"EventCounts"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"DateRangeInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"FilterInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"eventCounts"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"dateRange"},"value":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}},{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventCountFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Event"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"properties"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventPropertyFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventCountFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventCount"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"count"}},{"kind":"Field","name":{"kind":"Name","value":"event"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventFields"}}]}}]}}]} as unknown as DocumentNode<EventCountsQuery, EventCountsQueryVariables>;
export const ChartDataDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ChartData"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"DateRangeInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"FilterInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"bucket"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"TimeBucket"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"limit"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"offset"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"dashboard"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"dateRange"},"value":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"dailyStats"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"bucket"},"value":{"kind":"Variable","name":{"kind":"Name","value":"bucket"}}},{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"Variable","name":{"kind":"Name","value":"limit"}}},{"kind":"Argument","name":{"kind":"Name","value":"offset"},"value":{"kind":"Variable","name":{"kind":"Name","value":"offset"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"DailyStatsFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DailyStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DailyStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"date"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"pageViews"}},{"kind":"Field","name":{"kind":"Name","value":"sessions"}}]}}]} as unknown as DocumentNode<ChartDataQuery, ChartDataQueryVariables>;
export const CreateSiteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateSite"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateSiteInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createSite"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"SiteDetailsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"trackCountry"}},{"kind":"Field","name":{"kind":"Name","value":"blockedIPs"}},{"kind":"Field","name":{"kind":"Name","value":"blockedCountries"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<CreateSiteMutation, CreateSiteMutationVariables>;
export const GeoIpStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GeoIPStatus"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"geoIPStatus"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"GeoIPStatusFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"GeoIPStatusFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPStatus"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"state"}},{"kind":"Field","name":{"kind":"Name","value":"dbPath"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"lastError"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]} as unknown as DocumentNode<GeoIpStatusQuery, GeoIpStatusQueryVariables>;
export const SiteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Site"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"site"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"SiteDetailsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"trackCountry"}},{"kind":"Field","name":{"kind":"Name","value":"blockedIPs"}},{"kind":"Field","name":{"kind":"Name","value":"blockedCountries"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<SiteQuery, SiteQueryVariables>;
export const UpdateSiteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateSite"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateSiteInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateSite"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"SiteDetailsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"trackCountry"}},{"kind":"Field","name":{"kind":"Name","value":"blockedIPs"}},{"kind":"Field","name":{"kind":"Name","value":"blockedCountries"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<UpdateSiteMutation, UpdateSiteMutationVariables>;
export const SitesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Sites"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"sites"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"SiteSummaryFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<SitesQuery, SitesQueryVariables>;