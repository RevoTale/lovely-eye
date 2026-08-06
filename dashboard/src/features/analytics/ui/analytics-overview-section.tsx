import { Clock, Eye, TrendingDown, Users } from 'lucide-react';
import type { FunctionComponent } from 'react';
import type { DashboardLoadState } from '@/features/analytics/model/dashboard-load-state';
import { formatDuration } from '@/features/analytics/model/dashboard-utils';
import { ActivePagesCard } from '@/features/analytics/ui/active-pages-card';
import OverviewChartSection from '@/features/analytics/ui/overview-chart-section';
import StatCard from '@/features/analytics/ui/stat-card';
import { useFragment as getFragmentData } from '@/shared/api/generated/fragment-masking';
import {
  type DashboardStatsFieldsFragment,
  type FilterInput,
  type RealtimeQuery,
  RealtimeStatsFieldsFragmentDoc,
} from '@/shared/api/generated/graphql';

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

const AnalyticsOverviewSection: FunctionComponent<AnalyticsOverviewSectionProps> = ({
  siteId,
  stats,
  dashboardState,
  realtime,
  dateRange,
  filter,
  chartBucket,
  onChartBucketChange,
}) => {
  const realtimeData =
    realtime === undefined ? undefined : getFragmentData(RealtimeStatsFieldsFragmentDoc, realtime);

  return (
    <>
      <div className='grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4'>
        <StatCard
          title='Total Visitors'
          value={stats.visitors.toLocaleString()}
          icon={Users}
          state={dashboardState}
        />
        <StatCard
          title='Page Views'
          value={stats.pageViews.toLocaleString()}
          icon={Eye}
          state={dashboardState}
        />
        <StatCard
          title='Avg. Session'
          value={formatDuration(stats.avgDuration)}
          icon={Clock}
          state={dashboardState}
        />
        <StatCard
          title='Bounce Rate'
          value={`${String(Math.round(stats.bounceRate))}%`}
          icon={TrendingDown}
          state={dashboardState}
        />
      </div>
      <OverviewChartSection
        siteId={siteId}
        dateRange={dateRange}
        filter={filter}
        bucket={chartBucket}
        onBucketChange={onChartBucketChange}
      />
      {realtimeData?.activePages !== undefined ? (
        <ActivePagesCard activePages={realtimeData.activePages} />
      ) : null}
    </>
  );
};

export default AnalyticsOverviewSection;
