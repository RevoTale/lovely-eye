/**
 * Lovely Eye Analytics Tracker
 * Lightweight, privacy-focused analytics tracking script
 */
type TrackInput = {
  name?: string;
  path?: string;
  referrer?: string;
  screen_width?: number;
  duration?: number;
  properties?: Record<string, unknown> | string;
  utm_source?: string;
  utm_medium?: string;
  utm_campaign?: string;
};

type TrackPayload = {
  site_key: string;
  name: string;
  path: string;
  properties: string;
  referrer: string;
  screen_width: number;
  duration: number;
  utm_source: string;
  utm_medium: string;
  utm_campaign: string;
};

type PayloadStringKey = 'name' | 'path' | 'referrer' | 'utm_source' | 'utm_medium' | 'utm_campaign';

type PayloadNumberKey = 'screen_width' | 'duration';

declare global {
  interface Window {
    lovelyEye?: {
      track: (data?: TrackInput) => void;
    };
  }
}

(() => {
  const script = document.currentScript as HTMLScriptElement | null;
  const siteKey = script?.getAttribute('data-site-key') ?? '';
  const apiUrl = script?.getAttribute('data-api-url') ?? script?.src?.replace(/\/[^/]*$/, '') ?? '';
  const includeQuery = script?.getAttribute('data-include-query') === 'true';

  if (!siteKey || !apiUrl) return;

  let lastPath = '';
  let pageStartTime = Date.now();

  const getPath = (): string =>
    includeQuery ? window.location.pathname + window.location.search : window.location.pathname;

  const getReferrer = (): string => {
    const ref = document.referrer;
    if (!ref) return '';
    try {
      const refUrl = new URL(ref);
      if (refUrl.hostname === window.location.hostname) return '';
      return ref;
    } catch {
      return ref;
    }
  };

  const assignStringOverride = (
    payload: TrackPayload,
    key: PayloadStringKey,
    value: string | undefined
  ): void => {
    if (typeof value === 'string') {
      payload[key] = value;
    }
  };

  const assignNumberOverride = (
    payload: TrackPayload,
    key: PayloadNumberKey,
    value: number | undefined
  ): void => {
    if (typeof value === 'number') {
      payload[key] = value;
    }
  };

  const getPropertiesValue = (properties: TrackInput['properties']): string | undefined => {
    if (typeof properties === 'string') {
      return properties;
    }

    if (properties !== undefined) {
      return JSON.stringify(properties);
    }

    return undefined;
  };

  const buildPayload = (data?: TrackInput): TrackPayload => {
    const params = new URLSearchParams(window.location.search);
    const payload: TrackPayload = {
      site_key: siteKey,
      name: '',
      path: getPath(),
      properties: '',
      referrer: getReferrer(),
      screen_width: window.innerWidth,
      duration: 0,
      utm_source: params.get('utm_source') || '',
      utm_medium: params.get('utm_medium') || '',
      utm_campaign: params.get('utm_campaign') || '',
    };

    if (!data) return payload;

    assignStringOverride(payload, 'name', data.name);
    assignStringOverride(payload, 'path', data.path);
    assignStringOverride(payload, 'referrer', data.referrer);
    assignNumberOverride(payload, 'screen_width', data.screen_width);
    assignNumberOverride(payload, 'duration', data.duration);
    assignStringOverride(payload, 'utm_source', data.utm_source);
    assignStringOverride(payload, 'utm_medium', data.utm_medium);
    assignStringOverride(payload, 'utm_campaign', data.utm_campaign);

    const properties = getPropertiesValue(data.properties);
    if (properties !== undefined) {
      payload.properties = properties;
    }

    return payload;
  };

  const send = (endpoint: string, data: TrackPayload): void => {
    const url = `${apiUrl}${endpoint}?site_key=${encodeURIComponent(siteKey)}`;
    const payload = JSON.stringify(data);

    if (navigator.sendBeacon) {
      const blob = new Blob([payload], { type: 'application/json' });
      navigator.sendBeacon(url, blob);
    } else {
      fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: payload,
        keepalive: true,
      }).catch(() => {});
    }
  };

  const track = (data?: TrackInput): void => {
    const payload = buildPayload(data);
    if (payload.path === lastPath && !payload.name && payload.duration === 0) return;
    lastPath = payload.path;
    send('/api/collect', payload);
    pageStartTime = Date.now();
  };

  const trackLeave = (): void => {
    const duration = Math.round((Date.now() - pageStartTime) / 1000);
    if (duration > 0 && duration < 3600) {
      track({ duration });
    }
  };

  const init = (): void => {
    track();

    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'hidden') {
        trackLeave();
      }
    });

    const originalPushState = history.pushState;
    history.pushState = function (this: History, ...args: Parameters<History['pushState']>) {
      originalPushState.apply(this, args);
      track();
    };

    const originalReplaceState = history.replaceState;
    history.replaceState = function (this: History, ...args: Parameters<History['replaceState']>) {
      originalReplaceState.apply(this, args);
      track();
    };

    window.addEventListener('popstate', () => {
      track();
    });
    window.addEventListener('beforeunload', trackLeave);
  };

  window.lovelyEye = { track };

  if (document.readyState === 'complete') init();
  else window.addEventListener('load', init);
})();

export {};
