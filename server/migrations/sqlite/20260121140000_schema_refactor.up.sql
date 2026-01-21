-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;

-- === STEP 1: Create clients table FIRST ===
CREATE TABLE `clients` (
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

-- === STEP 2: Populate clients from page_views ===
INSERT OR IGNORE INTO `clients` (`site_id`, `hash`)
SELECT DISTINCT pv.`site_id`, pv.`visitor_id` AS `hash`
FROM `page_views` pv;

-- === STEP 3: Migrate sessions table ===
-- create "new_sessions" table
CREATE TABLE `new_sessions` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `client_id` integer NOT NULL,
  `enter_time` integer NOT NULL,
  `enter_hour` integer NOT NULL,
  `enter_day` integer NOT NULL,
  `enter_path` varchar NOT NULL,
  `exit_time` integer NOT NULL,
  `exit_hour` integer NOT NULL,
  `exit_day` integer NOT NULL,
  `exit_path` varchar NOT NULL,
  `referrer` varchar NULL,
  `utm_source` varchar NULL,
  `utm_medium` varchar NULL,
  `utm_campaign` varchar NULL,
  `duration` integer NOT NULL DEFAULT 0,
  `page_view_count` integer NOT NULL DEFAULT 0,
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Populate sessions with data from page_views
INSERT INTO `new_sessions` (
  `id`, `site_id`, `client_id`,
  `enter_time`, `enter_hour`, `enter_day`, `enter_path`,
  `exit_time`, `exit_hour`, `exit_day`, `exit_path`,
  `referrer`, `utm_source`, `utm_medium`, `utm_campaign`,
  `duration`, `page_view_count`
)
SELECT
  s.`id`,
  s.`site_id`,
  (SELECT c.`id` FROM `clients` c, `page_views` pv WHERE c.`hash` = pv.`visitor_id` AND pv.`session_id` = s.`id` LIMIT 1) as `client_id`,
  COALESCE((SELECT CAST(strftime('%s', pv.`created_at`) AS INTEGER) FROM `page_views` pv WHERE pv.`session_id` = s.`id` ORDER BY pv.`created_at` ASC LIMIT 1), 0) as `enter_time`,
  COALESCE((SELECT CAST(strftime('%s', pv.`created_at`) AS INTEGER) / 3600 FROM `page_views` pv WHERE pv.`session_id` = s.`id` ORDER BY pv.`created_at` ASC LIMIT 1), 0) as `enter_hour`,
  COALESCE((SELECT CAST(strftime('%s', pv.`created_at`) AS INTEGER) / 86400 FROM `page_views` pv WHERE pv.`session_id` = s.`id` ORDER BY pv.`created_at` ASC LIMIT 1), 0) as `enter_day`,
  COALESCE((SELECT pv.`path` FROM `page_views` pv WHERE pv.`session_id` = s.`id` ORDER BY pv.`created_at` ASC LIMIT 1), '/') as `enter_path`,
  COALESCE((SELECT CAST(strftime('%s', pv.`created_at`) AS INTEGER) FROM `page_views` pv WHERE pv.`session_id` = s.`id` ORDER BY pv.`created_at` DESC LIMIT 1), 0) as `exit_time`,
  COALESCE((SELECT CAST(strftime('%s', pv.`created_at`) AS INTEGER) / 3600 FROM `page_views` pv WHERE pv.`session_id` = s.`id` ORDER BY pv.`created_at` DESC LIMIT 1), 0) as `exit_hour`,
  COALESCE((SELECT CAST(strftime('%s', pv.`created_at`) AS INTEGER) / 86400 FROM `page_views` pv WHERE pv.`session_id` = s.`id` ORDER BY pv.`created_at` DESC LIMIT 1), 0) as `exit_day`,
  COALESCE((SELECT pv.`path` FROM `page_views` pv WHERE pv.`session_id` = s.`id` ORDER BY pv.`created_at` DESC LIMIT 1), '/') as `exit_path`,
  s.`referrer`,
  s.`utm_source`,
  s.`utm_medium`,
  s.`utm_campaign`,
  IFNULL(s.`duration`, 0) as `duration`,
  COALESCE((SELECT COUNT(*) FROM `page_views` pv WHERE pv.`session_id` = s.`id`), 0) as `page_view_count`
FROM `sessions` s;

-- drop "sessions" table after copying rows
DROP TABLE `sessions`;
-- rename temporary table "new_sessions" to "sessions"
ALTER TABLE `new_sessions` RENAME TO `sessions`;

-- === STEP 4: Migrate events table ===
-- create "new_events" table
CREATE TABLE `new_events` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `session_id` integer NOT NULL,
  `time` integer NOT NULL,
  `hour` integer NOT NULL,
  `day` integer NOT NULL,
  `path` varchar NOT NULL,
  `name` varchar NOT NULL,
  `type` integer NOT NULL,
  `definition_id` integer NULL,
  CONSTRAINT `0` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`definition_id`) REFERENCES `event_definitions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- copy rows from old table "events" to new temporary table "new_events"
-- Convert created_at timestamp to unix time and calculate hour/day buckets
INSERT INTO `new_events` (`id`, `session_id`, `time`, `hour`, `day`, `path`, `name`, `type`, `definition_id`)
SELECT
  `id`,
  `session_id`,
  CAST(strftime('%s', `created_at`) AS INTEGER) AS `time`,
  CAST(strftime('%s', `created_at`) AS INTEGER) / 3600 AS `hour`,
  CAST(strftime('%s', `created_at`) AS INTEGER) / 86400 AS `day`,
  COALESCE(`path`, ''),
  `name`,
  1 AS `type`,
  NULL AS `definition_id`
FROM `events`;

-- Migrate page_views to events table with Type=0 (pageview)
INSERT INTO `new_events` (`session_id`, `time`, `hour`, `day`, `path`, `name`, `type`, `definition_id`)
SELECT
  pv.`session_id`,
  CAST(strftime('%s', pv.`created_at`) AS INTEGER) AS `time`,
  CAST(strftime('%s', pv.`created_at`) AS INTEGER) / 3600 AS `hour`,
  CAST(strftime('%s', pv.`created_at`) AS INTEGER) / 86400 AS `day`,
  pv.`path`,
  COALESCE(pv.`title`, pv.`path`) AS `name`,
  0 AS `type`,
  NULL AS `definition_id`
FROM `page_views` pv
WHERE pv.`session_id` IS NOT NULL;

-- drop "events" table after copying rows
DROP TABLE `events`;
-- rename temporary table "new_events" to "events"
ALTER TABLE `new_events` RENAME TO `events`;

-- drop "page_views" table
DROP TABLE `page_views`;

-- === STEP 5: Modify event_definition_fields ===
-- create "new_event_definition_fields" table
CREATE TABLE `new_event_definition_fields` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `event_definition_id` integer NOT NULL,
  `key` varchar NOT NULL,
  `type` integer NOT NULL,
  `required` boolean NOT NULL DEFAULT false,
  `max_length` integer NOT NULL DEFAULT 500,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`event_definition_id`) REFERENCES `event_definitions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- copy rows from old table "event_definition_fields" to new temporary table "new_event_definition_fields"
