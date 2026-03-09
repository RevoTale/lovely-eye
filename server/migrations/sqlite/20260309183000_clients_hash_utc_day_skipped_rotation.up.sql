-- repoint duplicate-client sessions to the canonical lowest client id per (site_id, hash)
WITH duplicate_clients AS (
  SELECT c.id AS duplicate_id, canon.keep_id
  FROM `clients` AS c
  INNER JOIN (
    SELECT `site_id`, `hash`, MIN(`id`) AS keep_id
    FROM `clients`
    GROUP BY `site_id`, `hash`
    HAVING COUNT(*) > 1
  ) AS canon
    ON canon.`site_id` = c.`site_id`
   AND canon.`hash` = c.`hash`
  WHERE c.`id` <> canon.keep_id
)
UPDATE `sessions`
SET `client_id` = (
  SELECT duplicate_clients.keep_id
  FROM duplicate_clients
  WHERE duplicate_clients.duplicate_id = `sessions`.`client_id`
)
WHERE `client_id` IN (SELECT duplicate_id FROM duplicate_clients);

-- delete duplicate client rows now that sessions point at the canonical client
WITH duplicate_clients AS (
  SELECT c.id AS duplicate_id
  FROM `clients` AS c
  INNER JOIN (
    SELECT `site_id`, `hash`, MIN(`id`) AS keep_id
    FROM `clients`
    GROUP BY `site_id`, `hash`
    HAVING COUNT(*) > 1
  ) AS canon
    ON canon.`site_id` = c.`site_id`
   AND canon.`hash` = c.`hash`
  WHERE c.`id` <> canon.keep_id
)
DELETE FROM `clients`
WHERE `id` IN (SELECT duplicate_id FROM duplicate_clients);

-- enforce one client row per site/hash and speed up active-session lookups
CREATE UNIQUE INDEX `clients_site_id_hash` ON `clients` (`site_id`, `hash`);
CREATE INDEX `sessions_site_id_client_id_exit_time` ON `sessions` (`site_id`, `client_id`, `exit_time` DESC);
