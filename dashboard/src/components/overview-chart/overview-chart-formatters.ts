const compactNumberFormatter = new Intl.NumberFormat('en-US', {
  notation: 'compact',
  maximumFractionDigits: 1,
});

export const formatOverviewAxisValue = (value: number): string => compactNumberFormatter.format(value);

export const formatOverviewValue = (value: number): string => value.toLocaleString();
