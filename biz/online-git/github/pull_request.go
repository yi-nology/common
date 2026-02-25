package github

import (
	"context"

	"github.com/google/go-github/v56/github"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// ListPullRequests 获取 Pull Request 列表
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

// GetPullRequest 获取指定的 Pull Request
func (p *Provider) GetPullRequest(ctx context.Context, number int) (*onlinegit.PullRequest, error) {
	pr, resp, err := p.client.PullRequests.Get(ctx, p.owner, p.repo, number)
	if err != nil {
		return nil, p.wrapError("GetPullRequest", resp, err)
	}
	return p.toPullRequest(pr), nil
}

// CreatePullRequest 创建 Pull Request
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

// UpdatePullRequest 更新 Pull Request
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

// MergePullRequest 合并 Pull Request
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

// ClosePullRequest 关闭 Pull Request
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

// GetPullRequestCommits 获取 Pull Request 的提交列表
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

// ListComments 获取 Pull Request 的评论列表
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

// CreateComment 创建评论
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

// UpdateComment 更新评论
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

// DeleteComment 删除评论
func (p *Provider) DeleteComment(ctx context.Context, commentID int64) error {
	resp, err := p.client.Issues.DeleteComment(ctx, p.owner, p.repo, commentID)
	if err != nil {
		return p.wrapError("DeleteComment", resp, err)
	}
	return nil
}
