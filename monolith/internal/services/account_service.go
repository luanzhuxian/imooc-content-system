package services

import (
	"context"
	"time"

	"imooc-content-system/internal/common"
	"imooc-content-system/internal/domain"
	"imooc-content-system/internal/repository"
)

// 创建实例 依赖注入
type AccountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepository: accountRepo,
	}
}

func (s *AccountService) RegisterUser(ctx context.Context, userID, password, nickname string) error {
	// 1. 业务规则验证 (例如密码强度检查等)
	// ...

	// 2. 加密密码
	hashedPassword, err := common.EncryptPassword(password)
	if err != nil {
		return err
	}

	// 3. 检查用户是否已存在
	exists, err := s.accountRepository.Exists(ctx, userID)
	if err != nil {
		return err
	}

	if exists {
		return common.ErrUserAlreadyExists
	}

	// 4. 创建账户
	currentTime := time.Now()
	account := &domain.Account{
		UserID:    userID,
		Password:  hashedPassword,
		Nickname:  nickname,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	// 5. 保存账户
	return s.accountRepository.Save(ctx, account)
}

func (s *AccountService) FindUser(ctx context.Context, userID string) (*domain.Account, error) {
	return s.accountRepository.FindByUserID(ctx, userID)
}
