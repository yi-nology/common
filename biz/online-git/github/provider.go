package github

import (
	"context"
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

func (p *Provider) GetRepository(ctx context.Context) (*onlinegit.Repository, error) {
	repo, resp, err := p.client.Repositories.Get(ctx, p.owner, p.repo)
	if err != nil {
		return nil, p.wrapError("GetRepository", resp, err)
	}

	return &onlinegit.Repository{
		ID:            repo.GetID(),
		Name:          repo.GetName(),
		FullName:      repo.GetFullName(),
		Description:   repo.GetDescription(),
		URL:           repo.GetHTMLURL(),
		CloneURL:      repo.GetCloneURL(),
		DefaultBranch: repo.GetDefaultBranch(),
		Private:       repo.GetPrivate(),
		Fork:          repo.GetFork(),
		CreatedAt:     repo.GetCreatedAt().Time,
		UpdatedAt:     repo.GetUpdatedAt().Time,
	}, nil
}

func (p *Provider) ListBranches(ctx context.Context, opts *onlinegit.ListOptions) ([]*onlinegit.Branch, error) {
	ghOpts := &github.BranchListOptions{
		ListOptions: github.ListOptions{
			Page:    opts.Page,
			PerPage: opts.PerPage,
		},
	}

	branches, resp, err := p.client.Repositories.ListBranches(ctx, p.owner, p.repo, ghOpts)
	if err != nil {
		return nil, p.wrapError("ListBranches", resp, err)
	}

	result := make([]*onlinegit.Branch, len(branches))
	for i, b := range branches {
		result[i] = &onlinegit.Branch{
			Name:      b.GetName(),
			CommitSHA: b.GetCommit().GetSHA(),
			Protected: b.GetProtected(),
		}
	}
	return result, nil
}

func (p *Provider) GetBranch(ctx context.Context, name string) (*onlinegit.Branch, error) {
	branch, resp, err := p.client.Repositories.GetBranch(ctx, p.owner, p.repo, name, 0)
	if err != nil {
		return nil, p.wrapError("GetBranch", resp, err)
	}

	return &onlinegit.Branch{
		Name:      branch.GetName(),
		CommitSHA: branch.GetCommit().GetSHA(),
		Protected: branch.GetProtected(),
		Commit:    p.toCommit(branch.GetCommit().Commit, branch.GetCommit().GetSHA()),
	}, nil
}

func (p *Provider) CreateBranch(ctx context.Context, name, sourceBranch string) (*onlinegit.Branch, error) {
	// 获取源分支的 SHA
	sourceBranchInfo, resp, err := p.client.Repositories.GetBranch(ctx, p.owner, p.repo, sourceBranch, 0)
	if err != nil {
		return nil, p.wrapError("CreateBranch:GetSourceBranch", resp, err)
	}

	sha := sourceBranchInfo.GetCommit().GetSHA()

	// 创建新分支的 ref
	ref := &github.Reference{
		Ref:    github.String("refs/heads/" + name),
		Object: &github.GitObject{SHA: github.String(sha)},
	}

	_, resp, err = p.client.Git.CreateRef(ctx, p.owner, p.repo, ref)
	if err != nil {
		return nil, p.wrapError("CreateBranch", resp, err)
	}

	return &onlinegit.Branch{
		Name:      name,
		CommitSHA: sha,
		Protected: false,
	}, nil
}

func (p *Provider) DeleteBranch(ctx context.Context, name string) error {
	ref := "heads/" + name
	resp, err := p.client.Git.DeleteRef(ctx, p.owner, p.repo, ref)
	if err != nil {
		return p.wrapError("DeleteBranch", resp, err)
	}
	return nil
}

func (p *Provider) SetBranchProtection(ctx context.Context, name string, rules *onlinegit.ProtectionRules) error {
	req := &github.ProtectionRequest{
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			RequiredApprovingReviewCount: rules.RequiredReviews,
			DismissStaleReviews:          rules.DismissStaleReviews,
			RequireCodeOwnerReviews:      rules.RequireCodeOwner,
		},
		EnforceAdmins:    rules.EnforceAdmins,
		AllowForcePushes: github.Bool(rules.AllowForcePush),
		AllowDeletions:   github.Bool(rules.AllowDeletions),
	}

	if len(rules.RequiredStatusChecks) > 0 {
		checks := make([]*github.RequiredStatusCheck, len(rules.RequiredStatusChecks))
		for i, c := range rules.RequiredStatusChecks {
			checks[i] = &github.RequiredStatusCheck{Context: c}
		}
		req.RequiredStatusChecks = &github.RequiredStatusChecks{
			Strict: true,
			Checks: checks,
		}
	}

	_, resp, err := p.client.Repositories.UpdateBranchProtection(ctx, p.owner, p.repo, name, req)
	if err != nil {
		return p.wrapError("SetBranchProtection", resp, err)
	}
	return nil
}

