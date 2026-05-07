package gitea

import (
	"code.gitea.io/sdk/gitea"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

// toPullRequest 转换 Pull Request
func (p *Provider) toPullRequest(pr *gitea.PullRequest) *onlinegit.PullRequest {
	state := onlinegit.PRStateOpen
	if pr.HasMerged {
		state = onlinegit.PRStateMerged
	} else if pr.State == gitea.StateClosed {
		state = onlinegit.PRStateClosed
	}

	result := &onlinegit.PullRequest{
		ID:           pr.ID,
		Number:       int(pr.Index),
		Title:        pr.Title,
		Body:         pr.Body,
		State:        state,
		SourceBranch: pr.Head.Ref,
		TargetBranch: pr.Base.Ref,
		URL:          pr.HTMLURL,
		Merged:       pr.HasMerged,
		CreatedAt:    *pr.Created,
		UpdatedAt:    *pr.Updated,
	}

	if pr.Poster != nil {
		result.Author = p.toUser(pr.Poster)
	}

	if pr.Merged != nil {
		result.MergedAt = *pr.Merged
	}

	if pr.Closed != nil {
		result.ClosedAt = *pr.Closed
	}

	if pr.MergedBy != nil {
		result.MergedBy = p.toUser(pr.MergedBy)
	}

	if len(pr.Labels) > 0 {
		result.Labels = make([]string, len(pr.Labels))
		for i, l := range pr.Labels {
			result.Labels[i] = l.Name
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
func (p *Provider) toUser(u *gitea.User) *onlinegit.User {
	return &onlinegit.User{
		ID:        u.ID,
		Login:     u.UserName,
		Name:      u.FullName,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
	}
}

// toBranchCommit 转换分支提交信息
func (p *Provider) toBranchCommit(c *gitea.PayloadCommit) *onlinegit.Commit {
	if c == nil {
		return nil
	}
	result := &onlinegit.Commit{
		SHA:     c.ID,
		Message: c.Message,
		URL:     c.URL,
	}

	if c.Author != nil {
		result.Author = &onlinegit.User{
			Name:  c.Author.Name,
			Email: c.Author.Email,
		}
	}

	if c.Committer != nil {
		result.Committer = &onlinegit.User{
			Name:  c.Committer.Name,
			Email: c.Committer.Email,
		}
	}

	return result
}

// toRepoCommit 转换仓库提交信息
func (p *Provider) toRepoCommit(c *gitea.Commit) *onlinegit.Commit {
	if c == nil {
		return nil
	}

	result := &onlinegit.Commit{
		SHA:     c.SHA,
		Message: c.RepoCommit.Message,
		URL:     c.HTMLURL,
	}

	if c.Author != nil {
		result.Author = &onlinegit.User{
			ID:        c.Author.ID,
			Login:     c.Author.UserName,
			Name:      c.Author.FullName,
			Email:     c.Author.Email,
			AvatarURL: c.Author.AvatarURL,
		}
	}

	if c.Committer != nil {
		result.Committer = &onlinegit.User{
			ID:        c.Committer.ID,
			Login:     c.Committer.UserName,
			Name:      c.Committer.FullName,
			Email:     c.Committer.Email,
			AvatarURL: c.Committer.AvatarURL,
		}
	}

	if c.CommitMeta != nil {
		result.CreatedAt = c.CommitMeta.Created
	}

	return result
}

// toComment 转换评论信息
func (p *Provider) toComment(c *gitea.Comment) *onlinegit.Comment {
	result := &onlinegit.Comment{
		ID:        c.ID,
		Body:      c.Body,
		URL:       c.HTMLURL,
		CreatedAt: c.Created,
		UpdatedAt: c.Updated,
	}

	if c.Poster != nil {
		result.Author = p.toUser(c.Poster)
	}

	return result
}
