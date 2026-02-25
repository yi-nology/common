package gitea

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

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
	client *gitea.Client
	owner  string
	repo   string
}

// NewProvider 创建 Gitea Provider
func NewProvider(cfg *onlinegit.ProviderConfig) (*Provider, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required for Gitea")
	}

	opts := []gitea.ClientOption{
		gitea.SetToken(cfg.Token),
	}

	if cfg.InsecureSkipTLS {
		httpClient := &http.Client{
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

	return &Provider{
		client: client,
		owner:  cfg.Owner,
		repo:   cfg.Repo,
	}, nil
}

func (p *Provider) GetPlatform() onlinegit.Platform {
	return onlinegit.PlatformGitea
}

func (p *Provider) GetRepository(ctx context.Context) (*onlinegit.Repository, error) {
	repo, resp, err := p.client.GetRepo(p.owner, p.repo)
	if err != nil {
		return nil, p.wrapError("GetRepository", resp, err)
	}

	return &onlinegit.Repository{
		ID:            repo.ID,
		Name:          repo.Name,
		FullName:      repo.FullName,
		Description:   repo.Description,
		URL:           repo.HTMLURL,
		CloneURL:      repo.CloneURL,
		DefaultBranch: repo.DefaultBranch,
		Private:       repo.Private,
		Fork:          repo.Fork,
		CreatedAt:     repo.Created,
		UpdatedAt:     repo.Updated,
	}, nil
}

func (p *Provider) ListBranches(ctx context.Context, opts *onlinegit.ListOptions) ([]*onlinegit.Branch, error) {
	giteaOpts := gitea.ListRepoBranchesOptions{
		ListOptions: gitea.ListOptions{
			Page:     opts.Page,
			PageSize: opts.PerPage,
		},
	}

	branches, resp, err := p.client.ListRepoBranches(p.owner, p.repo, giteaOpts)
	if err != nil {
		return nil, p.wrapError("ListBranches", resp, err)
	}

	result := make([]*onlinegit.Branch, len(branches))
	for i, b := range branches {
		result[i] = &onlinegit.Branch{
			Name:      b.Name,
			CommitSHA: b.Commit.ID,
			Protected: b.Protected,
		}
	}
	return result, nil
}

func (p *Provider) GetBranch(ctx context.Context, name string) (*onlinegit.Branch, error) {
	branch, resp, err := p.client.GetRepoBranch(p.owner, p.repo, name)
	if err != nil {
		return nil, p.wrapError("GetBranch", resp, err)
	}

	return &onlinegit.Branch{
		Name:      branch.Name,
		CommitSHA: branch.Commit.ID,
		Protected: branch.Protected,
		Commit:    p.toBranchCommit(branch.Commit),
	}, nil
}

func (p *Provider) CreateBranch(ctx context.Context, name, sourceBranch string) (*onlinegit.Branch, error) {
	opts := gitea.CreateBranchOption{
		BranchName:    name,
		OldBranchName: sourceBranch,
	}

	branch, resp, err := p.client.CreateBranch(p.owner, p.repo, opts)
	if err != nil {
		return nil, p.wrapError("CreateBranch", resp, err)
	}

	return &onlinegit.Branch{
		Name:      branch.Name,
		CommitSHA: branch.Commit.ID,
		Protected: branch.Protected,
	}, nil
}

func (p *Provider) DeleteBranch(ctx context.Context, name string) error {
	deleted, resp, err := p.client.DeleteRepoBranch(p.owner, p.repo, name)
	if err != nil {
		return p.wrapError("DeleteBranch", resp, err)
	}
	if !deleted {
		return onlinegit.NewProviderError(onlinegit.PlatformGitea, "DeleteBranch", fmt.Errorf("failed to delete branch"), "")
	}
	return nil
}

func (p *Provider) SetBranchProtection(ctx context.Context, name string, rules *onlinegit.ProtectionRules) error {
	opts := gitea.CreateBranchProtectionOption{
		BranchName:             name,
		EnablePush:             true,
		EnablePushWhitelist:    false,
		RequiredApprovals:      int64(rules.RequiredReviews),
		EnableStatusCheck:      len(rules.RequiredStatusChecks) > 0,
		StatusCheckContexts:    rules.RequiredStatusChecks,
		DismissStaleApprovals:  rules.DismissStaleReviews,
		BlockOnRejectedReviews: true,
		BlockOnOutdatedBranch:  true,
	}

	_, resp, err := p.client.CreateBranchProtection(p.owner, p.repo, opts)
	if err != nil {
		return p.wrapError("SetBranchProtection", resp, err)
	}
	return nil
}

func (p *Provider) UnsetBranchProtection(ctx context.Context, name string) error {
	resp, err := p.client.DeleteBranchProtection(p.owner, p.repo, name)
	if err != nil {
		return p.wrapError("UnsetBranchProtection", resp, err)
	}
	return nil
}

func (p *Provider) ListPullRequests(ctx context.Context, state onlinegit.PRState, opts *onlinegit.ListOptions) ([]*onlinegit.PullRequest, error) {
	giteaState := gitea.StateAll
	switch state {
	case onlinegit.PRStateOpen:
		giteaState = gitea.StateOpen
	case onlinegit.PRStateClosed:
		giteaState = gitea.StateClosed
	}

	giteaOpts := gitea.ListPullRequestsOptions{
		State: giteaState,
		ListOptions: gitea.ListOptions{
			Page:     opts.Page,
			PageSize: opts.PerPage,
		},
	}

	prs, resp, err := p.client.ListRepoPullRequests(p.owner, p.repo, giteaOpts)
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
	pr, resp, err := p.client.GetPullRequest(p.owner, p.repo, int64(number))
	if err != nil {
		return nil, p.wrapError("GetPullRequest", resp, err)
	}
	return p.toPullRequest(pr), nil
}

func (p *Provider) CreatePullRequest(ctx context.Context, req *onlinegit.CreatePRRequest) (*onlinegit.PullRequest, error) {
	opts := gitea.CreatePullRequestOption{
		Head:   req.SourceBranch,
		Base:   req.TargetBranch,
		Title:  req.Title,
		Body:   req.Body,
		Labels: []int64{}, // Gitea 需要 label IDs
	}

	if len(req.Assignees) > 0 {
		opts.Assignees = req.Assignees
	}

	pr, resp, err := p.client.CreatePullRequest(p.owner, p.repo, opts)
	if err != nil {
		return nil, p.wrapError("CreatePullRequest", resp, err)
	}

	return p.toPullRequest(pr), nil
}

func (p *Provider) UpdatePullRequest(ctx context.Context, number int, title, body string) (*onlinegit.PullRequest, error) {
	opts := gitea.EditPullRequestOption{
		Title: title,
		Body:  &body,
	}

	pr, resp, err := p.client.EditPullRequest(p.owner, p.repo, int64(number), opts)
	if err != nil {
		return nil, p.wrapError("UpdatePullRequest", resp, err)
	}
	return p.toPullRequest(pr), nil
}

func (p *Provider) MergePullRequest(ctx context.Context, number int, opts *onlinegit.MergeOptions) error {
	mergeStyle := gitea.MergeStyleMerge
	if opts != nil {
		switch opts.Method {
		case onlinegit.MergeMethodSquash:
			mergeStyle = gitea.MergeStyleSquash
		case onlinegit.MergeMethodRebase:
			mergeStyle = gitea.MergeStyleRebase
		}
	}

	mergeOpts := gitea.MergePullRequestOption{
		Style: mergeStyle,
	}

	if opts != nil {
		if opts.CommitTitle != "" {
			mergeOpts.Title = opts.CommitTitle
		}
		if opts.CommitMessage != "" {
			mergeOpts.Message = opts.CommitMessage
		}
		if opts.DeleteBranch {
			mergeOpts.DeleteBranchAfterMerge = true
		}
	}

	merged, resp, err := p.client.MergePullRequest(p.owner, p.repo, int64(number), mergeOpts)
	if err != nil {
		return p.wrapError("MergePullRequest", resp, err)
	}
	if !merged {
		return onlinegit.NewProviderError(onlinegit.PlatformGitea, "MergePullRequest", onlinegit.ErrNotMergeable, "")
	}
	return nil
}

func (p *Provider) ClosePullRequest(ctx context.Context, number int) error {
	opts := gitea.EditPullRequestOption{
		State: &[]gitea.StateType{gitea.StateClosed}[0],
	}

	_, resp, err := p.client.EditPullRequest(p.owner, p.repo, int64(number), opts)
	if err != nil {
		return p.wrapError("ClosePullRequest", resp, err)
	}
	return nil
}

func (p *Provider) GetPullRequestCommits(ctx context.Context, number int) ([]*onlinegit.Commit, error) {
	commits, resp, err := p.client.ListPullRequestCommits(p.owner, p.repo, int64(number), gitea.ListPullRequestCommitsOptions{})
	if err != nil {
		return nil, p.wrapError("GetPullRequestCommits", resp, err)
	}

	result := make([]*onlinegit.Commit, len(commits))
	for i, c := range commits {
		result[i] = p.toRepoCommit(c)
	}
	return result, nil
}

func (p *Provider) CompareBranches(ctx context.Context, base, head string) (*onlinegit.CompareResult, error) {
	// Gitea 没有直接的 Compare API，我们通过 commits 差异来模拟
	// 获取 head 分支的提交历史
	commits, _, err := p.client.ListRepoCommits(p.owner, p.repo, gitea.ListCommitOptions{
		SHA: head,
		ListOptions: gitea.ListOptions{
			Page:     1,
			PageSize: 100,
		},
	})
	if err != nil {
		return nil, p.wrapError("CompareBranches:ListCommits", nil, err)
	}

	// 获取 base 分支的最新提交
	baseBranch, _, err := p.client.GetRepoBranch(p.owner, p.repo, base)
	if err != nil {
		return nil, p.wrapError("CompareBranches:GetBaseBranch", nil, err)
	}

	// 找出不在 base 分支中的提交
	var diffCommits []*onlinegit.Commit
	baseSHA := baseBranch.Commit.ID
	for _, c := range commits {
		if c.SHA == baseSHA {
			break
		}
		diffCommits = append(diffCommits, p.toRepoCommit(c))
	}

	return &onlinegit.CompareResult{
		BaseBranch:   base,
		HeadBranch:   head,
		TotalCommits: len(diffCommits),
		Commits:      diffCommits,
		DiffStat: &onlinegit.DiffStat{
			ChangedFiles: 0, // Gitea API 不直接提供
		},
	}, nil
}

func (p *Provider) ListComments(ctx context.Context, prNumber int) ([]*onlinegit.Comment, error) {
	comments, resp, err := p.client.ListIssueComments(p.owner, p.repo, int64(prNumber), gitea.ListIssueCommentOptions{})
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
	opts := gitea.CreateIssueCommentOption{
		Body: body,
	}

	comment, resp, err := p.client.CreateIssueComment(p.owner, p.repo, int64(prNumber), opts)
	if err != nil {
		return nil, p.wrapError("CreateComment", resp, err)
	}
	return p.toComment(comment), nil
}

func (p *Provider) UpdateComment(ctx context.Context, commentID int64, body string) (*onlinegit.Comment, error) {
	opts := gitea.EditIssueCommentOption{
		Body: body,
	}

	comment, resp, err := p.client.EditIssueComment(p.owner, p.repo, commentID, opts)
	if err != nil {
		return nil, p.wrapError("UpdateComment", resp, err)
	}
	return p.toComment(comment), nil
}

func (p *Provider) DeleteComment(ctx context.Context, commentID int64) error {
	resp, err := p.client.DeleteIssueComment(p.owner, p.repo, commentID)
	if err != nil {
		return p.wrapError("DeleteComment", resp, err)
	}
	return nil
}

func (p *Provider) GetCommit(ctx context.Context, sha string) (*onlinegit.Commit, error) {
	commit, resp, err := p.client.GetSingleCommit(p.owner, p.repo, sha)
	if err != nil {
		return nil, p.wrapError("GetCommit", resp, err)
	}
	return p.toRepoCommit(commit), nil
}

func (p *Provider) ListCommits(ctx context.Context, branch string, opts *onlinegit.ListOptions) ([]*onlinegit.Commit, error) {
	giteaOpts := gitea.ListCommitOptions{
		SHA: branch,
		ListOptions: gitea.ListOptions{
			Page:     opts.Page,
			PageSize: opts.PerPage,
		},
	}

	commits, resp, err := p.client.ListRepoCommits(p.owner, p.repo, giteaOpts)
	if err != nil {
		return nil, p.wrapError("ListCommits", resp, err)
	}

	result := make([]*onlinegit.Commit, len(commits))
	for i, c := range commits {
		result[i] = p.toRepoCommit(c)
	}
	return result, nil
}

// ==================== CI/CD Pipeline 管理 ====================
// Gitea 不直接支持 CI/CD Pipeline API，以下方法返回 ErrNotSupported

func (p *Provider) TriggerPipeline(ctx context.Context, opts *onlinegit.TriggerPipelineOptions) (*onlinegit.Pipeline, error) {
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "TriggerPipeline", onlinegit.ErrNotSupported, "Gitea does not support Pipeline API directly")
}

