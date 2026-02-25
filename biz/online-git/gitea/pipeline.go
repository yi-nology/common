package gitea

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"code.gitea.io/sdk/gitea"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// workflowDispatchRequest workflow dispatch 请求体
type workflowDispatchRequest struct {
	Ref    string            `json:"ref"`
	Inputs map[string]string `json:"inputs,omitempty"`
}

// TriggerPipeline 触发 Gitea Actions Workflow
// 使用 Gitea REST API: POST /api/v1/repos/{owner}/{repo}/actions/workflows/{workflow}/dispatches
// 需要 Gitea 1.23+ 版本支持
func (p *Provider) TriggerPipeline(ctx context.Context, opts *onlinegit.TriggerPipelineOptions) (*onlinegit.Pipeline, error) {
	if opts == nil || opts.Ref == "" {
		return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "TriggerPipeline", fmt.Errorf("ref is required"), "")
	}

	// 获取 workflow 文件名，默认使用 ci.yml
	workflowFile := "ci.yml"
	if opts.Variables != nil {
		if wf, ok := opts.Variables["workflow_file"]; ok {
			workflowFile = wf
			delete(opts.Variables, "workflow_file")
		}
	}

	// 构建请求体
	reqBody := workflowDispatchRequest{
		Ref:    opts.Ref,
		Inputs: opts.Variables,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "TriggerPipeline", err, "failed to marshal request body")
	}

	// 构建 API URL
	apiURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/actions/workflows/%s/dispatches",
		p.baseURL, p.owner, p.repo, workflowFile)

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "TriggerPipeline", err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+p.token)

	// 发送请求
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "TriggerPipeline", err, "failed to send request")
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, p.wrapHTTPError("TriggerPipeline", resp.StatusCode, string(body))
	}

	// 返回一个占位的 Pipeline 对象，表示触发成功
	// Gitea dispatch API 不返回 run ID，需要后续查询
	return &onlinegit.Pipeline{
		Ref:    opts.Ref,
		Status: onlinegit.PipelineStatusPending,
		Source: onlinegit.PipelineSourceWeb,
	}, nil
}

// GetPipeline 获取 Gitea Actions Workflow Run 详情
func (p *Provider) GetPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	run, resp, err := p.client.GetRepoActionRun(p.owner, p.repo, pipelineID)
	if err != nil {
		return nil, p.wrapError("GetPipeline", resp, err)
	}
	return p.toActionWorkflowRunPipeline(run), nil
}

// ListPipelines 获取 Gitea Actions Workflow Run 列表
func (p *Provider) ListPipelines(ctx context.Context, opts *onlinegit.ListPipelineOptions) ([]*onlinegit.Pipeline, error) {
	listOpts := gitea.ListRepoActionRunsOptions{
		ListOptions: gitea.ListOptions{
			Page:     opts.Page,
			PageSize: opts.PerPage,
		},
	}

	runsResp, resp, err := p.client.ListRepoActionRuns(p.owner, p.repo, listOpts)
	if err != nil {
		return nil, p.wrapError("ListPipelines", resp, err)
	}

	result := make([]*onlinegit.Pipeline, len(runsResp.WorkflowRuns))
	for i, run := range runsResp.WorkflowRuns {
		result[i] = p.toActionWorkflowRunPipeline(run)
	}
	return result, nil
}

// CancelPipeline 取消 Gitea Actions Workflow Run
// 注意：Gitea SDK 目前不直接支持取消 workflow run
func (p *Provider) CancelPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "CancelPipeline", onlinegit.ErrNotSupported, "Gitea SDK does not support cancel workflow run API yet")
}

// RetryPipeline 重试 Gitea Actions Workflow Run
// 注意：Gitea SDK 目前不直接支持重试 workflow run
func (p *Provider) RetryPipeline(ctx context.Context, pipelineID int64) (*onlinegit.Pipeline, error) {
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitea, "RetryPipeline", onlinegit.ErrNotSupported, "Gitea SDK does not support retry workflow run API yet")
}

