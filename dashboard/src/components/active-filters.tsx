import React from 'react';
import { Badge } from '@/components/ui';
import { Link } from '@/router';

interface FilterSearch {
  referrer?: string;
  device?: string;
  page?: string;
  country?: string;
}

interface ActiveFiltersProps {
  siteId: string;
  search: FilterSearch;
}

export function ActiveFilters({ siteId, search }: ActiveFiltersProps): React.JSX.Element | null {
  const hasFilters = Boolean(search.referrer ?? search.device ?? search.page ?? search.country);
  if (!hasFilters) {
    return null;
  }

  return (
    <div className="flex items-center gap-2 flex-wrap">
      <span className="text-sm text-muted-foreground">Filtered by:</span>
      {search.referrer && (
        <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Referrer: {search.referrer}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      )}
      {search.device && (
        <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Device: {search.device}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      )}
      {search.page && (
        <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Page: {search.page}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      )}
      {search.country && (
        <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Country: {search.country}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      )}
      <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
        <Badge variant="outline" className="cursor-pointer hover:bg-accent text-xs">
          Clear all
        </Badge>
      </Link>
    </div>
  );
}
