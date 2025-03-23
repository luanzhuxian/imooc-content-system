package controllers

import (
	"imooc-content-system/internal/utils"
	"net/http"

	"imooc-content-system/internal/services"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountService *services.AccountService
}

// 创建实例 依赖注入
func NewAccountController(accountService *services.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

type RegisterRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type RegisterResponse struct {
	Message string `json:"message" binding:"required"`
}

// Controller负责HTTP处理
func (c *AccountController) Register(ctx *gin.Context) {
	// 1. 解析和验证请求
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 调用Service层
	err := c.accountService.RegisterUser(ctx, req.UserID, req.Password, req.Nickname)

	// 3. 处理结果并返回响应
	if err != nil {
		// 根据不同错误类型返回不同状态码
		if err == utils.ErrUserAlreadyExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account already exists"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &RegisterResponse{
			Message: "register success, user_id: " + req.UserID + ", nickname: " + req.Nickname,
		},
	})
}

func (c *AccountController) FindByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	account, err := c.accountService.FindUser(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": account,
	})
}