func (p *Provider) UnsetBranchProtection(ctx context.Context, name string) error {
	resp, err := p.client.Repositories.RemoveBranchProtection(ctx, p.owner, p.repo, name)
	if err != nil {
		return p.wrapError("UnsetBranchProtection", resp, err)
	}
	return nil
}

func (p *Provider) ListPullRequests(ctx context.Context, state onlinegit.PRState, opts *onlinegit.ListOptions) ([]*onlinegit.PullRequest, error) {
	ghState := "all"
	if state != onlinegit.PRStateAll {
		ghState = string(state)
	}

	ghOpts := &github.PullRequestListOptions{
		State: ghState,
		ListOptions: github.ListOptions{
			Page:    opts.Page,
			PerPage: opts.PerPage,
		},
	}

	prs, resp, err := p.client.PullRequests.List(ctx, p.owner, p.repo, ghOpts)
	if err != nil {
		return nil, p.wrapError("ListPullRequests", resp, err)
	}

	result := make([]*onlinegit.PullRequest, len(prs))
	for i, pr := range prs {
		result[i] = p.toPullRequest(pr)
	}
	return result, nil
}

func (p *Provider) GetPullRequest(ctx context.Context, number int) (*onlinegit.PullRequest, error) {
	pr, resp, err := p.client.PullRequests.Get(ctx, p.owner, p.repo, number)
	if err != nil {
		return nil, p.wrapError("GetPullRequest", resp, err)
	}
	return p.toPullRequest(pr), nil
}

func (p *Provider) CreatePullRequest(ctx context.Context, req *onlinegit.CreatePRRequest) (*onlinegit.PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(req.Title),
		Body:  github.String(req.Body),
		Head:  github.String(req.SourceBranch),
		Base:  github.String(req.TargetBranch),
		Draft: github.Bool(req.Draft),
	}

	pr, resp, err := p.client.PullRequests.Create(ctx, p.owner, p.repo, newPR)
	if err != nil {
		return nil, p.wrapError("CreatePullRequest", resp, err)
	}

	// 添加标签
	if len(req.Labels) > 0 {
		_, _, _ = p.client.Issues.AddLabelsToIssue(ctx, p.owner, p.repo, pr.GetNumber(), req.Labels)
	}

	// 添加指派人
	if len(req.Assignees) > 0 {
		_, _, _ = p.client.Issues.AddAssignees(ctx, p.owner, p.repo, pr.GetNumber(), req.Assignees)
	}

	return p.toPullRequest(pr), nil
}

func (p *Provider) UpdatePullRequest(ctx context.Context, number int, title, body string) (*onlinegit.PullRequest, error) {
	update := &github.PullRequest{
		Title: github.String(title),
		Body:  github.String(body),
	}

	pr, resp, err := p.client.PullRequests.Edit(ctx, p.owner, p.repo, number, update)
	if err != nil {
		return nil, p.wrapError("UpdatePullRequest", resp, err)
	}
	return p.toPullRequest(pr), nil
}

