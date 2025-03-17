package data

import (
	"content_manage/internal/biz"
	"context"
	"fmt"
	"hash/fnv"
	"math/big"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

const contentNumTables = 4

type contentRepo struct {
	data *Data
	log  *log.Helper
}

// NewContentRepo .
func NewContentRepo(data *Data, logger log.Logger) biz.ContentRepo {
	return &contentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

type ContentDetail struct {
	ID             int64         `gorm:"column:id;primary_key"`  // 自增ID
	ContentID      string        `gorm:"column:content_id"`      // 内容ID
	Title          string        `gorm:"column:title"`           // 内容标题
	Description    string        `gorm:"column:description"`     // 内容描述
	Author         string        `gorm:"column:author"`          // 作者
	VideoURL       string        `gorm:"column:video_url"`       // 视频播放URL
	Thumbnail      string        `gorm:"column:thumbnail"`       // 封面图URL
	Category       string        `gorm:"column:category"`        // 内容分类
	Duration       time.Duration `gorm:"column:duration"`        // 内容时长
	Resolution     string        `gorm:"column:resolution"`      // 分辨率 如720p、1080p
	FileSize       int64         `gorm:"column:fileSize"`        // 文件大小
	Format         string        `gorm:"column:format"`          // 文件格式 如MP4、AVI
	Quality        int32         `gorm:"column:quality"`         // 视频质量 1-高清 2-标清
	ApprovalStatus int32         `gorm:"column:approval_status"` // 审核状态 1-审核中 2-审核通过 3-审核不通过
	UpdatedAt      time.Time     `gorm:"column:updated_at"`      // 内容更新时间
	CreatedAt      time.Time     `gorm:"column:created_at"`      // 内容创建时间
}

//func (c ContentDetail) TableName() string {
//	return "cms_content.t_content_details"
//}

type IdxContentDetail struct {
	ID        int64     `gorm:"column:id;primary_key"` // 自增ID
	ContentID string    `gorm:"column:content_id"`     // 内容ID
	Title     string    `gorm:"column:title"`          // 内容标题
	Author    string    `gorm:"column:author"`         // 作者
	UpdatedAt time.Time `gorm:"column:updated_at"`     // 内容更新时间
	CreatedAt time.Time `gorm:"column:created_at"`     // 内容创建时间
}

func (c IdxContentDetail) TableName() string {
	return "cms_content.t_idx_content_details"
}

func getContentDetailsTable(contentID string) string {
	tableIndex := getContentTableIndex(contentID)
	table := fmt.Sprintf("cms_content.t_content_details_%d", tableIndex)
	log.Infof("content_id = %s, table = %s", contentID, table)
	return table
}

func getContentTableIndex(uuid string) int {
	// 计算UUID的哈希值
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(uuid))
	hashValue := hash.Sum32()
	// 将哈希值映射到表的索引范围内
	// 创建一个大整数对象
	bigNum := big.NewInt(int64(hashValue))
	// 创建一个大整数对象
	bigModulo := big.NewInt(contentNumTables)
	// 使用大整数对象的Mod方法，计算哈希值对表数量取模的结果。这将确保结果在表的索引范围内。
	tableIndex := bigNum.Mod(bigNum, bigModulo).Int64()
	return int(tableIndex)
}

func (c *contentRepo) Create(ctx context.Context, content *biz.Content) (int64, error) {
	c.log.Infof("contentRepo Create content = %+v", content)
	db := c.data.db
	idx := IdxContentDetail{
		ContentID: content.ContentID,
		Title:     content.Title,
		Author:    content.Author,
	}
	if err := db.Create(&idx).Error; err != nil {
		return 0, err
	}
	detail := ContentDetail{
		Title:          content.Title,
		ContentID:      content.ContentID,
		Description:    content.Description,
		Author:         content.Author,
		VideoURL:       content.VideoURL,
		Thumbnail:      content.Thumbnail,
		Category:       content.Category,
		Duration:       content.Duration,
		Resolution:     content.Resolution,
		FileSize:       content.FileSize,
		Format:         content.Format,
		Quality:        content.Quality,
		ApprovalStatus: content.ApprovalStatus,
	}
	if err := db.Table(getContentDetailsTable(content.ContentID)).Create(&detail).Error; err != nil {
		c.log.Errorf("content create error = %v", err)
		return 0, err
	}
	return idx.ID, nil
}

func (c *contentRepo) Update(ctx context.Context, id int64, content *biz.Content) error {
	db := c.data.db
	var idx IdxContentDetail
	if err := db.Where("id = ?", id).First(&idx).Error; err != nil {
		return err
	}
	detail := ContentDetail{
		ContentID:      content.ContentID,
		Title:          content.Title,
		Description:    content.Description,
		Author:         content.Author,
		VideoURL:       content.VideoURL,
		Thumbnail:      content.Thumbnail,
		Category:       content.Category,
		Duration:       content.Duration,
		Resolution:     content.Resolution,
		FileSize:       content.FileSize,
		Format:         content.Format,
		Quality:        content.Quality,
		ApprovalStatus: content.ApprovalStatus,
	}
	if err := db.Table(getContentDetailsTable(idx.ContentID)).Where("content_id = ?", idx.ContentID).
		Updates(&detail).Error; err != nil {
		c.log.WithContext(ctx).Errorf("content update error = %v", err)
		return err
	}
	return nil
}

func (c *contentRepo) IsExist(ctx context.Context, id int64) (bool, error) {
	db := c.data.db
	var detail IdxContentDetail
	err := db.Where("id = ?", id).First(&detail).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		c.log.WithContext(ctx).Errorf("ContentDao isExist = [%v]", err)
		return false, err
	}
	return true, nil
}

