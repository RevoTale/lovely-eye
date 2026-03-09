import type { DailyStatsFieldsFragment } from '@/gql/graphql';
import { parseChartDate } from '@/lib/chart-date';
import type { OverviewPoint } from '@/components/overview-chart/overview-chart-series';

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
