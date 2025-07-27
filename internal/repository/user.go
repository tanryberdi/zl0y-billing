package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"zl0y-billing/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(login, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO users (login, password_hash, balance)
		VALUES ($1, $2, 10000) -- 100.00 in cents as a starting balance
		RETURNING id, login, password_hash, balance, created_at
	`

	var user models.User
	err := r.db.QueryRow(query, login, passwordHash).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByLogin(login string) (*models.User, error) {
	query := `
		SELECT id, login, password_hash, balance, created_at
		FROM users
		WHERE login = $1
	`

	var user models.User
	err := r.db.QueryRow(query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found") // User not found
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT id, login, password_hash, balance, created_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found") // User not found
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserBalance(userID int, newBalance int) error {
	query := `
		UPDATE users
		SET balance = $2
		WHERE id = $1
	`

	result, err := r.db.Exec(query, userID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}

func (r *UserRepository) DeductBalance(userID, amount int) error {
	query := `
		UPDATE users
		SET balance = balance - $2
		WHERE id = $1 AND balance >= $2
		RETURNING balance
	`

	var newBalance int
	err := r.db.QueryRow(query, userID, amount).Scan(&newBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("insufficient balance or user not found")
		}
		return fmt.Errorf("failed to deduct balance: %w", err)
	}

	return nil
}
