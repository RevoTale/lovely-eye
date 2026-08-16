import { createFileRoute, Outlet } from '@tanstack/react-router';
import { rememberSite } from '@/features/sites/model/recent-site';

export const Route = createFileRoute('/_auth/sites/$siteId')({
  beforeLoad: ({ params }) => rememberSite(params.siteId),
  component: () => <Outlet />,
});
