import { parseChartDate } from '@/features/analytics/model/chart-date';
import type { OverviewPoint } from '@/features/analytics/ui/overview-chart/overview-chart-series';
import type { DailyStatsFieldsFragment } from '@/shared/api/generated/graphql';

export const buildOverviewChartData = (stats: DailyStatsFieldsFragment[]): OverviewPoint[] =>
  stats
    .map((stat) => {
      const timestamp = parseChartDate(stat.date);
      if (timestamp === null) {
        return null;
      }

      return {
        timestamp,
        visitors: stat.visitors,
        pageViews: stat.pageViews,
        sessions: stat.sessions,
      };
    })
    .filter((point): point is OverviewPoint => point !== null)
    .sort((left, right) => left.timestamp - right.timestamp);
