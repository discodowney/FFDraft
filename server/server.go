package server

import (
	"fmt"
	"go-app/server/apis"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Custom logger middleware
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

// SetupRouter initializes and returns a configured Gin router
func SetupRouter(db *sqlx.DB) *gin.Engine {
	// Create a new gin router without default middleware
	r := gin.New()

	// Add custom logger middleware
	r.Use(customLogger())
	r.Use(gin.Recovery())

	// Initialize handlers
	userHandler := apis.NewUserHandler()
	playerHandler := apis.NewPlayerHandler()
	teamHandler := apis.NewTeamHandler(db)

	// API routes
	api := r.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.GET("", userHandler.ListUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Player routes
		players := api.Group("/players")
		{
			players.GET("", playerHandler.ListPlayers)
			players.GET("/:id", playerHandler.GetPlayer)
			players.GET("/:id/stats", playerHandler.GetPlayerStats)
		}

		// Team routes
		teams := api.Group("/teams")
		{
			teams.GET("", teamHandler.ListTeams)
			teams.GET("/:id", teamHandler.GetTeam)
		}

		// League routes with nested team routes
		leagues := api.Group("/leagues")
		{
			leagues.GET("/:leagueId/teams", teamHandler.GetTeamsByLeague)
		}

		// Team routes with nested player routes
		teams.GET("/:teamId/players", playerHandler.GetPlayersByTeam)
	}

	return r
}
