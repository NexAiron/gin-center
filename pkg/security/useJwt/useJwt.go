package useJwt

import (
	"errors"
	"fmt"
	security_errors "gin-center/pkg/security/errors"
	security_types "gin-center/pkg/security/types"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidToken = security_errors.ErrInvalidToken
	ErrExpiredToken = security_errors.ErrExpiredToken
	ErrRevokedToken = security_errors.ErrRevokedToken
	ErrEmptyToken   = security_errors.ErrEmptyToken
)

// JWTConfig JWT配置
type JWTConfig struct {
	// SecretKey JWT签名密钥（自动转换为字节数组）
	SecretKey  string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
	// Issuer Token签发者
	Issuer string `mapstructure:"issuer"`
	// AccessTokenLifetime 访问令牌生命周期
	AccessTokenLifetime time.Duration `mapstructure:"access_token_lifetime"`
	// RefreshTokenLifetime 刷新令牌生命周期
	RefreshTokenLifetime time.Duration `mapstructure:"refresh_token_lifetime"`
	// BlacklistCleanupTick 黑名单清理间隔
	BlacklistCleanupTick time.Duration `mapstructure:"blacklist_cleanup_tick"`
	// 签名方法配置
	SigningMethod     jwt.SigningMethod
	blacklistedTokens sync.Map
}

type BlacklistedToken struct {
	Expiry time.Time
}

// NewJWTConfig 创建新的JWT配置实例
func NewJWTConfig(cfg *JWTConfig) *JWTConfig {
	if cfg == nil || cfg.SecretKey == "" || cfg.Issuer == "" || cfg.AccessTokenLifetime <= 0 || cfg.RefreshTokenLifetime <= 0 || cfg.BlacklistCleanupTick <= 0 {
		return nil
	}
	jwtConfig := &JWTConfig{
		SecretKey:            cfg.SecretKey,
		Issuer:               cfg.Issuer,
		AccessTokenLifetime:  cfg.AccessTokenLifetime,
		RefreshTokenLifetime: cfg.RefreshTokenLifetime,
		BlacklistCleanupTick: cfg.BlacklistCleanupTick,
		SigningMethod:        jwt.SigningMethodHS256,
	}
	go jwtConfig.cleanupBlacklist()
	return jwtConfig
}

// GenerateTokenWithClaims 生成包含自定义声明的JWT令牌
func (c *JWTConfig) GenerateTokenWithClaims(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.SecretKey))
	if err != nil {
		return "", fmt.Errorf("生成令牌失败: %w", err)
	}
	return tokenString, nil
}

func (c *JWTConfig) GenerateTokenPair(userID string, username, role string, fingerprint string) (*security_types.TokenPair, error) {
	accessClaims := &security_types.UserClaims{
		BaseClaims: security_types.BaseClaims{
			Username: username,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(c.SecretKey))
	if err != nil {
		return nil, err
	}
	refreshClaims := &security_types.UserClaims{
		BaseClaims: security_types.BaseClaims{
			Username: username,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(c.SecretKey))
	if err != nil {
		return nil, err
	}
	return &security_types.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
func (c *JWTConfig) ParseToken(tokenString string) (*security_types.UserClaims, error) {
	if tokenString == "" {
		return nil, ErrEmptyToken
	}
	if _, exists := c.blacklistedTokens.Load(tokenString); exists {
		return nil, ErrRevokedToken
	}
	token, err := jwt.ParseWithClaims(tokenString, &security_types.JWTUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	jwtClaims, ok := token.Claims.(*security_types.JWTUserClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	
	// 将JWTUserClaims转换为UserClaims
	userClaims := &security_types.UserClaims{
		BaseClaims: jwtClaims.BaseClaims,
		UserID: jwtClaims.UserID,
	}
	return userClaims, nil
}

func (c *JWTConfig) RevokeToken(tokenString string) error {
	claims, err := c.ParseToken(tokenString)
	if err != nil && err != ErrRevokedToken {
		return err
	}
	c.blacklistedTokens.Store(tokenString, BlacklistedToken{
		Expiry: claims.ExpiresAt.Time,
	})
	return nil
}

func (c *JWTConfig) RefreshToken(refreshToken string, fingerprint string) (*security_types.TokenPair, error) {
	claims, err := c.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return c.GenerateTokenPair(claims.BaseClaims.ID, claims.Username, "", fingerprint)
}
func (c *JWTConfig) cleanupBlacklist() {
	ticker := time.NewTicker(c.BlacklistCleanupTick)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		c.blacklistedTokens.Range(func(key, value interface{}) bool {
			if token, ok := value.(BlacklistedToken); ok {
				if token.Expiry.Before(now) {
					c.blacklistedTokens.Delete(key)
				}
			}
			return true
		})
	}
}
