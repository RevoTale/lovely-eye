-- add position to site_domains for ordered domains
ALTER TABLE `site_domains` ADD COLUMN `position` integer NOT NULL DEFAULT 0;

WITH ranked AS (
  SELECT sd.id AS id,
         (SELECT COUNT(*) - 1 FROM site_domains sd2 WHERE sd2.site_id = sd.site_id AND sd2.id <= sd.id) AS position
  FROM site_domains sd
)
UPDATE site_domains
SET position = (SELECT position FROM ranked WHERE ranked.id = site_domains.id);

PRAGMA foreign_keys=off;

-- rebuild sites table without domain column
CREATE TABLE `sites_new` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `user_id` integer NOT NULL,
  `name` varchar NOT NULL,
  `public_key` varchar NOT NULL,
  `track_country` boolean NOT NULL DEFAULT false,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);

INSERT INTO `sites_new` (`id`, `user_id`, `name`, `public_key`, `track_country`, `created_at`, `updated_at`)
SELECT `id`, `user_id`, `name`, `public_key`, `track_country`, `created_at`, `updated_at` FROM `sites`;

DROP TABLE `sites`;
ALTER TABLE `sites_new` RENAME TO `sites`;

CREATE UNIQUE INDEX `sites_public_key` ON `sites` (`public_key`);

PRAGMA foreign_keys=on;
