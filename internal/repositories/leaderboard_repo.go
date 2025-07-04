package repositories

import (
	"good-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaderboardRepository struct {
	DB *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) *LeaderboardRepository {
	return &LeaderboardRepository{DB: db}
}

// We should only get the users who are competing in a tournament
// GetGlobalLeaderboard fetches the top users globally based on level.
func (r *LeaderboardRepository) GetGlobalLeaderboard() ([]models.User, error) {
	var users []models.User

	err := r.DB.
		Joins("INNER JOIN tournament_participants tp ON users.id = tp.user_id").
		Order("users.level DESC").
		Limit(1000).
		Find(&users).Error
	return users, err
}

// We should only get the users who are competing in a tournament
// GetCountryLeaderboard fetches the top users in a specific country based on level.
func (r *LeaderboardRepository) GetCountryLeaderboard(country string) ([]models.User, error) {
	var users []models.User

	err := r.DB.
		Joins("INNER JOIN tournament_participants tp ON users.id = tp.user_id").
		Where("users.country = ?", country).
		Order("users.level DESC").
		Limit(1000).
		Find(&users).Error
	return users, err
}

// GetTournamentRank fetches a user's rank in a specific tournament.
func (r *LeaderboardRepository) GetTournamentRank(userID uuid.UUID, tournamentID uuid.UUID) (int, error) {
	var userLevel int

	// ✅ Get user's level by joining `users` with `tournament_participants`
	err := r.DB.Table("tournament_participants tp").
		Select("u.level").
		Joins("JOIN users u ON u.id = tp.user_id").
		Where("tp.user_id = ? AND tp.tournament_id = ?", userID, tournamentID).
		Scan(&userLevel).Error
	if err != nil {
		return 0, err
	}

	var rank int64

	// ✅ Count how many users have a **higher level** in the same tournament
	err = r.DB.Table("tournament_participants tp").
		Joins("JOIN users u ON u.id = tp.user_id").
		Where("tp.tournament_id = ? AND u.level > ?", tournamentID, userLevel).
		Count(&rank).Error
	if err != nil {
		return 0, err
	}

	// ✅ Rank = Users with higher levels + 1 (user's position)
	return int(rank) + 1, nil
}
