
import { Badge, Progress } from '@/components/ui';
import type { DeviceStatsFieldsFragment } from '@/gql/graphql';
import { Monitor, Smartphone } from 'lucide-react';
import { BoardCard, BoardCardSkeleton } from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';

interface DevicesCardProps {
  devices: DeviceStatsFieldsFragment[];
  total: number;
  totalVisitors: number;
  page: number;
  pageSize: number;
  siteId: string;
  onPageChange: (page: number) => void;
  loading?: boolean;
}

const EMPTY_COUNT = 0;
const PERCENT_MULTIPLIER = 100;
const PERCENT_PRECISION = 1;

export const DevicesCard = ({
  devices,
  total,
  totalVisitors,
  page,
  pageSize,
  siteId,
  onPageChange,
  loading = false,
}: DevicesCardProps): React.ReactNode => {
  if (loading) {
    return <BoardCardSkeleton title="Device Types" icon={Monitor} />;
  }

  return (
    <BoardCard
      title="Device Types"
      icon={Monitor}
      pagination={{ page, pageSize, total, onPageChange }}
    >
      {devices.length > EMPTY_COUNT ? (
        <div className="grid gap-4 sm:grid-cols-2">
          {devices.map((deviceStat, index) => {
            const percentage =
              totalVisitors > EMPTY_COUNT
                ? (deviceStat.visitors / totalVisitors) * PERCENT_MULTIPLIER
                : EMPTY_COUNT;

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
                      {percentage.toFixed(PERCENT_PRECISION)}%
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
