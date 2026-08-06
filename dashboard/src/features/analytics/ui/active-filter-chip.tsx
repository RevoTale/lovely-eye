import { Link } from '@tanstack/react-router';
import type { FunctionComponent, ReactNode } from 'react';
import {
  type AnalyticsFilterKey,
  removeAnalyticsFilter,
} from '@/features/analytics/model/analytics-search';
import { Badge } from '@/shared/ui/badge';

interface ActiveFilterChipProps {
  field: AnalyticsFilterKey;
  label: string;
  siteId: string;
  value: string;
  children: ReactNode;
}

const ActiveFilterChip: FunctionComponent<ActiveFilterChipProps> = ({
  field,
  label,
  siteId,
  value,
  children,
}) => (
  <Link
    key={`${field}-${value}`}
    to='/sites/$siteId/analytics'
    params={{ siteId }}
    search={(prev) => removeAnalyticsFilter(prev, field, value)}
  >
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
