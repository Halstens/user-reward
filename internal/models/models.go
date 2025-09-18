package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Reward int    `json:"reward"`
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
