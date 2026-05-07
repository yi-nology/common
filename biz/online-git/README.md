# Online Git SDK

统一的在线 Git 平台操作 SDK，支持 GitHub、GitLab、Gitea 三大平台，提供一致的接口抽象。

> 注意：本包用于**在线 Git 平台 API 操作**（PR、分支管理、评论等）。如需本地 Git 仓库操作，请使用 `common/git` 包。

## 安装

```go
import (
    onlinegit "github.com/yi-nology/common/biz/online-git"

    // 按需导入平台实现（通过 init() 自动注册）
    _ "github.com/yi-nology/common/biz/online-git/github"
    _ "github.com/yi-nology/common/biz/online-git/gitlab"
    _ "github.com/yi-nology/common/biz/online-git/gitea"
)
```

## 快速开始

### GitHub

```go
package main

import (
    "context"
    "fmt"

    onlinegit "github.com/yi-nology/common/biz/online-git"
    _ "github.com/yi-nology/common/biz/online-git/github"
)

func main() {
    provider, err := onlinegit.NewGitProvider(&onlinegit.ProviderConfig{
        Platform: onlinegit.PlatformGitHub,
        Token:    "ghp_xxxxxxxxxxxx",
        Owner:    "your-org",
        Repo:     "your-repo",
    })
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // 获取仓库信息
    repo, err := provider.GetRepository(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Printf("仓库: %s (%s)\n", repo.FullName, repo.DefaultBranch)

    // 列出 PR
    prs, err := provider.ListPullRequests(ctx, onlinegit.PRStateOpen, &onlinegit.ListOptions{Page: 1, PerPage: 10})
    if err != nil {
        panic(err)
    }
    for _, pr := range prs {
        fmt.Printf("#%d %s [%s]\n", pr.Number, pr.Title, pr.State)
    }
}
```

### GitLab（私有化部署）

```go
provider, err := onlinegit.NewGitProvider(&onlinegit.ProviderConfig{
    Platform:        onlinegit.PlatformGitLab,
    BaseURL:         "https://gitlab.example.com",
    Token:           "glpat-xxxxxxxxxxxx",
    Owner:           "your-group",
    Repo:            "your-project",
    InsecureSkipTLS: true, // 跳过自签名证书验证
})
```

### Gitea

```go
provider, err := onlinegit.NewGitProvider(&onlinegit.ProviderConfig{
    Platform: onlinegit.PlatformGitea,
    BaseURL:  "https://gitea.example.com",
    Token:    "your-gitea-token",
    Owner:    "your-org",
    Repo:     "your-repo",
})
```

## 配置说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `Platform` | `Platform` | 是 | 平台类型：`github`、`gitlab`、`gitea` |
| `BaseURL` | `string` | 否* | API 地址。GitHub 默认 `https://api.github.com`，GitLab 默认 `https://gitlab.com`，Gitea 必填 |
| `Token` | `string` | 是 | 平台访问令牌 |
| `Owner` | `string` | 是 | 仓库所有者/组织/群组 |
| `Repo` | `string` | 是 | 仓库名称 |
| `InsecureSkipTLS` | `bool` | 否 | 跳过 TLS 证书验证，用于私有化部署自签名证书场景 |

## API 列表

### 仓库操作
| 方法 | 说明 |
|------|------|
| `GetRepository(ctx)` | 获取仓库信息 |

### 分支管理
| 方法 | 说明 |
|------|------|
| `ListBranches(ctx, opts)` | 列出所有分支 |
| `GetBranch(ctx, name)` | 获取分支详情 |
| `CreateBranch(ctx, name, sourceBranch)` | 从指定分支创建新分支 |
| `DeleteBranch(ctx, name)` | 删除分支 |
| `SetBranchProtection(ctx, name, rules)` | 设置分支保护规则 |
| `UnsetBranchProtection(ctx, name)` | 取消分支保护 |

