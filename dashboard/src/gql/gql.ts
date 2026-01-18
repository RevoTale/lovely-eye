/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
    "mutation RefreshGeoIPDatabase {\n  refreshGeoIPDatabase {\n    state\n    dbPath\n    source\n    lastError\n    updatedAt\n  }\n}": typeof types.RefreshGeoIpDatabaseDocument,
    "mutation DeleteSite($id: ID!) {\n  deleteSite(id: $id)\n}": typeof types.DeleteSiteDocument,
    "query EventDefinitions($siteId: ID!) {\n  eventDefinitions(siteId: $siteId) {\n    id\n    name\n    createdAt\n    updatedAt\n    fields {\n      id\n      key\n      type\n      required\n      maxLength\n    }\n  }\n}\n\nmutation UpsertEventDefinition($siteId: ID!, $input: EventDefinitionInput!) {\n  upsertEventDefinition(siteId: $siteId, input: $input) {\n    id\n    name\n    createdAt\n    updatedAt\n    fields {\n      id\n      key\n      type\n      required\n      maxLength\n    }\n  }\n}\n\nmutation DeleteEventDefinition($siteId: ID!, $name: String!) {\n  deleteEventDefinition(siteId: $siteId, name: $name)\n}": typeof types.EventDefinitionsDocument,
    "mutation RegenerateSiteKey($id: ID!) {\n  regenerateSiteKey(id: $id) {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}": typeof types.RegenerateSiteKeyDocument,
    "query GeoIPCountries($search: String) {\n  geoIPCountries(search: $search) {\n    code\n    name\n  }\n}": typeof types.GeoIpCountriesDocument,
    "query Me {\n  me {\n    id\n    username\n    role\n    createdAt\n    sites {\n      id\n      domains\n      name\n      publicKey\n      createdAt\n    }\n  }\n}\n\nmutation Login($input: LoginInput!) {\n  login(input: $input) {\n    user {\n      id\n      username\n      role\n    }\n  }\n}\n\nmutation Register($input: RegisterInput!) {\n  register(input: $input) {\n    user {\n      id\n      username\n      role\n    }\n  }\n}\n\nmutation Logout {\n  logout\n}": typeof types.MeDocument,
    "query Dashboard($siteId: ID!, $dateRange: DateRangeInput, $filter: FilterInput) {\n  dashboard(siteId: $siteId, dateRange: $dateRange, filter: $filter) {\n    visitors\n    pageViews\n    sessions\n    bounceRate\n    avgDuration\n    topPages {\n      path\n      views\n      visitors\n    }\n    topReferrers {\n      referrer\n      visitors\n    }\n    browsers {\n      browser\n      visitors\n    }\n    devices {\n      device\n      visitors\n    }\n    countries {\n      country\n      visitors\n    }\n    dailyStats {\n      date\n      visitors\n      pageViews\n      sessions\n    }\n  }\n}\n\nquery Realtime($siteId: ID!) {\n  realtime(siteId: $siteId) {\n    visitors\n    activePages {\n      path\n      visitors\n    }\n  }\n}\n\nquery Events($siteId: ID!, $dateRange: DateRangeInput, $limit: Int, $offset: Int) {\n  events(siteId: $siteId, dateRange: $dateRange, limit: $limit, offset: $offset) {\n    total\n    events {\n      id\n      name\n      path\n      createdAt\n      properties {\n        key\n        value\n      }\n    }\n  }\n}": typeof types.DashboardDocument,
    "mutation CreateSite($input: CreateSiteInput!) {\n  createSite(input: $input) {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}\n\nquery GeoIPStatus {\n  geoIPStatus {\n    state\n    dbPath\n    source\n    lastError\n    updatedAt\n  }\n}\n\nquery Site($id: ID!) {\n  site(id: $id) {\n    id\n    domains\n    name\n    publicKey\n    trackCountry\n    blockedIPs\n    blockedCountries\n    createdAt\n  }\n}\n\nmutation UpdateSite($id: ID!, $input: UpdateSiteInput!) {\n  updateSite(id: $id, input: $input) {\n    id\n    domains\n    name\n    publicKey\n    trackCountry\n    blockedIPs\n    blockedCountries\n    createdAt\n  }\n}": typeof types.CreateSiteDocument,
    "query Sites {\n  sites {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}": typeof types.SitesDocument,
};
const documents: Documents = {
    "mutation RefreshGeoIPDatabase {\n  refreshGeoIPDatabase {\n    state\n    dbPath\n    source\n    lastError\n    updatedAt\n  }\n}": types.RefreshGeoIpDatabaseDocument,
    "mutation DeleteSite($id: ID!) {\n  deleteSite(id: $id)\n}": types.DeleteSiteDocument,
    "query EventDefinitions($siteId: ID!) {\n  eventDefinitions(siteId: $siteId) {\n    id\n    name\n    createdAt\n    updatedAt\n    fields {\n      id\n      key\n      type\n      required\n      maxLength\n    }\n  }\n}\n\nmutation UpsertEventDefinition($siteId: ID!, $input: EventDefinitionInput!) {\n  upsertEventDefinition(siteId: $siteId, input: $input) {\n    id\n    name\n    createdAt\n    updatedAt\n    fields {\n      id\n      key\n      type\n      required\n      maxLength\n    }\n  }\n}\n\nmutation DeleteEventDefinition($siteId: ID!, $name: String!) {\n  deleteEventDefinition(siteId: $siteId, name: $name)\n}": types.EventDefinitionsDocument,
    "mutation RegenerateSiteKey($id: ID!) {\n  regenerateSiteKey(id: $id) {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}": types.RegenerateSiteKeyDocument,
    "query GeoIPCountries($search: String) {\n  geoIPCountries(search: $search) {\n    code\n    name\n  }\n}": types.GeoIpCountriesDocument,
    "query Me {\n  me {\n    id\n    username\n    role\n    createdAt\n    sites {\n      id\n      domains\n      name\n      publicKey\n      createdAt\n    }\n  }\n}\n\nmutation Login($input: LoginInput!) {\n  login(input: $input) {\n    user {\n      id\n      username\n      role\n    }\n  }\n}\n\nmutation Register($input: RegisterInput!) {\n  register(input: $input) {\n    user {\n      id\n      username\n      role\n    }\n  }\n}\n\nmutation Logout {\n  logout\n}": types.MeDocument,
    "query Dashboard($siteId: ID!, $dateRange: DateRangeInput, $filter: FilterInput) {\n  dashboard(siteId: $siteId, dateRange: $dateRange, filter: $filter) {\n    visitors\n    pageViews\n    sessions\n    bounceRate\n    avgDuration\n    topPages {\n      path\n      views\n      visitors\n    }\n    topReferrers {\n      referrer\n      visitors\n    }\n    browsers {\n      browser\n      visitors\n    }\n    devices {\n      device\n      visitors\n    }\n    countries {\n      country\n      visitors\n    }\n    dailyStats {\n      date\n      visitors\n      pageViews\n      sessions\n    }\n  }\n}\n\nquery Realtime($siteId: ID!) {\n  realtime(siteId: $siteId) {\n    visitors\n    activePages {\n      path\n      visitors\n    }\n  }\n}\n\nquery Events($siteId: ID!, $dateRange: DateRangeInput, $limit: Int, $offset: Int) {\n  events(siteId: $siteId, dateRange: $dateRange, limit: $limit, offset: $offset) {\n    total\n    events {\n      id\n      name\n      path\n      createdAt\n      properties {\n        key\n        value\n      }\n    }\n  }\n}": types.DashboardDocument,
    "mutation CreateSite($input: CreateSiteInput!) {\n  createSite(input: $input) {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}\n\nquery GeoIPStatus {\n  geoIPStatus {\n    state\n    dbPath\n    source\n    lastError\n    updatedAt\n  }\n}\n\nquery Site($id: ID!) {\n  site(id: $id) {\n    id\n    domains\n    name\n    publicKey\n    trackCountry\n    blockedIPs\n    blockedCountries\n    createdAt\n  }\n}\n\nmutation UpdateSite($id: ID!, $input: UpdateSiteInput!) {\n  updateSite(id: $id, input: $input) {\n    id\n    domains\n    name\n    publicKey\n    trackCountry\n    blockedIPs\n    blockedCountries\n    createdAt\n  }\n}": types.CreateSiteDocument,
    "query Sites {\n  sites {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}": types.SitesDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "mutation RefreshGeoIPDatabase {\n  refreshGeoIPDatabase {\n    state\n    dbPath\n    source\n    lastError\n    updatedAt\n  }\n}"): (typeof documents)["mutation RefreshGeoIPDatabase {\n  refreshGeoIPDatabase {\n    state\n    dbPath\n    source\n    lastError\n    updatedAt\n  }\n}"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "mutation DeleteSite($id: ID!) {\n  deleteSite(id: $id)\n}"): (typeof documents)["mutation DeleteSite($id: ID!) {\n  deleteSite(id: $id)\n}"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "query EventDefinitions($siteId: ID!) {\n  eventDefinitions(siteId: $siteId) {\n    id\n    name\n    createdAt\n    updatedAt\n    fields {\n      id\n      key\n      type\n      required\n      maxLength\n    }\n  }\n}\n\nmutation UpsertEventDefinition($siteId: ID!, $input: EventDefinitionInput!) {\n  upsertEventDefinition(siteId: $siteId, input: $input) {\n    id\n    name\n    createdAt\n    updatedAt\n    fields {\n      id\n      key\n      type\n      required\n      maxLength\n    }\n  }\n}\n\nmutation DeleteEventDefinition($siteId: ID!, $name: String!) {\n  deleteEventDefinition(siteId: $siteId, name: $name)\n}"): (typeof documents)["query EventDefinitions($siteId: ID!) {\n  eventDefinitions(siteId: $siteId) {\n    id\n    name\n    createdAt\n    updatedAt\n    fields {\n      id\n      key\n      type\n      required\n      maxLength\n    }\n  }\n}\n\nmutation UpsertEventDefinition($siteId: ID!, $input: EventDefinitionInput!) {\n  upsertEventDefinition(siteId: $siteId, input: $input) {\n    id\n    name\n    createdAt\n    updatedAt\n    fields {\n      id\n      key\n      type\n      required\n      maxLength\n    }\n  }\n}\n\nmutation DeleteEventDefinition($siteId: ID!, $name: String!) {\n  deleteEventDefinition(siteId: $siteId, name: $name)\n}"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "mutation RegenerateSiteKey($id: ID!) {\n  regenerateSiteKey(id: $id) {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}"): (typeof documents)["mutation RegenerateSiteKey($id: ID!) {\n  regenerateSiteKey(id: $id) {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "query GeoIPCountries($search: String) {\n  geoIPCountries(search: $search) {\n    code\n    name\n  }\n}"): (typeof documents)["query GeoIPCountries($search: String) {\n  geoIPCountries(search: $search) {\n    code\n    name\n  }\n}"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "query Me {\n  me {\n    id\n    username\n    role\n    createdAt\n    sites {\n      id\n      domains\n      name\n      publicKey\n      createdAt\n    }\n  }\n}\n\nmutation Login($input: LoginInput!) {\n  login(input: $input) {\n    user {\n      id\n      username\n      role\n    }\n  }\n}\n\nmutation Register($input: RegisterInput!) {\n  register(input: $input) {\n    user {\n      id\n      username\n      role\n    }\n  }\n}\n\nmutation Logout {\n  logout\n}"): (typeof documents)["query Me {\n  me {\n    id\n    username\n    role\n    createdAt\n    sites {\n      id\n      domains\n      name\n      publicKey\n      createdAt\n    }\n  }\n}\n\nmutation Login($input: LoginInput!) {\n  login(input: $input) {\n    user {\n      id\n      username\n      role\n    }\n  }\n}\n\nmutation Register($input: RegisterInput!) {\n  register(input: $input) {\n    user {\n      id\n      username\n      role\n    }\n  }\n}\n\nmutation Logout {\n  logout\n}"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "query Dashboard($siteId: ID!, $dateRange: DateRangeInput, $filter: FilterInput) {\n  dashboard(siteId: $siteId, dateRange: $dateRange, filter: $filter) {\n    visitors\n    pageViews\n    sessions\n    bounceRate\n    avgDuration\n    topPages {\n      path\n      views\n      visitors\n    }\n    topReferrers {\n      referrer\n      visitors\n    }\n    browsers {\n      browser\n      visitors\n    }\n    devices {\n      device\n      visitors\n    }\n    countries {\n      country\n      visitors\n    }\n    dailyStats {\n      date\n      visitors\n      pageViews\n      sessions\n    }\n  }\n}\n\nquery Realtime($siteId: ID!) {\n  realtime(siteId: $siteId) {\n    visitors\n    activePages {\n      path\n      visitors\n    }\n  }\n}\n\nquery Events($siteId: ID!, $dateRange: DateRangeInput, $limit: Int, $offset: Int) {\n  events(siteId: $siteId, dateRange: $dateRange, limit: $limit, offset: $offset) {\n    total\n    events {\n      id\n      name\n      path\n      createdAt\n      properties {\n        key\n        value\n      }\n    }\n  }\n}"): (typeof documents)["query Dashboard($siteId: ID!, $dateRange: DateRangeInput, $filter: FilterInput) {\n  dashboard(siteId: $siteId, dateRange: $dateRange, filter: $filter) {\n    visitors\n    pageViews\n    sessions\n    bounceRate\n    avgDuration\n    topPages {\n      path\n      views\n      visitors\n    }\n    topReferrers {\n      referrer\n      visitors\n    }\n    browsers {\n      browser\n      visitors\n    }\n    devices {\n      device\n      visitors\n    }\n    countries {\n      country\n      visitors\n    }\n    dailyStats {\n      date\n      visitors\n      pageViews\n      sessions\n    }\n  }\n}\n\nquery Realtime($siteId: ID!) {\n  realtime(siteId: $siteId) {\n    visitors\n    activePages {\n      path\n      visitors\n    }\n  }\n}\n\nquery Events($siteId: ID!, $dateRange: DateRangeInput, $limit: Int, $offset: Int) {\n  events(siteId: $siteId, dateRange: $dateRange, limit: $limit, offset: $offset) {\n    total\n    events {\n      id\n      name\n      path\n      createdAt\n      properties {\n        key\n        value\n      }\n    }\n  }\n}"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "mutation CreateSite($input: CreateSiteInput!) {\n  createSite(input: $input) {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}\n\nquery GeoIPStatus {\n  geoIPStatus {\n    state\n    dbPath\n    source\n    lastError\n    updatedAt\n  }\n}\n\nquery Site($id: ID!) {\n  site(id: $id) {\n    id\n    domains\n    name\n    publicKey\n    trackCountry\n    blockedIPs\n    blockedCountries\n    createdAt\n  }\n}\n\nmutation UpdateSite($id: ID!, $input: UpdateSiteInput!) {\n  updateSite(id: $id, input: $input) {\n    id\n    domains\n    name\n    publicKey\n    trackCountry\n    blockedIPs\n    blockedCountries\n    createdAt\n  }\n}"): (typeof documents)["mutation CreateSite($input: CreateSiteInput!) {\n  createSite(input: $input) {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}\n\nquery GeoIPStatus {\n  geoIPStatus {\n    state\n    dbPath\n    source\n    lastError\n    updatedAt\n  }\n}\n\nquery Site($id: ID!) {\n  site(id: $id) {\n    id\n    domains\n    name\n    publicKey\n    trackCountry\n    blockedIPs\n    blockedCountries\n    createdAt\n  }\n}\n\nmutation UpdateSite($id: ID!, $input: UpdateSiteInput!) {\n  updateSite(id: $id, input: $input) {\n    id\n    domains\n    name\n    publicKey\n    trackCountry\n    blockedIPs\n    blockedCountries\n    createdAt\n  }\n}"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "query Sites {\n  sites {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}"): (typeof documents)["query Sites {\n  sites {\n    id\n    domains\n    name\n    publicKey\n    createdAt\n  }\n}"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;