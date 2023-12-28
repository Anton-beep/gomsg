package tests

import (
	"github.com/stretchr/testify/assert"
	"gomsg/pkg/db"
	"gomsg/pkg/models"
	"strconv"
	"testing"
	"time"
)

func getTimestampMessage() models.Message {
	var newMessage models.Message
	newMessage.Text = "text" + strconv.Itoa(int(time.Now().UnixNano()))
	newMessage.ChatID = -1
	newMessage.SenderID = -1
	newMessage.Timestamp = int(time.Now().Unix())
	return newMessage
}

func TestCreateNewMessage(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newSender := getTimestampUser()
	newSenderID, err := dbApi.CreateNewUser(newSender)
	assert.NoError(t, err)

	newChat := getTimestampChat()
	newChat.UsersIDs = append(newChat.UsersIDs, newSenderID)
	newChatID, err := dbApi.CreateNewChat(newChat)
	assert.NoError(t, err)

	newMessage := getTimestampMessage()
	newMessage.SenderID = newSenderID
	newMessage.ChatID = newChatID
	newID, err := dbApi.CreateNewMessage(newMessage)
	assert.NoError(t, err)

	res, err := dbApi.DeleteMessage(newID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteMessage(newID)
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	res, err = dbApi.DeleteChat(newChatID)
	assert.NoError(t, err)

	res, err = dbApi.DeleteUserByUsername(newSender.Username)
	assert.NoError(t, err)
}

func TestGetMessagesByChatID(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newSender := getTimestampUser()
	newSenderID, err := dbApi.CreateNewUser(newSender)
	assert.NoError(t, err)

	newChat := getTimestampChat()
	newChat.UsersIDs = append(newChat.UsersIDs, newSenderID)
	newChatID, err := dbApi.CreateNewChat(newChat)
	assert.NoError(t, err)

	newMessage := getTimestampMessage()
	newMessage.SenderID = newSenderID
	newMessage.ChatID = newChatID
	newID, err := dbApi.CreateNewMessage(newMessage)
	assert.NoError(t, err)

	messages, err := dbApi.GetMessagesByChatID(newMessage.ChatID, 1, int(time.Now().Unix()))
	assert.NoError(t, err)
	assert.Equal(t, newMessage.Text, messages[0].Text)

	res, err := dbApi.DeleteMessage(newID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteChat(newChatID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteUserByUsername(newSender.Username)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestEditMessage(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newSender := getTimestampUser()
	newSenderID, err := dbApi.CreateNewUser(newSender)
	assert.NoError(t, err)

	newChat := getTimestampChat()
	newChat.UsersIDs = append(newChat.UsersIDs, newSenderID)
	newChatID, err := dbApi.CreateNewChat(newChat)
	assert.NoError(t, err)

	newMessage := getTimestampMessage()
	newMessage.SenderID = newSenderID
	newMessage.ChatID = newChatID
	newID, err := dbApi.CreateNewMessage(newMessage)
	assert.NoError(t, err)

	res, err := dbApi.EditMessage(newID, "newText")
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	dbMessage, err := dbApi.GetMessageByID(newID)
	assert.NoError(t, err)
	assert.Equal(t, "newText", dbMessage.Text)

	res, err = dbApi.DeleteMessage(newID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteChat(newChatID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteUserByUsername(newSender.Username)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestGetUpdatesMessages(t *testing.T) {
	dbApi, err := db.NewDb()
	assert.NoError(t, err)

	newSender := getTimestampUser()
	newSenderID, err := dbApi.CreateNewUser(newSender)
	assert.NoError(t, err)

	newChat := getTimestampChat()
	newChat.UsersIDs = append(newChat.UsersIDs, newSenderID)
	newChatID, err := dbApi.CreateNewChat(newChat)
	assert.NoError(t, err)

	newMessage := getTimestampMessage()
	newMessage.SenderID = newSenderID
	newMessage.ChatID = newChatID
	newID, err := dbApi.CreateNewMessage(newMessage)
	assert.NoError(t, err)

	messages, err := dbApi.GetMessageUpdates(newSenderID, newMessage.Timestamp-10)
	assert.NoError(t, err)
	assert.Equal(t, newMessage.Text, messages[0].Text)

	res, err := dbApi.DeleteMessage(newID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteChat(newChatID)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = dbApi.DeleteUserByUsername(newSender.Username)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}