func (p *Provider) MergePullRequest(ctx context.Context, number int, opts *onlinegit.MergeOptions) error {
	method := "merge"
	if opts != nil && opts.Method != "" {
		method = string(opts.Method)
	}

	mergeOpts := &github.PullRequestOptions{
		MergeMethod: method,
	}
	if opts != nil {
		if opts.CommitTitle != "" {
			mergeOpts.CommitTitle = opts.CommitTitle
		}
	}

	commitMsg := ""
	if opts != nil && opts.CommitMessage != "" {
		commitMsg = opts.CommitMessage
	}

	_, resp, err := p.client.PullRequests.Merge(ctx, p.owner, p.repo, number, commitMsg, mergeOpts)
	if err != nil {
		return p.wrapError("MergePullRequest", resp, err)
	}

	// 删除源分支
	if opts != nil && opts.DeleteBranch {
		pr, _, err := p.client.PullRequests.Get(ctx, p.owner, p.repo, number)
		if err == nil && pr.Head != nil && pr.Head.Ref != nil {
			_ = p.DeleteBranch(ctx, *pr.Head.Ref)
		}
	}

	return nil
}

func (p *Provider) ClosePullRequest(ctx context.Context, number int) error {
	state := "closed"
	update := &github.PullRequest{
		State: &state,
	}

	_, resp, err := p.client.PullRequests.Edit(ctx, p.owner, p.repo, number, update)
	if err != nil {
		return p.wrapError("ClosePullRequest", resp, err)
	}
	return nil
}

func (p *Provider) GetPullRequestCommits(ctx context.Context, number int) ([]*onlinegit.Commit, error) {
	commits, resp, err := p.client.PullRequests.ListCommits(ctx, p.owner, p.repo, number, nil)
	if err != nil {
		return nil, p.wrapError("GetPullRequestCommits", resp, err)
	}

	result := make([]*onlinegit.Commit, len(commits))
	for i, c := range commits {
		result[i] = p.toCommit(c.Commit, c.GetSHA())
	}
	return result, nil
}

func (p *Provider) CompareBranches(ctx context.Context, base, head string) (*onlinegit.CompareResult, error) {
	comparison, resp, err := p.client.Repositories.CompareCommits(ctx, p.owner, p.repo, base, head, nil)
	if err != nil {
		return nil, p.wrapError("CompareBranches", resp, err)
	}

	commits := make([]*onlinegit.Commit, len(comparison.Commits))
	for i, c := range comparison.Commits {
		commits[i] = p.toCommit(c.Commit, c.GetSHA())
	}

	files := make([]*onlinegit.FileChange, len(comparison.Files))
	for i, f := range comparison.Files {
		files[i] = &onlinegit.FileChange{
			Filename:     f.GetFilename(),
			Status:       onlinegit.FileChangeType(f.GetStatus()),
			Additions:    f.GetAdditions(),
			Deletions:    f.GetDeletions(),
			Changes:      f.GetChanges(),
			PreviousName: f.GetPreviousFilename(),
			Patch:        f.GetPatch(),
		}
	}

	return &onlinegit.CompareResult{
		BaseBranch:   base,
		HeadBranch:   head,
		AheadBy:      comparison.GetAheadBy(),
		BehindBy:     comparison.GetBehindBy(),
		TotalCommits: comparison.GetTotalCommits(),
		Commits:      commits,
		Files:        files,
		DiffStat: &onlinegit.DiffStat{
			Additions:    0,
			Deletions:    0,
			ChangedFiles: len(comparison.Files),
		},
	}, nil
}

