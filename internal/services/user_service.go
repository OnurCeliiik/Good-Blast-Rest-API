package services

import (
	"errors"
	"fmt"
	"good-api/internal/cache"
	"good-api/internal/models"
	"good-api/internal/repositories"

	"github.com/google/uuid"
)

/*
Service layer is responsible for business logic.
It processes data before passing it to the repository or returning it to the handler.
It sits between the handler and the repository to ensure separation of concerns.
*/

// Service calls the repository to get or modify data.

type UserService struct {
	repo *repositories.UserRepository // Uses the repository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{repo: userRepo}
}

// CreateUser validates and creates a new user.
func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	// Ensure username is unique
	existingUser, _ := s.repo.GetUserByUsername(user.Username)
	if existingUser != nil {
		return nil, errors.New("username is already taken")
	}

	// Assign a new UUID
	user.ID = uuid.New()

	if user.Level == 0 {
		user.Level = 1 // Default value if not provided
	}
	if user.Coins == 0 {
		user.Coins = 1000 // Default value if not provided
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

// UpdateUser updates a user's details.
func (s *UserService) UpdateUser(user *models.User) (*models.User, error) {
	// Ensure user exists
	existingUser, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	fmt.Println("Attempting to update user with ID: ", user.ID)

	existingUser.Level = user.Level
	existingUser.Coins = user.Coins

	updatedUser, err := s.repo.UpdateUser(existingUser)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

// DeleteUser removes a user from the system
func (s *UserService) DeleteUser(userID uuid.UUID) error {
	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.repo.DeleteUser(userID)
}

// IncreaseLevel increments the user's level.
func (s *UserService) IncreaseLevel(userID uuid.UUID) error {
	var user models.User

	if err := s.repo.DB.First(&user, "id = ?", userID).Error; err != nil {
		return errors.New("user not found")
	}

	user.Level += 1
	user.Coins += 100

	if err := s.repo.DB.Save(&user).Error; err != nil {
		return errors.New("failed to update user's level and coins")
	}

	// Sync Redis leaderboard so the user ranking updates
	tournament, err := s.repo.GetUserTournament(userID)
	if err == nil && tournament != nil {
		cache.AddUserToLeaderboard(tournament.ID, user.ID, user.Level)
	}
	return nil
}
