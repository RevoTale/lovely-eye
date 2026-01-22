
import { useState, type ReactElement } from 'react';
import { useMutation } from '@apollo/client/react';
import { DeleteSiteDocument, SitesDocument } from '@/gql/graphql';
import { DangerZoneCard } from '@/components/site-form/danger-zone-card';

interface DangerZoneSectionProps {
  siteId: string;
  onDeleted: () => void;
}

export function DangerZoneSection({
  siteId,
  onDeleted,
}: DangerZoneSectionProps): ReactElement {
  const SITES_PAGE_SIZE = 100;
  const SITES_PAGE_OFFSET = 0;
  const sitesPaging = { limit: SITES_PAGE_SIZE, offset: SITES_PAGE_OFFSET };
  const [actionError, setActionError] = useState('');
  const [confirmingDelete, setConfirmingDelete] = useState(false);
  const [deleteSite, { loading: deleting }] = useMutation(DeleteSiteDocument, {
    refetchQueries: [{ query: SitesDocument, variables: { paging: sitesPaging } }],
    onCompleted: () => {
      onDeleted();
    },
  });

  const handleDelete = async (): Promise<void> => {
    setActionError('');
    try {
      await deleteSite({
        variables: { id: siteId },
      });
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to delete site');
    } finally {
      setConfirmingDelete(false);
    }
  };

  return (
    <div className="space-y-2">
      <DangerZoneCard
        deleting={deleting}
        confirming={confirmingDelete}
        onDeleteRequest={() => {
          setConfirmingDelete(true);
        }}
        onConfirmDelete={() => {
          void handleDelete();
        }}
        onCancelDelete={() => {
          setConfirmingDelete(false);
        }}
      />
      {actionError === '' ? null : (
        <p className="text-xs text-destructive">
          {actionError}
        </p>
      )}
    </div>
  );
}
