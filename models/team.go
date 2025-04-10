package models

import "time"

type Team struct {
	ID         int       `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	ExternalId int       `db:"external_id" json:"-"` // Not returned in JSON
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
