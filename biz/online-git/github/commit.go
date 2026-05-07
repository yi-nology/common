package github

import (
	"context"

	"github.com/google/go-github/v56/github"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// GetCommit 获取指定的提交
func (p *Provider) GetCommit(ctx context.Context, sha string) (*onlinegit.Commit, error) {
	commit, resp, err := p.client.Repositories.GetCommit(ctx, p.owner, p.repo, sha, nil)
	if err != nil {
		return nil, p.wrapError("GetCommit", resp, err)
	}
	return p.toCommit(commit.Commit, commit.GetSHA()), nil
}

// ListCommits 获取提交列表
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
