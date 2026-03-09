-- modify "clients" table
ALTER TABLE "public"."clients"
  ALTER COLUMN "device" DROP DEFAULT,
  ALTER COLUMN "browser" DROP DEFAULT,
  ALTER COLUMN "os" DROP DEFAULT,
  ALTER COLUMN "screen_size" DROP DEFAULT;

ALTER TABLE "public"."clients"
  ALTER COLUMN "device" TYPE smallint USING CASE
    WHEN "device" IS NULL OR btrim("device") = '' THEN 0
    WHEN lower(btrim("device")) = 'console' THEN 1
    WHEN lower(btrim("device")) = 'desktop' THEN 2
    WHEN lower(btrim("device")) = 'mobile' THEN 3
    WHEN lower(btrim("device")) = 'other' THEN 4
    WHEN lower(btrim("device")) IN ('smart-tv', 'smart tv', 'smarttv') THEN 5
    WHEN lower(btrim("device")) = 'tablet' THEN 6
    WHEN lower(btrim("device")) = 'watch' THEN 7
    ELSE 4
  END,
  ALTER COLUMN "browser" TYPE smallint USING CASE
    WHEN "browser" IS NULL OR btrim("browser") = '' THEN 0
    WHEN lower(btrim("browser")) = 'other' THEN 1
    WHEN lower(btrim("browser")) = 'android webview' THEN 2
    WHEN lower(btrim("browser")) IN ('chrome', 'chrome mobile ios', 'chrome mobile') THEN 3
    WHEN lower(btrim("browser")) = 'duckduckgo' THEN 4
    WHEN lower(btrim("browser")) IN ('edge', 'edge mobile') THEN 5
    WHEN lower(btrim("browser")) = 'facebook in-app browser' THEN 6
    WHEN lower(btrim("browser")) IN ('firefox', 'firefox mobile') THEN 7
    WHEN lower(btrim("browser")) = 'instagram in-app browser' THEN 8
    WHEN lower(btrim("browser")) = 'internet explorer' THEN 9
    WHEN lower(btrim("browser")) = 'miui browser' THEN 10
    WHEN lower(btrim("browser")) IN ('opera', 'opera mobi') THEN 11
    WHEN lower(btrim("browser")) = 'playstation browser' THEN 12
    WHEN lower(btrim("browser")) IN ('safari', 'mobile safari') THEN 13
    WHEN lower(btrim("browser")) = 'samsung internet' THEN 14
    WHEN lower(btrim("browser")) = 'uc browser' THEN 15
    WHEN lower(btrim("browser")) = 'vivaldi' THEN 16
    WHEN lower(btrim("browser")) = 'xbox browser' THEN 17
    WHEN lower(btrim("browser")) = 'yandex browser' THEN 18
    ELSE 1
  END,
  ALTER COLUMN "os" TYPE smallint USING CASE
    WHEN "os" IS NULL OR btrim("os") = '' THEN 0
    WHEN lower(btrim("os")) = 'other' THEN 1
    WHEN lower(btrim("os")) = 'android' THEN 2
    WHEN lower(btrim("os")) IN ('chromeos', 'chrome os') THEN 3
    WHEN lower(btrim("os")) = 'ios' THEN 4
    WHEN lower(btrim("os")) = 'ipados' THEN 5
    WHEN lower(btrim("os")) = 'linux' THEN 6
    WHEN lower(btrim("os")) IN ('macos', 'mac os x', 'os x') THEN 7
    WHEN lower(btrim("os")) = 'playstation os' THEN 8
    WHEN lower(btrim("os")) = 'wear os' THEN 9
    WHEN lower(btrim("os")) IN ('watchos', 'watch os') THEN 10
    WHEN lower(btrim("os")) = 'windows' THEN 11
    WHEN lower(btrim("os")) = 'xbox os' THEN 12
    ELSE 1
  END,
  ALTER COLUMN "screen_size" TYPE smallint USING CASE
    WHEN "screen_size" IS NULL OR btrim("screen_size") = '' THEN 0
    WHEN lower(btrim("screen_size")) = 'watch' THEN 1
    WHEN lower(btrim("screen_size")) = 'xs' THEN 2
    WHEN lower(btrim("screen_size")) = 'sm' THEN 3
    WHEN lower(btrim("screen_size")) = 'md' THEN 4
    WHEN lower(btrim("screen_size")) = 'lg' THEN 5
    WHEN lower(btrim("screen_size")) = 'xl' THEN 6
    WHEN lower(btrim("screen_size")) ~ '^[0-9]+x[0-9]+$' THEN CASE
      WHEN split_part(lower(btrim("screen_size")), 'x', 1)::integer < 320 THEN 1
      WHEN split_part(lower(btrim("screen_size")), 'x', 1)::integer < 576 THEN 2
      WHEN split_part(lower(btrim("screen_size")), 'x', 1)::integer < 768 THEN 3
      WHEN split_part(lower(btrim("screen_size")), 'x', 1)::integer < 992 THEN 4
      WHEN split_part(lower(btrim("screen_size")), 'x', 1)::integer < 1200 THEN 5
      ELSE 6
    END
    WHEN lower(btrim("screen_size")) ~ '^[0-9]+$' THEN CASE
      WHEN lower(btrim("screen_size"))::integer < 320 THEN 1
      WHEN lower(btrim("screen_size"))::integer < 576 THEN 2
      WHEN lower(btrim("screen_size"))::integer < 768 THEN 3
      WHEN lower(btrim("screen_size"))::integer < 992 THEN 4
      WHEN lower(btrim("screen_size"))::integer < 1200 THEN 5
      ELSE 6
    END
    ELSE 0
  END;

ALTER TABLE "public"."clients"
  ALTER COLUMN "device" SET DEFAULT 0,
  ALTER COLUMN "device" SET NOT NULL,
  ALTER COLUMN "browser" SET DEFAULT 0,
  ALTER COLUMN "browser" SET NOT NULL,
  ALTER COLUMN "os" SET DEFAULT 0,
  ALTER COLUMN "os" SET NOT NULL,
  ALTER COLUMN "screen_size" SET DEFAULT 0,
  ALTER COLUMN "screen_size" SET NOT NULL;
