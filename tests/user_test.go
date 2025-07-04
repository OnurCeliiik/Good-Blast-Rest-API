package tests

import (
	"bytes"
	"encoding/json"
	"good-api/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	router := SetupRouter()

	newUser := models.User{
		Username: "created_user",
		Level:    10,
		Coins:    1000,
		Country:  "Turkey",
	}

	payload, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/users/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var created models.User
	err := json.Unmarshal(rec.Body.Bytes(), &created)
	assert.NoError(t, err)
	assert.Equal(t, newUser.Username, created.Username)
	check := assert.NotEmpty(t, created.ID)
	if check != true {
		log.Fatalln("Problem")
	}
}

func TestGetUser(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	user, _ := SeedTestData(db)

	req, _ := http.NewRequest("GET", "/users/"+user.ID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result models.User
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, user.ID, result.ID)
	check := assert.Equal(t, user.Username, result.Username)
	if check != true {
		log.Fatalln("Problem")
	}
}

func TestGetAllUsers(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	SeedTestData(db)

	req, _ := http.NewRequest("GET", "/users/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var users []models.User
	err := json.Unmarshal(rec.Body.Bytes(), &users)
	assert.NoError(t, err)
	check := assert.GreaterOrEqual(t, len(users), 1)
	if check != true {
		log.Fatalln("Problem!")
	}
}

func TestUpdateUser(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	user, _ := SeedTestData(db)

	update := models.User{
		Level: 30,
		Coins: 2000,
	}

	payload, _ := json.Marshal(update)

	req, _ := http.NewRequest("PUT", "/users/"+user.ID.String(), bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var updated models.User
	db.First(&updated, "id = ?", user.ID)
	assert.Equal(t, 30, updated.Level)
	check := assert.Equal(t, 2000, updated.Coins)
	if check != true {
		log.Fatalln("Problem")
	}
}

func TestDeleteUser(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()
	user, _ := SeedTestData(db)

	req, _ := http.NewRequest("DELETE", "/users/"+user.ID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var deleted models.User
	err := db.First(&deleted, "id = ?", user.ID).Error
	check := assert.Error(t, err)
	if check != true {
		log.Fatalln("Problem")
	}
}
