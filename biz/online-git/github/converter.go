package github

import (
	"github.com/google/go-github/v56/github"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// toPullRequest 转换 Pull Request
func (p *Provider) toPullRequest(pr *github.PullRequest) *onlinegit.PullRequest {
	state := onlinegit.PRStateOpen
	if pr.GetMerged() {
		state = onlinegit.PRStateMerged
	} else if pr.GetState() == "closed" {
		state = onlinegit.PRStateClosed
	}

	result := &onlinegit.PullRequest{
		ID:           pr.GetID(),
		Number:       pr.GetNumber(),
		Title:        pr.GetTitle(),
		Body:         pr.GetBody(),
		State:        state,
		SourceBranch: pr.GetHead().GetRef(),
		TargetBranch: pr.GetBase().GetRef(),
		URL:          pr.GetHTMLURL(),
		Merged:       pr.GetMerged(),
		CreatedAt:    pr.GetCreatedAt().Time,
		UpdatedAt:    pr.GetUpdatedAt().Time,
	}

	if pr.GetUser() != nil {
		result.Author = p.toUser(pr.GetUser())
	}

	if pr.MergedAt != nil {
		result.MergedAt = pr.GetMergedAt().Time
	}

	if pr.ClosedAt != nil {
		result.ClosedAt = pr.GetClosedAt().Time
	}

	if pr.MergedBy != nil {
		result.MergedBy = p.toUser(pr.MergedBy)
	}

	if len(pr.Labels) > 0 {
		result.Labels = make([]string, len(pr.Labels))
		for i, l := range pr.Labels {
			result.Labels[i] = l.GetName()
		}
	}

	if len(pr.Assignees) > 0 {
		result.Assignees = make([]*onlinegit.User, len(pr.Assignees))
		for i, a := range pr.Assignees {
			result.Assignees[i] = p.toUser(a)
		}
	}

	return result
}

// toUser 转换用户信息
func (p *Provider) toUser(u *github.User) *onlinegit.User {
	return &onlinegit.User{
		ID:        u.GetID(),
		Login:     u.GetLogin(),
		Name:      u.GetName(),
		Email:     u.GetEmail(),
		AvatarURL: u.GetAvatarURL(),
	}
}

// toCommit 转换提交信息
func (p *Provider) toCommit(c *github.Commit, sha string) *onlinegit.Commit {
	if c == nil {
		return nil
	}

	result := &onlinegit.Commit{
		SHA:     sha,
		Message: c.GetMessage(),
		URL:     c.GetURL(),
	}

	if c.Author != nil {
		result.Author = &onlinegit.User{
			Name:  c.Author.GetName(),
			Email: c.Author.GetEmail(),
		}
		if c.Author.Date != nil {
			result.CreatedAt = c.Author.Date.Time
		}
	}

	if c.Committer != nil {
		result.Committer = &onlinegit.User{
			Name:  c.Committer.GetName(),
			Email: c.Committer.GetEmail(),
		}
	}

	return result
}

// toComment 转换评论信息
func (p *Provider) toComment(c *github.IssueComment) *onlinegit.Comment {
	result := &onlinegit.Comment{
		ID:        c.GetID(),
		Body:      c.GetBody(),
		URL:       c.GetHTMLURL(),
		CreatedAt: c.GetCreatedAt().Time,
		UpdatedAt: c.GetUpdatedAt().Time,
	}

	if c.User != nil {
		result.Author = p.toUser(c.User)
	}

	return result
}
