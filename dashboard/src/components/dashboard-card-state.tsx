import { Loader2 } from 'lucide-react';
import type { FunctionComponent, ReactNode } from 'react';
import { cn } from '@/lib/utils';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface DashboardCardStateProps {
  state: DashboardLoadState;
  skeleton: ReactNode;
  children: ReactNode;
  className?: string;
  overlayLabel?: string;
}

const DashboardCardState: FunctionComponent<DashboardCardStateProps> = ({
  state,
  skeleton,
  children,
  className,
  overlayLabel = 'Updating',
}) => {
  if (state === 'initial') {
    return <div className={className}>{skeleton}</div>;
  }

  return (
    <div className="relative">
      <div className={cn(className, state === 'refreshing' && 'opacity-80 transition-opacity duration-200')}>
        {children}
      </div>
      {state === 'refreshing' ? (
        <div className="pointer-events-none absolute right-3 top-3">
          <div className="inline-flex items-center gap-2 rounded-full bg-background/95 px-3 py-1 text-xs text-muted-foreground shadow-sm">
            <Loader2 className="h-3 w-3 animate-spin" />
            <span>{overlayLabel}</span>
          </div>
        </div>
      ) : null}
    </div>
  );
};

export default DashboardCardState;
