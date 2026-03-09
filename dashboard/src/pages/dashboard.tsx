import type { FunctionComponent } from 'react';
import { ActiveFilters } from '@/components/active-filters';
import AnalyticsContent from '@/components/analytics-content';
import { DashboardHeader } from '@/components/dashboard-header';
import { DashboardEmptyState, DashboardLoading, DashboardNotFound } from '@/components/dashboard-states';
import { TimeRangeCard } from '@/components/time-range-card';
import { useDashboardData, PAGE_SIZES } from '@/hooks/use-dashboard-data';
import { useDashboardPageState } from '@/hooks/use-dashboard-page-state';
import { isInitialLoadState } from '@/lib/dashboard-load-state';

const DashboardPage: FunctionComponent = () => {
  const state = useDashboardPageState();
  const dashboard = useDashboardData({
    siteId: state.siteId,
    dateRange: state.dateRange,
    filter: state.filterInput,
    eventsPage: state.eventsPage,
    eventsCountsPage: state.eventsCountsPage,
    topPagesPage: state.topPagesPage,
    referrersPage: state.referrersPage,
    devicesPage: state.devicesPage,
    osPage: state.osPage,
    countriesPage: state.countriesPage,
  });

  if (isInitialLoadState(dashboard.siteState)) return <DashboardLoading />;
  if (dashboard.site === null || dashboard.site === undefined) return <DashboardNotFound />;
  if (dashboard.stats === undefined && dashboard.dashboardState === 'ready') return <DashboardEmptyState />;

  return (
    <div className="space-y-8">
      <DashboardHeader site={dashboard.site} siteId={state.siteId} realtime={dashboard.realtime} />
      <TimeRangeCard preset={state.preset} fromDate={state.fromDate} toDate={state.toDate} fromTime={state.fromTime} toTime={state.toTime} onPresetChange={state.setPreset} onApplyRange={state.applyCustomRange} />
      <ActiveFilters siteId={state.siteId} search={state.decodedSearch} />
      <AnalyticsContent siteId={state.siteId} stats={dashboard.stats} dashboardState={dashboard.dashboardState} realtime={dashboard.realtime} dateRange={state.dateRangeForChart} filter={state.filterInput} chartBucket={state.statsBucket} onChartBucketChange={state.setStatsBucket} eventsResult={dashboard.eventsResult} eventsState={dashboard.eventsState} eventCounts={dashboard.eventCounts} eventCountsState={dashboard.eventCountsState} eventsPage={state.eventsPage} eventsCountsPage={state.eventsCountsPage} topPagesPage={state.topPagesPage} referrersPage={state.referrersPage} devicesPage={state.devicesPage} osPage={state.osPage} countriesPage={state.countriesPage} onPageChange={state.setPage} pageSizes={PAGE_SIZES} />
    </div>
  );
};

export default DashboardPage;
export { DashboardPage };
