import { add, format } from 'date-fns';
import { enUS } from 'date-fns/locale';
import { Calendar as CalendarIcon } from 'lucide-react';
import { forwardRef, useMemo, useState } from 'react';
import { cn } from '@/shared/lib/utils';
import { Button } from '@/shared/ui/button';
import { Calendar as CalendarComponent } from '@/shared/ui/calendar';
import { Popover, PopoverContent, PopoverTrigger } from '@/shared/ui/popover';
import { calendarClassNames } from './calendar-class-names';
import MonthCaption from './month-caption';
import TimePicker from './time-picker';
import { genMonths, genYears, normalizedLocale } from './time-utils';
import type { DateTimePickerProps, DateTimePickerRef } from './types';

const DateTimePicker = forwardRef<DateTimePickerRef, DateTimePickerProps>(
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

    return (
      <Popover>
        <PopoverTrigger asChild disabled={disabled}>
          <Button
            variant='outline'
            className={cn(
              'w-full justify-start text-left font-normal',
              !value && 'text-muted-foreground',
              className
            )}
            ref={ref}
          >
            <CalendarIcon className='mr-2 h-4 w-4' />
            {value ? (
              format(value, hourCycle === 24 ? hourFormats.hour24 : hourFormats.hour12, {
                locale: mergedLocale,
              })
            ) : (
              <span>{placeholder}</span>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className='w-auto p-0'>
          <CalendarComponent
            mode='single'
            selected={value}
            month={month}
            onSelect={handleSelect}
            onMonthChange={handleSelect}
            locale={locale}
            buttonVariant='outline'
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
            <div className='border-t border-border p-3'>
              <TimePicker
                {...(onChange === undefined ? {} : { onChange })}
                date={value ?? null}
                hourCycle={hourCycle}
                granularity={granularity}
              />
            </div>
          )}
        </PopoverContent>
      </Popover>
    );
  }
);

DateTimePicker.displayName = 'DateTimePicker';

export default DateTimePicker;
