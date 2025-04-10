package models

import "time"

// IncidentType represents the type of incident that occurred
type IncidentType string

const (
	IncidentTypeGoal          IncidentType = "goal"
	IncidentTypeAssist        IncidentType = "assist"
	IncidentTypeYellowCard    IncidentType = "yellow_card"
	IncidentTypeRedCard       IncidentType = "red_card"
	IncidentTypeSubstitution  IncidentType = "substitution"
	IncidentTypeCleanSheet    IncidentType = "clean_sheet"
	IncidentTypePenaltyScored IncidentType = "penalty_scored"
	IncidentTypePenaltyMissed IncidentType = "penalty_missed"
	IncidentTypePenaltySaved  IncidentType = "penalty_saved"
	IncidentTypeOwnGoal       IncidentType = "own_goal"
)

// MatchIncident represents an event that occurred during a match
type MatchIncident struct {
	ID          int          `db:"id" json:"id"`
	MatchID     int          `db:"match_id" json:"match_id"`
	PlayerID    int          `db:"player_id" json:"player_id"`
	Type        IncidentType `db:"type" json:"type"`
	Minute      int          `db:"minute" json:"minute"`
	Description string       `db:"description" json:"description"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
}
