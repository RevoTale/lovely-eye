import type { Locale } from 'date-fns/locale';
import type { ComponentProps } from 'react';
import type { Calendar as CalendarComponent } from '@/shared/ui/calendar';

type CalendarComponentProps = ComponentProps<typeof CalendarComponent>;

export type Period = 'AM' | 'PM';
export type TimePickerType = 'minutes' | 'seconds' | 'hours' | '12hours';
export type Granularity = 'day' | 'hour' | 'minute' | 'second';

export interface TimePickerProps {
  date?: Date | null;
  onChange?: (date: Date | undefined) => void;
  hourCycle?: 12 | 24;
  granularity?: Granularity;
}

export interface TimePickerRef {
  minuteRef: HTMLInputElement | null;
  hourRef: HTMLInputElement | null;
  secondRef: HTMLInputElement | null;
}

export type DateTimePickerProps = {
  value?: Date;
  onChange?: (date: Date | undefined) => void;
  disabled?: boolean;
  hourCycle?: 12 | 24;
  placeholder?: string;
  yearRange?: number;
  displayFormat?: { hour24?: string; hour12?: string };
  granularity?: Granularity;
  className?: string;
} & Pick<CalendarComponentProps, 'locale' | 'weekStartsOn' | 'showWeekNumber' | 'showOutsideDays'>;

export type DateTimePickerRef = HTMLButtonElement;

export type DateTimeLocale = Pick<Locale, 'options' | 'localize' | 'formatLong'>;
