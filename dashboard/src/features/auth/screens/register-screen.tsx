import { Link, useNavigate } from '@tanstack/react-router';
import { type FormEvent, type ReactElement, useState } from 'react';
import { useAuth } from '@/features/auth/model/auth-context';
import { AuthShell } from '@/features/auth/ui/auth-shell';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';

export const RegisterScreen = (): ReactElement => {
  const MIN_PASSWORD_LENGTH = 8;
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { register, authMode, bootstrapError } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = (e: FormEvent): void => {
    e.preventDefault();
    setError(null);

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    if (password.length < MIN_PASSWORD_LENGTH) {
      setError(`Password must be at least ${MIN_PASSWORD_LENGTH} characters`);
      return;
    }

    setIsSubmitting(true);

    void register({ username, password })
      .then(() => {
        void navigate({ to: '/' });
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Registration failed');
      })
      .finally(() => {
        setIsSubmitting(false);
      });
  };

  const hasError = error !== null && error !== '';
  const isInitialSetup = authMode === 'register-only';

  if (bootstrapError !== null) {
    return (
      <AuthShell title='Authentication unavailable' description={bootstrapError}>
        <p className='text-center text-sm text-muted-foreground'>Refresh the page to retry.</p>
      </AuthShell>
    );
  }

  if (authMode === 'login-only') {
    return (
      <AuthShell title='Redirecting to sign in' description='Registration is currently disabled.'>
        <p className='text-center text-sm text-muted-foreground'>Please wait...</p>
      </AuthShell>
    );
  }

  return (
    <AuthShell
      title={isInitialSetup ? 'Create the initial admin account' : 'Create an account'}
      description={
        isInitialSetup
          ? 'No users exist yet. This first account will become the admin.'
          : 'Register to start tracking your websites'
      }
      footer={
        authMode === 'login-and-register' ? (
          <>
            Already have an account?{' '}
            <Link to='/login' className='text-primary hover:underline'>
              Sign in
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
            placeholder='Choose a username'
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
            placeholder='Create a password'
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
            }}
            required
            autoComplete='new-password'
          />
        </div>
        <div className='space-y-2'>
          <Label htmlFor='confirmPassword'>Confirm Password</Label>
          <Input
            id='confirmPassword'
            type='password'
            placeholder='Confirm your password'
            value={confirmPassword}
            onChange={(e) => {
              setConfirmPassword(e.target.value);
            }}
            required
            autoComplete='new-password'
          />
        </div>
        <Button type='submit' className='w-full' disabled={isSubmitting}>
          {isSubmitting ? 'Creating account...' : 'Create account'}
        </Button>
      </form>
    </AuthShell>
  );
};
