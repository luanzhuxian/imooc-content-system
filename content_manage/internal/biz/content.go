package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
)

// Content is a Content model.
type Content struct {
	ID             int64         `json:"id"`              // 内容标题
	ContentID      string        `json:"content_id"`      // 内容ID
	Title          string        `json:"title"`           // 内容标题
	VideoURL       string        `json:"video_url"`       // 视频播放URL
	Author         string        `json:"author"`          // 作者
	Description    string        `json:"description"`     // 内容描述
	Thumbnail      string        `json:"thumbnail"`       // 封面图URL
	Category       string        `json:"category"`        // 内容分类
	Duration       time.Duration `json:"duration"`        // 内容时长
	Resolution     string        `json:"resolution"`      // 分辨率 如720p、1080p
	FileSize       int64         `json:"fileSize"`        // 文件大小
	Format         string        `json:"format"`          // 文件格式 如MP4、AVI
	Quality        int32         `json:"quality"`         // 视频质量 1-高清 2-标清
	ApprovalStatus int32         `json:"approval_status"` // 审核状态 1-审核中 2-审核通过 3-审核不通过
	UpdatedAt      time.Time     `json:"updated_at"`      // 内容更新时间
	CreatedAt      time.Time     `json:"created_at"`      // 内容创建时间
}

type ContextIndex struct {
	ID        int64  `json:"id"`         // 自增ID
	ContentID string `json:"content_id"` // 内容ID
}

type FindParams struct {
	ID       int64
	Author   string
	Title    string
	Page     int32
	PageSize int32
}

// ContentRepo is a Content repo.
type ContentRepo interface {
	// Create 内容创建
	Create(ctx context.Context, c *Content) (int64, error)
	// Update 内容更新
	Update(ctx context.Context, id int64, c *Content) error
	// IsExist 内容是否存在
	IsExist(ctx context.Context, contentID int64) (bool, error)
	// Delete 删除内容
	Delete(ctx context.Context, id int64) error
	// Find 查找内容
	Find(ctx context.Context, params *FindParams) ([]*Content, int64, error)
	// FindIndex 查找内容关联的索引id
	FindIndex(ctx context.Context, params *FindParams) ([]*ContextIndex, int64, error)
	// First 查询指定ID内容详情
	First(ctx context.Context, idx *ContextIndex) (*Content, error)
}

// ContentUsecase is a Content usecase.
type ContentUsecase struct {
	repo ContentRepo
	log  *log.Helper
}

// NewContentUsecase new a Content usecase.
func NewContentUsecase(repo ContentRepo, logger log.Logger) *ContentUsecase {
	return &ContentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateContent creates a Content
func (uc *ContentUsecase) CreateContent(ctx context.Context, c *Content) (int64, error) {
	uc.log.WithContext(ctx).Infof("CreateContent: %+v", c)
	return uc.repo.Create(ctx, c)
}

// UpdateContent update a Content.
func (uc *ContentUsecase) UpdateContent(ctx context.Context, c *Content) error {
	uc.log.WithContext(ctx).Infof("UpdateContent: %+v", c)
	return uc.repo.Update(ctx, c.ID, c)
}

// DeleteContent delete a Content.
func (uc *ContentUsecase) DeleteContent(ctx context.Context, id int64) error {
	repo := uc.repo
	ok, err := repo.IsExist(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("内容不存在")
	}
	return repo.Delete(ctx, id)
}

// FindContent find Content.
func (uc *ContentUsecase) FindContent(ctx context.Context, params *FindParams) ([]*Content, int64, error) {
	repo := uc.repo
	indices, total, err := repo.FindIndex(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	var eg errgroup.Group
	contents := make([]*Content, len(indices), len(indices))
	for index, idx := range indices {
		tempIndex := index
		tempIdx := idx
		eg.Go(func() error {
			content, err := repo.First(ctx, tempIdx)
			if err != nil {
				return err
			}
			contents[tempIndex] = &Content{
				ID:             tempIdx.ID,
				ContentID:      content.ContentID,
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
				UpdatedAt:      content.UpdatedAt,
				CreatedAt:      content.CreatedAt,
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}
	return contents, total, err
}
