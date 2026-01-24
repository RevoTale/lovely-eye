
import { useCallback, useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle, ChartContainer, ChartLegend, ChartLegendContent, ChartTooltip, ChartTooltipContent, Progress, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui';
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts';
import { TrendingUp } from 'lucide-react';
import { CHART_CONFIG, CHART_MARGIN, TICK_MARGIN } from '@/lib/chart-config';
import { ChartSkeleton } from '@/components/chart-skeleton';
import { useChartDataLoader } from '@/hooks/use-chart-data-loader';
import type { FilterInput } from '@/gql/graphql';

interface OverviewChartSectionProps {
  siteId: string;
  dateRange: { from: Date; to: Date } | null;
  filter: FilterInput | null;
  bucket: 'daily' | 'hourly';
  onBucketChange: (bucket: 'daily' | 'hourly') => void;
}

const EMPTY_COUNT = 0;
const DATE_PART_LENGTH = 10;
const DATE_SEPARATOR = '-';
const TIME_SEPARATOR = ':';
const DATE_TIME_SEPARATOR = 'T';
const SPACE_SEPARATOR = ' ';
const YEAR_LENGTH = 4;
const MONTH_LENGTH = 2;
const DAY_LENGTH = 2;
const HOURS_LENGTH = 2;
const MINUTES_LENGTH = 2;
const SECONDS_LENGTH = 2;
const SECONDS_THRESHOLD = 1_000_000_000_000;
const MS_IN_SECOND = 1000;
const PROGRESS_MIN = 0;
const UTC_SUFFIX = ' UTC';
const UTC_MARKER = 'Z';
const TZ_OFFSET_LENGTH = 5;
const TZ_OFFSET_WITH_COLON_LENGTH = 6;
const TZ_COLON_INDEX = 3;
const NOT_FOUND = -1;
const NON_EMPTY_LENGTH = 1;
const INDEX_START = 0;
const INDEX_STEP = 1;
const CHAR_CODE_ZERO = 48;
const CHAR_CODE_NINE = 57;

const buildTimeFormatter = (bucket: 'daily' | 'hourly'): Intl.DateTimeFormat => {
  if (bucket === 'hourly') {
    return new Intl.DateTimeFormat('en-US', {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      timeZone: 'UTC',
    });
  }
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    timeZone: 'UTC',
  });
};

const normalizeEpoch = (value: number): number => (value < SECONDS_THRESHOLD ? value * MS_IN_SECOND : value);

const isDigits = (value: string): boolean => {
  if (value.length < NON_EMPTY_LENGTH) return false;
  for (let index = INDEX_START; index < value.length; index += INDEX_STEP) {
    const charCode = value.charCodeAt(index);
    if (charCode < CHAR_CODE_ZERO || charCode > CHAR_CODE_NINE) {
      return false;
    }
  }
  return true;
};

const isDateOnly = (value: string): boolean => {
  if (value.length !== DATE_PART_LENGTH) return false;
  const [year, month, day] = value.split(DATE_SEPARATOR);
  return year?.length === YEAR_LENGTH
    && month?.length === MONTH_LENGTH
    && day?.length === DAY_LENGTH
    && isDigits(year)
    && isDigits(month)
    && isDigits(day);
};

const isTimePart = (value: string): boolean => {
  const [hours, minutes, secondsWithZone] = value.split(TIME_SEPARATOR);
  const seconds = secondsWithZone?.slice(INDEX_START, SECONDS_LENGTH);
  return hours?.length === HOURS_LENGTH
    && minutes?.length === MINUTES_LENGTH
    && seconds?.length === SECONDS_LENGTH
    && isDigits(hours)
    && isDigits(minutes)
    && isDigits(seconds);
};

const normalizeTimezone = (value: string): string => {
  if (value.endsWith(UTC_SUFFIX)) {
    return `${value.slice(INDEX_START, -UTC_SUFFIX.length)}${UTC_MARKER}`;
  }
  if (value.endsWith(UTC_MARKER)) {
    return value;
  }
  const signIndex = Math.max(value.lastIndexOf('+'), value.lastIndexOf('-'));
  if (signIndex === NOT_FOUND) return value;
  const offset = value.slice(signIndex);
  if (offset.length === TZ_OFFSET_LENGTH && !offset.includes(TIME_SEPARATOR)) {
    return `${value.slice(INDEX_START, signIndex)}${offset.slice(INDEX_START, TZ_COLON_INDEX)}:${offset.slice(TZ_COLON_INDEX)}`;
  }
  if (offset.length === TZ_OFFSET_WITH_COLON_LENGTH) {
    return value;
  }
  return value;
};

const normalizeDateTime = (value: string): string | null => {
  const separator = value.includes(DATE_TIME_SEPARATOR) ? DATE_TIME_SEPARATOR : SPACE_SEPARATOR;
  const [datePart, timePart] = value.split(separator);
  if (datePart === undefined || datePart === '' || timePart === undefined || timePart === '') {
    return null;
  }
  if (!isDateOnly(datePart) || !isTimePart(timePart)) {
    return null;
  }
  return normalizeTimezone(`${datePart}${DATE_TIME_SEPARATOR}${timePart}`);
};

