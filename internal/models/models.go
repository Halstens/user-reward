package models

import (
	"errors"
	"time"
)

type OperationType string

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

const (
	subscribe_telegram OperationType = "subscribe_telegram"
	follow_twitter     OperationType = "follow_twitter"
	referral_signup    OperationType = "referral_signup"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID            int           `json:"id"`
	OperationType OperationType `json:"operationType"`
	Reward        int           `json:"reward"`
}

type UserTask struct {
	UserID    int       `json:"user_id"`
	TaskID    int       `json:"task_id"`
	Completed time.Time `json:"completed"`
}

type Referral struct {
	ReferrerID int       `json:"referrer_id"`
	RefereeID  int       `json:"referee_id"`
	CreatedAt  time.Time `json:"created_at"`
}
