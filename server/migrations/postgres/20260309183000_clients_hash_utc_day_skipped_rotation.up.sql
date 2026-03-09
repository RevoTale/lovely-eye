-- repoint duplicate-client sessions to the canonical lowest client id per (site_id, hash)
WITH duplicate_clients AS (
  SELECT c.id AS duplicate_id, canon.keep_id
  FROM "public"."clients" AS c
  INNER JOIN (
    SELECT "site_id", "hash", MIN("id") AS keep_id
    FROM "public"."clients"
    GROUP BY "site_id", "hash"
    HAVING COUNT(*) > 1
  ) AS canon
    ON canon."site_id" = c."site_id"
   AND canon."hash" = c."hash"
  WHERE c."id" <> canon.keep_id
)
UPDATE "public"."sessions" AS s
SET "client_id" = duplicate_clients.keep_id
FROM duplicate_clients
WHERE s."client_id" = duplicate_clients.duplicate_id;

-- delete duplicate client rows now that sessions point at the canonical client
WITH duplicate_clients AS (
  SELECT c.id AS duplicate_id
  FROM "public"."clients" AS c
  INNER JOIN (
    SELECT "site_id", "hash", MIN("id") AS keep_id
    FROM "public"."clients"
    GROUP BY "site_id", "hash"
    HAVING COUNT(*) > 1
  ) AS canon
    ON canon."site_id" = c."site_id"
   AND canon."hash" = c."hash"
  WHERE c."id" <> canon.keep_id
)
DELETE FROM "public"."clients" AS c
USING duplicate_clients
WHERE c."id" = duplicate_clients.duplicate_id;

-- enforce one client row per site/hash and speed up active-session lookups
ALTER TABLE "public"."clients" ADD CONSTRAINT "clients_site_id_hash" UNIQUE ("site_id", "hash");
CREATE INDEX "sessions_site_id_client_id_exit_time" ON "public"."sessions" ("site_id", "client_id", "exit_time" DESC);
