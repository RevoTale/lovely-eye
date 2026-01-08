import React from 'react';

interface LogoProps {
  className?: string;
  size?: number;
}

export function Logo({ className, size = 24 }: LogoProps): React.JSX.Element {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 64 64"
      fill="none"
      width={size}
      height={size}
      className={className}
      aria-label="Lovely Eye logo"
    >
      {/* Outer mechanical ring */}
      <circle cx="32" cy="32" r="28" fill="#1e293b" stroke="#475569" strokeWidth="2"/>
      {/* Inner mechanical ring */}
      <circle cx="32" cy="32" r="22" fill="#0f172a" stroke="#334155" strokeWidth="1.5"/>
      {/* Aperture blades */}
      <path d="M32 14 L38 22 L32 22 Z" fill="#334155"/>
      <path d="M32 14 L26 22 L32 22 Z" fill="#334155"/>
      <path d="M50 32 L42 38 L42 32 Z" fill="#334155"/>
      <path d="M50 32 L42 26 L42 32 Z" fill="#334155"/>
      <path d="M32 50 L26 42 L32 42 Z" fill="#334155"/>
      <path d="M32 50 L38 42 L32 42 Z" fill="#334155"/>
      <path d="M14 32 L22 26 L22 32 Z" fill="#334155"/>
      <path d="M14 32 L22 38 L22 32 Z" fill="#334155"/>
      {/* Glowing iris */}
      <circle cx="32" cy="32" r="12" fill="url(#robotIris)"/>
      {/* Core glow */}
      <circle cx="32" cy="32" r="6" fill="#60a5fa"/>
      <circle cx="32" cy="32" r="4" fill="#93c5fd"/>
      {/* Light reflection */}
      <circle cx="35" cy="29" r="2" fill="white" opacity="0.7"/>
      <defs>
        <radialGradient id="robotIris" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#3b82f6"/>
          <stop offset="70%" stopColor="#1d4ed8"/>
          <stop offset="100%" stopColor="#1e40af"/>
        </radialGradient>
      </defs>
    </svg>
  );
}
