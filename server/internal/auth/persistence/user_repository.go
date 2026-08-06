package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

var _ auth.UserStore = (*Repository)(nil)

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *Repository) CreateForRegistration(
	ctx context.Context,
	user *auth.StoredUser,
	allowRegistration bool,
) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := r.lockUsersForBootstrap(ctx, tx); err != nil {
			return err
		}

		hasUsers, err := anyUsers(ctx, tx)
		if err != nil {
			return err
		}
		if hasUsers && !allowRegistration {
			return auth.ErrRegistrationDisabled
		}

		exists, err := usernameExists(ctx, tx, user.Username)
		if err != nil {
			return err
		}
		if exists {
			return auth.ErrUserExists
		}

		stored := &User{
			Username:     user.Username,
			PasswordHash: user.PasswordHash,
			Role:         "user",
		}
		if !hasUsers {
			stored.Role = "admin"
		}

		if _, err := tx.NewInsert().Model(stored).Exec(ctx); err != nil {
			return fmt.Errorf("insert registered user: %w", err)
		}
		copyStoredUser(user, stored)
		return nil
	})
	if err != nil {
		return fmt.Errorf("create registered user: %w", err)
	}
	return nil
}

func (r *Repository) CreateInitialAdminIfNoUsers(
	ctx context.Context,
	user *auth.StoredUser,
) (bool, error) {
	created := false
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := r.lockUsersForBootstrap(ctx, tx); err != nil {
			return err
		}

		hasUsers, err := anyUsers(ctx, tx)
		if err != nil {
			return err
		}
		if hasUsers {
			return nil
		}

		stored := &User{
			Username:     user.Username,
			PasswordHash: user.PasswordHash,
			Role:         "admin",
		}
		if _, err := tx.NewInsert().Model(stored).Exec(ctx); err != nil {
			return fmt.Errorf("insert initial admin user: %w", err)
		}
		copyStoredUser(user, stored)
		created = true
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("create initial admin if no users: %w", err)
	}
	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*auth.StoredUser, error) {
	user := new(User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("scan user by id: %w", auth.ErrUserNotFound)
		}
		return nil, fmt.Errorf("scan user by id: %w", err)
	}
	return storedUser(user), nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*auth.StoredUser, error) {
	user := new(User)
	err := r.db.NewSelect().Model(user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("scan user by username: %w", auth.ErrUserNotFound)
		}
		return nil, fmt.Errorf("scan user by username: %w", err)
	}
	return storedUser(user), nil
}

func (r *Repository) HasUsers(ctx context.Context) (bool, error) {
	exists, err := r.db.NewSelect().Model((*User)(nil)).Where("1 = 1").Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("check users exist: %w", err)
	}
	return exists, nil
}

func storedUser(user *User) *auth.StoredUser {
	return &auth.StoredUser{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
	}
}

func copyStoredUser(destination *auth.StoredUser, source *User) {
	*destination = *storedUser(source)
}

func (r *Repository) lockUsersForBootstrap(ctx context.Context, tx bun.Tx) error {
	if isPostgresDialect(r.db) {
		if _, err := tx.ExecContext(ctx, `LOCK TABLE users IN EXCLUSIVE MODE`); err != nil {
			return fmt.Errorf("lock users table: %w", err)
		}
	}
	if isSQLiteDialect(r.db) {
		if _, err := tx.ExecContext(ctx, `UPDATE users SET updated_at = updated_at WHERE id = -1`); err != nil {
			return fmt.Errorf("lock users table: %w", err)
		}
	}
	return nil
}

func usernameExists(ctx context.Context, tx bun.Tx, username string) (bool, error) {
	var user User
	err := tx.NewSelect().
		Model(&user).
		Where("username = ?", username).
		Limit(1).
		Scan(ctx)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("check username: %w", err)
}

func anyUsers(ctx context.Context, tx bun.Tx) (bool, error) {
	exists, err := tx.NewSelect().
		Model((*User)(nil)).
		Where("1 = 1").
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("check users exist: %w", err)
	}
	return exists, nil
}

func isPostgresDialect(db *bun.DB) bool {
	name := fmt.Sprint(db.Dialect().Name())
	return name == "pg" || name == "postgres" || name == "postgresql"
}

func isSQLiteDialect(db *bun.DB) bool {
	name := fmt.Sprint(db.Dialect().Name())
	return name == "sqlite" || name == "sqlite3"
}
