package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gomsg/pkg/models"
	"net/http"
	"slices"
	"time"
)

// not private

func (a *API) Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type registerReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *API) Register(c *gin.Context) {
	var data registerReq
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		zap.L().Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := a.db.GetUserByUsername(data.Username)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user != nil {
		zap.L().Debug("user with this username already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user with this username already exists"})
		return
	}

	var newUser models.User
	newUser.Username = data.Username
	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	newUser.Token = hashedPassword

	newUserID, err := a.db.CreateNewUser(newUser)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"userID":  newUserID,
		"token":   hashedPassword,
	})
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *API) Login(c *gin.Context) {
	var data loginReq
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		zap.L().Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := a.db.GetUserByUsername(data.Username)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		zap.L().Debug("user with this username does not exist")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user with this username does not exist"})
		return
	}

	if !(CheckPasswordHash(data.Password, user.Token)) {
		zap.L().Debug("wrong password")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"userID":  user.UserID,
		"token":   user.Token})
}

// private

type getUsersChatsRequest struct {
	UserID int `json:"userID" binding:"required"`
}

func (a *API) GetUsersChats(c *gin.Context) {
	var data getUsersChatsRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		zap.L().Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := a.db.GetUserByID(data.UserID)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		zap.L().Debug("user with this userID does not exist")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user with this userID does not exist"})
		return
	}

	chats, err := a.db.GetChatsByUserID(user.UserID)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"chats":   chats})
}

type getMessagesByChatID struct {
	UserID    int `json:"userID" binding:"required"`
	ChatID    int `json:"chatID" binding:"required"`
	Quantity  int `json:"quantity" binding:"-"`
	Timestamp int `json:"timestamp" binding:"-"`
}

// GetMessagesByChatID by chatID and quantity (optional, default = 10, max = 50).
// Timestamps of messages are smaller than given (optional, default = now).
func (a *API) GetMessagesByChatID(c *gin.Context) {
	var data getMessagesByChatID
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		zap.L().Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if data.Quantity == 0 {
		data.Quantity = 10
	}
	if data.Timestamp == 0 {
		data.Timestamp = int(time.Now().Unix())
	}

	chat, err := a.db.GetChatByID(data.ChatID)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if chat == nil {
		zap.L().Debug("chat with this chatID does not exist")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "chat with this chatID does not exist"})
		return
	}
	if !(slices.Contains(chat.UsersIDs, data.UserID)) {
		zap.L().Debug("user with this userID is not in this chat")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user with this userID is not in this chat"})
		return
	}

	messages, err := a.db.GetMessagesByChatID(data.ChatID, data.Quantity, data.Timestamp)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "OK",
		"messages": messages},
	)
}

type limitedDataUser struct {
	Username string
	Picture  string
	Status   string
}

type getInfoByUserReq struct {
	DestUserID int `json:"destUserID" binding:"required"`
}

func (a *API) GetInfoAboutUser(c *gin.Context) {
	var data getInfoByUserReq
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		zap.L().Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := a.db.GetUserByID(data.DestUserID)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		zap.L().Debug("user with this destUserID does not exist")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user with this destUserID does not exist"})
		return
	}

	output := limitedDataUser{Username: user.Username, Picture: user.Picture, Status: user.Status}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"user":    output})
}

type editMessageReq struct {
	UserID    int    `json:"userID" binding:"required"`
	MessageID int    `json:"messageID" binding:"required"`
	NewText   string `json:"newText" binding:"required"`
}

func (a *API) EditMessage(c *gin.Context) {
	var data editMessageReq
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		zap.L().Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	msg, err := a.db.GetMessageByID(data.MessageID)
	if err != nil || msg == nil {
		zap.L().Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "message does not exist"})
		return
	}

	if msg.SenderID != data.UserID {
		zap.L().Debug("you are not sender of this message")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you are not sender of this message"})
		return
	}

	res, err := a.db.EditMessage(data.MessageID, data.NewText)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !res {
		zap.L().Error("message with this messageID does not exist")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "message with this messageID does not exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK"},
	)
}

type editStatusReq struct {
	UserID    int    `json:"userID" binding:"required"`
	NewStatus string `json:"newStatus" binding:"required"`
}

func (a *API) EditStatus(c *gin.Context) {
	var data editStatusReq
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	res, err := a.db.EditStatus(data.UserID, data.NewStatus)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !res {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

type createNewMessageReq struct {
	ChatID int    `json:"chatID" binding:"required"`
	UserID int    `json:"userID" binding:"required"`
	Text   string `json:"text" binding:"required"`
}

func (a *API) CreateNewMessage(c *gin.Context) {
	var data createNewMessageReq
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		zap.L().Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	message := models.Message{ChatID: data.ChatID, SenderID: data.UserID, Text: data.Text}

	messageID, err := a.db.CreateNewMessage(message)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "OK",
		"messageID": messageID,
	})
}

type createNewChatReq struct {
	UsersIDs []int  `json:"usersIDs" binding:"required"`
	ChatName string `json:"chatName" binding:"required"`
}

func (a *API) CreateNewChat(c *gin.Context) {
	var data createNewChatReq
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		zap.L().Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	chat := models.Chat{UsersIDs: data.UsersIDs, ChatName: data.ChatName}

	chatID, err := a.db.CreateNewChat(chat)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"chatID":  chatID,
	})
}

type getMessageUpdatesReq struct {
	UserID    int `json:"userID" binding:"required"`
	Timestamp int `json:"timestamp" binding:"required"`
}

// GetMessageUpdates timestamp is the moment when user got messages last time
func (a *API) GetMessageUpdates(c *gin.Context) {
	var data getMessageUpdatesReq
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		zap.L().Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messages, err := a.db.GetMessageUpdates(data.UserID, data.Timestamp)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"updates": messages,
	})
}
