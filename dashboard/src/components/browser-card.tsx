import { Compass } from 'lucide-react';
import type { FunctionComponent } from 'react';
import BoardCard from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';
import { Badge, Progress, Skeleton } from '@/components/ui';
import type { BrowserStatsFieldsFragment } from '@/gql/graphql';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface BrowserCardProps {
  browsers: BrowserStatsFieldsFragment[];
  totalVisitors: number;
  siteId: string;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;
const PERCENT_MULTIPLIER = 100;
const PERCENT_PRECISION = 1;
const SKELETON_ROWS = 5;

const BrowserCard: FunctionComponent<BrowserCardProps> = ({ browsers, totalVisitors, siteId, state = 'ready' }) => (
  <BoardCard
    title="Browsers"
    icon={Compass}
    state={state}
    overlayLabel="Refreshing browsers"
    skeleton={<div className="space-y-3">{Array.from({ length: SKELETON_ROWS }, (_, index) => <div key={index} className="space-y-2"><div className="flex items-center justify-between"><Skeleton className="h-4 w-28" /><Skeleton className="h-5 w-16" /></div><Skeleton className="h-2 w-full" /></div>)}</div>}
  >
    {browsers.length > EMPTY_COUNT ? (
      <div className="space-y-3">
        {browsers.map((browserStat) => {
          const percentage = totalVisitors > EMPTY_COUNT ? (browserStat.visitors / totalVisitors) * PERCENT_MULTIPLIER : EMPTY_COUNT;
          return (
            <div key={browserStat.browser}>
              <div className="mb-1 flex items-center justify-between gap-2">
                <FilterLink siteId={siteId} filterKey="browser" value={browserStat.browser} className="truncate text-sm font-medium hover:text-primary hover:underline">
                  {browserStat.browser}
                </FilterLink>
                <div className="flex items-center gap-2">
                  <Badge variant="secondary">{browserStat.visitors.toLocaleString()}</Badge>
                  <span className="text-xs text-muted-foreground">{percentage.toFixed(PERCENT_PRECISION)}%</span>
                </div>
              </div>
              <Progress value={percentage} className="h-2" />
            </div>
          );
        })}
      </div>
    ) : <ListEmptyState title="No browser data yet" />}
  </BoardCard>
);

export default BrowserCard;
