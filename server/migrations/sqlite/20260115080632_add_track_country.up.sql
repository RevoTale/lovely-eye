-- add column "track_country" to table: "sites"
ALTER TABLE `sites` ADD COLUMN `track_country` boolean NOT NULL DEFAULT false;
