
import { Badge, Progress } from '@/components/ui';
import type { BrowserStatsFieldsFragment } from '@/gql/graphql';
import { Compass } from 'lucide-react';
import { BoardCard, BoardCardSkeleton } from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';

interface BrowserCardProps {
  browsers: BrowserStatsFieldsFragment[];
  totalVisitors: number;
  siteId: string;
  loading?: boolean;
}

const EMPTY_COUNT = 0;
const PERCENT_MULTIPLIER = 100;
const PERCENT_PRECISION = 1;

export const BrowserCard = ({
  browsers,
  totalVisitors,
  siteId,
  loading = false,
}: BrowserCardProps): React.ReactNode => {
  if (loading) {
    return <BoardCardSkeleton title="Browsers" icon={Compass} />;
  }

  return (
    <BoardCard title="Browsers" icon={Compass}>
      {browsers.length > EMPTY_COUNT ? (
        <div className="space-y-3">
          {browsers.map((browserStat) => {
            const percentage =
              totalVisitors > EMPTY_COUNT
                ? (browserStat.visitors / totalVisitors) * PERCENT_MULTIPLIER
                : EMPTY_COUNT;

            return (
              <div key={browserStat.browser}>
                <div className="mb-1 flex items-center justify-between gap-2">
                  <FilterLink
                    siteId={siteId}
                    filterKey="browser"
                    value={browserStat.browser}
                    className="truncate text-sm font-medium hover:text-primary hover:underline"
                  >
                    {browserStat.browser}
                  </FilterLink>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary">
                      {browserStat.visitors.toLocaleString()}
                    </Badge>
                    <span className="text-xs text-muted-foreground">
                      {percentage.toFixed(PERCENT_PRECISION)}%
                    </span>
                  </div>
                </div>
                <Progress value={percentage} className="h-2" />
              </div>
            );
          })}
        </div>
      ) : (
        <ListEmptyState title="No browser data yet" />
      )}
    </BoardCard>
  );
}
