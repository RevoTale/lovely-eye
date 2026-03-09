-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_clients" table
CREATE TABLE `new_clients` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `hash` varchar NOT NULL,
  `country` varchar NULL,
  `device` varchar NULL,
  `browser` varchar NULL,
  `os` varchar NULL,
  `screen_size` varchar NULL,
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- copy rows from current table "clients" to new temporary table "new_clients"
INSERT INTO `new_clients` (`id`, `site_id`, `hash`, `country`, `device`, `browser`, `os`, `screen_size`)
SELECT
  `id`,
  `site_id`,
  `hash`,
  `country`,
  CASE `device`
    WHEN 1 THEN 'console'
    WHEN 2 THEN 'desktop'
    WHEN 3 THEN 'mobile'
    WHEN 4 THEN 'other'
    WHEN 5 THEN 'smart-tv'
    WHEN 6 THEN 'tablet'
    WHEN 7 THEN 'watch'
    ELSE NULL
  END AS `device`,
  CASE `browser`
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
  END AS `browser`,
  CASE `os`
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
  END AS `os`,
  CASE `screen_size`
    WHEN 1 THEN 'watch'
    WHEN 2 THEN 'xs'
    WHEN 3 THEN 'sm'
    WHEN 4 THEN 'md'
    WHEN 5 THEN 'lg'
    WHEN 6 THEN 'xl'
    ELSE NULL
  END AS `screen_size`
FROM `clients`;
-- drop "clients" table after copying rows
DROP TABLE `clients`;
-- rename temporary table "new_clients" to "clients"
ALTER TABLE `new_clients` RENAME TO `clients`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
