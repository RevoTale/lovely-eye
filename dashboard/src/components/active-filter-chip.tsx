import type { FunctionComponent, ReactNode } from 'react';
import { Badge } from '@/components/ui';
import { removeFilterValue, updateFilterSearch } from '@/lib/filter-utils';
import { Link } from '@/router';

export type ActiveFilterField =
  | 'browser'
  | 'country'
  | 'device'
  | 'eventName'
  | 'eventPath'
  | 'os'
  | 'page'
  | 'referrer';

interface ActiveFilterChipProps {
  field: ActiveFilterField;
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
    to="/sites/$siteId"
    params={{ siteId }}
    search={(prev) => ({
      ...updateFilterSearch(prev, field, removeFilterValue(prev[field], value)),
    })}
  >
    <Badge variant="secondary" className="flex cursor-pointer items-center gap-1 hover:bg-secondary/80">
      <span className="text-xs">
        {label}: {children}
      </span>
      <span className="ml-1 text-xs">x</span>
    </Badge>
  </Link>
);

export default ActiveFilterChip;
