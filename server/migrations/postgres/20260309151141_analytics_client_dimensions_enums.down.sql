-- reverse: modify "clients" table
ALTER TABLE "public"."clients"
  ALTER COLUMN "device" DROP DEFAULT,
  ALTER COLUMN "browser" DROP DEFAULT,
  ALTER COLUMN "os" DROP DEFAULT,
  ALTER COLUMN "screen_size" DROP DEFAULT;

ALTER TABLE "public"."clients"
  ALTER COLUMN "device" TYPE character varying(10) USING CASE "device"
    WHEN 1 THEN 'console'
    WHEN 2 THEN 'desktop'
    WHEN 3 THEN 'mobile'
    WHEN 4 THEN 'other'
    WHEN 5 THEN 'smart-tv'
    WHEN 6 THEN 'tablet'
    WHEN 7 THEN 'watch'
    ELSE NULL
  END,
  ALTER COLUMN "browser" TYPE character varying(32) USING CASE "browser"
    WHEN 1 THEN 'Other'
    WHEN 2 THEN 'Android WebView'
    WHEN 3 THEN 'Chrome'
    WHEN 4 THEN 'DuckDuckGo'
    WHEN 5 THEN 'Edge'
    WHEN 6 THEN 'Facebook In-App Browser'
    WHEN 7 THEN 'Firefox'
    WHEN 8 THEN 'Instagram In-App Browser'
    WHEN 9 THEN 'Internet Explorer'
    WHEN 10 THEN 'MIUI Browser'
    WHEN 11 THEN 'Opera'
    WHEN 12 THEN 'PlayStation Browser'
    WHEN 13 THEN 'Safari'
    WHEN 14 THEN 'Samsung Internet'
    WHEN 15 THEN 'UC Browser'
    WHEN 16 THEN 'Vivaldi'
    WHEN 17 THEN 'Xbox Browser'
    WHEN 18 THEN 'Yandex Browser'
    ELSE NULL
  END,
  ALTER COLUMN "os" TYPE character varying(32) USING CASE "os"
    WHEN 1 THEN 'Other'
    WHEN 2 THEN 'Android'
    WHEN 3 THEN 'ChromeOS'
    WHEN 4 THEN 'iOS'
    WHEN 5 THEN 'iPadOS'
    WHEN 6 THEN 'Linux'
    WHEN 7 THEN 'macOS'
    WHEN 8 THEN 'PlayStation OS'
    WHEN 9 THEN 'Wear OS'
    WHEN 10 THEN 'watchOS'
    WHEN 11 THEN 'Windows'
    WHEN 12 THEN 'Xbox OS'
    ELSE NULL
  END,
  ALTER COLUMN "screen_size" TYPE character varying(16) USING CASE "screen_size"
    WHEN 1 THEN 'watch'
    WHEN 2 THEN 'xs'
    WHEN 3 THEN 'sm'
    WHEN 4 THEN 'md'
    WHEN 5 THEN 'lg'
    WHEN 6 THEN 'xl'
    ELSE NULL
  END;

ALTER TABLE "public"."clients"
  ALTER COLUMN "device" DROP NOT NULL,
  ALTER COLUMN "browser" DROP NOT NULL,
  ALTER COLUMN "os" DROP NOT NULL,
  ALTER COLUMN "screen_size" DROP NOT NULL;
