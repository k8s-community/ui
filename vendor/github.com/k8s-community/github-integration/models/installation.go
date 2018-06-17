package models

import (
	"time"
)

// Possible types of sources
const (
	SourceGitHub = "github"
)

//go:generate reform

//reform:installations
type Installation struct {
	ID             int64  `reform:"id,pk"`
	Username       string `reform:"username"`
	InstallationID int    `reform:"installation_id"`

	CreatedAt time.Time `reform:"created_at"`
	UpdatedAt time.Time `reform:"updated_at"`
}

// BeforeInsert set CreatedAt and UpdatedAt.
func (i *Installation) BeforeInsert() error {
	i.CreatedAt = time.Now().UTC().Truncate(time.Second)
	i.UpdatedAt = i.CreatedAt
	return nil
}

// BeforeUpdate set UpdatedAt.
func (i *Installation) BeforeUpdate() error {
	i.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	return nil
}
