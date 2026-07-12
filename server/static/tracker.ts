/**
 * Lovely Eye Analytics Tracker
 * Lightweight, privacy-focused analytics tracking script
 */
type TrackInput = {
  name?: string;
  path?: string;
  referrer?: string;
  properties?: Record<string, unknown> | string;
  utm_source?: string;
  utm_medium?: string;
  utm_campaign?: string;
};

type TrackPayload = {
  name?: string;
  path: string;
  properties?: string;
  referrer?: string;
  exit?: true;
  utm_source?: string;
  utm_medium?: string;
  utm_campaign?: string;
};

type PayloadStringKey = 'name' | 'path' | 'referrer' | 'utm_source' | 'utm_medium' | 'utm_campaign';

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
  let exitSent = false;

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

  const getPropertiesValue = (properties: TrackInput['properties']): string | undefined => {
    if (typeof properties === 'string') {
      return properties;
    }

    if (properties !== undefined) {
      return JSON.stringify(properties);
    }

    return undefined;
  };

  const assignAttribution = (payload: TrackPayload): void => {
    const params = new URLSearchParams(window.location.search);
    const referrer = getReferrer();
    if (referrer) payload.referrer = referrer;

    const utmSource = params.get('utm_source');
    const utmMedium = params.get('utm_medium');
    const utmCampaign = params.get('utm_campaign');
    if (utmSource) payload.utm_source = utmSource;
    if (utmMedium) payload.utm_medium = utmMedium;
    if (utmCampaign) payload.utm_campaign = utmCampaign;
  };

  const buildPayload = (data?: TrackInput, includeAttribution = false): TrackPayload => {
    const payload: TrackPayload = { path: getPath() };

    if (includeAttribution) {
      assignAttribution(payload);
    }

    if (!data) return payload;

    assignStringOverride(payload, 'name', data.name);
    assignStringOverride(payload, 'path', data.path);
    assignStringOverride(payload, 'referrer', data.referrer);
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
      const blob = new Blob([payload], { type: 'text/plain;charset=UTF-8' });
      navigator.sendBeacon(url, blob);
    } else {
      fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'text/plain;charset=UTF-8' },
        body: payload,
        keepalive: true,
      }).catch(() => {});
    }
  };

  const track = (data?: TrackInput): void => {
    const payload = buildPayload(data, lastPath === '' && !data?.name);
    if (payload.path === lastPath && !payload.name) return;
    lastPath = payload.path;
    exitSent = false;
    send('/api/collect', payload);
  };

  const trackExit = (): void => {
    if (exitSent) return;
    const path = getPath();
    if (!path) return;
    exitSent = true;
    send('/api/collect', { path, exit: true });
  };

  const init = (): void => {
    track();

    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'hidden') {
        trackExit();
      } else {
        exitSent = false;
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
    window.addEventListener('pagehide', trackExit);
  };

  window.lovelyEye = { track };

  if (document.readyState === 'complete') init();
  else window.addEventListener('load', init);
})();

export {};
