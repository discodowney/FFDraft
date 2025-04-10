package models

import "time"

type UserTeamPlayer struct {
	ID         int       `db:"id" json:"id"`
	UserTeamID int       `db:"user_team_id" json:"user_team_id"`
	PlayerID   int       `db:"player_id" json:"player_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
