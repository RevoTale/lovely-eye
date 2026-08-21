import { forwardRef } from 'react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select';
import { display12HourValue, setDateByType } from './time-utils';
import type { Period } from './types';

interface PeriodSelectorProps {
  period: Period;
  setPeriod?: (period: Period) => void;
  date?: Date | null;
  onDateChange?: (date: Date | undefined) => void;
  onRightFocus?: () => void;
  onLeftFocus?: () => void;
}

const TimePeriodSelect = forwardRef<HTMLButtonElement, PeriodSelectorProps>(
  ({ date, onDateChange, onLeftFocus, onRightFocus, period, setPeriod }, ref) => {
    const handleValueChange = (value: Period): void => {
      setPeriod?.(value);
      if (!date) return;
      const tempDate = new Date(date);
      const hours = display12HourValue(date.getHours());
      onDateChange?.(setDateByType(tempDate, hours, '12hours', period === 'AM' ? 'PM' : 'AM'));
    };

    return (
      <div className='flex h-10 items-center'>
        <Select
          defaultValue={period}
          onValueChange={(value) => {
            if (value !== null) handleValueChange(value);
          }}
        >
          <SelectTrigger
            ref={ref}
            className='w-[65px] focus:bg-accent focus:text-accent-foreground'
            onKeyDown={(event) => {
              if (event.key === 'ArrowRight') onRightFocus?.();
              if (event.key === 'ArrowLeft') onLeftFocus?.();
            }}
          >
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value='AM'>AM</SelectItem>
            <SelectItem value='PM'>PM</SelectItem>
          </SelectContent>
        </Select>
      </div>
    );
  }
);

TimePeriodSelect.displayName = 'TimePeriodSelect';

export default TimePeriodSelect;
