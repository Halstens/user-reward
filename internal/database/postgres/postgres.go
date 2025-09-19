package postgress

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/user-reward/internal/models"
)

type RewardRepository struct {
	DB *sqlx.DB
}

func (wr *RewardRepository) GetUserById(id int) (*models.User, error) {
	stmt := `SELECT id, name, balance, created_at FROM users
		WHERE id = $1`
	row := wr.DB.QueryRow(stmt, id)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Balance, &user.CreatedAt)
	return user, err
}

func (wr *RewardRepository) GetTopList() ([]*models.User, error) {
	var top_users []*models.User
	stmt := `SELECT name, balance FROM users
		ORDER BY balance DESC LIMIT 10`
	err := wr.DB.Select(&top_users, stmt)
	if err != nil {
		return nil, err
	}
	return top_users, err

}

func (wr *RewardRepository) IsTaskCompleted(ctx context.Context, userID int, taskType string) (bool, error) {
	query := `SELECT COUNT(*) FROM completed_tasks WHERE user_id = $1 AND task_type = $2`
	var count int
	err := wr.DB.QueryRowContext(ctx, query, userID, taskType).Scan(&count)
	return count > 0, err
}

func (wr *RewardRepository) MarkTaskCompleted(ctx context.Context, userID int, taskType string) error {
	query := `INSERT INTO completed_tasks (user_id, task_type) VALUES ($1, $2)`
	_, err := wr.DB.ExecContext(ctx, query, userID, taskType)
	return err
}

func (wr *RewardRepository) UpdateUserBalance(id int, amount int, opType models.OperationType) error {
	var query string
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	tx, err := wr.DB.Begin()
	if err != nil {
		return fmt.Errorf("fail transaction: %w", err)
	}

	defer tx.Rollback()

	switch opType {
	case "subscribe_telegram":
		query = "UPDATE users SET balance = balance + $1 WHERE id = $2 RETURNING balance"
	case "follow_twitter":
		query = "UPDATE users SET balance = balance + $1 WHERE id = $2 RETURNING balance"
	case "referral_signup":
		query = "UPDATE users SET balance = balance + $1 WHERE id = $2 RETURNING balance"
	default:
		return fmt.Errorf("invalid op type")
	}

	var newBalance int
	err = tx.QueryRow(query, amount, id).Scan(&newBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("there is no such user")
		}
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return err
}

func (wr *RewardRepository) IsReferralExists(ctx context.Context, refereeID int) (bool, error) {
	query := `SELECT COUNT(*) FROM referrals WHERE referee_id = $1`
	var count int
	err := wr.DB.QueryRowContext(ctx, query, refereeID).Scan(&count)
	return count > 0, err
}

func (wr *RewardRepository) CreateReferral(ctx context.Context, referrerID, refereeID int) error {
	query := `INSERT INTO referrals (referrer_id, referee_id) VALUES ($1, $2)`
	_, err := wr.DB.ExecContext(ctx, query, referrerID, refereeID)
	return err
}

func (wr *RewardRepository) ProcessReferral(ctx context.Context, referrerID, refereeID, reward int) error {
	tx, err := wr.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Добавляем запись о реферале
	_, err = tx.ExecContext(ctx,
		`INSERT INTO referrals (referrer_id, referee_id) VALUES ($1, $2)`,
		referrerID, refereeID)
	if err != nil {
		return fmt.Errorf("failed to create referral: %w", err)
	}

	// Начисляем награду рефереру
	_, err = tx.ExecContext(ctx,
		`UPDATE users SET balance = balance + $1 WHERE id = $2`,
		reward, referrerID)
	if err != nil {
		return fmt.Errorf("failed to update referrer balance: %w", err)
	}

	// Отмечаем задание "referral_signup" как выполненное для реферала
	_, err = tx.ExecContext(ctx,
		`INSERT INTO completed_tasks (user_id, task_type) VALUES ($1, 'referral_signup')`,
		refereeID)
	if err != nil {
		return fmt.Errorf("failed to mark referral task: %w", err)
	}

	return tx.Commit()
}

func (wr *RewardRepository) GetUserByUsername(ctx context.Context, username string) (*UserWithAuth, error) {
	query := `SELECT id, name, balance, password_hash 
			  FROM users 
			  WHERE name = $1`

	row := wr.DB.QueryRowContext(ctx, query, username)

	var user UserWithAuth
	err := row.Scan(&user.ID, &user.Name, &user.Balance, &user.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

type UserWithAuth struct {
	ID           int
	Name         string
	Balance      int
	PasswordHash string
}
