package models

import "time"

//go:generate reform

//reform:builds
type Build struct {
	ID         int64  `reform:"id,pk" json:"-"`
	UUID       string `reform:"uuid" json:"uuid"`
	Username   string `reform:"username" json:"username"`
	Repository string `reform:"repository" json:"repository"`
	Commit     string `reform:"commit" json:"commit"`
	Passed     bool   `reform:"passed" json:"passed"`
	Log        string `reform:"log" json:"log"`

	CreatedAt time.Time `reform:"created_at" json:"created_at"`
	UpdatedAt time.Time `reform:"updated_at" json:"updated_at"`
}

// BeforeInsert set CreatedAt and UpdatedAt.
func (b *Build) BeforeInsert() error {
	b.CreatedAt = time.Now().UTC().Truncate(time.Second)
	b.UpdatedAt = b.CreatedAt
	return nil
}

// BeforeUpdate set UpdatedAt.
func (b *Build) BeforeUpdate() error {
	b.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	return nil
}
