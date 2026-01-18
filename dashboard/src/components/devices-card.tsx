import React from 'react';
import { Badge, Progress } from '@/components/ui';
import type { DeviceStats } from '@/gql/graphql';
import { Monitor, Smartphone } from 'lucide-react';
import { BoardCard } from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';

interface DevicesCardProps {
  devices: DeviceStats[];
  total: number;
  totalVisitors: number;
  page: number;
  pageSize: number;
  siteId: string;
  onPageChange: (page: number) => void;
}

export function DevicesCard({
  devices,
  total,
  totalVisitors,
  page,
  pageSize,
  siteId,
  onPageChange,
}: DevicesCardProps): React.JSX.Element {
  return (
    <BoardCard
      title="Device Types"
      icon={Monitor}
      pagination={{ page, pageSize, total, onPageChange }}
    >
      {devices.length > 0 ? (
        <div className="grid gap-4 sm:grid-cols-2">
          {devices.map((deviceStat, index) => {
            const percentage = totalVisitors > 0 ? (deviceStat.visitors / totalVisitors) * 100 : 0;

            return (
              <div key={index} className="space-y-2">
                <div className="flex items-center justify-between">
                  <FilterLink
                    siteId={siteId}
                    filterKey="device"
                    value={deviceStat.device}
                    className="flex items-center gap-2 hover:text-primary cursor-pointer"
                  >
                    {deviceStat.device === 'desktop' ? (
                      <Monitor className="h-5 w-5 text-primary" />
                    ) : (
                      <Smartphone className="h-5 w-5 text-primary" />
                    )}
                    <span className="text-sm font-medium capitalize hover:underline">{deviceStat.device}</span>
                  </FilterLink>
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
      ) : (
        <ListEmptyState title="No device data yet" />
      )}
    </BoardCard>
  );
}
