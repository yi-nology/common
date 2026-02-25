package github

import (
	"context"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// GetRepository 获取仓库信息
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
