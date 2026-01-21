-- reverse: drop "page_views" table
CREATE TABLE "public"."page_views" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "session_id" bigint NULL,
  "visitor_id" character varying NOT NULL,
  "path" character varying NOT NULL,
  "title" character varying NULL,
  "referrer" character varying NULL,
  "duration" bigint NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "page_views_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "public"."sessions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "page_views_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "public"."sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- reverse: drop "daily_stats" table
CREATE TABLE "public"."daily_stats" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "date" timestamptz NOT NULL,
  "visitors" bigint NULL DEFAULT 0,
  "page_views" bigint NULL DEFAULT 0,
  "sessions" bigint NULL DEFAULT 0,
  "bounce_rate" double precision NULL DEFAULT 0,
  "avg_duration" double precision NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "daily_stats_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "public"."sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- reverse: modify "sessions" table
ALTER TABLE "public"."sessions" DROP CONSTRAINT "sessions_client_id_fkey", DROP COLUMN "page_view_count", DROP COLUMN "exit_path", DROP COLUMN "exit_day", DROP COLUMN "exit_hour", DROP COLUMN "exit_time", DROP COLUMN "enter_path", DROP COLUMN "enter_day", DROP COLUMN "enter_hour", DROP COLUMN "enter_time", DROP COLUMN "client_id", ADD COLUMN "event_only" boolean NOT NULL DEFAULT false, ADD COLUMN "is_bounce" boolean NULL DEFAULT true, ALTER COLUMN "duration" DROP NOT NULL, ADD COLUMN "page_views" bigint NULL DEFAULT 0, ADD COLUMN "screen_size" character varying NULL, ADD COLUMN "os" character varying NULL, ADD COLUMN "browser" character varying NULL, ADD COLUMN "device" character varying NULL, ADD COLUMN "city" character varying NULL, ADD COLUMN "country" character varying NULL, ALTER COLUMN "utm_campaign" TYPE character varying, ALTER COLUMN "utm_medium" TYPE character varying, ALTER COLUMN "utm_source" TYPE character varying, ALTER COLUMN "referrer" TYPE character varying, ADD COLUMN "exit_page" character varying NULL, ADD COLUMN "entry_page" character varying NULL, ADD COLUMN "last_seen_at" timestamptz NOT NULL, ADD COLUMN "started_at" timestamptz NOT NULL, ADD COLUMN "visitor_id" character varying NOT NULL;
-- reverse: create "event_data" table
DROP TABLE "public"."event_data";
-- reverse: modify "event_definition_fields" table
ALTER TABLE "public"."event_definition_fields" ALTER COLUMN "type" TYPE character varying, ALTER COLUMN "key" TYPE character varying;
-- reverse: modify "events" table
ALTER TABLE "public"."events" DROP CONSTRAINT "events_definition_id_fkey", DROP COLUMN "definition_id", DROP COLUMN "type", DROP COLUMN "day", DROP COLUMN "hour", DROP COLUMN "time", ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, ADD COLUMN "properties" character varying NULL, ALTER COLUMN "path" TYPE character varying, ALTER COLUMN "path" DROP NOT NULL, ALTER COLUMN "name" TYPE character varying, ADD COLUMN "visitor_id" character varying NOT NULL, ALTER COLUMN "session_id" DROP NOT NULL, ADD COLUMN "site_id" bigint NOT NULL;
-- reverse: create "clients" table
DROP TABLE "public"."clients";
-- reverse: modify "site_domains" table
ALTER TABLE "public"."site_domains" ALTER COLUMN "position" TYPE integer;
-- reverse: modify "site_blocked_ips" table
ALTER TABLE "public"."site_blocked_ips" ADD CONSTRAINT "site_blocked_ips_site_id_ip_key" UNIQUE ("site_id", "ip");
-- reverse: modify "site_blocked_countries" table
ALTER TABLE "public"."site_blocked_countries" ADD CONSTRAINT "site_blocked_countries_site_id_country_code_key" UNIQUE ("site_id", "country_code");
