package tests

import (
	"good-api/internal/cache"
	"good-api/internal/handlers"
	"good-api/internal/models"
	"good-api/internal/repositories"
	"good-api/internal/services"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global test database instance
var testDB *gorm.DB
var once sync.Once

// SetupTestDB initializes the test PostgreSQL database.
func SetupTestDB() *gorm.DB {
	once.Do(func() {
		dsn := "host=match3-postgres user=postgres password=password dbname=testdb port=5432 sslmode=disable"

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to test database: %v", err)
		}

		// Enable UUID extension
		err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
		if err != nil {
			log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
		}

		// Apply database migrations
		err = db.AutoMigrate(&models.User{}, &models.Tournament{}, &models.TournamentParticipant{})
		if err != nil {
			log.Fatalf("Failed to migrate test database: %v", err)
		}

		testDB = db
	})

	return testDB
}

// SetupTestRedis ensures Redis is running for tests.
func SetupTestRedis() {
	cache.InitRedis()
}

// SeedTestData inserts test users, tournaments, and participants before tests run.
func SeedTestData(db *gorm.DB) (models.User, models.Tournament) {
	// Clean up previous test data

	db.Exec("DELETE FROM tournament_participants")
	db.Exec("DELETE FROM tournaments")
	db.Exec("DELETE FROM users")

	now := time.Now().UTC()
	startTime := now.Add(-1 * time.Hour)
	endTime := now.Add(10 * time.Hour)

	// Insert a test tournament
	tournament := models.Tournament{
		ID:        uuid.New(),
		StartTime: startTime,
		EndTime:   endTime,
		UserCount: 1,
		MaxUsers:  35,
		IsActive:  true,
	}
	db.Create(&tournament)

	// Insert a test user
	user := models.User{
		ID:       uuid.New(),
		Username: "test_user",
		Coins:    1000,
		Level:    15,
		Country:  "Turkey",
	}
	db.Create(&user)

	// Assign user to tournament
	participant := models.TournamentParticipant{
		ID:           uuid.New(),
		TournamentID: tournament.ID,
		UserID:       user.ID,
		Level:        user.Level,
	}
	db.Create(&participant)

	log.Println("Test data seeded successfully.")
	return user, tournament
}

func SetupRouter() *gin.Engine {
	db := SetupTestDB()
	SetupTestRedis()

	// repo
	userRepo := repositories.NewUserRepository(db)
	tournamentRepo := repositories.NewTournamentRepository(db)
	leaderboardRepo := repositories.NewLeaderboardRepository(db)

	// services
	userService := services.NewUserService(userRepo)
	tournamentService := services.NewTournamentService(tournamentRepo, userRepo)
	leaderboardService := services.NewLeaderboardService(leaderboardRepo)

	// Handlers
	userHandler := handlers.NewUserHandlerwithService(userRepo, userService)
	userJustHandler := handlers.NewUserHandlerwithRepo(userRepo)
	tournamentHandler := handlers.NewTournamentHandler(tournamentService, tournamentRepo)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService, leaderboardRepo)

	// Routes
	router := gin.Default()
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.GET("/:id", userJustHandler.GetUser)
		userRoutes.GET("/", userJustHandler.GetAllUsers)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	tournamentRoutes := router.Group("/tournaments")
	{
		tournamentRoutes.POST("/enter/:id", tournamentHandler.EnterTournament)
		tournamentRoutes.GET("/:id", tournamentHandler.GetTournament)
		tournamentRoutes.GET("/", tournamentHandler.GetAllTournaments)
		tournamentRoutes.POST("/update-score/:id", tournamentHandler.UpdateScore)
		tournamentRoutes.POST("/finish/:id", tournamentHandler.FinishTournament)
		tournamentRoutes.POST("/finish-all", tournamentHandler.FinishAllTournaments)
	}

	leaderboardRoutes := router.Group("/leaderboard")
	{
		leaderboardRoutes.GET("/global", leaderboardHandler.GetGlobalLeaderboard)
		leaderboardRoutes.GET("/country", leaderboardHandler.GetCountryLeaderboard)
		leaderboardRoutes.GET("/tournament", leaderboardHandler.GetTournamentLeaderboard)
		leaderboardRoutes.GET("/tournament/rank", leaderboardHandler.GetTournamentRank)
	}
	return router

}
