/* eslint-disable */
/** Internal type. DO NOT USE DIRECTLY. */
type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
/** Internal type. DO NOT USE DIRECTLY. */
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type CreateSiteInput = {
  domains: Array<string>;
  name: string;
};

export type DateRangeInput = {
  from: string | null | undefined;
  to: string | null | undefined;
};

export type EventDefinitionFieldInput = {
  key: string;
  maxLength: number | null | undefined;
  required: boolean;
  type: EventFieldType;
};

export type EventDefinitionInput = {
  fields: Array<EventDefinitionFieldInput>;
  name: string;
};

export type EventFieldType =
  | 'BOOLEAN'
  | 'INT'
  | 'STRING';

export type EventType =
  | 'PAGE_VIEW'
  | 'PREDEFINED';

export type FilterInput = {
  /** Filter by browser type */
  browser: Array<string> | null | undefined;
  /** Filter by ISO country code */
  country: Array<string> | null | undefined;
  /** Filter by device type (desktop, mobile, tablet, smart-tv, console, watch) */
  device: Array<string> | null | undefined;
  /** Filter by event definition ID */
  eventDefinitionId: Array<string | number> | null | undefined;
  /** Filter by event name */
  eventName: Array<string> | null | undefined;
  /** Filter by event path */
  eventPath: Array<string> | null | undefined;
  /** Filter by event type (page view or predefined) */
  eventType: Array<EventType> | null | undefined;
  /** Filter by operating system */
  os: Array<string> | null | undefined;
  /** Filter by page path */
  page: Array<string> | null | undefined;
  /** Filter by specific referrer */
  referrer: Array<string> | null | undefined;
};

export type GeoIpState =
  | 'DISABLED'
  | 'DOWNLOADING'
  | 'ERROR'
  | 'MISSING'
  | 'READY';

export type LoginInput = {
  password: string;
  username: string;
};

export type PagingInput = {
  limit: number;
  offset: number;
};

export type RegisterInput = {
  password: string;
  username: string;
};

export type TimeBucket =
  | 'DAILY'
  | 'HOURLY';

export type UpdateSiteInput = {
  /** Full list of blocked country codes */
  blockedCountries: Array<string> | null | undefined;
  /** Full list of blocked IPs */
  blockedIPs: Array<string> | null | undefined;
  /** Full list of tracked domains (includes primary) */
  domains: Array<string> | null | undefined;
  name: string;
  trackCountry: boolean | null | undefined;
};

export type DashboardQueryVariables = Exact<{
  siteId: string | number;
  dateRange?: DateRangeInput | null | undefined;
  filter?: FilterInput | null | undefined;
  topPagesPaging: PagingInput;
  referrersPaging: PagingInput;
  browsersPaging: PagingInput;
  devicesPaging: PagingInput;
  osPaging: PagingInput;
  countriesPaging: PagingInput;
}>;


export type DashboardQuery = { __typename: 'Query', dashboard: (
    { __typename: 'DashboardStats' }
    & { ' $fragmentRefs'?: { 'DashboardStatsFieldsFragment': DashboardStatsFieldsFragment } }
  ) };

export type RealtimeQueryVariables = Exact<{
  siteId: string | number;
  activePagesPaging: PagingInput;
}>;


export type RealtimeQuery = { __typename: 'Query', realtime: (
    { __typename: 'RealtimeStats' }
    & { ' $fragmentRefs'?: { 'RealtimeStatsFieldsFragment': RealtimeStatsFieldsFragment } }
  ) };

export type EventsQueryVariables = Exact<{
  siteId: string | number;
  dateRange?: DateRangeInput | null | undefined;
  filter?: FilterInput | null | undefined;
  paging: PagingInput;
}>;


export type EventsQuery = { __typename: 'Query', events: { __typename: 'EventsResult', total: number, events: Array<(
      { __typename: 'Event' }
      & { ' $fragmentRefs'?: { 'EventFieldsFragment': EventFieldsFragment } }
    )> } };

