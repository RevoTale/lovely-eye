-- modify "site_blocked_countries" table
ALTER TABLE "public"."site_blocked_countries" DROP CONSTRAINT "site_blocked_countries_site_id_country_code_key";
-- modify "site_blocked_ips" table
ALTER TABLE "public"."site_blocked_ips" DROP CONSTRAINT "site_blocked_ips_site_id_ip_key";
-- modify "site_domains" table
ALTER TABLE "public"."site_domains" ALTER COLUMN "position" TYPE bigint;

-- create "clients" table
CREATE TABLE "public"."clients" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "hash" character varying(64) NOT NULL,
  "country" character varying(2) NULL,
  "device" character varying(10) NULL,
  "browser" character varying(32) NULL,
  "os" character varying(32) NULL,
  "screen_size" character varying(16) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "clients_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "public"."sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- === DATA MIGRATION STEP 1: Create clients from page_views ===
INSERT INTO "public"."clients" ("site_id", "hash")
SELECT DISTINCT pv."site_id", pv."visitor_id" AS "hash"
FROM "public"."page_views" pv
ON CONFLICT DO NOTHING;

-- === SESSIONS TABLE MIGRATION ===
-- Add new columns to sessions
ALTER TABLE "public"."sessions"
  ADD COLUMN "client_id" bigint,
  ADD COLUMN "enter_time" bigint,
  ADD COLUMN "enter_hour" bigint,
  ADD COLUMN "enter_day" bigint,
  ADD COLUMN "enter_path" character varying(2048),
  ADD COLUMN "exit_time" bigint,
  ADD COLUMN "exit_hour" bigint,
  ADD COLUMN "exit_day" bigint,
  ADD COLUMN "exit_path" character varying(2048),
  ADD COLUMN "page_view_count" bigint DEFAULT 0;

-- === DATA MIGRATION STEP 2: Populate sessions with data from page_views ===
UPDATE "public"."sessions" s SET
  "client_id" = (SELECT c."id" FROM "public"."clients" c, "public"."page_views" pv WHERE c."hash" = pv."visitor_id" AND pv."session_id" = s."id" LIMIT 1),
  "enter_time" = (SELECT EXTRACT(epoch FROM pv."created_at")::bigint FROM "public"."page_views" pv WHERE pv."session_id" = s."id" ORDER BY pv."created_at" ASC LIMIT 1),
  "enter_hour" = (SELECT EXTRACT(epoch FROM pv."created_at")::bigint / 3600 FROM "public"."page_views" pv WHERE pv."session_id" = s."id" ORDER BY pv."created_at" ASC LIMIT 1),
  "enter_day" = (SELECT EXTRACT(epoch FROM pv."created_at")::bigint / 86400 FROM "public"."page_views" pv WHERE pv."session_id" = s."id" ORDER BY pv."created_at" ASC LIMIT 1),
  "enter_path" = (SELECT pv."path" FROM "public"."page_views" pv WHERE pv."session_id" = s."id" ORDER BY pv."created_at" ASC LIMIT 1),
  "exit_time" = (SELECT EXTRACT(epoch FROM pv."created_at")::bigint FROM "public"."page_views" pv WHERE pv."session_id" = s."id" ORDER BY pv."created_at" DESC LIMIT 1),
  "exit_hour" = (SELECT EXTRACT(epoch FROM pv."created_at")::bigint / 3600 FROM "public"."page_views" pv WHERE pv."session_id" = s."id" ORDER BY pv."created_at" DESC LIMIT 1),
  "exit_day" = (SELECT EXTRACT(epoch FROM pv."created_at")::bigint / 86400 FROM "public"."page_views" pv WHERE pv."session_id" = s."id" ORDER BY pv."created_at" DESC LIMIT 1),
  "exit_path" = (SELECT pv."path" FROM "public"."page_views" pv WHERE pv."session_id" = s."id" ORDER BY pv."created_at" DESC LIMIT 1),
  "page_view_count" = (SELECT COUNT(*) FROM "public"."page_views" pv WHERE pv."session_id" = s."id")
WHERE EXISTS (SELECT 1 FROM "public"."page_views" pv WHERE pv."session_id" = s."id");

