-- add position to site_domains for ordered domains
ALTER TABLE "site_domains" ADD COLUMN "position" integer NOT NULL DEFAULT 0;

WITH ranked AS (
  SELECT id, ROW_NUMBER() OVER (PARTITION BY site_id ORDER BY id) - 1 AS position
  FROM site_domains
)
UPDATE site_domains sd
SET position = ranked.position
FROM ranked
WHERE sd.id = ranked.id;

ALTER TABLE "sites" DROP CONSTRAINT "sites_domain_key";
ALTER TABLE "sites" DROP COLUMN "domain";
