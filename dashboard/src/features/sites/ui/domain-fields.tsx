import { Plus, X } from 'lucide-react';
import type { DomainField } from '@/features/sites/model/use-create-site-form';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';

interface DomainFieldsProps {
  domains: DomainField[];
  onAdd: () => void;
  onChange: (id: number, value: string) => void;
  onRemove: (id: number) => void;
}

export function DomainFields({
  domains,
  onAdd,
  onChange,
  onRemove,
}: DomainFieldsProps): React.ReactNode {
  return (
    <div className='space-y-2'>
      <Label htmlFor='primary-domain'>Domains</Label>
      <div className='space-y-2'>
        {domains.map((domain, index) => (
          <div key={domain.id} className='flex items-center gap-2'>
            <Input
              id={index === 0 ? 'primary-domain' : `domain-${index}`}
              placeholder={index === 0 ? 'example.com' : 'blog.example.com'}
              value={domain.value}
              onChange={(event) => onChange(domain.id, event.target.value)}
              required={index === 0}
            />
            {domains.length > 1 ? (
              <Button
                type='button'
                variant='outline'
                size='icon'
                onClick={() => onRemove(domain.id)}
                aria-label='Remove domain'
              >
                <X className='size-4' />
              </Button>
            ) : null}
          </div>
        ))}
      </div>
      <Button type='button' variant='outline' size='sm' onClick={onAdd}>
        <Plus className='size-4' />
        Add domain
      </Button>
      <p className='text-xs text-muted-foreground'>
        Add domains without https://. Every domain contributes to the same site analytics.
      </p>
    </div>
  );
}