func (p *Provider) ListComments(ctx context.Context, prNumber int) ([]*onlinegit.Comment, error) {
	comments, resp, err := p.client.Issues.ListComments(ctx, p.owner, p.repo, prNumber, nil)
	if err != nil {
		return nil, p.wrapError("ListComments", resp, err)
	}

	result := make([]*onlinegit.Comment, len(comments))
	for i, c := range comments {
		result[i] = p.toComment(c)
	}
	return result, nil
}

func (p *Provider) CreateComment(ctx context.Context, prNumber int, body string) (*onlinegit.Comment, error) {
	comment := &github.IssueComment{
		Body: github.String(body),
	}

	created, resp, err := p.client.Issues.CreateComment(ctx, p.owner, p.repo, prNumber, comment)
	if err != nil {
		return nil, p.wrapError("CreateComment", resp, err)
	}
	return p.toComment(created), nil
}

func (p *Provider) UpdateComment(ctx context.Context, commentID int64, body string) (*onlinegit.Comment, error) {
	comment := &github.IssueComment{
		Body: github.String(body),
	}

	updated, resp, err := p.client.Issues.EditComment(ctx, p.owner, p.repo, commentID, comment)
	if err != nil {
		return nil, p.wrapError("UpdateComment", resp, err)
	}
	return p.toComment(updated), nil
}

func (p *Provider) DeleteComment(ctx context.Context, commentID int64) error {
	resp, err := p.client.Issues.DeleteComment(ctx, p.owner, p.repo, commentID)
	if err != nil {
		return p.wrapError("DeleteComment", resp, err)
	}
	return nil
}

func (p *Provider) GetCommit(ctx context.Context, sha string) (*onlinegit.Commit, error) {
	commit, resp, err := p.client.Repositories.GetCommit(ctx, p.owner, p.repo, sha, nil)
	if err != nil {
		return nil, p.wrapError("GetCommit", resp, err)
	}
	return p.toCommit(commit.Commit, commit.GetSHA()), nil
}

func (p *Provider) ListCommits(ctx context.Context, branch string, opts *onlinegit.ListOptions) ([]*onlinegit.Commit, error) {
	ghOpts := &github.CommitsListOptions{
		SHA: branch,
		ListOptions: github.ListOptions{
			Page:    opts.Page,
			PerPage: opts.PerPage,
		},
	}

	commits, resp, err := p.client.Repositories.ListCommits(ctx, p.owner, p.repo, ghOpts)
	if err != nil {
		return nil, p.wrapError("ListCommits", resp, err)
	}

	result := make([]*onlinegit.Commit, len(commits))
	for i, c := range commits {
		result[i] = p.toCommit(c.Commit, c.GetSHA())
	}
	return result, nil
}

// ==================== CI/CD Pipeline 管理 (GitHub Actions) ====================

// TriggerPipeline 触发 GitHub Actions Workflow
// 通过 workflow_dispatch 事件触发指定的 workflow
func (p *Provider) TriggerPipeline(ctx context.Context, opts *onlinegit.TriggerPipelineOptions) (*onlinegit.Pipeline, error) {
	if opts == nil || opts.Ref == "" {
		return nil, onlinegit.NewProviderError(onlinegit.PlatformGitHub, "TriggerPipeline", fmt.Errorf("ref is required"), "")
	}

	// GitHub Actions 需要指定 workflow 文件名，默认使用 ci.yml
	workflowFile := "ci.yml"
	if opts.Variables != nil {
		if wf, ok := opts.Variables["workflow_file"]; ok {
			workflowFile = wf
			delete(opts.Variables, "workflow_file")
		}
	}

	// 构建 inputs 参数
	inputs := make(map[string]interface{})
	for k, v := range opts.Variables {
		inputs[k] = v
	}

	event := github.CreateWorkflowDispatchEventRequest{
		Ref:    opts.Ref,
		Inputs: inputs,
	}

	resp, err := p.client.Actions.CreateWorkflowDispatchEventByFileName(ctx, p.owner, p.repo, workflowFile, event)
	if err != nil {
		return nil, p.wrapError("TriggerPipeline", resp, err)
	}

	// GitHub Actions 的 dispatch 不直接返回 run，需要查询最新的 run
	// 等待一小段时间后查询
	listOpts := &github.ListWorkflowRunsOptions{
		Branch: opts.Ref,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	}

	runs, _, err := p.client.Actions.ListRepositoryWorkflowRuns(ctx, p.owner, p.repo, listOpts)
	if err != nil || len(runs.WorkflowRuns) == 0 {
		// 返回一个占位的 Pipeline 对象，表示触发成功
		return &onlinegit.Pipeline{
			Ref:    opts.Ref,
			Status: onlinegit.PipelineStatusPending,
			Source: onlinegit.PipelineSourceAPI,
		}, nil
	}

	return p.toWorkflowRunPipeline(runs.WorkflowRuns[0]), nil
}

