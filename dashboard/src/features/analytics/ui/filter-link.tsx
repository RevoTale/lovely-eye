import { Link } from '@tanstack/react-router';
import {
  type AnalyticsFilterKey,
  addAnalyticsFilter,
} from '@/features/analytics/model/analytics-search';

interface FilterLinkProps {
  siteId: string;
  filterKey: AnalyticsFilterKey;
  value: string;
  className?: string;
  children: React.ReactNode;
}

export const FilterLink = ({
  siteId,
  filterKey,
  value,
  className,
  children,
}: FilterLinkProps): React.ReactNode => (
  <Link
    to='/sites/$siteId/analytics'
    params={{ siteId }}
    search={(prev) => addAnalyticsFilter(prev, filterKey, value)}
    className={className}
  >
    {children}
  </Link>
);
