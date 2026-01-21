import React from 'react';
import { Globe } from 'lucide-react';
import { Badge, Progress } from '@/components/ui';
import type { CountryStats } from '@/gql/graphql';
import { BoardCard, BoardCardSkeleton } from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';

interface CountryCardProps {
  countries: CountryStats[];
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

export function CountryCard({ countries, total, totalVisitors, page, pageSize, siteId, onPageChange, loading = false }: CountryCardProps): React.JSX.Element {
  if (loading) {
    return <BoardCardSkeleton title="Countries" icon={Globe} />;
  }

  return (
    <BoardCard
      title="Countries"
      icon={Globe}
      pagination={{ page, pageSize, total, onPageChange, align: 'center' }}
    >
      <div className="space-y-3">
        {countries.length > EMPTY_COUNT ? (
          countries.map((countryStat, index) => {
            const percentage =
              totalVisitors > EMPTY_COUNT
                ? (countryStat.visitors / totalVisitors) * PERCENT_MULTIPLIER
                : EMPTY_COUNT;

            return (
              <div key={index}>
                <div className="flex items-center justify-between mb-1">
                  <FilterLink
                    siteId={siteId}
                    filterKey="country"
                    value={countryStat.country}
                    className="text-sm font-medium truncate max-w-[200px] hover:text-primary hover:underline cursor-pointer"
                  >
                    {countryStat.country}
                  </FilterLink>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary">{countryStat.visitors.toLocaleString()}</Badge>
                    <span className="text-xs text-muted-foreground">
                      {percentage.toFixed(PERCENT_PRECISION)}%
                    </span>
                  </div>
                </div>
                <Progress value={percentage} className="h-2" />
              </div>
            );
          })
        ) : (
          <ListEmptyState title="No country data yet" />
        )}
      </div>
    </BoardCard>
  );
}
