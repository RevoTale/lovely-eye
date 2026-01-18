import React from 'react';
import { useMutation } from '@apollo/client/react';
import { DeleteSiteDocument, SitesDocument } from '@/gql/graphql';
import { DangerZoneCard } from '@/components/site-form/danger-zone-card';

interface DangerZoneSectionProps {
  siteId: string;
  siteName: string;
  onDeleted: () => void;
}

export function DangerZoneSection({
  siteId,
  siteName,
  onDeleted,
}: DangerZoneSectionProps): React.JSX.Element {
  const [actionError, setActionError] = React.useState('');
  const [deleteSite, { loading: deleting }] = useMutation(DeleteSiteDocument, {
    refetchQueries: [{ query: SitesDocument }],
    onCompleted: () => {
      onDeleted();
    },
  });

  const handleDelete = async (): Promise<void> => {
    if (!window.confirm(`Are you sure you want to delete "${siteName}"? This action cannot be undone.`)) {
      return;
    }

    setActionError('');
    try {
      await deleteSite({
        variables: { id: siteId },
      });
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to delete site');
    }
  };

  return (
    <div className="space-y-2">
      <DangerZoneCard
        deleting={deleting}
        onDelete={() => {
          void handleDelete();
        }}
      />
      {actionError ? (
        <p className="text-xs text-destructive">
          {actionError}
        </p>
      ) : null}
    </div>
  );
}
