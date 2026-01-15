import React from 'react';
import { Badge, Card, CardContent, CardHeader, CardTitle, Progress } from '@/components/ui';
import type { DeviceStats } from '@/generated/graphql';
import { Link } from '@/router';
import { addFilterValue } from '@/lib/filter-utils';
import { Monitor, Smartphone } from 'lucide-react';
import { PaginationControls } from '@/components/pagination-controls';

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
          {devices.map((deviceStat, index) => {
            const percentage = totalVisitors > 0 ? (deviceStat.visitors / totalVisitors) * 100 : 0;

            return (
              <div key={index} className="space-y-2">
                <div className="flex items-center justify-between">
                  <Link
                    to="/sites/$siteId"
                    params={{ siteId }}
                    search={(prev) => ({
                      ...prev,
                      device: addFilterValue(prev.device, deviceStat.device),
                    })}
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
        {total > pageSize ? (
          <div className="mt-4">
            <PaginationControls
              page={page}
              pageSize={pageSize}
              total={total}
              onPageChange={onPageChange}
            />
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}
