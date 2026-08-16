package persistence

import (
	"time"

	"github.com/uptrace/bun"
)

type Definition struct {
	bun.BaseModel `bun:"table:event_definitions,alias:ed"`

	ID        int64     `bun:"id,pk,autoincrement"`
	SiteID    int64     `bun:"site_id,notnull,unique:event_definitions_site_id_name"`
	Name      string    `bun:"name,notnull,unique:event_definitions_site_id_name"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`

	Fields []*Field `bun:"rel:has-many,join:id=event_definition_id"`
}

type FieldType int8

const (
	FieldTypeString FieldType = iota
	FieldTypeInt
	FieldTypeFloat
	FieldTypeBool
)

type Field struct {
	bun.BaseModel `bun:"table:event_definition_fields,alias:edf"`

	ID                int64     `bun:"id,pk,autoincrement"`
	EventDefinitionID int64     `bun:"event_definition_id,notnull"`
	Key               string    `bun:"key,notnull,type:varchar(64)"`
	Type              FieldType `bun:"type,notnull"`
	Required          bool      `bun:"required,notnull,default:false"`
	MaxLength         int       `bun:"max_length,notnull,default:500"`
	CreatedAt         time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt         time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
