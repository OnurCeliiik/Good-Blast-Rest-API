package routes

/*
The routes layer is responsible for defining HTTP endpoints and connecting them
to the corresponding handler functions.
It acts as a bridge between the incoming API requests and the handler layer.
*/
/*
It will define which HTTP methods (GET, POST, PUT, DELETE) are used for each action.
It will map URLs to handler functions.
It will register all the endpoints when the app starts.
*/

import (
	"github.com/gin-gonic/gin"

	"good-api/internal/handlers"
)

// SetupRoutes defines all API routes and connects them to handlers.
func SetupRoutes(router *gin.Engine, userHandler *handlers.UserHandler, userJustHandler *handlers.UserHandler, tournamentHandler *handlers.TournamentHandler, leaderboardHandler *handlers.LeaderboardHandler) {

	// User routes
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUser)      // Create a user
		userRoutes.GET("/:id", userJustHandler.GetUser)   // Get a user by ID
		userRoutes.GET("/", userJustHandler.GetAllUsers)  // Get all users
		userRoutes.PUT("/:id", userHandler.UpdateUser)    // Update user
		userRoutes.DELETE("/:id", userHandler.DeleteUser) // Delete user

	}

	// Tournament routes
	tournamentRoutes := router.Group("/tournaments")
	{
		tournamentRoutes.GET("/", tournamentHandler.GetAllTournaments)               // Get all tournaments
		tournamentRoutes.POST("/enter/:id", tournamentHandler.EnterTournament)       // Enter a tournament
		tournamentRoutes.GET("/:id", tournamentHandler.GetTournament)                // Get tournament details
		tournamentRoutes.POST("/finish/:id", tournamentHandler.FinishTournament)     // Manually finish a tournament
		tournamentRoutes.POST("/finish-all", tournamentHandler.FinishAllTournaments) // Manually finish all tournaments
		tournamentRoutes.PUT("/update-score/:id", tournamentHandler.UpdateScore)     // Update user's level in a tournament
	}

	// Leaderboard routes
	leaderboardRoutes := router.Group("/leaderboard")
	{
		leaderboardRoutes.GET("/global", leaderboardHandler.GetGlobalLeaderboard)   // will get users who compete in any tournament and rank them globally.
		leaderboardRoutes.GET("/country", leaderboardHandler.GetCountryLeaderboard) // will get users who compete in any tournament and rank them according to country we choose.
		leaderboardRoutes.GET("/tournament", leaderboardHandler.GetTournamentLeaderboard)
		leaderboardRoutes.GET("/tournament/rank", leaderboardHandler.GetTournamentRank)
	}

}
