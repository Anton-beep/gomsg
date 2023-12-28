package tests

import (
	"github.com/stretchr/testify/assert"
	"gomsg/pkg/db"
	"gomsg/pkg/models"
	"strconv"
	"testing"
	"time"
)

func getTimestampUser() models.User {
	var newUser models.User
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))
	newUser.Username = "user-" + timestamp
	newUser.Token = "user-" + timestamp
	return newUser
}

func TestCreateAndDeleteUser(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newUser := getTimestampUser()

	_, err = dbApi.CreateNewUser(newUser)
	assert.NoError(t, err)

	res, err := dbApi.DeleteUserByUsername(newUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteUserByUsername(newUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, false, res)
}

func TestGetUser(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newUser := getTimestampUser()
	_, err = dbApi.CreateNewUser(newUser)
	assert.NoError(t, err)

	// by username
	dbUser, err := dbApi.GetUserByUsername(newUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, newUser.Username, dbUser.Username)
	assert.Equal(t, newUser.Token, dbUser.Token)

	// by id
	dbUser2, err := dbApi.GetUserByID(dbUser.UserID)
	assert.NoError(t, err)
	assert.Equal(t, dbUser2, dbUser)

	// delete user after test
	_, err = dbApi.DeleteUserByUsername(newUser.Username)
	assert.NoError(t, err)
}

func TestEditStatus(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newUser := getTimestampUser()
	userID, err := dbApi.CreateNewUser(newUser)
	assert.NoError(t, err)

	newStatus := "new status"
	res, err := dbApi.EditStatus(userID, newStatus)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	dbUser, err := dbApi.GetUserByUsername(newUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, newStatus, dbUser.Status)

	res, err = dbApi.DeleteUserByUsername(newUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestEditPicture(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newUser := getTimestampUser()
	userID, err := dbApi.CreateNewUser(newUser)
	assert.NoError(t, err)

	newPicture := "new picture"
	res, err := dbApi.EditPicture(userID, newPicture)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	dbUser, err := dbApi.GetUserByUsername(newUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, newPicture, dbUser.Picture)

	res, err = dbApi.DeleteUserByUsername(newUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}