export type EventCountsQueryVariables = Exact<{
  siteId: string | number;
  dateRange?: DateRangeInput | null | undefined;
  filter?: FilterInput | null | undefined;
  paging: PagingInput;
}>;


export type EventCountsQuery = { __typename: 'Query', eventCounts: { __typename: 'EventCountsResult', total: number, items: Array<(
      { __typename: 'EventCount' }
      & { ' $fragmentRefs'?: { 'EventCountFieldsFragment': EventCountFieldsFragment } }
    )> } };

export type ChartDataQueryVariables = Exact<{
  siteId: string | number;
  dateRange?: DateRangeInput | null | undefined;
  filter?: FilterInput | null | undefined;
  bucket?: TimeBucket | null | undefined;
  paging: PagingInput;
}>;


export type ChartDataQuery = { __typename: 'Query', dashboard: { __typename: 'DashboardStats', dailyStats: Array<(
      { __typename: 'DailyStats' }
      & { ' $fragmentRefs'?: { 'DailyStatsFieldsFragment': DailyStatsFieldsFragment } }
    )> } };

export type PageStatsFieldsFragment = { __typename: 'PageStats', path: string, views: number, visitors: number } & { ' $fragmentName'?: 'PageStatsFieldsFragment' };

export type ReferrerStatsFieldsFragment = { __typename: 'ReferrerStats', referrer: string, visitors: number } & { ' $fragmentName'?: 'ReferrerStatsFieldsFragment' };

export type BrowserStatsFieldsFragment = { __typename: 'BrowserStats', browser: string, visitors: number } & { ' $fragmentName'?: 'BrowserStatsFieldsFragment' };

export type DeviceStatsFieldsFragment = { __typename: 'DeviceStats', device: string, visitors: number } & { ' $fragmentName'?: 'DeviceStatsFieldsFragment' };

export type OperatingSystemStatsFieldsFragment = { __typename: 'OperatingSystemStats', os: string, visitors: number } & { ' $fragmentName'?: 'OperatingSystemStatsFieldsFragment' };

export type CountryStatsFieldsFragment = { __typename: 'CountryStats', visitors: number, country: (
    { __typename: 'Country' }
    & { ' $fragmentRefs'?: { 'CountryFieldsFragment': CountryFieldsFragment } }
  ) } & { ' $fragmentName'?: 'CountryStatsFieldsFragment' };

export type ActivePageStatsFieldsFragment = { __typename: 'ActivePageStats', path: string, visitors: number } & { ' $fragmentName'?: 'ActivePageStatsFieldsFragment' };

export type DailyStatsFieldsFragment = { __typename: 'DailyStats', date: string, visitors: number, pageViews: number, sessions: number } & { ' $fragmentName'?: 'DailyStatsFieldsFragment' };

export type EventPropertyFieldsFragment = { __typename: 'EventProperty', key: string, value: string } & { ' $fragmentName'?: 'EventPropertyFieldsFragment' };

export type EventFieldsFragment = { __typename: 'Event', id: string, name: string, path: string, createdAt: string, definition: { __typename: 'EventDefinition', id: string, name: string } | null, properties: Array<(
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
    )> }, operatingSystems: { __typename: 'PagedOperatingSystemStats', total: number, totalVisitors: number, items: Array<(
      { __typename: 'OperatingSystemStats' }
      & { ' $fragmentRefs'?: { 'OperatingSystemStatsFieldsFragment': OperatingSystemStatsFieldsFragment } }
    )> }, countries: { __typename: 'PagedCountryStats', total: number, totalVisitors: number, items: Array<(
      { __typename: 'CountryStats' }
      & { ' $fragmentRefs'?: { 'CountryStatsFieldsFragment': CountryStatsFieldsFragment } }
    )> } } & { ' $fragmentName'?: 'DashboardStatsFieldsFragment' };

export type MeQueryVariables = Exact<{ [key: string]: never; }>;


