package gitlab

import (
	"context"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// ==================== CI/CD Pipeline 管理 ====================

// TriggerPipeline 触发 Pipeline
func (p *Provider) TriggerPipeline(ctx context.Context, opts *onlinegit.TriggerPipelineOptions) (*onlinegit.Pipeline, error) {
	createOpts := &gitlab.CreatePipelineOptions{
		Ref: gitlab.Ptr(opts.Ref),
	}

	// 添加变量
	if len(opts.Variables) > 0 {
		vars := make([]*gitlab.PipelineVariableOptions, 0, len(opts.Variables))
		for key, value := range opts.Variables {
			vars = append(vars, &gitlab.PipelineVariableOptions{
				Key:   gitlab.Ptr(key),
				Value: gitlab.Ptr(value),
			})
		}
		createOpts.Variables = &vars
	}

	pipeline, resp, err := p.client.Pipelines.CreatePipeline(p.projectID, createOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("TriggerPipeline", resp, err)
	}

	return p.toPipeline(pipeline), nil
}

// GetPipeline 获取 Pipeline 详情
func (p *Provider) GetPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	pipeline, resp, err := p.client.Pipelines.GetPipeline(p.projectID, pipelineID, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("GetPipeline", resp, err)
	}
	return p.toPipeline(pipeline), nil
}

// ListPipelines 获取 Pipeline 列表
func (p *Provider) ListPipelines(ctx context.Context, opts *onlinegit.ListPipelineOptions) ([]*onlinegit.Pipeline, error) {
	listOpts := &gitlab.ListProjectPipelinesOptions{}

	if opts != nil {
		if opts.Ref != "" {
			listOpts.Ref = gitlab.Ptr(opts.Ref)
		}
		if opts.Status != "" {
			listOpts.Status = gitlab.Ptr(gitlab.BuildStateValue(string(opts.Status)))
		}
		if opts.Username != "" {
			listOpts.Username = gitlab.Ptr(opts.Username)
		}
		if opts.OrderBy != "" {
			listOpts.OrderBy = gitlab.Ptr(opts.OrderBy)
		}
		if opts.Sort != "" {
			listOpts.Sort = gitlab.Ptr(opts.Sort)
		}
		if opts.Page > 0 || opts.PerPage > 0 {
			listOpts.ListOptions = gitlab.ListOptions{
				Page:    int64(opts.Page),
				PerPage: int64(opts.PerPage),
			}
		}
	}

	pipelines, resp, err := p.client.Pipelines.ListProjectPipelines(p.projectID, listOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("ListPipelines", resp, err)
	}

	result := make([]*onlinegit.Pipeline, len(pipelines))
	for i, pl := range pipelines {
		result[i] = p.toPipelineBasic(pl)
	}
	return result, nil
}

// CancelPipeline 取消 Pipeline
func (p *Provider) CancelPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	pipeline, resp, err := p.client.Pipelines.CancelPipelineBuild(p.projectID, pipelineID, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("CancelPipeline", resp, err)
	}
	return p.toPipeline(pipeline), nil
}

// RetryPipeline 重试 Pipeline
func (p *Provider) RetryPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	pipeline, resp, err := p.client.Pipelines.RetryPipelineBuild(p.projectID, pipelineID, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("RetryPipeline", resp, err)
	}
	return p.toPipeline(pipeline), nil
}

// ListPipelineJobs 获取 Pipeline 作业列表
func (p *Provider) ListPipelineJobs(ctx context.Context, pipelineID int64) ([]*onlinegit.PipelineJob, error) {
	jobs, resp, err := p.client.Jobs.ListPipelineJobs(p.projectID, pipelineID, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("ListPipelineJobs", resp, err)
	}

	result := make([]*onlinegit.PipelineJob, len(jobs))
	for i, job := range jobs {
		result[i] = p.toPipelineJob(job)
	}
	return result, nil
}

// Helper methods

func (p *Provider) toPipeline(pl *gitlab.Pipeline) *onlinegit.Pipeline {
	result := &onlinegit.Pipeline{
		ID:             int64(pl.ID),
		IID:            pl.IID,
		ProjectID:      int64(pl.ProjectID),
		Status:         onlinegit.PipelineStatus(pl.Status),
		Source:         onlinegit.PipelineSource(pl.Source),
		Ref:            pl.Ref,
		SHA:            pl.SHA,
		WebURL:         pl.WebURL,
		Duration:       pl.Duration,
		QueuedDuration: pl.QueuedDuration,
	}

	if pl.CreatedAt != nil {
		result.CreatedAt = *pl.CreatedAt
	}
	if pl.UpdatedAt != nil {
		result.UpdatedAt = *pl.UpdatedAt
	}
	if pl.StartedAt != nil {
		result.StartedAt = pl.StartedAt
	}
	if pl.FinishedAt != nil {
		result.FinishedAt = pl.FinishedAt
	}
	if pl.User != nil {
		result.User = &onlinegit.User{
			ID:        int64(pl.User.ID),
			Login:     pl.User.Username,
			Name:      pl.User.Name,
			AvatarURL: pl.User.AvatarURL,
		}
	}

	return result
}

func (p *Provider) toPipelineBasic(pl *gitlab.PipelineInfo) *onlinegit.Pipeline {
	result := &onlinegit.Pipeline{
		ID:        int64(pl.ID),
		IID:       pl.IID,
		ProjectID: int64(pl.ProjectID),
		Status:    onlinegit.PipelineStatus(pl.Status),
		Source:    onlinegit.PipelineSource(pl.Source),
		Ref:       pl.Ref,
		SHA:       pl.SHA,
		WebURL:    pl.WebURL,
	}
	if pl.CreatedAt != nil {
		result.CreatedAt = *pl.CreatedAt
	}
	if pl.UpdatedAt != nil {
		result.UpdatedAt = *pl.UpdatedAt
	}
	return result
}

func (p *Provider) toPipelineJob(job *gitlab.Job) *onlinegit.PipelineJob {
	result := &onlinegit.PipelineJob{
		ID:           int64(job.ID),
		Name:         job.Name,
		Stage:        job.Stage,
		Status:       onlinegit.PipelineStatus(job.Status),
		Ref:          job.Ref,
		WebURL:       job.WebURL,
		Duration:     job.Duration,
		AllowFailure: job.AllowFailure,
	}
	if job.CreatedAt != nil {
		result.CreatedAt = *job.CreatedAt
	}
	if job.StartedAt != nil {
		result.StartedAt = job.StartedAt
	}
	if job.FinishedAt != nil {
		result.FinishedAt = job.FinishedAt
	}
	if job.User != nil {
		result.User = &onlinegit.User{
			ID:        int64(job.User.ID),
			Login:     job.User.Username,
			Name:      job.User.Name,
			AvatarURL: job.User.AvatarURL,
		}
	}
	return result
}