func (c *contentRepo) Delete(ctx context.Context, id int64) error {
	db := c.data.db
	// 查询索引表信息
	var idx IdxContentDetail
	if err := db.Where("id = ?", id).First(&idx).Error; err != nil {
		return err
	}
	// 删除索引信息
	err := db.Where("id = ?", id).
		Delete(&IdxContentDetail{}).Error
	if err != nil {
		c.log.WithContext(ctx).Errorf("content delete error = %v", err)
		return err
	}
	// 删除详情信息
	err = db.Table(getContentDetailsTable(idx.ContentID)).
		Where("content_id = ?", idx.ContentID).
		Delete(&ContentDetail{}).Error
	if err != nil {
		c.log.WithContext(ctx).Errorf("content delete error = %v", err)
		return err
	}
	return nil
}

func (c *contentRepo) First(ctx context.Context, idx *biz.ContextIndex) (*biz.Content, error) {
	db := c.data.db
	var detail ContentDetail
	if err := db.Table(getContentDetailsTable(idx.ContentID)).
		Where("content_id = ?", idx.ContentID).First(&detail).Error; err != nil {
		c.log.WithContext(ctx).Errorf("content first error = %v", err)
		return nil, err
	}
	content := &biz.Content{
		ID:             detail.ID,
		ContentID:      detail.ContentID,
		Title:          detail.Title,
		VideoURL:       detail.VideoURL,
		Author:         detail.Author,
		Description:    detail.Description,
		Thumbnail:      detail.Thumbnail,
		Category:       detail.Category,
		Duration:       detail.Duration,
		Resolution:     detail.Resolution,
		FileSize:       detail.FileSize,
		Format:         detail.Format,
		Quality:        detail.Quality,
		ApprovalStatus: detail.ApprovalStatus,
		UpdatedAt:      detail.UpdatedAt,
		CreatedAt:      detail.CreatedAt,
	}
	return content, nil
}

// FindIndex 查询内容索引
func (c *contentRepo) FindIndex(ctx context.Context, params *biz.FindParams) ([]*biz.ContextIndex, int64, error) {
	// 构造查询条件
	query := c.data.db.Model(&IdxContentDetail{})
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.Author != "" {
		query = query.Where("author = ?", params.Author)
	}
	if params.Title != "" {
		query = query.Where("title = ?", params.Title)
	}
	// 总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var page, pageSize = 1, 10
	if params.Page > 0 {
		page = int(params.Page)
	}
	if params.PageSize > 0 {
		pageSize = int(params.PageSize)
	}
	offset := (page - 1) * pageSize
	var results []*IdxContentDetail
	if err := query.Offset(offset).Limit(pageSize).
		Find(&results).Error; err != nil {
		c.log.WithContext(ctx).Errorf("content find error = %v", err)
		return nil, 0, err
	}
	var contents []*biz.ContextIndex
	for _, r := range results {
		contents = append(contents, &biz.ContextIndex{
			ID:        r.ID,
			ContentID: r.ContentID,
		})
	}
	return contents, total, nil
}

func (c *contentRepo) Find(ctx context.Context, params *biz.FindParams) ([]*biz.Content, int64, error) {
	// 构造查询条件
	query := c.data.db.Model(&ContentDetail{})
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.Author != "" {
		query = query.Where("author = ?", params.Author)
	}
	if params.Title != "" {
		query = query.Where("title = ?", params.Title)
	}
	// 总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var page, pageSize = 1, 10
	if params.Page > 0 {
		page = int(params.Page)
	}
	if params.PageSize > 0 {
		pageSize = int(params.PageSize)
	}
	offset := (page - 1) * pageSize
	var results []*ContentDetail
	if err := query.Offset(offset).Limit(pageSize).
		Find(&results).Error; err != nil {
		c.log.WithContext(ctx).Errorf("content find error = %v", err)
		return nil, 0, err
	}
	var contents []*biz.Content
	for _, r := range results {
		contents = append(contents, &biz.Content{
			ID:             r.ID,
			Title:          r.Title,
			VideoURL:       r.VideoURL,
			Author:         r.Author,
			Description:    r.Description,
			Thumbnail:      r.Thumbnail,
			Category:       r.Category,
			Duration:       r.Duration,
			Resolution:     r.Resolution,
			FileSize:       r.FileSize,
			Format:         r.Format,
			Quality:        r.Quality,
			ApprovalStatus: r.ApprovalStatus,
			UpdatedAt:      r.UpdatedAt,
			CreatedAt:      r.CreatedAt,
		})
	}
	return contents, total, nil
}
