package team

import (
	"net/http"
	"strconv"

	"go-app/services/player"
	"go-app/services/team"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type TeamHandler struct {
	teamService   team.TeamService
	playerService player.PlayerService
}

// NewTeamHandler creates a new TeamHandler instance
func NewTeamHandler(db *sqlx.DB) *TeamHandler {
	return &TeamHandler{
		teamService:   team.NewTeamService(db),
		playerService: player.NewPlayerService(db),
	}
}

// GetTeam handles GET /api/teams/:id
func (h *TeamHandler) GetTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid team ID",
		})
		return
	}

	team, err := h.teamService.GetTeam(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve team",
		})
		return
	}

	c.JSON(http.StatusOK, team)
}

// ListTeams handles GET /api/teams with query parameters
func (h *TeamHandler) ListTeams(c *gin.Context) {
	// Get query parameters
	query := c.Request.URL.Query()

	// Initialize filter
	filter := &team.TeamFilter{}

	// Name filter (partial match)
	if name := query.Get("name"); name != "" {
		filter.Name = name
	}

	// External ID filter
	if externalID := query.Get("external_id"); externalID != "" {
		filter.ExternalID = externalID
	}

	teams, err := h.teamService.ListTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve teams",
		})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// GetTeamPlayers handles GET /api/teams/:id/players
func (h *TeamHandler) GetTeamPlayers(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := strconv.Atoi(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid team ID",
		})
		return
	}

	players, err := h.playerService.GetPlayersByTeam(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve team players",
		})
		return
	}

	c.JSON(http.StatusOK, players)
}
