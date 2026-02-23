package onlinegit

import (
	"errors"
	"fmt"
)

// 预定义错误
var (
	ErrNotFound        = errors.New("resource not found")
	ErrUnauthorized    = errors.New("unauthorized: invalid or expired token")
	ErrForbidden       = errors.New("forbidden: insufficient permissions")
	ErrConflict        = errors.New("conflict: resource already exists")
	ErrRateLimit       = errors.New("rate limit exceeded")
	ErrBadRequest      = errors.New("bad request: invalid parameters")
	ErrNotMergeable    = errors.New("pull request is not mergeable")
	ErrBranchProtected = errors.New("branch is protected")
	ErrInvalidPlatform = errors.New("invalid or unsupported platform")
	ErrInvalidConfig   = errors.New("invalid configuration")
)

// ProviderError 平台特定错误
type ProviderError struct {
	Platform Platform
	Op       string
	Err      error
	Message  string
}

func (e *ProviderError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("[%s] %s: %s - %v", e.Platform, e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s: %v", e.Platform, e.Op, e.Err)
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError 创建平台错误
func NewProviderError(platform Platform, op string, err error, message string) *ProviderError {
	return &ProviderError{
		Platform: platform,
		Op:       op,
		Err:      err,
		Message:  message,
	}
}

// IsNotFound 检查是否为资源不存在错误
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsUnauthorized 检查是否为认证错误
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden 检查是否为权限错误
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsConflict 检查是否为冲突错误
func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

// IsRateLimit 检查是否为限流错误
func IsRateLimit(err error) bool {
	return errors.Is(err, ErrRateLimit)
}
