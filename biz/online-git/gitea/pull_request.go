package gitea

import (
	"context"

	"code.gitea.io/sdk/gitea"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// ListPullRequests 获取 Pull Request 列表
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

// GetPullRequest 获取指定的 Pull Request
func (p *Provider) GetPullRequest(ctx context.Context, number int) (*onlinegit.PullRequest, error) {
	pr, resp, err := p.client.GetPullRequest(p.owner, p.repo, int64(number))
	if err != nil {
		return nil, p.wrapError("GetPullRequest", resp, err)
	}
	return p.toPullRequest(pr), nil
}

// CreatePullRequest 创建 Pull Request
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

// UpdatePullRequest 更新 Pull Request
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

// MergePullRequest 合并 Pull Request
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

// ClosePullRequest 关闭 Pull Request
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

// GetPullRequestCommits 获取 Pull Request 的提交列表
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

// ListComments 获取 Pull Request 的评论列表
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

// CreateComment 创建评论
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

// UpdateComment 更新评论
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

// DeleteComment 删除评论
func (p *Provider) DeleteComment(ctx context.Context, commentID int64) error {
	resp, err := p.client.DeleteIssueComment(p.owner, p.repo, commentID)
	if err != nil {
		return p.wrapError("DeleteComment", resp, err)
	}
	return nil
}
