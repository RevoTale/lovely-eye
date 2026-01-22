-- modify "event_definitions" table
ALTER TABLE "public"."event_definitions" ADD CONSTRAINT "event_definitions_site_id_name" UNIQUE ("site_id", "name");
-- modify "site_domains" table
ALTER TABLE "public"."site_domains" DROP CONSTRAINT "site_domains_domain_key", ADD CONSTRAINT "site_domains_site_id_domain" UNIQUE ("site_id", "domain");
