package gitea

import (
	"context"
	"fmt"

	"code.gitea.io/sdk/gitea"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// ListBranches 获取分支列表
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

// GetBranch 获取指定分支
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

// CreateBranch 创建分支
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

// DeleteBranch 删除分支
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

// SetBranchProtection 设置分支保护
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

// UnsetBranchProtection 取消分支保护
func (p *Provider) UnsetBranchProtection(ctx context.Context, name string) error {
	resp, err := p.client.DeleteBranchProtection(p.owner, p.repo, name)
	if err != nil {
		return p.wrapError("UnsetBranchProtection", resp, err)
	}
	return nil
}

// CompareBranches 比较两个分支
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
