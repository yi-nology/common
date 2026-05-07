package gitlab

import (
	gitlab "gitlab.com/gitlab-org/api/client-go"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

func (p *Provider) toMergeRequest(mr *gitlab.MergeRequest) *onlinegit.PullRequest {
	state := onlinegit.PRStateOpen
	switch mr.State {
	case "merged":
		state = onlinegit.PRStateMerged
	case "closed":
		state = onlinegit.PRStateClosed
	}

	result := &onlinegit.PullRequest{
		ID:           mr.IID,
		Number:       int(mr.IID),
		Title:        mr.Title,
		Body:         mr.Description,
		State:        state,
		SourceBranch: mr.SourceBranch,
		TargetBranch: mr.TargetBranch,
		URL:          mr.WebURL,
		Merged:       mr.State == "merged",
		CreatedAt:    *mr.CreatedAt,
		UpdatedAt:    *mr.UpdatedAt,
	}

	if mr.Author != nil {
		result.Author = &onlinegit.User{
			ID:        int64(mr.Author.ID),
			Login:     mr.Author.Username,
			Name:      mr.Author.Name,
			AvatarURL: mr.Author.AvatarURL,
		}
	}

	if mr.MergedAt != nil {
		result.MergedAt = *mr.MergedAt
	}

	if mr.ClosedAt != nil {
		result.ClosedAt = *mr.ClosedAt
	}

	if mr.MergedBy != nil {
		result.MergedBy = &onlinegit.User{
			ID:        int64(mr.MergedBy.ID),
			Login:     mr.MergedBy.Username,
			Name:      mr.MergedBy.Name,
			AvatarURL: mr.MergedBy.AvatarURL,
		}
	}

	if len(mr.Labels) > 0 {
		result.Labels = mr.Labels
	}

	return result
}

func (p *Provider) toBasicMergeRequest(mr *gitlab.BasicMergeRequest) *onlinegit.PullRequest {
	state := onlinegit.PRStateOpen
	switch mr.State {
	case "merged":
		state = onlinegit.PRStateMerged
	case "closed":
		state = onlinegit.PRStateClosed
	}

	result := &onlinegit.PullRequest{
		ID:           mr.IID,
		Number:       int(mr.IID),
		Title:        mr.Title,
		Body:         mr.Description,
		State:        state,
		SourceBranch: mr.SourceBranch,
		TargetBranch: mr.TargetBranch,
		URL:          mr.WebURL,
		Merged:       mr.State == "merged",
		CreatedAt:    *mr.CreatedAt,
		UpdatedAt:    *mr.UpdatedAt,
	}

	if mr.Author != nil {
		result.Author = &onlinegit.User{
			ID:        int64(mr.Author.ID),
			Login:     mr.Author.Username,
			Name:      mr.Author.Name,
			AvatarURL: mr.Author.AvatarURL,
		}
	}

	if mr.MergedAt != nil {
		result.MergedAt = *mr.MergedAt
	}

	if mr.ClosedAt != nil {
		result.ClosedAt = *mr.ClosedAt
	}

	if mr.MergedBy != nil {
		result.MergedBy = &onlinegit.User{
			ID:        int64(mr.MergedBy.ID),
			Login:     mr.MergedBy.Username,
			Name:      mr.MergedBy.Name,
			AvatarURL: mr.MergedBy.AvatarURL,
		}
	}

	if len(mr.Labels) > 0 {
		result.Labels = mr.Labels
	}

	return result
}

func (p *Provider) toCommitFromBranch(c *gitlab.Commit) *onlinegit.Commit {
	if c == nil {
		return nil
	}
	return &onlinegit.Commit{
		SHA:     c.ID,
		Message: c.Message,
		Author: &onlinegit.User{
			Name:  c.AuthorName,
			Email: c.AuthorEmail,
		},
		URL:       c.WebURL,
		CreatedAt: *c.CreatedAt,
	}
}

func (p *Provider) toNote(n *gitlab.Note) *onlinegit.Comment {
	result := &onlinegit.Comment{
		ID:        int64(n.ID),
		Body:      n.Body,
		CreatedAt: *n.CreatedAt,
		UpdatedAt: *n.UpdatedAt,
	}

	// Author 是值类型，检查 ID 是否非零
	if n.Author.ID != 0 {
		result.Author = &onlinegit.User{
			ID:        int64(n.Author.ID),
			Login:     n.Author.Username,
			Name:      n.Author.Name,
			AvatarURL: n.Author.AvatarURL,
		}
	}

	return result
}
