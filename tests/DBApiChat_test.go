package tests

import (
	"github.com/stretchr/testify/assert"
	"gomsg/pkg/db"
	"gomsg/pkg/models"
	"slices"
	"strconv"
	"testing"
	"time"
)

func getTimestampChat() models.Chat {
	var newChat models.Chat
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))
	newChat.ChatName = "chatName" + timestamp
	newChat.UsersIDs = []int{-1, -2}
	return newChat
}

func TestCreateAndDeleteChat(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newChat := getTimestampChat()
	newId, err := dbApi.CreateNewChat(newChat)
	assert.NoError(t, err)

	res, err := dbApi.DeleteChat(newId)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteUserByUsername(newChat.ChatName)
	assert.NoError(t, err)
	assert.Equal(t, false, res)
}

func TestGetChatsByUserID(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newChat := getTimestampChat()
	newID, err := dbApi.CreateNewChat(newChat)
	assert.NoError(t, err)

	chats, err := dbApi.GetChatsByUserID(newChat.UsersIDs[0])
	assert.NoError(t, err)
	chatsIDs := make([]int, 0, len(chats))
	for _, chat := range chats {
		chatsIDs = append(chatsIDs, chat.UsersIDs...)
	}
	assert.True(t, slices.Contains(chatsIDs, newChat.UsersIDs[0]))

	chats, err = dbApi.GetChatsByUserID(-3)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(chats))

	res, err := dbApi.DeleteChat(newID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestEditChatNameByChatID(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newChat := getTimestampChat()
	newID, err := dbApi.CreateNewChat(newChat)
	assert.NoError(t, err)

	res, err := dbApi.EditChatNameByChatID(newID, "newChatName")
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	dbChat, err := dbApi.GetChatByID(newID)
	assert.NoError(t, err)
	assert.Equal(t, "newChatName", dbChat.ChatName)

	res, err = dbApi.DeleteChat(newID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestEditUserIDsByChatID(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newChat := getTimestampChat()
	newID, err := dbApi.CreateNewChat(newChat)
	assert.NoError(t, err)

	res, err := dbApi.AddUserToChat(newID, -3)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	dbChat, err := dbApi.GetChatByID(newID)
	assert.NoError(t, err)
	assert.Equal(t, []int{-1, -2, -3}, dbChat.UsersIDs)

	err = dbApi.RemoveUserFromChat(newID, -3)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteChat(newID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}
