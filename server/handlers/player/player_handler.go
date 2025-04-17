package player

import (
	"net/http"
	"strconv"

	"go-app/models"
	"go-app/services/player"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type PlayerHandler struct {
	playerService player.PlayerService
}

// NewPlayerHandler creates a new PlayerHandler instance
func NewPlayerHandler(db *sqlx.DB) *PlayerHandler {
	return &PlayerHandler{
		playerService: player.NewPlayerService(db),
	}
}

// GetPlayer handles GET /api/players/:id
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID",
		})
		return
	}

	player, err := h.playerService.GetPlayer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve player",
		})
		return
	}

	c.JSON(http.StatusOK, player)
}

// ListPlayers handles GET /api/players with query parameters
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	// Get query parameters
	query := c.Request.URL.Query()

	// Initialize filter
	filter := &player.PlayerFilter{}

	// Position filter
	if position := query.Get("position"); position != "" {
		// Validate position
		pos := models.Position(position)
		if err := h.playerService.ValidatePosition(pos); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid position. Must be one of: GK, DEF, MID, FWD",
			})
			return
		}
		filter.Position = pos
	}

	// Team filter
	if teamID := query.Get("team_id"); teamID != "" {
		if id, err := strconv.Atoi(teamID); err == nil {
			filter.TeamID = id
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid team_id",
			})
			return
		}
	}

	// First name filter (partial match)
	if firstName := query.Get("first_name"); firstName != "" {
		filter.FirstName = firstName
	}

	// Last name filter (partial match)
	if lastName := query.Get("last_name"); lastName != "" {
		filter.LastName = lastName
	}

	players, err := h.playerService.ListPlayers(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve players",
		})
		return
	}

	c.JSON(http.StatusOK, players)
}

// GetPlayersByTeam handles GET /api/teams/:teamId/players
func (h *PlayerHandler) GetPlayersByTeam(c *gin.Context) {
	teamIDStr := c.Param("teamId")
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

// GetPlayerStats handles GET /api/players/:id/stats
func (h *PlayerHandler) GetPlayerStats(c *gin.Context) {
	playerIDStr := c.Param("id")
	playerID, err := strconv.Atoi(playerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID",
		})
		return
	}

	stats, err := h.playerService.GetPlayerStats(playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve player stats",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