const parseChartDate = (value: string): number | null => {
  const trimmed = value.trim();
  const numeric = Number(trimmed);
  if (!Number.isNaN(numeric)) {
    return normalizeEpoch(numeric);
  }
  if (isDateOnly(trimmed)) {
    const parsed = Date.parse(`${trimmed}T00:00:00Z`);
    return Number.isNaN(parsed) ? null : parsed;
  }
  const normalized = normalizeDateTime(trimmed);
  if (normalized !== null) {
    const parsed = Date.parse(normalized);
    return Number.isNaN(parsed) ? null : parsed;
  }
  const parsed = Date.parse(trimmed);
  return Number.isNaN(parsed) ? null : parsed;
};

export function OverviewChartSection({ siteId, dateRange, filter, bucket, onBucketChange }: OverviewChartSectionProps): React.JSX.Element | null {
  const { loadedData, loading, loadingMore, progress, expectedCount } = useChartDataLoader({ siteId, dateRange, filter, bucket });

  const formatters = useMemo(() => ({
    daily: buildTimeFormatter('daily'),
    hourly: buildTimeFormatter('hourly'),
  }), []);

  const formatTimestamp = useCallback((timestamp: number): string => {
    const formatter = bucket === 'hourly' ? formatters.hourly : formatters.daily;
    return formatter.format(timestamp);
  }, [bucket, formatters]);

  const chartData = useMemo(() => (
    loadedData
      .map((stat) => {
        const timestamp = parseChartDate(stat.date);
        if (timestamp === null) {
          return null;
        }
        return {
          timestamp,
          visitors: stat.visitors,
          pageViews: stat.pageViews,
          sessions: stat.sessions,
        };
      })
      .filter((item): item is NonNullable<typeof item> => item !== null)
      .sort((a, b) => a.timestamp - b.timestamp)
  ), [loadedData]);

  if (loading && loadedData.length === EMPTY_COUNT) {
    return <ChartSkeleton />;
  }

  if (chartData.length === EMPTY_COUNT && !loading) {
    return null;
  }

  const showProgress = (loading && loadedData.length > EMPTY_COUNT) || loadingMore;
  const { length: loadedCount } = loadedData;
  const progressLabel = expectedCount === null
    ? `Loaded ${loadedCount.toLocaleString()} points`
    : `Loaded ${Math.min(loadedCount, expectedCount).toLocaleString()} of ${expectedCount.toLocaleString()} points`;

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="flex flex-col gap-4">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <CardTitle className="flex items-center gap-2">
            <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
              <TrendingUp className="h-4 w-4 text-primary" />
            </div>
            Analytics Overview
          </CardTitle>
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">Granularity</span>
            <Select value={bucket} onValueChange={(value) => {
              if (value === 'daily' || value === 'hourly') {
                onBucketChange(value);
              }
            }}>
              <SelectTrigger className="h-8 w-[140px]">
                <SelectValue placeholder="Daily" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="daily">Daily</SelectItem>
                <SelectItem value="hourly">Hourly</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
        {showProgress ? (
          <div className="space-y-2">
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              <span>{loadingMore ? 'Loading more data' : 'Loading chart data'}</span>
              <span>{progressLabel}</span>
            </div>
            <Progress value={progress ?? PROGRESS_MIN} className="h-2" />
          </div>
        ) : null}
      </CardHeader>
      <CardContent>
        <ChartContainer config={CHART_CONFIG} className="h-[300px] w-full">
          <AreaChart
            data={chartData}
            margin={CHART_MARGIN}
          >
            <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
            <XAxis
              dataKey="timestamp"
              type="number"
              scale="time"
              tickFormatter={(value) => formatTimestamp(Number(value))}
              tickLine={false}
              axisLine={false}
              tickMargin={TICK_MARGIN}
              className="text-xs"
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={TICK_MARGIN}
              className="text-xs"
            />
            <ChartTooltip content={
              <ChartTooltipContent
                labelFormatter={(value) => {
                  if (typeof value === 'number') {
                    return formatTimestamp(value);
                  }
                  if (typeof value === 'string') {
                    const timestamp = parseChartDate(value);
                    return timestamp === null ? value : formatTimestamp(timestamp);
                  }
                  return String(value);
                }}
              />
            } />
            <ChartLegend content={<ChartLegendContent />} />
            <Area
              type="monotone"
              dataKey="visitors"
              stackId="1"
              stroke="var(--color-visitors)"
              fill="var(--color-visitors)"
              fillOpacity={0.6}
            />
            <Area
              type="monotone"
              dataKey="pageViews"
              stackId="2"
              stroke="var(--color-pageViews)"
              fill="var(--color-pageViews)"
              fillOpacity={0.6}
            />
            <Area
              type="monotone"
              dataKey="sessions"
              stackId="3"
              stroke="var(--color-sessions)"
              fill="var(--color-sessions)"
              fillOpacity={0.6}
            />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
