import React from 'react';
import { Button } from '@/components/ui';

interface PaginationControlsProps {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
}

export function PaginationControls({
  page,
  pageSize,
  total,
  onPageChange,
}: PaginationControlsProps): React.JSX.Element {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const clampedPage = Math.min(Math.max(page, 1), totalPages);
  const start = total === 0 ? 0 : (clampedPage - 1) * pageSize + 1;
  const end = total === 0 ? 0 : Math.min(clampedPage * pageSize, total);

  return (
    <div className="flex flex-wrap items-center justify-between gap-2 text-xs text-muted-foreground">
      <span>
        {total === 0 ? 'No results' : `Showing ${start}-${end} of ${total}`}
      </span>
      <div className="flex items-center gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={clampedPage <= 1}
          onClick={() => {
            onPageChange(Math.max(1, clampedPage - 1));
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
            onPageChange(Math.min(totalPages, clampedPage + 1));
          }}
        >
          Next
        </Button>
      </div>
    </div>
  );
}
