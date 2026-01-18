import { gql } from '@apollo/client';

export const ME_QUERY = gql`
  query Me {
    me {
      id
      username
      role
      createdAt
      sites {
        id
        domains
        name
        publicKey
        createdAt
      }
    }
  }
`;

export const SITES_QUERY = gql`
  query Sites {
    sites {
      id
      domains
      name
      publicKey
      createdAt
    }
  }
`;

export const SITE_QUERY = gql`
  query Site($id: ID!) {
    site(id: $id) {
      id
      domains
      name
      publicKey
      trackCountry
      blockedIPs
      blockedCountries
      createdAt
    }
  }
`;

export const DASHBOARD_QUERY = gql`
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

export const REALTIME_QUERY = gql`
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

export const EVENTS_QUERY = gql`
  query Events($siteId: ID!, $dateRange: DateRangeInput, $limit: Int, $offset: Int) {
    events(siteId: $siteId, dateRange: $dateRange, limit: $limit, offset: $offset) {
      total
      events {
        id
        name
        path
        createdAt
        properties {
          key
          value
        }
      }
    }
  }
`;

export const LOGIN_MUTATION = gql`
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

export const REGISTER_MUTATION = gql`
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

export const LOGOUT_MUTATION = gql`
  mutation Logout {
    logout
  }
`;

export const CREATE_SITE_MUTATION = gql`
  mutation CreateSite($input: CreateSiteInput!) {
    createSite(input: $input) {
      id
      domains
      name
      publicKey
      createdAt
    }
  }
`;

export const UPDATE_SITE_MUTATION = gql`
  mutation UpdateSite($id: ID!, $input: UpdateSiteInput!) {
    updateSite(id: $id, input: $input) {
      id
      domains
      name
      publicKey
      trackCountry
      blockedIPs
      blockedCountries
      createdAt
    }
  }
`;

export const EVENT_DEFINITIONS_QUERY = gql`
  query EventDefinitions($siteId: ID!) {
    eventDefinitions(siteId: $siteId) {
      id
      name
      createdAt
      updatedAt
      fields {
        id
        key
        type
        required
        maxLength
      }
    }
  }
`;

export const UPSERT_EVENT_DEFINITION_MUTATION = gql`
  mutation UpsertEventDefinition($siteId: ID!, $input: EventDefinitionInput!) {
    upsertEventDefinition(siteId: $siteId, input: $input) {
      id
      name
      createdAt
      updatedAt
      fields {
        id
        key
        type
        required
        maxLength
      }
    }
  }
`;

export const DELETE_EVENT_DEFINITION_MUTATION = gql`
  mutation DeleteEventDefinition($siteId: ID!, $name: String!) {
    deleteEventDefinition(siteId: $siteId, name: $name)
  }
`;

export const GEOIP_STATUS_QUERY = gql`
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

export const GEOIP_COUNTRIES_QUERY = gql`
  query GeoIPCountries($search: String) {
    geoIPCountries(search: $search) {
      code
      name
    }
  }
`;

export const REFRESH_GEOIP_MUTATION = gql`
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

export const DELETE_SITE_MUTATION = gql`
  mutation DeleteSite($id: ID!) {
    deleteSite(id: $id)
  }
`;

export const REGENERATE_SITE_KEY_MUTATION = gql`
  mutation RegenerateSiteKey($id: ID!) {
    regenerateSiteKey(id: $id) {
      id
      domains
      name
      publicKey
      createdAt
    }
  }
`;
