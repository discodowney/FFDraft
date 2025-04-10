package apis

import (
	"net/http"
	"strconv"

	"go-app/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type TeamHandler struct {
	teamService *services.TeamService
}

// NewTeamHandler creates a new TeamHandler instance
func NewTeamHandler(db *sqlx.DB) *TeamHandler {
	return &TeamHandler{
		teamService: services.NewTeamService(db),
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

	team, err := h.teamService.GetTeam(id)
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

	// Build filter map
	filters := make(map[string]interface{})

	// Name filter (partial match)
	if name := query.Get("name"); name != "" {
		filters["name"] = name
	}

	teams, err := h.teamService.ListTeams(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve teams",
		})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// GetTeamsByLeague handles GET /api/leagues/:leagueId/teams
func (h *TeamHandler) GetTeamsByLeague(c *gin.Context) {
	leagueIDStr := c.Param("leagueId")
	leagueID, err := strconv.Atoi(leagueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid league ID",
		})
		return
	}

	teams, err := h.teamService.GetTeamsByLeague(leagueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve league teams",
		})
		return
	}

	c.JSON(http.StatusOK, teams)
}
