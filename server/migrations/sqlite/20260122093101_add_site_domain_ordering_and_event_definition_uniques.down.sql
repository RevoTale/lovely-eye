-- reverse: create index "site_domains_site_id_domain" to table: "site_domains"
DROP INDEX `site_domains_site_id_domain`;
-- reverse: drop index "site_domains_domain" from table: "site_domains"
CREATE UNIQUE INDEX `site_domains_domain` ON `site_domains` (`domain`);
-- reverse: create index "event_definitions_site_id_name" to table: "event_definitions"
DROP INDEX `event_definitions_site_id_name`;
