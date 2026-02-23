package onlinegit

import (
	"fmt"
	"sync"
)

// ProviderFactory Provider 工厂函数类型
type ProviderFactory func(cfg *ProviderConfig) (GitProvider, error)

var (
	providersMu sync.RWMutex
	providers   = make(map[Platform]ProviderFactory)
)

// RegisterProvider 注册 Provider 工厂函数
// 各平台实现包在 init() 中调用此函数注册自己
func RegisterProvider(platform Platform, factory ProviderFactory) {
	providersMu.Lock()
	defer providersMu.Unlock()
	providers[platform] = factory
}

// NewGitProvider 根据配置创建对应平台的 Provider
func NewGitProvider(cfg *ProviderConfig) (GitProvider, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("%w: token is required", ErrInvalidConfig)
	}

	if cfg.Owner == "" || cfg.Repo == "" {
		return nil, fmt.Errorf("%w: owner and repo are required", ErrInvalidConfig)
	}

	providersMu.RLock()
	factory, ok := providers[cfg.Platform]
	providersMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPlatform, cfg.Platform)
	}

	return factory(cfg)
}

// ValidatePlatform 验证平台类型是否有效
func ValidatePlatform(platform string) bool {
	switch Platform(platform) {
	case PlatformGitHub, PlatformGitLab, PlatformGitea:
		return true
	default:
		return false
	}
}

// GetSupportedPlatforms 返回支持的平台列表
func GetSupportedPlatforms() []Platform {
	return []Platform{
		PlatformGitHub,
		PlatformGitLab,
		PlatformGitea,
	}
}

// GetDefaultBaseURL 返回平台的默认 API 地址
func GetDefaultBaseURL(platform Platform) string {
	switch platform {
	case PlatformGitHub:
		return "https://api.github.com"
	case PlatformGitLab:
		return "https://gitlab.com"
	case PlatformGitea:
		return "" // Gitea 没有公共实例，必须指定 BaseURL
	default:
		return ""
	}
}
