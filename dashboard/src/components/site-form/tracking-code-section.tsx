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

  const { trackingScript, trackingSnippet } = React.useMemo(() => {
    const basePath = window.__ENV__?.BASE_PATH ?? '';
    const trackerUrl = `${window.location.origin}${basePath}/tracker.js`;

    const scriptTag = `<script
  defer
  src="${trackerUrl}"
  data-site-id="${publicKey}"
></script>`;
    const scriptSnippet = `(function () {
  var script = document.createElement('script');
  script.defer = true;
  script.src = '${trackerUrl}';
  script.dataset.siteKey = '${publicKey}';
  document.head.appendChild(script);
})();`;

    return { trackingScript: scriptTag, trackingSnippet: scriptSnippet };
  }, [publicKey]);

  return (
    <div className="space-y-2">
      <TrackingCodeCard
        publicKey={publicKey}
        trackingScript={trackingScript}
        trackingSnippet={trackingSnippet}
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
