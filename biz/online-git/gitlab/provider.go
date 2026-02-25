package gitlab

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

func init() {
	onlinegit.RegisterProvider(onlinegit.PlatformGitLab, func(cfg *onlinegit.ProviderConfig) (onlinegit.GitProvider, error) {
		return NewProvider(cfg)
	})
}

// Provider GitLab 平台实现
type Provider struct {
	client    *gitlab.Client
	owner     string
	repo      string
	projectID string // owner/repo 格式
}

// NewProvider 创建 GitLab Provider
func NewProvider(cfg *onlinegit.ProviderConfig) (*Provider, error) {
	var client *gitlab.Client
	var err error

	httpClient := &http.Client{}
	if cfg.InsecureSkipTLS {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	opts := []gitlab.ClientOptionFunc{
		gitlab.WithHTTPClient(httpClient),
	}

	// 支持私有化部署
	if cfg.BaseURL != "" && cfg.BaseURL != "https://gitlab.com" {
		parsedURL, parseErr := url.Parse(cfg.BaseURL)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid base URL: %w", parseErr)
		}
		opts = append(opts, gitlab.WithBaseURL(parsedURL.String()))
	}

	client, err = gitlab.NewClient(cfg.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}

	return &Provider{
		client:    client,
		owner:     cfg.Owner,
		repo:      cfg.Repo,
		projectID: cfg.Owner + "/" + cfg.Repo,
	}, nil
}

func (p *Provider) GetPlatform() onlinegit.Platform {
	return onlinegit.PlatformGitLab
}
