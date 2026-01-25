import { createContext, useContext, useCallback, useRef, useEffect, type RefObject, type ReactNode } from 'react';
import { useQuery, useMutation, useApolloClient } from '@apollo/client/react';
import {
  MeDocument,
  LoginDocument,
  RegisterDocument,
  LogoutDocument,
  AuthUserDetailsFieldsFragmentDoc,
  type AuthUserDetailsFieldsFragment,
  type LoginInput,
  type RegisterInput,
} from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';

type AuthUser = AuthUserDetailsFieldsFragment;

const SITES_PAGE_SIZE = 100;
const SITES_PAGE_OFFSET = 0;

export interface AuthContextType {
  user: AuthUser | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (input: LoginInput) => Promise<void>;
  register: (input: RegisterInput) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

interface AuthProviderProps {
  children: ReactNode;
  authErrorHandlerRef?: RefObject<(() => void) | null>;
}

export const AuthProvider = ({ children, authErrorHandlerRef }: AuthProviderProps): React.ReactNode => {
  const authErrorHandledRef = useRef(false);
  const client = useApolloClient();

  const { loading: meLoading, data: meData, refetch } = useQuery(MeDocument, {
    variables: {
      sitesPaging: {
        limit: SITES_PAGE_SIZE,
        offset: SITES_PAGE_OFFSET,
      },
    },
    fetchPolicy: 'network-only',
    errorPolicy: 'ignore',
  });

  const [loginMutation] = useMutation(LoginDocument);
  const [registerMutation] = useMutation(RegisterDocument);
  const [logoutMutation] = useMutation(LogoutDocument);

  const handleAuthError = useCallback(() => {
    if (authErrorHandledRef.current) {
      return;
    }
    authErrorHandledRef.current = true;
    void client.clearStore();
    void refetch();
  }, [client, refetch]);

  useEffect(() => {
    if (authErrorHandlerRef !== undefined) {
      const nextRef = authErrorHandlerRef;
      nextRef.current = handleAuthError;
    }
  }, [authErrorHandlerRef, handleAuthError]);

  const login = useCallback(async (input: LoginInput) => {
    await loginMutation({ variables: { input } });
    authErrorHandledRef.current = false;
    await refetch();
  }, [loginMutation, refetch]);

  const register = useCallback(async (input: RegisterInput) => {
    await registerMutation({ variables: { input } });
    authErrorHandledRef.current = false;
    await refetch();
  }, [registerMutation, refetch]);

  const logout = useCallback(async () => {
    await logoutMutation();
    await refetch();
  }, [logoutMutation, refetch]);

  const userData = meData?.me;
  const user =
    userData !== null && userData !== undefined
      ? getFragmentData(AuthUserDetailsFieldsFragmentDoc, userData)
      : null;

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
  if (context === null) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