### PR/MR 管理
| 方法 | 说明 |
|------|------|
| `ListPullRequests(ctx, state, opts)` | 列出合并请求（支持按状态过滤） |
| `GetPullRequest(ctx, number)` | 获取合并请求详情 |
| `CreatePullRequest(ctx, req)` | 创建合并请求 |
| `UpdatePullRequest(ctx, number, title, body)` | 更新合并请求 |
| `MergePullRequest(ctx, number, opts)` | 合并 PR（支持 merge/squash/rebase） |
| `ClosePullRequest(ctx, number)` | 关闭合并请求 |
| `GetPullRequestCommits(ctx, number)` | 获取 PR 包含的提交 |

### 分支比对
| 方法 | 说明 |
|------|------|
| `CompareBranches(ctx, base, head)` | 比较两个分支的差异 |

### 评论
| 方法 | 说明 |
|------|------|
| `ListComments(ctx, prNumber)` | 列出 PR 评论 |
| `CreateComment(ctx, prNumber, body)` | 添加评论 |
| `UpdateComment(ctx, commentID, body)` | 更新评论 |
| `DeleteComment(ctx, commentID)` | 删除评论 |

### 提交历史
| 方法 | 说明 |
|------|------|
| `GetCommit(ctx, sha)` | 获取提交详情 |
| `ListCommits(ctx, branch, opts)` | 列出分支提交历史 |

## 错误处理

SDK 提供了统一的错误类型和判断函数：

```go
repo, err := provider.GetRepository(ctx)
if err != nil {
    if onlinegit.IsNotFound(err) {
        fmt.Println("仓库不存在")
    } else if onlinegit.IsUnauthorized(err) {
        fmt.Println("Token 无效或已过期")
    } else if onlinegit.IsForbidden(err) {
        fmt.Println("权限不足")
    } else if onlinegit.IsRateLimit(err) {
        fmt.Println("请求频率超限")
    } else {
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

### 预定义错误
| 错误 | 说明 |
|------|------|
| `ErrNotFound` | 资源不存在 |
| `ErrUnauthorized` | 认证失败（Token 无效或过期） |
| `ErrForbidden` | 权限不足 |
| `ErrConflict` | 资源冲突（已存在） |
| `ErrRateLimit` | API 请求频率超限 |
| `ErrBadRequest` | 请求参数错误 |
| `ErrNotMergeable` | PR 无法合并 |
| `ErrBranchProtected` | 分支受保护 |
| `ErrInvalidPlatform` | 不支持的平台类型 |
| `ErrInvalidConfig` | 配置无效 |

## 工厂模式

SDK 使用工厂模式 + `init()` 自动注册，只需空导入对应平台包即可：

```go
// 只使用 GitHub
_ "github.com/yi-nology/common/biz/online-git/github"

// 使用全部平台
_ "github.com/yi-nology/common/biz/online-git/github"
_ "github.com/yi-nology/common/biz/online-git/gitlab"
_ "github.com/yi-nology/common/biz/online-git/gitea"
```

也可以动态查询支持的平台：

```go
platforms := onlinegit.GetSupportedPlatforms() // [github, gitlab, gitea]
valid := onlinegit.ValidatePlatform("github")  // true
url := onlinegit.GetDefaultBaseURL(onlinegit.PlatformGitHub) // https://api.github.com
```

## 文件结构

```
online-git/
├── models.go        # 数据模型定义（Repository, Branch, PR, Commit 等）
├── provider.go      # GitProvider 统一接口定义
├── errors.go        # 错误类型和判断函数
├── factory.go       # 工厂模式，Provider 注册与创建
├── github/
│   └── provider.go  # GitHub 平台实现
├── gitlab/
│   └── provider.go  # GitLab 平台实现
└── gitea/
    └── provider.go  # Gitea 平台实现
```

## 依赖

| 依赖 | 说明 |
|------|------|
| `github.com/google/go-github/v56` | GitHub REST API 客户端 |
| `gitlab.com/gitlab-org/api/client-go` | GitLab REST API 客户端 |
| `code.gitea.io/sdk/gitea` | Gitea REST API 客户端 |
| `golang.org/x/oauth2` | OAuth2 认证支持 |