func (p *Provider) GetPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "GetPipeline", onlinegit.ErrNotSupported, "Gitea does not support Pipeline API directly")
}

func (p *Provider) ListPipelines(ctx context.Context, opts *onlinegit.ListPipelineOptions) ([]*onlinegit.Pipeline, error) {
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "ListPipelines", onlinegit.ErrNotSupported, "Gitea does not support Pipeline API directly")
}

func (p *Provider) CancelPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "CancelPipeline", onlinegit.ErrNotSupported, "Gitea does not support Pipeline API directly")
}

func (p *Provider) RetryPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "RetryPipeline", onlinegit.ErrNotSupported, "Gitea does not support Pipeline API directly")
}

func (p *Provider) ListPipelineJobs(ctx context.Context, pipelineID int64) ([]*onlinegit.PipelineJob, error) {
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "ListPipelineJobs", onlinegit.ErrNotSupported, "Gitea does not support Pipeline API directly")
}

// Helper methods

func (p *Provider) toPullRequest(pr *gitea.PullRequest) *onlinegit.PullRequest {
	state := onlinegit.PRStateOpen
	if pr.HasMerged {
		state = onlinegit.PRStateMerged
	} else if pr.State == gitea.StateClosed {
		state = onlinegit.PRStateClosed
	}

	result := &onlinegit.PullRequest{
		ID:           pr.ID,
		Number:       int(pr.Index),
		Title:        pr.Title,
		Body:         pr.Body,
		State:        state,
		SourceBranch: pr.Head.Ref,
		TargetBranch: pr.Base.Ref,
		URL:          pr.HTMLURL,
		Merged:       pr.HasMerged,
		CreatedAt:    *pr.Created,
		UpdatedAt:    *pr.Updated,
	}

	if pr.Poster != nil {
		result.Author = p.toUser(pr.Poster)
	}

	if pr.Merged != nil {
		result.MergedAt = *pr.Merged
	}

	if pr.Closed != nil {
		result.ClosedAt = *pr.Closed
	}

	if pr.MergedBy != nil {
		result.MergedBy = p.toUser(pr.MergedBy)
	}

	if len(pr.Labels) > 0 {
		result.Labels = make([]string, len(pr.Labels))
		for i, l := range pr.Labels {
			result.Labels[i] = l.Name
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

func (p *Provider) toUser(u *gitea.User) *onlinegit.User {
	return &onlinegit.User{
		ID:        u.ID,
		Login:     u.UserName,
		Name:      u.FullName,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
	}
}

func (p *Provider) toBranchCommit(c *gitea.PayloadCommit) *onlinegit.Commit {
	if c == nil {
		return nil
	}
	result := &onlinegit.Commit{
		SHA:     c.ID,
		Message: c.Message,
		URL:     c.URL,
	}

	if c.Author != nil {
		result.Author = &onlinegit.User{
			Name:  c.Author.Name,
			Email: c.Author.Email,
		}
	}

	if c.Committer != nil {
		result.Committer = &onlinegit.User{
			Name:  c.Committer.Name,
			Email: c.Committer.Email,
		}
	}

	return result
}

func (p *Provider) toRepoCommit(c *gitea.Commit) *onlinegit.Commit {
	if c == nil {
		return nil
	}

	result := &onlinegit.Commit{
		SHA:     c.SHA,
		Message: c.RepoCommit.Message,
		URL:     c.HTMLURL,
	}

	if c.Author != nil {
		result.Author = &onlinegit.User{
			ID:        c.Author.ID,
			Login:     c.Author.UserName,
			Name:      c.Author.FullName,
			Email:     c.Author.Email,
			AvatarURL: c.Author.AvatarURL,
		}
	}

	if c.Committer != nil {
		result.Committer = &onlinegit.User{
			ID:        c.Committer.ID,
			Login:     c.Committer.UserName,
			Name:      c.Committer.FullName,
			Email:     c.Committer.Email,
			AvatarURL: c.Committer.AvatarURL,
		}
	}

	if c.CommitMeta != nil {
		result.CreatedAt = c.CommitMeta.Created
	}

	return result
}

func (p *Provider) toComment(c *gitea.Comment) *onlinegit.Comment {
	result := &onlinegit.Comment{
		ID:        c.ID,
		Body:      c.Body,
		URL:       c.HTMLURL,
		CreatedAt: c.Created,
		UpdatedAt: c.Updated,
	}

	if c.Poster != nil {
		result.Author = p.toUser(c.Poster)
	}

	return result
}

func (p *Provider) wrapError(op string, resp *gitea.Response, err error) error {
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrNotFound, "")
		case http.StatusUnauthorized:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrUnauthorized, "")
		case http.StatusForbidden:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrForbidden, "")
		case http.StatusConflict:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrConflict, err.Error())
		case http.StatusTooManyRequests:
			return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, onlinegit.ErrRateLimit, "")
		}
	}
	return onlinegit.NewProviderError(onlinegit.PlatformGitea, op, err, "")
}
