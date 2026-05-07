package github

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

func init() {
	onlinegit.RegisterProvider(onlinegit.PlatformGitHub, func(cfg *onlinegit.ProviderConfig) (onlinegit.GitProvider, error) {
		return NewProvider(cfg)
	})
}

// Provider GitHub 平台实现
type Provider struct {
	client *github.Client
	owner  string
	repo   string
}

// NewProvider 创建 GitHub Provider
func NewProvider(cfg *onlinegit.ProviderConfig) (*Provider, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)

	transport := &http.Transport{}
	if cfg.InsecureSkipTLS {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	tc := &http.Client{
		Transport: &oauth2.Transport{
			Source: ts,
			Base:   transport,
		},
	}

	client := github.NewClient(tc)

	// 支持 GitHub Enterprise
	if cfg.BaseURL != "" && cfg.BaseURL != "https://api.github.com" {
		var err error
		parsedURL, err := url.Parse(cfg.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid base URL: %w", err)
		}
		// 确保 URL 以 / 结尾
		if !strings.HasSuffix(parsedURL.Path, "/") {
			parsedURL.Path += "/"
		}
		client.BaseURL = parsedURL
	}

	return &Provider{
		client: client,
		owner:  cfg.Owner,
		repo:   cfg.Repo,
	}, nil
}

func (p *Provider) GetPlatform() onlinegit.Platform {
	return onlinegit.PlatformGitHub
}
