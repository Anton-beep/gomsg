package api

import (
	"github.com/gin-contrib/cors"
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

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/api/ping", a.Pong)
	router.POST("/api/register", a.Register)
	router.POST("/api/login", a.Login)

	authGroup := router.Group("/api/private")
	authGroup.Use(a.AuthMiddleware())

	authGroup.POST("getChats", a.GetUsersChats)
	authGroup.POST("getMessagesByChatID", a.GetMessagesByChatID)
	authGroup.POST("getInfoUser", a.GetInfoAboutUser)
	authGroup.POST("editMessage", a.EditMessage)
	authGroup.POST("editStatus", a.EditStatus)
	authGroup.POST("createMessage", a.CreateNewMessage)
	authGroup.POST("createChat", a.CreateNewChat)
	authGroup.POST("getUpdatesMessage", a.GetMessageUpdates)

	return router
}
