import React from 'react';
import { useParams } from '@tanstack/react-router';
import { useQuery } from '@apollo/client';
import { DASHBOARD_QUERY, REALTIME_QUERY, SITE_QUERY } from '@/graphql';
import { Card, CardContent, CardHeader, CardTitle, Badge, Skeleton, Progress } from '@/components/ui';
import { Users, Eye, Clock, TrendingDown, Globe, Monitor, Smartphone, Activity, ArrowUpRight, ArrowDownRight } from 'lucide-react';
import type { DashboardStats, Site, RealtimeStats } from '@/generated/graphql';
import { siteDetailRoute } from '@/router';

function StatCard({
  title,
  value,
  icon: Icon,
  suffix,
  trend,
  trendValue
}: {
  title: string;
  value: string | number;
  icon: React.ElementType;
  suffix?: string;
  trend?: 'up' | 'down';
  trendValue?: string;
}): React.JSX.Element {
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
          {suffix ? <span className="text-sm font-normal text-muted-foreground ml-1">{suffix}</span> : null}
        </div>
        {trend && trendValue ? (
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

function formatDuration(seconds: number): string {
  if (seconds < 60) {
    return `${String(Math.round(seconds))}s`;
  }
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.round(seconds % 60);
  return `${String(minutes)}m ${String(remainingSeconds)}s`;
}

export function DashboardPage(): React.JSX.Element {
  const { siteId } = useParams({ from: siteDetailRoute.id });

  const { data: siteData, loading: siteLoading } = useQuery(SITE_QUERY, {
    variables: { id: siteId },
    skip: !siteId,
  });

  const { data: dashboardData, loading: dashboardLoading } = useQuery(DASHBOARD_QUERY, {
    variables: { siteId },
    skip: !siteId,
    pollInterval: 60000, // Refresh every minute
  });

  const { data: realtimeData } = useQuery(REALTIME_QUERY, {
    variables: { siteId },
    skip: !siteId,
    pollInterval: 5000, // Refresh every 5 seconds
  });

  if (siteLoading || dashboardLoading) {
    return (
      <div className="space-y-6">
        <div className="space-y-2">
          <Skeleton className="h-8 w-64" />
          <Skeleton className="h-4 w-48" />
        </div>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-4 w-24" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-8 w-32" />
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  const site = siteData?.site as Site | undefined;
  const stats = dashboardData?.dashboard as DashboardStats | undefined;
  const realtime = realtimeData?.realtime as RealtimeStats | undefined;

  if (!site) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-destructive">Site not found</div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">{site.name}</h1>
          <p className="text-muted-foreground mt-1">{site.domain}</p>
        </div>
        {realtime ? (
          <Badge variant="outline" className="flex items-center gap-2 px-3 py-2">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            <Activity className="h-4 w-4" />
            <span className="font-semibold">{realtime.visitors}</span>
            <span className="text-muted-foreground">online</span>
          </Badge>
        ) : null}
      </div>

      {/* Stats grid */}
      {stats ? (
        <>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <StatCard
              title="Total Visitors"
              value={stats.visitors.toLocaleString()}
              icon={Users}
            />
            <StatCard
              title="Page Views"
              value={stats.pageViews.toLocaleString()}
              icon={Eye}
            />
            <StatCard
              title="Avg. Session"
              value={formatDuration(stats.avgDuration)}
              icon={Clock}
            />
            <StatCard
              title="Bounce Rate"
              value={`${String(Math.round(stats.bounceRate))}%`}
              icon={TrendingDown}
            />
          </div>

          {/* Top pages and referrers */}
          <div className="grid gap-6 md:grid-cols-2">
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
                  {stats.topPages.length > 0 ? (
                    stats.topPages.map((page, index) => (
                      <div key={index}>
                        <div className="flex items-center justify-between mb-1">
                          <span className="text-sm font-medium truncate max-w-[200px]">{page.path}</span>
                          <Badge variant="secondary" className="ml-2">
                            {page.views.toLocaleString()}
                          </Badge>
                        </div>
                        <Progress value={stats.topPages[0] ? (page.views / stats.topPages[0].views) * 100 : 0} className="h-2" />
                      </div>
                    ))
                  ) : (
                    <p className="text-sm text-muted-foreground text-center py-4">No page data yet</p>
                  )}
                </div>
              </CardContent>
            </Card>

            <Card className="hover:shadow-md transition-shadow">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
                    <Globe className="h-4 w-4 text-primary" />
                  </div>
                  Top Referrers
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {stats.topReferrers.length > 0 ? (
                    stats.topReferrers.map((ref, index) => (
                      <div key={index}>
                        <div className="flex items-center justify-between mb-1">
                          <span className="text-sm font-medium truncate max-w-[200px]">{ref.referrer || 'Direct'}</span>
                          <Badge variant="secondary" className="ml-2">
                            {ref.visitors.toLocaleString()}
                          </Badge>
                        </div>
                        <Progress value={stats.topReferrers[0] ? (ref.visitors / stats.topReferrers[0].visitors) * 100 : 0} className="h-2" />
                      </div>
                    ))
                  ) : (
                    <p className="text-sm text-muted-foreground text-center py-4">No referrer data yet</p>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Device breakdown */}
          <Card className="hover:shadow-md transition-shadow">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
                  <Monitor className="h-4 w-4 text-primary" />
                </div>
                Device Types
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 sm:grid-cols-2">
                {stats.devices.map((deviceStat, index) => {
                  const totalVisitors = stats.devices.reduce((sum, d) => sum + d.visitors, 0);
                  const percentage = totalVisitors > 0 ? (deviceStat.visitors / totalVisitors) * 100 : 0;

                  return (
                    <div key={index} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          {deviceStat.device === 'desktop' ? (
                            <Monitor className="h-5 w-5 text-primary" />
                          ) : (
                            <Smartphone className="h-5 w-5 text-primary" />
                          )}
                          <span className="text-sm font-medium capitalize">{deviceStat.device}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <Badge variant="secondary">
                            {deviceStat.visitors.toLocaleString()}
                          </Badge>
                          <span className="text-sm text-muted-foreground">
                            {percentage.toFixed(1)}%
                          </span>
                        </div>
                      </div>
                      <Progress value={percentage} className="h-2" />
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
        </>
      ) : (
        <Card>
          <CardContent className="py-12">
            <div className="text-center text-muted-foreground">
              <Eye className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>No analytics data yet. Add the tracking script to start collecting data.</p>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
