package repositories

import (
	"errors"
	"good-api/internal/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
Responsible for database operations (CRUD)
Acts as a middleman between the database(db.go) and the business logic(services)

Repository layer isolates database logic so that other parts of the app
don't directly interact with gorm.DB
*/

/*
Define a struct for the repository
This struct holds a database connection
We use a pointer (*gorm.DB) so wedon't copy the database object every time.
*/
type UserRepository struct {
	DB *gorm.DB
}

// This function initializes the repository and stores the db connection inside it.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create a user
func (repo *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := repo.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Find a user by ID
func (repo *UserRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := repo.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername fetches a user by their username
func (repo *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := repo.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		log.Printf("unexpected error fetching user: %v", err)
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := repo.DB.Find(&users).Error
	return users, err
}

// Update a User
func (repo *UserRepository) UpdateUser(user *models.User) (*models.User, error) {
	if err := repo.DB.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Delete a User
func (repo *UserRepository) DeleteUser(userID uuid.UUID) error {
	return repo.DB.Delete(&models.User{}, userID).Error // Deletes the user by id.
}

// AddCoins adds coins to a user's balance.
func (repo *UserRepository) AddCoins(userID uuid.UUID, amount int) error {
	var user models.User
	if err := repo.DB.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	user.Coins += amount

	if err := repo.DB.Save(&user).Error; err != nil {
		return errors.New("failed to update user balance")
	}

	return nil
}

func (repo *UserRepository) GetUserTournament(userID uuid.UUID) (*models.Tournament, error) {
	var participant models.TournamentParticipant

	// Check if user is in a tournament
	err := repo.DB.Where("user_id = ?", userID).First(&participant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var tournament models.Tournament
	err = repo.DB.Where("id = ? ", participant.TournamentID).First(&tournament).Error
	if err != nil {
		return nil, err
	}
	return &tournament, err
}
