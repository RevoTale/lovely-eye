-- create "event_definitions" table
CREATE TABLE "public"."event_definitions" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "name" character varying NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "event_definitions_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "public"."sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "event_definition_fields" table
CREATE TABLE "public"."event_definition_fields" (
  "id" bigserial NOT NULL,
  "event_definition_id" bigint NOT NULL,
  "key" character varying NOT NULL,
  "type" character varying NOT NULL,
  "required" boolean NOT NULL DEFAULT false,
  "max_length" bigint NOT NULL DEFAULT 500,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "event_definition_fields_event_definition_id_fkey" FOREIGN KEY ("event_definition_id") REFERENCES "public"."event_definitions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