export type MeQuery = { __typename: 'Query', me: (
    { __typename: 'User' }
    & { ' $fragmentRefs'?: { 'AuthUserDetailsFieldsFragment': AuthUserDetailsFieldsFragment } }
  ) | null, registrationStatus: { __typename: 'RegistrationStatus', hasUsers: boolean, allowRegistration: boolean } };

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

export type AuthUserDetailsFieldsFragment = { __typename: 'User', id: string, username: string, role: string, createdAt: string } & { ' $fragmentName'?: 'AuthUserDetailsFieldsFragment' };

export type EventDefinitionsQueryVariables = Exact<{
  siteId: string | number;
  paging: PagingInput;
}>;


export type EventDefinitionsQuery = { __typename: 'Query', eventDefinitions: Array<(
    { __typename: 'EventDefinition' }
    & { ' $fragmentRefs'?: { 'EventDefinitionFieldsFragment': EventDefinitionFieldsFragment } }
  )> };

export type UpsertEventDefinitionMutationVariables = Exact<{
  siteId: string | number;
  input: EventDefinitionInput;
}>;


export type UpsertEventDefinitionMutation = { __typename: 'Mutation', upsertEventDefinition: (
    { __typename: 'EventDefinition' }
    & { ' $fragmentRefs'?: { 'EventDefinitionFieldsFragment': EventDefinitionFieldsFragment } }
  ) };

export type DeleteEventDefinitionMutationVariables = Exact<{
  siteId: string | number;
  name: string;
}>;


export type DeleteEventDefinitionMutation = { __typename: 'Mutation', deleteEventDefinition: boolean };

export type EventDefinitionFieldFieldsFragment = { __typename: 'EventDefinitionField', id: string, key: string, type: EventFieldType, required: boolean, maxLength: number } & { ' $fragmentName'?: 'EventDefinitionFieldFieldsFragment' };

export type EventDefinitionFieldsFragment = { __typename: 'EventDefinition', id: string, name: string, createdAt: string, updatedAt: string, fields: Array<(
    { __typename: 'EventDefinitionField' }
    & { ' $fragmentRefs'?: { 'EventDefinitionFieldFieldsFragment': EventDefinitionFieldFieldsFragment } }
  )> } & { ' $fragmentName'?: 'EventDefinitionFieldsFragment' };

export type RefreshGeoIpDatabaseMutationVariables = Exact<{ [key: string]: never; }>;


export type RefreshGeoIpDatabaseMutation = { __typename: 'Mutation', refreshGeoIPDatabase: (
    { __typename: 'GeoIPStatus' }
    & { ' $fragmentRefs'?: { 'CountryTrackingGeoIpStatusFieldsFragment': CountryTrackingGeoIpStatusFieldsFragment } }
  ) };

export type CountryTrackingGeoIpStatusFieldsFragment = { __typename: 'GeoIPStatus', state: GeoIpState, dbPath: string, source: string | null, lastError: string | null, updatedAt: string | null } & { ' $fragmentName'?: 'CountryTrackingGeoIpStatusFieldsFragment' };

export type DeleteSiteMutationVariables = Exact<{
  id: string | number;
}>;


export type DeleteSiteMutation = { __typename: 'Mutation', deleteSite: boolean };

export type GeoIpStatusQueryVariables = Exact<{ [key: string]: never; }>;


export type GeoIpStatusQuery = { __typename: 'Query', geoIPStatus: (
    { __typename: 'GeoIPStatus' }
    & { ' $fragmentRefs'?: { 'GeoIpStatusFieldsFragment': GeoIpStatusFieldsFragment } }
  ) };

export type SiteQueryVariables = Exact<{
  id: string | number;
}>;


export type SiteQuery = { __typename: 'Query', site: (
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'SiteDetailsFieldsFragment': SiteDetailsFieldsFragment } }
  ) | null };

export type UpdateSiteMutationVariables = Exact<{
  id: string | number;
  input: UpdateSiteInput;
}>;


export type UpdateSiteMutation = { __typename: 'Mutation', updateSite: (
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'SiteDetailsFieldsFragment': SiteDetailsFieldsFragment } }
  ) };

