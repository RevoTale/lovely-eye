-- reverse: create index "site_domains_domain" to table: "site_domains"
DROP INDEX `site_domains_domain`;
-- reverse: create "site_domains" table
DROP TABLE `site_domains`;
-- reverse: create "site_blocked_ips" table
DROP TABLE `site_blocked_ips`;
-- reverse: create "site_blocked_countries" table
DROP TABLE `site_blocked_countries`;
-- reverse: create "event_data" table
DROP TABLE `event_data`;
-- reverse: create "event_definition_fields" table
DROP TABLE `event_definition_fields`;
-- reverse: create "events" table
DROP TABLE `events`;
-- reverse: create "event_definitions" table
DROP TABLE `event_definitions`;
-- reverse: create "sessions" table
DROP TABLE `sessions`;
-- reverse: create "clients" table
DROP TABLE `clients`;
-- reverse: create index "sites_public_key" to table: "sites"
DROP INDEX `sites_public_key`;
-- reverse: create "sites" table
DROP TABLE `sites`;
-- reverse: create index "users_username" to table: "users"
DROP INDEX `users_username`;
-- reverse: create "users" table
DROP TABLE `users`;
