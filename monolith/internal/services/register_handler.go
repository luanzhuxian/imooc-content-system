package services

import (
	"fmt"
	"imooc-content-system/internal/dao"
	"imooc-content-system/internal/model"
	"imooc-content-system/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type RegisterRsp struct {
	Message string `json:"message" binding:"required"`
}

func (c *CmsApp) Register(ctx *gin.Context) {
	var req RegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 加密密码
	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("hashedPassword = %s\n", hashedPassword)
	// 账号校验
	accountDao := dao.NewAccountDao(c.db)
	isExist, err := accountDao.IsExist(req.UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if isExist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account already exists"})
		return
	}
	// 持久化
	currentTime := time.Now()
	account := model.Account{
		UserID:   req.UserID,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Ct:       currentTime,
		Ut:       currentTime,
	}
	err = accountDao.Create(account)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("register req = %+v , hashedPassword = [%s]\n", req, hashedPassword)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &RegisterRsp{
			Message: fmt.Sprintf("register success, user_id: %s, nickname: %s", req.UserID, req.Nickname),
		},
	})
}
