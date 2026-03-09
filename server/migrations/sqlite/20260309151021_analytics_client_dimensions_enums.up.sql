-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_clients" table
CREATE TABLE `new_clients` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `hash` varchar NOT NULL,
  `country` varchar NULL,
  `device` integer NOT NULL DEFAULT 0,
  `browser` integer NOT NULL DEFAULT 0,
  `os` integer NOT NULL DEFAULT 0,
  `screen_size` integer NOT NULL DEFAULT 0,
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- copy rows from old table "clients" to new temporary table "new_clients"
INSERT INTO `new_clients` (`id`, `site_id`, `hash`, `country`, `device`, `browser`, `os`, `screen_size`)
SELECT
  `id`,
  `site_id`,
  `hash`,
  `country`,
  CASE
    WHEN `device` IS NULL OR trim(`device`) = '' THEN 0
    WHEN lower(trim(`device`)) = 'console' THEN 1
    WHEN lower(trim(`device`)) = 'desktop' THEN 2
    WHEN lower(trim(`device`)) = 'mobile' THEN 3
    WHEN lower(trim(`device`)) = 'other' THEN 4
    WHEN lower(trim(`device`)) IN ('smart-tv', 'smart tv', 'smarttv') THEN 5
    WHEN lower(trim(`device`)) = 'tablet' THEN 6
    WHEN lower(trim(`device`)) = 'watch' THEN 7
    ELSE 4
  END AS `device`,
  CASE
    WHEN `browser` IS NULL OR trim(`browser`) = '' THEN 0
    WHEN lower(trim(`browser`)) = 'other' THEN 1
    WHEN lower(trim(`browser`)) = 'android webview' THEN 2
    WHEN lower(trim(`browser`)) IN ('chrome', 'chrome mobile ios', 'chrome mobile') THEN 3
    WHEN lower(trim(`browser`)) = 'duckduckgo' THEN 4
    WHEN lower(trim(`browser`)) IN ('edge', 'edge mobile') THEN 5
    WHEN lower(trim(`browser`)) = 'facebook in-app browser' THEN 6
    WHEN lower(trim(`browser`)) IN ('firefox', 'firefox mobile') THEN 7
    WHEN lower(trim(`browser`)) = 'instagram in-app browser' THEN 8
    WHEN lower(trim(`browser`)) = 'internet explorer' THEN 9
    WHEN lower(trim(`browser`)) = 'miui browser' THEN 10
    WHEN lower(trim(`browser`)) IN ('opera', 'opera mobi') THEN 11
    WHEN lower(trim(`browser`)) = 'playstation browser' THEN 12
    WHEN lower(trim(`browser`)) IN ('safari', 'mobile safari') THEN 13
    WHEN lower(trim(`browser`)) = 'samsung internet' THEN 14
    WHEN lower(trim(`browser`)) = 'uc browser' THEN 15
    WHEN lower(trim(`browser`)) = 'vivaldi' THEN 16
    WHEN lower(trim(`browser`)) = 'xbox browser' THEN 17
    WHEN lower(trim(`browser`)) = 'yandex browser' THEN 18
    ELSE 1
  END AS `browser`,
  CASE
    WHEN `os` IS NULL OR trim(`os`) = '' THEN 0
    WHEN lower(trim(`os`)) = 'other' THEN 1
    WHEN lower(trim(`os`)) = 'android' THEN 2
    WHEN lower(trim(`os`)) IN ('chromeos', 'chrome os') THEN 3
    WHEN lower(trim(`os`)) = 'ios' THEN 4
    WHEN lower(trim(`os`)) = 'ipados' THEN 5
    WHEN lower(trim(`os`)) = 'linux' THEN 6
    WHEN lower(trim(`os`)) IN ('macos', 'mac os x', 'os x') THEN 7
    WHEN lower(trim(`os`)) = 'playstation os' THEN 8
    WHEN lower(trim(`os`)) = 'wear os' THEN 9
    WHEN lower(trim(`os`)) IN ('watchos', 'watch os') THEN 10
    WHEN lower(trim(`os`)) = 'windows' THEN 11
    WHEN lower(trim(`os`)) = 'xbox os' THEN 12
    ELSE 1
  END AS `os`,
  CASE
    WHEN `screen_size` IS NULL OR trim(`screen_size`) = '' THEN 0
    WHEN lower(trim(`screen_size`)) = 'watch' THEN 1
    WHEN lower(trim(`screen_size`)) = 'xs' THEN 2
    WHEN lower(trim(`screen_size`)) = 'sm' THEN 3
    WHEN lower(trim(`screen_size`)) = 'md' THEN 4
    WHEN lower(trim(`screen_size`)) = 'lg' THEN 5
    WHEN lower(trim(`screen_size`)) = 'xl' THEN 6
    WHEN instr(lower(trim(`screen_size`)), 'x') > 0
      AND substr(lower(trim(`screen_size`)), 1, instr(lower(trim(`screen_size`)), 'x') - 1) GLOB '[0-9]*'
      THEN CASE
        WHEN CAST(substr(lower(trim(`screen_size`)), 1, instr(lower(trim(`screen_size`)), 'x') - 1) AS integer) < 320 THEN 1
        WHEN CAST(substr(lower(trim(`screen_size`)), 1, instr(lower(trim(`screen_size`)), 'x') - 1) AS integer) < 576 THEN 2
        WHEN CAST(substr(lower(trim(`screen_size`)), 1, instr(lower(trim(`screen_size`)), 'x') - 1) AS integer) < 768 THEN 3
        WHEN CAST(substr(lower(trim(`screen_size`)), 1, instr(lower(trim(`screen_size`)), 'x') - 1) AS integer) < 992 THEN 4
        WHEN CAST(substr(lower(trim(`screen_size`)), 1, instr(lower(trim(`screen_size`)), 'x') - 1) AS integer) < 1200 THEN 5
        ELSE 6
      END
    WHEN trim(`screen_size`) GLOB '[0-9]*'
      THEN CASE
        WHEN CAST(trim(`screen_size`) AS integer) < 320 THEN 1
        WHEN CAST(trim(`screen_size`) AS integer) < 576 THEN 2
        WHEN CAST(trim(`screen_size`) AS integer) < 768 THEN 3
        WHEN CAST(trim(`screen_size`) AS integer) < 992 THEN 4
        WHEN CAST(trim(`screen_size`) AS integer) < 1200 THEN 5
        ELSE 6
      END
    ELSE 0
  END AS `screen_size`
FROM `clients`;
-- drop "clients" table after copying rows
DROP TABLE `clients`;
-- rename temporary table "new_clients" to "clients"
ALTER TABLE `new_clients` RENAME TO `clients`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
