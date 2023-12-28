package api

import (
	"github.com/gin-gonic/gin"
	"gomsg/pkg/db"
)

type API struct {
	db *db.APIDB
}

func NewApi(_db *db.APIDB) *API {
	newApi := &API{_db}
	return newApi
}

func (a *API) Start() *gin.Engine {
	router := gin.Default()

	router.GET("/api/ping", a.Pong)
	router.GET("/api/register", a.Register)
	router.GET("/api/login", a.Login)

	authGroup := router.Group("/api/private")
	authGroup.Use(a.AuthMiddleware())

	authGroup.GET("getChats", a.GetUsersChats)
	authGroup.GET("getMessagesByChatID", a.GetMessagesByChatID)
	authGroup.GET("getInfoUser", a.GetInfoAboutUser)
	authGroup.GET("editMessage", a.EditMessage)
	authGroup.GET("editStatus", a.EditStatus)
	authGroup.GET("createMessage", a.CreateNewMessage)
	authGroup.GET("createChat", a.CreateNewChat)
	authGroup.GET("getUpdatesMessage", a.GetMessageUpdates)

	return router
}
