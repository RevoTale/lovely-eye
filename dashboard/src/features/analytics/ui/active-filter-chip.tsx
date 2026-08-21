import { Link } from '@tanstack/react-router';
import type { FunctionComponent, ReactNode } from 'react';
import type { AnalyticsSearch } from '@/features/analytics/model/analytics-search';
import { Badge } from '@/shared/ui/badge';

interface ActiveFilterChipProps {
  label: string;
  siteId: string;
  remove: (search: AnalyticsSearch) => AnalyticsSearch;
  children: ReactNode;
}

const ActiveFilterChip: FunctionComponent<ActiveFilterChipProps> = ({
  label,
  siteId,
  remove,
  children,
}) => (
  <Link to='/sites/$siteId/analytics' params={{ siteId }} search={remove}>
    <Badge
      variant='secondary'
      className='flex cursor-pointer items-center gap-1 hover:bg-secondary/80'
    >
      <span className='text-xs'>
        {label}: {children}
      </span>
      <span className='ml-1 text-xs'>x</span>
    </Badge>
  </Link>
);

export default ActiveFilterChip;
