# ADR-0003: Use a feature-oriented Go modular monolith

## Status

Accepted

## Date

2026-08-06

## Context

The Go server currently groups code into generic models, repositories, services, handlers, GraphQL, and command packages. These layers expose persistence result types and aliases across boundaries, mix HTTP cookie behavior with authentication application logic, and make a feature change span unrelated package buckets.

Lovely Eye is one deployable application. It does not need distributed-service boundaries or ceremonial Clean Architecture layers, but it does need explicit ownership, contained change blast radius, data accuracy, and idiomatic composition.

This ADR records decisions D-003, D-005, D-016, and D-017.

## Decision

Use a feature-oriented modular monolith:

```text
server/internal/
  app/
  auth/
  site/
  analytics/
  event/
  geoip/
  transport/
  platform/
  seed/
```

- Feature packages own application behavior, feature types, and narrow consumer-side interfaces.
- `transport` owns HTTP, GraphQL, dashboard-serving, cookie, and transport-model mapping.
- `platform` owns database and other infrastructure adapters.
- `app` is the explicit composition root and owns startup and lifecycle cleanup.
- `cmd` packages only choose a command, load validated configuration, construct the application, and run it.
- Constructor composition is explicit and idiomatic. Wiring may be handwritten or deterministically generated when graph complexity justifies it.
- Globals, service locators, hidden runtime reflection, and runtime DI containers are prohibited.
- Persistence row structs and Bun/database tags remain inside persistence adapters.
- A feature introduces a separate type only when it protects a business invariant or represents a materially different contract. Mechanical one-to-one domain copies of every row are prohibited.
- Persistence types never leak into feature, GraphQL, or HTTP contracts.
- Deterministic code generation is allowed when it has a clear handwritten source of truth, reproducible project commands, isolated output, and a stale-artifact check.

## Alternatives considered

### Formal Clean Architecture layers

Rejected because generic use-case, entity, gateway, and presenter layers would add ceremony without an independent deployment or domain need.

### Keep persistence models as universal application models

Rejected because ORM tags and database result shapes would continue leaking through feature and transport contracts.

### Map every database row to a duplicate domain model

Rejected because mapping without an invariant or shape difference adds code without creating a meaningful boundary.

### Runtime dependency-injection container

Rejected because it hides construction and lifecycle errors that ordinary Go constructors can express directly.

## Consequences

- Backend migration proceeds one feature at a time while tests remain green.
- Interfaces live with consumers and remain narrow.
- Transport and persistence adapters perform explicit mapping only where justified.
- Generic legacy packages and transitional aliases must be deleted by the end of normalization.
- Package-direction checks and focused tests become part of the quality gate.
