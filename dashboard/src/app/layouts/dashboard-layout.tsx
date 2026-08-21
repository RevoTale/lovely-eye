import { Outlet } from '@tanstack/react-router';
import { LayoutDashboard, LogOut } from 'lucide-react';
import type { ReactElement } from 'react';
import { Link, useNavigate } from '@/app/router';
import { ThemeToggle } from '@/app/ui/theme-toggle';
import { useAuth } from '@/features/auth/model/auth-context';
import { Avatar, AvatarFallback } from '@/shared/ui/avatar';
import { Button } from '@/shared/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu';
import { Logo } from '@/shared/ui/logo';
import { Separator } from '@/shared/ui/separator';

const LOGO_SIZE = 32;
const USERNAME_INITIALS_START = 0;
const USERNAME_INITIALS_END = 2;

export const DashboardLayout = (): ReactElement => {
  const { user, logout, unauthenticatedRoute } = useAuth();
  const navigate = useNavigate();

  const handleLogout = (): void => {
    void logout().then(() => {
      void navigate({ to: unauthenticatedRoute });
    });
  };

  const initials =
    user?.username.slice(USERNAME_INITIALS_START, USERNAME_INITIALS_END).toUpperCase() ?? 'U';

  return (
    <div className='min-h-screen bg-background'>
      <header className='sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 shadow-sm'>
        <div className='mx-auto flex h-14 w-full max-w-screen-2xl items-center gap-2 px-4 sm:h-16 sm:px-6 lg:px-8'>
          <Link to='/' className='flex items-center gap-1.5 mr-0 sm:gap-2 sm:mr-8'>
            <Logo size={LOGO_SIZE} />
            <span className='font-bold text-lg hidden sm:inline'>Lovely Eye</span>
          </Link>

          <div className='flex items-center gap-1.5 ml-auto sm:gap-2'>
            <ThemeToggle />
            <nav className='flex items-center gap-2 sm:gap-3'>
              <Link
                to='/sites'
                className='flex items-center gap-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors'
              >
                <LayoutDashboard className='h-4 w-4' />
                <span className='hidden xs:inline-block'>Sites</span>
              </Link>
            </nav>

            <Separator orientation='vertical' className='h-5 mx-1.5 sm:h-6 sm:mx-2' />

            <DropdownMenu>
              <DropdownMenuTrigger
                render={
                  <Button
                    variant='ghost'
                    className='relative h-9 w-9 rounded-full'
                    aria-label='Open user menu'
                  />
                }
              >
                <Avatar className='h-9 w-9'>
                  <AvatarFallback className='bg-primary/10 text-primary font-semibold'>
                    {initials}
                  </AvatarFallback>
                </Avatar>
              </DropdownMenuTrigger>
              <DropdownMenuContent className='w-56' align='end'>
                <DropdownMenuGroup>
                  <DropdownMenuLabel className='font-normal'>
                    <div className='flex flex-col space-y-1'>
                      <p className='text-sm font-medium leading-none'>{user?.username}</p>
                      <p className='text-xs leading-none text-muted-foreground capitalize'>
                        {user?.role}
                      </p>
                    </div>
                  </DropdownMenuLabel>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onClick={handleLogout}
                  className='text-destructive focus:text-destructive'
                >
                  <LogOut className='mr-2 h-4 w-4' />
                  <span>Log out</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      <main className='mx-auto w-full min-w-0 max-w-screen-2xl px-4 py-6 sm:px-6 sm:py-8 lg:px-8'>
        <Outlet />
      </main>
    </div>
  );
};
