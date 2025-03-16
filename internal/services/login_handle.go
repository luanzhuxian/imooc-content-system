package services

import (
	"context"
	"fmt"
	"imooc-content-system/internal/dao"
	"imooc-content-system/internal/utils"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginReq struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRsp struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Nickname  string `json:"nickname"`
}

func (c *CmsApp) Login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var (
		userID   = req.UserID
		password = req.Password
	)
	// 实例化dao
	accountDao := dao.NewAccountDao(c.db)
	account, err := accountDao.FirstByUserID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "please input correct account id"})
		return
	}
	// 对密码加密比较
	if err := bcrypt.CompareHashAndPassword(
		[]byte(account.Password),
		[]byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "please input correct password"})
		return
	}
	sessionID, err := c.generateSessionID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &LoginRsp{
			SessionID: sessionID,
			UserID:    account.UserID,
			Nickname:  account.Nickname,
		},
	})
	return
}

func (c *CmsApp) generateSessionID(ctx context.Context, userID string) (string, error) {
	// 登录成功后，在Redis中创建一个临时的会话令牌，用于用户认证和会话管理
	sessionID := uuid.New().String()
	// key : session_id:{user_id} val : session_id  20s
	sessionKey := utils.GetSessionKey(userID)
	err := c.rdb.Set(ctx, sessionKey, sessionID, time.Hour*8).Err()
	if err != nil {
		fmt.Printf("rdb set error = %v \n", err)
		return "", err
	}
	authKey := utils.GetAuthKey(sessionID)
	// 为当前会话设置一个过期时间，值为当前时间戳，过期时间为 8 小时，用来判断会话是否过期
	err = c.rdb.Set(ctx, authKey, time.Now().Unix(), time.Hour*8).Err() 
	if err != nil {
		fmt.Printf("rdb set error = %v \n", err)
		return "", err
	}
	return sessionID, nil
}
