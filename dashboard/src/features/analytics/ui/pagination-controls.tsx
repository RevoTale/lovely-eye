import { useLayoutEffect, useRef } from 'react';
import type { DashboardLoadState } from '@/features/analytics/model/dashboard-load-state';
import { Button } from '@/shared/ui/button';

interface PaginationControlsProps {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;
const MIN_PAGE = 1;
const PAGE_INCREMENT = 1;
const VIEWPORT_START = 0;

export const PaginationControls = ({
  page,
  pageSize,
  total,
  onPageChange,
  state = 'ready',
}: PaginationControlsProps): React.ReactNode => {
  const rootRef = useRef<HTMLDivElement>(null);
  const previousPageRef = useRef(page);
  const pendingVisibilityCheckRef = useRef(false);
  const totalPages = Math.max(MIN_PAGE, Math.ceil(total / pageSize));
  const clampedPage = Math.min(Math.max(page, MIN_PAGE), totalPages);
  const start =
    total === EMPTY_COUNT
      ? EMPTY_COUNT
      : (clampedPage - PAGE_INCREMENT) * pageSize + PAGE_INCREMENT;
  const end = total === EMPTY_COUNT ? EMPTY_COUNT : Math.min(clampedPage * pageSize, total);

  useLayoutEffect(() => {
    if (previousPageRef.current !== page) {
      previousPageRef.current = page;
      pendingVisibilityCheckRef.current = true;
    }
    if (state !== 'ready' || !pendingVisibilityCheckRef.current) return;

    pendingVisibilityCheckRef.current = false;
    const card = rootRef.current?.closest<HTMLElement>('[data-slot="card"]');
    if (card === undefined || card === null) return;
    const { top } = card.getBoundingClientRect();
    if (top < VIEWPORT_START || top >= window.innerHeight) {
      card.scrollIntoView({ block: 'start' });
    }
  }, [page, state]);

  return (
    <div
      ref={rootRef}
      className='flex flex-wrap items-center justify-between gap-2 text-xs text-muted-foreground'
    >
      <span>{total === EMPTY_COUNT ? 'No results' : `Showing ${start}-${end} of ${total}`}</span>
      <div className='flex items-center gap-2'>
        <Button
          type='button'
          variant='outline'
          size='sm'
          disabled={clampedPage <= MIN_PAGE}
          onClick={() => {
            onPageChange(Math.max(MIN_PAGE, clampedPage - PAGE_INCREMENT));
          }}
        >
          Prev
        </Button>
        <span>
          Page {clampedPage} of {totalPages}
        </span>
        <Button
          type='button'
          variant='outline'
          size='sm'
          disabled={clampedPage >= totalPages}
          onClick={() => {
            onPageChange(Math.min(totalPages, clampedPage + PAGE_INCREMENT));
          }}
        >
          Next
        </Button>
      </div>
    </div>
  );
};
