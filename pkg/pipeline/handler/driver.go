package handler

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/pipeline/handler/drivers"
	"github.com/rancher/types/config"
)

var Drivers map[string]Driver

type Driver interface {
	Execute(apiContext types.APIContext) (int, error)
}

func RegisterDrivers(Management *config.ManagementContext) {
	Drivers = map[string]Driver{}
	Drivers[drivers.GITHUB_WEBHOOK_HEADER] = drivers.GithubDriver{Management: Management}
}
