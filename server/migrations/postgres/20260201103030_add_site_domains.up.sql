-- create "site_domains" table
CREATE TABLE "site_domains" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "domain" character varying NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "site_domains_domain_key" UNIQUE ("domain"),
  CONSTRAINT "site_domains_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- backfill existing site domains
INSERT INTO "site_domains" ("site_id", "domain", "created_at", "updated_at")
SELECT "id", "domain", "created_at", "updated_at" FROM "sites";
