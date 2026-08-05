import * as React from 'react';
import { useImperativeHandle, useMemo, useRef, useState } from 'react';
import { add, format } from 'date-fns';
import { enUS } from 'date-fns/locale';
import { Calendar as CalendarIcon } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Calendar as CalendarComponent } from '@/components/ui/calendar';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import { calendarClassNames } from './calendar-class-names';
import MonthCaption from './month-caption';
import TimePicker from './time-picker';
import { genMonths, genYears, normalizedLocale } from './time-utils';
import type { DateTimePickerProps, DateTimePickerRef } from './types';

const DateTimePicker = React.forwardRef<Partial<DateTimePickerRef>, DateTimePickerProps>(
  (
    {
      className,
      disabled = false,
      displayFormat,
      granularity = 'second',
      hourCycle = 24,
      locale = enUS,
      onChange,
      placeholder = 'Pick a date',
      value,
      yearRange = 50,
      ...props
    },
    ref
  ) => {
    const [month, setMonth] = useState<Date>(value ?? new Date());
    const buttonRef = useRef<HTMLButtonElement>(null);
    const monthOptions = useMemo(() => genMonths(normalizedLocale(locale)), [locale]);
    const yearOptions = useMemo(() => genYears(yearRange), [yearRange]);
    const mergedLocale = useMemo(() => normalizedLocale(locale), [locale]);
    const showWeekNumber = props.showWeekNumber;
    const hourFormats = {
      hour24:
        displayFormat?.hour24 ??
        `PPP HH:mm${!granularity || granularity === 'second' ? ':ss' : ''}`,
      hour12:
        displayFormat?.hour12 ??
        `PP hh:mm${!granularity || granularity === 'second' ? ':ss' : ''} b`,
    };

    const handleSelect = (newDay: Date | undefined): void => {
      if (!newDay) return;
      if (!value) {
        onChange?.(newDay);
        setMonth(newDay);
        return;
      }
      const diffInDays = (newDay.getTime() - value.getTime()) / (1000 * 60 * 60 * 24);
      const newDateFull = add(value, { days: Math.ceil(diffInDays) });
      onChange?.(newDateFull);
      setMonth(newDateFull);
    };

    useImperativeHandle(
      ref,
      () => ({
        ...buttonRef.current,
        value,
      }),
      [value]
    );

    return (
      <Popover>
        <PopoverTrigger asChild disabled={disabled}>
          <Button
            variant="outline"
            className={cn('w-full justify-start text-left font-normal', !value && 'text-muted-foreground', className)}
            ref={buttonRef}
          >
            <CalendarIcon className="mr-2 h-4 w-4" />
            {value ? (
              format(value, hourCycle === 24 ? hourFormats.hour24 : hourFormats.hour12, { locale: mergedLocale })
            ) : (
              <span>{placeholder}</span>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0">
          <CalendarComponent
            mode="single"
            selected={value}
            month={month}
            onSelect={handleSelect}
            onMonthChange={handleSelect}
            locale={locale}
            buttonVariant="outline"
            classNames={calendarClassNames(showWeekNumber)}
            components={{
              MonthCaption: (captionProps) => (
                <MonthCaption
                  {...captionProps}
                  monthOptions={monthOptions}
                  onSelectDate={handleSelect}
                  yearOptions={yearOptions}
                />
              ),
            }}
            {...props}
          />
          {granularity !== 'day' && (
            <div className="border-t border-border p-3">
              <TimePicker onChange={onChange} date={value} hourCycle={hourCycle} granularity={granularity} />
            </div>
          )}
        </PopoverContent>
      </Popover>
    );
  }
);

DateTimePicker.displayName = 'DateTimePicker';

export default DateTimePicker;
