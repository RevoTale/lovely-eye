import React from 'react';
import { Badge, Progress } from '@/components/ui';
import type { PageStats } from '@/gql/graphql';
import { Globe } from 'lucide-react';
import { BoardCard } from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';

interface TopPagesCardProps {
  pages: PageStats[];
  total: number;
  page: number;
  pageSize: number;
  siteId: string;
  onPageChange: (page: number) => void;
}

export function TopPagesCard({
  pages,
  total,
  page,
  pageSize,
  siteId,
  onPageChange,
}: TopPagesCardProps): React.JSX.Element {
  const maxViews = pages[0] ? pages[0].views : 0;

  return (
    <BoardCard
      title="Top Pages"
      icon={Globe}
      pagination={{ page, pageSize, total, onPageChange }}
    >
      <div className="space-y-3">
        {pages.length > 0 ? (
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
              <Progress value={maxViews ? (pageStat.views / maxViews) * 100 : 0} className="h-2" />
            </div>
          ))
        ) : (
          <ListEmptyState title="No page data yet" />
        )}
      </div>
    </BoardCard>
  );
}
