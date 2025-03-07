package user_repo

import (
	"context"
	"errors"
	"fmt"
	infraErrors "gin-center/infrastructure/errors"
	UserModel "gin-center/internal/domain/model/user"
	"gin-center/internal/types/constants"

	base_repository "gin-center/infrastructure/repository/base_repository"

	"gorm.io/gorm"
)

type UserRepository struct {
	*base_repository.GenericRepository[UserModel.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		GenericRepository: base_repository.NewGenericRepository[UserModel.User](db),
	}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*UserModel.User, error) {
	var user UserModel.User
	if err := r.GenericRepository.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Register(ctx context.Context, user *UserModel.User) error {
	_, err := base_repository.WithTx(ctx, func(txCtx context.Context) (*UserModel.User, error) {
		if exists, err := r.isUsernameExists(txCtx, user.Username); err != nil {
			return nil, fmt.Errorf("检查用户名是否存在失败: %w", err)
		} else if exists {
			return nil, infraErrors.ErrUsernameExists
		}

		if err := r.GenericRepository.DB.WithContext(txCtx).Create(user).Error; err != nil {
			return nil, fmt.Errorf("创建用户失败: %w", err)
		}
		return user, nil
	})
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *UserModel.User) (*UserModel.User, error) {
	return base_repository.WithTx(ctx, func(txCtx context.Context) (*UserModel.User, error) {
		existingUser, err := r.FindByID(txCtx, user.ID)
		if err != nil {
			return nil, err
		}

		if user.Username != existingUser.Username {
			if exists, err := r.isUsernameExists(txCtx, user.Username); err != nil {
				return nil, fmt.Errorf("检查用户名是否存在失败: %w", err)
			} else if exists {
				return nil, infraErrors.ErrUsernameExists
			}
		}

		if err := r.GenericRepository.Update(txCtx, user); err != nil {
			return nil, fmt.Errorf("更新用户失败: %w", err)
		}
		return user, nil
	})
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	if err := r.GenericRepository.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrUserNotFound
		}
		return fmt.Errorf("删除用户失败: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*UserModel.User, error) {
	user, err := r.GenericRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, page, pageSize int, query map[string]interface{}) ([]*UserModel.User, int64, error) {
	var users []*UserModel.User
	var total int64

	tx := r.GenericRepository.DB.WithContext(ctx).Model(&UserModel.User{})
	for key, value := range query {
		switch key {
		case "username":
			tx = tx.Where("username LIKE ?", fmt.Sprintf("%%%v%%", value))
		}
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计用户数量失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := tx.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	return users, total, nil
}

func (r *UserRepository) UpdateAvatar(ctx context.Context, userID uint, avatarPath string) error {
	result := r.GenericRepository.DB.WithContext(ctx).Model(&UserModel.User{}).Where("id = ?", userID).Update("avatar", avatarPath)
	if result.Error != nil {
		return fmt.Errorf("更新用户头像失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return constants.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) isUsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.GenericRepository.DB.WithContext(ctx).Model(&UserModel.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
