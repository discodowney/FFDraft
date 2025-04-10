package models

import "time"

type Match struct {
	ID         int       `db:"id" json:"id"`
	LeagueID   int       `db:"league_id" json:"league_id"`
	HomeTeamID int       `db:"home_team_id" json:"home_team_id"`
	AwayTeamID int       `db:"away_team_id" json:"away_team_id"`
	MatchDate  time.Time `db:"match_date" json:"match_date"`
	HomeScore  int       `db:"home_score" json:"home_score"`
	AwayScore  int       `db:"away_score" json:"away_score"`
	Status     string    `db:"status" json:"status"` // scheduled, in_progress, completed, postponed
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
