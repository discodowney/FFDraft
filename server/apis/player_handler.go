package apis

import (
	"net/http"
	"strconv"

	"go-app/models"
	"go-app/services"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerService *services.PlayerService
}

// NewPlayerHandler creates a new PlayerHandler instance
func NewPlayerHandler() *PlayerHandler {
	return &PlayerHandler{
		playerService: services.NewPlayerService(),
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

	// Build filter map
	filters := make(map[string]interface{})

	// Position filter
	if position := query.Get("position"); position != "" {
		// Validate position
		if err := h.playerService.ValidatePosition(models.Position(position)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid position. Must be one of: GK, DEF, MID, FWD",
			})
			return
		}
		filters["position"] = position
	}

	// Team filter
	if teamID := query.Get("team_id"); teamID != "" {
		if id, err := strconv.Atoi(teamID); err == nil {
			filters["team_id"] = id
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid team_id",
			})
			return
		}
	}

	// Name filter (partial match)
	if name := query.Get("name"); name != "" {
		filters["name"] = name
	}

	players, err := h.playerService.ListPlayers(filters)
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
