export type GraphQLErrorCode =
  | 'BAD_USER_INPUT'
  | 'UNAUTHENTICATED'
  | 'FORBIDDEN'
  | 'NOT_FOUND'
  | 'CONFLICT'
  | 'INTERNAL_SERVER_ERROR';

export function getGraphQLErrorCode(extensions: unknown): GraphQLErrorCode | null {
  if (typeof extensions !== 'object' || extensions === null || !('code' in extensions)) {
    return null;
  }

  switch (extensions.code) {
    case 'BAD_USER_INPUT':
    case 'UNAUTHENTICATED':
    case 'FORBIDDEN':
    case 'NOT_FOUND':
    case 'CONFLICT':
    case 'INTERNAL_SERVER_ERROR':
      return extensions.code;
    default:
      return null;
  }
}
