package tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gomsg/pkg/api"
	"gomsg/pkg/db"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupApi(t *testing.T) (*gin.Engine, *db.APIDB) {
	gin.SetMode(gin.TestMode)
	newDb, err := db.NewDb()
	assert.NoError(t, err)
	newAPI := api.NewApi(newDb)
	return newAPI.Start(), newDb
}

func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := api.HashPassword(password)
	assert.NoError(t, err)
	assert.True(t, api.CheckPasswordHash(password, hashedPassword))
}

func TestPing(t *testing.T) {
	router, _ := setupApi(t)

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/ping", nil)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"message":"pong"}`, recorder.Body.String())
}

func TestRegister(t *testing.T) {
	router, dbAPI := setupApi(t)

	// delete user if exists
	dbAPI.DeleteUserByUsername("user-1")

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/register", strings.NewReader(
		"{\"username\": \"user-1\", \"password\": \"user-1\"}"),
	)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// delete user after test
	dbAPI.DeleteUserByUsername("user-1")
}

func TestLogin(t *testing.T) {
	router, dbAPI := setupApi(t)

	// create user if not exists
	userRawPassword := "password"

	user := getTimestampUser()
	user.Token, _ = api.HashPassword(userRawPassword)
	dbAPI.CreateNewUser(user)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/login", strings.NewReader(
		fmt.Sprintf("{\"username\": \"%v\", \"password\": \"%v\"}", user.Username, userRawPassword),
	))
	router.ServeHTTP(recorder, req)

	fmt.Println(recorder.Body.String())
	assert.Equal(t, http.StatusOK, recorder.Code)

	// delete user after test
	dbAPI.DeleteUserByUsername(user.Username)
}
