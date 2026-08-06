export type DashboardLoadState = 'initial' | 'refreshing' | 'ready';

interface QueryDisplayState<T> {
  data: T | undefined;
  state: DashboardLoadState;
}

export const isInitialLoadState = (state: DashboardLoadState): boolean => state === 'initial';

export const isRefreshingLoadState = (state: DashboardLoadState): boolean => state === 'refreshing';

export const resolveDashboardLoadState = <T>(
  data: T | undefined,
  previousData: T | undefined,
  loading: boolean
): QueryDisplayState<T> => {
  const displayData = data ?? previousData;
  if (loading && displayData === undefined) {
    return { data: undefined, state: 'initial' };
  }
  if (loading) {
    return { data: displayData, state: 'refreshing' };
  }
  return { data: displayData, state: 'ready' };
};
