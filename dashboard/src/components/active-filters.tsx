import React from 'react';
import { Badge } from '@/components/ui';
import { Link } from '@/router';
import { normalizeFilterValue, removeFilterValue, updateFilterSearch } from '@/lib/filter-utils';

interface FilterSearch {
  referrer?: string | string[] | undefined;
  device?: string | string[] | undefined;
  page?: string | string[] | undefined;
  country?: string | string[] | undefined;
}

interface ActiveFiltersProps {
  siteId: string;
  search: FilterSearch;
}

const EMPTY_COUNT = 0;

export function ActiveFilters({ siteId, search }: ActiveFiltersProps): React.JSX.Element | null {
  const referrers = normalizeFilterValue(search.referrer);
  const devices = normalizeFilterValue(search.device);
  const pages = normalizeFilterValue(search.page);
  const countries = normalizeFilterValue(search.country);
  const hasFilters =
    referrers.length > EMPTY_COUNT ||
    devices.length > EMPTY_COUNT ||
    pages.length > EMPTY_COUNT ||
    countries.length > EMPTY_COUNT;
  if (!hasFilters) {
    return null;
  }

  return (
    <div className="flex items-center gap-2 flex-wrap">
      <span className="text-sm text-muted-foreground">Filtered by:</span>
      {referrers.map((referrer) => (
        <Link
          key={`referrer-${referrer}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'referrer', removeFilterValue(prev.referrer, referrer)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Referrer: {referrer}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      {devices.map((device) => (
        <Link
          key={`device-${device}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'device', removeFilterValue(prev.device, device)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Device: {device}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      {pages.map((page) => (
        <Link
          key={`page-${page}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'page', removeFilterValue(prev.page, page)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Page: {page}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      {countries.map((country) => (
        <Link
          key={`country-${country}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'country', removeFilterValue(prev.country, country)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Country: {country}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
        <Badge variant="outline" className="cursor-pointer hover:bg-accent text-xs">
          Clear all
        </Badge>
      </Link>
    </div>
  );
}
