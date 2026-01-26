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

declare global {
  interface Window {
    lovelyEye?: {
      track: (data?: TrackInput) => void;
    };
  }
}

(() => {
  'use strict';

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
      utm_campaign: params.get('utm_campaign') || ''
    };

    if (!data) return payload;

    const {
      name,
      path,
      referrer,
      screen_width,
      duration,
      properties,
      utm_source,
      utm_medium,
      utm_campaign
    } = data ?? {};

    if (typeof name === 'string') payload.name = name;
    if (typeof path === 'string') payload.path = path;
    if (typeof referrer === 'string') payload.referrer = referrer;
    if (typeof screen_width === 'number') payload.screen_width = screen_width;
    if (typeof duration === 'number') payload.duration = duration;
    if (typeof utm_source === 'string') payload.utm_source = utm_source;
    if (typeof utm_medium === 'string') payload.utm_medium = utm_medium;
    if (typeof utm_campaign === 'string') payload.utm_campaign = utm_campaign;
    if (typeof properties === 'string') payload.properties = properties;
    else if (properties) payload.properties = JSON.stringify(properties);

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
        keepalive: true
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
    history.pushState = function () {
      originalPushState.apply(this, arguments as unknown as [data: unknown, title: string, url?: string | URL | null]);
      track();
    };

    const originalReplaceState = history.replaceState;
    history.replaceState = function () {
      originalReplaceState.apply(this, arguments as unknown as [data: unknown, title: string, url?: string | URL | null]);
      track();
    };

    window.addEventListener('popstate', track);
    window.addEventListener('beforeunload', trackLeave);
  };

  window.lovelyEye = { track };

  if (document.readyState === 'complete') init();
  else window.addEventListener('load', init);
})();

export {};
