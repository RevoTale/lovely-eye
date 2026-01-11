// GraphQL operations for the dashboard
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
        domain
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
      domain
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
      domain
      name
      publicKey
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
      domain
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
      domain
      name
      publicKey
      createdAt
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
      domain
      name
      publicKey
      createdAt
    }
  }
`;
