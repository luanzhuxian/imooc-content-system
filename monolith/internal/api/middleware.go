package api

import (
	"net/http"

	"imooc-content-system/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const SessionKey = "session_id"

type SessionAuthMiddleware struct {
	rdb *redis.Client
}

func (m *SessionAuthMiddleware) Auth(ctx *gin.Context) {
	sessionID := ctx.GetHeader(SessionKey)
	if sessionID == "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, "session is id null")
	}
	authKey := utils.GetAuthKey(sessionID)
	loginTime, err := m.rdb.Get(ctx, authKey).Result()
	if err != nil && err != redis.Nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "session auth error")
	}
	if loginTime == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "session is expired")
	}
	ctx.Next()
}
