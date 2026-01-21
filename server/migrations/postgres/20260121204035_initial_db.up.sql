-- create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "username" character varying NOT NULL,
  "password_hash" character varying NOT NULL,
  "role" character varying NOT NULL DEFAULT 'user',
  "email" character varying NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_username_key" UNIQUE ("username")
);
-- create "sites" table
CREATE TABLE "public"."sites" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "name" character varying NOT NULL,
  "public_key" character varying NOT NULL,
  "track_country" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "sites_public_key_key" UNIQUE ("public_key"),
  CONSTRAINT "sites_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
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
-- create "sessions" table
CREATE TABLE "public"."sessions" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "client_id" bigint NOT NULL,
  "enter_time" bigint NOT NULL,
  "enter_hour" bigint NOT NULL,
  "enter_day" bigint NOT NULL,
  "enter_path" character varying(2048) NOT NULL,
  "exit_time" bigint NOT NULL,
  "exit_hour" bigint NOT NULL,
  "exit_day" bigint NOT NULL,
  "exit_path" character varying(2048) NOT NULL,
  "referrer" character varying(2048) NULL,
  "utm_source" character varying(128) NULL,
  "utm_medium" character varying(128) NULL,
  "utm_campaign" character varying(256) NULL,
  "duration" bigint NOT NULL DEFAULT 0,
  "page_view_count" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "sessions_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "sessions_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "public"."sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "events" table
CREATE TABLE "public"."events" (
  "id" bigserial NOT NULL,
  "session_id" bigint NOT NULL,
  "time" bigint NOT NULL,
  "hour" bigint NOT NULL,
  "day" bigint NOT NULL,
  "path" character varying(2048) NOT NULL,
  "name" character varying(256) NOT NULL,
  "type" smallint NOT NULL,
  "definition_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "events_definition_id_fkey" FOREIGN KEY ("definition_id") REFERENCES "public"."event_definitions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "events_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "public"."sessions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "event_definition_fields" table
CREATE TABLE "public"."event_definition_fields" (
  "id" bigserial NOT NULL,
  "event_definition_id" bigint NOT NULL,
  "key" character varying(64) NOT NULL,
  "type" smallint NOT NULL,
  "required" boolean NOT NULL DEFAULT false,
  "max_length" bigint NOT NULL DEFAULT 500,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "event_definition_fields_event_definition_id_fkey" FOREIGN KEY ("event_definition_id") REFERENCES "public"."event_definitions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
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
-- create "site_blocked_countries" table
CREATE TABLE "public"."site_blocked_countries" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "country_code" character varying NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "site_blocked_countries_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "public"."sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "site_blocked_ips" table
CREATE TABLE "public"."site_blocked_ips" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "ip" character varying NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "site_blocked_ips_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "public"."sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "site_domains" table
CREATE TABLE "public"."site_domains" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "domain" character varying NOT NULL,
  "position" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "site_domains_domain_key" UNIQUE ("domain"),
  CONSTRAINT "site_domains_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "public"."sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
