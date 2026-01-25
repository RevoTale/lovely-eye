-- remove unused event type column
PRAGMA foreign_keys=OFF;

CREATE TABLE `events_new` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `session_id` integer NOT NULL,
  `time` integer NOT NULL,
  `hour` integer NOT NULL,
  `day` integer NOT NULL,
  `path` varchar NOT NULL,
  `definition_id` integer NULL,
  CONSTRAINT `0` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`definition_id`) REFERENCES `event_definitions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);

INSERT INTO `events_new` (`id`, `session_id`, `time`, `hour`, `day`, `path`, `definition_id`)
  SELECT `id`, `session_id`, `time`, `hour`, `day`, `path`, `definition_id` FROM `events`;

DROP TABLE `events`;
ALTER TABLE `events_new` RENAME TO `events`;

PRAGMA foreign_keys=ON;
