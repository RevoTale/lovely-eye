
import { Badge, Progress } from '@/components/ui';
import type { OperatingSystemStatsFieldsFragment } from '@/gql/graphql';
import { Laptop } from 'lucide-react';
import { BoardCard, BoardCardSkeleton } from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';

interface OSCardProps {
  operatingSystems: OperatingSystemStatsFieldsFragment[];
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

export const OSCard = ({
  operatingSystems,
  total,
  totalVisitors,
  page,
  pageSize,
  siteId,
  onPageChange,
  loading = false,
}: OSCardProps): React.ReactNode => {
  if (loading) {
    return <BoardCardSkeleton title="Operating Systems" icon={Laptop} />;
  }

  return (
    <BoardCard
      title="Operating Systems"
      icon={Laptop}
      pagination={{ page, pageSize, total, onPageChange }}
    >
      {operatingSystems.length > EMPTY_COUNT ? (
        <div className="space-y-3">
          {operatingSystems.map((operatingSystem) => {
            const percentage =
              totalVisitors > EMPTY_COUNT
                ? (operatingSystem.visitors / totalVisitors) * PERCENT_MULTIPLIER
                : EMPTY_COUNT;

            return (
              <div key={operatingSystem.os}>
                <div className="mb-1 flex items-center justify-between gap-2">
                  <FilterLink
                    siteId={siteId}
                    filterKey="os"
                    value={operatingSystem.os}
                    className="truncate text-sm font-medium hover:text-primary hover:underline"
                  >
                    {operatingSystem.os}
                  </FilterLink>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary">
                      {operatingSystem.visitors.toLocaleString()}
                    </Badge>
                    <span className="text-xs text-muted-foreground">
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
        <ListEmptyState title="No operating system data yet" />
      )}
    </BoardCard>
  );
}
