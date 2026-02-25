package gitlab

import (
	"context"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// GetRepository 获取仓库信息
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
