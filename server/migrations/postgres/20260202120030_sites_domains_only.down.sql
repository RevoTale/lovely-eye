ALTER TABLE "sites" ADD COLUMN "domain" character varying;

UPDATE "sites" s
SET "domain" = COALESCE((
  SELECT sd.domain
  FROM site_domains sd
  WHERE sd.site_id = s.id
  ORDER BY sd.position ASC, sd.id ASC
  LIMIT 1
), '');

ALTER TABLE "sites" ALTER COLUMN "domain" SET NOT NULL;
ALTER TABLE "sites" ADD CONSTRAINT "sites_domain_key" UNIQUE ("domain");

ALTER TABLE "site_domains" DROP COLUMN "position";
