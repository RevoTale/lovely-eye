import {
  type Dispatch,
  forwardRef,
  type InputHTMLAttributes,
  type KeyboardEvent,
  type SetStateAction,
  useEffect,
  useMemo,
  useState,
} from 'react';
import { cn } from '@/shared/lib/utils';
import { Input } from '@/shared/ui/input';
import { getArrowByType, getDateByType, setDateByType } from './time-utils';
import type { Period, TimePickerType } from './types';

interface TimePickerInputProps extends InputHTMLAttributes<HTMLInputElement> {
  picker: TimePickerType;
  date?: Date | null;
  onDateChange?: (date: Date | undefined) => void;
  period?: Period;
  onRightFocus?: () => void;
  onLeftFocus?: () => void;
}

const TimePickerInput = forwardRef<HTMLInputElement, TimePickerInputProps>(
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
    const [flag, setFlag] = useState(false);
    const [prevIntKey, setPrevIntKey] = useState('0');
    useEffect(() => {
      if (!flag) return undefined;
      const timer = setTimeout(() => setFlag(false), 2000);
      return () => clearTimeout(timer);
    }, [flag]);
    const calculatedValue = useMemo(() => getDateByType(date, picker), [date, picker]);

    const handleKeyDown = (event: KeyboardEvent<HTMLInputElement>): void => {
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
        inputMode='decimal'
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
  setFlag: Dispatch<SetStateAction<boolean>>
): void {
  const step = key === 'ArrowUp' ? 1 : -1;
  const newValue = getArrowByType(calculatedValue, step, picker);
  setFlag(false);
  onDateChange?.(setDateByType(date ? new Date(date) : new Date(), newValue, picker, period));
}

export default TimePickerInput;
