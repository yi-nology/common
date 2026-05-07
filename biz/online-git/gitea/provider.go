package gitea

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"code.gitea.io/sdk/gitea"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

func init() {
	onlinegit.RegisterProvider(onlinegit.PlatformGitea, func(cfg *onlinegit.ProviderConfig) (onlinegit.GitProvider, error) {
		return NewProvider(cfg)
	})
}

// Provider Gitea 平台实现
type Provider struct {
	client     *gitea.Client
	httpClient *http.Client
	baseURL    string
	token      string
	owner      string
	repo       string
}

// NewProvider 创建 Gitea Provider
func NewProvider(cfg *onlinegit.ProviderConfig) (*Provider, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required for Gitea")
	}

	opts := []gitea.ClientOption{
		gitea.SetToken(cfg.Token),
	}

	httpClient := &http.Client{}
	if cfg.InsecureSkipTLS {
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		opts = append(opts, gitea.SetHTTPClient(httpClient))
	}

	client, err := gitea.NewClient(cfg.BaseURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gitea client: %w", err)
	}

	// 确保 baseURL 没有尾随斜杠
	baseURL := strings.TrimSuffix(cfg.BaseURL, "/")

	return &Provider{
		client:     client,
		httpClient: httpClient,
		baseURL:    baseURL,
		token:      cfg.Token,
		owner:      cfg.Owner,
		repo:       cfg.Repo,
	}, nil
}

func (p *Provider) GetPlatform() onlinegit.Platform {
	return onlinegit.PlatformGitea
}
