import { useEffect, type ReactElement } from 'react';
import { Outlet } from '@tanstack/react-router';
import { useAuth } from '@/hooks';
import { Link, useNavigate } from '@/router';
import {
  Button,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  Avatar,
  AvatarFallback,
  Separator,
} from '@/components/ui';
import { LogOut, LayoutDashboard } from 'lucide-react';
import { Logo } from '@/components/logo';
import { ThemeToggle } from '@/components/theme-toggle';

const LOGO_SIZE = 32;
const USERNAME_INITIALS_START = 0;
const USERNAME_INITIALS_END = 2;

export const DashboardLayout = (): ReactElement => {
  const { user, logout, isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();

  const handleLogout = (): void => {
    void logout().then(() => {
      void navigate({ to: '/login' });
    });
  };

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      void navigate({ to: '/login' });
    }
  }, [isAuthenticated, isLoading, navigate]);

  const initials =
    user?.username.slice(USERNAME_INITIALS_START, USERNAME_INITIALS_END).toUpperCase() ?? 'U';

  return (
    <div className="min-h-screen bg-background">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 shadow-sm">
        <div className="container flex items-center gap-2 py-1.5 sm:gap-2 sm:py-2 sm:h-16">
          <Link to="/" className="flex items-center gap-1.5 mr-0 sm:gap-2 sm:mr-8">
            <Logo size={LOGO_SIZE} />
            <span className="font-bold text-lg hidden sm:inline">Lovely Eye</span>
          </Link>

          <div className="flex items-center gap-1.5 ml-auto sm:gap-2">
            <ThemeToggle />
            <nav className="flex items-center gap-2 sm:gap-3">
              <Link
                to="/"
                className="flex items-center gap-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
              >
                <LayoutDashboard className="h-4 w-4" />
                <span className="hidden xs:inline-block">Sites</span>
              </Link>
            </nav>

            <Separator orientation="vertical" className="h-5 mx-1.5 sm:h-6 sm:mx-2" />

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-9 w-9 rounded-full">
                  <Avatar className="h-9 w-9">
                    <AvatarFallback className="bg-primary/10 text-primary font-semibold">
                      {initials}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56" align="end" forceMount>
                <DropdownMenuLabel className="font-normal">
                  <div className="flex flex-col space-y-1">
                    <p className="text-sm font-medium leading-none">{user?.username}</p>
                    <p className="text-xs leading-none text-muted-foreground capitalize">
                      {user?.role}
                    </p>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                {
                  /* TODO implement later
                  <DropdownMenuItem>
                  <User className="mr-2 h-4 w-4" />
                  <span>Profile</span>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Settings className="mr-2 h-4 w-4" />
                  <span>Settings</span>
                </DropdownMenuItem>
                */
                }
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout} className="text-destructive focus:text-destructive">
                  <LogOut className="mr-2 h-4 w-4" />
                  <span>Log out</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      <main className="container py-8">
        <Outlet />
      </main>
    </div>
  );
}
