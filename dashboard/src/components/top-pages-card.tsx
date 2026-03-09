import type { FunctionComponent } from 'react';
import { Globe } from 'lucide-react';
import { Badge, Progress, Skeleton } from '@/components/ui';
import BoardCard from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';
import type { PageStatsFieldsFragment } from '@/gql/graphql';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface TopPagesCardProps {
  pages: PageStatsFieldsFragment[];
  total: number;
  page: number;
  pageSize: number;
  siteId: string;
  onPageChange: (page: number) => void;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;
const FIRST_INDEX = 0;
const PERCENT_MULTIPLIER = 100;
const SKELETON_ROWS = 5;

const TopPagesCard: FunctionComponent<TopPagesCardProps> = ({ pages, total, page, pageSize, siteId, onPageChange, state = 'ready' }) => {
  const maxViews = pages.length > EMPTY_COUNT ? pages[FIRST_INDEX]?.views ?? EMPTY_COUNT : EMPTY_COUNT;

  return (
    <BoardCard
      title="Top Pages"
      icon={Globe}
      state={state}
      pagination={{ page, pageSize, total, onPageChange }}
      overlayLabel="Refreshing pages"
      skeleton={<div className="space-y-3">{Array.from({ length: SKELETON_ROWS }, (_, index) => <div key={index} className="space-y-2"><div className="flex items-center justify-between"><Skeleton className="h-4 w-32" /><Skeleton className="h-5 w-12" /></div><Skeleton className="h-2 w-full" /></div>)}</div>}
    >
      <div className="space-y-3">
        {pages.length > EMPTY_COUNT ? pages.map((pageStat) => (
          <div key={pageStat.path}>
            <div className="mb-1 flex items-center justify-between">
              <FilterLink siteId={siteId} filterKey="page" value={pageStat.path} className="max-w-[200px] cursor-pointer truncate text-sm font-medium hover:text-primary hover:underline">
                {pageStat.path}
              </FilterLink>
              <Badge variant="secondary" className="ml-2">{pageStat.views.toLocaleString()}</Badge>
            </div>
            <Progress value={maxViews > EMPTY_COUNT ? (pageStat.views / maxViews) * PERCENT_MULTIPLIER : EMPTY_COUNT} className="h-2" />
          </div>
        )) : <ListEmptyState title="No page data yet" />}
      </div>
    </BoardCard>
  );
};

export default TopPagesCard;
