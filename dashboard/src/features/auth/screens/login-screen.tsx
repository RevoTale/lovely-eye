import { Link, useNavigate } from '@tanstack/react-router';
import { type FormEvent, type ReactElement, useState } from 'react';
import { useAuth } from '@/features/auth/model/auth-context';
import { AuthShell } from '@/features/auth/ui/auth-shell';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';

export const LoginScreen = (): ReactElement => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { login, authMode, bootstrapError } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent): void => {
    e.preventDefault();
    setError(null);
    setIsSubmitting(true);

    void login({ username, password })
      .then(() => {
        void navigate({ to: '/' });
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Login failed');
      })
      .finally(() => {
        setIsSubmitting(false);
      });
  };

  const hasError = error !== null && error !== '';
  const canRegister = authMode === 'login-and-register';

  if (bootstrapError !== null) {
    return (
      <AuthShell title='Authentication unavailable' description={bootstrapError}>
        <p className='text-center text-sm text-muted-foreground'>Refresh the page to retry.</p>
      </AuthShell>
    );
  }

  if (authMode === 'register-only') {
    return (
      <AuthShell
        title='Redirecting to setup'
        description='No users exist yet, so registration is required first.'
      >
        <p className='text-center text-sm text-muted-foreground'>Please wait...</p>
      </AuthShell>
    );
  }

  return (
    <AuthShell
      title='Welcome back'
      description='Enter your credentials to access the dashboard'
      footer={
        canRegister ? (
          <>
            Don&apos;t have an account?{' '}
            <Link to='/register' className='text-primary hover:underline'>
              Register
            </Link>
          </>
        ) : undefined
      }
    >
      <form onSubmit={handleSubmit} className='space-y-4'>
        {hasError ? (
          <div className='bg-destructive/10 text-destructive text-sm p-3 rounded-md'>{error}</div>
        ) : null}
        <div className='space-y-2'>
          <Label htmlFor='username'>Username</Label>
          <Input
            id='username'
            type='text'
            placeholder='Enter your username'
            value={username}
            onChange={(e) => {
              setUsername(e.target.value);
            }}
            required
            autoComplete='username'
          />
        </div>
        <div className='space-y-2'>
          <Label htmlFor='password'>Password</Label>
          <Input
            id='password'
            type='password'
            placeholder='Enter your password'
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
            }}
            required
            autoComplete='current-password'
          />
        </div>
        <Button type='submit' className='w-full' disabled={isSubmitting}>
          {isSubmitting ? 'Signing in...' : 'Sign in'}
        </Button>
      </form>
    </AuthShell>
  );
};
