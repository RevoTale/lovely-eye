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
  `domain` varchar NOT NULL,
  `name` varchar NOT NULL,
  `public_key` varchar NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "sites_domain" to table: "sites"
CREATE UNIQUE INDEX `sites_domain` ON `sites` (`domain`);
-- create index "sites_public_key" to table: "sites"
CREATE UNIQUE INDEX `sites_public_key` ON `sites` (`public_key`);
-- create "daily_stats" table
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
-- create "sessions" table
CREATE TABLE `sessions` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `visitor_id` varchar NOT NULL,
  `started_at` timestamp NOT NULL,
  `last_seen_at` timestamp NOT NULL,
  `entry_page` varchar NULL,
  `exit_page` varchar NULL,
  `referrer` varchar NULL,
  `utm_source` varchar NULL,
  `utm_medium` varchar NULL,
  `utm_campaign` varchar NULL,
  `country` varchar NULL,
  `city` varchar NULL,
  `device` varchar NULL,
  `browser` varchar NULL,
  `os` varchar NULL,
  `screen_size` varchar NULL,
  `page_views` integer NULL DEFAULT 0,
  `duration` integer NULL DEFAULT 0,
  `is_bounce` boolean NULL DEFAULT true,
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "events" table
CREATE TABLE `events` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `session_id` integer NULL,
  `visitor_id` varchar NOT NULL,
  `name` varchar NOT NULL,
  `path` varchar NULL,
  `properties` varchar NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "page_views" table
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
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
