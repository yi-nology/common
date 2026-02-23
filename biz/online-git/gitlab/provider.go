package gitlab

import (
	"context"
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

func (p *Provider) GetRepository(ctx context.Context) (*onlinegit.Repository, error) {
	project, resp, err := p.client.Projects.GetProject(p.projectID, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("GetRepository", resp, err)
	}

	return &onlinegit.Repository{
		ID:            int64(project.ID),
		Name:          project.Name,
		FullName:      project.PathWithNamespace,
		Description:   project.Description,
		URL:           project.WebURL,
		CloneURL:      project.HTTPURLToRepo,
		DefaultBranch: project.DefaultBranch,
		Private:       project.Visibility != gitlab.PublicVisibility,
		Fork:          project.ForkedFromProject != nil,
		CreatedAt:     *project.CreatedAt,
	}, nil
}

func (p *Provider) ListBranches(ctx context.Context, opts *onlinegit.ListOptions) ([]*onlinegit.Branch, error) {
	glOpts := &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    int64(opts.Page),
			PerPage: int64(opts.PerPage),
		},
	}

	branches, resp, err := p.client.Branches.ListBranches(p.projectID, glOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("ListBranches", resp, err)
	}

	result := make([]*onlinegit.Branch, len(branches))
	for i, b := range branches {
		result[i] = &onlinegit.Branch{
			Name:      b.Name,
			CommitSHA: b.Commit.ID,
			Protected: b.Protected,
			Default:   b.Default,
		}
	}
	return result, nil
}

func (p *Provider) GetBranch(ctx context.Context, name string) (*onlinegit.Branch, error) {
	branch, resp, err := p.client.Branches.GetBranch(p.projectID, name, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("GetBranch", resp, err)
	}

	return &onlinegit.Branch{
		Name:      branch.Name,
		CommitSHA: branch.Commit.ID,
		Protected: branch.Protected,
		Default:   branch.Default,
		Commit:    p.toCommitFromBranch(branch.Commit),
	}, nil
}

func (p *Provider) CreateBranch(ctx context.Context, name, sourceBranch string) (*onlinegit.Branch, error) {
	opts := &gitlab.CreateBranchOptions{
		Branch: gitlab.Ptr(name),
		Ref:    gitlab.Ptr(sourceBranch),
	}

	branch, resp, err := p.client.Branches.CreateBranch(p.projectID, opts, gitlab.WithContext(ctx))
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
	resp, err := p.client.Branches.DeleteBranch(p.projectID, name, gitlab.WithContext(ctx))
	if err != nil {
		return p.wrapError("DeleteBranch", resp, err)
	}
	return nil
}

func (p *Provider) SetBranchProtection(ctx context.Context, name string, rules *onlinegit.ProtectionRules) error {
	// 使用 ProtectedBranches API
	opts := &gitlab.ProtectRepositoryBranchesOptions{
		Name: gitlab.Ptr(name),
	}

	if rules.AllowForcePush {
		opts.AllowForcePush = gitlab.Ptr(true)
	}

	_, resp, err := p.client.ProtectedBranches.ProtectRepositoryBranches(p.projectID, opts, gitlab.WithContext(ctx))
	if err != nil {
		return p.wrapError("SetBranchProtection", resp, err)
	}
	return nil
}

func (p *Provider) UnsetBranchProtection(ctx context.Context, name string) error {
	resp, err := p.client.ProtectedBranches.UnprotectRepositoryBranches(p.projectID, name, gitlab.WithContext(ctx))
	if err != nil {
		return p.wrapError("UnsetBranchProtection", resp, err)
	}
	return nil
}

func (p *Provider) ListPullRequests(ctx context.Context, state onlinegit.PRState, opts *onlinegit.ListOptions) ([]*onlinegit.PullRequest, error) {
	glState := ""
	switch state {
	case onlinegit.PRStateOpen:
		glState = "opened"
	case onlinegit.PRStateClosed:
		glState = "closed"
	case onlinegit.PRStateMerged:
		glState = "merged"
	case onlinegit.PRStateAll:
		glState = "all"
	}

	glOpts := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.Ptr(glState),
		ListOptions: gitlab.ListOptions{
			Page:    int64(opts.Page),
			PerPage: int64(opts.PerPage),
		},
	}

	mrs, resp, err := p.client.MergeRequests.ListProjectMergeRequests(p.projectID, glOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("ListPullRequests", resp, err)
	}

	result := make([]*onlinegit.PullRequest, len(mrs))
	for i, mr := range mrs {
		result[i] = p.toBasicMergeRequest(mr)
	}
	return result, nil
}

