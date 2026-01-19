-- create "site_blocked_ips" table
CREATE TABLE "site_blocked_ips" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "ip" character varying NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "site_blocked_ips_site_id_ip_key" UNIQUE ("site_id", "ip"),
  CONSTRAINT "site_blocked_ips_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "site_blocked_countries" table
CREATE TABLE "site_blocked_countries" (
  "id" bigserial NOT NULL,
  "site_id" bigint NOT NULL,
  "country_code" character varying NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "site_blocked_countries_site_id_country_code_key" UNIQUE ("site_id", "country_code"),
  CONSTRAINT "site_blocked_countries_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "sites" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
