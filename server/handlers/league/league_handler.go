package league

import (
	"net/http"
	"strconv"

	"go-app/models"
	"go-app/services/league"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type LeagueHandler struct {
	leagueService league.LeagueService
}

// NewLeagueHandler creates a new LeagueHandler instance
func NewLeagueHandler(db *sqlx.DB) *LeagueHandler {
	return &LeagueHandler{
		leagueService: league.NewLeagueService(db),
	}
}

// GetLeague handles GET /api/leagues/:id
func (h *LeagueHandler) GetLeague(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid league ID",
		})
		return
	}

	league, err := h.leagueService.GetLeague(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve league",
		})
		return
	}

	c.JSON(http.StatusOK, league)
}

// ListLeagues handles GET /api/leagues
func (h *LeagueHandler) ListLeagues(c *gin.Context) {
	leagues, err := h.leagueService.ListLeagues()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve leagues",
		})
		return
	}

	c.JSON(http.StatusOK, leagues)
}

// CreateLeague handles POST /api/leagues
func (h *LeagueHandler) CreateLeague(c *gin.Context) {
	var league models.League
	if err := c.ShouldBindJSON(&league); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	createdLeague, err := h.leagueService.CreateLeague(&league)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create league",
		})
		return
	}

	c.JSON(http.StatusCreated, createdLeague)
}

// UpdateLeague handles PUT /api/leagues/:id
func (h *LeagueHandler) UpdateLeague(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid league ID",
		})
		return
	}

	var league models.League
	if err := c.ShouldBindJSON(&league); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	league.ID = id
	updatedLeague, err := h.leagueService.UpdateLeague(&league)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update league",
		})
		return
	}

	c.JSON(http.StatusOK, updatedLeague)
}

// DeleteLeague handles DELETE /api/leagues/:id
func (h *LeagueHandler) DeleteLeague(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid league ID",
		})
		return
	}

	err = h.leagueService.DeleteLeague(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete league",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetLeagueByCode handles GET /api/leagues/code/:code
func (h *LeagueHandler) GetLeagueByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "League code is required",
		})
		return
	}

	league, err := h.leagueService.GetLeagueByCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve league",
		})
		return
	}

	if league == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "League not found",
		})
		return
	}

	c.JSON(http.StatusOK, league)
}
