package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"rwiratama.com/m/internal/models"
)

type UserRepository struct {
	pool PgxPool
}

func NewUserRepository(pool PgxPool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, username, email string) (*models.User, error) {
	query := `
		INSERT INTO users (username, email)
		VALUES ($1, $2)
		RETURNING uid, username, email, is_admin, created_at, updated_at
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, username, email).Scan(
		&user.UID,
		&user.Username,
		&user.Email,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// CreateWithPassword inserts a new user with a hashed password and admin flag.
// The caller is responsible for hashing the password (bcrypt) before calling.
func (r *UserRepository) CreateWithPassword(ctx context.Context, username, email, passwordHash string, isAdmin bool) (*models.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, is_admin)
		VALUES ($1, $2, $3, $4)
		RETURNING uid, username, email, is_admin, created_at, updated_at
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, username, email, passwordHash, isAdmin).Scan(
		&user.UID,
		&user.Username,
		&user.Email,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GetByID retrieves a user by UID
func (r *UserRepository) GetByID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	query := `
		SELECT uid, username, email, is_admin, created_at, updated_at
		FROM users
		WHERE uid = $1
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, uid).Scan(
		&user.UID,
		&user.Username,
		&user.Email,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetCredentialsByUsername retrieves a user's password hash plus auth flags.
// Used by the auth handler when verifying login credentials. The returned
// struct is not exposed via the API.
func (r *UserRepository) GetCredentialsByUsername(ctx context.Context, username string) (*models.UserCredentials, error) {
	query := `
		SELECT uid, username, email, password_hash, is_admin
		FROM users
		WHERE username = $1
	`

	var creds models.UserCredentials
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&creds.UID,
		&creds.Username,
		&creds.Email,
		&creds.PasswordHash,
		&creds.IsAdmin,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user credentials: %w", err)
	}

	return &creds, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT uid, username, email, is_admin, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.User, error) {
		var user models.User
		err := row.Scan(
			&user.UID,
			&user.Username,
			&user.Email,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		return user, err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to collect users: %w", err)
	}

	return users, nil
}

// Update updates a user's information
func (r *UserRepository) Update(ctx context.Context, uid uuid.UUID, username, email string) (*models.User, error) {
	query := `
		UPDATE users
		SET username = $2, email = $3, updated_at = CURRENT_TIMESTAMP
		WHERE uid = $1
		RETURNING uid, username, email, is_admin, created_at, updated_at
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, uid, username, email).Scan(
		&user.UID,
		&user.Username,
		&user.Email,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

// Delete removes a user from the database
func (r *UserRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	query := `DELETE FROM users WHERE uid = $1`

	result, err := r.pool.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
