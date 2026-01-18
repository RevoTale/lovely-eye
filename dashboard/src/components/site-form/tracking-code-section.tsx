import React from 'react';
import { useMutation } from '@apollo/client/react';
import { RegenerateSiteKeyDocument } from '@/gql/graphql';
import { TrackingCodeCard } from '@/components/site-form/tracking-code-card';

interface TrackingCodeSectionProps {
  siteId: string;
  publicKey: string;
  onViewAnalytics: () => void;
}

export function TrackingCodeSection({
  siteId,
  publicKey,
  onViewAnalytics,
}: TrackingCodeSectionProps): React.JSX.Element {
  const [actionError, setActionError] = React.useState('');
  const [regenerateKey, { loading: regenerating }] = useMutation(RegenerateSiteKeyDocument);

  const handleRegenerateKey = async (): Promise<void> => {
    if (!window.confirm('Are you sure you want to regenerate the site key? The old key will stop working.')) {
      return;
    }

    setActionError('');
    try {
      await regenerateKey({
        variables: { id: siteId },
      });
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to regenerate the site key');
    }
  };

  const trackingScript = React.useMemo(() => {
    const basePath = window.__ENV__?.BASE_PATH ?? '';
    const trackerUrl = `${window.location.origin}${basePath}/tracker.js`;

    return `<script
  defer
  src="${trackerUrl}"
  data-site-key="${publicKey}"
></script>`;
  }, [publicKey]);

  return (
    <div className="space-y-2">
      <TrackingCodeCard
        publicKey={publicKey}
        trackingScript={trackingScript}
        regenerating={regenerating}
        onRegenerateKey={() => {
          void handleRegenerateKey();
        }}
        onViewAnalytics={onViewAnalytics}
      />
      {actionError ? (
        <p className="text-xs text-destructive">
          {actionError}
        </p>
      ) : null}
    </div>
  );
}
