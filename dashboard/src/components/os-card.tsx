import { Laptop } from 'lucide-react';
import type { FunctionComponent } from 'react';
import BoardCard from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';
import { Badge, Progress, Skeleton } from '@/components/ui';
import type { OperatingSystemStatsFieldsFragment } from '@/gql/graphql';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface OSCardProps {
  operatingSystems: OperatingSystemStatsFieldsFragment[];
  total: number;
  totalVisitors: number;
  page: number;
  pageSize: number;
  siteId: string;
  onPageChange: (page: number) => void;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;
const PERCENT_MULTIPLIER = 100;
const PERCENT_PRECISION = 1;

const OSCard: FunctionComponent<OSCardProps> = ({ operatingSystems, total, totalVisitors, page, pageSize, siteId, onPageChange, state = 'ready' }) => (
  <BoardCard
    title="Operating Systems"
    icon={Laptop}
    state={state}
    pagination={{ page, pageSize, total, onPageChange }}
    overlayLabel="Refreshing operating systems"
    skeleton={<div className="space-y-3">{Array.from({ length: 5 }, (_, index) => <div key={index} className="space-y-2"><div className="flex items-center justify-between"><Skeleton className="h-4 w-28" /><Skeleton className="h-5 w-16" /></div><Skeleton className="h-2 w-full" /></div>)}</div>}
  >
    {operatingSystems.length > EMPTY_COUNT ? (
      <div className="space-y-3">
        {operatingSystems.map((operatingSystem) => {
          const percentage = totalVisitors > EMPTY_COUNT ? (operatingSystem.visitors / totalVisitors) * PERCENT_MULTIPLIER : EMPTY_COUNT;
          return (
            <div key={operatingSystem.os}>
              <div className="mb-1 flex items-center justify-between gap-2">
                <FilterLink siteId={siteId} filterKey="os" value={operatingSystem.os} className="truncate text-sm font-medium hover:text-primary hover:underline">
                  {operatingSystem.os}
                </FilterLink>
                <div className="flex items-center gap-2">
                  <Badge variant="secondary">{operatingSystem.visitors.toLocaleString()}</Badge>
                  <span className="text-xs text-muted-foreground">{percentage.toFixed(PERCENT_PRECISION)}%</span>
                </div>
              </div>
              <Progress value={percentage} className="h-2" />
            </div>
          );
        })}
      </div>
    ) : <ListEmptyState title="No operating system data yet" />}
  </BoardCard>
);

export default OSCard;
