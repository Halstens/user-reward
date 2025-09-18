package postgress

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/user-reward/internal/models"
)

type RewardRepository struct {
	DB *sqlx.DB
}

func (wr *RewardRepository) GetUserById(id string) (*models.User, error) {
	stmt := `SELECT id, name, balance, created_at FROM users
		WHERE id = ?`
	row := wr.DB.QueryRow(stmt, id)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Balance, &user.CreatedAt)
	return user, err
}

func (wr *RewardRepository) GetTopList(id string) ([]*models.User, error) {
	var top_users []*models.User
	stmt := `SELECT name, balance FROM users
		ORDER BY balance DESC LIMIT 10`
	err := wr.DB.Select(&top_users, stmt, id)
	if err != nil {
		return nil, err
	}
	return top_users, err

}

func (wr *RewardRepository) UpdateUserBalance(id string, amount int, opType models.OperationType) error {
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
