package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
)

type authReq struct {
	UserID int    `json:"userID" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func (a *API) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data authReq
		if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		user, err := a.db.GetUserByID(data.UserID)
		if err != nil {
			zap.L().Error(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		if data.Token != user.Token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "wrong token"})
			return
		}
		c.Next()
	}
}
