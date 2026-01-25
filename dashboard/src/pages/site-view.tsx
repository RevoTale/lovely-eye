
import { useParams, useSearch } from '@tanstack/react-router';
import { SiteFormPage } from './site-form';
import { DashboardPage } from './dashboard';
import { siteDetailRoute } from '@/router';

export const SiteViewPage = (): React.ReactNode => {
  const { siteId } = useParams({ from: siteDetailRoute.id });
  const search = useSearch({ from: siteDetailRoute.id });

  const showForm = siteId === 'new' || (search as { view?: string }).view === 'settings';

  if (showForm) {
    return <SiteFormPage />;
  }

  return <DashboardPage />;
}
