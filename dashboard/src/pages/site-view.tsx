import React from 'react';
import { useParams, useSearch } from '@tanstack/react-router';
import { SiteFormPage } from './site-form';
import { DashboardPage } from './dashboard';
import { siteDetailRoute } from '@/router';

export function SiteViewPage(): React.JSX.Element {
  const { siteId } = useParams({ from: siteDetailRoute.id });
  const search = useSearch({ from: siteDetailRoute.id });

  // Show form page if creating new site or if view=settings
  const showForm = siteId === 'new' || (search as { view?: string }).view === 'settings';

  if (showForm) {
    return <SiteFormPage />;
  }

  return <DashboardPage />;
}
