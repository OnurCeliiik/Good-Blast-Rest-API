package handlers

import (
	"good-api/internal/repositories"
	"good-api/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LeaderboardHandler struct {
	LeaderboardService    *services.LeaderboardService
	LeaderboardRepository *repositories.LeaderboardRepository
}

// NewLeaderboardHandler initializes the leaderboard handler.
func NewLeaderboardHandler(ls *services.LeaderboardService, lr *repositories.LeaderboardRepository) *LeaderboardHandler {
	return &LeaderboardHandler{
		LeaderboardService:    ls,
		LeaderboardRepository: lr,
	}
}

// @Summary Get Global Leaderboard
// @Description It gets the top 1000 users from the global leaderboard
// @Tags Leaderboards
// @Accept json
// @Produce json
// @Success 201 {object} []models.User
// @Failure 400 {object} map[string]string
// @Router /leaderboard/ [get]
func (h *LeaderboardHandler) GetGlobalLeaderboard(c *gin.Context) {
	leaderboard, err := h.LeaderboardRepository.GetGlobalLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
	// GetGlobalLeaderboard handles GET /leaderboard/global?limit=1000
}

// @Summary Get Country Leaderboard
// @Description It gets the top users from the specified country
// @Tags Leaderboards
// @Accept json
// @Produce json
// @Success 201 {object} []models.User
// @Failure 400 {object} map[string]string
// @Router /leaderboard/ [get]
func (h *LeaderboardHandler) GetCountryLeaderboard(c *gin.Context) {
	country := c.Query("country")
	if country == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Country is required"})
		return
	}

	leaderboard, err := h.LeaderboardRepository.GetCountryLeaderboard(country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
	// GetCountryLeaderboard handles GET /leaderboard/country?country=Turkey&limit=1000
}

// @Summary Get Tourmament Leaderboard
// @Description It gets the top players from the specified tournament
// @Tags Leaderboards
// @Accept json
// @Produce json
// @Success 201 {object} []models.User
// @Failure 400 {object} map[string]string
// @Router /leaderboard/ [get]
func (h *LeaderboardHandler) GetTournamentLeaderboard(c *gin.Context) {
	tournamentIDParam := c.Query("tournament_id")
	if tournamentIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tournament ID is required"})
		return
	}

	limitParam := c.DefaultQuery("limit", "1000")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	leaderboard, err := h.LeaderboardService.GetTournamentLeaderboard(tournamentIDParam, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
	// GetTournamentLeaderboard handles GET /leaderboard/tournament?tournament_id=xyz&limit=1000
}

// @Summary Get Tournament Rank
// @Description It gets the rank of the user from the tournament they are in
// @Tags Leaderboards
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /leaderboard/ [get]
func (h *LeaderboardHandler) GetTournamentRank(c *gin.Context) {
	userIDParam := c.Query("user_id")
	tournamentIDParam := c.Query("tournament_id")
	if userIDParam == "" || tournamentIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Tournament ID are required"})
		return
	}

	rank, err := h.LeaderboardService.GetTournamentRank(userIDParam, tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rank": rank})
	// GetTournamentRank handles GET /leaderboard/tournament/rank?user_id=xyz&tournament_id=xyz
}
