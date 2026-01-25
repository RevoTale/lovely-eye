-- restore event name and type columns for rollbacks
ALTER TABLE "public"."events" ADD COLUMN "name" character varying(256) NOT NULL DEFAULT '';
ALTER TABLE "public"."events" ADD COLUMN "type" smallint NOT NULL DEFAULT 0;
UPDATE "public"."events" SET "name" = "path" WHERE "definition_id" IS NULL;
UPDATE "public"."events" SET "name" = ed.name FROM "public"."event_definitions" ed WHERE "events"."definition_id" = ed.id;
UPDATE "public"."events" SET "type" = 1 WHERE "definition_id" IS NOT NULL;
ALTER TABLE "public"."events" ALTER COLUMN "name" DROP DEFAULT;
ALTER TABLE "public"."events" ALTER COLUMN "type" DROP DEFAULT;
