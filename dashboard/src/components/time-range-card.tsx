import React, { useEffect, useMemo, useState } from 'react';
import { Button, Card, CardContent, CardHeader, CardTitle, Tabs, TabsList, TabsTrigger, Label } from '@/components/ui';
import type { DatePreset } from '@/lib/date-range';
import { formatDateInput, isValidDateInput, isValidTimeInput, normalizeTimeInput } from '@/lib/date-range';
import { DateTimePicker } from '@/components/ui/datetime-picker';

interface TimeRangeCardProps {
  preset: DatePreset;
  fromDate: string;
  toDate: string;
  fromTime: string;
  toTime: string;
  onPresetChange: (preset: DatePreset) => void;
  onApplyRange: (range: { fromDate: string; toDate: string; fromTime: string; toTime: string }) => boolean;
}

export function TimeRangeCard({
  preset,
  fromDate,
  toDate,
  fromTime,
  toTime,
  onPresetChange,
  onApplyRange,
}: TimeRangeCardProps): React.JSX.Element {
  const isCustom = preset === 'custom';
  const displayFromTime = normalizeTimeInput(fromTime, '00:00');
  const displayToTime = normalizeTimeInput(toTime, '23:59');
  const [draftFromDate, setDraftFromDate] = useState<Date | undefined>(undefined);
  const [draftToDate, setDraftToDate] = useState<Date | undefined>(undefined);
  const [submitAttempted, setSubmitAttempted] = useState(false);
  const [justApplied, setJustApplied] = useState(false);

  useEffect(() => {
    const parseDraft = (dateValue: string, timeValue: string): Date | undefined => {
      if (!isValidDateInput(dateValue) || !isValidTimeInput(timeValue)) return undefined;
      const candidate = new Date(`${dateValue}T${timeValue}:00`);
      return Number.isNaN(candidate.getTime()) ? undefined : candidate;
    };
    setDraftFromDate(parseDraft(fromDate, fromTime));
    setDraftToDate(parseDraft(toDate, toTime));
    setSubmitAttempted(false);
    setJustApplied(false);
  }, [fromDate, fromTime, preset, toDate, toTime]);

  const showRange = Boolean(fromDate && toDate);
  const draftHasValidInputs = Boolean(draftFromDate && draftToDate);
  const draftRangeValid = useMemo(() => {
    if (!draftHasValidInputs) return false;
    const from = draftFromDate ?? new Date();
    const to = draftToDate ?? new Date();
    return from <= to;
  }, [draftFromDate, draftHasValidInputs, draftToDate]);
  const canApply = draftHasValidInputs && draftRangeValid;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm">Time Range</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <label className="text-xs text-muted-foreground">Quick presets</label>
          <div className="flex flex-wrap items-center gap-3">
            <Tabs value={preset} onValueChange={(value) => {
              onPresetChange(value as DatePreset);
            }}>
              <TabsList>
                <TabsTrigger value="all">All time</TabsTrigger>
                <TabsTrigger value="7d">7d</TabsTrigger>
                <TabsTrigger value="30d">30d</TabsTrigger>
                <TabsTrigger value="90d">90d</TabsTrigger>
                <TabsTrigger value="custom">Custom</TabsTrigger>
              </TabsList>
            </Tabs>
          </div>
          {showRange ? (
            <p className="text-xs text-muted-foreground">
              Showing {fromDate} {displayFromTime} → {toDate} {displayToTime}
            </p>
          ) : null}
        </div>
        <div className={`overflow-hidden transition-all duration-200 ${isCustom ? 'max-h-[420px] opacity-100' : 'max-h-0 opacity-0'}`}>
          <div className="mt-4 space-y-4 rounded-lg border bg-muted/30 p-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-3">
                <Label htmlFor="date-from" className="px-1">
                  From date
                </Label>
                <DateTimePicker
                  value={draftFromDate}
                  onChange={setDraftFromDate}
                  className="w-full"
                  locale={undefined}
                  weekStartsOn={undefined}
                  showWeekNumber={undefined}
                  showOutsideDays={undefined}
                  granularity="minute"
                />
              </div>
              <div className="space-y-3">
                <Label htmlFor="date-to" className="px-1">
                  To date
                </Label>
                <DateTimePicker
                  value={draftToDate}
                  onChange={setDraftToDate}
                  className="w-full"
                  locale={undefined}
                  weekStartsOn={undefined}
                  showWeekNumber={undefined}
                  showOutsideDays={undefined}
                  granularity="minute"
                />
              </div>
            </div>
            <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
              <div className="text-xs text-muted-foreground">
                Use 24-hour time. Example: 2025-01-10 09:30.
              </div>
              <Button
                type="button"
                onClick={() => {
                  setSubmitAttempted(true);
                  if (!canApply) return;
                  if (!draftFromDate || !draftToDate) return;
                  const applied = onApplyRange({
                    fromDate: formatDateInput(draftFromDate),
                    toDate: formatDateInput(draftToDate),
                    fromTime: `${String(draftFromDate.getHours()).padStart(2, '0')}:${String(draftFromDate.getMinutes()).padStart(2, '0')}`,
                    toTime: `${String(draftToDate.getHours()).padStart(2, '0')}:${String(draftToDate.getMinutes()).padStart(2, '0')}`,
                  });
                  if (applied) {
                    setSubmitAttempted(false);
                    setJustApplied(true);
                  }
                }}
                disabled={!canApply}
              >
                {justApplied ? 'Applied' : 'Apply range'}
              </Button>
            </div>
          </div>
        </div>
        {submitAttempted && !draftHasValidInputs ? (
          <p className="text-xs text-destructive">Select both start and end dates before applying.</p>
        ) : null}
        {submitAttempted && draftHasValidInputs && !draftRangeValid ? (
          <p className="text-xs text-destructive">Start date/time must be before end date/time.</p>
        ) : null}
      </CardContent>
    </Card>
  );
}
