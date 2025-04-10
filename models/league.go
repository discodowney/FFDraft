package models

import "time"

type League struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Country   string    `db:"country" json:"country"`
	Season    string    `db:"season" json:"season"`
	Code      string    `db:"code" json:"code"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
