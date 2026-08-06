import { Loader2, X } from 'lucide-react';
import type { Dispatch, FunctionComponent, SetStateAction } from 'react';
import { normalizeIPInput } from '@/features/site-settings/ui/utils';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { EMPTY_STRING, MAX_IPS } from './constants';
import type { BlockedIPEntry } from './types';

interface BlockedIPSectionProps {
  blockedIPCount: number;
  blockedIPs: BlockedIPEntry[];
  ipActionError: string;
  newIPError: string;
  newIPValue: string;
  savingBlockedIPs: boolean;
  onAddIP: () => void;
  onRemoveBlockedIP: (value: string) => void;
  setNewIPError: Dispatch<SetStateAction<string>>;
  setNewIPValue: Dispatch<SetStateAction<string>>;
}

const BlockedIPSection: FunctionComponent<BlockedIPSectionProps> = ({
  blockedIPCount,
  blockedIPs,
  ipActionError,
  newIPError,
  newIPValue,
  onAddIP,
  onRemoveBlockedIP,
  savingBlockedIPs,
  setNewIPError,
  setNewIPValue,
}) => (
  <div className='space-y-3'>
    <div className='flex items-center justify-between'>
      <Label>Blocked IPs</Label>
      <span className='text-xs text-muted-foreground'>
        {blockedIPCount}/{MAX_IPS}
      </span>
    </div>
    <div className='space-y-2'>
      {blockedIPs.length === 0 ? (
        <span className='text-xs text-muted-foreground'>No blocked IPs yet.</span>
      ) : (
        blockedIPs.map((entry) => (
          <div key={entry.id} className='flex items-center gap-2'>
            <Input value={entry.value} readOnly />
            <Button
              type='button'
              variant='outline'
              size='icon'
              disabled={savingBlockedIPs}
              onClick={() => onRemoveBlockedIP(entry.value)}
              aria-label='Remove blocked IP'
            >
              <X className='h-4 w-4' />
            </Button>
          </div>
        ))
      )}
      <div className='space-y-1'>
        <div className='flex items-center gap-2'>
          <Input
            placeholder='203.0.113.10'
            value={newIPValue}
            onChange={(event) => {
              setNewIPValue(normalizeIPInput(event.currentTarget.value));
              if (newIPError !== EMPTY_STRING) setNewIPError(EMPTY_STRING);
            }}
            disabled={savingBlockedIPs}
          />
          <Button
            type='button'
            variant='outline'
            onClick={onAddIP}
            disabled={savingBlockedIPs || blockedIPCount >= MAX_IPS}
          >
            {savingBlockedIPs ? <Loader2 className='h-4 w-4 animate-spin' /> : 'Add IP'}
          </Button>
        </div>
        {newIPError === EMPTY_STRING ? null : (
          <span className='text-xs text-destructive'>{newIPError}</span>
        )}
      </div>
    </div>
    {savingBlockedIPs ? (
      <p className='text-xs text-muted-foreground'>Updating blocked IPs...</p>
    ) : null}
    {ipActionError === EMPTY_STRING ? null : (
      <p className='text-xs text-destructive'>{ipActionError}</p>
    )}
  </div>
);

export default BlockedIPSection;
