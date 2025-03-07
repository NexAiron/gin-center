package UserModel

import (
	baseModel "gin-center/internal/domain/model/base"
	"time"
)

type User struct {
	baseModel.BaseModel
	Username    string     `gorm:"uniqueIndex" json:"username"`
	Password    string     `json:"password"`
	Phone       string     `json:"phone"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `json:"last_login_ip"`
	UserType    string     `json:"user_type" validate:"required,oneof=admin regular guest"`
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
		BaseModel: baseModel.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
func (User) TableName() string {
	return "normal_users"
}