// GetPipeline 获取 GitHub Actions Workflow Run 详情
func (p *Provider) GetPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	run, resp, err := p.client.Actions.GetWorkflowRunByID(ctx, p.owner, p.repo, pipelineID)
	if err != nil {
		return nil, p.wrapError("GetPipeline", resp, err)
	}
	return p.toWorkflowRunPipeline(run), nil
}

// ListPipelines 获取 GitHub Actions Workflow Run 列表
func (p *Provider) ListPipelines(ctx context.Context, opts *onlinegit.ListPipelineOptions) ([]*onlinegit.Pipeline, error) {
	listOpts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			Page:    opts.Page,
			PerPage: opts.PerPage,
		},
	}

	if opts.Ref != "" {
		listOpts.Branch = opts.Ref
	}
	if opts.Status != "" {
		listOpts.Status = string(opts.Status)
	}
	if opts.Username != "" {
		listOpts.Actor = opts.Username
	}

	runs, resp, err := p.client.Actions.ListRepositoryWorkflowRuns(ctx, p.owner, p.repo, listOpts)
	if err != nil {
		return nil, p.wrapError("ListPipelines", resp, err)
	}

	result := make([]*onlinegit.Pipeline, len(runs.WorkflowRuns))
	for i, run := range runs.WorkflowRuns {
		result[i] = p.toWorkflowRunPipeline(run)
	}
	return result, nil
}

// CancelPipeline 取消 GitHub Actions Workflow Run
func (p *Provider) CancelPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	resp, err := p.client.Actions.CancelWorkflowRunByID(ctx, p.owner, p.repo, pipelineID)
	if err != nil {
		return nil, p.wrapError("CancelPipeline", resp, err)
	}

	// 获取更新后的状态
	return p.GetPipeline(ctx, pipelineID)
}

// RetryPipeline 重试 GitHub Actions Workflow Run
func (p *Provider) RetryPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	resp, err := p.client.Actions.RerunWorkflowByID(ctx, p.owner, p.repo, pipelineID)
	if err != nil {
		return nil, p.wrapError("RetryPipeline", resp, err)
	}

	// 获取更新后的状态
	return p.GetPipeline(ctx, pipelineID)
}

// ListPipelineJobs 获取 GitHub Actions Workflow Run 的作业列表
func (p *Provider) ListPipelineJobs(ctx context.Context, pipelineID int64) ([]*onlinegit.PipelineJob, error) {
	jobs, resp, err := p.client.Actions.ListWorkflowJobs(ctx, p.owner, p.repo, pipelineID, nil)
	if err != nil {
		return nil, p.wrapError("ListPipelineJobs", resp, err)
	}

	result := make([]*onlinegit.PipelineJob, len(jobs.Jobs))
	for i, job := range jobs.Jobs {
		result[i] = p.toWorkflowJob(job)
	}
	return result, nil
}

