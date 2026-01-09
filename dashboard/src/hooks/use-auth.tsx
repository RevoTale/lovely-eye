import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react';
import { useMutation, useQuery } from '@apollo/client';
import { ME_QUERY, LOGIN_MUTATION, LOGOUT_MUTATION, REGISTER_MUTATION } from '@/graphql';
import type { User, LoginInput, RegisterInput } from '@/generated/graphql';

export interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (input: LoginInput) => Promise<void>;
  register: (input: RegisterInput) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }): React.JSX.Element {
  const [user, setUser] = useState<User | null>(null);
  
  const { loading: meLoading, data: meData, refetch } = useQuery(ME_QUERY, {
    fetchPolicy: 'network-only',
    errorPolicy: 'ignore',
  });

  const [loginMutation] = useMutation(LOGIN_MUTATION);
  const [registerMutation] = useMutation(REGISTER_MUTATION);
  const [logoutMutation] = useMutation(LOGOUT_MUTATION);

  useEffect(() => {
    if (meData !== undefined) {
      setUser(meData?.me ? (meData.me as User) : null);
    }
  }, [meData]);

  const login = useCallback(async (input: LoginInput) => {
    const result = await loginMutation({ variables: { input } });
    if (result.data?.login?.user) {
      setUser(result.data.login.user as User);
    }
  }, [loginMutation]);

  const register = useCallback(async (input: RegisterInput) => {
    const result = await registerMutation({ variables: { input } });
    if (result.data?.register?.user) {
      setUser(result.data.register.user as User);
    }
  }, [registerMutation]);

  const logout = useCallback(async () => {
    await logoutMutation();
    setUser(null);
    await refetch();
  }, [logoutMutation, refetch]);

  const value: AuthContextType = {
    user,
    isLoading: meLoading,
    isAuthenticated: user !== null,
    login,
    register,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
