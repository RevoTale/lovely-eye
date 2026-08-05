import * as React from 'react';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { getArrowByType, getDateByType, setDateByType } from './time-utils';
import type { Period, TimePickerType } from './types';

interface TimePickerInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  picker: TimePickerType;
  date?: Date | null;
  onDateChange?: (date: Date | undefined) => void;
  period?: Period;
  onRightFocus?: () => void;
  onLeftFocus?: () => void;
}

const TimePickerInput = React.forwardRef<HTMLInputElement, TimePickerInputProps>(
  (
    {
      className,
      date = new Date(new Date().setHours(0, 0, 0, 0)),
      id,
      name,
      onChange,
      onDateChange,
      onKeyDown,
      onLeftFocus,
      onRightFocus,
      period,
      picker,
      type = 'tel',
      value,
      ...props
    },
    ref
  ) => {
    const [flag, setFlag] = React.useState(false);
    const [prevIntKey, setPrevIntKey] = React.useState('0');
    React.useEffect(() => {
      if (!flag) return undefined;
      const timer = setTimeout(() => setFlag(false), 2000);
      return () => clearTimeout(timer);
    }, [flag]);
    const calculatedValue = React.useMemo(() => getDateByType(date, picker), [date, picker]);

    const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>): void => {
      if (event.key === 'Tab') return;
      event.preventDefault();
      if (event.key === 'ArrowRight') onRightFocus?.();
      if (event.key === 'ArrowLeft') onLeftFocus?.();
      if (event.key === 'ArrowUp' || event.key === 'ArrowDown') {
        applyArrowKey(event.key, calculatedValue, picker, date, period, onDateChange, setFlag);
      }
      if (event.key >= '0' && event.key <= '9') {
        if (picker === '12hours') setPrevIntKey(event.key);
        const newValue = calculateNewValue(event.key, picker, calculatedValue, flag, prevIntKey);
        if (flag) onRightFocus?.();
        setFlag((prev) => !prev);
        onDateChange?.(setDateByType(date ? new Date(date) : new Date(), newValue, picker, period));
      }
    };

    return (
      <Input
        ref={ref}
        id={id || picker}
        name={name || picker}
        className={cn(
          'w-[48px] text-center font-mono text-base tabular-nums caret-transparent focus:bg-accent focus:text-accent-foreground [&::-webkit-inner-spin-button]:appearance-none',
          className
        )}
        value={value || calculatedValue}
        onChange={(event) => {
          event.preventDefault();
          onChange?.(event);
        }}
        type={type}
        inputMode="decimal"
        onKeyDown={(event) => {
          onKeyDown?.(event);
          handleKeyDown(event);
        }}
        {...props}
      />
    );
  }
);

TimePickerInput.displayName = 'TimePickerInput';

function calculateNewValue(
  key: string,
  picker: TimePickerType,
  calculatedValue: string,
  flag: boolean,
  prevIntKey: string
): string {
  if (picker === '12hours' && flag && calculatedValue.slice(1, 2) === '1' && prevIntKey === '0') {
    return `0${key}`;
  }
  return !flag ? `0${key}` : calculatedValue.slice(1, 2) + key;
}

function applyArrowKey(
  key: string,
  calculatedValue: string,
  picker: TimePickerType,
  date: Date | null | undefined,
  period: Period | undefined,
  onDateChange: ((date: Date | undefined) => void) | undefined,
  setFlag: React.Dispatch<React.SetStateAction<boolean>>
): void {
  const step = key === 'ArrowUp' ? 1 : -1;
  const newValue = getArrowByType(calculatedValue, step, picker);
  setFlag(false);
  onDateChange?.(setDateByType(date ? new Date(date) : new Date(), newValue, picker, period));
}

export default TimePickerInput;
