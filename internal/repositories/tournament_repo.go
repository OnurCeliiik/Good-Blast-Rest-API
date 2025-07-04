package repositories

import (
	"errors"
	"fmt"
	"good-api/internal/models"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TournamentRepository struct {
	DB *gorm.DB
}

func NewTournamentRepository(db *gorm.DB) *TournamentRepository {
	return &TournamentRepository{DB: db}
}

// Create a new tournament
func (repo *TournamentRepository) NewTournament() (*models.Tournament, error) {
	var count int64
	repo.DB.Model(&models.Tournament{}).Count(&count) // Count existing tournaments
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	endTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 0, 0, time.UTC)

	tournament := &models.Tournament{
		ID:        uuid.New(),
		Name:      fmt.Sprintf("tournament_%d", count+1),
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  true,
		UserCount: 0,
		MaxUsers:  35,
	}

	if err := repo.DB.Create(tournament).Error; err != nil {
		return nil, err
	}
	return tournament, nil
}

// Fetch an active tournament
func (repo *TournamentRepository) GetActiveTournament() (*models.Tournament, error) {
	var tournament models.Tournament
	err := repo.DB.Where("is_active = ? AND user_count < ?", true, 35).First(&tournament).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tournament, err
}

// Get user's tournament
func (repo *TournamentRepository) GetUserTournament(userID uuid.UUID) (*models.Tournament, error) {
	var participant models.TournamentParticipant
	err := repo.DB.Where("user_id = ?", userID).First(&participant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var tournament models.Tournament
	err = repo.DB.Where("id = ?", participant.TournamentID).First(&tournament).Error
	if err != nil {
		return nil, err
	}
	fmt.Println("Looking up tournament participant for user: ", userID)
	return &tournament, nil
}

// Add a participant to a tournament
func (repo *TournamentRepository) AddParticipant(tournamentID, userID uuid.UUID) error {
	var user models.User

	// Get the current level
	if err := repo.DB.First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("failed to find user for level info: %w", err)
	}

	// Create participant with user's level
	participant := &models.TournamentParticipant{
		ID:           uuid.New(),
		UserID:       userID,
		TournamentID: tournamentID,
		Level:        user.Level,
	}

	if err := repo.DB.Create(participant).Error; err != nil {
		return err
	}

	// Increase tournament user count
	return repo.DB.Model(&models.Tournament{}).
		Where("id = ?", tournamentID).
		Update("user_count", gorm.Expr("user_count + 1")).Error
}

// Increase user score in a tournament
func (repo *TournamentRepository) IncreaseUserLevel(tournamentID, userID uuid.UUID) error {
	return repo.DB.Model(&models.TournamentParticipant{}).
		Where("tournament_id = ? AND user_id = ?", tournamentID, userID).
		Update("level", gorm.Expr("level + 1")).Error
}

// Update user coins
func (repo *TournamentRepository) UpdateUserCoins(userID uuid.UUID, coins int) error {
	return repo.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("coins", gorm.Expr("coins + ?", coins)).Error
}

// Get tournament by ID
func (repo *TournamentRepository) GetTournamentByID(tournamentID uuid.UUID) (*models.Tournament, error) {
	var tournament models.Tournament
	err := repo.DB.Where("id = ?", tournamentID).First(&tournament).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tournament, err
}

// Get all tournaments (ordered by start time)
func (repo *TournamentRepository) GetAllTournaments() ([]models.Tournament, error) {
	var tournaments []models.Tournament

	if repo.DB == nil {
		fmt.Println("Error: repo.DB is nil")
		return nil, errors.New("database connection is nil")
	}

	fmt.Println("Running SQL Query: SELECT * FROM tournaments ORDER BY start_time DESC")

	err := repo.DB.Order("start_time DESC").Find(&tournaments).Error

	if err != nil {
		log.Printf("Error fetching tournaments: %v", err)
		return nil, err
	}

	if len(tournaments) == 0 {
		log.Println("No tournaments found in the database")
	}
	fmt.Println("Successfully fetched tournaments: ", tournaments)
	return tournaments, nil
}

// Finish a tournament
func (repo *TournamentRepository) FinishTournament(tournamentID uuid.UUID) error {
	return repo.DB.Model(&models.Tournament{}).
		Where("id = ?", tournamentID).
		Update("is_active", false).Error
}

// Get top 1000 players across all tournaments (global ranking)
func (repo *TournamentRepository) GetTopGlobalPlayers() ([]models.User, error) {
	var users []models.User

	err := repo.DB.
		Joins("INNER JOIN tournament_participants tp ON users.id = tp.user_id").
		Order("users.level DESC").
		Limit(1000).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
