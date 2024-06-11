package main

import (
	"flag"
	"fmt"

	"github.com/delyr1c/dechoric/src/domain/domainRpc/domainRpc"
	"github.com/delyr1c/dechoric/src/domain/domainRpc/internal/config"
	"github.com/delyr1c/dechoric/src/domain/domainRpc/internal/server"
	"github.com/delyr1c/dechoric/src/domain/domainRpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/domainrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		domainRpc.RegisterDomainRpcServer(grpcServer, server.NewDomainRpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
