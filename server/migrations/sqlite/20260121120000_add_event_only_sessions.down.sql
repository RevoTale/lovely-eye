PRAGMA foreign_keys=off;

CREATE TABLE `sessions_new` (
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

INSERT INTO `sessions_new` (
  `id`, `site_id`, `visitor_id`, `started_at`, `last_seen_at`, `entry_page`, `exit_page`,
  `referrer`, `utm_source`, `utm_medium`, `utm_campaign`, `country`, `city`, `device`,
  `browser`, `os`, `screen_size`, `page_views`, `duration`, `is_bounce`
)
SELECT
  `id`, `site_id`, `visitor_id`, `started_at`, `last_seen_at`, `entry_page`, `exit_page`,
  `referrer`, `utm_source`, `utm_medium`, `utm_campaign`, `country`, `city`, `device`,
  `browser`, `os`, `screen_size`, `page_views`, `duration`, `is_bounce`
FROM `sessions`;

DROP TABLE `sessions`;
ALTER TABLE `sessions_new` RENAME TO `sessions`;

PRAGMA foreign_keys=on;
