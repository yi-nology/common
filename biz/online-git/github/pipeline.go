package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v56/github"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// TriggerPipeline 触发 GitHub Actions Workflow
// 通过 workflow_dispatch 事件触发指定的 workflow
func (p *Provider) TriggerPipeline(ctx context.Context, opts *onlinegit.TriggerPipelineOptions) (*onlinegit.Pipeline, error) {
	if opts == nil || opts.Ref == "" {
		return nil, onlinegit.NewProviderError(onlinegit.PlatformGitHub, "TriggerPipeline", fmt.Errorf("ref is required"), "")
	}

	// GitHub Actions 需要指定 workflow 文件名，默认使用 ci.yml
	workflowFile := "ci.yml"
	if opts.Variables != nil {
		if wf, ok := opts.Variables["workflow_file"]; ok {
			workflowFile = wf
			delete(opts.Variables, "workflow_file")
		}
	}

	// 构建 inputs 参数
	inputs := make(map[string]interface{})
	for k, v := range opts.Variables {
		inputs[k] = v
	}

	event := github.CreateWorkflowDispatchEventRequest{
		Ref:    opts.Ref,
		Inputs: inputs,
	}

	resp, err := p.client.Actions.CreateWorkflowDispatchEventByFileName(ctx, p.owner, p.repo, workflowFile, event)
	if err != nil {
		return nil, p.wrapError("TriggerPipeline", resp, err)
	}

	// GitHub Actions 的 dispatch 不直接返回 run，需要查询最新的 run
	// 等待一小段时间后查询
	listOpts := &github.ListWorkflowRunsOptions{
		Branch: opts.Ref,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	}

	runs, _, err := p.client.Actions.ListRepositoryWorkflowRuns(ctx, p.owner, p.repo, listOpts)
	if err != nil || len(runs.WorkflowRuns) == 0 {
		// 返回一个占位的 Pipeline 对象，表示触发成功
		return &onlinegit.Pipeline{
			Ref:    opts.Ref,
			Status: onlinegit.PipelineStatusPending,
			Source: onlinegit.PipelineSourceAPI,
		}, nil
	}

	return p.toWorkflowRunPipeline(runs.WorkflowRuns[0]), nil
}

// GetPipeline 获取 GitHub Actions Workflow Run 详情
func (p *Provider) GetPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	run, resp, err := p.client.Actions.GetWorkflowRunByID(ctx, p.owner, p.repo, pipelineID)
	if err != nil {
		return nil, p.wrapError("GetPipeline", resp, err)
	}
	return p.toWorkflowRunPipeline(run), nil
}

// ListPipelines 获取 GitHub Actions Workflow Run 列表
func (p *Provider) ListPipelines(ctx context.Context, opts *onlinegit.ListPipelineOptions) ([]*onlinegit.Pipeline, error) {
	listOpts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			Page:    opts.Page,
			PerPage: opts.PerPage,
		},
	}

	if opts.Ref != "" {
		listOpts.Branch = opts.Ref
	}
	if opts.Status != "" {
		listOpts.Status = string(opts.Status)
	}
	if opts.Username != "" {
		listOpts.Actor = opts.Username
	}

	runs, resp, err := p.client.Actions.ListRepositoryWorkflowRuns(ctx, p.owner, p.repo, listOpts)
	if err != nil {
		return nil, p.wrapError("ListPipelines", resp, err)
	}

	result := make([]*onlinegit.Pipeline, len(runs.WorkflowRuns))
	for i, run := range runs.WorkflowRuns {
		result[i] = p.toWorkflowRunPipeline(run)
	}
	return result, nil
}

// CancelPipeline 取消 GitHub Actions Workflow Run
func (p *Provider) CancelPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	resp, err := p.client.Actions.CancelWorkflowRunByID(ctx, p.owner, p.repo, pipelineID)
	if err != nil {
		return nil, p.wrapError("CancelPipeline", resp, err)
	}

	// 获取更新后的状态
	return p.GetPipeline(ctx, pipelineID)
}

// RetryPipeline 重试 GitHub Actions Workflow Run
func (p *Provider) RetryPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	resp, err := p.client.Actions.RerunWorkflowByID(ctx, p.owner, p.repo, pipelineID)
	if err != nil {
		return nil, p.wrapError("RetryPipeline", resp, err)
	}

	// 获取更新后的状态
	return p.GetPipeline(ctx, pipelineID)
}

// ListPipelineJobs 获取 GitHub Actions Workflow Run 的作业列表
func (p *Provider) ListPipelineJobs(ctx context.Context, pipelineID int64) ([]*onlinegit.PipelineJob, error) {
	jobs, resp, err := p.client.Actions.ListWorkflowJobs(ctx, p.owner, p.repo, pipelineID, nil)
	if err != nil {
		return nil, p.wrapError("ListPipelineJobs", resp, err)
	}

	result := make([]*onlinegit.PipelineJob, len(jobs.Jobs))
	for i, job := range jobs.Jobs {
		result[i] = p.toWorkflowJob(job)
	}
	return result, nil
}

