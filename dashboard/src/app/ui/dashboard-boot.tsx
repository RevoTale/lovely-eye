import { AuthShell } from '@/features/auth/ui/auth-shell';

export const DashboardBoot = (): React.ReactNode => (
  <AuthShell title='Loading dashboard' description='Preparing your dashboard.'>
    <p className='text-center text-sm text-muted-foreground'>Please wait...</p>
  </AuthShell>
);
