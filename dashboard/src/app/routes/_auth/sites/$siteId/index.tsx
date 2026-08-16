import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/sites/$siteId/')({
  beforeLoad: ({ params }) => {
    throw redirect({ to: '/sites/$siteId/analytics', params });
  },
});
