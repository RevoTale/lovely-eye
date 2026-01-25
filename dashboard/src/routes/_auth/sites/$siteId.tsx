import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';
import { z } from 'zod';

const filterValue = z.union([z.string(), z.array(z.string())]).optional();
const searchSchema = z.object({
  view: z.string().optional(),
  preset: z.enum(['7d', '30d', '90d', 'custom', 'all']).optional(),
  from: z.string().optional(),
  to: z.string().optional(),
  fromTime: z.string().optional(),
  toTime: z.string().optional(),
  statsBucket: z.enum(['daily', 'hourly']).optional(),
  eventsPage: z.string().optional(),
  eventsCountsPage: z.string().optional(),
  topPagesPage: z.string().optional(),
  referrersPage: z.string().optional(),
  devicesPage: z.string().optional(),
  countriesPage: z.string().optional(),
  referrer: filterValue,
  device: filterValue,
  page: filterValue,
  country: filterValue,
  eventName: filterValue,
  eventPath: filterValue,
});

export const Route = createFileRoute('/_auth/sites/$siteId')({
  validateSearch: (search: Record<string, unknown>) => {
    const parsed = searchSchema.safeParse(search);
    return parsed.success ? parsed.data : {};
  },
  component: lazyRouteComponent(async () => {
    const module = await import('../../../pages/site-view');
    return { default: module.SiteViewPage };
  }),
});
