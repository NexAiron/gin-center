package validator

import (
	infraErrors "gin-center/infrastructure/errors"
	"gin-center/internal/types/auth"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return auth.NewAuthError(infraErrors.ErrCodeInvalidUsername)
	}
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return auth.NewAuthError(infraErrors.ErrCodeUsernameValidationFailed)
		}
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 72 {
		return auth.NewAuthError(infraErrors.ErrCodeInvalidPassword)
	}
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return auth.NewAuthError(infraErrors.ErrCodeInvalidPassword)
	}
	return nil
}
