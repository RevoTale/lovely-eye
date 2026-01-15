import React from 'react';
import { cn } from '@/lib/utils';

interface ListEmptyStateProps {
  title: string;
  description?: string;
  icon?: React.ElementType;
  className?: string;
}

export function ListEmptyState({
  title,
  description,
  icon: Icon,
  className,
}: ListEmptyStateProps): React.JSX.Element {
  return (
    <div className={cn('text-center py-6 text-muted-foreground', className)}>
      {Icon ? <Icon className="h-10 w-10 mx-auto mb-3 opacity-50" /> : null}
      <p className="text-sm">{title}</p>
      {description ? <p className="text-xs mt-1">{description}</p> : null}
    </div>
  );
}
