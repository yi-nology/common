package gitlab

import (
	"net/http"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// wrapError 包装 GitLab SDK 错误
func (p *Provider) wrapError(op string, resp *gitlab.Response, err error) error {
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrNotFound, "")
		case http.StatusUnauthorized:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrUnauthorized, "")
		case http.StatusForbidden:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrForbidden, "")
		case http.StatusConflict:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrConflict, err.Error())
		case http.StatusTooManyRequests:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrRateLimit, "")
		}
	}
	return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, err, "")
}
