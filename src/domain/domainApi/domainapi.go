package main

import (
	"flag"
	"fmt"

	"github.com/delyr1c/dechoric/src/domain/domainApi/internal/config"
	"github.com/delyr1c/dechoric/src/domain/domainApi/internal/handler"
	"github.com/delyr1c/dechoric/src/domain/domainApi/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/domainapi-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}