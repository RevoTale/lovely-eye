-- remove unused event name and type columns
ALTER TABLE "public"."events" DROP COLUMN "name";
ALTER TABLE "public"."events" DROP COLUMN "type";
