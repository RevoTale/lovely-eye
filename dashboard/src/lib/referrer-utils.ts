const DIRECT_REFERRER_LABEL = 'Direct / None';

export const formatReferrer = (referrer: string | null): string => {
  if (referrer === null || referrer === '') {
    return DIRECT_REFERRER_LABEL;
  }

  try {
    const url = new URL(referrer);
    return url.hostname.replace('www.', '');
  } catch {
    return referrer;
  }
};

export const getReferrerIcon = (referrer: string | null): string => {
  if (referrer === null || referrer === '') return '🔗';

  const hostname = formatReferrer(referrer).toLowerCase();
  if (hostname.includes('google')) return '🔍';
  if (hostname.includes('facebook') || hostname.includes('fb.')) return '👥';
  if (hostname.includes('twitter') || hostname.includes('x.com')) return '🐦';
  if (hostname.includes('linkedin')) return '💼';
  if (hostname.includes('github')) return '💻';
  if (hostname.includes('youtube')) return '📹';
  if (hostname.includes('reddit')) return '🤖';
  return '🌐';
};
