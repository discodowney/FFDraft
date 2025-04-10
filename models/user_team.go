package models

import "time"

type UserTeam struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	UserID    int       `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
