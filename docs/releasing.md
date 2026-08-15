# Release runbook

This runbook defines how maintainers publish Lovely Eye source and multi-architecture container
releases without hiding operator action behind generated changelog entries.

## Release contract

Lovely Eye follows [Semantic Versioning 2.0.0](https://semver.org/). The public release surface
includes the tracker script and collect payload, environment variables and validation, database
migrations, container filesystem/runtime behavior, public HTTP paths, and documented Docker tags.
The bundled dashboard GraphQL schema remains internal unless a separate external API contract is
accepted later.

- Patch: backward-compatible bug or security correction.
- Minor: backward-compatible capability or operational improvement.
- Major: an operator, tracker integrator, or documented consumer must change something.

Release Please parses Conventional Commits and creates the version PR, tag, changelog, and a draft
GitHub Release. Use `!` or a `BREAKING CHANGE:` footer only when the public release surface is
actually incompatible. Generated notes are an index of changes, not a substitute for an upgrade
guide.

## Required release artifacts

Every release has:

- a SemVer Git tag and GitHub Release;
- exact, minor-channel, major-channel, and `latest` Docker tags;
- an immutable multi-architecture image digest appended to the GitHub Release;
- a reviewed `CHANGELOG.md` entry;
- passing CI, production build, and migration checks.

A major release, or any minor release with operator action, also requires
`docs/releases/<tag>.md`. The document must state:

- why the version has that SemVer level;
- action required before deployment;
- externally visible behavior and metric-definition changes;
- configuration additions, removals, validation, and default changes;
- database migration, backup, rollback, and expected downtime behavior;
- tracker/API compatibility;
- supported image platforms and databases;
- known issues and intentionally deferred risks.

When that file exists, `.github/workflows/release.yml` uses it as the GitHub Release body after the
container has been pushed. Otherwise, the generated Release Please body is retained. In both cases,
the workflow appends the published image tag and immutable digest. The release remains a draft until
tests pass and the multi-architecture image is available; the final workflow step then publishes it.
An image-build failure therefore leaves a maintainer-visible draft instead of a public partial
release.

## Before merging the release PR

1. Confirm the proposed version matches the public impact, not the diff size.
2. Rewrite noisy or duplicated generated changelog entries into user-observable language.
3. For a major/material release, confirm the versioned release note filename matches the proposed tag.
4. Verify `UPGRADING.md`, README examples, environment tables, tracker documentation, migration docs,
   and security notes describe the release candidate.
5. Run the complete repository test matrix, both database migration cycles, runtime-base-path matrix,
   production image health check, dependency audits, and performance guardrails.
6. Confirm rollback commands and the previous exact image tag are known before publishing.
7. Merge the release PR only during a window in which someone can monitor the release and answer
   operator reports.

Release Please recommends squash-merging feature PRs so temporary implementation commits do not
become misleading release entries. Its official workflow is documented at
https://github.com/googleapis/release-please.

## Docker tag policy

For release `v2.3.4`, the workflow publishes:

- `v2.3.4`: exact release; never intentionally overwritten;
- `v2.3`: moves only for patch releases in that minor line;
- `v2`: moves for backward-compatible releases in that major line;
- `latest`: moves to the newest stable release, including a new major.

Tell production operators to pin an exact version or digest. Major and minor aliases are convenience
channels with controlled movement; `latest` is not a compatibility promise. Docker documents that
tags can move while digests are immutable:
https://docs.docker.com/reference/cli/docker/image/pull/#pull-an-image-by-digest-immutable-identifier.

The release action uses Docker Metadata Action's SemVer patterns and OCI labels. The tag behavior and
generated metadata are documented at https://github.com/docker/metadata-action#semver. OCI-defined
source, version, revision, documentation, and license keys are specified at
https://github.com/opencontainers/image-spec/blob/main/annotations.md.

## Post-release verification

1. Wait for the release workflow, tests, and both architecture image manifests to finish. Confirm
   that the release changed from draft to public only after the image was published.
2. Confirm the GitHub Release contains the intended user notes and immutable digest.
3. Pull the exact tag on an empty test host and perform the documented upgrade from the previous
   stable release with representative persistent data.
4. Check `/health`, startup and migration logs, sign-in, site selection, analytics, settings, one
   real tracker request, and proxy-derived client identity.
5. Compare error rate, response latency, collection volume, CPU, memory, and database growth with the
   pre-release baseline.
6. If a rollback trigger is reached, stop promotion of moving tags in downstream deployments and
   follow `UPGRADING.md` using the previous exact tag or digest.

GitHub documents releases as the durable place for release notes and attached artifacts:
https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases.
