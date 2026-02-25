package gitlab

import (
	"context"
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	onlinegit "github.com/yi-nology/common/biz/online-git"
)

func (p *Provider) CompareBranches(ctx context.Context, base, head string) (*onlinegit.CompareResult, error) {
	opts := &gitlab.CompareOptions{
		From: gitlab.Ptr(base),
		To:   gitlab.Ptr(head),
	}

	compare, resp, err := p.client.Repositories.Compare(p.projectID, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("CompareBranches", resp, err)
	}

	commits := make([]*onlinegit.Commit, len(compare.Commits))
	for i, c := range compare.Commits {
		commits[i] = &onlinegit.Commit{
			SHA:     c.ID,
			Message: c.Message,
			Author: &onlinegit.User{
				Name:  c.AuthorName,
				Email: c.AuthorEmail,
			},
			CreatedAt: *c.CreatedAt,
		}
	}

	files := make([]*onlinegit.FileChange, len(compare.Diffs))
	for i, d := range compare.Diffs {
		status := onlinegit.FileChangeModified
		if d.NewFile {
			status = onlinegit.FileChangeAdded
		} else if d.DeletedFile {
			status = onlinegit.FileChangeDeleted
		} else if d.RenamedFile {
			status = onlinegit.FileChangeRenamed
		}

		files[i] = &onlinegit.FileChange{
			Filename:     d.NewPath,
			Status:       status,
			PreviousName: d.OldPath,
			Patch:        d.Diff,
		}
	}

	return &onlinegit.CompareResult{
		BaseBranch:   base,
		HeadBranch:   head,
		TotalCommits: len(commits),
		Commits:      commits,
		Files:        files,
		DiffStat: &onlinegit.DiffStat{
			ChangedFiles: len(files),
		},
	}, nil
}

func (p *Provider) ListComments(ctx context.Context, prNumber int) ([]*onlinegit.Comment, error) {
	notes, resp, err := p.client.Notes.ListMergeRequestNotes(p.projectID, int64(prNumber), nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("ListComments", resp, err)
	}

	result := make([]*onlinegit.Comment, len(notes))
	for i, n := range notes {
		result[i] = p.toNote(n)
	}
	return result, nil
}

func (p *Provider) CreateComment(ctx context.Context, prNumber int, body string) (*onlinegit.Comment, error) {
	opts := &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.Ptr(body),
	}

	note, resp, err := p.client.Notes.CreateMergeRequestNote(p.projectID, int64(prNumber), opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("CreateComment", resp, err)
	}
	return p.toNote(note), nil
}

func (p *Provider) UpdateComment(ctx context.Context, commentID int64, body string) (*onlinegit.Comment, error) {
	// GitLab Notes API 需要 MR 号，这里需要额外处理
	return nil, onlinegit.NewProviderError(onlinegit.PlatformGitLab, "UpdateComment", fmt.Errorf("not implemented: need MR number"), "")
}

func (p *Provider) DeleteComment(ctx context.Context, commentID int64) error {
	// GitLab Notes API 需要 MR 号
	return onlinegit.NewProviderError(onlinegit.PlatformGitLab, "DeleteComment", fmt.Errorf("not implemented: need MR number"), "")
}

func (p *Provider) GetCommit(ctx context.Context, sha string) (*onlinegit.Commit, error) {
	commit, resp, err := p.client.Commits.GetCommit(p.projectID, sha, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("GetCommit", resp, err)
	}

	return &onlinegit.Commit{
		SHA:     commit.ID,
		Message: commit.Message,
		Author: &onlinegit.User{
			Name:  commit.AuthorName,
			Email: commit.AuthorEmail,
		},
		Committer: &onlinegit.User{
			Name:  commit.CommitterName,
			Email: commit.CommitterEmail,
		},
		URL:       commit.WebURL,
		CreatedAt: *commit.CreatedAt,
	}, nil
}

func (p *Provider) ListCommits(ctx context.Context, branch string, opts *onlinegit.ListOptions) ([]*onlinegit.Commit, error) {
	glOpts := &gitlab.ListCommitsOptions{
		RefName: gitlab.Ptr(branch),
		ListOptions: gitlab.ListOptions{
			Page:    int64(opts.Page),
			PerPage: int64(opts.PerPage),
		},
	}

	commits, resp, err := p.client.Commits.ListCommits(p.projectID, glOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, p.wrapError("ListCommits", resp, err)
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
			URL:       c.WebURL,
			CreatedAt: *c.CreatedAt,
		}
	}
	return result, nil
}
