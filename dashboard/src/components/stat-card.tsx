import React from 'react';
import { Card, CardContent, CardHeader, CardTitle, Skeleton } from '@/components/ui';
import { ArrowDownRight, ArrowUpRight } from 'lucide-react';

interface StatCardProps {
  title: string;
  value: string | number;
  icon: React.ElementType;
  suffix?: string;
  trend?: 'up' | 'down';
  trendValue?: string;
}

export function StatCard({
  title,
  value,
  icon: Icon,
  suffix,
  trend,
  trendValue,
}: StatCardProps): React.JSX.Element {
  const hasSuffix = suffix !== undefined && suffix !== '';
  const hasTrend = trend !== undefined && trendValue !== undefined && trendValue !== '';

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
        <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
          <Icon className="h-4 w-4 text-primary" />
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold tracking-tight">
          {value}
          {hasSuffix ? (
            <span className="text-sm font-normal text-muted-foreground ml-1">{suffix}</span>
          ) : null}
        </div>
        {hasTrend ? (
          <div className="flex items-center gap-1 mt-1">
            {trend === 'up' ? (
              <ArrowUpRight className="h-3 w-3 text-green-500" />
            ) : (
              <ArrowDownRight className="h-3 w-3 text-red-500" />
            )}
            <span className={`text-xs ${trend === 'up' ? 'text-green-500' : 'text-red-500'}`}>
              {trendValue}
            </span>
            <span className="text-xs text-muted-foreground">vs last period</span>
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}

export function StatCardSkeleton({ title, icon: Icon }: { title: string; icon: React.ElementType }): React.JSX.Element {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
        <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
          <Icon className="h-4 w-4 text-primary" />
        </div>
      </CardHeader>
      <CardContent>
        <Skeleton className="h-8 w-32" />
      </CardContent>
    </Card>
  );
}
