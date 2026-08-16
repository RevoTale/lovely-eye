import { cn } from '@/shared/lib/utils';
import { buttonVariants } from '@/shared/ui/button';

export function calendarClassNames(showWeekNumber: boolean | undefined): Record<string, string> {
  return {
    months: 'flex flex-col justify-center space-y-4 sm:flex-row sm:space-y-0',
    month: 'flex flex-col items-center space-y-4',
    month_caption: 'relative flex items-center justify-center pt-1',
    caption_label: 'text-sm font-medium',
    nav: 'flex items-center space-x-1',
    button_previous: cn(
      buttonVariants({ variant: 'outline' }),
      'absolute left-5 top-5 h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100'
    ),
    button_next: cn(
      buttonVariants({ variant: 'outline' }),
      'absolute right-5 top-5 h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100'
    ),
    month_grid: 'w-full border-collapse space-y-1',
    weekdays: cn('flex', showWeekNumber && 'justify-end'),
    weekday: 'text-muted-foreground w-9 rounded-md text-[0.8rem] font-normal',
    week: 'mt-2 flex w-full',
    day: 'relative h-9 w-9 rounded-1 p-0 text-center text-sm [&:has([aria-selected].day-range-end)]:rounded-r-md [&:has([aria-selected].day-outside)]:bg-accent/50 [&:has([aria-selected])]:bg-accent first:[&:has([aria-selected])]:rounded-l-md last:[&:has([aria-selected])]:rounded-r-md focus-within:relative focus-within:z-20',
    day_button: cn(
      buttonVariants({ variant: 'ghost' }),
      'h-9 w-9 rounded-l-md rounded-r-md p-0 font-normal aria-selected:opacity-100'
    ),
    range_end: 'day-range-end',
    selected:
      'rounded-l-md rounded-r-md bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground',
    today: 'bg-accent text-accent-foreground',
    outside:
      'day-outside text-muted-foreground opacity-50 aria-selected:bg-accent/50 aria-selected:text-muted-foreground aria-selected:opacity-30',
    disabled: 'text-muted-foreground opacity-50',
    range_middle: 'aria-selected:bg-accent aria-selected:text-accent-foreground',
    hidden: 'invisible',
  };
}
