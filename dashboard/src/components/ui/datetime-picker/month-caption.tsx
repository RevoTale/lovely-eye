import type { FunctionComponent } from 'react';
import type { MonthCaptionProps } from '@daypicker/react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

interface DateTimeMonthCaptionProps extends MonthCaptionProps {
  monthOptions: Array<{ value: number; label: string }>;
  yearOptions: Array<{ value: number; label: string }>;
  onSelectDate: (date: Date | undefined) => void;
}

const MonthCaption: FunctionComponent<DateTimeMonthCaptionProps> = ({
  calendarMonth,
  monthOptions,
  onSelectDate,
  yearOptions,
}) => (
  <div className="inline-flex gap-2">
    <Select
      defaultValue={calendarMonth.date.getMonth().toString()}
      onValueChange={(value) => {
        const newDate = new Date(calendarMonth.date);
        newDate.setMonth(Number.parseInt(value, 10));
        onSelectDate(newDate);
      }}
    >
      <SelectTrigger className="w-fit gap-1 border-none p-0 focus:bg-accent focus:text-accent-foreground">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        {monthOptions.map((monthOption) => (
          <SelectItem key={monthOption.value} value={monthOption.value.toString()}>
            {monthOption.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
    <Select
      defaultValue={calendarMonth.date.getFullYear().toString()}
      onValueChange={(value) => {
        const newDate = new Date(calendarMonth.date);
        newDate.setFullYear(Number.parseInt(value, 10));
        onSelectDate(newDate);
      }}
    >
      <SelectTrigger className="w-fit gap-1 border-none p-0 focus:bg-accent focus:text-accent-foreground">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        {yearOptions.map((yearOption) => (
          <SelectItem key={yearOption.value} value={yearOption.value.toString()}>
            {yearOption.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  </div>
);

export default MonthCaption;
