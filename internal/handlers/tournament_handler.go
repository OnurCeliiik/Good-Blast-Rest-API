package handlers

import (
	"good-api/internal/repositories"
	"good-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TournamentHandler struct {
	TournamentService *services.TournamentService
	TournamentRepo    *repositories.TournamentRepository
	UserRepo          *repositories.UserRepository
}

// NewTournamentHandler creates a new TournamentHandler.
func NewTournamentHandler(ts *services.TournamentService, tr *repositories.TournamentRepository) *TournamentHandler {
	return &TournamentHandler{
		TournamentService: ts,
		TournamentRepo:    tr,
	}
}

// @Summary Enter Tournament
// @Description It enters the user to the tournament
// @Tags Tournaments
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tournaments/ [post]
func (h *TournamentHandler) EnterTournament(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	tournament, err := h.TournamentService.EnterTournament(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "User entered the tournament successfully",
		"tournament_id": tournament.ID,
		"start_time":    tournament.StartTime,
		"end_time":      tournament.EndTime,
		"user_count":    tournament.UserCount,
	})
}

// @Summary Get All Tournament
// @Description It gets all tournaments
// @Tags Tournaments
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tournaments/ [get]
func (h *TournamentHandler) GetAllTournaments(c *gin.Context) {
	tournaments, err := h.TournamentRepo.GetAllTournaments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tournaments": tournaments})
}

// @Summary Get Tournament
// @Description It gets a single tournament
// @Tags Tournaments
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tournaments/ [get]
func (h *TournamentHandler) GetTournament(c *gin.Context) {
	tournamentIDStr := c.Param("id")
	tournamentID, err := uuid.Parse(tournamentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID format"})
		return
	}
	tournament, err := h.TournamentService.GetTournamentByID(tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	c.JSON(http.StatusOK, tournament)
}

// @Summary Finish Tournament
// @Description It finishes a single tournament
// @Tags Tournaments
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tournaments/ [post]
func (h *TournamentHandler) FinishTournament(c *gin.Context) {
	tournamentIDStr := c.Param("id")
	tournamentID, err := uuid.Parse(tournamentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID format"})
		return
	}

	err = h.TournamentService.FinishTournament(tournamentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish tournament"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tournament finished successfully"})
}

// @Summary Finish All Tournaments
// @Description It finishes all tournaments
// @Tags Tournaments
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tournaments/ [post]
func (h *TournamentHandler) FinishAllTournaments(c *gin.Context) {
	err := h.TournamentService.FinishAllTournaments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish all tournaments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All tournaments finished successfully"})
}

// @Summary Update Score
// @Description It updates the score of the user.
// @Tags Tournaments
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tournaments/ [put]
func (h *TournamentHandler) UpdateScore(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = h.TournamentService.UpdateScore(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Score updated successfully"})
}
