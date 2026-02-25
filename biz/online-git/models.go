package onlinegit

import "time"

// Platform 定义 Git 平台类型
type Platform string

const (
	PlatformGitHub Platform = "github"
	PlatformGitLab Platform = "gitlab"
	PlatformGitea  Platform = "gitea"
)

// PRState 定义 PR/MR 状态
type PRState string

const (
	PRStateOpen   PRState = "open"
	PRStateClosed PRState = "closed"
	PRStateMerged PRState = "merged"
	PRStateAll    PRState = "all"
)

// FileChangeType 定义文件变更类型
type FileChangeType string

const (
	FileChangeAdded    FileChangeType = "added"
	FileChangeModified FileChangeType = "modified"
	FileChangeDeleted  FileChangeType = "deleted"
	FileChangeRenamed  FileChangeType = "renamed"
)

// MergeMethod 定义合并方式
type MergeMethod string

const (
	MergeMethodMerge  MergeMethod = "merge"
	MergeMethodSquash MergeMethod = "squash"
	MergeMethodRebase MergeMethod = "rebase"
)

// Repository 仓库信息
type Repository struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	Description   string    `json:"description"`
	URL           string    `json:"url"`
	CloneURL      string    `json:"clone_url"`
	DefaultBranch string    `json:"default_branch"`
	Private       bool      `json:"private"`
	Fork          bool      `json:"fork"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Branch 分支信息
type Branch struct {
	Name      string  `json:"name"`
	CommitSHA string  `json:"commit_sha"`
	Protected bool    `json:"protected"`
	Default   bool    `json:"default"`
	Commit    *Commit `json:"commit,omitempty"`
}

// Commit 提交信息
type Commit struct {
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Author    *User     `json:"author"`
	Committer *User     `json:"committer,omitempty"`
	URL       string    `json:"url"`
	Parents   []string  `json:"parents,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// User 用户信息
type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// PullRequest 合并请求/PR
type PullRequest struct {
	ID           int64     `json:"id"`
	Number       int       `json:"number"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	State        PRState   `json:"state"`
	SourceBranch string    `json:"source_branch"`
	TargetBranch string    `json:"target_branch"`
	Author       *User     `json:"author"`
	Assignees    []*User   `json:"assignees,omitempty"`
	Labels       []string  `json:"labels,omitempty"`
	URL          string    `json:"url"`
	Merged       bool      `json:"merged"`
	MergedAt     time.Time `json:"merged_at,omitempty"`
	MergedBy     *User     `json:"merged_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at,omitempty"`
}

// Comment 评论
type Comment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	Author    *User     `json:"author"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CompareResult 分支比对结果
type CompareResult struct {
	BaseBranch   string        `json:"base_branch"`
	HeadBranch   string        `json:"head_branch"`
	AheadBy      int           `json:"ahead_by"`
	BehindBy     int           `json:"behind_by"`
	TotalCommits int           `json:"total_commits"`
	Commits      []*Commit     `json:"commits,omitempty"`
	Files        []*FileChange `json:"files,omitempty"`
	DiffStat     *DiffStat     `json:"diff_stat"`
}

// DiffStat 差异统计
type DiffStat struct {
	Additions    int `json:"additions"`
	Deletions    int `json:"deletions"`
	ChangedFiles int `json:"changed_files"`
}

// FileChange 文件变更
type FileChange struct {
	Filename     string         `json:"filename"`
	Status       FileChangeType `json:"status"`
	Additions    int            `json:"additions"`
	Deletions    int            `json:"deletions"`
	Changes      int            `json:"changes"`
	PreviousName string         `json:"previous_name,omitempty"`
	Patch        string         `json:"patch,omitempty"`
}

// ProtectionRules 分支保护规则
type ProtectionRules struct {
	RequiredReviews      int      `json:"required_reviews"`
	DismissStaleReviews  bool     `json:"dismiss_stale_reviews"`
	RequireCodeOwner     bool     `json:"require_code_owner"`
	RequiredStatusChecks []string `json:"required_status_checks,omitempty"`
	EnforceAdmins        bool     `json:"enforce_admins"`
	AllowForcePush       bool     `json:"allow_force_push"`
	AllowDeletions       bool     `json:"allow_deletions"`
}

// CreatePRRequest 创建 PR 请求参数
type CreatePRRequest struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	SourceBranch string   `json:"source_branch"`
	TargetBranch string   `json:"target_branch"`
	Labels       []string `json:"labels,omitempty"`
	Assignees    []string `json:"assignees,omitempty"`
	Draft        bool     `json:"draft"`
}

