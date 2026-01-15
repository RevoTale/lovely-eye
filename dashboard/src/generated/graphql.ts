import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
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
  domain: Scalars['String']['input'];
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
  countries: Array<CountryStats>;
  dailyStats: Array<DailyStats>;
  devices: Array<DeviceStats>;
  pageViews: Scalars['Int']['output'];
  sessions: Scalars['Int']['output'];
  topPages: Array<PageStats>;
  topReferrers: Array<ReferrerStats>;
  visitors: Scalars['Int']['output'];
};

export type DateRangeInput = {
  from?: InputMaybe<Scalars['Time']['input']>;
  to?: InputMaybe<Scalars['Time']['input']>;
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
  country?: InputMaybe<Scalars['String']['input']>;
  /** Filter by device type (desktop, mobile, tablet) */
  device?: InputMaybe<Scalars['String']['input']>;
  /** Filter by page path */
  page?: InputMaybe<Scalars['String']['input']>;
  /** Filter by specific referrer */
  referrer?: InputMaybe<Scalars['String']['input']>;
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
};


export type MutationCreateSiteArgs = {
  input: CreateSiteInput;
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

export type PageStats = {
  __typename: 'PageStats';
  path: Scalars['String']['output'];
  views: Scalars['Int']['output'];
  visitors: Scalars['Int']['output'];
};

export type Query = {
  __typename: 'Query';
  dashboard: DashboardStats;
  /** Get events for a site with pagination */
  events: EventsResult;
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


export type QueryEventsArgs = {
  dateRange: InputMaybe<DateRangeInput>;
  limit: InputMaybe<Scalars['Int']['input']>;
  offset: InputMaybe<Scalars['Int']['input']>;
  siteId: Scalars['ID']['input'];
};


export type QueryRealtimeArgs = {
  siteId: Scalars['ID']['input'];
};


export type QuerySiteArgs = {
  id: Scalars['ID']['input'];
};

export type RealtimeStats = {
  __typename: 'RealtimeStats';
  /** Active pages with visitor count */
  activePages: Array<ActivePageStats>;
  /** Visitors active in last 5 minutes */
  visitors: Scalars['Int']['output'];
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

export type Site = {
  __typename: 'Site';
  createdAt: Scalars['Time']['output'];
  domain: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  /** Used in tracking script */
  publicKey: Scalars['String']['output'];
  /** Enable country tracking (requires GeoIP database) */
  trackCountry: Scalars['Boolean']['output'];
};

export type TokenPayload = {
  __typename: 'TokenPayload';
  accessToken: Scalars['String']['output'];
  refreshToken: Scalars['String']['output'];
};

export type UpdateSiteInput = {
  name: Scalars['String']['input'];
  trackCountry?: InputMaybe<Scalars['Boolean']['input']>;
};

export type User = {
  __typename: 'User';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  role: Scalars['String']['output'];
  sites: Maybe<Array<Site>>;
  username: Scalars['String']['output'];
};

export type MeQueryVariables = Exact<{ [key: string]: never; }>;


export type MeQuery = { __typename: 'Query', me: { __typename: 'User', id: string, username: string, role: string, createdAt: string, sites: Array<{ __typename: 'Site', id: string, domain: string, name: string, publicKey: string, createdAt: string }> | null } | null };

export type SitesQueryVariables = Exact<{ [key: string]: never; }>;


export type SitesQuery = { __typename: 'Query', sites: Array<{ __typename: 'Site', id: string, domain: string, name: string, publicKey: string, createdAt: string }> };

export type SiteQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type SiteQuery = { __typename: 'Query', site: { __typename: 'Site', id: string, domain: string, name: string, publicKey: string, trackCountry: boolean, createdAt: string } | null };

export type DashboardQueryVariables = Exact<{
  siteId: Scalars['ID']['input'];
  dateRange: InputMaybe<DateRangeInput>;
  filter: InputMaybe<FilterInput>;
}>;


export type DashboardQuery = { __typename: 'Query', dashboard: { __typename: 'DashboardStats', visitors: number, pageViews: number, sessions: number, bounceRate: number, avgDuration: number, topPages: Array<{ __typename: 'PageStats', path: string, views: number, visitors: number }>, topReferrers: Array<{ __typename: 'ReferrerStats', referrer: string, visitors: number }>, browsers: Array<{ __typename: 'BrowserStats', browser: string, visitors: number }>, devices: Array<{ __typename: 'DeviceStats', device: string, visitors: number }>, countries: Array<{ __typename: 'CountryStats', country: string, visitors: number }>, dailyStats: Array<{ __typename: 'DailyStats', date: string, visitors: number, pageViews: number, sessions: number }> } };

export type RealtimeQueryVariables = Exact<{
  siteId: Scalars['ID']['input'];
}>;


export type RealtimeQuery = { __typename: 'Query', realtime: { __typename: 'RealtimeStats', visitors: number, activePages: Array<{ __typename: 'ActivePageStats', path: string, visitors: number }> } };

export type LoginMutationVariables = Exact<{
  input: LoginInput;
}>;


export type LoginMutation = { __typename: 'Mutation', login: { __typename: 'AuthPayload', user: { __typename: 'User', id: string, username: string, role: string } } };

export type RegisterMutationVariables = Exact<{
  input: RegisterInput;
}>;


export type RegisterMutation = { __typename: 'Mutation', register: { __typename: 'AuthPayload', user: { __typename: 'User', id: string, username: string, role: string } } };

export type LogoutMutationVariables = Exact<{ [key: string]: never; }>;


export type LogoutMutation = { __typename: 'Mutation', logout: boolean };

export type CreateSiteMutationVariables = Exact<{
  input: CreateSiteInput;
}>;


export type CreateSiteMutation = { __typename: 'Mutation', createSite: { __typename: 'Site', id: string, domain: string, name: string, publicKey: string, createdAt: string } };

export type UpdateSiteMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  input: UpdateSiteInput;
}>;


export type UpdateSiteMutation = { __typename: 'Mutation', updateSite: { __typename: 'Site', id: string, domain: string, name: string, publicKey: string, trackCountry: boolean, createdAt: string } };

export type GeoIpStatusQueryVariables = Exact<{ [key: string]: never; }>;


export type GeoIpStatusQuery = { __typename: 'Query', geoIPStatus: { __typename: 'GeoIPStatus', state: string, dbPath: string, source: string | null, lastError: string | null, updatedAt: string | null } };

export type RefreshGeoIpDatabaseMutationVariables = Exact<{ [key: string]: never; }>;


export type RefreshGeoIpDatabaseMutation = { __typename: 'Mutation', refreshGeoIPDatabase: { __typename: 'GeoIPStatus', state: string, dbPath: string, source: string | null, lastError: string | null, updatedAt: string | null } };

export type DeleteSiteMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteSiteMutation = { __typename: 'Mutation', deleteSite: boolean };

export type RegenerateSiteKeyMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type RegenerateSiteKeyMutation = { __typename: 'Mutation', regenerateSiteKey: { __typename: 'Site', id: string, domain: string, name: string, publicKey: string, createdAt: string } };


export const MeDocument = gql`
    query Me {
  me {
    id
    username
    role
    createdAt
    sites {
      id
      domain
      name
      publicKey
      createdAt
    }
  }
}
    `;

/**
 * __useMeQuery__
 *
 * To run a query within a React component, call `useMeQuery` and pass it any options that fit your needs.
 * When your component renders, `useMeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMeQuery({
 *   variables: {
 *   },
 * });
 */
export function useMeQuery(baseOptions?: Apollo.QueryHookOptions<MeQuery, MeQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<MeQuery, MeQueryVariables>(MeDocument, options);
      }
export function useMeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<MeQuery, MeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<MeQuery, MeQueryVariables>(MeDocument, options);
        }
// @ts-ignore
export function useMeSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<MeQuery, MeQueryVariables>): Apollo.UseSuspenseQueryResult<MeQuery, MeQueryVariables>;
export function useMeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MeQuery, MeQueryVariables>): Apollo.UseSuspenseQueryResult<MeQuery | undefined, MeQueryVariables>;
export function useMeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MeQuery, MeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<MeQuery, MeQueryVariables>(MeDocument, options);
        }
