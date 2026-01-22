-- reverse: modify "site_domains" table
ALTER TABLE "public"."site_domains" DROP CONSTRAINT "site_domains_site_id_domain", ADD CONSTRAINT "site_domains_domain_key" UNIQUE ("domain");
-- reverse: modify "event_definitions" table
ALTER TABLE "public"."event_definitions" DROP CONSTRAINT "event_definitions_site_id_name";
