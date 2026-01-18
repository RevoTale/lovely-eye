import React from 'react';
import { Globe } from 'lucide-react';
import { Badge, Progress } from '@/components/ui';
import type { CountryStats } from '@/gql/graphql';
import { BoardCard } from '@/components/board-card';
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
}

export function CountryCard({ countries, total, totalVisitors, page, pageSize, siteId, onPageChange }: CountryCardProps): React.JSX.Element {
  return (
    <BoardCard
      title="Countries"
      icon={Globe}
      pagination={{ page, pageSize, total, onPageChange, align: 'center' }}
    >
      <div className="space-y-3">
        {countries.length > 0 ? (
          countries.map((countryStat, index) => {
            const percentage = totalVisitors > 0 ? (countryStat.visitors / totalVisitors) * 100 : 0;

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
                    <span className="text-xs text-muted-foreground">{percentage.toFixed(1)}%</span>
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