export type GeoIpStatusFieldsFragment = { __typename: 'GeoIPStatus', state: GeoIpState, dbPath: string, source: string | null, lastError: string | null, updatedAt: string | null } & { ' $fragmentName'?: 'GeoIpStatusFieldsFragment' };

export type SiteDetailsFieldsFragment = { __typename: 'Site', id: string, domains: Array<string>, name: string, publicKey: string, trackCountry: boolean, blockedIPs: Array<string>, blockedCountries: Array<string>, createdAt: string } & { ' $fragmentName'?: 'SiteDetailsFieldsFragment' };

export type RegenerateSiteKeyMutationVariables = Exact<{
  id: string | number;
}>;


export type RegenerateSiteKeyMutation = { __typename: 'Mutation', regenerateSiteKey: (
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'TrackingCodeSiteSummaryFieldsFragment': TrackingCodeSiteSummaryFieldsFragment } }
  ) };

export type TrackingCodeSiteSummaryFieldsFragment = { __typename: 'Site', id: string, domains: Array<string>, name: string, publicKey: string, createdAt: string } & { ' $fragmentName'?: 'TrackingCodeSiteSummaryFieldsFragment' };

export type GeoIpCountriesQueryVariables = Exact<{
  search?: string | null | undefined;
  paging: PagingInput;
}>;


export type GeoIpCountriesQuery = { __typename: 'Query', geoIPCountries: Array<(
    { __typename: 'Country' }
    & { ' $fragmentRefs'?: { 'CountryFieldsFragment': CountryFieldsFragment } }
  )> };

export type SitesQueryVariables = Exact<{
  paging: PagingInput;
}>;


export type SitesQuery = { __typename: 'Query', sites: Array<(
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'SiteSummaryFieldsFragment': SiteSummaryFieldsFragment } }
  )> };

export type CreateSiteMutationVariables = Exact<{
  input: CreateSiteInput;
}>;


export type CreateSiteMutation = { __typename: 'Mutation', createSite: (
    { __typename: 'Site' }
    & { ' $fragmentRefs'?: { 'SiteSummaryFieldsFragment': SiteSummaryFieldsFragment } }
  ) };

export type SiteSummaryFieldsFragment = { __typename: 'Site', id: string, domains: Array<string>, name: string, publicKey: string, createdAt: string } & { ' $fragmentName'?: 'SiteSummaryFieldsFragment' };

export type CountryFieldsFragment = { __typename: 'Country', code: string, name: string } & { ' $fragmentName'?: 'CountryFieldsFragment' };

export type CountriesByCodeQueryVariables = Exact<{
  codes?: Array<string> | string | null | undefined;
  paging: PagingInput;
}>;


export type CountriesByCodeQuery = { __typename: 'Query', geoIPCountries: Array<(
    { __typename: 'Country' }
    & { ' $fragmentRefs'?: { 'CountryFieldsFragment': CountryFieldsFragment } }
  )> };

