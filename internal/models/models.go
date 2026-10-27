package models

import (
	"time"

	"github.com/IFA-01/messenger/internal/repository/queries"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	LastSeen     time.Time `json:"last_seen"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func DatabaseUserToUser(dbUser queries.User) User {
	return User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
		Username:     dbUser.Username,
		Nickname:     dbUser.Nickname,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		LastSeen:     dbUser.LastSeen.Time,
	}
}
