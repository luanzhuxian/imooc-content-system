package services

import (
	"imooc-content-system/internal/dao"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ContentFindReq struct {
	ID       int    `json:"id"`        // 内容ID
	Author   string `json:"author"`    // 作者
	Title    string `json:"title"`     // 标题
	Page     int    `json:"page"`      // 页
	PageSize int    `json:"page_size"` // 页大小
}

type Content struct {
	ID             int           `json:"id"`                        // 内容ID
	Title          string        `json:"title"`                     // 内容标题
	VideoURL       string        `json:"video_url" `                // 视频播放URL
	Author         string        `json:"author" binding:"required"` // 作者
	Description    string        `json:"description"`               // 内容描述
	Thumbnail      string        `json:"thumbnail"`                 // 封面图URL
	Category       string        `json:"category"`                  // 内容分类
	Duration       time.Duration `json:"duration"`                  // 内容时长
	Resolution     string        `json:"resolution"`                // 分辨率 如720p、1080p
	FileSize       int64         `json:"fileSize"`                  // 文件大小
	Format         string        `json:"format"`                    // 文件格式 如MP4、AVI
	Quality        int           `json:"quality"`                   // 视频质量 1-高清 2-标清
	ApprovalStatus int           `json:"approval_status"`
}

type ContentFindRsp struct {
	Message  string    `json:"message"`
	Contents []Content `json:"contents"`
	Total    int64     `json:"total"`
}

func (c *CmsApp) ContentFind(ctx *gin.Context) {
	var req ContentFindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentDao := dao.NewContentDao(c.db)
	contentList, total, err := contentDao.Find(&dao.FindParams{
		ID:       req.ID,
		Author:   req.Author,
		Title:    req.Title,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// ???
	// 这段代码创建新切片而不直接返回contentList主要有以下几个重要原因：
	// 1. 数据模型转换（Model-DTO转换）
	// model.ContentDetail是数据库模型 - 在DAO层使用
	// services.Content是DTO (Data Transfer Object) - 用于API响应
	// 这种分层设计遵循了"关注点分离"的原则
	contents := make([]Content, 0, len(contentList))
	for _, content := range contentList {
		contents = append(contents, Content{
			ID:             content.ID,
			Title:          content.Title,
			VideoURL:       content.VideoURL,
			Author:         content.Author,
			Description:    content.Description,
			Thumbnail:      content.Thumbnail,
			Category:       content.Category,
			Duration:       content.Duration,
			Resolution:     content.Resolution,
			FileSize:       content.FileSize,
			Format:         content.Format,
			Quality:        content.Quality,
			ApprovalStatus: content.ApprovalStatus,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &ContentFindRsp{
			Message:  "ok",
			Contents: contents,
			Total:    total,
		},
	})
}
