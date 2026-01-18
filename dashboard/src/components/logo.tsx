import React from 'react';

interface LogoProps {
  className?: string;
  size?: number;
}

const DEFAULT_LOGO_SIZE = 24;
const LOGO_SRC = '/favicon.svg';

export function Logo({ className, size = DEFAULT_LOGO_SIZE }: LogoProps): React.JSX.Element {
  return (
    <img
      src={LOGO_SRC}
      alt="Lovely Eye"
      width={size}
      height={size}
      className={className}
    />
  );
}
