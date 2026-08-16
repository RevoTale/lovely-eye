package persistence

import "github.com/uptrace/bun"

// These table-only rows keep complete site deletion inside the site adapter
// without importing sibling feature persistence packages.
type ownedClient struct {
	bun.BaseModel `bun:"table:clients,alias:c"`
}

type ownedSession struct {
	bun.BaseModel `bun:"table:sessions,alias:s"`
}

type ownedEvent struct {
	bun.BaseModel `bun:"table:events,alias:e"`
}

type ownedEventData struct {
	bun.BaseModel `bun:"table:event_data,alias:evd"`
}

type ownedEventDefinition struct {
	bun.BaseModel `bun:"table:event_definitions,alias:ed"`
}

type ownedEventDefinitionField struct {
	bun.BaseModel `bun:"table:event_definition_fields,alias:edf"`
}
