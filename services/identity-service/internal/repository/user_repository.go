package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, role string, limit, offset int) ([]*domain.User, error)
	ListByRole(ctx context.Context, role domain.UserRole) ([]*domain.User, error)
	Count(ctx context.Context) (int, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, username, display_name, password_hash, bio, avatar_url, banner_url, verified, role, status, categories, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		user.ID,
		user.Email,
		user.Username,
		user.DisplayName,
		user.PasswordHash,
		user.Bio,
		user.AvatarURL,
		user.BannerURL,
		user.Verified,
		user.Role,
		user.Status,
		user.Categories,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if isPGUniqueViolation(err) {
			if strings.Contains(err.Error(), "email") {
				return domain.ErrEmailTaken
			}
			if strings.Contains(err.Error(), "username") {
				return domain.ErrUsernameTaken
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, bio, avatar_url, banner_url,
		       verified, role, status, categories, created_at, updated_at
		FROM users WHERE id = $1`

	user := &domain.User{}
	var categories []string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.PasswordHash,
		&user.Bio, &user.AvatarURL, &user.BannerURL,
		&user.Verified, &user.Role, &user.Status, &categories,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	user.Categories = categories
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, bio, avatar_url, banner_url,
		       verified, role, status, categories, created_at, updated_at
		FROM users WHERE email = $1`

	user := &domain.User{}
	var categories []string
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.PasswordHash,
		&user.Bio, &user.AvatarURL, &user.BannerURL,
		&user.Verified, &user.Role, &user.Status, &categories,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	user.Categories = categories
	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, bio, avatar_url, banner_url,
		       verified, role, status, categories, created_at, updated_at
		FROM users WHERE username = $1`

	user := &domain.User{}
	var categories []string
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.PasswordHash,
		&user.Bio, &user.AvatarURL, &user.BannerURL,
		&user.Verified, &user.Role, &user.Status, &categories,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	user.Categories = categories
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			email = $2, username = $3, display_name = $4, password_hash = $5,
			bio = $6, avatar_url = $7, banner_url = $8, verified = $9,
			role = $10, status = $11, categories = $12
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.Username, user.DisplayName, user.PasswordHash,
		user.Bio, user.AvatarURL, user.BannerURL,
		user.Verified, user.Role, user.Status, user.Categories,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET status = 'banned', updated_at = NOW() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Search(ctx context.Context, query string, role string, limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var rows pgx.Rows
	var err error

	if query != "" && role != "" {
		rows, err = r.pool.Query(ctx, `
			SELECT id, email, username, display_name, password_hash, bio, avatar_url, banner_url,
			       verified, role, status, categories, created_at, updated_at
			FROM users
			WHERE (username ILIKE '%' || $1 || '%' OR display_name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
			AND role = $2
			ORDER BY created_at DESC LIMIT $3 OFFSET $4`, query, role, limit, offset)
	} else if query != "" {
		rows, err = r.pool.Query(ctx, `
			SELECT id, email, username, display_name, password_hash, bio, avatar_url, banner_url,
			       verified, role, status, categories, created_at, updated_at
			FROM users
			WHERE username ILIKE '%' || $1 || '%' OR display_name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%'
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`, query, limit, offset)
	} else if role != "" {
		rows, err = r.pool.Query(ctx, `
			SELECT id, email, username, display_name, password_hash, bio, avatar_url, banner_url,
			       verified, role, status, categories, created_at, updated_at
			FROM users
			WHERE role = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`, role, limit, offset)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, email, username, display_name, password_hash, bio, avatar_url, banner_url,
			       verified, role, status, categories, created_at, updated_at
			FROM users
			ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var categories []string
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.PasswordHash,
			&user.Bio, &user.AvatarURL, &user.BannerURL,
			&user.Verified, &user.Role, &user.Status, &categories,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		user.Categories = categories
		users = append(users, user)
	}
	if users == nil {
		users = []*domain.User{}
	}
	return users, nil
}

func (r *userRepository) ListByRole(ctx context.Context, role domain.UserRole) ([]*domain.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, email, username, display_name, password_hash, bio, avatar_url, banner_url,
		       verified, role, status, categories, created_at, updated_at
		FROM users WHERE role = $1
		ORDER BY created_at DESC`, role)
	if err != nil {
		return nil, fmt.Errorf("failed to list users by role: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var categories []string
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.PasswordHash,
			&user.Bio, &user.AvatarURL, &user.BannerURL,
			&user.Verified, &user.Role, &user.Status, &categories,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		user.Categories = categories
		users = append(users, user)
	}
	if users == nil {
		users = []*domain.User{}
	}
	return users, nil
}

func (r *userRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

func isPGUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
