package main

import (
	"good-api/internal/cache"
	"good-api/internal/database"
	"good-api/internal/handlers"
	"good-api/internal/repositories"
	"good-api/internal/routes"
	"good-api/internal/services"
	"log"

	_ "good-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Good Blast Match 3 REST API
// @version 1.0
// @description API backend for the game
// @host localhost:8080
// @BasePath /
func main() {

	// Initialize Database
	database, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Initialize Redis
	cache.InitRedis()

	// Initialize User components
	userRepo := repositories.NewUserRepository(database)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandlerwithService(userRepo, userService)
	userJustHandler := handlers.NewUserHandlerwithRepo(userRepo)

	// Initialize Tournament components
	tournamentRepo := repositories.NewTournamentRepository(database)
	tournamentService := services.NewTournamentService(tournamentRepo, userRepo)
	tournamentHandler := handlers.NewTournamentHandler(tournamentService, tournamentRepo)

	// Initialize Leaderboard components
	leaderboardRepo := repositories.NewLeaderboardRepository(database)
	leaderboardService := services.NewLeaderboardService(leaderboardRepo)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService, leaderboardRepo)

	go cache.SyncLeaderboardsToDB(tournamentRepo, userService)

	// Setup Router
	router := gin.Default()
	routes.SetupRoutes(router, userHandler, userJustHandler, tournamentHandler, leaderboardHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start Server
	router.Run(":8080")
}
