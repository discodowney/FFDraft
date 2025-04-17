package v1

import (
	"fmt"
	"time"

	"go-app/server/handlers/player"
	"go-app/server/handlers/team"
	"go-app/server/handlers/user"

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
	// Apply custom logger middleware
	r.Use(customLogger())
	r.Use(gin.Recovery())

	// Player routes
	players := r.Group("/players")
	{
		players.GET("", h.playerHandler.ListPlayers)
		players.GET("/:id", h.playerHandler.GetPlayer)
		players.GET("/team/:team_id", h.playerHandler.GetPlayersByTeam)
		players.GET("/:id/stats", h.playerHandler.GetPlayerStats)
	}

	// Team routes
	teams := r.Group("/teams")
	{
		teams.GET("", h.teamHandler.ListTeams)
		teams.GET("/:id", h.teamHandler.GetTeam)
	}

	// User routes
	users := r.Group("/users")
	{
		users.GET("", h.userHandler.ListUsers)
		users.GET("/:id", h.userHandler.GetUser)
		users.POST("", h.userHandler.CreateUser)
		users.PUT("/:id", h.userHandler.UpdateUser)
		users.DELETE("/:id", h.userHandler.DeleteUser)
	}
}

// StartServer initializes and starts the HTTP server
func StartServer(db *sqlx.DB) {
	// Initialize Gin router
	router := gin.Default()

	// Create API v1 handler
	v1Handler := NewHandler(db)

	// Register v1 routes
	v1 := router.Group("/api/v1")
	v1Handler.RegisterRoutes(v1)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := router.Run(port); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
