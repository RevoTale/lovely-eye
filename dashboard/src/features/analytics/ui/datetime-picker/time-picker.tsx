import { Clock } from 'lucide-react';
import { forwardRef, useImperativeHandle, useRef, useState } from 'react';
import TimePeriodSelect from './time-period-select';
import TimePickerInput from './time-picker-input';
import type { Period, TimePickerProps, TimePickerRef } from './types';

const TimePicker = forwardRef<TimePickerRef, TimePickerProps>(
  ({ date, granularity = 'second', hourCycle = 24, onChange }, ref) => {
    const minuteRef = useRef<HTMLInputElement>(null);
    const hourRef = useRef<HTMLInputElement>(null);
    const secondRef = useRef<HTMLInputElement>(null);
    const periodRef = useRef<HTMLButtonElement>(null);
    const [period, setPeriod] = useState<Period>(date && date.getHours() >= 12 ? 'PM' : 'AM');
    const inputProps = {
      date: date ?? null,
      ...(onChange === undefined ? {} : { onDateChange: onChange }),
    };

    useImperativeHandle(ref, () => ({
      minuteRef: minuteRef.current,
      hourRef: hourRef.current,
      secondRef: secondRef.current,
    }));

    return (
      <div className='flex items-center justify-center gap-2'>
        <label htmlFor='datetime-picker-hour-input' className='cursor-pointer'>
          <Clock className='mr-2 h-4 w-4' />
        </label>
        <TimePickerInput
          {...inputProps}
          picker={hourCycle === 24 ? 'hours' : '12hours'}
          id='datetime-picker-hour-input'
          ref={hourRef}
          period={period}
          onRightFocus={() => minuteRef.current?.focus()}
        />
        {(granularity === 'minute' || granularity === 'second') && (
          <>
            :
            <TimePickerInput
              {...inputProps}
              picker='minutes'
              ref={minuteRef}
              onLeftFocus={() => hourRef.current?.focus()}
              onRightFocus={() => secondRef.current?.focus()}
            />
          </>
        )}
        {granularity === 'second' && (
          <>
            :
            <TimePickerInput
              {...inputProps}
              picker='seconds'
              ref={secondRef}
              onLeftFocus={() => minuteRef.current?.focus()}
              onRightFocus={() => periodRef.current?.focus()}
            />
          </>
        )}
        {hourCycle === 12 && (
          <div className='grid gap-1 text-center'>
            <TimePeriodSelect
              period={period}
              setPeriod={setPeriod}
              date={date ?? null}
              onDateChange={(nextDate) => {
                onChange?.(nextDate);
                setPeriod(nextDate && nextDate.getHours() >= 12 ? 'PM' : 'AM');
              }}
              ref={periodRef}
              onLeftFocus={() => secondRef.current?.focus()}
            />
          </div>
        )}
      </div>
    );
  }
);

TimePicker.displayName = 'TimePicker';

export default TimePicker;
