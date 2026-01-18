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
  const [confirmingRegenerate, setConfirmingRegenerate] = React.useState(false);
  const [regenerateKey, { loading: regenerating }] = useMutation(RegenerateSiteKeyDocument);

  const handleRegenerateKey = async (): Promise<void> => {
    setActionError('');
    try {
      await regenerateKey({
        variables: { id: siteId },
      });
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to regenerate the site key');
    } finally {
      setConfirmingRegenerate(false);
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
        confirmingRegenerate={confirmingRegenerate}
        onRegenerateRequest={() => {
          setConfirmingRegenerate(true);
        }}
        onConfirmRegenerate={() => {
          void handleRegenerateKey();
        }}
        onCancelRegenerate={() => {
          setConfirmingRegenerate(false);
        }}
        onViewAnalytics={onViewAnalytics}
      />
      {actionError === '' ? null : (
        <p className="text-xs text-destructive">
          {actionError}
        </p>
      )}
    </div>
  );
}
