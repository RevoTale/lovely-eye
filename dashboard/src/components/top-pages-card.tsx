
import { Badge, Progress } from '@/components/ui';
import type { PageStatsFieldsFragment } from '@/gql/graphql';
import { Globe } from 'lucide-react';
import { BoardCard, BoardCardSkeleton } from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';

interface TopPagesCardProps {
  pages: PageStatsFieldsFragment[];
  total: number;
  page: number;
  pageSize: number;
  siteId: string;
  onPageChange: (page: number) => void;
  loading?: boolean;
}

const EMPTY_COUNT = 0;
const FIRST_INDEX = 0;
const PERCENT_MULTIPLIER = 100;

export const TopPagesCard = ({
  pages,
  total,
  page,
  pageSize,
  siteId,
  onPageChange,
  loading = false,
}: TopPagesCardProps): React.ReactNode => {
  const maxViews = pages.length > EMPTY_COUNT ? pages[FIRST_INDEX]?.views ?? EMPTY_COUNT : EMPTY_COUNT;

  if (loading) {
    return <BoardCardSkeleton title="Top Pages" icon={Globe} />;
  }

  return (
    <BoardCard
      title="Top Pages"
      icon={Globe}
      pagination={{ page, pageSize, total, onPageChange }}
    >
      <div className="space-y-3">
        {pages.length > EMPTY_COUNT ? (
          pages.map((pageStat, index) => (
            <div key={index}>
              <div className="flex items-center justify-between mb-1">
                <FilterLink
                  siteId={siteId}
                  filterKey="page"
                  value={pageStat.path}
                  className="text-sm font-medium truncate max-w-[200px] hover:text-primary hover:underline cursor-pointer"
                >
                  {pageStat.path}
                </FilterLink>
                <Badge variant="secondary" className="ml-2">
                  {pageStat.views.toLocaleString()}
                </Badge>
              </div>
              <Progress
                value={
                  maxViews > EMPTY_COUNT
                    ? (pageStat.views / maxViews) * PERCENT_MULTIPLIER
                    : EMPTY_COUNT
                }
                className="h-2"
              />
            </div>
          ))
        ) : (
          <ListEmptyState title="No page data yet" />
        )}
      </div>
    </BoardCard>
  );
}
