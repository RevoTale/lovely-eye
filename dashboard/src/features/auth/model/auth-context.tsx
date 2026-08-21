import { useApolloClient, useMutation, useQuery } from '@apollo/client/react';
import {
  createContext,
  type ReactNode,
  type RefObject,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from 'react';
import {
  type AuthUserDetailsFieldsFragment,
  AuthUserDetailsFieldsFragmentDoc,
  LoginDocument,
  type LoginInput,
  LogoutDocument,
  MeDocument,
  RegisterDocument,
  type RegisterInput,
} from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';

type AuthUser = AuthUserDetailsFieldsFragment;
type AuthMode = 'register-only' | 'login-only' | 'login-and-register';

const AUTH_STATUS_ERROR_MESSAGE =
  'Unable to load authentication status. Refresh the page and try again.';
const AUTH_RESPONSE_ERROR_MESSAGE = 'Authentication completed without a user.';

export interface AuthContextValue {
  user: AuthUser | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  authMode: AuthMode;
  bootstrapError: string | null;
  unauthenticatedRoute: '/login' | '/register';
  login: (input: LoginInput) => Promise<void>;
  register: (input: RegisterInput) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

interface AuthProviderProps {
  children: ReactNode;
  authErrorHandlerRef?: RefObject<(() => void) | null>;
}

export const AuthProvider = ({
  children,
  authErrorHandlerRef,
}: AuthProviderProps): React.ReactNode => {
  const authErrorHandledRef = useRef(false);
  const [userOverride, setUserOverride] = useState<AuthUser | null | undefined>(undefined);
  const [retainedAuthMode, setRetainedAuthMode] = useState<AuthMode>('login-only');
  const client = useApolloClient();

  const {
    loading: meLoading,
    data: meData,
    error: meError,
    refetch,
  } = useQuery(MeDocument, {
    fetchPolicy: 'network-only',
    errorPolicy: 'all',
  });
  const registrationStatus = meData?.registrationStatus;

  const [loginMutation] = useMutation(LoginDocument);
  const [registerMutation] = useMutation(RegisterDocument);
  const [logoutMutation] = useMutation(LogoutDocument);

  const handleAuthError = useCallback(() => {
    if (authErrorHandledRef.current) {
      return;
    }
    authErrorHandledRef.current = true;
    setUserOverride(undefined);
    void client.clearStore();
    void refetch();
  }, [client, refetch]);

  useEffect(() => {
    if (authErrorHandlerRef !== undefined) {
      const nextRef = authErrorHandlerRef;
      nextRef.current = handleAuthError;
    }
  }, [authErrorHandlerRef, handleAuthError]);

  const login = useCallback(
    async (input: LoginInput) => {
      const result = await loginMutation({ variables: { input } });
      const userData = result.data?.login.user;
      if (userData === null || userData === undefined) throw new Error(AUTH_RESPONSE_ERROR_MESSAGE);
      const nextUser = readFragment(AuthUserDetailsFieldsFragmentDoc, userData);
      await client.clearStore();
      setUserOverride(nextUser);
      authErrorHandledRef.current = false;
    },
    [client, loginMutation]
  );

  const register = useCallback(
    async (input: RegisterInput) => {
      const result = await registerMutation({ variables: { input } });
      const userData = result.data?.register.user;
      if (userData === null || userData === undefined) throw new Error(AUTH_RESPONSE_ERROR_MESSAGE);
      const nextUser = readFragment(AuthUserDetailsFieldsFragmentDoc, userData);
      await client.clearStore();
      setUserOverride(nextUser);
      if (registrationStatus !== null && registrationStatus !== undefined) {
        setRetainedAuthMode(getAuthMode(true, registrationStatus.allowRegistration));
      }
      authErrorHandledRef.current = false;
    },
    [client, registerMutation, registrationStatus]
  );

  const logout = useCallback(async () => {
    await logoutMutation();
    await client.clearStore();
    setUserOverride(null);
  }, [client, logoutMutation]);

  const userData = meData?.me;
  const queriedUser =
    userData !== null && userData !== undefined
      ? readFragment(AuthUserDetailsFieldsFragmentDoc, userData)
      : null;
  const queriedAuthMode =
    registrationStatus !== null && registrationStatus !== undefined
      ? getAuthMode(registrationStatus.hasUsers, registrationStatus.allowRegistration)
      : undefined;
  useEffect(() => {
    if (queriedAuthMode !== undefined) setRetainedAuthMode(queriedAuthMode);
  }, [queriedAuthMode]);
  const authMode = queriedAuthMode ?? retainedAuthMode;
  const user = userOverride === undefined ? queriedUser : userOverride;
  const bootstrapError =
    userOverride === undefined && !meLoading && registrationStatus == null && meError !== undefined
      ? AUTH_STATUS_ERROR_MESSAGE
      : null;
  const unauthenticatedRoute = getUnauthenticatedRoute(authMode);

  const value: AuthContextValue = {
    user,
    isLoading: userOverride === undefined && meLoading,
    isAuthenticated: user !== null,
    authMode,
    bootstrapError,
    unauthenticatedRoute,
    login,
    register,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (context === null) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

function getAuthMode(hasUsers: boolean, allowRegistration: boolean): AuthMode {
  if (!hasUsers) {
    return 'register-only';
  }
  if (!allowRegistration) {
    return 'login-only';
  }
  return 'login-and-register';
}

function getUnauthenticatedRoute(authMode: AuthMode): '/login' | '/register' {
  if (authMode === 'register-only') {
    return '/register';
  }
  return '/login';
}