// toWorkflowRunPipeline 将 GitHub WorkflowRun 转换为统一的 Pipeline 结构
func (p *Provider) toWorkflowRunPipeline(run *github.WorkflowRun) *onlinegit.Pipeline {
	pipeline := &onlinegit.Pipeline{
		ID:     run.GetID(),
		IID:    int64(run.GetRunNumber()),
		Ref:    run.GetHeadBranch(),
		SHA:    run.GetHeadSHA(),
		WebURL: run.GetHTMLURL(),
		Status: p.mapWorkflowRunStatus(run.GetStatus(), run.GetConclusion()),
		Source: p.mapWorkflowRunEvent(run.GetEvent()),
	}

	if run.CreatedAt != nil {
		pipeline.CreatedAt = run.CreatedAt.Time
	}
	if run.UpdatedAt != nil {
		pipeline.UpdatedAt = run.UpdatedAt.Time
	}
	if run.RunStartedAt != nil {
		startedAt := run.RunStartedAt.Time
		pipeline.StartedAt = &startedAt
	}

	if run.Actor != nil {
		pipeline.User = &onlinegit.User{
			ID:        run.Actor.GetID(),
			Login:     run.Actor.GetLogin(),
			Name:      run.Actor.GetName(),
			AvatarURL: run.Actor.GetAvatarURL(),
		}
	}

	return pipeline
}

// mapWorkflowRunStatus 映射 GitHub Workflow Run 状态到统一状态
func (p *Provider) mapWorkflowRunStatus(status, conclusion string) onlinegit.PipelineStatus {
	switch status {
	case "queued":
		return onlinegit.PipelineStatusPending
	case "in_progress":
		return onlinegit.PipelineStatusRunning
	case "completed":
		switch conclusion {
		case "success":
			return onlinegit.PipelineStatusSuccess
		case "failure":
			return onlinegit.PipelineStatusFailed
		case "cancelled":
			return onlinegit.PipelineStatusCanceled
		case "skipped":
			return onlinegit.PipelineStatusSkipped
		default:
			return onlinegit.PipelineStatusFailed
		}
	case "waiting":
		return onlinegit.PipelineStatusWaitingForResource
	default:
		return onlinegit.PipelineStatusCreated
	}
}

// mapWorkflowRunEvent 映射 GitHub Workflow Run 事件到 Pipeline Source
func (p *Provider) mapWorkflowRunEvent(event string) onlinegit.PipelineSource {
	switch event {
	case "push":
		return onlinegit.PipelineSourcePush
	case "workflow_dispatch":
		return onlinegit.PipelineSourceWeb
	case "schedule":
		return onlinegit.PipelineSourceSchedule
	case "pull_request", "pull_request_target":
		return onlinegit.PipelineSourceMergeRequestEvent
	case "repository_dispatch":
		return onlinegit.PipelineSourceAPI
	default:
		return onlinegit.PipelineSourcePush
	}
}

// toWorkflowJob 将 GitHub WorkflowJob 转换为统一的 PipelineJob 结构
func (p *Provider) toWorkflowJob(job *github.WorkflowJob) *onlinegit.PipelineJob {
	pipelineJob := &onlinegit.PipelineJob{
		ID:     job.GetID(),
		Name:   job.GetName(),
		Status: p.mapWorkflowJobStatus(job.GetStatus(), job.GetConclusion()),
		WebURL: job.GetHTMLURL(),
	}

	if job.StartedAt != nil {
		startedAt := job.StartedAt.Time
		pipelineJob.StartedAt = &startedAt
	}
	if job.CompletedAt != nil {
		completedAt := job.CompletedAt.Time
		pipelineJob.FinishedAt = &completedAt
		if job.StartedAt != nil {
			pipelineJob.Duration = job.CompletedAt.Sub(job.StartedAt.Time).Seconds()
		}
	}

	return pipelineJob
}

