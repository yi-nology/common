package gitea

import (
	"context"

	"code.gitea.io/sdk/gitea"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// GetCommit 获取指定的提交
func (p *Provider) GetCommit(ctx context.Context, sha string) (*onlinegit.Commit, error) {
	commit, resp, err := p.client.GetSingleCommit(p.owner, p.repo, sha)
	if err != nil {
		return nil, p.wrapError("GetCommit", resp, err)
	}
	return p.toRepoCommit(commit), nil
}

// ListCommits 获取提交列表
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
