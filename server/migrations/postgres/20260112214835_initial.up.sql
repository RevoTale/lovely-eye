-- create "users" table
CREATE TABLE "users" (
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
CREATE TABLE "sites" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "domain" character varying NOT NULL,
  "name" character varying NOT NULL,
  "public_key" character varying NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "sites_domain_key" UNIQUE ("domain"),
  CONSTRAINT "sites_public_key_key" UNIQUE ("public_key"),
  CONSTRAINT "sites_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "daily_stats" table
CREATE TABLE "daily_stats" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "date" timestamptz NOT NULL,
  "visitors" bigint NULL DEFAULT 0,
  "page_views" bigint NULL DEFAULT 0,
  "sessions" bigint NULL DEFAULT 0,
  "bounce_rate" double precision NULL DEFAULT 0,
  "avg_duration" double precision NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "daily_stats_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "sessions" table
CREATE TABLE "sessions" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "visitor_id" character varying NOT NULL,
  "started_at" timestamptz NOT NULL,
  "last_seen_at" timestamptz NOT NULL,
  "entry_page" character varying NULL,
  "exit_page" character varying NULL,
  "referrer" character varying NULL,
  "utm_source" character varying NULL,
  "utm_medium" character varying NULL,
  "utm_campaign" character varying NULL,
  "country" character varying NULL,
  "city" character varying NULL,
  "device" character varying NULL,
  "browser" character varying NULL,
  "os" character varying NULL,
  "screen_size" character varying NULL,
  "page_views" bigint NULL DEFAULT 0,
  "duration" bigint NULL DEFAULT 0,
  "is_bounce" boolean NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "sessions_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "events" table
CREATE TABLE "events" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "session_id" bigint NULL,
  "visitor_id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "path" character varying NULL,
  "properties" character varying NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "events_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "events_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "page_views" table
CREATE TABLE "page_views" (
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
  CONSTRAINT "page_views_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "page_views_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
