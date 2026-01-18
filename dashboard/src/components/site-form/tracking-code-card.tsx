import React from 'react';
import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Input, Label } from '@/components/ui';
import { CheckCircle2, Copy, RefreshCw } from 'lucide-react';

interface TrackingCodeCardProps {
  publicKey: string;
  trackingScript: string;
  regenerating: boolean;
  onRegenerateKey: () => void;
  onViewAnalytics: () => void;
}

export function TrackingCodeCard({
  publicKey,
  trackingScript,
  regenerating,
  onRegenerateKey,
  onViewAnalytics,
}: TrackingCodeCardProps): React.JSX.Element {
  const [copied, setCopied] = React.useState(false);

  const handleCopy = async (value: string): Promise<void> => {
    await navigator.clipboard.writeText(value);
    setCopied(true);
    setTimeout(() => {
      setCopied(false);
    }, 2000);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Tracking Code</CardTitle>
        <CardDescription>
          Add this script to your website&apos;s HTML to start tracking
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label>Site Key</Label>
          <div className="flex gap-2">
            <Input
              value={publicKey}
              readOnly
              className="font-mono text-sm"
            />
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => {
                void handleCopy(publicKey);
              }}
            >
              {copied ? (
                <CheckCircle2 className="h-4 w-4 text-green-500" />
              ) : (
                <Copy className="h-4 w-4" />
              )}
            </Button>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={onRegenerateKey}
              disabled={regenerating}
            >
              <RefreshCw className={`h-4 w-4 ${regenerating ? 'animate-spin' : ''}`} />
            </Button>
          </div>
          <p className="text-xs text-muted-foreground">
            Used by the tracker to associate events with this site.
          </p>
        </div>

        <div className="space-y-2">
          <Label>Tracking Script</Label>
          <div className="relative">
            <pre className="p-4 bg-muted rounded-lg overflow-x-auto text-xs border">
              <code>{trackingScript}</code>
            </pre>
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="absolute top-2 right-2"
              onClick={() => {
                void handleCopy(trackingScript);
              }}
            >
              {copied ? (
                <CheckCircle2 className="h-4 w-4 text-green-500" />
              ) : (
                <Copy className="h-4 w-4" />
              )}
            </Button>
          </div>
          <p className="text-xs text-muted-foreground">
            Add this script to the &lt;head&gt; section of your website
          </p>
        </div>

        <div className="pt-4">
          <Button
            variant="outline"
            onClick={onViewAnalytics}
          >
            View Analytics
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
