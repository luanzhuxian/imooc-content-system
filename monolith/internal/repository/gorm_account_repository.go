package repository

import (
	"context"
	"imooc-content-system/internal/domain"
	"imooc-content-system/internal/model"

	"gorm.io/gorm"
)

type AccountGormRepository struct {
	db *gorm.DB
}

// 创建实例 依赖注入
// 实现AccountRepository接口
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &AccountGormRepository{db: db}
}

func (r *AccountGormRepository) FindByUserID(ctx context.Context, userID string) (*domain.Account, error) {
	var accountModel model.Account
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&accountModel).Error
	if err != nil {
		return nil, err
	}
	
	// Convert from DB model to domain model
	return &domain.Account{
		ID:       accountModel.ID,
		UserID:   accountModel.UserID,
		Password: accountModel.Password,
		Nickname: accountModel.Nickname,
		CreatedAt: accountModel.Ct,
		UpdatedAt: accountModel.Ut,
	}, nil
}

func (r *AccountGormRepository) Save(ctx context.Context, account *domain.Account) error {
	// Convert from domain model to DB model
	accountModel := model.Account{
		UserID:   account.UserID,
		Password: account.Password,
		Nickname: account.Nickname,
		Ct:       account.CreatedAt,
		Ut:       account.UpdatedAt,
	}
	
	return r.db.WithContext(ctx).Create(&accountModel).Error
}

func (r *AccountGormRepository) Exists(ctx context.Context, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Account{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}