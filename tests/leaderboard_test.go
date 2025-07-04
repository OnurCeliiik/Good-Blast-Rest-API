package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeaderboardGlobal(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/leaderboard/global", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println("All good mate")
}

func TestLeaderboardCountry(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/leaderboard/country?country=Turkey", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println("All good mate")
}

func TestTournamentLeaderboard(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	_, tournament := SeedTestData(db)

	url := fmt.Sprintf("/leaderboard/tournament?tournament_id=%s", tournament.ID.String())
	req, _ := http.NewRequest("GET", url, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println("All good mate")
}

func TestGetTournamentRank(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	user, tournament := SeedTestData(db)

	url := fmt.Sprintf("/leaderboard/tournament/rank?user_id=%s&tournament_id=%s", user.ID.String(), tournament.ID.String())
	req, _ := http.NewRequest("GET", url, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println("All good mate")
}
