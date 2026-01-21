const FIRST_PAGE = 1;
const PAGINATION_KEYS = new Set(['eventsPage', 'topPagesPage', 'referrersPage', 'devicesPage', 'countriesPage']);

export function clearPaginationParams(prev: Record<string, unknown>): Record<string, unknown> {
  return Object.fromEntries(
    Object.entries(prev).filter(([key]) => !PAGINATION_KEYS.has(key))
  );
}

export function updatePageParam(prev: Record<string, unknown>, key: string, nextPage: number): Record<string, unknown> {
  if (nextPage <= FIRST_PAGE) {
    return Object.fromEntries(
      Object.entries(prev).filter(([entryKey]) => entryKey !== key)
    );
  }
  return { ...prev, [key]: String(nextPage) };
}
