package models

import "time"

// Position represents a player's position on the field
type Position string

const (
	PositionGK  Position = "GK"  // Goalkeeper
	PositionDEF Position = "DEF" // Defender
	PositionMID Position = "MID" // Midfielder
	PositionFWD Position = "FWD" // Forward
)

type Player struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Position  Position  `db:"position" json:"position"`
	TeamID    int       `db:"team_id" json:"team_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
