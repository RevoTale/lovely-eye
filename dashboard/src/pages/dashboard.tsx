import React from 'react';
import { useParams, useSearch } from '@tanstack/react-router';
import { useQuery } from '@apollo/client';
import { DASHBOARD_QUERY, REALTIME_QUERY, SITE_QUERY } from '@/graphql';
import { Card, CardContent, CardHeader, CardTitle, Badge, Skeleton, Progress, ChartContainer, ChartTooltip, ChartTooltipContent, ChartLegend, ChartLegendContent, type ChartConfig } from '@/components/ui';
import { Users, Eye, Clock, TrendingDown, Globe, Monitor, Smartphone, Activity, ArrowUpRight, ArrowDownRight, Settings, TrendingUp } from 'lucide-react';
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts';
import type { DashboardStats, Site, RealtimeStats } from '@/generated/graphql';
import { siteDetailRoute, Link } from '@/router';
import { ReferrersCard } from '@/components/referrers-card';
import { ActivePagesCard } from '@/components/active-pages-card';

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
  const search = useSearch({ from: siteDetailRoute.id }) as { view?: string; referrer?: string; device?: string; page?: string };

  const { data: siteData, loading: siteLoading } = useQuery(SITE_QUERY, {
    variables: { id: siteId },
    skip: !siteId,
  });

  // Build filter object from URL parameters
  const filter = {
    ...(search.referrer && { referrer: search.referrer }),
    ...(search.device && { device: search.device }),
    ...(search.page && { page: search.page }),
  };

  const { data: dashboardData, loading: dashboardLoading } = useQuery(DASHBOARD_QUERY, {
    variables: {
      siteId,
      filter: Object.keys(filter).length > 0 ? filter : undefined,
    },
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
        <div className="flex items-center gap-3">
          <Link to="/sites/$siteId" params={{ siteId }} search={{ view: 'settings' }}>
            <Badge variant="outline" className="flex items-center gap-2 px-3 py-2 cursor-pointer hover:bg-accent">
              <Settings className="h-4 w-4" />
              <span>Settings</span>
            </Badge>
          </Link>
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
      </div>

      {/* Active Filters */}
      {(search.referrer || search.device || search.page) && (
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-sm text-muted-foreground">Filtered by:</span>
          {search.referrer && (
            <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
              <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
                <span className="text-xs">Referrer: {search.referrer}</span>
                <span className="ml-1 text-xs">×</span>
              </Badge>
            </Link>
          )}
          {search.device && (
            <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
              <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
                <span className="text-xs">Device: {search.device}</span>
                <span className="ml-1 text-xs">×</span>
              </Badge>
            </Link>
          )}
          {search.page && (
            <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
              <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
                <span className="text-xs">Page: {search.page}</span>
                <span className="ml-1 text-xs">×</span>
              </Badge>
            </Link>
          )}
          <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
            <Badge variant="outline" className="cursor-pointer hover:bg-accent text-xs">
              Clear all
            </Badge>
          </Link>
        </div>
      )}

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

          {/* Analytics Chart */}
          {stats.dailyStats && stats.dailyStats.length > 0 ? (
            <Card className="hover:shadow-md transition-shadow">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
                    <TrendingUp className="h-4 w-4 text-primary" />
                  </div>
                  Analytics Overview
                </CardTitle>
              </CardHeader>
              <CardContent>
                <ChartContainer
                  config={{
                    visitors: {
                      label: 'Visitors',
                      color: 'hsl(var(--primary))',
                    },
                    pageViews: {
                      label: 'Page Views',
                      color: 'hsl(var(--chart-2))',
                    },
                    sessions: {
                      label: 'Sessions',
                      color: 'hsl(var(--chart-3))',
                    },
                  } satisfies ChartConfig}
                  className="h-[300px] w-full"
                >
                  <AreaChart
                    data={stats.dailyStats.map(stat => ({
                      date: new Date(stat.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
                      visitors: stat.visitors,
                      pageViews: stat.pageViews,
                      sessions: stat.sessions,
                    }))}
                    margin={{ top: 10, right: 10, left: 0, bottom: 0 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                    <XAxis
                      dataKey="date"
                      tickLine={false}
                      axisLine={false}
                      tickMargin={8}
                      className="text-xs"
                    />
                    <YAxis
                      tickLine={false}
                      axisLine={false}
                      tickMargin={8}
                      className="text-xs"
                    />
                    <ChartTooltip content={<ChartTooltipContent />} />
                    <ChartLegend content={<ChartLegendContent />} />
                    <Area
                      type="monotone"
                      dataKey="visitors"
                      stackId="1"
                      stroke="var(--color-visitors)"
                      fill="var(--color-visitors)"
                      fillOpacity={0.6}
                    />
                    <Area
                      type="monotone"
                      dataKey="pageViews"
                      stackId="2"
                      stroke="var(--color-pageViews)"
                      fill="var(--color-pageViews)"
                      fillOpacity={0.6}
                    />
                    <Area
                      type="monotone"
                      dataKey="sessions"
                      stackId="3"
                      stroke="var(--color-sessions)"
                      fill="var(--color-sessions)"
                      fillOpacity={0.6}
                    />
                  </AreaChart>
                </ChartContainer>
              </CardContent>
            </Card>
          ) : null}

          {/* Active Pages (Realtime) */}
          {realtime && realtime.activePages ? (
            <ActivePagesCard activePages={realtime.activePages} />
          ) : null}

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
                          <Link
                            to="/sites/$siteId"
                            params={{ siteId }}
                            search={{ page: page.path }}
                            className="text-sm font-medium truncate max-w-[200px] hover:text-primary hover:underline cursor-pointer"
                          >
                            {page.path}
                          </Link>
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

            <ReferrersCard
              referrers={stats.topReferrers}
              totalVisitors={stats.visitors}
              siteId={siteId}
            />
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
                        <Link
                          to="/sites/$siteId"
                          params={{ siteId }}
                          search={{ device: deviceStat.device }}
                          className="flex items-center gap-2 hover:text-primary cursor-pointer"
                        >
                          {deviceStat.device === 'desktop' ? (
                            <Monitor className="h-5 w-5 text-primary" />
                          ) : (
                            <Smartphone className="h-5 w-5 text-primary" />
                          )}
                          <span className="text-sm font-medium capitalize hover:underline">{deviceStat.device}</span>
                        </Link>
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
