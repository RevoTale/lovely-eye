
import { cn } from '@/lib/utils';

interface ListEmptyStateProps {
  title: string;
  description?: string;
  icon?: React.ElementType;
  className?: string;
}

export const ListEmptyState = ({
  title,
  description,
  icon: Icon,
  className,
}: ListEmptyStateProps): React.ReactNode => {
  const hasDescription = description !== undefined && description !== '';
  const hasIcon = Icon !== undefined;

  return (
    <div className={cn('text-center py-6 text-muted-foreground', className)}>
      {hasIcon ? <Icon className="h-10 w-10 mx-auto mb-3 opacity-50" /> : null}
      <p className="text-sm">{title}</p>
      {hasDescription ? <p className="text-xs mt-1">{description}</p> : null}
    </div>
  );
}
