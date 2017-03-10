package sshgit

import (
	contextpkg "context"
	"net/url"
	"path"
	"strings"

	"github.com/docker/containerd/reference"
	"github.com/docker/containerd/remotes"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

type sshgitResolver struct{}

var _ remotes.Resolver = &sshgitResolver{}

//NewGITResolver funtion
func NewGITResolver() remotes.Resolver {
	return &sshgitResolver{}
}

func (r *sshgitResolver) Resolve(ctx contextpkg.Context, ref string) (string, ocispec.Descriptor, remotes.Fetcher, error) {
	refspec, err := reference.Parse(ref)
	if err != nil {
		return "", ocispec.Descriptor{}, nil, err
	}

	var (
		base url.URL
	)

	switch refspec.Hostname() {
	case "git":
		base.Scheme = "ssh"
		base.Host = "localhost:22"
		prefix := strings.TrimPrefix(refspec.Locator, "git/")
		base.Path = path.Join("/git-server/repos", prefix)
	default:
		return "", ocispec.Descriptor{}, nil, errors.Errorf("unsupported locator: %q", refspec.Locator)
	}

	return "", ocispec.Descriptor{}, nil, errors.Errorf("From Git: %v %v %v, base : %v", ref, refspec.Locator, refspec.Object, base)
}
