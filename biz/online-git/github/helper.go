package github

import (
	"net/http"

	"github.com/google/go-github/v56/github"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// wrapError 包装 GitHub SDK 错误
func (p *Provider) wrapError(op string, resp *github.Response, err error) error {
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return onlinegit.NewProviderError(onlinegit.PlatformGitHub, op, onlinegit.ErrNotFound, "")
		case http.StatusUnauthorized:
			return onlinegit.NewProviderError(onlinegit.PlatformGitHub, op, onlinegit.ErrUnauthorized, "")
		case http.StatusForbidden:
			if resp.Rate.Remaining == 0 {
				return onlinegit.NewProviderError(onlinegit.PlatformGitHub, op, onlinegit.ErrRateLimit, "")
			}
			return onlinegit.NewProviderError(onlinegit.PlatformGitHub, op, onlinegit.ErrForbidden, "")
		case http.StatusConflict, http.StatusUnprocessableEntity:
			return onlinegit.NewProviderError(onlinegit.PlatformGitHub, op, onlinegit.ErrConflict, err.Error())
		}
	}
	return onlinegit.NewProviderError(onlinegit.PlatformGitHub, op, err, "")
}
