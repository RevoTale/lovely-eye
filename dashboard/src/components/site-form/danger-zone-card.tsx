import React from 'react';
import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui';
import { Trash2 } from 'lucide-react';

interface DangerZoneCardProps {
  deleting: boolean;
  onDelete: () => void;
}

export function DangerZoneCard({
  deleting,
  onDelete,
}: DangerZoneCardProps): React.JSX.Element {
  return (
    <Card className="border-destructive/50">
      <CardHeader>
        <CardTitle className="text-destructive">Danger Zone</CardTitle>
        <CardDescription>
          Irreversible actions that affect this site
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Button
          variant="destructive"
          onClick={onDelete}
          disabled={deleting}
        >
          <Trash2 className="h-4 w-4 mr-2" />
          {deleting ? 'Deleting...' : 'Delete Site'}
        </Button>
        <p className="text-xs text-muted-foreground mt-2">
          This will permanently delete all analytics data for this site
        </p>
      </CardContent>
    </Card>
  );
}
