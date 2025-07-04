package services

import (
	"good-api/internal/cache"
	"good-api/internal/repositories"

	"github.com/google/uuid"
)

type LeaderboardService struct {
	LeaderboardRepo *repositories.LeaderboardRepository
}

func NewLeaderboardService(repo *repositories.LeaderboardRepository) *LeaderboardService {
	return &LeaderboardService{LeaderboardRepo: repo}
}

// GetTournamentLeaderboard fetches the leaderboard of a specific tournament.
func (s *LeaderboardService) GetTournamentLeaderboard(tournamentID string, limit int) ([]string, error) {
	tID, err := uuid.Parse(tournamentID)
	if err != nil {
		return nil, err
	}
	return cache.GetTournamentLeaderboard(tID, limit)
}

// GetTournamentRank fetches the rank of a user in a tournament.
func (s *LeaderboardService) GetTournamentRank(userID string, tournamentID string) (int, error) {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}
	tID, err := uuid.Parse(tournamentID)
	if err != nil {
		return 0, err
	}
	return s.LeaderboardRepo.GetTournamentRank(uID, tID)
}
