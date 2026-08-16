# Upgrading Lovely Eye

This guide is for operators updating a self-hosted Lovely Eye container. Read the release notes for
every version between the version you run and the version you plan to deploy. Major releases may
change the tracker, configuration, or other externally observable contracts.

## Choose an image reference

Lovely Eye publishes these multi-architecture tags for `linux/amd64` and `linux/arm64`:

| Reference | Movement | Use |
| --- | --- | --- |
| `v2.0.0` | Never intentionally moved | Recommended for controlled production upgrades |
| `v2.0` | Receives patch releases in the `2.0.x` line | Patch channel |
| `v2` | Receives backward-compatible `2.x` releases | Major channel |
| `latest` | Moves across major versions | Evaluation only; do not automate production upgrades from it |
| `@sha256:...` | Immutable registry content | Maximum reproducibility and auditability |

Docker tags are mutable registry references; a digest identifies exact content. Use an exact version
or the digest published in the corresponding GitHub Release when you need deterministic rollout and
rollback. See Docker's [digest documentation](https://docs.docker.com/reference/cli/docker/image/pull/#pull-an-image-by-digest-immutable-identifier).

## Before every upgrade

1. Read the target [GitHub Release](https://github.com/RevoTale/lovely-eye/releases) and every
   intervening release note.
2. Record the currently deployed image reference and digest:

   ```bash
   docker compose images
   docker image inspect ghcr.io/revotale/lovely-eye:<current-version> --format '{{index .RepoDigests 0}}'
   ```

3. Back up the database before the new container starts. Lovely Eye applies pending migrations
   automatically during startup, before serving traffic.
4. Preserve `JWT_SECRET` and `ANALYTICS_IDENTITY_SECRET`. Changing either value has authentication or
   analytics-identity consequences unrelated to the image upgrade.
5. Check the target release for configuration validation changes, proxy requirements, tracker changes,
   and database disk-space requirements.

## Back up the database

Stop the Lovely Eye application while taking an SQLite file backup so the copied database and its
journal files form one consistent snapshot.

```bash
docker compose stop lovely-eye
mkdir -p backups/pre-upgrade
docker compose cp lovely-eye:/app/data/. backups/pre-upgrade/
docker compose start lovely-eye
```

If `/app/data` is a host bind mount, copy or snapshot that host directory instead. Keep the file mode
and ownership metadata when your backup system supports it.

For PostgreSQL, stop the application and take a logical or infrastructure-native snapshot while the
database remains running. Adapt the database name and user to your Compose configuration:

```bash
docker compose stop lovely-eye
mkdir -p backups
docker compose exec -T lovely-eye-db pg_dump -U lovely -d lovely_eye -Fc > backups/lovely-eye-pre-upgrade.dump
docker compose start lovely-eye
```

Test restoration periodically. A backup that has never been restored is not a verified rollback path.

## Upgrade with Docker Compose

1. Change the image to the exact target version. Do not change the database image in the same rollout
   unless the release notes explicitly require it.

   ```yaml
   services:
     lovely-eye:
       image: ghcr.io/revotale/lovely-eye:v2.0.0
   ```

2. Pull and replace only the application service:

   ```bash
   docker compose pull lovely-eye
   docker compose up -d --no-deps lovely-eye
   ```

3. Verify startup, migration state, and health:

   ```bash
   docker compose logs --tail=100 lovely-eye
   curl --fail --silent --show-error http://localhost:8080/health
   ```

   Startup logs must report either the applied migration group or `No new migrations to run` before
   the server begins accepting traffic.

4. Sign in and smoke-test site selection, analytics, settings, and one real tracker request. A healthy
   endpoint proves that the process, dashboard files, and database are available; it does not prove
   that reverse-proxy client identity or collection traffic is configured correctly.

## Roll back

Stop the rollout when startup fails, health remains unhealthy, collection errors rise, analytics
identity becomes implausible, or a critical admin flow fails.

For a code-only rollback, restore the previous exact image reference and recreate the application
container:

```bash
docker compose stop lovely-eye
# Edit compose.yaml and restore the previous exact image tag or digest.
docker compose up -d --no-deps lovely-eye
```

Only use code-only rollback when the target release explicitly states that its migrated schema remains
compatible with the previous image. Otherwise, stop both application versions, restore the pre-upgrade
database backup, restore the previous image reference, and start the previous version. Restoring a
backup discards writes accepted after that backup; decide which rollback mode is appropriate before
the rollout starts.

## Version-specific guides

- [Version 2.0.0](docs/releases/v2.0.0.md)
