package svc

import (
	"github.com/delyr1c/dechoric/src/domain/domainApi/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
