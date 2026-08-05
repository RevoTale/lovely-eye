
import { Card, CardContent, CardHeader, Skeleton, SKELETON_KEYS } from '@/components/ui';

const STAT_PLACEHOLDER_COUNT = 4;

export const DashboardLoading = (): React.ReactNode => (
  <div className="space-y-6">
    <div className="space-y-2">
      <Skeleton className="h-8 w-64" />
      <Skeleton className="h-4 w-48" />
    </div>
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {SKELETON_KEYS.slice(0, STAT_PLACEHOLDER_COUNT).map((key) => (
        <Card key={key}>
          <CardHeader>
            <Skeleton className="h-4 w-24" />
          </CardHeader>
          <CardContent>
            <Skeleton className="h-8 w-32" />
          </CardContent>
        </Card>
      ))}
    </div>
  </div>
)

export const DashboardNotFound = (): React.ReactNode => (
  <div className="flex items-center justify-center min-h-[400px]">
    <div className="text-destructive">Site not found</div>
  </div>
)

export const DashboardEmptyState = (): React.ReactNode => (
  <Card>
    <CardContent className="py-12">
      <div className="text-center text-muted-foreground">
        <p>No analytics data yet. Add the tracking script to start collecting data.</p>
      </div>
    </CardContent>
  </Card>
)
