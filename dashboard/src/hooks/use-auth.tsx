import { createContext, useContext, useState, useCallback, useRef, type MutableRefObject, type ReactNode } from 'react';
import { useMutation, useQuery, useApolloClient } from '@apollo/client';
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

interface AuthProviderProps {
  children: ReactNode;
  authErrorHandlerRef?: MutableRefObject<(() => void) | null>;
}

export function AuthProvider({ children, authErrorHandlerRef }: AuthProviderProps): React.JSX.Element {
  const [overrideUser, setOverrideUser] = useState<User | null>(null);
  const authErrorHandledRef = useRef(false);
  const client = useApolloClient();

  const { loading: meLoading, data: meData, refetch } = useQuery(ME_QUERY, {
    fetchPolicy: 'network-only',
    errorPolicy: 'ignore',
  });

  const [loginMutation] = useMutation(LOGIN_MUTATION);
  const [registerMutation] = useMutation(REGISTER_MUTATION);
  const [logoutMutation] = useMutation(LOGOUT_MUTATION);

  const handleAuthError = useCallback(() => {
    if (authErrorHandledRef.current) {
      return;
    }
    authErrorHandledRef.current = true;
    setOverrideUser(null);
    void client.clearStore();
  }, [client]);

  const handlerRef = authErrorHandlerRef;
  if (handlerRef) {
    handlerRef.current = handleAuthError;
  }

  const login = useCallback(async (input: LoginInput) => {
    const result = await loginMutation({ variables: { input } });
    if (result.data?.login?.user) {
      authErrorHandledRef.current = false;
      setOverrideUser(result.data.login.user as User);
      await refetch();
    }
  }, [loginMutation, refetch]);

  const register = useCallback(async (input: RegisterInput) => {
    const result = await registerMutation({ variables: { input } });
    if (result.data?.register?.user) {
      authErrorHandledRef.current = false;
      setOverrideUser(result.data.register.user as User);
      await refetch();
    }
  }, [registerMutation, refetch]);

  const logout = useCallback(async () => {
    await logoutMutation();
    setOverrideUser(null);
    await refetch();
  }, [logoutMutation, refetch]);

  const hasMeData = meData !== undefined;
  const meUser = hasMeData ? (meData?.me ? (meData.me as User) : null) : null;
  const user = hasMeData ? meUser : overrideUser;

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
