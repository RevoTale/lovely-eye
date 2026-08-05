package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
)

var (
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserRegistrationDisabled = errors.New("user registration disabled")
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *UserRepository) CreateForRegistration(ctx context.Context, user *models.User, allowRegistration bool) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := r.lockUsersForBootstrap(ctx, tx); err != nil {
			return err
		}

		hasUsers, err := anyUsers(ctx, tx)
		if err != nil {
			return err
		}
		if hasUsers && !allowRegistration {
			return ErrUserRegistrationDisabled
		}

		exists, err := usernameExists(ctx, tx, user.Username)
		if err != nil {
			return err
		}
		if exists {
			return ErrUserAlreadyExists
		}

		user.Role = "user"
		if !hasUsers {
			user.Role = "admin"
		}

		if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
			return fmt.Errorf("insert registered user: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("create registered user: %w", err)
	}
	return nil
}

func (r *UserRepository) CreateInitialAdminIfNoUsers(ctx context.Context, user *models.User) (bool, error) {
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

		user.Role = "admin"
		if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
			return fmt.Errorf("insert initial admin user: %w", err)
		}
		created = true
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("create initial admin if no users: %w", err)
	}
	return created, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("scan user by id: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().Model(user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("scan user by username: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*models.User)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.NewSelect().Model(&users).Limit(limit).Offset(offset).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("scan user list: %w", err)
	}
	return users, nil
}

func (r *UserRepository) lockUsersForBootstrap(ctx context.Context, tx bun.Tx) error {
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
	var user models.User
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
		Model((*models.User)(nil)).
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