INSERT INTO `new_event_definition_fields` (`id`, `event_definition_id`, `key`, `type`, `required`, `max_length`, `created_at`, `updated_at`)
SELECT `id`, `event_definition_id`, `key`, `type`, `required`, `max_length`, `created_at`, `updated_at` FROM `event_definition_fields`;
-- drop "event_definition_fields" table after copying rows
DROP TABLE `event_definition_fields`;
-- rename temporary table "new_event_definition_fields" to "event_definition_fields"
ALTER TABLE `new_event_definition_fields` RENAME TO `event_definition_fields`;

-- === STEP 6: Create event_data table ===
CREATE TABLE `event_data` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `event_id` integer NOT NULL,
  `field_id` integer NOT NULL,
  `value` varchar NOT NULL,
  CONSTRAINT `0` FOREIGN KEY (`field_id`) REFERENCES `event_definition_fields` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- === STEP 7: Clean up old indexes and drop daily_stats ===
-- drop "daily_stats" table
DROP TABLE IF EXISTS `daily_stats`;
-- drop index "site_blocked_ips_site_id_ip" from table: "site_blocked_ips"
DROP INDEX IF EXISTS `site_blocked_ips_site_id_ip`;
-- drop index "site_blocked_countries_site_id_country" from table: "site_blocked_countries"
DROP INDEX IF EXISTS `site_blocked_countries_site_id_country`;

-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