export type MeQueryHookResult = ReturnType<typeof useMeQuery>;
export type MeLazyQueryHookResult = ReturnType<typeof useMeLazyQuery>;
export type MeSuspenseQueryHookResult = ReturnType<typeof useMeSuspenseQuery>;
export type MeQueryResult = Apollo.QueryResult<MeQuery, MeQueryVariables>;
export const SitesDocument = gql`
    query Sites {
  sites {
    id
    domain
    name
    publicKey
    createdAt
  }
}
    `;

/**
 * __useSitesQuery__
 *
 * To run a query within a React component, call `useSitesQuery` and pass it any options that fit your needs.
 * When your component renders, `useSitesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSitesQuery({
 *   variables: {
 *   },
 * });
 */
export function useSitesQuery(baseOptions?: Apollo.QueryHookOptions<SitesQuery, SitesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<SitesQuery, SitesQueryVariables>(SitesDocument, options);
      }
export function useSitesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SitesQuery, SitesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<SitesQuery, SitesQueryVariables>(SitesDocument, options);
        }
// @ts-ignore
export function useSitesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<SitesQuery, SitesQueryVariables>): Apollo.UseSuspenseQueryResult<SitesQuery, SitesQueryVariables>;
export function useSitesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SitesQuery, SitesQueryVariables>): Apollo.UseSuspenseQueryResult<SitesQuery | undefined, SitesQueryVariables>;
export function useSitesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SitesQuery, SitesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<SitesQuery, SitesQueryVariables>(SitesDocument, options);
        }
