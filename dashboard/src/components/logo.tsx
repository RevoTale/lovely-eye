import React from 'react';
import logoSvg from '/favicon.svg';

interface LogoProps {
  className?: string;
  size?: number;
}

export function Logo({ className, size = 24 }: LogoProps): React.JSX.Element {
  return (
    <img
      src={logoSvg}
      alt="Lovely Eye"
      width={size}
      height={size}
      className={className}
    />
  );
}
