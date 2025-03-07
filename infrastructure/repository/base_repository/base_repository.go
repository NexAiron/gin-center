package base_repository

import (
	"context"
	"fmt"
	zaplogger "gin-center/infrastructure/zaplogger"
	"gin-center/internal/types/models/base"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type contextKey string

const (
	dbContextKey contextKey = "db"
	txContextKey contextKey = "tx"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*T, error)
	FindAll(ctx context.Context) ([]T, error)
	FindWithPagination(ctx context.Context, page, size int) ([]T, int64, error)
}

type GenericRepository[T base.Model] struct {
	DB *gorm.DB
}

func NewGenericRepository[T base.Model](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{
		DB: db,
	}
}

func (r *GenericRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *GenericRepository[T]) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(new(T), id).Error
}

func (r *GenericRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GenericRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GenericRepository[T]) FindWithPagination(ctx context.Context, page, size int) ([]T, int64, error) {
	var entities []T
	var total int64
	query := r.DB.WithContext(ctx)

	err := query.Model(new(T)).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * size).Limit(size).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func WithTx[T any](ctx context.Context, fn func(ctx context.Context) (T, error)) (T, error) {
	db := ctx.Value(dbContextKey).(*gorm.DB)
	tx := db.Begin()
	if tx.Error != nil {
		var zero T
		return zero, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	logger := zaplogger.NewServiceLogger()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.LogError("panic in transaction", zap.Skip(), zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	ctxWithTx := context.WithValue(ctx, txContextKey, tx)
	result, err := fn(ctxWithTx)

	if err != nil {
		logger.LogError("transaction failed", zap.Skip(), zap.Error(err))
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			logger.LogError("failed to rollback transaction", zap.Skip(), zap.Error(rollbackErr))
			return result, fmt.Errorf("transaction failed: %w (rollback failed: %v)", err, rollbackErr)
		}
		return result, err
	}

	if err := tx.Commit().Error; err != nil {
		logger.LogError("failed to commit transaction", zap.Skip(), zap.Error(err))
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			logger.LogError("failed to rollback transaction after commit error", zap.Skip(), zap.Error(rollbackErr))
			return result, fmt.Errorf("failed to commit transaction: %w (rollback failed: %v)", err, rollbackErr)
		}
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}
