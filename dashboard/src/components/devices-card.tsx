import { Gamepad2, Monitor, Smartphone, Tablet, Tv, Watch } from 'lucide-react';
import type { ElementType, FunctionComponent } from 'react';
import BoardCard from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';
import { Badge, Progress, Skeleton } from '@/components/ui';
import type { DeviceStatsFieldsFragment } from '@/gql/graphql';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface DevicesCardProps {
  devices: DeviceStatsFieldsFragment[];
  total: number;
  totalVisitors: number;
  page: number;
  pageSize: number;
  siteId: string;
  onPageChange: (page: number) => void;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;
const PERCENT_MULTIPLIER = 100;
const PERCENT_PRECISION = 1;

const getDeviceIcon = (device: string): ElementType => {
  switch (device) {
    case 'desktop': return Monitor;
    case 'tablet': return Tablet;
    case 'smart-tv': return Tv;
    case 'console': return Gamepad2;
    case 'watch': return Watch;
    default: return Smartphone;
  }
};

const DevicesCard: FunctionComponent<DevicesCardProps> = ({ devices, total, totalVisitors, page, pageSize, siteId, onPageChange, state = 'ready' }) => (
  <BoardCard
    title="Device Types"
    icon={Monitor}
    state={state}
    pagination={{ page, pageSize, total, onPageChange }}
    overlayLabel="Refreshing devices"
    skeleton={<div className="grid gap-4 sm:grid-cols-2">{Array.from({ length: 4 }, (_, index) => <div key={index} className="space-y-2"><div className="flex items-center justify-between"><Skeleton className="h-5 w-24" /><Skeleton className="h-5 w-16" /></div><Skeleton className="h-2 w-full" /></div>)}</div>}
  >
    {devices.length > EMPTY_COUNT ? (
      <div className="grid gap-4 sm:grid-cols-2">
        {devices.map((deviceStat) => {
          const percentage = totalVisitors > EMPTY_COUNT ? (deviceStat.visitors / totalVisitors) * PERCENT_MULTIPLIER : EMPTY_COUNT;
          const DeviceIcon = getDeviceIcon(deviceStat.device);
          return (
            <div key={deviceStat.device} className="space-y-2">
              <div className="flex items-center justify-between">
                <FilterLink siteId={siteId} filterKey="device" value={deviceStat.device} className="flex cursor-pointer items-center gap-2 hover:text-primary">
                  <DeviceIcon className="h-5 w-5 text-primary" />
                  <span className="text-sm font-medium capitalize hover:underline">{deviceStat.device}</span>
                </FilterLink>
                <div className="flex items-center gap-2">
                  <Badge variant="secondary">{deviceStat.visitors.toLocaleString()}</Badge>
                  <span className="text-sm text-muted-foreground">{percentage.toFixed(PERCENT_PRECISION)}%</span>
                </div>
              </div>
              <Progress value={percentage} className="h-2" />
            </div>
          );
        })}
      </div>
    ) : <ListEmptyState title="No device data yet" />}
  </BoardCard>
);

export default DevicesCard;