// mapWorkflowJobStatus 映射 GitHub Workflow Job 状态
func (p *Provider) mapWorkflowJobStatus(status, conclusion string) onlinegit.PipelineStatus {
	switch status {
	case "queued":
		return onlinegit.PipelineStatusPending
	case "in_progress":
		return onlinegit.PipelineStatusRunning
	case "completed":
		switch conclusion {
		case "success":
			return onlinegit.PipelineStatusSuccess
		case "failure":
			return onlinegit.PipelineStatusFailed
		case "cancelled":
			return onlinegit.PipelineStatusCanceled
		case "skipped":
			return onlinegit.PipelineStatusSkipped
		default:
			return onlinegit.PipelineStatusFailed
		}
	default:
		return onlinegit.PipelineStatusCreated
	}
}

// Helper methods

func (p *Provider) toPullRequest(pr *github.PullRequest) *onlinegit.PullRequest {
	state := onlinegit.PRStateOpen
	if pr.GetMerged() {
		state = onlinegit.PRStateMerged
	} else if pr.GetState() == "closed" {
		state = onlinegit.PRStateClosed
	}

	result := &onlinegit.PullRequest{
		ID:           pr.GetID(),
		Number:       pr.GetNumber(),
		Title:        pr.GetTitle(),
		Body:         pr.GetBody(),
		State:        state,
		SourceBranch: pr.GetHead().GetRef(),
		TargetBranch: pr.GetBase().GetRef(),
		URL:          pr.GetHTMLURL(),
		Merged:       pr.GetMerged(),
		CreatedAt:    pr.GetCreatedAt().Time,
		UpdatedAt:    pr.GetUpdatedAt().Time,
	}

	if pr.GetUser() != nil {
		result.Author = p.toUser(pr.GetUser())
	}

	if pr.MergedAt != nil {
		result.MergedAt = pr.GetMergedAt().Time
	}

	if pr.ClosedAt != nil {
		result.ClosedAt = pr.GetClosedAt().Time
	}

	if pr.MergedBy != nil {
		result.MergedBy = p.toUser(pr.MergedBy)
	}

	if len(pr.Labels) > 0 {
		result.Labels = make([]string, len(pr.Labels))
		for i, l := range pr.Labels {
			result.Labels[i] = l.GetName()
		}
	}

	if len(pr.Assignees) > 0 {
		result.Assignees = make([]*onlinegit.User, len(pr.Assignees))
		for i, a := range pr.Assignees {
			result.Assignees[i] = p.toUser(a)
		}
	}

	return result
}

func (p *Provider) toUser(u *github.User) *onlinegit.User {
	return &onlinegit.User{
		ID:        u.GetID(),
		Login:     u.GetLogin(),
		Name:      u.GetName(),
		Email:     u.GetEmail(),
		AvatarURL: u.GetAvatarURL(),
	}
}

func (p *Provider) toCommit(c *github.Commit, sha string) *onlinegit.Commit {
	if c == nil {
		return nil
	}

	result := &onlinegit.Commit{
		SHA:     sha,
		Message: c.GetMessage(),
		URL:     c.GetURL(),
	}

	if c.Author != nil {
		result.Author = &onlinegit.User{
			Name:  c.Author.GetName(),
			Email: c.Author.GetEmail(),
		}
		if c.Author.Date != nil {
			result.CreatedAt = c.Author.Date.Time
		}
	}

	if c.Committer != nil {
		result.Committer = &onlinegit.User{
			Name:  c.Committer.GetName(),
			Email: c.Committer.GetEmail(),
		}
	}

	return result
}

func (p *Provider) toComment(c *github.IssueComment) *onlinegit.Comment {
	result := &onlinegit.Comment{
		ID:        c.GetID(),
		Body:      c.GetBody(),
		URL:       c.GetHTMLURL(),
		CreatedAt: c.GetCreatedAt().Time,
		UpdatedAt: c.GetUpdatedAt().Time,
	}

	if c.User != nil {
		result.Author = p.toUser(c.User)
	}

	return result
}

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