export type SitesQueryHookResult = ReturnType<typeof useSitesQuery>;
export type SitesLazyQueryHookResult = ReturnType<typeof useSitesLazyQuery>;
export type SitesSuspenseQueryHookResult = ReturnType<typeof useSitesSuspenseQuery>;
export type SitesQueryResult = Apollo.QueryResult<SitesQuery, SitesQueryVariables>;
export const SiteDocument = gql`
    query Site($id: ID!) {
  site(id: $id) {
    id
    domain
    name
    publicKey
    trackCountry
    createdAt
  }
}
    `;

/**
 * __useSiteQuery__
 *
 * To run a query within a React component, call `useSiteQuery` and pass it any options that fit your needs.
 * When your component renders, `useSiteQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSiteQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useSiteQuery(baseOptions: Apollo.QueryHookOptions<SiteQuery, SiteQueryVariables> & ({ variables: SiteQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<SiteQuery, SiteQueryVariables>(SiteDocument, options);
      }
export function useSiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SiteQuery, SiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<SiteQuery, SiteQueryVariables>(SiteDocument, options);
        }
// @ts-ignore
export function useSiteSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<SiteQuery, SiteQueryVariables>): Apollo.UseSuspenseQueryResult<SiteQuery, SiteQueryVariables>;
export function useSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SiteQuery, SiteQueryVariables>): Apollo.UseSuspenseQueryResult<SiteQuery | undefined, SiteQueryVariables>;
export function useSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SiteQuery, SiteQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<SiteQuery, SiteQueryVariables>(SiteDocument, options);
        }
export type SiteQueryHookResult = ReturnType<typeof useSiteQuery>;
export type SiteLazyQueryHookResult = ReturnType<typeof useSiteLazyQuery>;
export type SiteSuspenseQueryHookResult = ReturnType<typeof useSiteSuspenseQuery>;
export type SiteQueryResult = Apollo.QueryResult<SiteQuery, SiteQueryVariables>;
export const DashboardDocument = gql`
    query Dashboard($siteId: ID!, $dateRange: DateRangeInput, $filter: FilterInput) {
  dashboard(siteId: $siteId, dateRange: $dateRange, filter: $filter) {
    visitors
    pageViews
    sessions
    bounceRate
    avgDuration
    topPages {
      path
      views
      visitors
    }
    topReferrers {
      referrer
      visitors
    }
    browsers {
      browser
      visitors
    }
    devices {
      device
      visitors
    }
    countries {
      country
      visitors
    }
    dailyStats {
      date
      visitors
      pageViews
      sessions
    }
  }
}
    `;

/**
 * __useDashboardQuery__
 *
 * To run a query within a React component, call `useDashboardQuery` and pass it any options that fit your needs.
 * When your component renders, `useDashboardQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useDashboardQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *      dateRange: // value for 'dateRange'
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useDashboardQuery(baseOptions: Apollo.QueryHookOptions<DashboardQuery, DashboardQueryVariables> & ({ variables: DashboardQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<DashboardQuery, DashboardQueryVariables>(DashboardDocument, options);
      }
export function useDashboardLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<DashboardQuery, DashboardQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<DashboardQuery, DashboardQueryVariables>(DashboardDocument, options);
        }
// @ts-ignore
export function useDashboardSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<DashboardQuery, DashboardQueryVariables>): Apollo.UseSuspenseQueryResult<DashboardQuery, DashboardQueryVariables>;
export function useDashboardSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<DashboardQuery, DashboardQueryVariables>): Apollo.UseSuspenseQueryResult<DashboardQuery | undefined, DashboardQueryVariables>;
export function useDashboardSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<DashboardQuery, DashboardQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<DashboardQuery, DashboardQueryVariables>(DashboardDocument, options);
        }
export type DashboardQueryHookResult = ReturnType<typeof useDashboardQuery>;
export type DashboardLazyQueryHookResult = ReturnType<typeof useDashboardLazyQuery>;
export type DashboardSuspenseQueryHookResult = ReturnType<typeof useDashboardSuspenseQuery>;
export type DashboardQueryResult = Apollo.QueryResult<DashboardQuery, DashboardQueryVariables>;
export const RealtimeDocument = gql`
    query Realtime($siteId: ID!) {
  realtime(siteId: $siteId) {
    visitors
    activePages {
      path
      visitors
    }
  }
}
    `;

/**
 * __useRealtimeQuery__
 *
 * To run a query within a React component, call `useRealtimeQuery` and pass it any options that fit your needs.
 * When your component renders, `useRealtimeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useRealtimeQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useRealtimeQuery(baseOptions: Apollo.QueryHookOptions<RealtimeQuery, RealtimeQueryVariables> & ({ variables: RealtimeQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<RealtimeQuery, RealtimeQueryVariables>(RealtimeDocument, options);
      }
export function useRealtimeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<RealtimeQuery, RealtimeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<RealtimeQuery, RealtimeQueryVariables>(RealtimeDocument, options);
        }
// @ts-ignore
export function useRealtimeSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<RealtimeQuery, RealtimeQueryVariables>): Apollo.UseSuspenseQueryResult<RealtimeQuery, RealtimeQueryVariables>;
export function useRealtimeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<RealtimeQuery, RealtimeQueryVariables>): Apollo.UseSuspenseQueryResult<RealtimeQuery | undefined, RealtimeQueryVariables>;
export function useRealtimeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<RealtimeQuery, RealtimeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<RealtimeQuery, RealtimeQueryVariables>(RealtimeDocument, options);
        }
export type RealtimeQueryHookResult = ReturnType<typeof useRealtimeQuery>;
export type RealtimeLazyQueryHookResult = ReturnType<typeof useRealtimeLazyQuery>;
export type RealtimeSuspenseQueryHookResult = ReturnType<typeof useRealtimeSuspenseQuery>;
export type RealtimeQueryResult = Apollo.QueryResult<RealtimeQuery, RealtimeQueryVariables>;
export const LoginDocument = gql`
    mutation Login($input: LoginInput!) {
  login(input: $input) {
    user {
      id
      username
      role
    }
  }
}
    `;
export type LoginMutationFn = Apollo.MutationFunction<LoginMutation, LoginMutationVariables>;

/**
 * __useLoginMutation__
 *
 * To run a mutation, you first call `useLoginMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useLoginMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [loginMutation, { data, loading, error }] = useLoginMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useLoginMutation(baseOptions?: Apollo.MutationHookOptions<LoginMutation, LoginMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<LoginMutation, LoginMutationVariables>(LoginDocument, options);
      }
export type LoginMutationHookResult = ReturnType<typeof useLoginMutation>;
export type LoginMutationResult = Apollo.MutationResult<LoginMutation>;
export type LoginMutationOptions = Apollo.BaseMutationOptions<LoginMutation, LoginMutationVariables>;
export const RegisterDocument = gql`
    mutation Register($input: RegisterInput!) {
  register(input: $input) {
    user {
      id
      username
      role
    }
  }
}
    `;
export type RegisterMutationFn = Apollo.MutationFunction<RegisterMutation, RegisterMutationVariables>;

/**
 * __useRegisterMutation__
 *
 * To run a mutation, you first call `useRegisterMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRegisterMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [registerMutation, { data, loading, error }] = useRegisterMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useRegisterMutation(baseOptions?: Apollo.MutationHookOptions<RegisterMutation, RegisterMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RegisterMutation, RegisterMutationVariables>(RegisterDocument, options);
      }
export type RegisterMutationHookResult = ReturnType<typeof useRegisterMutation>;
export type RegisterMutationResult = Apollo.MutationResult<RegisterMutation>;
export type RegisterMutationOptions = Apollo.BaseMutationOptions<RegisterMutation, RegisterMutationVariables>;
export const LogoutDocument = gql`
    mutation Logout {
  logout
}
    `;
export type LogoutMutationFn = Apollo.MutationFunction<LogoutMutation, LogoutMutationVariables>;

/**
 * __useLogoutMutation__
 *
 * To run a mutation, you first call `useLogoutMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useLogoutMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [logoutMutation, { data, loading, error }] = useLogoutMutation({
 *   variables: {
 *   },
 * });
 */
