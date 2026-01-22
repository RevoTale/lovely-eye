
import { Badge, Progress } from '@/components/ui';
import { Globe, ExternalLink, TrendingUp } from 'lucide-react';
import type { ReferrerStats } from '@/gql/graphql';
import { BoardCard, BoardCardSkeleton } from '@/components/board-card';
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
  loading?: boolean;
}

const EMPTY_COUNT = 0;
const FALLBACK_MAX_VISITORS = 1;
const PERCENT_MULTIPLIER = 100;
const PERCENT_PRECISION = 1;
const ZERO_PERCENT = '0.0';
const DIRECT_LABEL = '(direct)';
const DIRECT_REFERRER_LABEL = 'Direct / None';
const FIRST_INDEX = 0;

export function ReferrersCard({
  referrers,
  totalCount,
  totalVisitors,
  siteId,
  page,
  pageSize,
  onPageChange,
  loading = false,
}: ReferrersCardProps): React.JSX.Element {
  if (loading) {
    return <BoardCardSkeleton title="Top Referrers" icon={Globe} />;
  }

  const formatReferrer = (referrer: string | null): string => {
    if (referrer === null || referrer === '') return DIRECT_REFERRER_LABEL;

    try {
      const url = new URL(referrer);
      return url.hostname.replace('www.', '');
    } catch {
      return referrer;
    }
  };

  const getReferrerIcon = (referrer: string | null): string => {
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

  const hasReferrers = referrers.length > EMPTY_COUNT;
  const topReferrer = hasReferrers ? referrers[FIRST_INDEX] : undefined;
  const maxVisitors =
    topReferrer === undefined ? FALLBACK_MAX_VISITORS:topReferrer.visitors ;

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
        {hasReferrers ? (
          referrers.map((ref, index) => {
            const percentage =
              totalVisitors > EMPTY_COUNT
                ? ((ref.visitors / totalVisitors) * PERCENT_MULTIPLIER).toFixed(PERCENT_PRECISION)
                : ZERO_PERCENT;
            const barWidth =
              maxVisitors > EMPTY_COUNT
                ? (ref.visitors / maxVisitors) * PERCENT_MULTIPLIER
                : EMPTY_COUNT;

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
                      value={ref.referrer === '' ? DIRECT_LABEL : ref.referrer}
                      className="text-sm font-medium truncate hover:text-primary hover:underline cursor-pointer"
                    >
                      {formatReferrer(ref.referrer)}
                    </FilterLink>
                    {ref.referrer === '' ? null : (
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
                    )}
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