func (p *Provider) GetPullRequest(ctx context.Context, number int) (*onlinegit.PullRequest, error) {
	mr, resp, err := p.client.MergeRequests.GetMergeRequest(p.projectID, int64(number), nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("GetPullRequest", resp, err)
	}
	return p.toMergeRequest(mr), nil
}

func (p *Provider) CreatePullRequest(ctx context.Context, req *onlinegit.CreatePRRequest) (*onlinegit.PullRequest, error) {
	opts := &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.Ptr(req.Title),
		Description:  gitlab.Ptr(req.Body),
		SourceBranch: gitlab.Ptr(req.SourceBranch),
		TargetBranch: gitlab.Ptr(req.TargetBranch),
	}

	if len(req.Labels) > 0 {
		opts.Labels = gitlab.Ptr(gitlab.LabelOptions(req.Labels))
	}

	mr, resp, err := p.client.MergeRequests.CreateMergeRequest(p.projectID, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("CreatePullRequest", resp, err)
	}

	return p.toMergeRequest(mr), nil
}

func (p *Provider) UpdatePullRequest(ctx context.Context, number int, title, body string) (*onlinegit.PullRequest, error) {
	opts := &gitlab.UpdateMergeRequestOptions{
		Title:       gitlab.Ptr(title),
		Description: gitlab.Ptr(body),
	}

	mr, resp, err := p.client.MergeRequests.UpdateMergeRequest(p.projectID, int64(number), opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("UpdatePullRequest", resp, err)
	}
	return p.toMergeRequest(mr), nil
}

func (p *Provider) MergePullRequest(ctx context.Context, number int, opts *onlinegit.MergeOptions) error {
	acceptOpts := &gitlab.AcceptMergeRequestOptions{}

	if opts != nil {
		if opts.CommitMessage != "" {
			acceptOpts.MergeCommitMessage = gitlab.Ptr(opts.CommitMessage)
		}
		if opts.Method == onlinegit.MergeMethodSquash {
			acceptOpts.Squash = gitlab.Ptr(true)
		}
		if opts.DeleteBranch {
			acceptOpts.ShouldRemoveSourceBranch = gitlab.Ptr(true)
		}
	}

	_, resp, err := p.client.MergeRequests.AcceptMergeRequest(p.projectID, int64(number), acceptOpts, gitlab.WithContext(ctx))
	if err != nil {
		return p.wrapError("MergePullRequest", resp, err)
	}
	return nil
}

func (p *Provider) ClosePullRequest(ctx context.Context, number int) error {
	opts := &gitlab.UpdateMergeRequestOptions{
		StateEvent: gitlab.Ptr("close"),
	}

	_, resp, err := p.client.MergeRequests.UpdateMergeRequest(p.projectID, int64(number), opts, gitlab.WithContext(ctx))
	if err != nil {
		return p.wrapError("ClosePullRequest", resp, err)
	}
	return nil
}

func (p *Provider) GetPullRequestCommits(ctx context.Context, number int) ([]*onlinegit.Commit, error) {
	commits, resp, err := p.client.MergeRequests.GetMergeRequestCommits(p.projectID, int64(number), nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("GetPullRequestCommits", resp, err)
	}

	result := make([]*onlinegit.Commit, len(commits))
	for i, c := range commits {
		result[i] = &onlinegit.Commit{
			SHA:     c.ID,
			Message: c.Message,
			Author: &onlinegit.User{
				Name:  c.AuthorName,
				Email: c.AuthorEmail,
			},
			CreatedAt: *c.CreatedAt,
		}
	}
	return result, nil
}

