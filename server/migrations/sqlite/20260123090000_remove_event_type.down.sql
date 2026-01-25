-- restore event name and type columns for rollbacks
ALTER TABLE `events` ADD COLUMN `name` varchar NOT NULL DEFAULT '';
ALTER TABLE `events` ADD COLUMN `type` integer NOT NULL DEFAULT 0;
UPDATE `events` SET `name` = `path` WHERE `definition_id` IS NULL;
UPDATE `events` SET `name` = (SELECT `name` FROM `event_definitions` WHERE `event_definitions`.`id` = `events`.`definition_id`) WHERE `definition_id` IS NOT NULL;
UPDATE `events` SET `type` = 1 WHERE `definition_id` IS NOT NULL;
