package git

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/yi-nology/common/utils/xlogger"
)

//帮我拥有如下功能的go程序
//1.获取全部分支
//2.获取全部标签
//3.将特定分支特定版本号进行打tag
//4.删除指定tag

type Git interface {
	Branches(log xlogger.Logger) ([]*plumbing.Reference, error)
	Tags(log xlogger.Logger) ([]*plumbing.Reference, error)
	CreateTag(log xlogger.Logger, ctx context.Context, pl plumbing.Reference, version string, sign *object.Signature) error
	DeleteTag(log xlogger.Logger, ctx context.Context, version string) error
}
