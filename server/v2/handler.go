package v2

import (
	"fmt"
	"time"

	"go-app/server/handlers/player"
	"go-app/server/handlers/team"
	"go-app/server/handlers/user"

	//"go-app/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	playerHandler *player.PlayerHandler
	teamHandler   *team.TeamHandler
	userHandler   *user.UserHandler
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		playerHandler: player.NewPlayerHandler(db),
		teamHandler:   team.NewTeamHandler(db),
		userHandler:   user.NewUserHandler(db),
	}
}

func customLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Format the log message
		logMsg := fmt.Sprintf("[%s] %s %d %v", method, path, statusCode, latency)
		fmt.Println(logMsg)
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.Use(customLogger())
	r.Use(gin.Recovery())
	// Apply rate limiting middleware
	//r.Use(middleware.RateLimit())
	// Apply custom logger middleware
	r.Use(customLogger())

	// Player routes with pagination and filtering
	players := r.Group("/players")
	{
		players.GET("", h.playerHandler.ListPlayers)
		players.GET("/:id", h.playerHandler.GetPlayer)
		players.GET("/team/:team_id", h.playerHandler.GetPlayersByTeam)
		players.GET("/:id/stats", h.playerHandler.GetPlayerStats)
		//players.GET("/search", h.playerHandler.SearchPlayers) // New search endpoint
	}

	// Team routes with enhanced features
	teams := r.Group("/teams")
	{
		teams.GET("", h.teamHandler.ListTeams)
		teams.GET("/:id", h.teamHandler.GetTeam)
		//teams.GET("/:id/roster", h.teamHandler.GetTeamRoster) // New roster endpoint
		//teams.GET("/:id/stats", h.teamHandler.GetTeamStats)   // New stats endpoint
	}

	// User routes with authentication
	users := r.Group("/users")
	{
		//users.POST("/register", h.userHandler.Register)
		//users.POST("/login", h.userHandler.Login)
		users.GET("", h.userHandler.ListUsers)
		users.GET("/:id", h.userHandler.GetUser)
		users.POST("", h.userHandler.CreateUser)
		users.PUT("/:id", h.userHandler.UpdateUser)
		users.DELETE("/:id", h.userHandler.DeleteUser)
		//users.GET("/me", h.userHandler.GetCurrentUser) // New current user endpoint
	}
}
