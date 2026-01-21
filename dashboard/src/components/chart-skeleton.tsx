import React from 'react';
import { Card, CardContent, CardHeader, Skeleton } from '@/components/ui';

export function ChartSkeleton(): React.JSX.Element {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-8 w-[140px]" />
      </CardHeader>
      <CardContent>
        <Skeleton className="h-[300px] w-full" />
      </CardContent>
    </Card>
  );
}