// ListPipelineJobs 获取 Gitea Actions Workflow Run 的作业列表
func (p *Provider) ListPipelineJobs(ctx context.Context, pipelineID int64) ([]*onlinegit.PipelineJob, error) {
	listOpts := gitea.ListRepoActionJobsOptions{
		ListOptions: gitea.ListOptions{
			Page:     1,
			PageSize: 100,
		},
	}

	jobsResp, resp, err := p.client.ListRepoActionRunJobs(p.owner, p.repo, pipelineID, listOpts)
	if err != nil {
		return nil, p.wrapError("ListPipelineJobs", resp, err)
	}

	result := make([]*onlinegit.PipelineJob, len(jobsResp.Jobs))
	for i, job := range jobsResp.Jobs {
		result[i] = p.toActionWorkflowJob(job)
	}
	return result, nil
}

// toActionWorkflowRunPipeline 将 Gitea ActionWorkflowRun 转换为统一的 Pipeline 结构
func (p *Provider) toActionWorkflowRunPipeline(run *gitea.ActionWorkflowRun) *onlinegit.Pipeline {
	pipeline := &onlinegit.Pipeline{
		ID:     run.ID,
		IID:    run.RunNumber,
		Ref:    run.HeadBranch,
		SHA:    run.HeadSha,
		WebURL: run.HTMLURL,
		Status: p.mapGiteaActionStatus(run.Status, run.Conclusion),
		Source: p.mapGiteaActionEvent(run.Event),
	}

	if !run.StartedAt.IsZero() {
		startedAt := run.StartedAt
		pipeline.StartedAt = &startedAt
		pipeline.CreatedAt = run.StartedAt // 使用 StartedAt 作为 CreatedAt
	}

	if !run.CompletedAt.IsZero() {
		completedAt := run.CompletedAt
		pipeline.FinishedAt = &completedAt
		pipeline.UpdatedAt = run.CompletedAt // 使用 CompletedAt 作为 UpdatedAt
	}

	// 计算持续时间
	if pipeline.StartedAt != nil && pipeline.FinishedAt != nil {
		pipeline.Duration = int64(pipeline.FinishedAt.Sub(*pipeline.StartedAt).Seconds())
	}

	if run.Actor != nil {
		pipeline.User = &onlinegit.User{
			ID:        run.Actor.ID,
			Login:     run.Actor.UserName,
			Name:      run.Actor.FullName,
			Email:     run.Actor.Email,
			AvatarURL: run.Actor.AvatarURL,
		}
	}

	return pipeline
}

// mapGiteaActionStatus 映射 Gitea Action 状态到统一状态
func (p *Provider) mapGiteaActionStatus(status, conclusion string) onlinegit.PipelineStatus {
	switch status {
	case "queued", "waiting":
		return onlinegit.PipelineStatusPending
	case "in_progress", "running":
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

// mapGiteaActionEvent 映射 Gitea Action 事件到 Pipeline Source
func (p *Provider) mapGiteaActionEvent(event string) onlinegit.PipelineSource {
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

// toActionWorkflowJob 将 Gitea ActionWorkflowJob 转换为统一的 PipelineJob 结构
func (p *Provider) toActionWorkflowJob(job *gitea.ActionWorkflowJob) *onlinegit.PipelineJob {
	pipelineJob := &onlinegit.PipelineJob{
		ID:     job.ID,
		Name:   job.Name,
		Status: p.mapGiteaJobStatus(job.Status, job.Conclusion),
		WebURL: job.HTMLURL,
	}

	if !job.StartedAt.IsZero() {
		startedAt := job.StartedAt
		pipelineJob.StartedAt = &startedAt
	}
	if !job.CompletedAt.IsZero() {
		completedAt := job.CompletedAt
		pipelineJob.FinishedAt = &completedAt
		if !job.StartedAt.IsZero() {
			pipelineJob.Duration = job.CompletedAt.Sub(job.StartedAt).Seconds()
		}
	}

	return pipelineJob
}

// mapGiteaJobStatus 映射 Gitea Job 状态
func (p *Provider) mapGiteaJobStatus(status, conclusion string) onlinegit.PipelineStatus {
	switch status {
	case "queued", "waiting":
		return onlinegit.PipelineStatusPending
	case "in_progress", "running":
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
