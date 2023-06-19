package git

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/yi-nology/common/utils/xlogger"
	"os"
)

func (i *Info) Branches(log xlogger.Logger) ([]*plumbing.Reference, error) {
	err := i.project.Fetch(&git.FetchOptions{Auth: i.publicKey})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, err
	}

	branches, err := i.project.Branches()
	if err != nil {
		log.Errorf("get branches is err %+v", err)
		return nil, err
	}
	var branchesResult []*plumbing.Reference
	branches.ForEach(func(ref *plumbing.Reference) error {
		fmt.Println(ref.Name())
		return nil
	})

	return branchesResult, nil
}

func (i *Info) Tags(log xlogger.Logger) ([]*plumbing.Reference, error) {
	tags, err := i.project.Tags()
	if err != nil {
		log.Errorf("get tags is err %+v", err)
		return nil, err
	}
	var tagsResult []*plumbing.Reference
	if err := tags.ForEach(func(ref *plumbing.Reference) error {
		if ref != nil {
			tagsResult = append(tagsResult, ref)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return tagsResult, nil
}

func (i *Info) CreateTag(log xlogger.Logger, ctx context.Context, pl plumbing.Reference, version string, sign *object.Signature) error {
	_, err := i.project.CreateTag(version, pl.Hash(), &git.CreateTagOptions{
		Tagger:  sign, // 这里可以填入你的信息
		Message: "this is my tag",
	})
	if err != nil {
		log.Errorf("CreateTag err!=nil err:%+v", err)
		return err
	}

	return i.push(ctx, log)
}

func (i *Info) DeleteTag(log xlogger.Logger, ctx context.Context, version string) error {
	err := i.project.DeleteTag(version)
	if err != nil {
		log.Errorf("DeleteTag %+v", err)
		return err
	}

	return i.push(ctx, log)
}

func (i *Info) push(ctx context.Context, log xlogger.Logger) error {
	if err := i.project.PushContext(ctx, &git.PushOptions{
		RemoteName: i.url,
		Progress:   os.Stdout,
		RefSpecs: []config.RefSpec{
			"+refs/tags/*:refs/tags/*",
		},
		Auth: i.publicKey,
	}); err != nil {
		log.Errorf("push err:%+v", err)
		return err
	}
	return nil
}
