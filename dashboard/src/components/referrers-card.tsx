import React from 'react';
import { Card, CardContent, CardHeader, CardTitle, Badge, Progress } from '@/components/ui';
import { Globe, ExternalLink, TrendingUp } from 'lucide-react';
import type { ReferrerStats } from '@/generated/graphql';
import { Link } from '@/router';

interface ReferrersCardProps {
  referrers: ReferrerStats[];
  totalVisitors: number;
  siteId: string;
}

export function ReferrersCard({ referrers, totalVisitors, siteId }: ReferrersCardProps): React.JSX.Element {
  const formatReferrer = (referrer: string | null): string => {
    if (!referrer) return 'Direct / None';

    try {
      const url = new URL(referrer);
      return url.hostname.replace('www.', '');
    } catch {
      return referrer;
    }
  };

  const getReferrerIcon = (referrer: string | null): string => {
    if (!referrer) return '🔗';

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

  const maxVisitors = referrers.length > 0 && referrers[0] ? referrers[0].visitors : 1;

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
              <Globe className="h-4 w-4 text-primary" />
            </div>
            Top Referrers
          </div>
          <Badge variant="secondary" className="flex items-center gap-1">
            <TrendingUp className="h-3 w-3" />
            {referrers.length}
          </Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {referrers.length > 0 ? (
            referrers.map((ref, index) => {
              const percentage = totalVisitors > 0
                ? ((ref.visitors / totalVisitors) * 100).toFixed(1)
                : '0.0';
              const barWidth = maxVisitors > 0
                ? (ref.visitors / maxVisitors) * 100
                : 0;

              return (
                <div key={index} className="space-y-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2 flex-1 min-w-0">
                      <span className="text-lg" role="img" aria-label="referrer icon">
                        {getReferrerIcon(ref.referrer)}
                      </span>
                      <Link
                        to="/sites/$siteId"
                        params={{ siteId }}
                        search={{ referrer: ref.referrer || '(direct)' }}
                        className="text-sm font-medium truncate hover:text-primary hover:underline cursor-pointer"
                      >
                        {formatReferrer(ref.referrer)}
                      </Link>
                      {ref.referrer && ref.referrer !== '' ? (
                        <a
                          href={ref.referrer}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-muted-foreground hover:text-primary transition-colors"
                          onClick={(e) => e.stopPropagation()}
                        >
                          <ExternalLink className="h-3 w-3" />
                        </a>
                      ) : null}
                    </div>
                    <div className="flex items-center gap-3 ml-2">
                      <Badge variant="outline" className="font-mono">
                        {ref.visitors.toLocaleString()}
                      </Badge>
                      <span className="text-xs text-muted-foreground w-12 text-right">
                        {percentage}%
                      </span>
                    </div>
                  </div>
                  <Progress value={barWidth} className="h-2" />
                </div>
              );
            })
          ) : (
            <div className="text-center py-8">
              <Globe className="h-12 w-12 mx-auto mb-3 text-muted-foreground opacity-50" />
              <p className="text-sm text-muted-foreground">No referrer data yet</p>
              <p className="text-xs text-muted-foreground mt-1">
                Referrers will appear as visitors arrive from external sources
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
