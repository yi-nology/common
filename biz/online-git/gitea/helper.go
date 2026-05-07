package gitea

import (
	"fmt"
	"net/http"

	"code.gitea.io/sdk/gitea"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// wrapError 包装 Gitea SDK 错误
func (p *Provider) wrapError(op string, resp *gitea.Response, err error) error {
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrNotFound, "")
		case http.StatusUnauthorized:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrUnauthorized, "")
		case http.StatusForbidden:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrForbidden, "")
		case http.StatusConflict:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrConflict, err.Error())
		case http.StatusTooManyRequests:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrRateLimit, "")
		}
	}
	return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, err, "")
}

// wrapHTTPError 处理直接 HTTP 调用的错误
func (p *Provider) wrapHTTPError(op string, statusCode int, body string) error {
	switch statusCode {
	case http.StatusNotFound:
		return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrNotFound, body)
	case http.StatusUnauthorized:
		return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrUnauthorized, body)
	case http.StatusForbidden:
		return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrForbidden, body)
	case http.StatusConflict:
		return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrConflict, body)
	case http.StatusTooManyRequests:
		return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrRateLimit, body)
	default:
		return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, fmt.Errorf("HTTP %d: %s", statusCode, body), "")
	}
}
