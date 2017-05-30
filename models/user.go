package models

import (
	"time"
)

// Possible types of sources
const (
	SourceGitHub = "github"
)

//go:generate reform

//reform:users
type User struct {
	ID          int64     `reform:"id,pk"`
	Name        string    `reform:"name"`
	Source      string    `reform:"source"`
	SessionID   *string   `reform:"session_id"`
	SessionData *string   `reform:"session_data"`
	CreatedAt   time.Time `reform:"created_at"`
	UpdatedAt   time.Time `reform:"updated_at"`
}

// BeforeInsert set CreatedAt and UpdatedAt.
func (u *User) BeforeInsert() error {
	u.CreatedAt = time.Now().UTC().Truncate(time.Second)
	u.UpdatedAt = u.CreatedAt
	return nil
}

// BeforeUpdate set UpdatedAt.
func (u *User) BeforeUpdate() error {
	u.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	return nil
}
