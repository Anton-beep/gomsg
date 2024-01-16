package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gomsg/pkg/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetChats(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	user.Token, _ = api.HashPassword("token")
	userID, _ := dbAPI.CreateNewUser(user)

	chat := getTimestampChat()
	chat.UsersIDs = append(chat.UsersIDs, userID)
	chatID, _ := dbAPI.CreateNewChat(chat)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\"}",
		userID, user.Token))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/getChats", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	fmt.Println(recorder.Body.String())

	res, err := dbAPI.DeleteChat(chatID)
	assert.NoError(t, err)
	assert.True(t, res)
	res, err = dbAPI.DeleteUserByUsername(user.Username)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestGetMessagesByChatIDAPI(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	chat := getTimestampChat()
	chat.UsersIDs = append(chat.UsersIDs, userID)
	chatID, err := dbAPI.CreateNewChat(chat)
	assert.NoError(t, err)

	message1 := getTimestampMessage()
	message1.ChatID = chatID
	message1.SenderID = userID
	message1ID, err := dbAPI.CreateNewMessage(message1)
	assert.NoError(t, err)

	message2 := getTimestampMessage()
	message2.ChatID = chatID
	message2.SenderID = userID
	message2ID, err := dbAPI.CreateNewMessage(message2)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"chatID\": %v}",
		userID, user.Token, chatID))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/getMessagesByChatID", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.True(t, strings.Contains(recorder.Body.String(), message1.Text))
	assert.True(t, strings.Contains(recorder.Body.String(), message2.Text))

	res, err := dbAPI.DeleteMessage(message1ID)
	assert.NoError(t, err)
	assert.True(t, res)
	res, err = dbAPI.DeleteMessage(message2ID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteChat(chatID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByUsername(user.Username)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestGetInfoUser(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	userDest := getTimestampUser()
	userDestID, err := dbAPI.CreateNewUser(userDest)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"destUserID\": %v}",
		userID, user.Token, userDestID))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/getInfoUser", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.True(t, strings.Contains(recorder.Body.String(), userDest.Username))

	res, err := dbAPI.DeleteUserByUsername(user.Username)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByUsername(userDest.Username)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestEditMessageAPI(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	chat := getTimestampChat()
	chat.UsersIDs = append(chat.UsersIDs, userID)
	chatID, err := dbAPI.CreateNewChat(chat)
	assert.NoError(t, err)

	message := getTimestampMessage()
	message.ChatID = chatID
	message.SenderID = userID
	messageID, err := dbAPI.CreateNewMessage(message)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"messageID\": %v, \"newText\": \"newText\"}",
		userID, user.Token, messageID))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/editMessage", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	editedMessage, err := dbAPI.GetMessageByID(messageID)
	assert.NoError(t, err)
	assert.Equal(t, "newText", editedMessage.Text)

	res, err := dbAPI.DeleteMessage(messageID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteChat(chatID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByUsername(user.Username)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestEditStatusApi(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"newStatus\": \"newStatus\"}",
		userID, user.Token))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/editStatus", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	userDB, err := dbAPI.GetUserByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, "newStatus", userDB.Status)

	res, err := dbAPI.DeleteUserByUsername(user.Username)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestCreateNewMessageApi(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	chat := getTimestampChat()
	chatID, err := dbAPI.CreateNewChat(chat)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"chatID\": %v, \"text\": \"newMessage\"}",
		userID, user.Token, chatID))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/createMessage", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	fmt.Println(recorder.Body)
	assert.Equal(t, http.StatusOK, recorder.Code)
	messages, err := dbAPI.GetMessagesByChatID(chatID, 1, int(time.Now().Unix()))
	assert.NoError(t, err)
	assert.Equal(t, "newMessage", messages[0].Text)

	res, err := dbAPI.DeleteMessage(messages[0].MessageID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteChat(chatID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByID(userID)
	assert.NoError(t, err)
	assert.True(t, res)
}

type answerCreateChat struct {
	ChatID int `json:"chatID"`
}

func TestCreateChatAPI(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"chatName\": \"chatName\", \"usersIDs\": [%v]}",
		userID, user.Token, userID))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/createChat", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var answer answerCreateChat
	err = json.Unmarshal(recorder.Body.Bytes(), &answer)
	assert.NoError(t, err)

	chat, err := dbAPI.GetChatByID(answer.ChatID)
	assert.NoError(t, err)
	assert.Equal(t, "chatName", chat.ChatName)

	res, err := dbAPI.DeleteChat(chat.ChatID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByID(userID)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestGetMessageUpdatesAPI(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	chat := getTimestampChat()
	chatID, err := dbAPI.CreateNewChat(chat)
	assert.NoError(t, err)

	message := getTimestampMessage()
	message.ChatID = chatID
	message.SenderID = userID
	messageID, err := dbAPI.CreateNewMessage(message)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"timestamp\": %v}",
		userID, user.Token, message.Timestamp))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/getUpdatesMessage", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.True(t, strings.Contains(recorder.Body.String(), message.Text))

	res, err := dbAPI.DeleteMessage(messageID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteChat(chatID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByID(userID)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestIsUserExistsAPI(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"username\": \"%v\"}",
		userID, user.Token, user.Username))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/isUserExists", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.True(t, strings.Contains(recorder.Body.String(), "true"))

	res, err := dbAPI.DeleteUserByID(userID)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestCreateChatByUsernames(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	user2 := getTimestampUser()
	user2ID, err := dbAPI.CreateNewUser(user2)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"chatName\": \"chatName\", \"usernames\": [\"%v\", \"%v\"]}",
		userID, user.Token, user.Username, user2.Username))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/createChatByUsernames", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var answer answerCreateChat
	err = json.Unmarshal(recorder.Body.Bytes(), &answer)
	assert.NoError(t, err)

	chat, err := dbAPI.GetChatByID(answer.ChatID)
	assert.True(t, chat != nil)
	assert.NoError(t, err)
	assert.Equal(t, "chatName", chat.ChatName)

	res, err := dbAPI.DeleteChat(chat.ChatID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByID(userID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByID(user2ID)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestGetUsersInChatAPI(t *testing.T) {
	router, dbAPI := setupApi(t)

	user := getTimestampUser()
	userID, err := dbAPI.CreateNewUser(user)
	assert.NoError(t, err)

	user2 := getTimestampUser()
	user2ID, err := dbAPI.CreateNewUser(user2)
	assert.NoError(t, err)

	chat := getTimestampChat()
	chat.UsersIDs = make([]int, 0)
	chat.UsersIDs = append(chat.UsersIDs, userID)
	chat.UsersIDs = append(chat.UsersIDs, user2ID)
	chat.UserNames = make([]string, 0)
	chat.UserNames = append(chat.UserNames, user.Username)
	chat.UserNames = append(chat.UserNames, user2.Username)
	chatID, err := dbAPI.CreateNewChat(chat)
	assert.NoError(t, err)

	reqBody := strings.NewReader(fmt.Sprintf("{\"userID\": %v, \"token\": \"%v\", \"chatID\": %v}",
		userID, user.Token, chatID))
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/private/getInfoChat", reqBody)
	assert.NoError(t, err)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	fmt.Println(recorder.Body.String())
	assert.True(t, strings.Contains(recorder.Body.String(), user.Username))
	assert.True(t, strings.Contains(recorder.Body.String(), user2.Username))

	res, err := dbAPI.DeleteChat(chatID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByID(userID)
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = dbAPI.DeleteUserByID(user2ID)
	assert.NoError(t, err)
	assert.True(t, res)
}
