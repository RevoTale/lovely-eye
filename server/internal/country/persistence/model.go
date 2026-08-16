package persistence

import "github.com/uptrace/bun"

type Country struct {
	bun.BaseModel `bun:"table:countries,alias:co"`

	Code string `bun:"code,pk,type:varchar(2)"`
	Name string `bun:"name,notnull,type:varchar(128)"`
}