// toWorkflowRunPipeline 将 GitHub WorkflowRun 转换为统一的 Pipeline 结构
func (p *Provider) toWorkflowRunPipeline(run *github.WorkflowRun) *onlinegit.Pipeline {
	pipeline := &onlinegit.Pipeline{
		ID:     run.GetID(),
		IID:    int64(run.GetRunNumber()),
		Ref:    run.GetHeadBranch(),
		SHA:    run.GetHeadSHA(),
		WebURL: run.GetHTMLURL(),
		Status: p.mapWorkflowRunStatus(run.GetStatus(), run.GetConclusion()),
		Source: p.mapWorkflowRunEvent(run.GetEvent()),
	}

	if run.CreatedAt != nil {
		pipeline.CreatedAt = run.CreatedAt.Time
	}
	if run.UpdatedAt != nil {
		pipeline.UpdatedAt = run.UpdatedAt.Time
	}
	if run.RunStartedAt != nil {
		startedAt := run.RunStartedAt.Time
		pipeline.StartedAt = &startedAt
	}

	if run.Actor != nil {
		pipeline.User = &onlinegit.User{
			ID:        run.Actor.GetID(),
			Login:     run.Actor.GetLogin(),
			Name:      run.Actor.GetName(),
			AvatarURL: run.Actor.GetAvatarURL(),
		}
	}

	return pipeline
}

// mapWorkflowRunStatus 映射 GitHub Workflow Run 状态到统一状态
func (p *Provider) mapWorkflowRunStatus(status, conclusion string) onlinegit.PipelineStatus {
	switch status {
	case "queued":
		return onlinegit.PipelineStatusPending
	case "in_progress":
		return onlinegit.PipelineStatusRunning
	case "completed":
		switch conclusion {
		case "success":
			return onlinegit.PipelineStatusSuccess
		case "failure":
			return onlinegit.PipelineStatusFailed
		case "cancelled":
			return onlinegit.PipelineStatusCanceled
		case "skipped":
			return onlinegit.PipelineStatusSkipped
		default:
			return onlinegit.PipelineStatusFailed
		}
	case "waiting":
		return onlinegit.PipelineStatusWaitingForResource
	default:
		return onlinegit.PipelineStatusCreated
	}
}

// mapWorkflowRunEvent 映射 GitHub Workflow Run 事件到 Pipeline Source
func (p *Provider) mapWorkflowRunEvent(event string) onlinegit.PipelineSource {
	switch event {
	case "push":
		return onlinegit.PipelineSourcePush
	case "workflow_dispatch":
		return onlinegit.PipelineSourceWeb
	case "schedule":
		return onlinegit.PipelineSourceSchedule
	case "pull_request", "pull_request_target":
		return onlinegit.PipelineSourceMergeRequestEvent
	case "repository_dispatch":
		return onlinegit.PipelineSourceAPI
	default:
		return onlinegit.PipelineSourcePush
	}
}

// toWorkflowJob 将 GitHub WorkflowJob 转换为统一的 PipelineJob 结构
func (p *Provider) toWorkflowJob(job *github.WorkflowJob) *onlinegit.PipelineJob {
	pipelineJob := &onlinegit.PipelineJob{
		ID:     job.GetID(),
		Name:   job.GetName(),
		Status: p.mapWorkflowJobStatus(job.GetStatus(), job.GetConclusion()),
		WebURL: job.GetHTMLURL(),
	}

	if job.StartedAt != nil {
		startedAt := job.StartedAt.Time
		pipelineJob.StartedAt = &startedAt
	}
	if job.CompletedAt != nil {
		completedAt := job.CompletedAt.Time
		pipelineJob.FinishedAt = &completedAt
		if job.StartedAt != nil {
			pipelineJob.Duration = job.CompletedAt.Sub(job.StartedAt.Time).Seconds()
		}
	}

	return pipelineJob
}

// mapWorkflowJobStatus 映射 GitHub Workflow Job 状态
func (p *Provider) mapWorkflowJobStatus(status, conclusion string) onlinegit.PipelineStatus {
	switch status {
	case "queued":
		return onlinegit.PipelineStatusPending
	case "in_progress":
		return onlinegit.PipelineStatusRunning
	case "completed":
		switch conclusion {
		case "success":
			return onlinegit.PipelineStatusSuccess
		case "failure":
			return onlinegit.PipelineStatusFailed
		case "cancelled":
			return onlinegit.PipelineStatusCanceled
		case "skipped":
			return onlinegit.PipelineStatusSkipped
		default:
			return onlinegit.PipelineStatusFailed
		}
	default:
		return onlinegit.PipelineStatusCreated
	}
}
