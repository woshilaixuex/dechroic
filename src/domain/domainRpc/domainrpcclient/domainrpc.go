// Code generated by goctl. DO NOT EDIT.
// Source: domainRpc.proto

package domainrpcclient

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/domainRpc/domainRpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Request  = domainRpc.Request
	Response = domainRpc.Response

	DomainRpc interface {
		Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	}

	defaultDomainRpc struct {
		cli zrpc.Client
	}
)

func NewDomainRpc(cli zrpc.Client) DomainRpc {
	return &defaultDomainRpc{
		cli: cli,
	}
}

func (m *defaultDomainRpc) Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	client := domainRpc.NewDomainRpcClient(m.cli.Conn())
	return client.Ping(ctx, in, opts...)
}