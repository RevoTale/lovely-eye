
import { Link } from '@/router';
import { addFilterValue } from '@/lib/filter-utils';

type FilterKey = 'referrer' | 'device' | 'page' | 'country' | 'eventName' | 'eventPath';

interface FilterLinkProps {
  siteId: string;
  filterKey: FilterKey;
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
    to="/sites/$siteId"
    params={{ siteId }}
    search={(prev) => ({
      ...prev,
      [filterKey]: addFilterValue((prev as Record<string, string | string[] | undefined>)[filterKey], value),
    })}
    className={className}
  >
    {children}
  </Link>
)