func (p *Provider) CompareBranches(ctx context.Context, base, head string) (*onlinegit.CompareResult, error) {
	opts := &gitlab.CompareOptions{
		From: gitlab.Ptr(base),
		To:   gitlab.Ptr(head),
	}

	compare, resp, err := p.client.Repositories.Compare(p.projectID, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("CompareBranches", resp, err)
	}

	commits := make([]*onlinegit.Commit, len(compare.Commits))
	for i, c := range compare.Commits {
		commits[i] = &onlinegit.Commit{
			SHA:     c.ID,
			Message: c.Message,
			Author: &onlinegit.User{
				Name:  c.AuthorName,
				Email: c.AuthorEmail,
			},
			CreatedAt: *c.CreatedAt,
		}
	}

	files := make([]*onlinegit.FileChange, len(compare.Diffs))
	for i, d := range compare.Diffs {
		status := onlinegit.FileChangeModified
		if d.NewFile {
			status = onlinegit.FileChangeAdded
		} else if d.DeletedFile {
			status = onlinegit.FileChangeDeleted
		} else if d.RenamedFile {
			status = onlinegit.FileChangeRenamed
		}

		files[i] = &onlinegit.FileChange{
			Filename:     d.NewPath,
			Status:       status,
			PreviousName: d.OldPath,
			Patch:        d.Diff,
		}
	}

	return &onlinegit.CompareResult{
		BaseBranch:   base,
		HeadBranch:   head,
		TotalCommits: len(commits),
		Commits:      commits,
		Files:        files,
		DiffStat: &onlinegit.DiffStat{
			ChangedFiles: len(files),
		},
	}, nil
}

func (p *Provider) ListComments(ctx context.Context, prNumber int) ([]*onlinegit.Comment, error) {
	notes, resp, err := p.client.Notes.ListMergeRequestNotes(p.projectID, int64(prNumber), nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("ListComments", resp, err)
	}

	result := make([]*onlinegit.Comment, len(notes))
	for i, n := range notes {
		result[i] = p.toNote(n)
	}
	return result, nil
}

func (p *Provider) CreateComment(ctx context.Context, prNumber int, body string) (*onlinegit.Comment, error) {
	opts := &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.Ptr(body),
	}

	note, resp, err := p.client.Notes.CreateMergeRequestNote(p.projectID, int64(prNumber), opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("CreateComment", resp, err)
	}
	return p.toNote(note), nil
}

func (p *Provider) UpdateComment(ctx context.Context, commentID int64, body string) (*onlinegit.Comment, error) {
	// GitLab Notes API 需要 MR 号，这里需要额外处理
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitLab, "UpdateComment", fmt.Errorf("not implemented: need MR number"), "")
}

func (p *Provider) DeleteComment(ctx context.Context, commentID int64) error {
	// GitLab Notes API 需要 MR 号
	return onlinegit.NewProviderError(onlinegit.PlatformGitLab, "DeleteComment", fmt.Errorf("not implemented: need MR number"), "")
}

func (p *Provider) GetCommit(ctx context.Context, sha string) (*onlinegit.Commit, error) {
	commit, resp, err := p.client.Commits.GetCommit(p.projectID, sha, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("GetCommit", resp, err)
	}

	return &onlinegit.Commit{
		SHA:     commit.ID,
		Message: commit.Message,
		Author: &onlinegit.User{
			Name:  commit.AuthorName,
			Email: commit.AuthorEmail,
		},
		Committer: &onlinegit.User{
			Name:  commit.CommitterName,
			Email: commit.CommitterEmail,
		},
		URL:       commit.WebURL,
		CreatedAt: *commit.CreatedAt,
	}, nil
}

func (p *Provider) ListCommits(ctx context.Context, branch string, opts *onlinegit.ListOptions) ([]*onlinegit.Commit, error) {
	glOpts := &gitlab.ListCommitsOptions{
		RefName: gitlab.Ptr(branch),
		ListOptions: gitlab.ListOptions{
			Page:    int64(opts.Page),
			PerPage: int64(opts.PerPage),
		},
	}

	commits, resp, err := p.client.Commits.ListCommits(p.projectID, glOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("ListCommits", resp, err)
	}

	result := make([]*onlinegit.Commit, len(commits))
	for i, c := range commits {
		result[i] = &onlinegit.Commit{
			SHA:     c.ID,
			Message: c.Message,
			Author: &onlinegit.User{
				Name:  c.AuthorName,
				Email: c.AuthorEmail,
			},
			URL:       c.WebURL,
			CreatedAt: *c.CreatedAt,
		}
	}
	return result, nil
}

// Helper methods

