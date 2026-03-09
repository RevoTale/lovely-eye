import { ArrowDownRight, ArrowUpRight } from 'lucide-react';
import type { ElementType, FunctionComponent } from 'react';
import DashboardCardState from '@/components/dashboard-card-state';
import { Card, CardContent, CardHeader, CardTitle, Skeleton } from '@/components/ui';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface StatCardProps {
  title: string;
  value: string | number;
  icon: ElementType;
  state?: DashboardLoadState;
  suffix?: string;
  trend?: 'up' | 'down';
  trendValue?: string;
}

const StatCard: FunctionComponent<StatCardProps> = ({ title, value, icon: Icon, state = 'ready', suffix, trend, trendValue }) => {
  const hasSuffix = suffix !== undefined && suffix !== '';
  const hasTrend = trend !== undefined && trendValue !== undefined && trendValue !== '';

  return (
    <Card className="transition-shadow hover:shadow-md">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
          <Icon className="h-4 w-4 text-primary" />
        </div>
      </CardHeader>
      <CardContent>
        <DashboardCardState state={state} skeleton={<Skeleton className="h-8 w-32" />} className="min-h-12" overlayLabel="Refreshing">
          <>
            <div className="text-2xl font-bold tracking-tight">
              {value}
              {hasSuffix ? <span className="ml-1 text-sm font-normal text-muted-foreground">{suffix}</span> : null}
            </div>
            {hasTrend ? (
              <div className="mt-1 flex items-center gap-1">
                {trend === 'up' ? <ArrowUpRight className="h-3 w-3 text-green-500" /> : <ArrowDownRight className="h-3 w-3 text-red-500" />}
                <span className={`text-xs ${trend === 'up' ? 'text-green-500' : 'text-red-500'}`}>{trendValue}</span>
                <span className="text-xs text-muted-foreground">vs last period</span>
              </div>
            ) : null}
          </>
        </DashboardCardState>
      </CardContent>
    </Card>
  );
};

export default StatCard;
