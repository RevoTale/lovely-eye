import { useMutation } from '@apollo/client/react';
import { Link, useNavigate } from '@tanstack/react-router';
import { ArrowLeft, Globe, Save } from 'lucide-react';
import { rememberSite } from '@/features/sites/model/recent-site';
import { useCreateSiteForm } from '@/features/sites/model/use-create-site-form';
import { SITES_PAGING } from '@/features/sites/model/use-sites';
import { DomainFields } from '@/features/sites/ui/domain-fields';
import {
  CreateSiteDocument,
  SiteSummaryFieldsFragmentDoc,
  SitesDocument,
} from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';
import { Button } from '@/shared/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';

export function CreateSiteScreen(): React.ReactNode {
  const navigate = useNavigate();
  const [createSite, { loading }] = useMutation(CreateSiteDocument, {
    refetchQueries: [{ query: SitesDocument, variables: { paging: SITES_PAGING } }],
  });
  const form = useCreateSiteForm({
    onCreate: async (name, domains) => {
      const { data } = await createSite({ variables: { input: { name, domains } } });
      if (data === undefined) throw new Error('Site creation returned no data');
      const site = readFragment(SiteSummaryFieldsFragmentDoc, data.createSite);
      rememberSite(site.id);
      await navigate({ to: '/sites/$siteId/analytics', params: { siteId: site.id } });
    },
  });

  return (
    <div className='mx-auto w-full max-w-3xl space-y-6'>
      <Button variant='outline' size='sm' asChild>
        <Link to='/sites'>
          <ArrowLeft className='size-4' />
          Back to Sites
        </Link>
      </Button>
      <div>
        <h1 className='text-3xl font-bold tracking-tight'>Add New Site</h1>
        <p className='mt-1 text-muted-foreground'>Create a site and connect all of its domains.</p>
      </div>
      <form onSubmit={(event) => void form.submit(event)}>
        <Card>
          <CardHeader>
            <CardTitle className='flex items-center gap-2'>
              <Globe className='size-5 text-primary' />
              Site Information
            </CardTitle>
            <CardDescription>Each domain is an alias of the same analytics site.</CardDescription>
          </CardHeader>
          <CardContent className='space-y-6'>
            <DomainFields
              domains={form.domains}
              onAdd={form.addDomain}
              onChange={form.changeDomain}
              onRemove={form.removeDomain}
            />
            <div className='space-y-2'>
              <Label htmlFor='site-name'>Site Name</Label>
              <Input
                id='site-name'
                placeholder='My Website'
                value={form.name}
                onChange={(event) => form.setName(event.target.value)}
                required
              />
            </div>
            {form.error === null ? null : <p className='text-sm text-destructive'>{form.error}</p>}
            <div className='flex flex-col gap-3 sm:flex-row'>
              <Button type='submit' className='w-full sm:w-auto' disabled={loading}>
                <Save className='size-4' />
                {loading ? 'Creating...' : 'Create Site'}
              </Button>
              <Button
                type='button'
                variant='outline'
                className='w-full sm:w-auto'
                disabled={loading}
                asChild
              >
                <Link to='/sites'>Cancel</Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      </form>
    </div>
  );
}
