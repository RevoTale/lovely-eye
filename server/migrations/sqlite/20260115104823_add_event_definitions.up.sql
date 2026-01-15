-- create "event_definitions" table
CREATE TABLE `event_definitions` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `site_id` integer NOT NULL,
  `name` varchar NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "event_definition_fields" table
CREATE TABLE `event_definition_fields` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `event_definition_id` integer NOT NULL,
  `key` varchar NOT NULL,
  `type` varchar NOT NULL,
  `required` boolean NOT NULL DEFAULT false,
  `max_length` integer NOT NULL DEFAULT 500,
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp),
  CONSTRAINT `0` FOREIGN KEY (`event_definition_id`) REFERENCES `event_definitions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
