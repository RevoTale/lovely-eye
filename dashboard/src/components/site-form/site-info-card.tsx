import type { ReactElement } from 'react';
import { Globe } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, Input, Label } from '@/components/ui';
import DomainFields from './site-info/domain-fields';
import SiteInfoActions from './site-info/site-info-actions';
import type { SiteInfoCardProps } from './site-info/types';
import { useSiteInfoForm } from './site-info/use-site-info-form';

export const SiteInfoCard = ({
  creating,
  initialDomains,
  initialName,
  isNew,
  onCancel,
  onCreate,
  onSaveDomains,
  updating,
}: SiteInfoCardProps): ReactElement => {
  const form = useSiteInfoForm({ initialDomains, initialName, isNew, onCreate, onSaveDomains });

  return (
    <form onSubmit={(event) => void form.handleSubmit(event)}>
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
              <Globe className="h-4 w-4 text-primary" />
            </div>
            Site Information
          </CardTitle>
          <CardDescription>
            {isNew ? 'Enter your website details' : 'Site configuration and tracking details'}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <DomainFields
            domains={form.domains}
            onAddDomain={form.addDomain}
            onDomainChange={form.handleDomainChange}
            onRemoveDomain={form.removeDomain}
          />
          <div className="space-y-2">
            <Label htmlFor="name">Site Name</Label>
            <Input
              id="name"
              placeholder="My Awesome Website"
              value={form.name}
              onChange={(event) => form.setName(event.target.value)}
              disabled={!isNew}
              required
            />
            <p className="text-xs text-muted-foreground">A friendly name to identify your site</p>
          </div>
          <SiteInfoActions
            creating={creating}
            hasDomainChanges={form.hasDomainChanges}
            isNew={isNew}
            onCancel={onCancel}
            updating={updating}
          />
          {form.formError === '' ? null : <p className="text-xs text-destructive">{form.formError}</p>}
        </CardContent>
      </Card>
    </form>
  );
};