func (p *Provider) toMergeRequest(mr *gitlab.MergeRequest) *onlinegit.PullRequest {
	state := onlinegit.PRStateOpen
	switch mr.State {
	case "merged":
		state = onlinegit.PRStateMerged
	case "closed":
		state = onlinegit.PRStateClosed
	}

	result := &onlinegit.PullRequest{
		ID:           mr.IID,
		Number:       int(mr.IID),
		Title:        mr.Title,
		Body:         mr.Description,
		State:        state,
		SourceBranch: mr.SourceBranch,
		TargetBranch: mr.TargetBranch,
		URL:          mr.WebURL,
		Merged:       mr.State == "merged",
		CreatedAt:    *mr.CreatedAt,
		UpdatedAt:    *mr.UpdatedAt,
	}

	if mr.Author != nil {
		result.Author = &onlinegit.User{
			ID:        int64(mr.Author.ID),
			Login:     mr.Author.Username,
			Name:      mr.Author.Name,
			AvatarURL: mr.Author.AvatarURL,
		}
	}

	if mr.MergedAt != nil {
		result.MergedAt = *mr.MergedAt
	}

	if mr.ClosedAt != nil {
		result.ClosedAt = *mr.ClosedAt
	}

	if mr.MergedBy != nil {
		result.MergedBy = &onlinegit.User{
			ID:        int64(mr.MergedBy.ID),
			Login:     mr.MergedBy.Username,
			Name:      mr.MergedBy.Name,
			AvatarURL: mr.MergedBy.AvatarURL,
		}
	}

	if len(mr.Labels) > 0 {
		result.Labels = mr.Labels
	}

	return result
}

func (p *Provider) toBasicMergeRequest(mr *gitlab.BasicMergeRequest) *onlinegit.PullRequest {
	state := onlinegit.PRStateOpen
	switch mr.State {
	case "merged":
		state = onlinegit.PRStateMerged
	case "closed":
		state = onlinegit.PRStateClosed
	}

	result := &onlinegit.PullRequest{
		ID:           mr.IID,
		Number:       int(mr.IID),
		Title:        mr.Title,
		Body:         mr.Description,
		State:        state,
		SourceBranch: mr.SourceBranch,
		TargetBranch: mr.TargetBranch,
		URL:          mr.WebURL,
		Merged:       mr.State == "merged",
		CreatedAt:    *mr.CreatedAt,
		UpdatedAt:    *mr.UpdatedAt,
	}

	if mr.Author != nil {
		result.Author = &onlinegit.User{
			ID:        int64(mr.Author.ID),
			Login:     mr.Author.Username,
			Name:      mr.Author.Name,
			AvatarURL: mr.Author.AvatarURL,
		}
	}

	if mr.MergedAt != nil {
		result.MergedAt = *mr.MergedAt
	}

	if mr.ClosedAt != nil {
		result.ClosedAt = *mr.ClosedAt
	}

	if mr.MergedBy != nil {
		result.MergedBy = &onlinegit.User{
			ID:        int64(mr.MergedBy.ID),
			Login:     mr.MergedBy.Username,
			Name:      mr.MergedBy.Name,
			AvatarURL: mr.MergedBy.AvatarURL,
		}
	}

	if len(mr.Labels) > 0 {
		result.Labels = mr.Labels
	}

	return result
}

func (p *Provider) toCommitFromBranch(c *gitlab.Commit) *onlinegit.Commit {
	if c == nil {
		return nil
	}
	return &onlinegit.Commit{
		SHA:     c.ID,
		Message: c.Message,
		Author: &onlinegit.User{
			Name:  c.AuthorName,
			Email: c.AuthorEmail,
		},
		URL:       c.WebURL,
		CreatedAt: *c.CreatedAt,
	}
}

func (p *Provider) toNote(n *gitlab.Note) *onlinegit.Comment {
	result := &onlinegit.Comment{
		ID:        int64(n.ID),
		Body:      n.Body,
		CreatedAt: *n.CreatedAt,
		UpdatedAt: *n.UpdatedAt,
	}

	// Author 是值类型，检查 ID 是否非零
	if n.Author.ID != 0 {
		result.Author = &onlinegit.User{
			ID:        int64(n.Author.ID),
			Login:     n.Author.Username,
			Name:      n.Author.Name,
			AvatarURL: n.Author.AvatarURL,
		}
	}

	return result
}

func (p *Provider) wrapError(op string, resp *gitlab.Response, err error) error {
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrNotFound, "")
		case http.StatusUnauthorized:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrUnauthorized, "")
		case http.StatusForbidden:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrForbidden, "")
		case http.StatusConflict:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrConflict, err.Error())
		case http.StatusTooManyRequests:
			return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, onlinegit.ErrRateLimit, "")
		}
	}
	return onlinegit.NewProviderError(onlinegit.PlatformGitLab, op, err, "")
}
