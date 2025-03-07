package use_request

import (
	"fmt"
	"gin-center/internal/types/models/structs"
	use_headers "gin-center/pkg/http/headers"
	use_response "gin-center/pkg/http/response"
	useJwt "gin-center/pkg/security/useJwt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func init() {
	// 初始化验证器
	validate.RegisterValidation("mobile", validateMobile)
	// 注册日期验证器
	validate.RegisterValidation("date", validateDate)
}

func BindAndValidate(c *gin.Context, req interface{}) error {
	if err := c.ShouldBind(req); err != nil {
		use_response.BadRequest(c, "")
		return err
	}
	setDefaultValues(req)
	err := validate.Struct(req)
	if err == nil {
		return nil
	}
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		use_response.BadRequest(c, "")
		return err
	}
	use_response.BadRequest(c, formatValidationError(validationErrors))
	return err
}
func formatValidationError(errs validator.ValidationErrors) string {
	var errMsgs []string
	for _, err := range errs {
		errMsg := getValidationErrorMessage(err)
		errMsgs = append(errMsgs, errMsg)
	}
	return strings.Join(errMsgs, "; ")
}

var validationErrorTemplates = map[string]string{
	"required": "%s?",
	"min":      "%s%s",
	"max":      "%s%s",
	"len":      "%s%s",
	"mobile":   "%s",
	"date":     "%s",
	"regexp":   "%s",
}

func getValidationErrorMessage(err validator.FieldError) string {
	if template, ok := validationErrorTemplates[err.Tag()]; ok {
		if err.Param() != "" {
			return fmt.Sprintf(template, err.Field(), err.Param())
		}
		return fmt.Sprintf(template, err.Field())
	}
	return fmt.Sprintf("%s", err.Field(), err.Tag())
}
func GetParam(c *gin.Context, key string, defaultValue string) string {
	value := c.Query(key)
	if value != "" {
		return value
	}
	value = c.Param(key)
	if value != "" {
		return value
	}
	return defaultValue
}
func GetParamInt(c *gin.Context, key string, defaultValue int) int {
	value := GetParam(c, key, "")
	if value == "" {
		return defaultValue
	}
	result, err := parseInt(value)
	if err != nil {
		return defaultValue
	}
	return result
}
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	return matched
}
func validateDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
func setDefaultValues(obj interface{}) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}
		defaultTag := typ.Field(i).Tag.Get("default")
		if defaultTag == "" {
			continue
		}
		setFieldDefaultValue(field, defaultTag)
	}
}
func setFieldDefaultValue(field reflect.Value, defaultTag string) {
	switch field.Kind() {
	case reflect.String:
		if field.String() == "" {
			field.SetString(defaultTag)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() == 0 {
			var intValue int
			if _, err := fmt.Sscanf(defaultTag, "%d", &intValue); err == nil {
				field.SetInt(int64(intValue))
			}
		}
	}
}
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
func GetTokenFromHeader(c *gin.Context) string {
	return use_headers.GetAuthorizationToken(c)
}
func ValidateToken(c *gin.Context, jwtConfig *useJwt.JWTConfig) (*structs.UserClaims, error) {
	token := GetTokenFromHeader(c)
	if token == "" {
		use_response.Unauthorized(c, "Token")
		return nil, useJwt.ErrInvalidToken
	}
	claims, err := jwtConfig.ParseToken(token)
	if err != nil {
		use_response.Unauthorized(c, "Token")
		return nil, err
	}
	return claims, nil
}
