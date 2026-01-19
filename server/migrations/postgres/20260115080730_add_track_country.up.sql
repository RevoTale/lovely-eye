-- modify "sites" table
ALTER TABLE "public"."sites" ADD COLUMN "track_country" boolean NOT NULL DEFAULT false;
