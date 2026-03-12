import type { ReactNode } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui';
import { Logo } from '@/components/logo';

interface AuthShellProps {
  title: string;
  description: string;
  children: ReactNode;
  footer?: ReactNode;
}

export const AuthShell = ({ title, description, children, footer }: AuthShellProps): ReactNode => (
  <div className="min-h-screen flex items-center justify-center bg-muted/50 px-4">
    <Card className="w-full max-w-md">
      <CardHeader className="space-y-1 text-center">
        <div className="flex justify-center mb-4">
          <div className="flex items-center gap-2">
            <Logo size={32} />
            <span className="text-2xl font-bold">Lovely Eye</span>
          </div>
        </div>
        <CardTitle className="text-2xl">{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        {children}
        {footer !== undefined ? (
          <div className="mt-4 text-center text-sm text-muted-foreground">
            {footer}
          </div>
        ) : null}
      </CardContent>
    </Card>
  </div>
);
