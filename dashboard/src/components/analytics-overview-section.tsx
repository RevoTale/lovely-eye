import { Clock, Eye, TrendingDown, Users } from 'lucide-react';
import type { FunctionComponent } from 'react';
import { ActivePagesCard } from '@/components/active-pages-card';
import OverviewChartSection from '@/components/overview-chart-section';
import StatCard from '@/components/stat-card';
import { RealtimeStatsFieldsFragmentDoc, type DashboardStatsFieldsFragment, type FilterInput, type RealtimeQuery } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';
import { formatDuration } from '@/lib/dashboard-utils';

interface AnalyticsOverviewSectionProps {
  siteId: string;
  stats: DashboardStatsFieldsFragment;
  dashboardState: DashboardLoadState;
  realtime: RealtimeQuery['realtime'] | undefined;
  dateRange: { from: Date; to: Date } | null;
  filter: FilterInput | null;
  chartBucket: 'daily' | 'hourly';
  onChartBucketChange: (bucket: 'daily' | 'hourly') => void;
}

const AnalyticsOverviewSection: FunctionComponent<AnalyticsOverviewSectionProps> = ({ siteId, stats, dashboardState, realtime, dateRange, filter, chartBucket, onChartBucketChange }) => {
  const realtimeData = realtime === undefined ? undefined : getFragmentData(RealtimeStatsFieldsFragmentDoc, realtime);

  return (
    <>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard title="Total Visitors" value={stats.visitors.toLocaleString()} icon={Users} state={dashboardState} />
        <StatCard title="Page Views" value={stats.pageViews.toLocaleString()} icon={Eye} state={dashboardState} />
        <StatCard title="Avg. Session" value={formatDuration(stats.avgDuration)} icon={Clock} state={dashboardState} />
        <StatCard title="Bounce Rate" value={`${String(Math.round(stats.bounceRate))}%`} icon={TrendingDown} state={dashboardState} />
      </div>
      <OverviewChartSection siteId={siteId} dateRange={dateRange} filter={filter} bucket={chartBucket} onBucketChange={onChartBucketChange} />
      {realtimeData?.activePages !== undefined ? <ActivePagesCard activePages={realtimeData.activePages} /> : null}
    </>
  );
};

export default AnalyticsOverviewSection;
