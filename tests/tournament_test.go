package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnterTournament(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	user, _ := SeedTestData(db)

	// First request for entering the tournament will be 200 given the conditions are met.
	req1, _ := http.NewRequest("POST", "/tournaments/enter/"+user.ID.String(), nil)
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)

	assert.Equal(t, http.StatusOK, rec1.Code, "First is good to go")

	/*
		// If user has entered a tournament before, no further entry is allowed until the existing tournament is concluded.
		req2, _ := http.NewRequest("POST", "/tournaments/enter/"+user.ID.String(), nil)
		rec2 := httptest.NewRecorder()
		router.ServeHTTP(rec2, req2)

		assert.Equal(t, http.StatusBadRequest, rec2.Code, "Second is no no")
	*/
}

func TestGetTournament(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	_, tournament := SeedTestData(db)

	req, _ := http.NewRequest("GET", "/tournaments/"+tournament.ID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllTournaments(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	SeedTestData(db)

	req, _ := http.NewRequest("GET", "/tournaments/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestUpdateScore(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	user, _ := SeedTestData(db)

	http.NewRequest("PUT", "/tournaments/enter/"+user.ID.String(), nil)
	req, _ := http.NewRequest("POST", "/tournaments/update-score/"+user.ID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestFinishTournament(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	_, tournament := SeedTestData(db)

	req, _ := http.NewRequest("POST", "/tournaments/finish/"+tournament.ID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestFinishAllTournaments(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()

	req, _ := http.NewRequest("POST", "/tournaments/finish-all", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

}
