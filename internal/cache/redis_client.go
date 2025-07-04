package cache

import (
	"context"
	"fmt"
	"os"
	"sync"

	"good-api/internal/models"
	"good-api/internal/repositories"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserLevelUpdater interface {
	IncreaseLevel(userID uuid.UUID) error
}

var ctx = context.Background()
var redisClient *redis.Client

// InitRedis initializes the Redis client.
func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "match3-redis:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
}

// AddUserToLeaderboard adds a user to a tournament leaderboard.
func AddUserToLeaderboard(tournamentID uuid.UUID, userID uuid.UUID, level int) {
	key := fmt.Sprintf("leaderboard:%s", tournamentID) // Each tournament has its own leaderboard.

	// Log when a user is added
	fmt.Printf("Adding user %s with level %d to leaderboard %s\n", userID, level, key)

	result, err := redisClient.ZAdd(ctx, key, []redis.Z{ // Redis sorted set.
		{
			Score:  float64(level),  // User's level determines ranking
			Member: userID.String(), // User's ID is stored.
		},
	}...).Result()

	if err != nil {
		fmt.Printf("Error adding user to leaderboard: %v\n", err)
	} else {
		fmt.Printf("ZADD result: %v\n", result)
	}
}

// GetTournamentLeaderboard retrieves tournament leaderboard sorted by level.
func GetTournamentLeaderboard(tournamentID uuid.UUID, limit int) ([]string, error) {
	return redisClient.ZRevRange(ctx, fmt.Sprintf("leaderboard:%s", tournamentID), 0, int64(limit-1)).Result()
}

// SyncLeaderboardsToDB syncs all tournament leaderboards to the database using concurrency.
func SyncLeaderboardsToDB(repo *repositories.TournamentRepository, userService UserLevelUpdater) {
	tournamentKeys, err := redisClient.Keys(ctx, "leaderboard:*").Result()
	if err != nil {
		fmt.Println("Error fetching leaderboard keys:", err)
		return
	}

	var wg sync.WaitGroup

	for _, key := range tournamentKeys {
		wg.Add(1)

		go func(tournamentKey string) {
			defer wg.Done()

			// Extract Tournament ID
			tournamentIDStr := tournamentKey[len("leaderboard:"):]
			tournamentID, err := uuid.Parse(tournamentIDStr)
			if err != nil {
				fmt.Println("Error parsing tournament ID:", err)
				return
			}

			fmt.Printf("Processing leaderboard for tournament: %s\n", tournamentID.String())

			// Fetch leaderboard from Redis
			leaderboard, err := redisClient.ZRevRangeWithScores(ctx, tournamentKey, 0, -1).Result()
			if err != nil {
				fmt.Println("Error fetching leaderboard data:", err)
				return
			}

			// Process top players using GORM ORM
			for index, entry := range leaderboard {
				userID, err := uuid.Parse(entry.Member.(string))
				if err != nil {
					fmt.Println("Skipping invalid user ID:", entry.Member)
					continue
				}

				reward := calculateReward(index + 1)

				if reward > 0 {
					// Use GORM to update user coins
					if err := repo.DB.Model(&models.User{}).
						Where("id = ?", userID).
						Update("coins", gorm.Expr("coins + ?", reward)).Error; err != nil {
						fmt.Println("Failed to update user coins:", err)
						continue
					}

					// Increase user level by 1 (if they placed in the top 10)
					if index < 10 {
						err = userService.IncreaseLevel(userID) // Dependency is injected, no import needed
						if err != nil {
							fmt.Println("Failed to increase user level:", err)
						}
					}
				}
			}

			// Delete tournament leaderboard from Redis after syncing
			redisClient.Del(ctx, tournamentKey)

		}(key)
	}

	wg.Wait() // Wait for all Go routines to complete
}

// Calculates the coin reward based on rank.
func calculateReward(rank int) int {
	switch {
	case rank == 1:
		return 5000
	case rank == 2:
		return 3000
	case rank == 3:
		return 2000
	case rank >= 4 && rank <= 10:
		return 1000
	default:
		return 0
	}
}

// Helper function to fetch environment variables.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func DeleteTournamentLeaderboard(tournamentID uuid.UUID) {
	key := fmt.Sprintf("leaderboard:%s", tournamentID)

	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		fmt.Println("Failed to delete tournament leaderboard from Redis: ", err)
	} else {
		fmt.Println("Tournament leaderboard deleted from Redis: ", key)
	}
}
