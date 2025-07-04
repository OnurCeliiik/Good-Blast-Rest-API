package services

import (
	"errors"
	"fmt"
	"good-api/internal/cache"
	"good-api/internal/models"
	"good-api/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type TournamentService struct {
	TournamentRepo *repositories.TournamentRepository
	UserRepo       *repositories.UserRepository
}

func NewTournamentService(tournamentRepo *repositories.TournamentRepository, userRepo *repositories.UserRepository) *TournamentService {
	if tournamentRepo == nil || userRepo == nil {
		panic("TournamentService: Repositories must not be nil")
	}
	return &TournamentService{
		TournamentRepo: tournamentRepo,
		UserRepo:       userRepo,
	}
}

// EnterTournament handles adding a user to a tournament.
func (service *TournamentService) EnterTournament(userID uuid.UUID) (*models.Tournament, error) {
	now := time.Now().UTC()
	cutOffTime := time.Date(now.Year(), now.Month(), now.Day(), 19, 0, 0, 0, time.UTC)
	if now.After(cutOffTime) {
		return nil, errors.New("tournament entry is closed mate")
	}

	// Check if the UserRepo is initialized
	if service.UserRepo == nil {
		return nil, errors.New("userrepo is not initialized")
	}
	// Fetch the user safely
	user, err := service.UserRepo.GetUserByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Check if user meets entry requirements
	if user.Level < 10 || user.Coins < 500 {
		return nil, errors.New("user does not meet entry requirements")
	}

	// Ensure TournamentRepo is not nill
	if service.TournamentRepo == nil {
		return nil, errors.New("tournamentrepo is not initialized")
	}

	// Check if user is already in a tournament
	if _, err := service.TournamentRepo.GetUserTournament(userID); err != nil {
		return nil, errors.New("user is already in a tournament")
	}

	// Find an active tournament with space
	tournament, err := service.TournamentRepo.GetActiveTournament()
	if err != nil || tournament == nil || tournament.UserCount >= tournament.MaxUsers {
		tournament, err = service.TournamentRepo.NewTournament()
		if err != nil {
			return nil, err
		}
	}

	// Deduct coins before entering tournament
	user.Coins -= 500
	_, err = service.UserRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	// Add user to tournament
	err = service.TournamentRepo.AddParticipant(tournament.ID, userID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("User %s entered tournament %s. Adding to Redis...\n", userID, tournament.ID)

	cache.AddUserToLeaderboard(tournament.ID, user.ID, user.Level)

	return tournament, nil
}

func (service *TournamentService) GetTournamentByID(tournamentID uuid.UUID) (*models.Tournament, error) {
	return service.TournamentRepo.GetTournamentByID(tournamentID)
}

func (service *TournamentService) UpdateScore(userID uuid.UUID) error {
	tournament, err := service.TournamentRepo.GetUserTournament(userID)
	if err != nil {
		return errors.New("user is not in a tournament")
	}

	// Increase the user's score.
	err = service.TournamentRepo.IncreaseUserLevel(tournament.ID, userID)
	if err != nil {
		return err
	}
	return nil
}

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

func (service *TournamentService) FinishTournament(tournamentID uuid.UUID) error {
	// Mark the tournament as finished
	err := service.TournamentRepo.FinishTournament(tournamentID)
	if err != nil {
		return err
	}

	// Fetch leaderboard from Redis
	leaderboard, err := cache.GetTournamentLeaderboard(tournamentID, 35)
	if err != nil {
		return err
	}

	// Process rewards for top players
	for rank, userIDstr := range leaderboard {
		userID, err := uuid.Parse(userIDstr)
		if err != nil {
			continue
		}

		// Determine the reward based on rank
		reward := calculateReward(rank + 1)

		err = service.TournamentRepo.UpdateUserCoins(userID, reward)
		if err != nil {
			fmt.Println("Failed to update coins for user: ", userID, err)
		}

		// Top 10 players get a level-up
		if rank < 10 {
			err = service.TournamentRepo.IncreaseUserLevel(tournamentID, userID)
			if err != nil {
				fmt.Println("Failed to update level for user:", userID, err)
			}
		}
	}
	cache.DeleteTournamentLeaderboard(tournamentID)

	fmt.Println("Tournament finished and rewards processed", tournamentID)
	return nil
}

func (service *TournamentService) FinishAllTournaments() error {
	var activeTournaments []models.Tournament
	err := service.TournamentRepo.DB.Where("is_active = ?", true).Find(&activeTournaments).Error
	if err != nil {
		return err
	}

	if len(activeTournaments) == 0 {
		fmt.Println("No active tournaments to finish")
		return nil
	}

	for _, tournament := range activeTournaments {
		err := service.FinishTournament(tournament.ID)
		if err != nil {
			fmt.Println("Failed to finish tournament: ", err)
		}
	}

	topPlayers, err := service.TournamentRepo.GetTopGlobalPlayers()
	if err != nil {
		fmt.Println("Failed to fetch global rankings: ", err)
		return err
	}

	fmt.Println("Top 1000 Players across all tournaments:")
	for rank, user := range topPlayers {
		fmt.Printf("%d. %s - Level: %d\n", rank+1, user.Username, user.Level)
	}
	return nil
}
