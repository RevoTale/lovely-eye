-- create "site_blocked_ips" table
CREATE TABLE `site_blocked_ips` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `ip` varchar NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "site_blocked_ips_site_id_ip" to table: "site_blocked_ips"
CREATE UNIQUE INDEX `site_blocked_ips_site_id_ip` ON `site_blocked_ips` (`site_id`, `ip`);
-- create "site_blocked_countries" table
CREATE TABLE `site_blocked_countries` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `country_code` varchar NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `1` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "site_blocked_countries_site_id_country" to table: "site_blocked_countries"
CREATE UNIQUE INDEX `site_blocked_countries_site_id_country` ON `site_blocked_countries` (`site_id`, `country_code`);
