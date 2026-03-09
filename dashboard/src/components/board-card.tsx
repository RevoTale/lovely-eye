import type { ElementType, FunctionComponent, ReactNode } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui';
import DashboardCardState from '@/components/dashboard-card-state';
import { PaginationControls } from '@/components/pagination-controls';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface PaginationProps {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  align?: 'start' | 'center';
}

interface BoardCardProps {
  title: string;
  icon: ElementType;
  children: ReactNode;
  skeleton: ReactNode;
  state?: DashboardLoadState;
  headerRight?: ReactNode;
  pagination?: PaginationProps;
  overlayLabel?: string;
  contentClassName?: string;
}

const BoardCard: FunctionComponent<BoardCardProps> = ({
  title,
  icon: Icon,
  children,
  skeleton,
  state = 'ready',
  headerRight,
  pagination,
  overlayLabel,
  contentClassName,
}) => {
  const showPagination = pagination !== undefined && pagination.total > pagination.pageSize;
  const paginationAlign = pagination?.align ?? 'start';

  return (
    <Card className="transition-shadow hover:shadow-md">
      <CardHeader>
        <CardTitle className={`flex items-center ${headerRight !== undefined ? 'justify-between' : 'gap-2'}`}>
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
              <Icon className="h-4 w-4 text-primary" />
            </div>
            {title}
          </div>
          {headerRight}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <DashboardCardState state={state} skeleton={skeleton} className={contentClassName} overlayLabel={overlayLabel}>
          <>
            {children}
            {pagination !== undefined ? (
              <div className={`mt-4 flex min-h-9 ${paginationAlign === 'center' ? 'justify-center' : 'justify-start'}`}>
                {showPagination ? (
                  <PaginationControls
                    page={pagination.page}
                    pageSize={pagination.pageSize}
                    total={pagination.total}
                    onPageChange={pagination.onPageChange}
                  />
                ) : null}
              </div>
            ) : null}
          </>
        </DashboardCardState>
      </CardContent>
    </Card>
  );
};

export default BoardCard;
