import { cn } from "@/lib/utils"

const SKELETON_KEYS = ['skeleton-1', 'skeleton-2', 'skeleton-3', 'skeleton-4', 'skeleton-5'] as const

const Skeleton = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => {
  return (
    <div
      className={cn("animate-pulse rounded-md bg-muted/80 dark:bg-muted", className)}
      {...props}
    />
  )
}

export { Skeleton, SKELETON_KEYS }