export function useLogoutMutation(baseOptions?: Apollo.MutationHookOptions<LogoutMutation, LogoutMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<LogoutMutation, LogoutMutationVariables>(LogoutDocument, options);
      }
export type LogoutMutationHookResult = ReturnType<typeof useLogoutMutation>;
export type LogoutMutationResult = Apollo.MutationResult<LogoutMutation>;
export type LogoutMutationOptions = Apollo.BaseMutationOptions<LogoutMutation, LogoutMutationVariables>;
export const CreateSiteDocument = gql`
    mutation CreateSite($input: CreateSiteInput!) {
  createSite(input: $input) {
    id
    domain
    name
    publicKey
    createdAt
  }
}
    `;
export type CreateSiteMutationFn = Apollo.MutationFunction<CreateSiteMutation, CreateSiteMutationVariables>;

/**
 * __useCreateSiteMutation__
 *
 * To run a mutation, you first call `useCreateSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createSiteMutation, { data, loading, error }] = useCreateSiteMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateSiteMutation(baseOptions?: Apollo.MutationHookOptions<CreateSiteMutation, CreateSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateSiteMutation, CreateSiteMutationVariables>(CreateSiteDocument, options);
      }
export type CreateSiteMutationHookResult = ReturnType<typeof useCreateSiteMutation>;
export type CreateSiteMutationResult = Apollo.MutationResult<CreateSiteMutation>;
export type CreateSiteMutationOptions = Apollo.BaseMutationOptions<CreateSiteMutation, CreateSiteMutationVariables>;
export const UpdateSiteDocument = gql`
    mutation UpdateSite($id: ID!, $input: UpdateSiteInput!) {
  updateSite(id: $id, input: $input) {
    id
    domain
    name
    publicKey
    trackCountry
    createdAt
  }
}
    `;
export type UpdateSiteMutationFn = Apollo.MutationFunction<UpdateSiteMutation, UpdateSiteMutationVariables>;

/**
 * __useUpdateSiteMutation__
 *
 * To run a mutation, you first call `useUpdateSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSiteMutation, { data, loading, error }] = useUpdateSiteMutation({
 *   variables: {
 *      id: // value for 'id'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateSiteMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSiteMutation, UpdateSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateSiteMutation, UpdateSiteMutationVariables>(UpdateSiteDocument, options);
      }
export type UpdateSiteMutationHookResult = ReturnType<typeof useUpdateSiteMutation>;
export type UpdateSiteMutationResult = Apollo.MutationResult<UpdateSiteMutation>;
export type UpdateSiteMutationOptions = Apollo.BaseMutationOptions<UpdateSiteMutation, UpdateSiteMutationVariables>;
export const GeoIpStatusDocument = gql`
    query GeoIPStatus {
  geoIPStatus {
    state
    dbPath
    source
    lastError
    updatedAt
  }
}
    `;

/**
 * __useGeoIpStatusQuery__
 *
 * To run a query within a React component, call `useGeoIpStatusQuery` and pass it any options that fit your needs.
 * When your component renders, `useGeoIpStatusQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGeoIpStatusQuery({
 *   variables: {
 *   },
 * });
 */
