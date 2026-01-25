
import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui';
import { Trash2 } from 'lucide-react';

interface DangerZoneCardProps {
  deleting: boolean;
  confirming: boolean;
  onDeleteRequest: () => void;
  onConfirmDelete: () => void;
  onCancelDelete: () => void;
}

export const DangerZoneCard = ({
  deleting,
  confirming,
  onDeleteRequest,
  onConfirmDelete,
  onCancelDelete,
}: DangerZoneCardProps): React.ReactNode => (
  <Card className="border-destructive/50">
    <CardHeader>
      <CardTitle className="text-destructive">Danger Zone</CardTitle>
      <CardDescription>
        Irreversible actions that affect this site
      </CardDescription>
    </CardHeader>
    <CardContent>
      {confirming ? (
        <div className="flex flex-wrap gap-2">
          <Button
            variant="destructive"
            onClick={onConfirmDelete}
            disabled={deleting}
          >
            <Trash2 className="h-4 w-4 mr-2" />
            {deleting ? 'Deleting...' : 'Confirm Delete'}
          </Button>
          <Button
            variant="outline"
            onClick={onCancelDelete}
            disabled={deleting}
          >
            Cancel
          </Button>
        </div>
      ) : (
        <Button
          variant="destructive"
          onClick={onDeleteRequest}
          disabled={deleting}
        >
          <Trash2 className="h-4 w-4 mr-2" />
          Delete Site
        </Button>
      )}
      <p className="text-xs text-muted-foreground mt-2">
        This will permanently delete all analytics data for this site
      </p>
    </CardContent>
  </Card>
)
