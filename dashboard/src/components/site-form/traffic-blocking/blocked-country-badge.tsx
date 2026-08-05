import type { FunctionComponent } from 'react';
import { X } from 'lucide-react';
import { Badge, Button } from '@/components/ui';

interface BlockedCountryBadgeProps {
  code: string;
  countryNameLookup: Map<string, string>;
  savingBlockedCountries: boolean;
  onRemoveBlockedCountry: (code: string) => void;
}

const BlockedCountryBadge: FunctionComponent<BlockedCountryBadgeProps> = ({
  code,
  countryNameLookup,
  onRemoveBlockedCountry,
  savingBlockedCountries,
}) => {
  const normalizedCode = code.trim().toUpperCase();
  const displayName = countryNameLookup.get(normalizedCode) ?? normalizedCode;
  const showCode = displayName.trim().toUpperCase() !== normalizedCode;

  return (
    <Badge variant="secondary" className="flex items-center gap-2">
      <span>{displayName}</span>
      {showCode ? <span className="text-xs text-muted-foreground">{code}</span> : null}
      <Button
        type="button"
        variant="ghost"
        size="icon"
        className="h-5 w-5"
        disabled={savingBlockedCountries}
        onClick={() => onRemoveBlockedCountry(code)}
        aria-label={`Remove ${code}`}
      >
        <X className="h-3 w-3" />
      </Button>
    </Badge>
  );
};

export default BlockedCountryBadge;