export function useGeoIpStatusQuery(baseOptions?: Apollo.QueryHookOptions<GeoIpStatusQuery, GeoIpStatusQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GeoIpStatusQuery, GeoIpStatusQueryVariables>(GeoIpStatusDocument, options);
      }
export function useGeoIpStatusLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GeoIpStatusQuery, GeoIpStatusQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GeoIpStatusQuery, GeoIpStatusQueryVariables>(GeoIpStatusDocument, options);
        }
// @ts-ignore
export function useGeoIpStatusSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GeoIpStatusQuery, GeoIpStatusQueryVariables>): Apollo.UseSuspenseQueryResult<GeoIpStatusQuery, GeoIpStatusQueryVariables>;
export function useGeoIpStatusSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GeoIpStatusQuery, GeoIpStatusQueryVariables>): Apollo.UseSuspenseQueryResult<GeoIpStatusQuery | undefined, GeoIpStatusQueryVariables>;
export function useGeoIpStatusSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GeoIpStatusQuery, GeoIpStatusQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GeoIpStatusQuery, GeoIpStatusQueryVariables>(GeoIpStatusDocument, options);
        }
export type GeoIpStatusQueryHookResult = ReturnType<typeof useGeoIpStatusQuery>;
export type GeoIpStatusLazyQueryHookResult = ReturnType<typeof useGeoIpStatusLazyQuery>;
export type GeoIpStatusSuspenseQueryHookResult = ReturnType<typeof useGeoIpStatusSuspenseQuery>;
export type GeoIpStatusQueryResult = Apollo.QueryResult<GeoIpStatusQuery, GeoIpStatusQueryVariables>;
export const RefreshGeoIpDatabaseDocument = gql`
    mutation RefreshGeoIPDatabase {
  refreshGeoIPDatabase {
    state
    dbPath
    source
    lastError
    updatedAt
  }
}
    `;
