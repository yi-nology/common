package gitlab

import (
	"context"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

func (p *Provider) ListPullRequests(ctx context.Context, state onlinegit.PRState, opts *onlinegit.ListOptions) ([]*onlinegit.PullRequest, error) {
	glState := ""
	switch state {
	case onlinegit.PRStateOpen:
		glState = "opened"
	case onlinegit.PRStateClosed:
		glState = "closed"
	case onlinegit.PRStateMerged:
		glState = "merged"
	case onlinegit.PRStateAll:
		glState = "all"
	}

	glOpts := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.Ptr(glState),
		ListOptions: gitlab.ListOptions{
			Page:    int64(opts.Page),
			PerPage: int64(opts.PerPage),
		},
	}

	mrs, resp, err := p.client.MergeRequests.ListProjectMergeRequests(p.projectID, glOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("ListPullRequests", resp, err)
	}

	result := make([]*onlinegit.PullRequest, len(mrs))
	for i, mr := range mrs {
		result[i] = p.toBasicMergeRequest(mr)
	}
	return result, nil
}

func (p *Provider) GetPullRequest(ctx context.Context, number int) (*onlinegit.PullRequest, error) {
	mr, resp, err := p.client.MergeRequests.GetMergeRequest(p.projectID, int64(number), nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("GetPullRequest", resp, err)
	}
	return p.toMergeRequest(mr), nil
}

func (p *Provider) CreatePullRequest(ctx context.Context, req *onlinegit.CreatePRRequest) (*onlinegit.PullRequest, error) {
	opts := &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.Ptr(req.Title),
		Description:  gitlab.Ptr(req.Body),
		SourceBranch: gitlab.Ptr(req.SourceBranch),
		TargetBranch: gitlab.Ptr(req.TargetBranch),
	}

	if len(req.Labels) > 0 {
		opts.Labels = gitlab.Ptr(gitlab.LabelOptions(req.Labels))
	}

	mr, resp, err := p.client.MergeRequests.CreateMergeRequest(p.projectID, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("CreatePullRequest", resp, err)
	}

	return p.toMergeRequest(mr), nil
}

func (p *Provider) UpdatePullRequest(ctx context.Context, number int, title, body string) (*onlinegit.PullRequest, error) {
	opts := &gitlab.UpdateMergeRequestOptions{
		Title:       gitlab.Ptr(title),
		Description: gitlab.Ptr(body),
	}

	mr, resp, err := p.client.MergeRequests.UpdateMergeRequest(p.projectID, int64(number), opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("UpdatePullRequest", resp, err)
	}
	return p.toMergeRequest(mr), nil
}

func (p *Provider) MergePullRequest(ctx context.Context, number int, opts *onlinegit.MergeOptions) error {
	acceptOpts := &gitlab.AcceptMergeRequestOptions{}

	if opts != nil {
		if opts.CommitMessage != "" {
			acceptOpts.MergeCommitMessage = gitlab.Ptr(opts.CommitMessage)
		}
		if opts.Method == onlinegit.MergeMethodSquash {
			acceptOpts.Squash = gitlab.Ptr(true)
		}
		if opts.DeleteBranch {
			acceptOpts.ShouldRemoveSourceBranch = gitlab.Ptr(true)
		}
	}

	_, resp, err := p.client.MergeRequests.AcceptMergeRequest(p.projectID, int64(number), acceptOpts, gitlab.WithContext(ctx))
	if err != nil {
		return p.wrapError("MergePullRequest", resp, err)
	}
	return nil
}

func (p *Provider) ClosePullRequest(ctx context.Context, number int) error {
	opts := &gitlab.UpdateMergeRequestOptions{
		StateEvent: gitlab.Ptr("close"),
	}

	_, resp, err := p.client.MergeRequests.UpdateMergeRequest(p.projectID, int64(number), opts, gitlab.WithContext(ctx))
	if err != nil {
		return p.wrapError("ClosePullRequest", resp, err)
	}
	return nil
}

func (p *Provider) GetPullRequestCommits(ctx context.Context, number int) ([]*onlinegit.Commit, error) {
	commits, resp, err := p.client.MergeRequests.GetMergeRequestCommits(p.projectID, int64(number), nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("GetPullRequestCommits", resp, err)
	}

	result := make([]*onlinegit.Commit, len(commits))
	for i, c := range commits {
		result[i] = &onlinegit.Commit{
			SHA:     c.ID,
			Message: c.Message,
			Author: &onlinegit.User{
				Name:  c.AuthorName,
				Email: c.AuthorEmail,
			},
			CreatedAt: *c.CreatedAt,
		}
	}
	return result, nil
}

