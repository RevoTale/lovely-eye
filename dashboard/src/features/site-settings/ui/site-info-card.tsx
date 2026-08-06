import { Globe } from 'lucide-react';
import type { ReactElement } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import DomainFields from './site-info/domain-fields';
import SiteInfoActions from './site-info/site-info-actions';
import type { SiteInfoCardProps } from './site-info/types';
import { useSiteInfoForm } from './site-info/use-site-info-form';

export const SiteInfoCard = ({
  initialDomains,
  initialName,
  onSaveDomains,
  updating,
}: SiteInfoCardProps): ReactElement => {
  const form = useSiteInfoForm({ initialDomains, initialName, onSaveDomains });

  return (
    <form onSubmit={(event) => void form.handleSubmit(event)}>
      <Card>
        <CardHeader>
          <CardTitle className='flex items-center gap-2'>
            <div className='flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10'>
              <Globe className='h-4 w-4 text-primary' />
            </div>
            Site Information
          </CardTitle>
          <CardDescription>Site configuration and tracking details</CardDescription>
        </CardHeader>
        <CardContent className='space-y-6'>
          <DomainFields
            domains={form.domains}
            onAddDomain={form.addDomain}
            onDomainChange={form.handleDomainChange}
            onRemoveDomain={form.removeDomain}
          />
          <div className='space-y-2'>
            <Label htmlFor='name'>Site Name</Label>
            <Input id='name' placeholder='My Awesome Website' value={form.name} disabled required />
            <p className='text-xs text-muted-foreground'>A friendly name to identify your site</p>
          </div>
          <SiteInfoActions hasDomainChanges={form.hasDomainChanges} updating={updating} />
          {form.formError === '' ? null : (
            <p className='text-xs text-destructive'>{form.formError}</p>
          )}
        </CardContent>
      </Card>
    </form>
  );
};
