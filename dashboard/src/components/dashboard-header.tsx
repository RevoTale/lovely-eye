import React from 'react';
import { Badge } from '@/components/ui';
import { Activity, Settings } from 'lucide-react';
import type { RealtimeStats, Site } from '@/generated/graphql';
import { Link } from '@/router';

interface DashboardHeaderProps {
  site: Site;
  siteId: string;
  realtime: RealtimeStats | undefined;
}

export function DashboardHeader({ site, siteId, realtime }: DashboardHeaderProps): React.JSX.Element {
  return (
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
  );
}
