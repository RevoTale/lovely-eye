import React from 'react';
import { Badge, Card, CardContent, CardHeader, CardTitle, Progress } from '@/components/ui';
import type { PageStats } from '@/generated/graphql';
import { Link } from '@/router';
import { addFilterValue } from '@/lib/filter-utils';
import { PaginationControls } from '@/components/pagination-controls';
import { Globe } from 'lucide-react';

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
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
            <Globe className="h-4 w-4 text-primary" />
          </div>
          Top Pages
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {pages.length > 0 ? (
            pages.map((pageStat, index) => (
              <div key={index}>
                <div className="flex items-center justify-between mb-1">
                  <Link
                    to="/sites/$siteId"
                    params={{ siteId }}
                    search={(prev) => ({
                      ...prev,
                      page: addFilterValue(prev.page, pageStat.path),
                    })}
                    className="text-sm font-medium truncate max-w-[200px] hover:text-primary hover:underline cursor-pointer"
                  >
                    {pageStat.path}
                  </Link>
                  <Badge variant="secondary" className="ml-2">
                    {pageStat.views.toLocaleString()}
                  </Badge>
                </div>
                <Progress value={maxViews ? (pageStat.views / maxViews) * 100 : 0} className="h-2" />
              </div>
            ))
          ) : (
            <p className="text-sm text-muted-foreground text-center py-4">No page data yet</p>
          )}
        </div>
        {total > pageSize ? (
          <div className="mt-4">
            <PaginationControls
              page={page}
              pageSize={pageSize}
              total={total}
              onPageChange={onPageChange}
            />
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}
