import type { FunctionComponent } from 'react';
import CountryCard from '@/components/country-card';
import ReferrersCard from '@/components/referrers-card';
import TopPagesCard from '@/components/top-pages-card';
import type { CountryStatsFieldsFragment, PageStatsFieldsFragment, ReferrerStatsFieldsFragment } from '@/gql/graphql';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface AnalyticsTrafficBreakdownSectionProps {
  siteId: string;
  dashboardState: DashboardLoadState;
  topPages: PageStatsFieldsFragment[];
  topPagesTotal: number;
  topPagesPage: number;
  topPagesPageSize: number;
  onTopPagesPageChange: (page: number) => void;
  referrers: ReferrerStatsFieldsFragment[];
  referrersTotal: number;
  referrersPage: number;
  referrersPageSize: number;
  totalVisitors: number;
  onReferrersPageChange: (page: number) => void;
  countries: CountryStatsFieldsFragment[];
  countriesTotal: number;
  countriesTotalVisitors: number;
  countriesPage: number;
  countriesPageSize: number;
  onCountriesPageChange: (page: number) => void;
}

const AnalyticsTrafficBreakdownSection: FunctionComponent<AnalyticsTrafficBreakdownSectionProps> = (props) => (
  <>
    <div className="grid gap-6 md:grid-cols-2">
      <TopPagesCard pages={props.topPages} total={props.topPagesTotal} page={props.topPagesPage} pageSize={props.topPagesPageSize} siteId={props.siteId} state={props.dashboardState} onPageChange={props.onTopPagesPageChange} />
      <ReferrersCard referrers={props.referrers} totalCount={props.referrersTotal} totalVisitors={props.totalVisitors} siteId={props.siteId} page={props.referrersPage} pageSize={props.referrersPageSize} state={props.dashboardState} onPageChange={props.onReferrersPageChange} />
    </div>
    <CountryCard countries={props.countries} total={props.countriesTotal} totalVisitors={props.countriesTotalVisitors} page={props.countriesPage} pageSize={props.countriesPageSize} siteId={props.siteId} state={props.dashboardState} onPageChange={props.onCountriesPageChange} />
  </>
);

export default AnalyticsTrafficBreakdownSection;
