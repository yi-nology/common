package gitea

import (
	"context"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// GetRepository 获取仓库信息
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
