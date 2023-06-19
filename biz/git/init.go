package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
)

type Info struct {
	project   *git.Repository
	url       string
	publicKey *ssh.PublicKeys
}

func NewMemory(url, sshKeyStr string, password ...string) *Info {
	pw := ""
	if len(password) != 0 {
		pw = password[0]
	}
	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKeyStr), pw)
	if err != nil {
		return nil
	}
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		Progress: os.Stdout,
		URL:      url,
		Auth:     publicKey,
	})
	if err != nil {
		return nil
	}
	return &Info{project: r, url: url, publicKey: publicKey}
}
