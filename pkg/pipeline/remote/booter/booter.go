package booter

import (
	"github.com/pkg/errors"
	"github.com/rancher/rancher/pkg/pipeline/remote"
	"github.com/rancher/rancher/pkg/pipeline/remote/github"
	"github.com/rancher/types/apis/management.cattle.io/v3"
)

func New(pipeline v3.ClusterPipeline, remoteType string) (remote.Remote, error) {
	switch remoteType {
	case "github":
		return github.New(pipeline)
	}
	return nil, errors.New("unsupported remote type")
}
