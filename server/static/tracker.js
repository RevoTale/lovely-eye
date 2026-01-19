/**
 * Lovely Eye Analytics Tracker
 * Lightweight, privacy-focused analytics tracking script
 */
(function() {
  'use strict';

  var script = document.currentScript;
  var siteKey = script.getAttribute('data-site-key');
  var apiUrl = script.getAttribute('data-api-url') || (script.src.replace(/\/[^/]*$/, ''));

  if (!siteKey) {
    console.warn('Lovely Eye: Missing data-site-key attribute');
    return;
  }

  var lastPath = null;
  var pageStartTime = Date.now();
  var includeQuery = script.getAttribute('data-include-query') === 'true';

  function getPath() {
    return includeQuery ? window.location.pathname + window.location.search : window.location.pathname;
  }

  function getReferrer() {
    var ref = document.referrer;
    if (!ref) return '';
    try {
      var refUrl = new URL(ref);
      if (refUrl.hostname === window.location.hostname) {
        return ''; // Internal referrer, treat as direct
      }
      return ref;
    } catch (e) {
      return ref;
    }
  }

  function getUTMParams() {
    var params = new URLSearchParams(window.location.search);
    return {
      utm_source: params.get('utm_source') || '',
      utm_medium: params.get('utm_medium') || '',
      utm_campaign: params.get('utm_campaign') || ''
    };
  }

  function send(endpoint, data) {
    var url = apiUrl + endpoint + '?site_key=' + encodeURIComponent(siteKey);

    // Use sendBeacon if available for reliability
    if (navigator.sendBeacon) {
      var blob = new Blob([JSON.stringify(data)], { type: 'application/json' });
      navigator.sendBeacon(url, blob);
    } else {
      // Fallback to fetch
      fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
        keepalive: true
      }).catch(function() {});
    }
  }

  function trackPageView() {
    var path = getPath();

    if (path === lastPath) return;
    lastPath = path;

    var utm = getUTMParams();

    send('/api/collect', {
      site_key: siteKey,
      path: path,
      title: document.title,
      referrer: getReferrer(),
      screen_width: window.innerWidth,
      utm_source: utm.utm_source,
      utm_medium: utm.utm_medium,
      utm_campaign: utm.utm_campaign
    });

    pageStartTime = Date.now();
  }

  function trackEvent(name, properties) {
    if (!name) return;

    send('/api/event', {
      site_key: siteKey,
      name: name,
      path: getPath(),
      properties: properties ? JSON.stringify(properties) : ''
    });
  }

  function trackLeave() {
    var duration = Math.round((Date.now() - pageStartTime) / 1000);
    if (duration > 0 && duration < 3600) { // Sanity check: max 1 hour
      send('/api/collect', {
        site_key: siteKey,
        path: getPath(),
        duration: duration
      });
    }
  }

  function init() {
    trackPageView();

    document.addEventListener('visibilitychange', function() {
      if (document.visibilityState === 'hidden') {
        trackLeave();
      }
    });

    var originalPushState = history.pushState;
    history.pushState = function() {
      originalPushState.apply(this, arguments);
      trackPageView();
    };

    var originalReplaceState = history.replaceState;
    history.replaceState = function() {
      originalReplaceState.apply(this, arguments);
      trackPageView();
    };

    window.addEventListener('popstate', trackPageView);

    window.addEventListener('beforeunload', trackLeave);
  }

  window.lovelyEye = {
    track: trackEvent,
    trackPageView: trackPageView
  };

  if (document.readyState === 'complete') {
    init();
  } else {
    window.addEventListener('load', init);
  }
})();
