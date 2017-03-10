package distproto

import (
	contextpkg "context"

	//"github.com/docker/containerd/cmd/dist/distproto/registrydocker"
	"github.com/docker/containerd/cmd/dist/distproto/sshgit"
	"github.com/docker/containerd/remotes"
)

// GetResolver prepares the resolver from the environment and options.
func GetResolver(ctx contextpkg.Context) (remotes.Resolver, error) {
	return sshgit.NewGITResolver(), nil
	//return registrydocker.NewDockerResolver(), nil
}