export type RefreshGeoIpDatabaseMutationFn = Apollo.MutationFunction<RefreshGeoIpDatabaseMutation, RefreshGeoIpDatabaseMutationVariables>;

/**
 * __useRefreshGeoIpDatabaseMutation__
 *
 * To run a mutation, you first call `useRefreshGeoIpDatabaseMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRefreshGeoIpDatabaseMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [refreshGeoIpDatabaseMutation, { data, loading, error }] = useRefreshGeoIpDatabaseMutation({
 *   variables: {
 *   },
 * });
 */
export function useRefreshGeoIpDatabaseMutation(baseOptions?: Apollo.MutationHookOptions<RefreshGeoIpDatabaseMutation, RefreshGeoIpDatabaseMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RefreshGeoIpDatabaseMutation, RefreshGeoIpDatabaseMutationVariables>(RefreshGeoIpDatabaseDocument, options);
      }
export type RefreshGeoIpDatabaseMutationHookResult = ReturnType<typeof useRefreshGeoIpDatabaseMutation>;
export type RefreshGeoIpDatabaseMutationResult = Apollo.MutationResult<RefreshGeoIpDatabaseMutation>;
export type RefreshGeoIpDatabaseMutationOptions = Apollo.BaseMutationOptions<RefreshGeoIpDatabaseMutation, RefreshGeoIpDatabaseMutationVariables>;
export const DeleteSiteDocument = gql`
    mutation DeleteSite($id: ID!) {
  deleteSite(id: $id)
}
    `;
export type DeleteSiteMutationFn = Apollo.MutationFunction<DeleteSiteMutation, DeleteSiteMutationVariables>;

/**
 * __useDeleteSiteMutation__
 *
 * To run a mutation, you first call `useDeleteSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteSiteMutation, { data, loading, error }] = useDeleteSiteMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useDeleteSiteMutation(baseOptions?: Apollo.MutationHookOptions<DeleteSiteMutation, DeleteSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteSiteMutation, DeleteSiteMutationVariables>(DeleteSiteDocument, options);
      }
export type DeleteSiteMutationHookResult = ReturnType<typeof useDeleteSiteMutation>;
export type DeleteSiteMutationResult = Apollo.MutationResult<DeleteSiteMutation>;
export type DeleteSiteMutationOptions = Apollo.BaseMutationOptions<DeleteSiteMutation, DeleteSiteMutationVariables>;
export const RegenerateSiteKeyDocument = gql`
    mutation RegenerateSiteKey($id: ID!) {
  regenerateSiteKey(id: $id) {
    id
    domain
    name
    publicKey
    createdAt
  }
}
    `;
export type RegenerateSiteKeyMutationFn = Apollo.MutationFunction<RegenerateSiteKeyMutation, RegenerateSiteKeyMutationVariables>;

/**
 * __useRegenerateSiteKeyMutation__
 *
 * To run a mutation, you first call `useRegenerateSiteKeyMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRegenerateSiteKeyMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [regenerateSiteKeyMutation, { data, loading, error }] = useRegenerateSiteKeyMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useRegenerateSiteKeyMutation(baseOptions?: Apollo.MutationHookOptions<RegenerateSiteKeyMutation, RegenerateSiteKeyMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RegenerateSiteKeyMutation, RegenerateSiteKeyMutationVariables>(RegenerateSiteKeyDocument, options);
      }
export type RegenerateSiteKeyMutationHookResult = ReturnType<typeof useRegenerateSiteKeyMutation>;
export type RegenerateSiteKeyMutationResult = Apollo.MutationResult<RegenerateSiteKeyMutation>;
export type RegenerateSiteKeyMutationOptions = Apollo.BaseMutationOptions<RegenerateSiteKeyMutation, RegenerateSiteKeyMutationVariables>;