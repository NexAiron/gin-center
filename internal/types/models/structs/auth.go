package structs

import (
	security_types "gin-center/pkg/security/types"
)

// 重用security_types包中的结构体和接口
type BaseClaims = security_types.BaseClaims
type Claims = security_types.Claims
type AdminClaims = security_types.AdminClaims
type JWTUserClaims = security_types.JWTUserClaims
type UserClaims = security_types.UserClaims
type TokenPair = security_types.TokenPair
type AuthResponse = security_types.AuthResponse
type UserInfo = security_types.UserInfo

// 保留原有的请求结构体，这些在security_types中没有
type LoginRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=8,max=72"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	Fingerprint  string `json:"fingerprint,omitempty"`
}
