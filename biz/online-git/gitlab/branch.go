package gitlab

import (
	"context"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

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
