
import { Button } from '@/components/ui';

interface PaginationControlsProps {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
}

const EMPTY_COUNT = 0;
const MIN_PAGE = 1;
const PAGE_INCREMENT = 1;

export const PaginationControls = ({
  page,
  pageSize,
  total,
  onPageChange,
}: PaginationControlsProps): React.ReactNode => {
  const totalPages = Math.max(MIN_PAGE, Math.ceil(total / pageSize));
  const clampedPage = Math.min(Math.max(page, MIN_PAGE), totalPages);
  const start =
    total === EMPTY_COUNT
      ? EMPTY_COUNT
      : (clampedPage - PAGE_INCREMENT) * pageSize + PAGE_INCREMENT;
  const end =
    total === EMPTY_COUNT ? EMPTY_COUNT : Math.min(clampedPage * pageSize, total);

  return (
    <div className="flex flex-wrap items-center justify-between gap-2 text-xs text-muted-foreground">
      <span>
        {total === EMPTY_COUNT ? 'No results' : `Showing ${start}-${end} of ${total}`}
      </span>
      <div className="flex items-center gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
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
          type="button"
          variant="outline"
          size="sm"
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
}
