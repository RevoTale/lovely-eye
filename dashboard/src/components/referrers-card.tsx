import { ExternalLink, Globe, TrendingUp } from 'lucide-react';
import type { FunctionComponent } from 'react';
import BoardCard from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';
import { Badge, Progress, Skeleton } from '@/components/ui';
import type { ReferrerStatsFieldsFragment } from '@/gql/graphql';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';
import { formatReferrer, getReferrerIcon } from '@/lib/referrer-utils';

interface ReferrersCardProps {
  referrers: ReferrerStatsFieldsFragment[];
  totalCount: number;
  totalVisitors: number;
  siteId: string;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;
const FALLBACK_MAX_VISITORS = 1;
const PERCENT_MULTIPLIER = 100;
const PERCENT_PRECISION = 1;
const ZERO_PERCENT = '0.0';
const DIRECT_LABEL = '(direct)';

const ReferrersCard: FunctionComponent<ReferrersCardProps> = ({ referrers, totalCount, totalVisitors, siteId, page, pageSize, onPageChange, state = 'ready' }) => {
  const maxVisitors = referrers.length > EMPTY_COUNT ? referrers[0]?.visitors ?? FALLBACK_MAX_VISITORS : FALLBACK_MAX_VISITORS;

  return (
    <BoardCard
      title="Top Referrers"
      icon={Globe}
      state={state}
      headerRight={<Badge variant="secondary" className="flex items-center gap-1"><TrendingUp className="h-3 w-3" />{totalCount}</Badge>}
      pagination={{ page, pageSize, total: totalCount, onPageChange }}
      overlayLabel="Refreshing referrers"
      skeleton={<div className="space-y-4">{Array.from({ length: 5 }, (_, index) => <div key={index} className="space-y-2"><div className="flex items-center justify-between"><Skeleton className="h-5 w-32" /><Skeleton className="h-5 w-16" /></div><Skeleton className="h-2 w-full" /></div>)}</div>}
    >
      <div className="space-y-4">
        {referrers.length > EMPTY_COUNT ? referrers.map((ref) => {
          const percentage = totalVisitors > EMPTY_COUNT ? ((ref.visitors / totalVisitors) * PERCENT_MULTIPLIER).toFixed(PERCENT_PRECISION) : ZERO_PERCENT;
          const barWidth = maxVisitors > EMPTY_COUNT ? (ref.visitors / maxVisitors) * PERCENT_MULTIPLIER : EMPTY_COUNT;
          return (
            <div key={`${ref.referrer}-${ref.visitors}`} className="space-y-2">
              <div className="flex items-center justify-between">
                <div className="flex min-w-0 flex-1 items-center gap-2">
                  <span className="text-lg" role="img" aria-label="referrer icon">{getReferrerIcon(ref.referrer)}</span>
                  <FilterLink siteId={siteId} filterKey="referrer" value={ref.referrer === '' ? DIRECT_LABEL : ref.referrer} className="truncate text-sm font-medium hover:text-primary hover:underline">
                    {formatReferrer(ref.referrer)}
                  </FilterLink>
                  {ref.referrer === '' ? null : <a href={ref.referrer} target="_blank" rel="noopener noreferrer" className="text-muted-foreground transition-colors hover:text-primary" onClick={(event) => { event.stopPropagation(); }}><ExternalLink className="h-3 w-3" /></a>}
                </div>
                <div className="ml-2 flex items-center gap-3">
                  <Badge variant="outline" className="font-mono">{ref.visitors.toLocaleString()}</Badge>
                  <span className="w-12 text-right text-xs text-muted-foreground">{percentage}%</span>
                </div>
              </div>
              <Progress value={barWidth} className="h-2" />
            </div>
          );
        }) : <ListEmptyState icon={Globe} title="No referrer data yet" description="Referrers will appear as visitors arrive from external sources" />}
      </div>
    </BoardCard>
  );
};

export default ReferrersCard;
