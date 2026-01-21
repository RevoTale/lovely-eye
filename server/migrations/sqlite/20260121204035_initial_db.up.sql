-- create "users" table
CREATE TABLE `users` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `username` varchar NOT NULL,
  `password_hash` varchar NOT NULL,
  `role` varchar NOT NULL DEFAULT 'user',
  `email` varchar NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp)
);
-- create index "users_username" to table: "users"
CREATE UNIQUE INDEX `users_username` ON `users` (`username`);
-- create "sites" table
CREATE TABLE `sites` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `user_id` integer NOT NULL,
  `name` varchar NOT NULL,
  `public_key` varchar NOT NULL,
  `track_country` boolean NOT NULL DEFAULT false,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "sites_public_key" to table: "sites"
CREATE UNIQUE INDEX `sites_public_key` ON `sites` (`public_key`);
-- create "clients" table
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
-- create "sessions" table
CREATE TABLE `sessions` (
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
-- create "event_definitions" table
CREATE TABLE `event_definitions` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `name` varchar NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "events" table
CREATE TABLE `events` (
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
-- create "event_definition_fields" table
CREATE TABLE `event_definition_fields` (
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
-- create "event_data" table
CREATE TABLE `event_data` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `event_id` integer NOT NULL,
  `field_id` integer NOT NULL,
  `value` varchar NOT NULL,
  CONSTRAINT `0` FOREIGN KEY (`field_id`) REFERENCES `event_definition_fields` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "site_blocked_countries" table
CREATE TABLE `site_blocked_countries` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `country_code` varchar NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "site_blocked_ips" table
CREATE TABLE `site_blocked_ips` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `ip` varchar NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "site_domains" table
CREATE TABLE `site_domains` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `domain` varchar NOT NULL,
  `position` integer NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "site_domains_domain" to table: "site_domains"
CREATE UNIQUE INDEX `site_domains_domain` ON `site_domains` (`domain`);