// MergeOptions 合并选项
type MergeOptions struct {
	Method        MergeMethod `json:"method"`
	CommitTitle   string      `json:"commit_title,omitempty"`
	CommitMessage string      `json:"commit_message,omitempty"`
	DeleteBranch  bool        `json:"delete_branch"`
}

// ListOptions 列表分页参数
type ListOptions struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// ProviderConfig Git 平台配置
type ProviderConfig struct {
	Platform        Platform `json:"platform"`
	BaseURL         string   `json:"base_url"`
	Token           string   `json:"token"`
	Owner           string   `json:"owner"`
	Repo            string   `json:"repo"`
	InsecureSkipTLS bool     `json:"insecure_skip_tls"` // 跳过 TLS 证书验证（用于私有化部署自签名证书）
}

// ==================== CI/CD Pipeline 相关 ====================

// PipelineStatus Pipeline 状态
type PipelineStatus string

const (
	PipelineStatusCreated            PipelineStatus = "created"
	PipelineStatusWaitingForResource PipelineStatus = "waiting_for_resource"
	PipelineStatusPreparing          PipelineStatus = "preparing"
	PipelineStatusPending            PipelineStatus = "pending"
	PipelineStatusRunning            PipelineStatus = "running"
	PipelineStatusSuccess            PipelineStatus = "success"
	PipelineStatusFailed             PipelineStatus = "failed"
	PipelineStatusCanceled           PipelineStatus = "canceled"
	PipelineStatusSkipped            PipelineStatus = "skipped"
	PipelineStatusManual             PipelineStatus = "manual"
	PipelineStatusScheduled          PipelineStatus = "scheduled"
)

// PipelineSource Pipeline 触发来源
type PipelineSource string

const (
	PipelineSourcePush              PipelineSource = "push"
	PipelineSourceWeb               PipelineSource = "web"
	PipelineSourceTrigger           PipelineSource = "trigger"
	PipelineSourceSchedule          PipelineSource = "schedule"
	PipelineSourceAPI               PipelineSource = "api"
	PipelineSourceMergeRequestEvent PipelineSource = "merge_request_event"
)

// Pipeline CI/CD Pipeline 信息
type Pipeline struct {
	ID             int64          `json:"id"`
	IID            int64          `json:"iid"` // GitLab 内部编号
	ProjectID      int64          `json:"project_id"`
	Status         PipelineStatus `json:"status"`
	Source         PipelineSource `json:"source"`
	Ref            string         `json:"ref"` // 分支/标签名
	SHA            string         `json:"sha"` // 提交 SHA
	WebURL         string         `json:"web_url"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	StartedAt      *time.Time     `json:"started_at,omitempty"`
	FinishedAt     *time.Time     `json:"finished_at,omitempty"`
	Duration       int64          `json:"duration"` // 秒
	QueuedDuration int64          `json:"queued_duration"`
	User           *User          `json:"user,omitempty"`
	CommitTitle    string         `json:"commit_title"`
}

// PipelineJob Pipeline 作业
type PipelineJob struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Stage        string         `json:"stage"`
	Status       PipelineStatus `json:"status"`
	Ref          string         `json:"ref"`
	WebURL       string         `json:"web_url"`
	CreatedAt    time.Time      `json:"created_at"`
	StartedAt    *time.Time     `json:"started_at,omitempty"`
	FinishedAt   *time.Time     `json:"finished_at,omitempty"`
	Duration     float64        `json:"duration"`
	User         *User          `json:"user,omitempty"`
	Pipeline     *Pipeline      `json:"pipeline,omitempty"`
	AllowFailure bool           `json:"allow_failure"`
}

// TriggerPipelineOptions 触发 Pipeline 选项
type TriggerPipelineOptions struct {
	Ref       string            `json:"ref"`       // 必需：分支/标签名
	Variables map[string]string `json:"variables"` // 可选：CI 变量
}

// ListPipelineOptions 查询 Pipeline 选项
type ListPipelineOptions struct {
	Ref      string         `json:"ref,omitempty"`      // 按分支过滤
	Status   PipelineStatus `json:"status,omitempty"`   // 按状态过滤
	Source   PipelineSource `json:"source,omitempty"`   // 按来源过滤
	Username string         `json:"username,omitempty"` // 按用户过滤
	OrderBy  string         `json:"order_by,omitempty"` // id/status/ref/updated_at/user_id
	Sort     string         `json:"sort,omitempty"`     // asc/desc
	Page     int            `json:"page,omitempty"`
	PerPage  int            `json:"per_page,omitempty"`
}
