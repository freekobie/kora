package postgres

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/freekobie/kora/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	conn *pgxpool.Pool
}

func NewUserStore(conn *pgxpool.Pool) model.UserStore {
	return &UserStore{
		conn: conn,
	}
}

// InsertUser implements model.UserStore.
func (u *UserStore) InsertUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, profile_photo, created_at, last_modified, verified)
		VALUES ($1, NULLIF($2,''), $3, $4, $5, $6, $7, $8);`

	_, err := u.conn.Exec(ctx, query,
		user.Id,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.ProfilePhoto,
		user.CreatedAt,
		user.LastModifed,
		user.Verified,
	)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return model.ErrDuplicateUser
		}
		slog.Error("failed to insert user", "error", err)
		return err
	}
	return nil
}

// DeleteUser implements model.UserStore.
func (u *UserStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1;`

	result, err := u.conn.Exec(ctx, query, id)
	if err != nil {
		slog.Error("failed delete user", "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

// GetUser implements model.UserStore.
func (u *UserStore) GetUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	query := `
		SELECT id, name, email, password_hash, profile_photo, created_at, last_modified, verified 
		FROM users 
		WHERE id = $1;`

	var user model.User
	err := u.conn.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePhoto,
		&user.CreatedAt,
		&user.LastModifed,
		&user.Verified,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, model.ErrNotFound
	}

	return user, nil
}

// GetUserByMail implements model.UserStore.
func (u *UserStore) GetUserByMail(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT id, name, email, password_hash, profile_photo, created_at, last_modified, verified 
		FROM users 
		WHERE email = $1;`

	var user model.User
	err := u.conn.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePhoto,
		&user.CreatedAt,
		&user.LastModifed,
		&user.Verified,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, model.ErrNotFound
	}

	return user, nil
}

// UpdateUser implements model.UserStore.
func (u *UserStore) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users 
		SET name = $1, email = $2, password_hash = $3, profile_photo = $4, last_modified = $5, verified = $6
		WHERE id = $7;`

	result, err := u.conn.Exec(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.ProfilePhoto,
		user.LastModifed,
		user.Verified,
		user.Id,
	)
	if err != nil {
		slog.Error("failed update user", "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

// InsertToken implements model.TokenStore.
func (t *UserStore) InsertToken(ctx context.Context, token *model.UserToken) error {
	query := `INSERT INTO user_tokens(token_hash, user_id, scope, expires_at)
	VALUES($1, $2, $3, $4);`

	_, err := t.conn.Exec(ctx, query, token.Hash, token.UserId, token.Scope, token.ExpiresAt)
	if err != nil {
		slog.Error("failed to insert token", "error", err)
		return err
	}

	return nil
}

// GetUserForToken implements model.UserStore
func (t *UserStore) GetUserForToken(ctx context.Context, tokenHash string, scope string, email string) (*model.User, error) {
	query := `SELECT
	users.id,
	users.name,
	users.email,
	users.password_hash,
	users.profile_photo,
	users.verified,
	users.created_at,
	users.last_modified
	FROM users
	JOIN user_tokens AS tokens
	ON users.id = tokens.user_id
	WHERE tokens.token_hash = $1
	AND tokens.scope = $2
	AND users.email = $3
	AND tokens.expires_at > now();
	`

	var user model.User
	row := t.conn.QueryRow(ctx, query, tokenHash, scope, email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash, &user.ProfilePhoto, &user.Verified, &user.CreatedAt, &user.LastModifed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		slog.Error("failed to fetch token", "error", err)
		return nil, err
	}

	return &user, nil
}

// DeleteToken implements model.TokenStore.
func (t *UserStore) DeleteToken(ctx context.Context, tokenHash, scope string) error {
	query := `DELETE FROM user_tokens WHERE token_hash = $1 AND scope = $2;`

	_, err := t.conn.Exec(ctx, query, tokenHash, scope)
	if err != nil {
		slog.Error("failed to delete user token", "error", err)
		return err
	}

	return nil
}