export const DailyStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DailyStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DailyStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"date"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"pageViews"}},{"kind":"Field","name":{"kind":"Name","value":"sessions"}}]}}]} as unknown as DocumentNode<DailyStatsFieldsFragment, unknown>;
export const EventPropertyFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}}]} as unknown as DocumentNode<EventPropertyFieldsFragment, unknown>;
export const EventFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Event"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"definition"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"properties"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventPropertyFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}}]} as unknown as DocumentNode<EventFieldsFragment, unknown>;
export const EventCountFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventCountFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventCount"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"count"}},{"kind":"Field","name":{"kind":"Name","value":"event"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Event"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"definition"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"properties"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventPropertyFields"}}]}}]}}]} as unknown as DocumentNode<EventCountFieldsFragment, unknown>;
export const ActivePageStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ActivePageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ActivePageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<ActivePageStatsFieldsFragment, unknown>;
export const RealtimeStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"RealtimeStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"RealtimeStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"activePages"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"activePagesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ActivePageStatsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ActivePageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ActivePageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<RealtimeStatsFieldsFragment, unknown>;
export const PageStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"PageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"views"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<PageStatsFieldsFragment, unknown>;
export const ReferrerStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ReferrerStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ReferrerStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referrer"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<ReferrerStatsFieldsFragment, unknown>;
export const BrowserStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"BrowserStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"BrowserStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"browser"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<BrowserStatsFieldsFragment, unknown>;
export const DeviceStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DeviceStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DeviceStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"device"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<DeviceStatsFieldsFragment, unknown>;
export const OperatingSystemStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"OperatingSystemStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"OperatingSystemStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"os"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<OperatingSystemStatsFieldsFragment, unknown>;
export const CountryFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Country"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]} as unknown as DocumentNode<CountryFieldsFragment, unknown>;
export const CountryStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"CountryStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"country"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryFields"}}]}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Country"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]} as unknown as DocumentNode<CountryStatsFieldsFragment, unknown>;
export const DashboardStatsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DashboardStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DashboardStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"pageViews"}},{"kind":"Field","name":{"kind":"Name","value":"sessions"}},{"kind":"Field","name":{"kind":"Name","value":"bounceRate"}},{"kind":"Field","name":{"kind":"Name","value":"avgDuration"}},{"kind":"Field","name":{"kind":"Name","value":"topPages"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"topPagesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PageStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"topReferrers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"referrersPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ReferrerStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"browsers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"browsersPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"BrowserStatsFields"}}]}},{"kind":"Field","name":{"kind":"Name","value":"devices"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"devicesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"DeviceStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"operatingSystems"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"osPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"OperatingSystemStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"countries"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"countriesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryStatsFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Country"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"PageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"views"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ReferrerStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ReferrerStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referrer"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"BrowserStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"BrowserStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"browser"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DeviceStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DeviceStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"device"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"OperatingSystemStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"OperatingSystemStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"os"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"CountryStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"country"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryFields"}}]}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}}]} as unknown as DocumentNode<DashboardStatsFieldsFragment, unknown>;
export const AuthUserFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}}]}}]} as unknown as DocumentNode<AuthUserFieldsFragment, unknown>;
export const AuthUserDetailsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<AuthUserDetailsFieldsFragment, unknown>;
export const EventDefinitionFieldFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFieldFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionField"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"required"}},{"kind":"Field","name":{"kind":"Name","value":"maxLength"}}]}}]} as unknown as DocumentNode<EventDefinitionFieldFieldsFragment, unknown>;
export const EventDefinitionFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinition"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"fields"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFieldFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFieldFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionField"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"required"}},{"kind":"Field","name":{"kind":"Name","value":"maxLength"}}]}}]} as unknown as DocumentNode<EventDefinitionFieldsFragment, unknown>;
export const CountryTrackingGeoIpStatusFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryTrackingGeoIPStatusFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPStatus"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"state"}},{"kind":"Field","name":{"kind":"Name","value":"dbPath"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"lastError"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]} as unknown as DocumentNode<CountryTrackingGeoIpStatusFieldsFragment, unknown>;
export const GeoIpStatusFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"GeoIPStatusFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPStatus"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"state"}},{"kind":"Field","name":{"kind":"Name","value":"dbPath"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"lastError"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]} as unknown as DocumentNode<GeoIpStatusFieldsFragment, unknown>;
export const SiteDetailsFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"trackCountry"}},{"kind":"Field","name":{"kind":"Name","value":"blockedIPs"}},{"kind":"Field","name":{"kind":"Name","value":"blockedCountries"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<SiteDetailsFieldsFragment, unknown>;
export const TrackingCodeSiteSummaryFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"TrackingCodeSiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<TrackingCodeSiteSummaryFieldsFragment, unknown>;
export const SiteSummaryFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<SiteSummaryFieldsFragment, unknown>;
export const DashboardDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Dashboard"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"DateRangeInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"FilterInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"topPagesPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"referrersPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"browsersPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"devicesPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"osPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"countriesPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"dashboard"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"dateRange"},"value":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"DashboardStatsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"PageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"views"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ReferrerStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ReferrerStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referrer"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"BrowserStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"BrowserStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"browser"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DeviceStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DeviceStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"device"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"OperatingSystemStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"OperatingSystemStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"os"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Country"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"CountryStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"country"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryFields"}}]}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DashboardStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DashboardStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"pageViews"}},{"kind":"Field","name":{"kind":"Name","value":"sessions"}},{"kind":"Field","name":{"kind":"Name","value":"bounceRate"}},{"kind":"Field","name":{"kind":"Name","value":"avgDuration"}},{"kind":"Field","name":{"kind":"Name","value":"topPages"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"topPagesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PageStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"topReferrers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"referrersPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ReferrerStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"browsers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"browsersPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"BrowserStatsFields"}}]}},{"kind":"Field","name":{"kind":"Name","value":"devices"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"devicesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"DeviceStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"operatingSystems"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"osPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"OperatingSystemStatsFields"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"countries"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"countriesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"totalVisitors"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryStatsFields"}}]}}]}}]}}]} as unknown as DocumentNode<DashboardQuery, DashboardQueryVariables>;
export const RealtimeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Realtime"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"activePagesPaging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"realtime"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"RealtimeStatsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ActivePageStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ActivePageStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"RealtimeStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"RealtimeStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"activePages"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"activePagesPaging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ActivePageStatsFields"}}]}}]}}]} as unknown as DocumentNode<RealtimeQuery, RealtimeQueryVariables>;
export const EventsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Events"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"DateRangeInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"FilterInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"events"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"dateRange"},"value":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}},{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"events"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Event"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"definition"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"properties"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventPropertyFields"}}]}}]}}]} as unknown as DocumentNode<EventsQuery, EventsQueryVariables>;
export const EventCountsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"EventCounts"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"DateRangeInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"FilterInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"eventCounts"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"dateRange"},"value":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}},{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventCountFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventPropertyFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventProperty"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"value"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Event"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"definition"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"properties"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventPropertyFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventCountFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventCount"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"count"}},{"kind":"Field","name":{"kind":"Name","value":"event"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventFields"}}]}}]}}]} as unknown as DocumentNode<EventCountsQuery, EventCountsQueryVariables>;
export const ChartDataDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ChartData"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"DateRangeInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"FilterInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"bucket"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"TimeBucket"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"dashboard"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"dateRange"},"value":{"kind":"Variable","name":{"kind":"Name","value":"dateRange"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"dailyStats"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"bucket"},"value":{"kind":"Variable","name":{"kind":"Name","value":"bucket"}}},{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"DailyStatsFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"DailyStatsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DailyStats"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"date"}},{"kind":"Field","name":{"kind":"Name","value":"visitors"}},{"kind":"Field","name":{"kind":"Name","value":"pageViews"}},{"kind":"Field","name":{"kind":"Name","value":"sessions"}}]}}]} as unknown as DocumentNode<ChartDataQuery, ChartDataQueryVariables>;
export const MeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"AuthUserDetailsFields"}}]}},{"kind":"Field","name":{"kind":"Name","value":"registrationStatus"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"hasUsers"}},{"kind":"Field","name":{"kind":"Name","value":"allowRegistration"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<MeQuery, MeQueryVariables>;
export const LoginDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Login"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"LoginInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"AuthUserFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}}]}}]} as unknown as DocumentNode<LoginMutation, LoginMutationVariables>;
export const RegisterDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Register"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"RegisterInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"register"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"AuthUserFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"AuthUserFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"User"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"role"}}]}}]} as unknown as DocumentNode<RegisterMutation, RegisterMutationVariables>;
export const LogoutDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Logout"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"logout"}}]}}]} as unknown as DocumentNode<LogoutMutation, LogoutMutationVariables>;
export const EventDefinitionsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"EventDefinitions"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"eventDefinitions"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFieldFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionField"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"required"}},{"kind":"Field","name":{"kind":"Name","value":"maxLength"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinition"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"fields"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFieldFields"}}]}}]}}]} as unknown as DocumentNode<EventDefinitionsQuery, EventDefinitionsQueryVariables>;
export const UpsertEventDefinitionDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpsertEventDefinition"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"upsertEventDefinition"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFieldFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinitionField"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"key"}},{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"required"}},{"kind":"Field","name":{"kind":"Name","value":"maxLength"}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"EventDefinitionFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EventDefinition"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"fields"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"EventDefinitionFieldFields"}}]}}]}}]} as unknown as DocumentNode<UpsertEventDefinitionMutation, UpsertEventDefinitionMutationVariables>;
export const DeleteEventDefinitionDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteEventDefinition"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteEventDefinition"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"siteId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"siteId"}}},{"kind":"Argument","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}}]}]}}]} as unknown as DocumentNode<DeleteEventDefinitionMutation, DeleteEventDefinitionMutationVariables>;
export const RefreshGeoIpDatabaseDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"RefreshGeoIPDatabase"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"refreshGeoIPDatabase"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryTrackingGeoIPStatusFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryTrackingGeoIPStatusFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPStatus"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"state"}},{"kind":"Field","name":{"kind":"Name","value":"dbPath"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"lastError"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]} as unknown as DocumentNode<RefreshGeoIpDatabaseMutation, RefreshGeoIpDatabaseMutationVariables>;
export const DeleteSiteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteSite"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteSite"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}]}}]} as unknown as DocumentNode<DeleteSiteMutation, DeleteSiteMutationVariables>;
export const GeoIpStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GeoIPStatus"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"geoIPStatus"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"GeoIPStatusFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"GeoIPStatusFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"GeoIPStatus"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"state"}},{"kind":"Field","name":{"kind":"Name","value":"dbPath"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"lastError"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]} as unknown as DocumentNode<GeoIpStatusQuery, GeoIpStatusQueryVariables>;
export const SiteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Site"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"site"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"SiteDetailsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"trackCountry"}},{"kind":"Field","name":{"kind":"Name","value":"blockedIPs"}},{"kind":"Field","name":{"kind":"Name","value":"blockedCountries"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<SiteQuery, SiteQueryVariables>;
export const UpdateSiteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateSite"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateSiteInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateSite"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"SiteDetailsFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteDetailsFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"trackCountry"}},{"kind":"Field","name":{"kind":"Name","value":"blockedIPs"}},{"kind":"Field","name":{"kind":"Name","value":"blockedCountries"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<UpdateSiteMutation, UpdateSiteMutationVariables>;
export const RegenerateSiteKeyDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"RegenerateSiteKey"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"regenerateSiteKey"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"TrackingCodeSiteSummaryFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"TrackingCodeSiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<RegenerateSiteKeyMutation, RegenerateSiteKeyMutationVariables>;
export const GeoIpCountriesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GeoIPCountries"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"search"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"geoIPCountries"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"search"},"value":{"kind":"Variable","name":{"kind":"Name","value":"search"}}},{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Country"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]} as unknown as DocumentNode<GeoIpCountriesQuery, GeoIpCountriesQueryVariables>;
export const SitesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Sites"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"sites"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"SiteSummaryFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<SitesQuery, SitesQueryVariables>;
export const CreateSiteDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateSite"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateSiteInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createSite"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"SiteSummaryFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"SiteSummaryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Site"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"domains"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"publicKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]} as unknown as DocumentNode<CreateSiteMutation, CreateSiteMutationVariables>;
export const CountriesByCodeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"CountriesByCode"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"codes"}},"type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"paging"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"PagingInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"geoIPCountries"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"codes"},"value":{"kind":"Variable","name":{"kind":"Name","value":"codes"}}},{"kind":"Argument","name":{"kind":"Name","value":"paging"},"value":{"kind":"Variable","name":{"kind":"Name","value":"paging"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"CountryFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"CountryFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Country"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]} as unknown as DocumentNode<CountriesByCodeQuery, CountriesByCodeQueryVariables>;