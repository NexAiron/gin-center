package security_types

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenPair 定义了访问令牌和刷新令牌的结构
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Claims 定义了JWT令牌的基本接口
type Claims interface {
	jwt.Claims
	GenerateToken(secretKey string) (string, error)
}

// JWTUserClaims 普通用户JWT声明结构
type JWTUserClaims struct {
	BaseClaims
	UserID      string `json:"user_id"`
	TokenType   string `json:"token_type"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

// Valid 实现jwt.Claims接口
func (c *JWTUserClaims) Valid() error {
	if err := c.ExpiresAt.Valid(); err != nil {
		return err
	}
	return nil
}

// GenerateToken 为JWTUserClaims实现Claims接口
func (c *JWTUserClaims) GenerateToken(secretKey string) (string, error) {
	return c.BaseClaims.GenerateToken(secretKey)
}

// UserClaims 普通用户声明结构
type UserClaims struct {
	BaseClaims
	UserID string `json:"user_id"`
}

// Valid 实现jwt.Claims接口
func (c *UserClaims) Valid() error {
	return c.BaseClaims.Valid()
}

// GenerateToken 为UserClaims实现Claims接口
func (c *UserClaims) GenerateToken(secretKey string) (string, error) {
	return c.BaseClaims.GenerateToken(secretKey)
}

// AdminClaims 管理员专用的Claims结构
type AdminClaims struct {
	BaseClaims
	AdminID string `json:"admin_id"`
	IsAdmin bool   `json:"is_admin"`
}

// Valid 实现jwt.Claims接口
func (c *AdminClaims) Valid() error {
	return c.BaseClaims.Valid()
}

// GenerateToken 为AdminClaims实现Claims接口
func (c *AdminClaims) GenerateToken(secretKey string) (string, error) {
	return c.BaseClaims.GenerateToken(secretKey)
}

// BaseClaims 基础声明结构
type BaseClaims struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	ExpiresAt *TimeStamp `json:"exp,omitempty"`
}

// GenerateToken 生成JWT令牌
func (c *BaseClaims) GenerateToken(secretKey string) (string, error) {
	if secretKey == "" {
		return "", fmt.Errorf("JWT密钥无效")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("生成JWT令牌失败: %w", err)
	}

	return tokenString, nil
}

// TimeStamp 用于JWT的时间戳
type TimeStamp struct {
	Time time.Time
}

// Valid 验证时间戳是否有效
func (t *TimeStamp) Valid() error {
	if t == nil {
		return nil
	}
	now := time.Now()
	if t.Time.Before(now) {
		return errors.New("token is expired")
	}
	return nil
}

// Valid 实现jwt.Claims接口，验证令牌的有效性
func (c *BaseClaims) Valid() error {
	if c.ExpiresAt == nil {
		return nil
	}
	return c.ExpiresAt.Valid()
}

// AuthResponse 认证响应结构
type AuthResponse struct {
	Tokens TokenPair `json:"tokens"`
	User   UserInfo  `json:"user"`
}

// UserInfo 用户信息接口
type UserInfo interface {
	GetID() string
	GetUsername() string
	GetPhone() string
}