package models

import "time"

type PlayerStats struct {
	ID            int       `db:"id" json:"id"`
	PlayerID      int       `db:"player_id" json:"player_id"`
	Goals         int       `db:"goals" json:"goals"`
	Assists       int       `db:"assists" json:"assists"`
	CleanSheets   int       `db:"clean_sheets" json:"clean_sheets"`
	Saves         int       `db:"saves" json:"saves"`
	YellowCards   int       `db:"yellow_cards" json:"yellow_cards"`
	RedCards      int       `db:"red_cards" json:"red_cards"`
	MinutesPlayed int       `db:"minutes_played" json:"minutes_played"`
	OwnGoals      int       `db:"own_goals" json:"own_goals"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
