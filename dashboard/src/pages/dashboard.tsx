import React from 'react';
import { useParams } from '@tanstack/react-router';
import { useQuery } from '@apollo/client';
import { DASHBOARD_QUERY, REALTIME_QUERY, SITE_QUERY } from '@/graphql';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui';
import { Users, Eye, Clock, TrendingDown, Globe, Monitor, Smartphone } from 'lucide-react';
import type { DashboardStats, Site, RealtimeStats } from '@/generated/graphql';
import { siteDetailRoute } from '@/router';

function StatCard({ 
  title, 
  value, 
  icon: Icon, 
  suffix 
}: { 
  title: string; 
  value: string | number; 
  icon: React.ElementType; 
  suffix?: string;
}): React.JSX.Element {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">
          {value}
          {suffix ? <span className="text-sm font-normal text-muted-foreground ml-1">{suffix}</span> : null}
        </div>
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
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-muted-foreground">Loading dashboard...</div>
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
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{site.name}</h1>
        <p className="text-muted-foreground">{site.domain}</p>
      </div>

      {/* Realtime visitors */}
      {realtime ? (
        <Card className="bg-primary/5 border-primary/20">
          <CardContent className="py-4">
            <div className="flex items-center gap-2">
              <span className="relative flex h-3 w-3">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
              </span>
              <span className="text-lg font-semibold">{realtime.visitors}</span>
              <span className="text-muted-foreground">visitors online now</span>
            </div>
          </CardContent>
        </Card>
      ) : null}

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
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Globe className="h-5 w-5" />
                  Top Pages
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {stats.topPages.length > 0 ? (
                    stats.topPages.map((page, index) => (
                      <div key={index} className="flex items-center justify-between">
                        <span className="text-sm truncate max-w-[200px]">{page.path}</span>
                        <span className="text-sm text-muted-foreground">{page.views.toLocaleString()}</span>
                      </div>
                    ))
                  ) : (
                    <p className="text-sm text-muted-foreground">No page data yet</p>
                  )}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Globe className="h-5 w-5" />
                  Top Referrers
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {stats.topReferrers.length > 0 ? (
                    stats.topReferrers.map((ref, index) => (
                      <div key={index} className="flex items-center justify-between">
                        <span className="text-sm truncate max-w-[200px]">{ref.referrer || 'Direct'}</span>
                        <span className="text-sm text-muted-foreground">{ref.visitors.toLocaleString()}</span>
                      </div>
                    ))
                  ) : (
                    <p className="text-sm text-muted-foreground">No referrer data yet</p>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Device breakdown */}
          <Card>
            <CardHeader>
              <CardTitle>Device Types</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex gap-8">
                {stats.devices.map((deviceStat, index) => (
                  <div key={index} className="flex items-center gap-2">
                    {deviceStat.device === 'desktop' ? (
                      <Monitor className="h-5 w-5 text-muted-foreground" />
                    ) : (
                      <Smartphone className="h-5 w-5 text-muted-foreground" />
                    )}
                    <span className="text-sm capitalize">{deviceStat.device}</span>
                    <span className="text-sm text-muted-foreground">
                      {deviceStat.visitors}
                    </span>
                  </div>
                ))}
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
