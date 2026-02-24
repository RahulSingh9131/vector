package repository

import (
	"context"
	"errors"
	"time"

	"github.com/RahulSingh9131/vector/internal/database"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db     database.DBTX
	server *server.Server
}

// NewUserRepository creates a new user repository
func NewUserRepository(s *server.Server) *UserRepository {
	return &UserRepository{
		db:     s.DB.Pool,
		server: s,
	}
}

// GetByID retrieves a user by their internal UUID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, clerk_user_id, email, first_name, last_name, avatar_url, 
		       is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.ClerkUserID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetByClerkID retrieves a user by their Clerk user ID
func (r *UserRepository) GetByClerkID(ctx context.Context, clerkUserID string) (*models.User, error) {
	query := `
		SELECT id, clerk_user_id, email, first_name, last_name, avatar_url, 
		       is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE clerk_user_id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, clerkUserID).Scan(
		&user.ID,
		&user.ClerkUserID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, clerk_user_id, email, first_name, last_name, avatar_url, 
		       is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.ClerkUserID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, params models.CreateUserParams) (*models.User, error) {
	query := `
		INSERT INTO users (clerk_user_id, email, first_name, last_name, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, clerk_user_id, email, first_name, last_name, avatar_url, 
		          is_active, last_login_at, created_at, updated_at
	`

	var user models.User
	err := r.db.QueryRow(
		ctx,
		query,
		params.ClerkUserID,
		params.Email,
		params.FirstName,
		params.LastName,
		params.AvatarURL,
	).Scan(
		&user.ID,
		&user.ClerkUserID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates a user's information
func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, params models.UpdateUserParams) (*models.User, error) {
	query := `
		UPDATE users
		SET email = COALESCE($2, email),
		    first_name = COALESCE($3, first_name),
		    last_name = COALESCE($4, last_name),
		    avatar_url = COALESCE($5, avatar_url),
		    is_active = COALESCE($6, is_active)
		WHERE id = $1
		RETURNING id, clerk_user_id, email, first_name, last_name, avatar_url, 
		          is_active, last_login_at, created_at, updated_at
	`

	var user models.User
	err := r.db.QueryRow(
		ctx,
		query,
		id,
		params.Email,
		params.FirstName,
		params.LastName,
		params.AvatarURL,
		params.IsActive,
	).Scan(
		&user.ID,
		&user.ClerkUserID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET last_login_at = $2
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id, time.Now())
	return err
}

// Delete soft deletes a user by setting is_active to false
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET is_active = false
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

// GetOrCreate retrieves a user by Clerk ID or creates them if they don't exist (atomic)
func (r *UserRepository) GetOrCreate(ctx context.Context, params models.CreateUserParams) (*models.User, error) {
	query := `
		INSERT INTO users (clerk_user_id, email, first_name, last_name, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (clerk_user_id) DO NOTHING
	`

	_, err := r.db.Exec(
		ctx,
		query,
		params.ClerkUserID,
		params.Email,
		params.FirstName,
		params.LastName,
		params.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	// Always fetch the user (whether just inserted or already existed)
	return r.GetByClerkID(ctx, params.ClerkUserID)
}
// List retrieves all users
func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, clerk_user_id, email, first_name, last_name, avatar_url, 
		       is_active, last_login_at, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.ClerkUserID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.AvatarURL,
			&u.IsActive,
			&u.LastLoginAt,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
