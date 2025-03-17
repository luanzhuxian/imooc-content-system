package services

import (
	"imooc-content-system/internal/dao"

	"net/http"

	"github.com/gin-gonic/gin"
)

type ContentDeleteReq struct {
	ID int `json:"id" binding:"required"` // 内容ID
}

type ContentDeleteRsp struct {
	Message string `json:"message"`
}

func (c *CmsApp) ContentDelete(ctx *gin.Context) {
	var req ContentDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contentDao := dao.NewContentDao(c.db)
	ok, err := contentDao.IsExist(req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "content not exist"})
		return
	}
	if err := contentDao.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// if err != nil {
	// 	return
	// }
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &ContentDeleteRsp{
			Message: "ok",
		},
	})
}
