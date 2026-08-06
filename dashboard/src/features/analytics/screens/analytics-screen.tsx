import type { FunctionComponent } from 'react';
import { isInitialLoadState } from '@/features/analytics/model/dashboard-load-state';
import { useDashboardData } from '@/features/analytics/model/use-dashboard-data';
import { useDashboardPageState } from '@/features/analytics/model/use-dashboard-page-state';
import { ActiveFilters } from '@/features/analytics/ui/active-filters';
import AnalyticsContent from '@/features/analytics/ui/analytics-content';
import { DashboardHeader } from '@/features/analytics/ui/dashboard-header';
import {
  DashboardEmptyState,
  DashboardLoading,
  DashboardNotFound,
} from '@/features/analytics/ui/dashboard-states';
import { TimeRangeCard } from '@/features/analytics/ui/time-range-card';

export const AnalyticsScreen: FunctionComponent = () => {
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
  if (dashboard.stats === undefined && dashboard.dashboardState === 'ready')
    return <DashboardEmptyState />;

  return (
    <div className='space-y-8'>
      <DashboardHeader site={dashboard.site} siteId={state.siteId} realtime={dashboard.realtime} />
      <TimeRangeCard
        preset={state.preset}
        fromDate={state.fromDate}
        toDate={state.toDate}
        fromTime={state.fromTime}
        toTime={state.toTime}
        onPresetChange={state.setPreset}
        onApplyRange={state.applyCustomRange}
      />
      <ActiveFilters siteId={state.siteId} search={state.decodedSearch} />
      <AnalyticsContent dashboard={dashboard} state={state} />
    </div>
  );
};

export default AnalyticsScreen;
