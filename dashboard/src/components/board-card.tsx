import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui';
import { PaginationControls } from '@/components/pagination-controls';

interface PaginationProps {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  align?: 'start' | 'center';
}

interface BoardCardProps {
  title: string;
  icon: React.ElementType;
  headerRight?: React.ReactNode;
  children: React.ReactNode;
  pagination?: PaginationProps;
}

export function BoardCard({
  title,
  icon: Icon,
  headerRight,
  children,
  pagination,
}: BoardCardProps): React.JSX.Element {
  const paginationAlign = pagination?.align ?? 'start';
  const hasHeaderRight = headerRight !== null && headerRight !== undefined;
  const hasPagination = pagination !== undefined;

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className={`flex items-center ${hasHeaderRight ? 'justify-between' : 'gap-2'}`}>
          <div className="flex items-center gap-2">
            <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
              <Icon className="h-4 w-4 text-primary" />
            </div>
            {title}
          </div>
          {headerRight}
        </CardTitle>
      </CardHeader>
      <CardContent>
        {children}
        {hasPagination && pagination.total > pagination.pageSize ? (
          <div className={`mt-4 flex ${paginationAlign === 'center' ? 'justify-center' : 'justify-start'}`}>
            <PaginationControls
              page={pagination.page}
              pageSize={pagination.pageSize}
              total={pagination.total}
              onPageChange={pagination.onPageChange}
            />
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}
