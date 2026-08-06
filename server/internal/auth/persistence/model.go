package persistence

import (
	"time"

	"github.com/uptrace/bun"
)

// User is the Bun row owned by the authentication persistence adapter.
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           int64     `bun:"id,pk,autoincrement"`
	Username     string    `bun:"username,unique,notnull"`
	PasswordHash string    `bun:"password_hash,notnull"`
	Role         string    `bun:"role,notnull,default:'user'"`
	Email        string    `bun:"email"`
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
