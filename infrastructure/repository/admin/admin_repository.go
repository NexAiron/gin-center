package admin

import (
	"context"
	"errors"
	base_repository "gin-center/infrastructure/repository/base_repository"
	AdminModel "gin-center/internal/domain/model/admin"

	"gorm.io/gorm"
)

type AdminRepository struct {
	*base_repository.GenericRepository[AdminModel.Admin]
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		GenericRepository: base_repository.NewGenericRepository[AdminModel.Admin](db),
	}
}
func (r *AdminRepository) Create(ctx context.Context, admin *AdminModel.Admin) error {
	return r.GenericRepository.Create(ctx, admin)
}
func (r *AdminRepository) FindByUsername(ctx context.Context, username string) (*AdminModel.Admin, error) {
	var admin AdminModel.Admin
	result := r.GenericRepository.DB.WithContext(ctx).Where("username = ?", username).First(&admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &admin, nil
}
func (r *AdminRepository) Update(ctx context.Context, admin *AdminModel.Admin) error {
	return r.GenericRepository.Update(ctx, admin)
}
func (r *AdminRepository) PaginateAdmins(ctx context.Context, page, pageSize int) ([]AdminModel.Admin, int64, error) {
	var admins []AdminModel.Admin
	var total int64
	if err := r.GenericRepository.DB.WithContext(ctx).Model(&AdminModel.Admin{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	result := r.GenericRepository.DB.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&admins)
	return admins, total, result.Error
}
func (r *AdminRepository) Delete(ctx context.Context, username string) error {
	return r.GenericRepository.DB.WithContext(ctx).Where("username = ?", username).Delete(&AdminModel.Admin{}).Error
}
