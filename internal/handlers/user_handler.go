package handlers

import (
	"good-api/internal/models"
	"good-api/internal/repositories"
	"good-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
This is the part that interacts with incoming HTTP requests.
It receives the request from the router.
Extracts data from the request. (JSON body or query parameters)
calls the service layer to process the request
Returns a proper response. (JSON + status code)
*/

type UserHandler struct {
	// UserService *services.UserService
	UserRepo    *repositories.UserRepository
	UserService *services.UserService
}

// First initializer for GET requests, uses only repository
func NewUserHandlerwithRepo(userRepo *repositories.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

// Second initializer, for operations that need business logic
func NewUserHandlerwithService(userRepo *repositories.UserRepository, userService *services.UserService) *UserHandler {
	return &UserHandler{
		UserRepo:    userRepo,
		UserService: userService,
	}
}

// @Summary Create user
// @Description Creates a new user with username, country, level, coins and ID
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/ [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	createdUser, err := h.UserService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// @Summary Get user
// @Description Gets the user with all info it has
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/ [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDstr := c.Param("id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	user, err := h.UserRepo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Gets all users
// @Description Gets all users with their info
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} []models.User
// @Failure 400 {object} map[string]string
// @Router /users/ [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.UserRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "There are no users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Update user
// @Description Updates the user with new information
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/ [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	if h.UserService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserService is not initialized"})
		return
	}

	userIDstr := c.Param("id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid input"})
		return
	}
	updateData.ID = userID

	updatedUser, err := h.UserService.UpdateUser(&updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "anlamsiz hata"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// @Summary Delete user
// @Description Deletes the user completely
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /users/ [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.UserService.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) IncreaseLevel(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.UserService.IncreaseLevel(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User level increased"})
}
