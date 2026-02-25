package onlinegit

import "context"

// GitProvider 定义统一的 Git 平台操作接口
// 所有平台实现（GitHub、GitLab、Gitea）都必须实现此接口
type GitProvider interface {
	// GetPlatform 返回平台类型
	GetPlatform() Platform

	// ==================== 仓库操作 ====================

	// GetRepository 获取仓库信息
	GetRepository(ctx context.Context) (*Repository, error)

	// ==================== 分支管理 ====================

	// ListBranches 列出所有分支
	ListBranches(ctx context.Context, opts *ListOptions) ([]*Branch, error)

	// GetBranch 获取分支详情
	GetBranch(ctx context.Context, name string) (*Branch, error)

	// CreateBranch 创建新分支
	// name: 新分支名称
	// sourceBranch: 源分支名称（从该分支创建）
	CreateBranch(ctx context.Context, name, sourceBranch string) (*Branch, error)

	// DeleteBranch 删除分支
	DeleteBranch(ctx context.Context, name string) error

	// SetBranchProtection 设置分支保护规则
	SetBranchProtection(ctx context.Context, name string, rules *ProtectionRules) error

	// UnsetBranchProtection 取消分支保护
	UnsetBranchProtection(ctx context.Context, name string) error

	// ==================== 合并请求/PR 管理 ====================

	// ListPullRequests 列出合并请求
	// state: open, closed, merged, all
	ListPullRequests(ctx context.Context, state PRState, opts *ListOptions) ([]*PullRequest, error)

	// GetPullRequest 获取合并请求详情
	GetPullRequest(ctx context.Context, number int) (*PullRequest, error)

	// CreatePullRequest 创建合并请求
	CreatePullRequest(ctx context.Context, req *CreatePRRequest) (*PullRequest, error)

	// UpdatePullRequest 更新合并请求
	UpdatePullRequest(ctx context.Context, number int, title, body string) (*PullRequest, error)

	// MergePullRequest 合并 PR
	MergePullRequest(ctx context.Context, number int, opts *MergeOptions) error

	// ClosePullRequest 关闭合并请求
	ClosePullRequest(ctx context.Context, number int) error

	// GetPullRequestCommits 获取 PR 包含的提交
	GetPullRequestCommits(ctx context.Context, number int) ([]*Commit, error)

	// ==================== 分支比对 ====================

	// CompareBranches 比较两个分支的差异
	// base: 基准分支
	// head: 对比分支
	CompareBranches(ctx context.Context, base, head string) (*CompareResult, error)

	// ==================== 评论功能 ====================

	// ListComments 列出 PR 评论
	ListComments(ctx context.Context, prNumber int) ([]*Comment, error)

	// CreateComment 添加评论
	CreateComment(ctx context.Context, prNumber int, body string) (*Comment, error)

	// UpdateComment 更新评论
	UpdateComment(ctx context.Context, commentID int64, body string) (*Comment, error)

	// DeleteComment 删除评论
	DeleteComment(ctx context.Context, commentID int64) error

	// ==================== 提交历史 ====================

	// GetCommit 获取提交详情
	GetCommit(ctx context.Context, sha string) (*Commit, error)

	// ListCommits 列出分支的提交历史
	ListCommits(ctx context.Context, branch string, opts *ListOptions) ([]*Commit, error)

	// ==================== CI/CD Pipeline 管理 ====================

	// TriggerPipeline 触发新的 Pipeline
	TriggerPipeline(ctx context.Context, opts *TriggerPipelineOptions) (*Pipeline, error)

	// GetPipeline 获取单个 Pipeline 详情
	GetPipeline(ctx context.Context, pipelineID int64) (*Pipeline, error)

	// ListPipelines 获取 Pipeline 列表
	ListPipelines(ctx context.Context, opts *ListPipelineOptions) ([]*Pipeline, error)

	// CancelPipeline 取消运行中的 Pipeline
	CancelPipeline(ctx context.Context, pipelineID int64) (*Pipeline, error)

	// RetryPipeline 重试失败的 Pipeline
	RetryPipeline(ctx context.Context, pipelineID int64) (*Pipeline, error)

	// ListPipelineJobs 获取 Pipeline 的作业列表
	ListPipelineJobs(ctx context.Context, pipelineID int64) ([]*PipelineJob, error)
}
