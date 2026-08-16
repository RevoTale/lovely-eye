import { Loader2, Save } from 'lucide-react';
import type { FunctionComponent } from 'react';
import { Button } from '@/shared/ui/button';

interface SiteInfoActionsProps {
  hasDomainChanges: boolean;
  updating: boolean;
}

const SiteInfoActions: FunctionComponent<SiteInfoActionsProps> = ({
  hasDomainChanges,
  updating,
}) => (
  <div className='flex gap-3 pt-4'>
    <Button type='submit' className='w-full sm:w-auto' disabled={updating || !hasDomainChanges}>
      {updating ? (
        <>
          <Loader2 className='mr-2 h-4 w-4 animate-spin' />
          Saving...
        </>
      ) : (
        <>
          <Save className='mr-2 h-4 w-4' />
          Save Domains
        </>
      )}
    </Button>
  </div>
);

export default SiteInfoActions;
