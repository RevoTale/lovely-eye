# Server performance and memory

This runbook defines how maintainers measure Lovely Eye server RAM, allocation churn, and CPU cost.
It covers the production Go process and the collection path; dashboard bundle budgets remain in
`docs/plans/architecture-normalization.md`.

## Interpret the memory numbers correctly

Lovely Eye is not only a reverse proxy. Every accepted collection request resolves current site
configuration, validates origin and blocking rules, parses the user agent, derives visitor identity,
and reads or writes analytics state. Dashboard requests also execute GraphQL and analytical queries.

Use all three measurements when investigating memory:

- cgroup memory: the container charge, including process memory and charged file cache;
- process RSS/PSS: resident executable, mapped files, stacks, and heap;
- private anonymous memory and Go heap: memory most directly affected by Go allocations and GC.

`docker stats` may exclude inactive file cache, so it can differ from `memory.current`. A changing RSS
alone is not proof of a leak. Confirm retained Go heap or continuously growing private anonymous
memory after equivalent GC and workload cycles.

The baseline below measures the application container with SQLite. A separate PostgreSQL container's
memory is outside that cgroup, and an enabled GeoIP database adds mapped/file-backed pages. Measure
those deployment components separately before selecting a whole-stack memory limit.

## 2026-08-15 baseline

Environment: Go 1.26.6, Linux arm64, SQLite, production Alpine image, one CPU limit, 128 MiB memory
limit, default Go runtime settings. The representative database had 80 clients, 159 sessions, 668
page views, and 239 predefined events.

| Measurement | Before | Current |
| --- | ---: | ---: |
| Warm collect benchmark | 198–206 us/op | 120–124 us/op |
| Allocated bytes per collect | 95.8–96.0 KiB/op | 45.2–45.3 KiB/op |
| Allocations per collect | 900 allocs/op | 565 allocs/op |
| 2,000-request runtime GC cycles | 91 | 39 |
| 2,000-request process high-water RSS | 23.5 MiB | 21.5 MiB |
| 2,000-request final cgroup memory | 28.38 MiB | 28.04 MiB |
| 2,000-request elapsed time | 3,044 ms | 3,044 ms |

The runtime workload used eight concurrent clients and a fresh cloned database for each image. Its
elapsed time is SQLite/load-generator bound and therefore intentionally recorded as neutral. The
reduction in GC cycles and process high-water RSS confirms lower runtime pressure without claiming a
latency gain that the end-to-end sample did not show.

Idle samples used roughly 3.2–4.5 MiB of private anonymous memory. Total cgroup observations ranged
from about 8 MiB before executable pages were faulted to about 25 MiB after warm-up. Do not treat the
lower cold value as a production memory requirement.

## Preserved hot-path invariants

- The HTTP handler resolves site configuration once and passes the resolved site to analytics.
- Analytics still rechecks bot, domain, blocked-IP, and blocked-country rules before persistence.
- Site configuration is read from the database on every request. Do not add a process cache without
  a multi-instance invalidation design and measured retained-memory budget.
- The collection lookup returns domains ordered by position, blocked IPs ordered by address, and
  blocked countries ordered by code on both SQLite and PostgreSQL.
- No `GOMEMLIMIT` is set by the application. The measured live heap was about 1 MiB with a 4 MiB GC
  goal; lowering the soft limit would add GC CPU without reducing the dominant mapped/file-backed
  memory. Operators may set `GOMEMLIMIT` below a container limit after workload-specific testing.
- Go 1.27 applies the container-aware `GOMAXPROCS` default. Do not add an automatic CPU-limit
  package or hardcode `GOMAXPROCS` without a same-quota latency and throughput comparison.

## Reproduce the benchmarks

Run the complete maintained baseline inside the Dev Container with no competing project workload:

```bash
./scripts/measure-baseline.sh
```

Run only the collection allocation benchmark:

```bash
cd server
go test -run '^$' -bench '^BenchmarkAnalyticsHandlerCollectPageView$' \
  -benchmem -benchtime=2s -count=10 ./internal/transport/http/collect
```

Capture CPU and allocation profiles without leaving a test binary in the repository:

```bash
cd server
go test -o /tmp/collect.test -run '^$' \
  -bench '^BenchmarkAnalyticsHandlerCollectPageView$' -benchtime=5s \
  -cpuprofile /tmp/collect-cpu.pprof -memprofile /tmp/collect-mem.pprof \
  ./internal/transport/http/collect
go tool pprof -top /tmp/collect-cpu.pprof
go tool pprof -top -sample_index=alloc_space /tmp/collect-mem.pprof
```

Go documents profile selection and interference in its
[diagnostics guide](https://go.dev/doc/diagnostics). `testing.B.ReportAllocs` and `-benchmem` are the
standard allocation measurements described by the
[`testing` package](https://pkg.go.dev/testing#B.ReportAllocs).

## Rejected and deferred work

- Replacing `map[string]*collectBucket` with `map[string]collectBucket` was reverted. Repeated keys
  slowed from about 42 to 44 ns/op, while rotating keys grew from 138 to 220 B/op because the larger
  map buckets outweighed removal of the pointed allocation.
- The collect and authentication limiter maps are time-cleaned but not hard-capped. A hard cap would
  bound retained memory under rotating-IP traffic, but it also defines how many active identities can
  be admitted before legitimate traffic is rejected or old limits are evicted. Set that capacity only
  after the maintainer accepts a traffic envelope and overflow policy.
- Do not add `sync.Pool`, JSON replacement, a site cache, or profile-guided optimization without a new
  profile proving that the target remains material and a same-environment result outside noise.

The Go GC guide documents that `GOMEMLIMIT` is a soft runtime limit and can trade memory for up to
substantial GC CPU when configured too low:
https://go.dev/doc/gc-guide. Container-aware CPU behavior is documented at
https://go.dev/blog/container-aware-gomaxprocs.
