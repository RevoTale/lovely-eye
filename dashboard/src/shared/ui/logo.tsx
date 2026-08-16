interface LogoProps {
  className?: string;
  size?: number;
}

const DEFAULT_LOGO_SIZE = 24;
// The document base is injected by the Go mount, so this resolves inside any runtime BASE_PATH.
const LOGO_SRC = './favicon.svg';

export const Logo = ({ className, size = DEFAULT_LOGO_SIZE }: LogoProps): React.ReactNode => (
  <img src={LOGO_SRC} alt='Lovely Eye' width={size} height={size} className={className} />
);
