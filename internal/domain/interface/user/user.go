package use_userInterface

import (
	"context"
	UserModel "gin-center/internal/domain/model/user"
	"gin-center/internal/types/models/structs"
	type_response "gin-center/internal/types/response"
)

type UserServiceInterface interface {
	Register(username, password string, extraFields ...interface{}) error
	// Login 用户登录
	Login(username, password string) (map[string]interface{}, error)
	ValidateToken(tokenString string) (*structs.UserClaims, error)
	GetUserByID(ctx context.Context, id uint) (*UserModel.User, error)
	UpdateUser(ctx context.Context, user *UserModel.User) error
	ListUsers(ctx context.Context, page, pageSize int, query map[string]interface{}) (*type_response.UserListResponse, error)
	UpdateUserAvatar(ctx context.Context, userID uint, avatarPath string) error
	UpdateUserProfile(ctx context.Context, userID uint, profile *type_response.UpdateUserProfileRequest) error
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
}
