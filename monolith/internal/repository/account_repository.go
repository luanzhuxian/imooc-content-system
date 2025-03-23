package repository

import (
	"context"
	"imooc-content-system/internal/domain"
)

// 定义AccountRepository接口
type AccountRepository interface {
	FindByUserID(ctx context.Context, userID string) (*domain.Account, error)
	Save(ctx context.Context, account *domain.Account) error
	Exists(ctx context.Context, userID string) (bool, error)
}