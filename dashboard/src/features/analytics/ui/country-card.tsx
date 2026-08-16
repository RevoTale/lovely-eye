import { Globe } from 'lucide-react';
import type { FunctionComponent } from 'react';
import type { DashboardLoadState } from '@/features/analytics/model/dashboard-load-state';
import BoardCard from '@/features/analytics/ui/board-card';
import { FilterLink } from '@/features/analytics/ui/filter-link';
import { ListEmptyState } from '@/features/analytics/ui/list-empty-state';
import {
  CountryFieldsFragmentDoc,
  type CountryStatsFieldsFragment,
} from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';
import { SKELETON_KEYS } from '@/shared/lib/skeleton';
import { Badge } from '@/shared/ui/badge';
import { Progress } from '@/shared/ui/progress';
import { Skeleton } from '@/shared/ui/skeleton';

interface CountryCardProps {
  countries: CountryStatsFieldsFragment[];
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

const CountryCard: FunctionComponent<CountryCardProps> = ({
  countries,
  total,
  totalVisitors,
  page,
  pageSize,
  siteId,
  onPageChange,
  state = 'ready',
}) => (
  <BoardCard
    title='Countries'
    icon={Globe}
    state={state}
    pagination={{ page, pageSize, total, onPageChange, align: 'center' }}
    overlayLabel='Refreshing countries'
    skeleton={
      <div className='space-y-3'>
        {SKELETON_KEYS.map((key) => (
          <div key={key} className='space-y-2'>
            <div className='flex items-center justify-between'>
              <Skeleton className='h-4 w-28' />
              <Skeleton className='h-5 w-16' />
            </div>
            <Skeleton className='h-2 w-full' />
          </div>
        ))}
      </div>
    }
  >
    <div className='space-y-3'>
      {countries.length > EMPTY_COUNT ? (
        countries.map((countryStat) => {
          const percentage =
            totalVisitors > EMPTY_COUNT
              ? (countryStat.visitors / totalVisitors) * PERCENT_MULTIPLIER
              : EMPTY_COUNT;
          const country = readFragment(CountryFieldsFragmentDoc, countryStat.country);
          return (
            <div key={country.code}>
              <div className='mb-1 flex min-w-0 items-center justify-between gap-2'>
                <FilterLink
                  siteId={siteId}
                  filterKey='country'
                  value={country.code}
                  className='min-w-0 flex-1 cursor-pointer truncate text-sm font-medium hover:text-primary hover:underline'
                >
                  {country.name}
                </FilterLink>
                <div className='flex shrink-0 items-center gap-2'>
                  <Badge variant='secondary'>{countryStat.visitors.toLocaleString()}</Badge>
                  <span className='text-xs text-muted-foreground'>
                    {percentage.toFixed(PERCENT_PRECISION)}%
                  </span>
                </div>
              </div>
              <Progress value={percentage} className='h-2' />
            </div>
          );
        })
      ) : (
        <ListEmptyState title='No country data yet' />
      )}
    </div>
  </BoardCard>
);

export default CountryCard;
