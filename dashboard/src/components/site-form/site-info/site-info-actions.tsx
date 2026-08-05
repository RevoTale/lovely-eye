import type { FunctionComponent } from 'react';
import { Loader2, Save } from 'lucide-react';
import { Button } from '@/components/ui';

interface SiteInfoActionsProps {
  creating: boolean;
  hasDomainChanges: boolean;
  isNew: boolean;
  updating: boolean;
  onCancel: () => void;
}

const SiteInfoActions: FunctionComponent<SiteInfoActionsProps> = ({
  creating,
  hasDomainChanges,
  isNew,
  onCancel,
  updating,
}) =>
  isNew ? (
    <div className="flex gap-3 pt-4">
      <Button type="submit" disabled={creating}>
        <Save className="mr-2 h-4 w-4" />
        {creating ? 'Creating...' : 'Create Site'}
      </Button>
      <Button type="button" variant="outline" onClick={onCancel} disabled={creating}>
        Cancel
      </Button>
    </div>
  ) : (
    <div className="flex gap-3 pt-4">
      <Button type="submit" disabled={updating || !hasDomainChanges}>
        {updating ? (
          <>
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            Saving...
          </>
        ) : (
          <>
            <Save className="mr-2 h-4 w-4" />
            Save Domains
          </>
        )}
      </Button>
    </div>
  );

export default SiteInfoActions;