-- Make new columns NOT NULL and add constraints
ALTER TABLE "public"."sessions"
  ALTER COLUMN "client_id" SET NOT NULL,
  ALTER COLUMN "enter_time" SET NOT NULL,
  ALTER COLUMN "enter_hour" SET NOT NULL,
  ALTER COLUMN "enter_day" SET NOT NULL,
  ALTER COLUMN "enter_path" SET NOT NULL,
  ALTER COLUMN "exit_time" SET NOT NULL,
  ALTER COLUMN "exit_hour" SET NOT NULL,
  ALTER COLUMN "exit_day" SET NOT NULL,
  ALTER COLUMN "exit_path" SET NOT NULL,
  ALTER COLUMN "page_view_count" SET NOT NULL,
  ALTER COLUMN "duration" SET NOT NULL,
  ALTER COLUMN "referrer" TYPE character varying(2048),
  ALTER COLUMN "utm_source" TYPE character varying(128),
  ALTER COLUMN "utm_medium" TYPE character varying(128),
  ALTER COLUMN "utm_campaign" TYPE character varying(256),
  ADD CONSTRAINT "sessions_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;

-- Drop old columns
ALTER TABLE "public"."sessions"
  DROP COLUMN "visitor_id",
  DROP COLUMN "started_at",
  DROP COLUMN "last_seen_at",
  DROP COLUMN "entry_page",
  DROP COLUMN "exit_page",
  DROP COLUMN "country",
  DROP COLUMN "city",
  DROP COLUMN "device",
  DROP COLUMN "browser",
  DROP COLUMN "os",
  DROP COLUMN "screen_size",
  DROP COLUMN "page_views",
  DROP COLUMN "is_bounce",
  DROP COLUMN "event_only";

-- === EVENTS TABLE MIGRATION ===
-- Add new columns to events
ALTER TABLE "public"."events"
  ADD COLUMN "time" bigint,
  ADD COLUMN "hour" bigint,
  ADD COLUMN "day" bigint,
  ADD COLUMN "type" smallint,
  ADD COLUMN "definition_id" bigint;

-- Populate time columns from existing created_at
UPDATE "public"."events" SET
  "time" = EXTRACT(epoch FROM "created_at")::bigint,
  "hour" = EXTRACT(epoch FROM "created_at")::bigint / 3600,
  "day" = EXTRACT(epoch FROM "created_at")::bigint / 86400,
  "type" = 1; -- All existing events are custom events (type=1)

-- Make columns NOT NULL
ALTER TABLE "public"."events"
  ALTER COLUMN "session_id" SET NOT NULL,
  ALTER COLUMN "time" SET NOT NULL,
  ALTER COLUMN "hour" SET NOT NULL,
  ALTER COLUMN "day" SET NOT NULL,
  ALTER COLUMN "type" SET NOT NULL,
  ALTER COLUMN "name" TYPE character varying(256),
  ALTER COLUMN "path" TYPE character varying(2048),
  ALTER COLUMN "path" SET NOT NULL,
  ADD CONSTRAINT "events_definition_id_fkey" FOREIGN KEY ("definition_id") REFERENCES "public"."event_definitions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;

-- Drop old columns
ALTER TABLE "public"."events"
  DROP COLUMN "site_id",
  DROP COLUMN "visitor_id",
  DROP COLUMN "properties",
  DROP COLUMN "created_at";

-- === DATA MIGRATION STEP 3: Migrate page_views to events ===
INSERT INTO "public"."events" ("session_id", "time", "hour", "day", "path", "name", "type", "definition_id")
SELECT
  pv."session_id",
  EXTRACT(epoch FROM pv."created_at")::bigint AS "time",
  EXTRACT(epoch FROM pv."created_at")::bigint / 3600 AS "hour",
  EXTRACT(epoch FROM pv."created_at")::bigint / 86400 AS "day",
  pv."path",
  COALESCE(pv."title", pv."path") AS "name",
  0 AS "type",
  NULL AS "definition_id"
FROM "public"."page_views" pv
WHERE pv."session_id" IS NOT NULL;

-- modify "event_definition_fields" table
ALTER TABLE "public"."event_definition_fields"
  ALTER COLUMN "key" TYPE character varying(64),
  ALTER COLUMN "type" TYPE smallint;

-- create "event_data" table
CREATE TABLE "public"."event_data" (
  "id" bigserial NOT NULL,
  "event_id" bigint NOT NULL,
  "field_id" bigint NOT NULL,
  "value" character varying(1024) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "event_data_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "public"."events" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "event_data_field_id_fkey" FOREIGN KEY ("field_id") REFERENCES "public"."event_definition_fields" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- drop "daily_stats" table
DROP TABLE "public"."daily_stats";
-- drop "page_views" table
DROP TABLE "public"."page_views";
