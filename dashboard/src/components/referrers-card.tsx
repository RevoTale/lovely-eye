import React from 'react';
import { Badge, Progress } from '@/components/ui';
import { Globe, ExternalLink, TrendingUp } from 'lucide-react';
import type { ReferrerStats } from '@/gql/graphql';
import { BoardCard } from '@/components/board-card';
import { FilterLink } from '@/components/filter-link';
import { ListEmptyState } from '@/components/list-empty-state';

interface ReferrersCardProps {
  referrers: ReferrerStats[];
  totalCount: number;
  totalVisitors: number;
  siteId: string;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

export function ReferrersCard({
  referrers,
  totalCount,
  totalVisitors,
  siteId,
  page,
  pageSize,
  onPageChange,
}: ReferrersCardProps): React.JSX.Element {
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

  const headerRight = (
    <Badge variant="secondary" className="flex items-center gap-1">
      <TrendingUp className="h-3 w-3" />
      {totalCount}
    </Badge>
  );

  return (
    <BoardCard
      title="Top Referrers"
      icon={Globe}
      headerRight={headerRight}
      pagination={{ page, pageSize, total: totalCount, onPageChange }}
    >
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
                    <FilterLink
                      siteId={siteId}
                      filterKey="referrer"
                      value={ref.referrer || '(direct)'}
                      className="text-sm font-medium truncate hover:text-primary hover:underline cursor-pointer"
                    >
                      {formatReferrer(ref.referrer)}
                    </FilterLink>
                    {ref.referrer && ref.referrer !== '' ? (
                      <a
                        href={ref.referrer}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-muted-foreground hover:text-primary transition-colors"
                        onClick={(e) => {
                          e.stopPropagation();
                        }}
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
          <ListEmptyState
            icon={Globe}
            title="No referrer data yet"
            description="Referrers will appear as visitors arrive from external sources"
          />
        )}
      </div>
    </BoardCard>
  );
}
