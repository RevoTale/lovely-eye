PRAGMA foreign_keys=off;

-- rebuild sites table with domain column restored
CREATE TABLE `sites_new` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `user_id` integer NOT NULL,
  `domain` varchar NOT NULL,
  `name` varchar NOT NULL,
  `public_key` varchar NOT NULL,
  `track_country` boolean NOT NULL DEFAULT false,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);

INSERT INTO `sites_new` (`id`, `user_id`, `domain`, `name`, `public_key`, `track_country`, `created_at`, `updated_at`)
SELECT s.id,
       s.user_id,
       COALESCE((
         SELECT sd.domain
         FROM site_domains sd
         WHERE sd.site_id = s.id
         ORDER BY sd.position ASC, sd.id ASC
         LIMIT 1
       ), ''),
       s.name,
       s.public_key,
       s.track_country,
       s.created_at,
       s.updated_at
FROM sites s;

DROP TABLE `sites`;
ALTER TABLE `sites_new` RENAME TO `sites`;

CREATE UNIQUE INDEX `sites_domain` ON `sites` (`domain`);
CREATE UNIQUE INDEX `sites_public_key` ON `sites` (`public_key`);

-- rebuild site_domains table without position column
CREATE TABLE `site_domains_new` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `domain` varchar NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);

INSERT INTO `site_domains_new` (`id`, `site_id`, `domain`, `created_at`, `updated_at`)
SELECT `id`, `site_id`, `domain`, `created_at`, `updated_at` FROM `site_domains`;

DROP TABLE `site_domains`;
ALTER TABLE `site_domains_new` RENAME TO `site_domains`;

CREATE UNIQUE INDEX `site_domains_domain` ON `site_domains` (`domain`);

PRAGMA foreign_keys=on;
