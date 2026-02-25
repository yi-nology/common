package github

import (
	"context"

	"github.com/google/go-github/v56/github"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// ListBranches 获取分支列表
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

// GetBranch 获取指定分支
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

// CreateBranch 创建分支
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

// DeleteBranch 删除分支
func (p *Provider) DeleteBranch(ctx context.Context, name string) error {
	ref := "heads/" + name
	resp, err := p.client.Git.DeleteRef(ctx, p.owner, p.repo, ref)
	if err != nil {
		return p.wrapError("DeleteBranch", resp, err)
	}
	return nil
}

// SetBranchProtection 设置分支保护
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

// UnsetBranchProtection 取消分支保护
func (p *Provider) UnsetBranchProtection(ctx context.Context, name string) error {
	resp, err := p.client.Repositories.RemoveBranchProtection(ctx, p.owner, p.repo, name)
	if err != nil {
		return p.wrapError("UnsetBranchProtection", resp, err)
	}
	return nil
}

// CompareBranches 比较两个分支
func (p *Provider) CompareBranches(ctx context.Context, base, head string) (*onlinegit.CompareResult, error) {
	comparison, resp, err := p.client.Repositories.CompareCommits(ctx, p.owner, p.repo, base, head, nil)
	if err != nil {
		return nil, p.wrapError("CompareBranches", resp, err)
	}

	commits := make([]*onlinegit.Commit, len(comparison.Commits))
	for i, c := range comparison.Commits {
		commits[i] = p.toCommit(c.Commit, c.GetSHA())
	}

	return &onlinegit.CompareResult{
		BaseBranch:   base,
		HeadBranch:   head,
		TotalCommits: comparison.GetTotalCommits(),
		Commits:      commits,
		DiffStat: &onlinegit.DiffStat{
			ChangedFiles: len(comparison.Files),
			Additions:    comparison.GetAheadBy(),
			Deletions:    comparison.GetBehindBy(),
		},
	}, nil
}
