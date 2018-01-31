package remote

import (
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"net/http"
)

type Remote interface {
	Type() string

	CanLogin() bool

	CanRepos() bool

	CanHook() bool

	//Login handle oauth login
	Login(redirectURL string, code string) (*v3.RemoteAccount, error)

	Repos(account *v3.RemoteAccount) ([]*v3.GitRepository, error)

	CreateHook()

	DeleteHook()

	ParseHook(r *http.Request)
}
