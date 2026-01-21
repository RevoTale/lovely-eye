-- create "site_domains" table
CREATE TABLE `site_domains` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `domain` varchar NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "site_domains_domain" to table: "site_domains"
CREATE UNIQUE INDEX `site_domains_domain` ON `site_domains` (`domain`);
-- backfill existing site domains
INSERT INTO `site_domains` (`site_id`, `domain`, `created_at`, `updated_at`)
SELECT `id`, `domain`, `created_at`, `updated_at` FROM `sites`;
