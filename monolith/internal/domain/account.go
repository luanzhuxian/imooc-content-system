package domain

import (
	"errors"
	"time"
)

type Account struct {
	ID       int64     `json:"id"`
	UserID   string    `json:"user_id"`
	Password string    `json:"password"`
	Nickname string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Account) Validate() error {
	if a.UserID == "" || a.Password == "" || a.Nickname == "" {
		return errors.New("invalid account data")
	}
	return nil
}