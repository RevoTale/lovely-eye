import type { Page, Route } from '@playwright/test';

interface DeferredOperation {
  seen: Promise<void>;
  release: () => void;
}

interface PendingOperation extends DeferredOperation {
  markSeen: () => void;
  waitForRelease: Promise<void>;
}

export interface GraphQLOperationController {
  blockNext: (operationName: string) => DeferredOperation;
}

export const installGraphQLOperationController = async (
  page: Page
): Promise<GraphQLOperationController> => {
  const pending = new Map<string, PendingOperation[]>();

  await page.route('**/graphql', async (route) => {
    const operationName = readOperationName(route);
    const queue = operationName === undefined ? undefined : pending.get(operationName);
    const operation = queue?.shift();
    if (operation !== undefined) {
      operation.markSeen();
      await operation.waitForRelease;
    }
    await route.continue();
  });

  return {
    blockNext: (operationName) => {
      const operation = createPendingOperation();
      const queue = pending.get(operationName) ?? [];
      queue.push(operation);
      pending.set(operationName, queue);
      return operation;
    },
  };
};

const createPendingOperation = (): PendingOperation => {
  let markSeen = (): void => undefined;
  let release = (): void => undefined;
  const seen = new Promise<void>((resolve) => {
    markSeen = resolve;
  });
  const waitForRelease = new Promise<void>((resolve) => {
    release = resolve;
  });
  return { markSeen, release, seen, waitForRelease };
};

const readOperationName = (route: Route): string | undefined => {
  const rawBody = route.request().postData();
  if (rawBody === null) return undefined;
  try {
    const payload: unknown = JSON.parse(rawBody);
    if (
      typeof payload === 'object' &&
      payload !== null &&
      'operationName' in payload &&
      typeof payload.operationName === 'string'
    ) {
      return payload.operationName;
    }
  } catch {
    return undefined;
  }
  return undefined;
};
