import { Globe } from 'lucide-react';
import type { FunctionComponent } from 'react';
import BoardCard from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';
import { Badge, Progress, Skeleton } from '@/components/ui';
import { CountryFieldsFragmentDoc, type CountryStatsFieldsFragment } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

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
const SKELETON_ROWS = 5;

const CountryCard: FunctionComponent<CountryCardProps> = ({ countries, total, totalVisitors, page, pageSize, siteId, onPageChange, state = 'ready' }) => (
  <BoardCard
    title="Countries"
    icon={Globe}
    state={state}
    pagination={{ page, pageSize, total, onPageChange, align: 'center' }}
    overlayLabel="Refreshing countries"
    skeleton={<div className="space-y-3">{Array.from({ length: SKELETON_ROWS }, (_, index) => <div key={index} className="space-y-2"><div className="flex items-center justify-between"><Skeleton className="h-4 w-28" /><Skeleton className="h-5 w-16" /></div><Skeleton className="h-2 w-full" /></div>)}</div>}
  >
    <div className="space-y-3">
      {countries.length > EMPTY_COUNT ? countries.map((countryStat) => {
        const percentage = totalVisitors > EMPTY_COUNT ? (countryStat.visitors / totalVisitors) * PERCENT_MULTIPLIER : EMPTY_COUNT;
        const country = getFragmentData(CountryFieldsFragmentDoc, countryStat.country);
        return (
          <div key={country.code}>
            <div className="mb-1 flex items-center justify-between">
              <FilterLink siteId={siteId} filterKey="country" value={country.code} className="max-w-[200px] cursor-pointer truncate text-sm font-medium hover:text-primary hover:underline">
                {country.name}
              </FilterLink>
              <div className="flex items-center gap-2">
                <Badge variant="secondary">{countryStat.visitors.toLocaleString()}</Badge>
                <span className="text-xs text-muted-foreground">{percentage.toFixed(PERCENT_PRECISION)}%</span>
              </div>
            </div>
            <Progress value={percentage} className="h-2" />
          </div>
        );
      }) : <ListEmptyState title="No country data yet" />}
    </div>
  </BoardCard>
);

export default CountryCard;
