-- create index "event_definitions_site_id_name" to table: "event_definitions"
CREATE UNIQUE INDEX `event_definitions_site_id_name` ON `event_definitions` (`site_id`, `name`);
-- drop index "site_domains_domain" from table: "site_domains"
DROP INDEX `site_domains_domain`;
-- create index "site_domains_site_id_domain" to table: "site_domains"
CREATE UNIQUE INDEX `site_domains_site_id_domain` ON `site_domains` (`site_id`, `domain`);
