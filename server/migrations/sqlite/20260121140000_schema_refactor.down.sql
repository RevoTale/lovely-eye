-- reverse: create "event_data" table
DROP TABLE `event_data`;
-- reverse: create "clients" table
DROP TABLE `clients`;
-- reverse: drop index "site_blocked_countries_site_id_country" from table: "site_blocked_countries"
CREATE UNIQUE INDEX `site_blocked_countries_site_id_country` ON `site_blocked_countries` (`site_id`, `country_code`);
-- reverse: drop index "site_blocked_ips_site_id_ip" from table: "site_blocked_ips"
CREATE UNIQUE INDEX `site_blocked_ips_site_id_ip` ON `site_blocked_ips` (`site_id`, `ip`);
-- reverse: create "new_event_definition_fields" table
DROP TABLE `new_event_definition_fields`;
-- reverse: drop "page_views" table
CREATE TABLE `page_views` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `session_id` integer NULL,
  `visitor_id` varchar NOT NULL,
  `path` varchar NOT NULL,
  `title` varchar NULL,
  `referrer` varchar NULL,
  `duration` integer NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `1` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- reverse: create "new_events" table
DROP TABLE `new_events`;
-- reverse: create "new_sessions" table
DROP TABLE `new_sessions`;
-- reverse: drop "daily_stats" table
CREATE TABLE `daily_stats` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `date` timestamp NOT NULL,
  `visitors` integer NULL DEFAULT 0,
  `page_views` integer NULL DEFAULT 0,
  `sessions` integer NULL DEFAULT 0,
  `bounce_rate` double precision NULL DEFAULT 0,
  `avg_duration` double precision NULL DEFAULT 0,
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
